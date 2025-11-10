package redis

import (
	"crypto/tls"
	"io"

	gcppropagator "github.com/GoogleCloudPlatform/opentelemetry-operations-go/propagator"
	"github.com/quincy0/go-kit/core/trace"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"

	"github.com/go-redis/redis/extra/redisotel"
	red "github.com/go-redis/redis/v8"
	"github.com/quincy0/go-kit/core/syncx"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

const (
	defaultDatabase = 0
	maxRetries      = 3
	idleConns       = 8
)

var clientManager = syncx.NewResourceManager()

func getClient(r *Redis) (*red.Client, error) {
	if len(r.ProjectId) > 0 {
		_, _ = tracerProvider(r.ProjectId)
	}

	val, err := clientManager.GetResource(r.Addr, func() (io.Closer, error) {
		var tlsConfig *tls.Config
		if r.tls {
			tlsConfig = &tls.Config{
				InsecureSkipVerify: true,
			}
		}
		store := red.NewClient(&red.Options{
			Addr:         r.Addr,
			Password:     r.Pass,
			DB:           r.DB,
			MaxRetries:   maxRetries,
			MinIdleConns: idleConns,
			TLSConfig:    tlsConfig,
		})
		store.AddHook(durationHook)

		if len(r.ProjectId) > 0 {
			store.AddHook(redisotel.TracingHook{})
		} else {
			store.AddHook(durationHook)
		}

		return store, nil
	})
	if err != nil {
		return nil, err
	}

	return val.(*red.Client), nil
}

func tracerProvider(projectId string) (*sdktrace.TracerProvider, error) {
	tp, _ := trace.GetCloudTraceProvider(projectId)
	if tp != nil {
		otel.SetTracerProvider(tp)
		otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
			// Putting the CloudTraceOneWayPropagator first means the TraceContext propagator
			// takes precedence if both the traceparent and the XCTC headers exist.
			gcppropagator.CloudTraceOneWayPropagator{},
			propagation.TraceContext{},
			propagation.Baggage{},
		))
	}
	return tp, nil
}

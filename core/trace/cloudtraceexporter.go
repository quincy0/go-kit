package trace

import (
	"sync"

	cloudtrace "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace"
	"go-kit/core/logx"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

var (
	googleCloudTraceProvider *sdktrace.TracerProvider
	googleCloudTraceExporter *cloudtrace.Exporter
	once                     sync.Once
)

func GetCloudTraceProvider(projectId string) (*sdktrace.TracerProvider, *cloudtrace.Exporter) {
	if len(projectId) == 0 {
		return nil, nil
	}
	var err error
	once.Do(func() {
		googleCloudTraceExporter, err = cloudtrace.New(cloudtrace.WithProjectID(projectId))
		if err != nil {
			logx.Errorw("new cloud trace exporter error", logx.Field("err", err))
		} else {
			googleCloudTraceProvider = sdktrace.NewTracerProvider(
				sdktrace.WithSampler(sdktrace.AlwaysSample()),
				sdktrace.WithBatcher(googleCloudTraceExporter),
			)
		}
	})

	return googleCloudTraceProvider, googleCloudTraceExporter
}

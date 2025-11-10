package sqlx

import (
	gcppropagator "github.com/GoogleCloudPlatform/opentelemetry-operations-go/propagator"
	"github.com/go-sql-driver/mysql"
	"github.com/quincy0/go-kit/core/logx"
	"github.com/quincy0/go-kit/core/trace"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
)

const (
	mysqlDriverName           = "mysql"
	duplicateEntryCode uint16 = 1062
)

var ProjectId = ""

// NewMysql returns a mysql connection.
func NewMysql(datasource string, opts ...SqlOption) SqlConn {
	opts = append(opts, withMysqlAcceptable())
	if len(ProjectId) > 0 {
		_, _ = tracerProvider(ProjectId)
	}

	return NewSqlConn(mysqlDriverName, datasource, opts...)
}

func mysqlAcceptable(err error) bool {
	if err == nil {
		return true
	}

	myerr, ok := err.(*mysql.MySQLError)
	if !ok {
		return false
	}

	switch myerr.Number {
	case duplicateEntryCode:
		return true
	default:
		return false
	}
}

func withMysqlAcceptable() SqlOption {
	return func(conn SqlConn) {
		if c, ok := conn.(*commonSqlConn); ok {
			c.accept = mysqlAcceptable
		} else {
			logx.Errorw("Error: provided SqlConn is not of type *commonSqlConn")
		}
	}
}

func tracerProvider(projectId string) (*tracesdk.TracerProvider, error) {
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

func SetProjectId(projectId string) {
	ProjectId = projectId
}

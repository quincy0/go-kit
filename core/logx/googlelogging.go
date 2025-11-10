package logx

import (
	"cloud.google.com/go/logging"
	"context"
	"fmt"
)

type GoogleLoggingWriter struct {
	logger *logging.Logger
	client *logging.Client
}

func NewGoogleLoggingWriter(projectId string, logId string) (Writer, error) {
	client, err := logging.NewClient(context.Background(), projectId)
	if err != nil {
		return nil, err
	}

	return &GoogleLoggingWriter{
		logger: client.Logger(logId),
		client: client,
	}, nil
}

func (g *GoogleLoggingWriter) Alert(v interface{}) {
	g.logger.Log(logging.Entry{
		Payload:  v,
		Severity: logging.Alert,
	})
}

func (g *GoogleLoggingWriter) Close() error {
	return g.client.Close()
}

func (g *GoogleLoggingWriter) Error(v interface{}, fields ...LogField) {
	g.logger.Log(logging.Entry{
		Payload:  fmt.Sprintf("%v%v", v, fields),
		Severity: logging.Error,
	})
}

func (g *GoogleLoggingWriter) Info(v interface{}, fields ...LogField) {
	g.logger.Log(logging.Entry{
		Payload:  fmt.Sprintf("%v%v", v, fields),
		Severity: logging.Info,
	})
}

func (g *GoogleLoggingWriter) Severe(v interface{}) {
	g.logger.Log(logging.Entry{
		Payload:  v,
		Severity: logging.Default,
	})
}

func (g *GoogleLoggingWriter) Slow(v interface{}, fields ...LogField) {
	g.logger.Log(logging.Entry{
		Payload:  fmt.Sprintf("%v%v", v, fields),
		Severity: logging.Notice,
	})
}

func (g *GoogleLoggingWriter) Stack(v interface{}) {
	g.logger.Log(logging.Entry{
		Payload:  v,
		Severity: logging.Info,
	})
}

func (g *GoogleLoggingWriter) Stat(v interface{}, fields ...LogField) {
	g.logger.Log(logging.Entry{
		Payload:  fmt.Sprintf("%v%v", v, fields),
		Severity: logging.Info,
	})
}

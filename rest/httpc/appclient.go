package httpc

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	gcppropagator "github.com/GoogleCloudPlatform/opentelemetry-operations-go/propagator"
	"go-kit/core/trace"
	"go-kit/rest"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"

	"go-kit/core/logx"
)

type App interface {
	Client() AppClient
}

type AppClient interface {
	WithHeader(header map[string]string) AppClient
	WithQuery(query map[string]string) AppClient
	WithData(data interface{}) AppClient
	DisabledReqLog() AppClient
	DisabledReplyLog() AppClient
	Post(uri string, reply interface{}) (AppClient, error)
	PostCtx(ctx context.Context, uri string, reply interface{}) (AppClient, error)
	Get(uri string, reply interface{}) (AppClient, error)
	GetCtx(ctx context.Context, uri string, reply interface{}) (AppClient, error)
	Put(uri string, reply interface{}) (AppClient, error)
	PutCtx(ctx context.Context, uri string, reply interface{}) (AppClient, error)
	Delete(uri string, reply interface{}) (AppClient, error)
	DeleteCtx(ctx context.Context, uri string, reply interface{}) (AppClient, error)
	HttpCode() int
}

type app struct {
	name string
	host string
	conf *rest.RestConf
}

type appClient struct {
	sync.Mutex

	app             *app
	header          map[string]string
	query           map[string]string
	data            interface{}
	projectId       string
	isWriteReqLog   bool // 是否开启请求参数日志
	isWriteReplyLog bool // 是否开启返回数据日志
	resp            *http.Response
}

func NewApp(name, host string) (App, error) {
	return &app{
		name: name,
		host: host,
	}, nil
}

func (a *app) WithConfig(c *rest.RestConf) *app {
	a.conf = c
	return a
}

func (a *app) Client() AppClient {
	return &appClient{
		app:             a,
		header:          make(map[string]string),
		isWriteReqLog:   true,
		isWriteReplyLog: true,
	}
}

func (a *appClient) WithHeader(header map[string]string) AppClient {
	a.Lock()
	a.header = header
	a.Unlock()

	return a
}

func (a *appClient) WithQuery(query map[string]string) AppClient {
	a.Lock()
	a.query = query
	a.Unlock()

	return a
}

func (a *appClient) WithData(data interface{}) AppClient {
	a.Lock()
	a.data = data
	a.Unlock()

	return a
}

func (a *appClient) DisabledReqLog() AppClient {
	a.isWriteReqLog = false
	return a
}

func (a *appClient) DisabledReplyLog() AppClient {
	a.isWriteReplyLog = false
	return a
}

func (a *appClient) Post(uri string, reply interface{}) (AppClient, error) {
	return a.PostCtx(context.Background(), uri, reply)
}

func (a *appClient) PostCtx(ctx context.Context, uri string, reply interface{}) (AppClient, error) {
	bodyBytes, err := a.doRequest(ctx, "POST", uri)
	if err != nil {
		return a, err
	}

	if len(bodyBytes) > 0 && reply != nil {
		if err = json.Unmarshal(bodyBytes, reply); err != nil {
			return a, fmt.Errorf("app client request json umnarshal failed, app:%s, uri:%s, err: %s", a.app.name, uri, err.Error())
		}
	}

	return a, nil
}

func (a *appClient) Get(uri string, reply interface{}) (AppClient, error) {
	return a.GetCtx(context.Background(), uri, reply)
}

func (a *appClient) GetCtx(ctx context.Context, uri string, reply interface{}) (AppClient, error) {
	bodyBytes, err := a.doRequest(ctx, "GET", uri)
	if err != nil {
		return a, err
	}

	if len(bodyBytes) > 0 && reply != nil {
		if err = json.Unmarshal(bodyBytes, reply); err != nil {
			return a, fmt.Errorf("app client request json umnarshal failed, app:%s, uri:%s, err: %s", a.app.name, uri, err.Error())
		}
	}

	return a, nil
}

func (a *appClient) Put(uri string, reply interface{}) (AppClient, error) {
	return a.PutCtx(context.Background(), uri, reply)
}

func (a *appClient) PutCtx(ctx context.Context, uri string, reply interface{}) (AppClient, error) {
	bodyBytes, err := a.doRequest(ctx, "PUT", uri)
	if err != nil {
		return a, err
	}

	if len(bodyBytes) > 0 && reply != nil {
		if err = json.Unmarshal(bodyBytes, reply); err != nil {
			return a, fmt.Errorf("app client request json umnarshal failed, app:%s, uri:%s, err: %s", a.app.name, uri, err.Error())
		}
	}

	return a, nil
}

func (a *appClient) Delete(uri string, reply interface{}) (AppClient, error) {
	return a.DeleteCtx(context.Background(), uri, reply)
}

func (a *appClient) DeleteCtx(ctx context.Context, uri string, reply interface{}) (AppClient, error) {
	bodyBytes, err := a.doRequest(ctx, "DELETE", uri)
	if err != nil {
		return a, err
	}

	if len(bodyBytes) > 0 && reply != nil {
		if err = json.Unmarshal(bodyBytes, reply); err != nil {
			return a, fmt.Errorf("app client request json umnarshal failed, app:%s, uri:%s, err: %s", a.app.name, uri, err.Error())
		}
	}

	return a, nil
}

func (a *appClient) HttpCode() int {
	if a.resp != nil {
		return a.resp.StatusCode
	}

	return 0
}

func (a *appClient) doRequest(ctx context.Context, method, uri string) ([]byte, error) {

	if a.app.conf != nil && a.app.conf.Telemetry.GoogleCloudTrace {
		installPropagators()
		shutdown, err := initTracer(a.app.conf.Telemetry.ProjectId)
		if err != nil {
			logx.WithContext(ctx).Errorw("test error", logx.Field("error", err))
		}
		defer shutdown()
	}

	url := a.app.host
	if len(uri) != 0 {
		url = fmt.Sprintf("%s/%s", a.app.host, uri)
	}

	if len(a.query) > 0 {
		kvRaw := make([]string, 0)
		for k, v := range a.query {
			kvRaw = append(kvRaw, fmt.Sprintf("%s=%s", k, v))
		}

		if strings.Index(url, "?") != -1 {
			url = fmt.Sprintf("%s&%s", url, strings.Join(kvRaw, "&"))
		} else {
			url = fmt.Sprintf("%s?%s", url, strings.Join(kvRaw, "&"))
		}
	}

	if a.isWriteReqLog {
		logx.WithContext(ctx).Infow("app-request-pre", logx.Field("app-name", a.app.name),
			logx.Field("url", url),
			logx.Field("header", a.header),
			logx.Field("data", a.data),
			logx.Field("query", a.query))
	} else {
		logx.WithContext(ctx).Infow("app-request-pre", logx.Field("app-name", a.app.name),
			logx.Field("url", url),
			logx.Field("header", a.header),
			logx.Field("query", a.query))
	}

	var data io.Reader

	if a.data != nil {
		dataBytes, err := json.Marshal(a.data)
		if err != nil {
			return nil, fmt.Errorf("app client request json marshal failed, app: %s, uri: %s, err: %s", a.app.name, uri, err.Error())
		}
		data = bytes.NewReader(dataBytes)
	}

	r, err := http.NewRequest(method, url, data)
	if err != nil {
		return nil, err
	}

	if method == http.MethodPost && a.data != nil {
		r.Header.Set("Content-Type", "application/json")
	}

	if len(a.header) > 0 {
		for k, v := range a.header {
			r.Header.Set(k, v)
		}
	}

	a.resp, err = DoRequest(r)
	if err != nil {
		return nil, fmt.Errorf("app client request failed, app: %s, uri: %s, err: %s", a.app.name, uri, err.Error())
	}

	if a.resp != nil {
		defer a.resp.Body.Close()
	}

	body, err := ioutil.ReadAll(a.resp.Body)
	if err != nil {
		logx.WithContext(ctx).Errorw("app-client-request-failed", logx.Field("app-name", a.app.name),
			logx.Field("err", err),
			logx.Field("uri", uri),
			logx.Field("param", a.data),
			logx.Field("query", a.query),
			logx.Field("body", string(body)))

		return nil, err
	}

	if a.isWriteReplyLog {
		logx.WithContext(ctx).Infow("app-reply", logx.Field("app-name", a.app.name),
			logx.Field("uri", uri),
			logx.Field("param", a.data),
			logx.Field("query", a.query),
			logx.Field("body", string(body)))
	}

	return body, nil
}

func initTracer(projectId string) (func(), error) {

	// Create Google Cloud Trace exporter to be able to retrieve
	// the collected spans.
	tp, _ := trace.GetCloudTraceProvider(projectId)
	if tp != nil {
		otel.SetTracerProvider(tp)
	}

	return func() {
		err := tp.Shutdown(context.Background())
		if err != nil {
			fmt.Printf("error shutting down trace provider: %+v", err)
		}
	}, nil
}

func installPropagators() {
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			// Putting the CloudTraceOneWayPropagator first means the TraceContext propagator
			// takes precedence if both the traceparent and the XCTC headers exist.
			gcppropagator.CloudTraceOneWayPropagator{},
			propagation.TraceContext{},
			propagation.Baggage{},
		))
}

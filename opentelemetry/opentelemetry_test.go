package opentelemetry

import (
	"context"
	"fmt"
	"net/http"
	"testing"
)

func TestSetupStdout(t *testing.T) {
	Setup(WithTraceExporter(&Exporter{
		ExporterEndpoint: "",
		Insecure:         false,
		Creds:            nil,
	}), WithMetricExporter(&Exporter{
		ExporterEndpoint: "",
		Insecure:         false,
		Creds:            nil,
	}))
	if TracerIsEnable() {
		t.Errorf("want tracer is enable is false before setup, got true")
		return
	}
	traceStop()
	metricStop()
	_, testspan1 := SpanStart(context.Background(), "testspan1")
	if testspan1 != nil {
		t.Errorf("want span nil")
	}
	Setup(WithTraceExporter(&Exporter{
		ExporterEndpoint: DefaultStdout,
		Insecure:         false,
		Creds:            nil,
	}), WithMetricExporter(&Exporter{
		ExporterEndpoint: DefaultStdout,
		Insecure:         false,
		Creds:            nil,
	}))
	if !TracerIsEnable() {
		t.Errorf("want tracer is enable is true, got false")
		return
	}
	ctx, testspan2 := SpanStart(context.Background(), "testspan2")
	if testspan2 == nil {
		t.Errorf("want span not nil")
	}
	traceid := GetTraceId(ctx)
	if traceid == "" {
		t.Errorf("want tracer id is not nil")
		return
	}
	Stop()
}

type testHandler struct {
}

func (t *testHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
}

type testResponse struct {
}

func (*testResponse) Header() http.Header {
	return map[string][]string{}
}

func (*testResponse) Write([]byte) (int, error) {
	return 1, nil
}

func (*testResponse) WriteHeader(statusCode int) {
}

func TestInitHttpHandler(t *testing.T) {
	testh := &testHandler{}
	h := InitHttpHandler(testh)
	if h == testh || h == nil {
		t.Errorf("get unexpect handler")
		return
	}

	req, err := http.NewRequest(http.MethodGet, "http://127.0.0.1:8080/healthx", nil)
	if err != nil {
		t.Errorf("new request error %v", err)
		return
	}
	h.ServeHTTP(&testResponse{}, req)

	req, err = http.NewRequest(http.MethodGet, "http://127.0.0.1:8380/heal11", nil)
	if err != nil {
		t.Errorf("new request error %v", err)
		return
	}
	h.ServeHTTP(&testResponse{}, req)
}

func TestGrpcTraceInterceptor(t *testing.T) {
	unary := GrpcUnaryTraceInterceptor()
	if unary == nil {
		t.Errorf("want unary interceptor not nil")
		return
	}
	stream := GrpcStreamTraceInterceptor()
	if stream == nil {
		t.Errorf("want stream interceptor not nil")
	}
}

func Test_errorHandler_Handle(t *testing.T) {
	errorHandler{}.Handle(fmt.Errorf("this is test errorHandler"))
}

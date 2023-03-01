package opentelemetry

import (
	"context"
	"fmt"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"testing"
)

var errString = ""

func TestWithTraceOptions(t *testing.T) {
	testConfig := &Config{
		traceExporter:      &Exporter{},
		metricExporter:     &Exporter{},
		metricReportPeriod: "",
		serviceInfo:        &ServiceInfo{},
		attributes:         nil,
		headers:            nil,
		idGenerator:        nil,
		resource:           nil,
		otelErrorHandler:   nil,
		traceBatchOptions:  nil,
	}
	testTraceExporter := &Exporter{
		ExporterEndpoint: "test.trace.exporter",
		Insecure:         true,
		Creds:            nil,
	}
	testMetricExporter := &Exporter{
		ExporterEndpoint: "test.metric.exporter",
		Insecure:         true,
		Creds:            nil,
	}
	testMetricReportPeriod := "33s"
	testServiceInfo := &ServiceInfo{
		Name:      "testService",
		Namespace: "testServiceNamespace",
		Version:   "testVersion",
	}
	testAttributes := map[string]string{
		"aa1": "bb2",
		"cc3": "dd4",
	}
	testHeaders := map[string]string{
		"hkey1": "hvalue1",
		"hkey2": "hvalue2",
	}
	testResourse := resource.NewSchemaless(attribute.KeyValue{"attrkey", attribute.StringValue("attrva")})
	opts := []Option{}
	opts = append(opts, WithTraceExporter(testTraceExporter), WithMetricExporter(testMetricExporter), WithSampleRatio(0.5),
		WithMetricReportPeriod(testMetricReportPeriod), WithServiceInfo(testServiceInfo), WithAttributes(testAttributes),
		WithHeaders(testHeaders), WithResource(testResourse), WithIDGenerator(&testIDGenerator{}), WithErrorHandler(&testErrorHandler{}),
		WithTraceBatchOptions([]sdktrace.BatchSpanProcessorOption{sdktrace.WithMaxQueueSize(100), sdktrace.WithMaxExportBatchSize(30)}))
	for _, opt := range opts {
		opt(testConfig)
	}
	if testConfig.traceExporter.ExporterEndpoint != "test.trace.exporter" {
		t.Errorf("want endpoint test.trace.exporter, get endpoint %v", testConfig.traceExporter.ExporterEndpoint)
		return
	}
	if testConfig.metricExporter.ExporterEndpoint != "test.metric.exporter" {
		t.Errorf("want endpoint test.trace.exporter, get endpoint %v", testConfig.metricExporter.ExporterEndpoint)
		return
	}
	if testConfig.metricReportPeriod != "33s" {
		t.Errorf("want metric report period 33s, get %v", testConfig.metricReportPeriod)
		return
	}
	if testConfig.serviceInfo.Name != "testService" || testConfig.serviceInfo.Namespace != "testServiceNamespace" ||
		testConfig.serviceInfo.Version != "testVersion" {
		t.Errorf("get unexpected service %v", testConfig.serviceInfo)
		return
	}
	for k, v := range testAttributes {
		if va, ok := testConfig.attributes[k]; !ok || va != v {
			t.Errorf("testConfig attributes %v is not equal testAttributes %v", testConfig.attributes, testAttributes)
			return
		}
	}
	for k, v := range testHeaders {
		if va, ok := testConfig.headers[k]; !ok || va != v {
			t.Errorf("testConfig attributes %v is not equal testAttributes %v", testConfig.headers, testHeaders)
			return
		}
	}

	idG := testConfig.idGenerator
	tid, _ := idG.NewIDs(context.Background())

	if tid.String() != "31323334353637383930313233343536" {
		t.Errorf("get unespected trace id %v", tid.String())
		return
	}
	if testConfig.resource.String() != "attrkey=attrva" {
		t.Errorf("get unespected string %v", testConfig.resource.String())
		return
	}
	otelErrorH := testConfig.otelErrorHandler
	otelErrorH.Handle(fmt.Errorf("this is test error"))
	if errString != "this is error handler" {
		t.Errorf("want errString 'this is error handler', got %v", errString)
		return
	}
	batchConfig := &sdktrace.BatchSpanProcessorOptions{}
	for _, batchOption := range testConfig.traceBatchOptions {
		batchOption(batchConfig)
	}
	if batchConfig.MaxQueueSize != 100 || batchConfig.MaxExportBatchSize != 30 {
		t.Errorf("want MaxQueueSize 100, MaxExportBatchSize 30, got %v", batchConfig)
		return
	}
	Setup(opts...)
}

type testIDGenerator struct {
}

func (*testIDGenerator) NewIDs(ctx context.Context) (trace.TraceID, trace.SpanID) {
	var ret [16]byte
	copy(ret[:], "1234567890123456")
	return ret, trace.SpanID{}
}

func (*testIDGenerator) NewSpanID(ctx context.Context, traceID trace.TraceID) trace.SpanID {
	return trace.SpanID{}
}

type testErrorHandler struct {
}

func (*testErrorHandler) Handle(error) {
	errString = "this is error handler"
}

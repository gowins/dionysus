package opentelemetry

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gowins/dionysus/healthy"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/contrib/instrumentation/host"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	otlpTraceGrpc "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	oteltrace "go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/encoding/gzip"
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	traceStop     = func() {}
	metricStop    = func() {}
	DefaultCred   = credentials.NewClientTLSFromCert(nil, "")
	tracerEnable  bool
	defaultTracer oteltrace.Tracer
)

const (
	DefaultStdout       = "stdout"
	grpcHealthyMethod   = "/healthy.HealthService/Health"
	healthHost          = "127.0.0.1:8080"
	instrumentationName = "github.com/gowins/dionysus/opentelemetry"
)

func Setup(opts ...Option) {
	config := &Config{
		traceExporter:      &Exporter{},
		metricExporter:     &Exporter{},
		metricReportPeriod: "",
		serviceInfo:        &ServiceInfo{},
		attributes:         map[string]string{},
		headers:            map[string]string{},
		idGenerator:        nil,
		resource:           &resource.Resource{},
		otelErrorHandler:   errorHandler{},
		traceBatcOptions:   []sdktrace.BatchSpanProcessorOption{},
	}
	for _, opt := range opts {
		opt(config)
	}
	otel.SetErrorHandler(config.otelErrorHandler)
	mergeResource(config)
	if config.traceExporter.ExporterEndpoint != "" {
		initTracer(config)
	}
	if config.metricExporter.ExporterEndpoint != "" {
		initMetric(config)
	}
}

func TracerIsEnable() bool {
	return tracerEnable
}

func InitHttpHandler(handler http.Handler) http.Handler {
	return otelhttp.NewHandler(handler, "dioHttp", otelhttp.WithFilter(func(request *http.Request) bool {
		if request.Host == healthHost && strings.HasPrefix(request.URL.Path, healthy.HealthGroupPath) {
			return false
		}
		return true
	}))
}

func GrpcUnaryTraceInterceptor() grpc.UnaryServerInterceptor {
	return otelgrpc.UnaryServerInterceptor(otelgrpc.WithInterceptorFilter(func(info *otelgrpc.InterceptorInfo) bool {
		if strings.HasPrefix(info.UnaryServerInfo.FullMethod, grpcHealthyMethod) {
			return false
		}
		return true
	}))
}

func GrpcStreamTraceInterceptor() grpc.StreamServerInterceptor {
	return otelgrpc.StreamServerInterceptor()
}

type errorHandler struct {
}

func (errorHandler) Handle(err error) {
	log.Errorf("opentelemetry tracer exporter error %v", err)
}

func Stop() {
	metricStop()
	traceStop()
}

func mergeResource(config *Config) {
	hostname, _ := os.Hostname()
	defaultResource := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(config.serviceInfo.Name),
		semconv.HostNameKey.String(hostname),
		semconv.ServiceNamespaceKey.String(config.serviceInfo.Namespace),
		semconv.ServiceVersionKey.String(config.serviceInfo.Version),
		semconv.ProcessPIDKey.Int(os.Getpid()),
		semconv.ProcessCommandKey.String(os.Args[0]),
	)
	var err error
	if config.resource, err = resource.Merge(defaultResource, config.resource); err != nil {
		log.Fatalf("resource merge with default resource error %v", err)
	}
	if config.resource, err = resource.Merge(resource.Environment(), config.resource); err != nil {
		log.Fatalf("resource merge with environment resource error %v", err)
	}

	resource.Environment()
}

func initTracer(config *Config) {
	traceExporter, err := initTracerExporter(config)
	if err != nil {
		log.Fatalf("new traceExporter error %v", err)
		return
	}
	traceProvider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(traceExporter, config.traceBatcOptions...),
		sdktrace.WithIDGenerator(config.idGenerator),
		sdktrace.WithResource(config.resource),
	)
	otel.SetTracerProvider(traceProvider)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	traceStop = func() {
		log.Infof("trace Provider is shutdown")
		if err = traceProvider.Shutdown(context.Background()); err != nil {
			log.Errorf("trace Provider shutdown error %v\n")
		}
	}
	tracerEnable = true
	name := instrumentationName + "/" + config.serviceInfo.Namespace + "/" + config.serviceInfo.Name
	defaultTracer = otel.GetTracerProvider().Tracer(name, oteltrace.WithInstrumentationVersion("v1.1.0"))
}

// SpanStart creates a span and a context.Context containing the newly-created span.
//
// If the context.Context provided in `ctx` contains a Span then the newly-created
// Span will be a child of that span, otherwise it will be a root span. This behavior
// can be overridden by providing `WithNewRoot()` as a SpanOption, causing the
// newly-created Span to be a root span even if `ctx` contains a Span.
//
// When creating a Span it is recommended to provide all known span attributes using
// the `WithAttributes()` SpanOption as samplers will only have access to the
// attributes provided when a Span is created.
//
// Any Span that is created MUST also be ended. This is the responsibility of the user.
// Implementations of this API may leak memory or other resources if Spans are not ended.
func SpanStart(ctx context.Context, spanName string, opts ...oteltrace.SpanStartOption) (context.Context, oteltrace.Span) {
	if !tracerEnable {
		return ctx, nil
	}
	return defaultTracer.Start(ctx, spanName, opts...)
}

func initTracerExporter(config *Config) (trace.SpanExporter, error) {
	if config.traceExporter.ExporterEndpoint == DefaultStdout {
		return stdouttrace.New(stdouttrace.WithPrettyPrint())
	} else if config.traceExporter.ExporterEndpoint != "" {
		traceSecureOption := otlpTraceGrpc.WithTLSCredentials(config.traceExporter.Creds)
		if config.traceExporter.Insecure {
			traceSecureOption = otlpTraceGrpc.WithInsecure()
		}
		return otlptrace.New(context.Background(),
			otlpTraceGrpc.NewClient(otlpTraceGrpc.WithEndpoint(config.traceExporter.ExporterEndpoint),
				traceSecureOption, otlpTraceGrpc.WithHeaders(config.headers), otlpTraceGrpc.WithCompressor(gzip.Name)))
	}
	return nil, fmt.Errorf("tracer exporter endpoint is nil, no exporter is inited")
}

func initMetric(config *Config) {
	metricsExporter, err := initMetricExporter(config)
	if err != nil {
		log.Fatalf("new metricsExporter error %v", err)
		return
	}
	period, err := time.ParseDuration(config.metricReportPeriod)
	if err != nil {
		period = time.Second * 30
	}
	reader := metric.NewPeriodicReader(metricsExporter, metric.WithInterval(period))
	meterProvider := metric.NewMeterProvider(metric.WithReader(reader), metric.WithResource(config.resource))
	global.SetMeterProvider(meterProvider)
	err = host.Start(host.WithMeterProvider(meterProvider))
	if err != nil {
		log.Fatalf("host metric start error %v", err)
	}
	err = runtime.Start(runtime.WithMeterProvider(meterProvider), runtime.WithMinimumReadMemStatsInterval(5*time.Second))
	if err != nil {
		log.Fatalf("golang runtime metric start error %v", err)
	}
	metricStop = func() {
		log.Infof("metrics Provider is shutdown")
		if err = metricsExporter.Shutdown(context.Background()); err != nil {
			log.Errorf("metrics Provider shutdown error %v", err)
		}
	}
}

func initMetricExporter(config *Config) (metric.Exporter, error) {
	if config.metricExporter.ExporterEndpoint == DefaultStdout {
		encoder := json.NewEncoder(os.Stdout)
		return stdoutmetric.New(stdoutmetric.WithEncoder(encoder))
	} else if config.metricExporter.ExporterEndpoint != "" {
		metricSecureOption := otlpmetricgrpc.WithTLSCredentials(config.metricExporter.Creds)
		if config.metricExporter.Insecure {
			metricSecureOption = otlpmetricgrpc.WithInsecure()
		}
		return otlpmetricgrpc.New(context.Background(), otlpmetricgrpc.WithEndpoint(config.metricExporter.ExporterEndpoint),
			metricSecureOption, otlpmetricgrpc.WithHeaders(config.headers), otlpmetricgrpc.WithCompressor(gzip.Name))
	}
	return nil, fmt.Errorf("metric exporter endpoint is nil, no exporter is inited")
}

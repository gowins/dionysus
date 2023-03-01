package opentelemetry

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc/credentials"
)

type Config struct {
	traceExporter      *Exporter
	metricExporter     *Exporter
	metricReportPeriod string
	serviceInfo        *ServiceInfo
	attributes         map[string]string
	headers            map[string]string
	idGenerator        sdktrace.IDGenerator
	resource           *resource.Resource
	otelErrorHandler   otel.ErrorHandler
	traceBatchOptions  []sdktrace.BatchSpanProcessorOption
	sampleRatio        float64
}

type Exporter struct {
	ExporterEndpoint string
	Insecure         bool
	Creds            credentials.TransportCredentials
}

type ServiceInfo struct {
	Name      string
	Namespace string
	Version   string
}

type Option func(*Config)

// WithTraceExporter sets export that spans will send to
func WithTraceExporter(traceExporter *Exporter) Option {
	return func(config *Config) {
		config.traceExporter = traceExporter
	}
}

// WithMetricExporter sets export that metric will send to
func WithMetricExporter(metricExporter *Exporter) Option {
	return func(config *Config) {
		config.metricExporter = metricExporter
	}
}

// WithMetricReportPeriod set period
func WithMetricReportPeriod(period string) Option {
	return func(config *Config) {
		config.metricReportPeriod = period
	}
}

// WithServiceInfo set service info
func WithServiceInfo(serviceInfo *ServiceInfo) Option {
	return func(config *Config) {
		config.serviceInfo = serviceInfo
	}
}

// WithAttributes set attributes which span send with
func WithAttributes(attributes map[string]string) Option {
	return func(config *Config) {
		config.attributes = attributes
	}
}

// WithHeaders set headers which span send with
func WithHeaders(headers map[string]string) Option {
	return func(config *Config) {
		config.headers = headers
	}
}

// WithIDGenerator set id generator for trace id and span id
func WithIDGenerator(idGenerator sdktrace.IDGenerator) Option {
	return func(config *Config) {
		config.idGenerator = idGenerator
	}
}

// WithResource set resource which span send with
func WithResource(resource *resource.Resource) Option {
	return func(config *Config) {
		config.resource = resource
	}
}

// WithErrorHandler set error handler when span send failed
func WithErrorHandler(errorHandler otel.ErrorHandler) Option {
	return func(config *Config) {
		config.otelErrorHandler = errorHandler
	}
}

// WithTraceBatchOptions set configuration settings for BatchSpanProcessor
func WithTraceBatchOptions(options []sdktrace.BatchSpanProcessorOption) Option {
	return func(config *Config) {
		config.traceBatchOptions = options
	}
}

// WithSampleRatio set ratio of traces samples
func WithSampleRatio(sampleRatio float64) Option {
	return func(config *Config) {
		config.sampleRatio = sampleRatio
	}
}

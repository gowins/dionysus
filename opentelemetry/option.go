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
	traceBatcOptions   []sdktrace.BatchSpanProcessorOption
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

func WithTraceExporter(traceExporter *Exporter) Option {
	return func(config *Config) {
		config.traceExporter = traceExporter
	}
}

func WithMetricExporter(metricExporter *Exporter) Option {
	return func(config *Config) {
		config.metricExporter = metricExporter
	}
}

func WithMetricReportPeriod(period string) Option {
	return func(config *Config) {
		config.metricReportPeriod = period
	}
}

func WithServiceInfo(serviceInfo *ServiceInfo) Option {
	return func(config *Config) {
		config.serviceInfo = serviceInfo
	}
}

func WithAttributes(attributes map[string]string) Option {
	return func(config *Config) {
		config.attributes = attributes
	}
}

func WithHeaders(headers map[string]string) Option {
	return func(config *Config) {
		config.headers = headers
	}
}

func WithIDGenerator(idGenerator sdktrace.IDGenerator) Option {
	return func(config *Config) {
		config.idGenerator = idGenerator
	}
}

func WithResource(resource *resource.Resource) Option {
	return func(config *Config) {
		config.resource = resource
	}
}

func WithErrorHandler(errorHandler otel.ErrorHandler) Option {
	return func(config *Config) {
		config.otelErrorHandler = errorHandler
	}
}

func WithTraceBatchOptions(options []sdktrace.BatchSpanProcessorOption) Option {
	return func(config *Config) {
		config.traceBatcOptions = options
	}
}

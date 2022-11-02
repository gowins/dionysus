package serverinterceptors

import (
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"google.golang.org/grpc"
)

func OpenTracingUnary(opts ...grpc_opentracing.Option) grpc.UnaryServerInterceptor {
	return grpc_opentracing.UnaryServerInterceptor(opts...)
}

func OpenTracingStream(opts ...grpc_opentracing.Option) grpc.StreamServerInterceptor {
	return grpc_opentracing.StreamServerInterceptor(opts...)
}

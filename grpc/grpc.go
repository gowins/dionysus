package grpc

import (
	"github.com/gowins/dionysus/grpc/balancer/resolver"
	"github.com/gowins/dionysus/grpc/client"
	logger "github.com/gowins/dionysus/log"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"

	// 默认注册中心为etcdv3
	_ "github.com/gowins/dionysus/grpc/registry/etcdv3"
	"google.golang.org/grpc"
)

func WithDefaultClientStreamInterceptors(interceptors ...grpc.StreamClientInterceptor) grpc.DialOption {
	return grpc.WithChainStreamInterceptor(append([]grpc.StreamClientInterceptor{
		grpc_opentracing.StreamClientInterceptor(),
	}, interceptors...)...)
}

func WithDefaultClientUnaryInterceptors(interceptors ...grpc.UnaryClientInterceptor) grpc.DialOption {
	return grpc.WithChainUnaryInterceptor(append([]grpc.UnaryClientInterceptor{
		grpc_opentracing.UnaryClientInterceptor(),
	}, interceptors...)...)
}

func NewClient(service string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	return client.New(service, opts...)
}

func GetClient(service string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	return client.Get(service, opts...)
}

// func NewServer() *Server { return server.New() }

func SetLog(log logger.Logger) {
	// server.SetLog(log)
	resolver.SetLog(log)
	// serverinterceptors.SetLog(log)
}

package client

import (
	"context"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"

	"github.com/gowins/dionysus/grpc/balancer/resolver"
	"github.com/gowins/dionysus/grpc/registry"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
)

func buildTarget(service string) string {
	// 注册中心为nil走直连模式
	if registry.Default == nil {
		return resolver.BuildDirectTarget([]string{service})
	}

	return resolver.BuildDiscovTarget([]string{registry.Default.String()}, service)
}

func New(service string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	target := buildTarget(service)
	conn, err := dial(target, append(defaultDialOpts, opts...)...)
	if err != nil {
		return nil, errors.Wrapf(err, "dial %s error\n", target)
	}
	return conn, nil
}

func NewConnWithTracer(service string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	opts = append(opts, grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
		grpc.WithStreamInterceptor(otelgrpc.StreamClientInterceptor()))
	return New(service, opts...)
}

func GetConnWithTracer(service string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	opts = append(opts, grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
		grpc.WithStreamInterceptor(otelgrpc.StreamClientInterceptor()))
	return Get(service, opts...)
}

// Get new grpc client
func Get(service string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	if val, ok := clients.Load(service); ok {
		if val.(*grpc.ClientConn).GetState() == connectivity.Ready {
			return val.(*grpc.ClientConn), nil
		}
	}

	mu.Lock()
	defer mu.Unlock()

	// 双检, 避免多次创建
	if val, ok := clients.Load(service); ok {
		if val.(*grpc.ClientConn).GetState() == connectivity.Ready {
			return val.(*grpc.ClientConn), nil
		}
	}

	conn, err := New(service, opts...)
	if err != nil {
		return nil, err
	}

	clients.Store(service, conn)
	return conn, nil
}

func dial(target string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()

	conn, err := grpc.DialContext(ctx, target, opts...)
	if err != nil {
		return nil, errors.Wrapf(err, "DialContext error, target:%s\n", target)
	}
	return conn, nil
}

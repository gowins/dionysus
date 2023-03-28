package clientinterceptors

import (
	"context"
	"google.golang.org/grpc"

	"github.com/afex/hystrix-go/hystrix"
	"github.com/gowins/dionysus/grpc/clientinterceptors/grpc_hystrix"
	"github.com/gowins/dionysus/grpc/ghystrix"
)

func BreakerUnary(cc hystrix.CommandConfig, cfgs ...ghystrix.HystrixCfg) grpc.UnaryClientInterceptor {
	return grpc_hystrix.UnaryClientInterceptor(cc, cfgs...)
}

func BreakerStream(cc hystrix.CommandConfig, cfgs ...ghystrix.HystrixCfg) grpc.StreamClientInterceptor {
	return grpc_hystrix.StreamClientInterceptor(cc, cfgs...)
}

func BreakerUnaryByService(serviceName string, cc hystrix.CommandConfig) grpc.UnaryClientInterceptor {
	hystrix.ConfigureCommand(serviceName, cc)
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		return hystrix.Do(serviceName, func() error { return invoker(ctx, method, req, reply, cc, opts...) }, nil)
	}
}

func BreakerStreamByService(serviceName string, cc hystrix.CommandConfig) grpc.StreamClientInterceptor {
	hystrix.ConfigureCommand(serviceName, cc)
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (cs grpc.ClientStream, err error) {
		return cs, hystrix.Do(serviceName, func() error { cs, err = streamer(ctx, desc, cc, method, opts...); return err }, nil)
	}
}

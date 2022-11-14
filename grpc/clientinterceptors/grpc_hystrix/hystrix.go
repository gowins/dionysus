package grpc_hystrix

import (
	"context"

	"github.com/afex/hystrix-go/hystrix"
	"github.com/gowins/dionysus/grpc/ghystrix"
	"google.golang.org/grpc"
)

func UnaryClientInterceptor(cc hystrix.CommandConfig, cfgs ...ghystrix.HystrixCfg) grpc.UnaryClientInterceptor {
	ghystrix.HystrixDefault(cc)
	for _, cfg := range cfgs {
		hystrix.ConfigureCommand(cfg.Name, cfg.Cfg)
	}
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		return hystrix.Do(method, func() error { return invoker(ctx, method, req, reply, cc, opts...) }, nil)
	}
}

func StreamClientInterceptor(cc hystrix.CommandConfig, cfgs ...ghystrix.HystrixCfg) grpc.StreamClientInterceptor {
	ghystrix.HystrixDefault(cc)
	for _, cfg := range cfgs {
		hystrix.ConfigureCommand(cfg.Name, cfg.Cfg)
	}
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (cs grpc.ClientStream, err error) {
		return cs, hystrix.Do(method, func() error { cs, err = streamer(ctx, desc, cc, method, opts...); return err }, nil)
	}
}

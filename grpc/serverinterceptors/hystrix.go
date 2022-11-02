package serverinterceptors

import (
	"context"

	"github.com/afex/hystrix-go/hystrix"
	"github.com/gowins/dionysus/grpc/ghystrix"
	"google.golang.org/grpc"
)

func HystrixRateLimitUnary(cc hystrix.CommandConfig, cfgs ...ghystrix.HystrixCfg) grpc.UnaryServerInterceptor {
	ghystrix.HystrixDefault(cc)
	for _, cfg := range cfgs {
		hystrix.ConfigureCommand(cfg.Name, cfg.Cfg)
	}
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		err = hystrix.Do(info.FullMethod, func() error {
			resp, err = handler(ctx, req)
			return err
		}, nil)
		return
	}
}

func HystrixRateLimitStream(cc hystrix.CommandConfig, cfgs ...ghystrix.HystrixCfg) grpc.StreamServerInterceptor {
	ghystrix.HystrixDefault(cc)
	for _, cfg := range cfgs {
		hystrix.ConfigureCommand(cfg.Name, cfg.Cfg)
	}
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		return hystrix.Do(info.FullMethod, func() error {
			return handler(srv, ss)
		}, nil)
	}
}

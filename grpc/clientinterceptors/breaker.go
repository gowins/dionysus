package clientinterceptors

import (
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

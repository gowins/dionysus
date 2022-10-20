package clientinterceptors

import (
	"google.golang.org/grpc"

	"github.com/gowins/dionysus/grpc/clientinterceptors/grpc_hystrix"
)

func BreakerUnary(opts ...grpc_hystrix.Option) grpc.UnaryClientInterceptor {
	return grpc_hystrix.UnaryClientInterceptor(opts...)
}

func BreakerStream(opts ...grpc_hystrix.Option) grpc.StreamClientInterceptor {
	return grpc_hystrix.StreamClientInterceptor(opts...)
}

package serverinterceptors

import (
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"google.golang.org/grpc"
)

func AuthUnary(authFunc grpc_auth.AuthFunc) grpc.UnaryServerInterceptor {
	return grpc_auth.UnaryServerInterceptor(authFunc)
}

func AuthStream(authFunc grpc_auth.AuthFunc) grpc.StreamServerInterceptor {
	return grpc_auth.StreamServerInterceptor(authFunc)
}

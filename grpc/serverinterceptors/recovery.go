package serverinterceptors

import (
	"context"

	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func RecoveryUnary(f grpc_recovery.RecoveryHandlerFuncContext) grpc.UnaryServerInterceptor {
	return grpc_recovery.UnaryServerInterceptor(grpc_recovery.WithRecoveryHandlerContext(f))
}

func RecoveryStream(f grpc_recovery.RecoveryHandlerFuncContext) grpc.StreamServerInterceptor {
	return grpc_recovery.StreamServerInterceptor(grpc_recovery.WithRecoveryHandlerContext(f))
}

func DefaultRecovery() grpc_recovery.RecoveryHandlerFuncContext {
	return func(ctx context.Context, p interface{}) (err error) {
		log.Errorf("dionysus grpc server panic %v", p)
		return status.Errorf(codes.Internal, "[grpc] server recovery error, err: %v", p)
	}
}

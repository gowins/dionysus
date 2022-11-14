package serverinterceptors

import (
	"github.com/grpc-ecosystem/go-grpc-middleware/ratelimit"
	"google.golang.org/grpc"
)

type RateLimit struct{}

func (r RateLimit) Limit() bool { return false }

func RateLimitUnary(limiter ratelimit.Limiter) grpc.UnaryServerInterceptor {
	return ratelimit.UnaryServerInterceptor(limiter)
}

func RateLimitStream(limiter ratelimit.Limiter) grpc.StreamServerInterceptor {
	return ratelimit.StreamServerInterceptor(limiter)
}

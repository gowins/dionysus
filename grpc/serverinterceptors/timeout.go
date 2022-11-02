package serverinterceptors

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	DefaultUnaryTimeout  = 10 * time.Second
	DefaultStreamTimeout = 10 * time.Second
)

// TimeoutUnary returns a new unary server interceptor for OpenTracing.
func TimeoutUnary(t time.Duration) grpc.UnaryServerInterceptor {
	defaultTimeOut := DefaultUnaryTimeout
	if t < DefaultUnaryTimeout {
		defaultTimeOut = t
	}

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if _, ok := ctx.Deadline(); !ok { //if ok is true, it is set by header grpc-timeout from client
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(ctx, defaultTimeOut)
			defer cancel()
		}

		// create a done channel to tell the request it's done
		done := make(chan struct{})

		// here you put the actual work needed for the request
		// and then send the doneChan with the status and body
		// to finish the request by writing the response
		var res interface{}
		var err error

		go func() {
			defer func() {
				if c := recover(); c != nil {
					log.Errorf("response request panic: %v", c)
					err = status.Errorf(codes.Internal, "response request panic: %v", c)
				}
				close(done)
			}()
			res, err = handler(ctx, req)
		}()

		// non-blocking select on two channels see if the request
		// times out or finishes
		select {

		// if the context is done it timed out or was canceled
		// so don't return anything
		case <-ctx.Done():
			return nil, status.Errorf(codes.DeadlineExceeded, "handler timeout")

		// if the request finished then finish the request by
		// writing the response
		case <-done:
			return res, err
		}
	}
}

// TimeoutStream returns a new streaming server interceptor for panic recovery.
func TimeoutStream(t time.Duration) grpc.StreamServerInterceptor {
	defaultTimeOut := DefaultStreamTimeout
	if t < DefaultStreamTimeout {
		defaultTimeOut = t
	}

	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		ctx := stream.Context()
		if _, ok := ctx.Deadline(); !ok { //if ok is true, it is set by header grpc-timeout from client
			if defaultTimeOut == 0 {
				return handler(srv, stream)
			}

			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(ctx, defaultTimeOut)
			defer cancel()
		}
		// create a done channel to tell the request it's done
		done := make(chan struct{})

		go func() {
			defer func() {
				if c := recover(); c != nil {
					log.Errorf("response request panic: %v", c)
					err = status.Errorf(codes.Internal, "response request panic: %v", c)
				}
				close(done)
			}()
			err = handler(srv, stream)
		}()

		// non-blocking select on two channels see if the request
		// times out or finishes
		select {

		// if the context is done it timed out or was canceled
		// so don't return anything
		case <-ctx.Done():
			return status.Errorf(codes.DeadlineExceeded, "handler timeout")

		// if the request finished then finish the request by
		// writing the response
		case <-done:
			return err
		}
	}
}

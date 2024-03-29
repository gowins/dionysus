package pool

import (
	"sync"
	"time"

	"google.golang.org/grpc"
)

type Option func(*GrpcPool)

func WithPoolSize(poolSize int) Option {
	return func(pool *GrpcPool) {
		pool.poolSize = poolSize
	}
}

func WithDialOptions(dialOptions []grpc.DialOption) Option {
	return func(pool *GrpcPool) {
		pool.dialOptions = dialOptions
	}
}

func WithScaleOption(scaleOption *ScaleOption) Option {
	return func(pool *GrpcPool) {
		pool.scaleOption = scaleOption
	}
}

const (
	// KeepAliveTime is the duration of time after which if the client doesn't see
	// any activity it pings the server to see if the transport is still alive.
	KeepAliveTime = time.Duration(10) * time.Second

	// KeepAliveTimeout is the duration of time for which the client waits after having
	// pinged for keepalive check and if no activity is seen even after that the connection
	// is closed.
	KeepAliveTimeout = time.Duration(3) * time.Second

	// InitialWindowSize we set it 256M is to provide system's throughput.
	InitialWindowSize = 1 << 28

	// InitialConnWindowSize we set it 256M is to provide system's throughput.
	InitialConnWindowSize = 1 << 28

	// MaxSendMsgSize set max gRPC request message poolSize sent to server.
	// If any request message poolSize is larger than current value, an error will be reported from gRPC.
	MaxSendMsgSize = 1 << 30

	// MaxRecvMsgSize set max gRPC receive message poolSize received from server.
	// If any message poolSize is larger than current value, an error will be reported from gRPC.
	MaxRecvMsgSize = 1 << 30
)

var DefaultScaleOption = &ScaleOption{
	Enable:          true,
	ScalePeriod:     time.Second * 30,
	MaxConn:         300,
	DesireMaxStream: 80,
}

var DefaultDialTimeout = 3 * time.Second

//var DefaultDialOpts = []grpc.DialOption{
//	grpc.WithTransportCredentials(insecure.NewCredentials()),
//	grpc.WithBlock(),
//	grpc.WithInitialWindowSize(InitialWindowSize),
//	grpc.WithInitialConnWindowSize(InitialConnWindowSize),
//	grpc.WithDefaultCallOptions(grpc.MaxCallSendMsgSize(MaxSendMsgSize)),
//	grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(MaxRecvMsgSize)),
//	grpc.WithKeepaliveParams(keepalive.ClientParameters{
//		Time:                KeepAliveTime,
//		Timeout:             KeepAliveTimeout,
//		PermitWithoutStream: true,
//	}),
//}

var defaultPoolSize = 3

var grpcPool sync.Map
var poolInit sync.Mutex

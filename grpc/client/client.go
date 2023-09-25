package client

import (
	"context"
	"fmt"
	"runtime"
	"sync/atomic"
	"time"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

const (
	// KeepAliveTime is the duration of time after which if the client doesn't see
	// any activity it pings the server to see if the transport is still alive.
	keepAliveTime = 10 * time.Second

	// KeepAliveTimeout is the duration of time for which the client waits after having
	// pinged for keepalive check and if no activity is seen even after that the connection
	// is closed.
	keepAliveTimeout = 3 * time.Second

	defaultDialDeadline = 2 * time.Second

	//default PoolSize
	defaultPoolCardinalSize = 10

	defaultMinIdleConns = 2

	//minutes
	defaultMaxConnAge = 16

	defaultIdleCheckFrequency = 1
)

var defaultDialOpts = []grpc.DialOption{
	grpc.WithTransportCredentials(insecure.NewCredentials()),
	grpc.WithBlock(),
	grpc.WithKeepaliveParams(keepalive.ClientParameters{
		Time:                keepAliveTime,
		Timeout:             keepAliveTimeout,
		PermitWithoutStream: true,
	}),
}

// GrpcPool implemented ClientConnInterface
type GrpcPool struct {
	conns       []*Conn
	next        uint32 // round-robin balancer
	poolCnt     *PoolController
	dialOptions []grpc.DialOption
	deadline    time.Duration
}

// Conn todo
type Conn struct {
	*grpc.ClientConn
}
type PoolController struct {
	// Maximum number of grpc connections.
	// Default is 10 connections per every available CPU as reported by runtime.GOMAXPROCS.
	PoolSize int
	// Maximum number of parallel stream use in ever grpc connections.
	// Default is runtime.GOMAXPROCS. Do not change this value if it is not specifically required,
	// but instead change PoolSize
	MaxCurrentStream int
	// Minimum number of idle connections which is useful when establishing
	// new connection is slow.
	MinIdleConns int
	// Connection age at which client retires (closes) the connection.
	// Default is to not close aged connections.
	MaxConnAge time.Duration
	// Frequency of idle checks made by idle connections reaper.
	// Default is 1 minute. -1 disables idle connections reaper,
	IdleCheckFrequency time.Duration
}

// Invoke sends the RPC request on the wire and returns after response is
func (gp *GrpcPool) Invoke(ctx context.Context, method string, args, reply interface{}, opts ...grpc.CallOption) error {
	return gp.pickLeastConn().Invoke(ctx, method, args, reply, opts...)
}

// NewStream begins a streaming RPC.
func (gp *GrpcPool) NewStream(ctx context.Context, desc *grpc.StreamDesc, method string, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	return gp.pickLeastConn().NewStream(ctx, desc, method, opts...)
}

// Close  Just back conn, GP no need close conn exclude scale according to PoolController
func (gp *GrpcPool) Close() {

}

func (gp *GrpcPool) pickLeastConn() *Conn {
	connsLen := uint32(len(gp.conns))
	nextIndex := atomic.AddUint32(&gp.next, 1)

	conn := gp.conns[nextIndex%connsLen]

	// if conn is not ready, choose a next ready conn
	if conn.GetState() != connectivity.Ready && conn.GetState() != connectivity.Idle {
		var i uint32
		for i = 0; i < connsLen; i++ {
			idx := (i + nextIndex) % connsLen
			if gp.conns[idx].GetState() == connectivity.Ready ||
				gp.conns[idx].GetState() == connectivity.Idle {
				return gp.conns[idx]
			}
		}
	}
	return conn
}

func combine[T any](o1 []T, o2 []T) []T {
	// we don't use append because o1 could have extra capacity whose
	// elements would be overwritten, which could cause inadvertent
	// sharing (and race conditions) between concurrent calls
	if len(o1) == 0 {
		return o2
	} else if len(o2) == 0 {
		return o1
	}
	ret := make([]T, len(o1)+len(o2))
	copy(ret, o1)
	copy(ret[len(o1):], o2)
	return ret
}

type Option = func(*GrpcPool)

func WithDialOptions(dp grpc.DialOption) Option {
	return func(gp *GrpcPool) {
		gp.dialOptions = append(gp.dialOptions, dp)
	}
}

func WithPoolControl(pc *PoolController) Option {
	return func(gp *GrpcPool) {
		gp.poolCnt = pc
	}
}
func WithDialDeadline(t time.Duration) Option {
	return func(gp *GrpcPool) {
		gp.deadline = t
	}
}

//Dial return GrpcPool with opts + default options

func Dial(target string, opts ...Option) (*GrpcPool, error) {
	if target == "" {
		return nil, errors.New("Invalid target")
	}
	gp := &GrpcPool{
		poolCnt: &PoolController{
			PoolSize:           defaultPoolCardinalSize * runtime.GOMAXPROCS(0),
			MaxCurrentStream:   runtime.GOMAXPROCS(0),
			MinIdleConns:       defaultMinIdleConns,
			MaxConnAge:         defaultMaxConnAge * time.Minute,
			IdleCheckFrequency: defaultIdleCheckFrequency,
		},
		deadline: defaultDialDeadline,
	}
	for _, opt := range opts {
		opt(gp)
	}
	//conn size todo: similarly keepalive programs
	gp.conns = make([]*Conn, gp.poolCnt.MinIdleConns)
	//conn timeout
	for i := 0; i < gp.poolCnt.MinIdleConns; i++ {
		dialOpts := combine(defaultDialOpts, gp.dialOptions)
		ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(gp.deadline))
		defer cancel()
		conn, err := grpc.DialContext(ctx, target, dialOpts...)
		if err != nil {
			return gp, fmt.Errorf("grpc dial target %v error %v", target, err)
		}
		gp.conns[i] = &Conn{
			conn,
		}
	}

	return gp, nil
}

package client

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gowins/dionysus/log"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

const (
	// KeepAliveTime is the duration of time after which is the client doesn't see
	// any activity it pings the server to see if the transport is still alive.
	keepAliveTime = 10 * time.Second

	// KeepAliveTimeout is the duration of time for which the client waits after having
	// pinged for keepalive check and if no activity is seen even after that the connection
	// is closed.
	keepAliveTimeout = 3 * time.Second

	defaultDialDeadline = 2 * time.Second

	//default max PoolSize
	defaultPoolCardinalSize = 10

	//default PoolSize
	defaultPoolInitSize = 3

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
	target      string
	conns       []*Conn
	poolSize    uint32
	next        uint32
	poolCnt     *PoolController
	dialOptions []grpc.DialOption
	deadline    time.Duration
	rw          sync.RWMutex
	// Notifications need to be scalable
	ntfChan chan struct{}
}

// Conn
// todo inflight less than maxcurrentstreams
type Conn struct {
	*grpc.ClientConn
	inflight int32
}

type PoolController struct {
	// Maximum number of grpc connections.
	// Default is 10 connections per every available CPU as reported by runtime.GOMAXPROCS.
	PoolSize uint32
	// Maximum number of parallel stream use in ever grpc connections.
	// Default is runtime.GOMAXPROCS. Do not change this value if it is not specifically required,
	// but instead change PoolSize
	MaxCurrentStream uint32
	// Minimum number of idle connections which is useful when establishing
	// new connection is slow.
	MinIdleConns uint32
	// Connection age at which client retires (closes) the connection.
	// Default is to not close aged connections.
	MaxConnAge time.Duration
	// Frequency of idle checks made by idle connections reaper.
	// Default is 1 minute. -1 disables idle connections reaper,
	IdleCheckFrequency time.Duration
}

// Invoke sends the RPC request on the wire and returns after response is
func (gp *GrpcPool) Invoke(ctx context.Context, method string, args, reply interface{}, opts ...grpc.CallOption) error {
	grpcConn, err := gp.pickLeastConn()
	if err != nil {
		return err
	}
	atomic.AddInt32(&grpcConn.inflight, 1)
	defer atomic.AddInt32(&grpcConn.inflight, -1)
	return grpcConn.Invoke(ctx, method, args, reply, opts...)
}

// NewStream begins a streaming RPC.
func (gp *GrpcPool) NewStream(ctx context.Context, desc *grpc.StreamDesc, method string, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	grpcConn, err := gp.pickLeastConn()
	if err != nil {
		return nil, err
	}
	atomic.AddInt32(&grpcConn.inflight, 1)
	defer atomic.AddInt32(&grpcConn.inflight, -1)
	return grpcConn.NewStream(ctx, desc, method, opts...)

}

func (gp *GrpcPool) Close() {

}

func (gp *GrpcPool) getConnsNum() uint32 {
	gp.rw.RLock()
	defer gp.rw.RUnlock()
	return uint32(len(gp.conns))
}
func (gp *GrpcPool) pickLeastConn() (*Conn, error) {
	gp.rw.RLock()
	defer gp.rw.RUnlock()

	nextIndex := atomic.AddUint32(&gp.next, 1)

	conn := gp.conns[nextIndex%gp.poolSize]

	// if conn is not ready, choose a next ready conn
	if conn.GetState() != connectivity.Ready && conn.GetState() != connectivity.Idle {
		var i uint32
		for i = 0; i < gp.poolSize; i++ {
			idx := (i + nextIndex) % gp.poolSize
			if gp.conns[idx].GetState() == connectivity.Ready ||
				gp.conns[idx].GetState() == connectivity.Idle {
				return gp.conns[idx], nil
			}
		}
	}
	return conn, nil
}
func (gp *GrpcPool) GetTotalUse() uint32 {
	var (
		totalUse uint32
		i        uint32
	)
	for ; i < gp.poolCnt.PoolSize; i++ {
		totalUse = totalUse + uint32(gp.conns[i].inflight)
	}
	return totalUse
}

func (gp *GrpcPool) scaleConn(connNum uint32) (uint32, error) {
	//todo reduce idle
	if connNum >= gp.poolCnt.PoolSize {
		return 0, fmt.Errorf("grpc dial target %s error: Setting PoolSize is %d, current conns aleady %d ",
			gp.target, gp.poolCnt.PoolSize, connNum)
	}
	var i uint32 = 0
	tc := make([]*Conn, 0, connNum-gp.poolSize)
	for ; i < connNum-gp.poolSize; i++ {
		ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(gp.deadline))
		defer cancel()
		conn, err := grpc.DialContext(ctx, gp.target, gp.dialOptions...)
		if err != nil {
			return 0, fmt.Errorf("grpc dial target %v error %v", gp.target, err)
		}
		tc = append(tc, &Conn{
			ClientConn: conn,
			inflight:   0,
		})
	}
	gp.conns = combine(tc, gp.conns)
	gp.poolSize = connNum
	return gp.poolSize, nil
}

func (gp *GrpcPool) ScalePool() error {
	log.Infof("grpc pool auto scaler start period %v", gp.poolCnt.IdleCheckFrequency)
	tk := time.NewTicker(gp.poolCnt.IdleCheckFrequency)
	r := func() error {
		gp.rw.Lock()
		totalUse := gp.GetTotalUse()
		if maxStreams := gp.poolCnt.PoolSize * gp.poolCnt.MaxCurrentStream; totalUse > maxStreams {
			deltaConn := (totalUse - maxStreams) / (gp.poolCnt.MaxCurrentStream / 2)
			_, err := gp.scaleConn(deltaConn)
			if err != nil {
				return err
			}
		}
		gp.rw.Unlock()
		return nil
	}
	for {
		select {
		case <-gp.ntfChan:
			tk.Reset(gp.poolCnt.IdleCheckFrequency)
			if err := r(); err != nil {
				return err
			}
		case <-tk.C:
			if err := r(); err != nil {
				return err
			}
		}
	}
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

func WithDialOptions(dp ...grpc.DialOption) Option {
	return func(gp *GrpcPool) {
		gp.dialOptions = append(gp.dialOptions, dp...)
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
		target: target,
		poolCnt: &PoolController{
			PoolSize:           uint32(defaultPoolCardinalSize * runtime.GOMAXPROCS(0)),
			MaxCurrentStream:   uint32(runtime.GOMAXPROCS(0)),
			MinIdleConns:       defaultMinIdleConns,
			MaxConnAge:         defaultMaxConnAge * time.Minute,
			IdleCheckFrequency: defaultIdleCheckFrequency,
		},
		deadline: defaultDialDeadline,
		ntfChan:  make(chan struct{}),
	}
	for _, opt := range opts {
		opt(gp)
	}
	//conn size todo: similarly keepalive programs
	//gp.conns = make([]*Conn, gp.poolCnt.PoolSize)
	//conn timeout
	gp.dialOptions = combine(defaultDialOpts, gp.dialOptions)
	if _, err := gp.scaleConn(defaultPoolInitSize); err != nil {
		return nil, fmt.Errorf("grpc dial target %v error %v", target, err)
	}

	return gp, nil
}

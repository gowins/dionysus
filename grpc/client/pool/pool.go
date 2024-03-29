package pool

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
)

type GrpcPool struct {
	conns       []*GrpcConn
	poolSize    int
	dialOptions []grpc.DialOption
	target      string
	rand        *rand.Rand
	scaleOption *ScaleOption
	sync.Locker
	stateUpdate sync.Locker
	isClosed    bool
}

type GrpcConn struct {
	conn     *grpc.ClientConn
	inflight int64
}

type GrpcPoolState struct {
	ConnStates  []GrpcConnState
	ReserveSize int
	Target      string
	ScaleOption ScaleOption
	IsClosed    bool
}

type GrpcConnState struct {
	connState string
	inflight  int64
}

type ScaleOption struct {
	Enable          bool
	ScalePeriod     time.Duration
	MaxConn         int
	DesireMaxStream int
}

var DefaultDialOpts = []grpc.DialOption{
	grpc.WithTransportCredentials(insecure.NewCredentials()),
	grpc.WithBlock(),
}

func InitGrpcPool(target string, opts ...Option) (*GrpcPool, error) {
	if target == "" {
		return nil, fmt.Errorf("grpc pool target should not be nil")
	}
	gp := &GrpcPool{
		poolSize:    defaultPoolSize,
		dialOptions: DefaultDialOpts,
		target:      target,
		Locker:      new(sync.Mutex),
		stateUpdate: new(sync.Mutex),
		rand:        rand.New(rand.NewSource(time.Now().Unix())),
		scaleOption: &ScaleOption{Enable: false, MaxConn: defaultPoolSize},
	}

	for _, opt := range opts {
		opt(gp)
	}

	if gp.scaleOption.MaxConn < gp.poolSize {
		gp.scaleOption.MaxConn = gp.poolSize
	}

	gp.conns = make([]*GrpcConn, gp.scaleOption.MaxConn)

	for i := 0; i < gp.poolSize; i++ {
		conn, err := grpcDialWithTimeout(gp.target, gp.dialOptions...)
		if err != nil {
			return gp, fmt.Errorf("grpc dial target %v error %v", gp.target, err)
		}
		gp.conns[i] = &GrpcConn{
			conn:     conn,
			inflight: 0,
		}
	}

	if gp.scaleOption.Enable {
		go gp.autoScalerRun()
	}
	return gp, nil
}

func GetGrpcPool(target string, opts ...Option) (*GrpcPool, error) {
	if val, ok := grpcPool.Load(target); ok {
		return val.(*GrpcPool), nil
	}

	poolInit.Lock()
	defer poolInit.Unlock()

	// 双检, 避免多次创建
	if val, ok := grpcPool.Load(target); ok {
		return val.(*GrpcPool), nil
	}

	gp, err := InitGrpcPool(target, opts...)
	if err != nil {
		return nil, err
	}

	grpcPool.Store(target, gp)
	return gp, nil
}

func grpcDialWithTimeout(target string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultDialTimeout)
	defer cancel()

	conn, err := grpc.DialContext(ctx, target, opts...)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func (gp *GrpcPool) Invoke(ctx context.Context, method string, args interface{}, reply interface{}, opts ...grpc.CallOption) error {
	grpcConn := gp.pickLeastConn()
	atomic.AddInt64(&grpcConn.inflight, 1)
	defer atomic.AddInt64(&grpcConn.inflight, -1)
	return grpcConn.conn.Invoke(ctx, method, args, reply, opts...)
}

func (gp *GrpcPool) NewStream(ctx context.Context, desc *grpc.StreamDesc, method string, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	grpcConn := gp.pickLeastConn()
	atomic.AddInt64(&grpcConn.inflight, 1)
	defer atomic.AddInt64(&grpcConn.inflight, -1)
	return grpcConn.conn.NewStream(ctx, desc, method, opts...)
}

func (gp *GrpcPool) pickLeastConn() *GrpcConn {
	gp.Lock()
	randIndex1 := gp.rand.Uint32()
	randIndex2 := gp.rand.Uint32()
	randIndex3 := gp.rand.Uint32()
	gp.Unlock()
	minIndex := randIndex1
	minInflight := gp.conns[int(minIndex)%gp.poolSize].inflight

	if minInflight > gp.conns[int(randIndex2)%gp.poolSize].inflight {
		minInflight = gp.conns[int(randIndex2)%gp.poolSize].inflight
		minIndex = randIndex2
	}

	if minInflight > gp.conns[int(randIndex3)%gp.poolSize].inflight {
		minInflight = gp.conns[int(randIndex3)%gp.poolSize].inflight
		minIndex = randIndex3
	}
	grpcConn := gp.conns[int(minIndex)%gp.poolSize]

	// if conn is not ready, choose a next ready conn
	if grpcConn.conn.GetState() != connectivity.Ready && grpcConn.conn.GetState() != connectivity.Idle {
		for i := 0; i < gp.poolSize; i++ {
			if gp.conns[(int(minIndex)+i)%gp.poolSize].conn.GetState() == connectivity.Ready ||
				gp.conns[(int(minIndex)+i)%gp.poolSize].conn.GetState() == connectivity.Idle {
				return gp.conns[(int(minIndex)+i)%gp.poolSize]
			}
		}
	}
	return grpcConn
}

func (gp *GrpcPool) autoScalerRun() {
	log.Infof("grpc pool auto scaler start period %v", gp.scaleOption.ScalePeriod)
	tk := time.NewTicker(gp.scaleOption.ScalePeriod)
	for {
		select {
		case <-tk.C:
			totalUse := gp.GetTotalUse()
			if totalUse > gp.poolSize*gp.scaleOption.DesireMaxStream {
				deltaConn := (totalUse - gp.poolSize*gp.scaleOption.DesireMaxStream) / (gp.scaleOption.DesireMaxStream / 2)
				gp.poolScaler(deltaConn)
			}
		}
	}
}

func (gp *GrpcPool) poolScaler(deltaConn int) {
	gp.stateUpdate.Lock()
	defer gp.stateUpdate.Unlock()
	if gp.isClosed {
		log.Infof("grpc pool is closed, will not scale")
		return
	}
	if deltaConn+gp.poolSize > gp.scaleOption.MaxConn {
		deltaConn = gp.scaleOption.MaxConn - gp.poolSize
	}

	if deltaConn <= 0 {
		log.Warnf("grpc conn reach max conn, be careful")
		return
	}

	for i := 0; i < deltaConn; i++ {
		conn, err := grpcDialWithTimeout(gp.target, gp.dialOptions...)
		if err != nil {
			log.Infof("grpc pool is scaler form %v to %v", gp.poolSize, gp.poolSize+i)
			gp.poolSize = gp.poolSize + i
			log.Errorf("grpc dial target %v error %v", gp.target, err)
			return
		}
		gp.conns[gp.poolSize+i] = &GrpcConn{
			conn:     conn,
			inflight: 0,
		}
	}
	log.Infof("grpc pool is scaler form %v to %v", gp.poolSize, gp.poolSize+deltaConn)
	gp.poolSize = gp.poolSize + deltaConn
}

func (gp *GrpcPool) GetTotalUse() int {
	var totalUse int
	for i := 0; i < gp.poolSize; i++ {
		totalUse = totalUse + int(gp.conns[i].inflight)
	}
	return totalUse
}

func (gp *GrpcPool) GetGrpcPoolState() *GrpcPoolState {
	connStates := make([]GrpcConnState, gp.poolSize)
	for i := 0; i < gp.poolSize; i++ {
		connStates[i] = GrpcConnState{
			connState: gp.conns[i].conn.GetState().String(),
			inflight:  gp.conns[i].inflight,
		}
	}

	return &GrpcPoolState{
		ConnStates:  connStates,
		ReserveSize: gp.poolSize,
		Target:      gp.target,
		IsClosed:    gp.isClosed,
		ScaleOption: ScaleOption{
			Enable:          gp.scaleOption.Enable,
			ScalePeriod:     gp.scaleOption.ScalePeriod,
			MaxConn:         gp.scaleOption.MaxConn,
			DesireMaxStream: gp.scaleOption.DesireMaxStream,
		},
	}
}

func (gp *GrpcPool) Closed() {
	gp.isClosed = true
	gp.stateUpdate.Lock()
	defer gp.stateUpdate.Unlock()
	for i := 0; i < gp.poolSize; i++ {
		if err := gp.conns[i].conn.Close(); err != nil {
			log.Errorf("grpc conn close error %v", err)
		}
	}
}

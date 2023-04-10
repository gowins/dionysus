package pool

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

const maxConn = 300

var (
	AutoScalePeriod = 30 * time.Second
	AutoScaler      = false
)

type GrpcPool struct {
	conns       []*GrpcConn
	ReserveSize int
	DialOptions []grpc.DialOption
	Target      string
	index       uint32
	rand        *rand.Rand
	autoScaler  bool
	sync.Locker
}

type GrpcConn struct {
	conn     *grpc.ClientConn
	inflight int64
}

func InitGrpcPool(target string, size int, dialOptions ...grpc.DialOption) (*GrpcPool, error) {
	if target == "" {
		return nil, fmt.Errorf("grpc pool target should not be nil")
	}
	if size < 3 {
		return nil, fmt.Errorf("grpc pool reserve size should not be <= 3")
	}
	gp := &GrpcPool{
		conns:       make([]*GrpcConn, maxConn),
		ReserveSize: size,
		DialOptions: dialOptions,
		Target:      target,
		Locker:      new(sync.Mutex),
		rand:        rand.New(rand.NewSource(time.Now().Unix())),
	}
	for i := 0; i < gp.ReserveSize; i++ {
		conn, err := grpc.Dial(gp.Target, gp.DialOptions...)
		if err != nil {
			return gp, fmt.Errorf("grpc dial target %v error %v", gp.Target, err)
		}
		gp.conns[i] = &GrpcConn{
			conn:     conn,
			inflight: 0,
		}
	}
	gp.autoScaler = AutoScaler
	if gp.autoScaler == true {
		go gp.autoScalerRun()
	}
	return gp, nil
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

// TODO make sure grpc conn state
func (gp *GrpcPool) pickLeastConn() *GrpcConn {
	gp.Lock()
	randIndex1 := gp.rand.Uint32()
	randIndex2 := gp.rand.Uint32()
	randIndex3 := gp.rand.Uint32()
	gp.Unlock()
	minIndex := randIndex1
	minInflight := gp.conns[int(minIndex)%gp.ReserveSize].inflight

	if minInflight > gp.conns[int(randIndex2)%gp.ReserveSize].inflight {
		minInflight = gp.conns[int(randIndex2)%gp.ReserveSize].inflight
		minIndex = randIndex2
	}

	if minInflight > gp.conns[int(randIndex3)%gp.ReserveSize].inflight {
		minInflight = gp.conns[int(randIndex3)%gp.ReserveSize].inflight
		minIndex = randIndex3
	}
	return gp.conns[int(minIndex)%gp.ReserveSize]
}

func (gp *GrpcPool) autoScalerRun() {
	tk := time.NewTicker(AutoScalePeriod)
	for {
		select {
		case <-tk.C:
			totalUse := gp.GetTotalUse()
			if totalUse > gp.ReserveSize*80 {
				deltaConn := (totalUse - gp.ReserveSize*80) / 50
				gp.poolScaler(deltaConn)
			}
		}
	}
}

func (gp *GrpcPool) poolScaler(deltaConn int) {
	if deltaConn < 1 {
		log.Errorf("deltaConn is %v, no need pool scaler", deltaConn)
		return
	}
	if deltaConn+gp.ReserveSize > maxConn {
		deltaConn = maxConn - gp.ReserveSize
	}

	for i := 0; i < deltaConn; i++ {
		conn, err := grpc.Dial(gp.Target, gp.DialOptions...)
		if err != nil {
			log.Infof("grpc pool is scaler form %v to %v", gp.ReserveSize, gp.ReserveSize+i)
			gp.ReserveSize = gp.ReserveSize + i
			log.Errorf("grpc dial target %v error %v", gp.Target, err)
			return
		}
		gp.conns[gp.ReserveSize+i] = &GrpcConn{
			conn:     conn,
			inflight: 0,
		}
	}
	log.Infof("grpc pool is scaler form %v to %v", gp.ReserveSize, gp.ReserveSize+deltaConn)
	gp.ReserveSize = gp.ReserveSize + deltaConn
}

func (gp *GrpcPool) GetTotalUse() int {
	var totalUse int
	for i := 0; i < gp.ReserveSize; i++ {
		totalUse = totalUse + int(gp.conns[i].inflight)
	}
	return totalUse
}

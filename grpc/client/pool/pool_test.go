package pool

import (
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"google.golang.org/grpc"
)

func TestInitGrpcPool(t *testing.T) {
	size := 33
	testGrpcPool := &GrpcPool{
		conns:       make([]*GrpcConn, size),
		ReserveSize: size,
		Locker:      new(sync.Mutex),
		rand:        rand.New(rand.NewSource(time.Now().Unix())),
	}
	for i := 0; i < size; i++ {
		testGrpcPool.conns[i] = &GrpcConn{
			conn: &grpc.ClientConn{},
		}
	}

	for i := 0; i < size*1000; i++ {
		gp := testGrpcPool.pickLeastConn()
		atomic.AddInt64(&gp.inflight, 1)
	}

	for i := 0; i < size; i++ {
		fmt.Printf("i: %v, inflight is %v\n", i, testGrpcPool.conns[i].inflight)
	}
}

package pool

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net"
	"testing"
	"time"

	"github.com/gowins/dionysus/grpc/client/pool/testpb"
)

/*
func TestInitGrpcPool(t *testing.T) {
	size := 10
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

	for i := 0; i < size*10000; i++ {
		go func() {
			gp := testGrpcPool.pickLeastConn()
			atomic.AddInt64(&gp.inflight, 1)
			defer atomic.AddInt64(&gp.inflight, -1)
			t := rand.Int() % 2000
			sleepTime := time.Millisecond * time.Duration(t)
			time.Sleep(sleepTime)
		}()
	}

	for {
		fmt.Printf("\n==================\n")
		for i := 0; i < size; i++ {
			fmt.Printf("i: %v, inflight is %v\n", i, testGrpcPool.conns[i].inflight)
		}
		fmt.Printf("\n****************\n")
		time.Sleep(time.Millisecond * 10)
	}
}
*/

type testgserver struct {
}

// SayHello1 implements helloworld.GreeterServer
func (s *testgserver) SayHelloTest1(ctx context.Context, in *testpb.HelloRequest) (*testpb.HelloReply, error) {
	time.Sleep(time.Hour)
	return &testpb.HelloReply{Message: "Hello Test1" + in.GetName()}, nil
}

// SayHello2 implements helloworld.GreeterServer
func (s *testgserver) SayHelloTest2(ctx context.Context, in *testpb.HelloRequest) (*testpb.HelloReply, error) {
	return &testpb.HelloReply{Message: "Hello Test2" + in.GetName()}, nil
}

// SayHello2 implements helloworld.GreeterServer
func (s *testgserver) SayHelloTest3(ctx context.Context, in *testpb.HelloRequest) (*testpb.HelloReply, error) {
	return &testpb.HelloReply{Message: "Hello Test3" + in.GetName()}, nil
}

func setupTestServer(serverDone chan struct{}, addr string) {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	testpb.RegisterGreeterServer(s, &testgserver{})
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()
	<-serverDone
	s.Stop()
}

func TestPoolScaler(t *testing.T) {
	addr := "127.0.0.1:8818"
	serverDone := make(chan struct{})
	go func() {
		setupTestServer(serverDone, addr)
	}()
	gPool, err := InitGrpcPool(addr, 3,
		grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		t.Errorf("grpc pool init dial error %v\n", err)
		return
	}
	c := testpb.NewGreeterClient(gPool)
	r, err := c.SayHelloTest1(context.Background(), &testpb.HelloRequest{Name: "nameing"})
	if err != nil {
		t.Errorf("could not greet: %v", err)
		return
	}
	fmt.Printf("Greeting: %s", r.GetMessage())
	close(serverDone)
	r, err = c.SayHelloTest1(context.Background(), &testpb.HelloRequest{Name: "nameing2"})
	if err == nil {
		t.Errorf("could not greet want error get nil")
		return
	}
}

func TestPoolScalerAdd(t *testing.T) {
	addr := "127.0.0.1:8817"
	serverDone := make(chan struct{})
	go func() {
		setupTestServer(serverDone, addr)
	}()
	AutoScaler = true
	AutoScalePeriod = time.Second * 10
	gPool, err := InitGrpcPool(addr, 30,
		grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		t.Errorf("grpc pool init dial error %v", err)
		return
	}
	go func() {
		for {
			totalUse := 0
			for i := 0; i < gPool.ReserveSize; i++ {
				totalUse = totalUse + int(gPool.conns[i].inflight)
				fmt.Printf("%v conn inflight is %v\n", i, gPool.conns[i].inflight)
			}
			fmt.Printf("reserveSize is %v, totalUse is %v\n", gPool.ReserveSize, totalUse)
			time.Sleep(3 * time.Second)
		}
	}()
	c := testpb.NewGreeterClient(gPool)
	for i := 0; i < 20; i++ {
		for j := 0; j < 3000; j++ {
			go func() {
				r, err := c.SayHelloTest1(context.Background(), &testpb.HelloRequest{Name: "nameing1"})
				if err != nil {
					t.Errorf("could not greet: %v", err)
					return
				}
				fmt.Printf("Greeting: %s\n", r.GetMessage())
			}()
		}
		time.Sleep(10 * time.Second)
	}

	close(serverDone)
}

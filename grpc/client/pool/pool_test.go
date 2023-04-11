package pool

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"net"
	"testing"
	"time"

	"github.com/gowins/dionysus/grpc/client/pool/testpb"
)

type testgserver struct {
}

// SayHello1 implements helloworld.GreeterServer
func (s *testgserver) SayHelloTest1(ctx context.Context, in *testpb.HelloRequest) (*testpb.HelloReply, error) {
	time.Sleep(time.Minute * 5)
	return &testpb.HelloReply{Message: "Hello Test1" + in.GetName()}, nil
}

// SayHello2 implements helloworld.GreeterServer
func (s *testgserver) SayHelloTest2(ctx context.Context, in *testpb.HelloRequest) (*testpb.HelloReply, error) {
	return &testpb.HelloReply{Message: "Hello Test2" + in.GetName()}, nil
}

// SayHello2 implements helloworld.GreeterServer
func (s *testgserver) SayHelloTest3(ctx context.Context, in *testpb.HelloRequest) (*testpb.HelloReply, error) {
	time.Sleep(time.Second * 5)
	return &testpb.HelloReply{Message: "Hello Test3" + in.GetName()}, nil
}

func setupTestServer(serverDone chan struct{}, addr string) {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	testEnforcementPolicy := keepalive.EnforcementPolicy{
		MinTime:             5 * time.Second, // If a client pings more than once every 5 seconds, terminate the connection
		PermitWithoutStream: true,            // Allow pings even when there are no active streams
	}
	s := grpc.NewServer(grpc.KeepaliveEnforcementPolicy(testEnforcementPolicy))
	testpb.RegisterGreeterServer(s, &testgserver{})
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()
	<-serverDone
	s.Stop()
}

func TestPoolScalerMax(t *testing.T) {
	addr := "127.0.0.1:8817"
	serverDone := make(chan struct{})
	defer close(serverDone)
	go func() {
		setupTestServer(serverDone, addr)
	}()
	DefaultScaleOption.ScalePeriod = 5 * time.Second
	gPool, err := GetGrpcPool(addr, WithScaleOption(DefaultScaleOption), WithReserveSize(30))
	if err != nil {
		t.Errorf("grpc pool init dial error %v", err)
		return
	}
	c := testpb.NewGreeterClient(gPool)
	for i := 0; i < 20; i++ {
		for j := 0; j < 2000; j++ {
			go func() {
				c.SayHelloTest1(context.Background(), &testpb.HelloRequest{Name: "nameing1"})
			}()
		}
		time.Sleep(10 * time.Millisecond)
	}

	time.Sleep(20 * time.Second)

	totalUse := gPool.GetTotalUse()
	if totalUse != 20*2000 {
		t.Errorf("want total use %v, get %v", 20*2000, totalUse)
		return
	}

	if gPool.poolSize != gPool.scaleOption.MaxConn {
		t.Errorf("want total poolSize %v equal MaxConn %v", gPool.poolSize, gPool.scaleOption.MaxConn)
		return
	}
}

func TestPoolWithoutScaler(t *testing.T) {
	addr := "127.0.0.1:8818"
	serverDone := make(chan struct{})
	defer close(serverDone)
	go func() {
		setupTestServer(serverDone, addr)
	}()
	gPool, err := GetGrpcPool(addr, WithReserveSize(18), WithDialOptions([]grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                KeepAliveTime,
			Timeout:             KeepAliveTimeout,
			PermitWithoutStream: true,
		}),
	}))
	if err != nil {
		t.Errorf("grpc pool init dial error %v", err)
		return
	}

	if gPool.poolSize != 18 {
		t.Errorf("want poolSize 18 get %v", gPool.poolSize)
		return
	}

	if len(gPool.dialOptions) != 3 {
		t.Errorf("want dialOptions 3 get %v", len(gPool.dialOptions))
	}

	if gPool.scaleOption.Enable {
		t.Errorf("want scale enable not true")
		return
	}

	c := testpb.NewGreeterClient(gPool)
	for i := 0; i < 18; i++ {
		for j := 0; j < 300; j++ {
			go func() {
				c.SayHelloTest1(context.Background(), &testpb.HelloRequest{Name: "nameing1"})
			}()
		}
		time.Sleep(10 * time.Millisecond)
	}

	time.Sleep(40 * time.Second)

	if gPool.poolSize != 18 {
		t.Errorf("want total poolSize 18, get %v", gPool.poolSize)
		return
	}

	for _, conn := range gPool.conns {
		fmt.Printf("conn inflight is %v\n", conn.inflight)
		if conn.inflight < 200 {
			t.Errorf("loadbalancer is not well")
			return
		}
	}
}

func TestPoolInitError(t *testing.T) {
	SetLog(log)
	_, err := GetGrpcPool("")
	if err == nil {
		t.Errorf("want error is not nil, get nil")
		return
	}

	_, err = GetGrpcPool("127.0.0.1:8888")
	if err == nil {
		t.Errorf("want error is not nil, get nil")
		return
	}
}

func TestPoolScaler(t *testing.T) {
	addr := "127.0.0.1:8819"
	serverDone := make(chan struct{})
	defer close(serverDone)
	go func() {
		setupTestServer(serverDone, addr)
	}()
	DefaultScaleOption.ScalePeriod = 5 * time.Second
	for i := 0; i < 30; i++ {
		gPool, err := GetGrpcPool(addr, WithScaleOption(DefaultScaleOption), WithReserveSize(30))
		if err != nil {
			t.Errorf("grpc pool init dial error %v", err)
			return
		}
		c := testpb.NewGreeterClient(gPool)
		for j := 0; j < 2000; j++ {
			go func() {
				rsp, err := c.SayHelloTest3(context.Background(), &testpb.HelloRequest{Name: "nameing1"})
				if err != nil || rsp.Message != "Hello Test3nameing1" {
					t.Errorf("get rsp failed")
					return
				}
			}()
		}
		fmt.Printf("get pool state %v\n", gPool.GetGrpcPoolState())
		time.Sleep(500 * time.Millisecond)
	}

	time.Sleep(25 * time.Second)
}

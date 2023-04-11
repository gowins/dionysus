package pool

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
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
	gPool, err := InitGrpcPool(addr, WithScaleOption(DefaultScaleOption), WithReserveSize(30))
	if err != nil {
		t.Errorf("grpc pool init dial error %v", err)
		return
	}
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
		time.Sleep(10 * time.Millisecond)
	}

	time.Sleep(20 * time.Second)

	totalUse := gPool.GetTotalUse()
	if totalUse != 20*3000 {
		t.Errorf("want total use %v, get %v", 20*3000, totalUse)
		return
	}

	if gPool.reserveSize != gPool.scaleOption.MaxConn {
		t.Errorf("want total reserveSize %v equal MaxConn %v", gPool.reserveSize, gPool.scaleOption.MaxConn)
		return
	}
}

package client

import (
	"context"
	"log"
	"net"
	"testing"
	"time"

	pb "github.com/gowins/dionysus/grpc/client/helloworld"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

// server is used to implement hello-world.GreeterServer.
type server struct {
	pb.UnimplementedGreeterServer
}

// SayHello implements hello-world.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: "Hello " + in.Name}, nil
}
func init() {
	lis = bufconn.Listen(bufSize)
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func TestSayHello(t *testing.T) {
	ctx := context.Background()
	conn, err := Dial("buffet",
		WithDialDeadline(time.Second),
		WithPoolControl(&PoolController{
			PoolSize:           4,
			MaxCurrentStream:   2,
			MinIdleConns:       2,
			MaxConnAge:         time.Second * 15,
			IdleCheckFrequency: time.Second * 30,
		}),
		WithDialOptions(grpc.WithTransportCredentials(insecure.NewCredentials())),
		WithDialOptions(grpc.WithContextDialer(bufDialer)))
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()
	client := pb.NewGreeterClient(conn)
	name := "Mr.Wong"
	resp, err := client.SayHello(ctx, &pb.HelloRequest{Name: name})
	if err != nil {
		t.Fatalf("SayHello failed: %v", err)
	}
	assert.Equal(t, resp.Message, "Hello "+name, "Say hello Fail")
}

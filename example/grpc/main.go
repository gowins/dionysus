package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gowins/dionysus"
	"github.com/gowins/dionysus/cmd"
	"github.com/gowins/dionysus/example/grpc/hw"
	"github.com/gowins/dionysus/grpc/server"
	"github.com/gowins/dionysus/grpc/serverinterceptors"
	"google.golang.org/grpc/metadata"
)

type gserver struct {
	hw.UnimplementedGreeterServer
}

// SayHello implements helloworld.GreeterServer
func (s *gserver) SayHello(ctx context.Context, in *hw.HelloRequest) (*hw.HelloReply, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	fmt.Printf("MetaData: %#v \n", md)
	log.Printf("Received: %v", in.GetName())
	return &hw.HelloReply{Message: "Hello " + in.GetName()}, nil
}

func main() {
	dio := dionysus.NewDio()
	cfg := server.DefaultCfg
	cfg.Address = ":8081"
	c := cmd.NewGrpcCmd(cmd.WithCfg(cfg))
	c.EnableDebug()
	// timeout interceptor
	c.AddUnaryServerInterceptors(serverinterceptors.TimeoutUnary(10 * time.Second))
	// recover interceptor
	c.AddUnaryServerInterceptors(serverinterceptors.RecoveryUnary(serverinterceptors.DefaultRecovery()))
	c.AddStreamServerInterceptors(serverinterceptors.RecoveryStream(serverinterceptors.DefaultRecovery()))
	// tacing interceptor
	c.AddUnaryServerInterceptors(serverinterceptors.OpenTracingUnary())
	c.AddStreamServerInterceptors(serverinterceptors.OpenTracingStream())
	// register grpc service
	c.RegisterGrpcService(hw.RegisterGreeterServer, &gserver{})
	dio.DioStart("grpc", c)
}

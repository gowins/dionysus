package main

import (
	"context"
	"fmt"
	otm "github.com/gowins/dionysus/opentelemetry"
	"github.com/gowins/dionysus/step"
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
	dio.PreRunStepsAppend(step.InstanceStep{
		StepName: "init trace",
		Func: func() error {
			otm.Setup(otm.WithServiceInfo(&otm.ServiceInfo{
				Name:      "testGrpc217",
				Namespace: "testGrpcNamespace217",
				Version:   "testGrpcVersion217",
			}), otm.WithTraceExporter(&otm.Exporter{
				ExporterEndpoint: otm.DefaultStdout,
				Insecure:         false,
				Creds:            otm.DefaultCred,
			}))
			return nil
		},
	})
	cfg := server.DefaultCfg
	cfg.Address = "127.0.0.1:8081"
	c := cmd.NewGrpcCmd(cmd.WithCfg(cfg))
	c.EnableDebug()
	// timeout interceptor
	c.AddUnaryServerInterceptors(serverinterceptors.TimeoutUnary(10 * time.Second))
	// recover interceptor
	c.AddUnaryServerInterceptors(serverinterceptors.RecoveryUnary(serverinterceptors.DefaultRecovery()))
	c.AddStreamServerInterceptors(serverinterceptors.RecoveryStream(serverinterceptors.DefaultRecovery()))
	// register grpc service
	c.RegisterGrpcService(hw.RegisterGreeterServer, &gserver{})
	dio.DioStart("grpc", c)
}

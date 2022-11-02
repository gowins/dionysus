# grpc

## client
1. client是对grpc client的封装, 主要关注于默认的参数的处理和拦截器的加载
2. client默认使用`p2c`的负载均衡器,  同时根据环境变量选择是否从注册中心加载服务, 连接服务
3. 使用
```go
package main

import (
	"context"
	"log"
	"sync"

	"github.com/gowins/dionysus/example/grpc/hw"
	"github.com/gowins/dionysus/grpc/client"
	xlog "github.com/gowins/dionysus/log"
	"google.golang.org/grpc"
)

func main() {
	xlog.Setup(xlog.SetProjectName("grpc-client"))
	var wg sync.WaitGroup
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			conn, err := client.New(":8081")
			if err != nil {
				log.Fatalf("did not connect: %v", err)
			}
			defer conn.Close()
			c := hw.NewGreeterClient(conn)
			// Contact the server and print out its response.
			r, err := c.SayHello(context.Background(), &hw.HelloRequest{Name: "nameing"})
			if err != nil {
				log.Printf("could not greet: %v", err)
				return
			}
			log.Printf("Greeting: %s", r.GetMessage())
		}()
	}
	wg.Wait()
}
```
4. 运行例子

go run example/grpc/client/main.go

5. 熔断
```go
clientinterceptors.BreakerUnary(hystrix.CommandConfig{})
```
* 若传递空的hystrix.CommandConfig，则所有的熔断配置走hystrix默认的
配置，若有自定义配置，则传非空hystrix.CommandConfig
* 可选参数，针对特定的rpc方法有自定义配置的，
```go
ghystrix.HystrixCfg{Name: "/hw.Greeter/SayHello", Cfg: hystrix.CommandConfig{MaxConcurrentRequests: 3}},
```
**Name字段必须是: /服务名/RPC方法名**

## server
1. server 是对grpc server的封装
2. server包含了默认参数处理, 拦截器加载, 服务注册逻辑和启动关闭流程
3. 使用
```go
package main

import (
	"context"
	"log"
	"time"

	"github.com/gowins/dionysus"
	"github.com/gowins/dionysus/cmd"
	"github.com/gowins/dionysus/example/grpc/hw"
	"github.com/gowins/dionysus/grpc/server"
	"github.com/gowins/dionysus/grpc/serverinterceptors"
)

type gserver struct {
	hw.UnimplementedGreeterServer
}

// SayHello implements helloworld.GreeterServer
func (s *gserver) SayHello(ctx context.Context, in *hw.HelloRequest) (*hw.HelloReply, error) {
	log.Printf("Received: %v", in.GetName())
	return &hw.HelloReply{Message: "Hello " + in.GetName()}, nil
}

func main() {
	dio := dionysus.NewDio()
	cfg := server.DefaultCfg
	cfg.Address = ":8081"
	c := cmd.NewGrpcCmd(cmd.WithCfg(cfg))
	c.EnableDebug()
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
```
4. 运行例子

go run example/grpc/main.go grpc

## 健康检查介绍
https://github.com/grpc/grpc/blob/master/doc/health-checking.md

package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"log"
	"sync"
	"time"

	"github.com/gowins/dionysus/example/grpc/hw"
	"github.com/gowins/dionysus/grpc/client/pool"
	xlog "github.com/gowins/dionysus/log"
	"google.golang.org/grpc/metadata"
)

var DefaultMaxRecvMsgSize = 1024 * 1024 * 4

// DefaultMaxSendMsgSize maximum message that client can send
// (4 MB).
var DefaultMaxSendMsgSize = 1024 * 1024 * 4

var clientParameters = keepalive.ClientParameters{
	Time:                10 * time.Second, // send pings every 10 seconds if there is no activity
	Timeout:             2 * time.Second,  // wait 2 second for ping ack before considering the connection dead
	PermitWithoutStream: true,             // send pings even without active streams
}

var defaultDialOpts = []grpc.DialOption{
	grpc.WithTransportCredentials(insecure.NewCredentials()),
	grpc.WithBlock(),
	grpc.WithKeepaliveParams(clientParameters),
	grpc.WithDefaultCallOptions(
		grpc.MaxCallRecvMsgSize(DefaultMaxRecvMsgSize),
		grpc.MaxCallSendMsgSize(DefaultMaxSendMsgSize)),
}

func main() {
	xlog.Setup(xlog.SetProjectName("grpc-client"))
	gPool, err := pool.InitGrpcPool("127.0.0.1:8081", 3, defaultDialOpts...)
	if err != nil {
		fmt.Printf("grpc pool init dial error %v\n", err)
	}
	var wg sync.WaitGroup
	for i := 0; i < 1; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			c := hw.NewGreeterClient(gPool)
			// Contact the server and print out its response.
			mdCtx := metadata.NewOutgoingContext(context.Background(), metadata.New(map[string]string{"k": "v"}))
			r, err := c.SayHello(mdCtx, &hw.HelloRequest{Name: "nameing"})
			if err != nil {
				log.Printf("could not greet: %v", err)
				return
			}
			log.Printf("Greeting: %s", r.GetMessage())
		}()
	}
	wg.Wait()
}

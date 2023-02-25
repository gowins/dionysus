package main

import (
	"context"
	otm "github.com/gowins/dionysus/opentelemetry"
	"log"
	"sync"

	"github.com/gowins/dionysus/example/grpc/hw"
	"github.com/gowins/dionysus/grpc/client"
	xlog "github.com/gowins/dionysus/log"
	"google.golang.org/grpc/metadata"
)

func main() {
	xlog.Setup(xlog.SetProjectName("grpc-client"))
	otm.Setup(otm.WithServiceInfo(&otm.ServiceInfo{
		Name:      "testGrpcCli217",
		Namespace: "testGrpcCliNamespace217",
		Version:   "testGrpcCliVersion217",
	}))
	var wg sync.WaitGroup
	for i := 0; i < 1; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			conn, err := client.NewConnWithTracer("127.0.0.1:8081")
			if err != nil {
				log.Fatalf("did not connect: %v", err)
			}
			defer conn.Close()
			c := hw.NewGreeterClient(conn)
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
	select {}
}

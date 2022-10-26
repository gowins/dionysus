package main

import (
	"context"
	"fmt"
	pb "github.com/gowins/dionysus/healthy/proto"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:1234", grpc.WithInsecure())
	if err != nil {
		fmt.Printf("grpc dial error %v\n", err)
	}
	defer conn.Close()

	client := pb.NewHealthServiceClient(conn)
	reply, err := client.HealthLivenessSet(context.TODO(), &pb.HealthyStatus{})
	if err != nil {
		fmt.Printf("health error %v\n", err)
	}
	fmt.Println(reply.Response)
}

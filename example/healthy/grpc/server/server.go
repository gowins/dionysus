package main

import (
	"context"
	"fmt"
	pb "github.com/gowins/dionysus/healthy/proto"
	"google.golang.org/grpc"
	"net"
)

type GrpcHealthyServer struct {
	livenessStatus  bool
	readinessStatus bool
	startupStatus   bool
}

func (h *GrpcHealthyServer) HealthLiveness(context.Context, *pb.HealthyRequest) (*pb.HealthyResponse, error) {
	if !h.livenessStatus {
		return &pb.HealthyResponse{}, fmt.Errorf("livenessStatus is closed")
	}
	return &pb.HealthyResponse{
		Response: "liveness",
	}, nil
}
func (h *GrpcHealthyServer) HealthLivenessSet(ctx context.Context, req *pb.HealthyStatus) (*pb.HealthyResponse, error) {
	h.livenessStatus = req.Status
	return &pb.HealthyResponse{
		Response: fmt.Sprintf("grpc liveness is set %v", req.Status),
	}, nil
}
func (h *GrpcHealthyServer) HealthReadiness(context.Context, *pb.HealthyRequest) (*pb.HealthyResponse, error) {
	if !h.readinessStatus {
		return &pb.HealthyResponse{}, fmt.Errorf("readinessStatus is closed")
	}
	return &pb.HealthyResponse{
		Response: "readiness",
	}, nil
}
func (h *GrpcHealthyServer) HealthReadinessSet(ctx context.Context, req *pb.HealthyStatus) (*pb.HealthyResponse, error) {
	h.readinessStatus = req.Status
	return &pb.HealthyResponse{
		Response: fmt.Sprintf("grpc readiness is set %v", req.Status),
	}, nil
}
func (h *GrpcHealthyServer) HealthStartup(context.Context, *pb.HealthyRequest) (*pb.HealthyResponse, error) {
	if !h.startupStatus {
		return &pb.HealthyResponse{}, fmt.Errorf("startupStatus is closed")
	}
	return &pb.HealthyResponse{
		Response: "startup",
	}, nil
}
func (h *GrpcHealthyServer) HealthStartupSet(ctx context.Context, req *pb.HealthyStatus) (*pb.HealthyResponse, error) {
	h.startupStatus = req.Status
	return &pb.HealthyResponse{
		Response: fmt.Sprintf("grpc startup is set %v", req.Status),
	}, nil
}

func main() {
	grpcServer := grpc.NewServer()
	pb.RegisterHealthServiceServer(grpcServer, &GrpcHealthyServer{
		livenessStatus:  true,
		readinessStatus: true,
		startupStatus:   true,
	})

	lis, err := net.Listen("tcp", ":1234")
	if err != nil {
		fmt.Printf("listen error %v\n", err)
		return
	}
	err = grpcServer.Serve(lis)
	if err != nil {
		fmt.Printf("grpc Server error %v\n", err)
	}
}

package healthy

import (
	"context"
	"fmt"
	pb "github.com/gowins/dionysus/healthy/proto"
	"google.golang.org/grpc"
	"net"
	"testing"
	"time"
)

func TestGetGrpcHealthyServer(t *testing.T) {
	grpcServer := GetGrpcHealthyServer()
	_, err := grpcServer.HealthLiveness(context.TODO(), &pb.HealthyRequest{})
	if err != nil {
		t.Errorf("want error nil get error %v", err)
		return
	}
	_, err = grpcServer.HealthReadiness(context.TODO(), &pb.HealthyRequest{})
	if err != nil {
		t.Errorf("want error nil get error %v", err)
		return
	}
	_, err = grpcServer.HealthStartup(context.TODO(), &pb.HealthyRequest{})
	if err != nil {
		t.Errorf("want error nil get error %v", err)
		return
	}
	_, err = grpcServer.HealthLivenessSet(context.TODO(), &pb.HealthyStatus{Status: false})
	if err != nil {
		t.Errorf("want error nil get error %v", err)
		return
	}
	_, err = grpcServer.HealthReadinessSet(context.TODO(), &pb.HealthyStatus{Status: false})
	if err != nil {
		t.Errorf("want error nil get error %v", err)
		return
	}
	_, err = grpcServer.HealthStartupSet(context.TODO(), &pb.HealthyStatus{Status: false})
	if err != nil {
		t.Errorf("want error nil get error %v", err)
		return
	}
	_, err = grpcServer.HealthLiveness(context.TODO(), &pb.HealthyRequest{})
	if err == nil {
		t.Errorf("want error not nil")
		return
	}
	_, err = grpcServer.HealthReadiness(context.TODO(), &pb.HealthyRequest{})
	if err == nil {
		t.Errorf("want error not nil")
		return
	}
	_, err = grpcServer.HealthStartup(context.TODO(), &pb.HealthyRequest{})
	if err == nil {
		t.Errorf("want error not nil")
		return
	}
}

func TestGrpcServer(t *testing.T) {
	grpcAddr := ":1234"
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
	go func() {
		err = grpcServer.Serve(lis)
		if err != nil {
			t.Errorf("grpc Server error %v\n", err)
			return
		}
		time.Sleep(8 * time.Second)
	}()
	time.Sleep(2 * time.Second)
	timeoutCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	conn, err := grpc.DialContext(timeoutCtx, "127.0.0.1:1234", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		t.Errorf("want error nil get error %v", err)
		return
	}
	defer conn.Close()
	c := pb.NewHealthServiceClient(conn)

	if err = CheckGrpcHealthy(grpcAddr, HealthLiveness); err != nil {
		t.Errorf("want error nil get error %v", err)
		return
	}
	err = CheckGrpcHealthy(grpcAddr, HealthReadiness)
	if err != nil {
		t.Errorf("want error nil get error %v", err)
		return
	}
	err = CheckGrpcHealthy(grpcAddr, HealthStartup)
	if err != nil {
		t.Errorf("want error nil get error %v", err)
		return
	}

	err = SetGrpcHealthyOpen(grpcAddr, HealthLiveness)
	if err != nil {
		t.Errorf("want error nil get error %v", err)
		return
	}
	err = SetGrpcHealthyOpen(grpcAddr, HealthReadiness)
	if err != nil {
		t.Errorf("want error nil get error %v", err)
		return
	}
	err = SetGrpcHealthyOpen(grpcAddr, HealthStartup)
	if err != nil {
		t.Errorf("want error nil get error %v", err)
		return
	}
	err = checkGrpcLiveness(c)
	if err != nil {
		t.Errorf("want error nil get error %v", err)
		return
	}
	err = checkGrpcReadiness(c)
	if err != nil {
		t.Errorf("want error nil get error %v", err)
		return
	}
	err = checkGrpcStartup(c)
	if err != nil {
		t.Errorf("want error nil get error %v", err)
		return
	}

	err = SetGrpcHealthyClose(grpcAddr, HealthLiveness)
	if err != nil {
		t.Errorf("want error nil get error %v", err)
		return
	}
	err = SetGrpcHealthyClose(grpcAddr, HealthReadiness)
	if err != nil {
		t.Errorf("want error nil get error %v", err)
		return
	}
	err = SetGrpcHealthyClose(grpcAddr, HealthStartup)
	if err != nil {
		t.Errorf("want error nil get error %v", err)
		return
	}

	err = checkGrpcLiveness(c)
	if err == nil {
		t.Errorf("want error not nil")
		return
	}
	err = checkGrpcReadiness(c)
	if err == nil {
		t.Errorf("want error not nil")
		return
	}
	err = checkGrpcStartup(c)
	if err == nil {
		t.Errorf("want error not nil")
		return
	}
	grpcServer.Stop()
}

func TestCheckHttpHealthyStat(t *testing.T) {
	err := CheckHttpHealthyStat("errurl", HealthLiveness)
	if err == nil {
		t.Errorf("want error not nil")
		return
	}
	err = CheckHttpHealthyStat("errurl", HealthReadiness)
	if err == nil {
		t.Errorf("want error not nil")
		return
	}
	err = CheckHttpHealthyStat("errurl", HealthStartup)
	if err == nil {
		t.Errorf("want error not nil")
		return
	}
}

func TestHealth_FileObserve(t *testing.T) {
	health := New()
	RegReadinessCheckers(func() error {
		return nil
	})
	RegLivenessCheckers(func() error {
		return fmt.Errorf("liveness check failed")
	})
	RegStartupCheckers(func() error {
		return nil
	})
	err := health.FileObserve(time.Second * 10)
	if err != nil {
		t.Errorf("want error nil get error %v", err)
	}
	err = health.Stat()
	if err == nil {
		t.Errorf("want error not nil")
	}
	err = CheckCtlHealthyStat(HealthLiveness)
	if err == nil {
		t.Errorf("want error not nil")
	}
}

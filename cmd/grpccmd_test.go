package cmd

import (
	"bytes"
	"io"
	"testing"
	"time"

	"github.com/gowins/dionysus/grpc/server"
	"github.com/gowins/dionysus/grpc/serverinterceptors"
	xlog "github.com/gowins/dionysus/log"
	grpcg "google.golang.org/grpc"
)

//go:norace
func TestGrpcCmd(t *testing.T) {
	xlog.Setup(xlog.SetProjectName("test"), xlog.WithWriter(io.Discard))
	c := NewGrpcCmd(WithCfg(server.DefaultCfg))
	c.EnableDebug()
	c.opts.debug = true
	// recover interceptor
	c.AddUnaryServerInterceptors(serverinterceptors.RecoveryUnary(serverinterceptors.DefaultRecovery()))
	c.AddStreamServerInterceptors(serverinterceptors.RecoveryStream(serverinterceptors.DefaultRecovery()))
	// tacing interceptor
	c.AddUnaryServerInterceptors(serverinterceptors.OpenTracingUnary())
	c.AddStreamServerInterceptors(serverinterceptors.OpenTracingStream())
	// register grpc service
	c.RegisterGrpcService(mockService, &mockerServer{})
	c.RegShutdownFunc(StopStep{StopFn: func() {}, StepName: "test"})
	co := c.GetCmd()
	if co == nil {
		t.Error("GetCmd return nil")
	}
	buffers := new(bytes.Buffer)
	co.SetOutput(buffers)
	go func() {
		co.Execute()
	}()
	time.Sleep(500 * time.Millisecond)
	c.RegisterGrpcService(mockService, &mockerServer{})
	c.GetShutdownFunc()()
}

var _ mocker = (*mockerServer)(nil)

type mocker interface {
	Mocker()
}

type mockerServer struct{}

func (m *mockerServer) Mocker() {}

func mockService(s grpcg.ServiceRegistrar, m mocker) {}

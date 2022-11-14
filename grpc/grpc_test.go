package grpc

import (
	"testing"
	"time"

	"github.com/gowins/dionysus/grpc/client"
	"github.com/gowins/dionysus/grpc/clientinterceptors"
	"github.com/gowins/dionysus/log"
	"github.com/smartystreets/goconvey/convey"
	ggrpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestClientStreamInterceptor(t *testing.T) {
	convey.Convey("client stream interceptor", t, func() {
		convey.So(func() {
			WithDefaultClientStreamInterceptors(clientinterceptors.RecoveryStream(clientinterceptors.DefaultRecovery()))
		}, convey.ShouldNotPanic)
	})
}

func TestClientUnaryInterceptor(t *testing.T) {
	convey.Convey("client unary interceptor", t, func() {
		convey.So(func() {
			WithDefaultClientUnaryInterceptors(clientinterceptors.RecoveryUnary(clientinterceptors.DefaultRecovery()))
		}, convey.ShouldNotPanic)
	})
}

//go:norace
func TestClient(t *testing.T) {
	convey.Convey("client", t, func() {
		convey.Convey("client dial failure", func() {
			client.DefaultTimeout = time.Millisecond * 500
			_, err := NewClient("", ggrpc.WithTransportCredentials(insecure.NewCredentials()))
			convey.So(err, convey.ShouldNotBeNil)
			_, err = GetClient("", ggrpc.WithTransportCredentials(insecure.NewCredentials()))
			convey.So(err, convey.ShouldNotBeNil)
		})
		s := NewServer()
		s.Init()
		go func() {
			s.Start()
		}()
		time.Sleep(time.Millisecond * 500)
		s.Stop()
	})
}

//go:norace
func TestGrpcServerStop(t *testing.T) {
	convey.Convey("grpc", t, func() {
		s := NewServer()
		err := s.Init()
		convey.So(err, convey.ShouldBeNil)
		go func() {
			s.Stop()
		}()
		s.Start()
	})
}

//go:norace
func TestSetLog(t *testing.T) {
	convey.Convey("set log", t, func() {
		log.Setup(log.SetProjectName("test"))
		SetLog(log.GetLogger())
	})
}

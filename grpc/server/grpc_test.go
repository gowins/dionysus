package server

import (
	"crypto/tls"
	"testing"
	"time"

	"github.com/gowins/dionysus/grpc/registry"
	"github.com/gowins/dionysus/grpc/serverinterceptors"
	"github.com/smartystreets/goconvey/convey"
	"google.golang.org/grpc"
)

var _ registry.Registry = (*mockRegistry)(nil)

type mockRegistry struct{}

func (m *mockRegistry) Init(opts ...registry.Option) error {
	return nil
}
func (m *mockRegistry) Register(*registry.Service, ...registry.RegisterOption) error {
	return nil
}
func (m *mockRegistry) Deregister(*registry.Service) error {
	return nil
}
func (m *mockRegistry) GetService(string) ([]*registry.Service, error) {
	return nil, nil
}
func (m *mockRegistry) ListServices() ([]*registry.Service, error) {
	return nil, nil
}
func (m *mockRegistry) Watch(...registry.WatchOption) (registry.Watcher, error) {
	return nil, nil
}
func (m *mockRegistry) String() string {
	return "mock"
}

func TestNewGrpcSever(t *testing.T) {
	convey.Convey("new grpc server", t, func() {
		convey.So(New(DefaultCfg), convey.ShouldNotBeNil)
	})
}

type mer interface {
	Mer()
}

//go:norace
func TestGrpcServer(t *testing.T) {
	convey.Convey("new grpc server", t, func() {
		cfg := DefaultCfg
		cfg.MaxMsgSize = 0
		cfg.TlsCfg = &tls.Config{}
		server := New(cfg)
		convey.So(server, convey.ShouldNotBeNil)
		server.AddUnaryServerInterceptors(serverinterceptors.RecoveryUnary(serverinterceptors.DefaultRecovery()))
		server.AddStreamServerInterceptors(serverinterceptors.RecoveryStream(serverinterceptors.DefaultRecovery()))
		convey.So(server.RegisterService(nil, 1), convey.ShouldNotBeNil)
		convey.So(server.RegisterService(1, nil), convey.ShouldNotBeNil)
		convey.So(server.RegisterService(func() {}, 1), convey.ShouldNotBeNil)
		convey.So(server.RegisterService(func(int, int) {}, 1), convey.ShouldNotBeNil)
		convey.So(server.RegisterService(func(*grpc.Server, mer) {}, &mockRegistry{}), convey.ShouldNotBeNil)
		convey.So(server.RegisterService(func(*grpc.Server, registry.Registry) {}, &mockRegistry{}), convey.ShouldBeNil)
		err := server.Init()
		convey.So(err, convey.ShouldBeNil)
		server.EnableDebug()
		go func() {
			server.Cfg.Address = ""
			_ = server.Start()
			server.Cfg.Address = "127.0.0.1:1"
			_ = server.Start()
			server.Cfg.Address = "127.0.0.1:"
			_ = server.Start()
		}()
		time.Sleep(time.Millisecond * 600)
		server.Stop()
		s := server
		go func() {
			s.Cfg.Address = "127.0.0.1:"
			err = s.Start()
		}()
		time.Sleep(time.Millisecond * 600)
		s.Stop()
	})
}

package server

import (
	"fmt"
	"net"
	"reflect"

	otm "github.com/gowins/dionysus/opentelemetry"
	"github.com/gowins/dionysus/recovery"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/channelz/service"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

type grpcServer struct {
	Cfg                Cfg
	srv                *grpc.Server
	handlers           []func()
	unaryInterceptors  []grpc.UnaryServerInterceptor
	streamInterceptors []grpc.StreamServerInterceptor
}

// EnableDebug
// https://github.com/grpc/grpc-experiments/tree/master/gdebug
func (g *grpcServer) EnableDebug() {
	reflection.Register(g.srv)
	service.RegisterChannelzServiceToServer(g.srv)
}

func (g *grpcServer) GetDefaultServerOpts() []grpc.ServerOption {
	return []grpc.ServerOption{
		grpc.MaxRecvMsgSize(g.Cfg.MaxMsgSize),
		grpc.MaxSendMsgSize(g.Cfg.MaxMsgSize),
		grpc.KeepaliveEnforcementPolicy(DefaultEnforcementPolicy),
		grpc.KeepaliveParams(DefaultServerParameters),
	}
}

// RegisterService 注册grpc handler
func (g *grpcServer) RegisterService(register interface{}, handler interface{}) (err error) {
	defer recovery.CheckErr(&err)

	if register == nil {
		return fmt.Errorf("[register] should not be nil")
	}

	if handler == nil {
		return fmt.Errorf("[handler] should not be nil")
	}

	tRegister := reflect.TypeOf(register)
	tHandler := reflect.TypeOf(handler)

	if tRegister.NumIn() != 2 {
		return fmt.Errorf("[register] input num not match")
	}

	// if tRegister.In(0).String() != "*grpc.Server" {
	registrar := new(grpc.ServiceRegistrar)
	if !tRegister.In(0).Implements(reflect.TypeOf(registrar).Elem()) {
		return fmt.Errorf("[register] input type error")
	}

	if !tHandler.Implements(tRegister.In(1)) {
		return fmt.Errorf("%s not implements interface %s", tHandler.String(), tRegister.In(1).String())
	}

	// 注册handler
	g.handlers = append(g.handlers, func() {
		reflect.ValueOf(register).Call([]reflect.Value{reflect.ValueOf(g.srv), reflect.ValueOf(handler)})
	})
	return nil
}

func (g *grpcServer) AddUnaryServerInterceptors(interceptors ...grpc.UnaryServerInterceptor) {
	g.unaryInterceptors = append(g.unaryInterceptors, interceptors...)
}

func (g *grpcServer) AddStreamServerInterceptors(interceptors ...grpc.StreamServerInterceptor) {
	g.streamInterceptors = append(g.streamInterceptors, interceptors...)
}

func (g *grpcServer) Init(opts ...grpc.ServerOption) (gErr error) {
	defer recovery.CheckErr(&gErr)

	opts1 := g.GetDefaultServerOpts()
	if g.Cfg.TlsCfg != nil {
		opts1 = append(opts1, grpc.Creds(credentials.NewTLS(g.Cfg.TlsCfg)))
	}

	if otm.TracerIsEnable() {
		log.Infof("[Dio] grpc use opentelemetry trace")
		g.unaryInterceptors = append([]grpc.UnaryServerInterceptor{otm.GrpcUnaryTraceInterceptor()}, g.unaryInterceptors...)
		g.streamInterceptors = append([]grpc.StreamServerInterceptor{otm.GrpcStreamTraceInterceptor()}, g.streamInterceptors...)
	}

	opts1 = append(opts1, grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(g.unaryInterceptors...)),
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(g.streamInterceptors...)))

	g.srv = grpc.NewServer(append(opts1, opts...)...)

	// 注册handle
	for i := range g.handlers {
		g.handlers[i]()
	}

	return
}

func (g *grpcServer) Start() (gErr error) {
	defer recovery.CheckErr(&gErr)

	if g.Cfg.Address == "" {
		return fmt.Errorf("[grpc] please set address")
	}

	ts, err := net.Listen("tcp", g.Cfg.Address)
	if err != nil {
		return errors.Wrapf(err, "net Listen error, addr:%s", g.Cfg.Address)
	}

	log.Infof("Server [grpc] Listening on %s", ts.Addr().String())
	g.Cfg.Address = ts.Addr().String()

	if err := g.srv.Serve(ts); err != nil {
		log.Errorf("[grpc] server stop error: %#v", err)
		return err
	}
	return nil
}

func (g *grpcServer) Stop() error {
	g.srv.GracefulStop()
	return nil
}

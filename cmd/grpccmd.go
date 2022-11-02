package cmd

import (
	"log"
	"os"

	"github.com/gowins/dionysus/grpc/server"
	"github.com/gowins/dionysus/healthy"
	pb "github.com/gowins/dionysus/healthy/proto"
	xlog "github.com/gowins/dionysus/log"
	"github.com/gowins/dionysus/runenv"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

const (
	GrpcUse = "grpc"
)

var (
	defaultGrpcAddr = ":8081"
	DefaultGrpcUse  = "grpc"
	DefaultShort    = "Run grpc server"
)

type grpcCmd struct {
	cobraCmd      *cobra.Command
	opts          *GrpcOptions
	serve         *server.GrpcServer
	shutdownSteps []StopStep
}

// GrpcOptFunc set GrpcOptions fileds
type GrpcOptFunc func(*GrpcOptions)

// GrpcOptions grpc options
type GrpcOptions struct {
	Cfg        server.Cfg
	ServerOpts []grpc.ServerOption
	debug      bool
}

// WithCfg set Cfg filed
func WithCfg(cfg server.Cfg) GrpcOptFunc {
	return func(opts *GrpcOptions) {
		opts.Cfg = cfg
	}
}

// WithGrpcServerOpt set ServerOpts filed
func WithGrpcServerOpt(serverOpts ...grpc.ServerOption) GrpcOptFunc {
	return func(opts *GrpcOptions) {
		opts.ServerOpts = serverOpts
	}
}

// NewGrpcCmd new grpcCmd
func NewGrpcCmd(opts ...GrpcOptFunc) *grpcCmd {
	c := &cobra.Command{
		Use:   GrpcUse,
		Short: DefaultShort,
	}
	c.SetOut(os.Stdout)
	c.SetErr(os.Stderr)
	grpcOpts := &GrpcOptions{}
	for _, opt := range opts {
		opt(grpcOpts)
	}
	return &grpcCmd{
		cobraCmd: c,
		opts:     grpcOpts,
		serve:    server.New(grpcOpts.Cfg),
	}
}

// AddUnaryServerInterceptors add unary server interceptors
func (c *grpcCmd) AddUnaryServerInterceptors(interceptors ...grpc.UnaryServerInterceptor) {
	c.serve.AddUnaryServerInterceptors(interceptors...)
}

// AddStreamServerInterceptors add stream server interceptors
func (c *grpcCmd) AddStreamServerInterceptors(interceptors ...grpc.StreamServerInterceptor) {
	c.serve.AddStreamServerInterceptors(interceptors...)
}

// GetCmd init cobra.Command and return it
func (c *grpcCmd) GetCmd() *cobra.Command {
	c.cobraCmd.PreRunE = func(_ *cobra.Command, _ []string) error {
		return nil
	}

	c.cobraCmd.RunE = func(_ *cobra.Command, _ []string) error {
		c.regHealthCheck()
		return c.start()
	}
	return c.cobraCmd
}

func (g *grpcCmd) RegShutdownFunc(stopSteps ...StopStep) {
	g.shutdownSteps = append(g.shutdownSteps, stopSteps...)
}

func (g *grpcCmd) GetShutdownFunc() StopFunc {
	return func() {
		for _, stopSteps := range g.shutdownSteps {
			xlog.Infof("run shutdown %v", stopSteps.StepName)
			stopSteps.StopFn()
		}
		g.serve.Stop()
	}
}
func (c *grpcCmd) start() error {
	if err := c.serve.Init(c.opts.ServerOpts...); err != nil {
		return err
	}
	if c.opts.debug {
		c.serve.EnableDebug()
	}
	return c.serve.Start()
}

func (c *grpcCmd) regHealthCheck() {
	c.RegisterGrpcService(pb.RegisterHealthServiceServer, healthy.GetGrpcHealthyServer())
}

// RegisterGrpcService register a grpc service
func (c *grpcCmd) RegisterGrpcService(fn any, implementor any) {
	if err := c.serve.RegisterService(fn, implementor); err != nil {
		log.Fatalf("[Dio] grpc RegisterService error, err: %v", err)
	}
}

// EnableDebug enable debug mode
func (c *grpcCmd) EnableDebug() {
	c.opts.debug = runenv.IsDev() || runenv.IsTest()
	grpc.EnableTracing = c.opts.debug
}

package server

import (
	"crypto/tls"
	"os"
	"time"

	"google.golang.org/grpc/keepalive"
)

type Cfg struct {
	MaxMsgSize int
	TlsCfg     *tls.Config
	Address    string
}

const (
	// DefaultMaxMsgSize define maximum message size that server can send
	// or receive.  Default value is 4MB.
	DefaultMaxMsgSize = 1024 * 1024 * 4
	// DefaultAddress default grpc server address
	DefaultAddress = ":8081"
)

var (
	DefaultServerParameters = keepalive.ServerParameters{
		MaxConnectionIdle:     30 * time.Second, // If a client is idle for 30 seconds, send a GOAWAY
		MaxConnectionAge:      55 * time.Second, // If any connection is alive for more than 55 seconds, send a GOAWAY
		MaxConnectionAgeGrace: 5 * time.Second,  // Allow 5 seconds for pending RPCs to complete before forcibly closing connections
		Time:                  10 * time.Second, // Ping the client if it is idle for 10 seconds to ensure the connection is still active
		Timeout:               2 * time.Second,  // Wait 2 second for the ping ack before assuming the connection is dead
	}
	DefaultEnforcementPolicy = keepalive.EnforcementPolicy{
		MinTime:             5 * time.Second, // If a client pings more than once every 5 seconds, terminate the connection
		PermitWithoutStream: true,            // Allow pings even when there are no active streams
	}
	DefaultCfg = Cfg{
		MaxMsgSize: DefaultMaxMsgSize,
	}
)

func New(cfgs ...Cfg) *GrpcServer {
	cfg := DefaultCfg
	if len(cfgs) > 0 {
		cfg = cfgs[0]
	}
	if cfg.MaxMsgSize == 0 {
		cfg.MaxMsgSize = DefaultMaxMsgSize
	}
	if cfg.Address == "" {
		cfg.Address = DefaultAddress
	}

	s := &grpcServer{
		Cfg: cfg,
	}

	return &GrpcServer{s}
}

type GrpcServer struct {
	*grpcServer
}

func getHostname() string {
	if name, err := os.Hostname(); err != nil {
		return "unknown"
	} else {
		return name
	}
}

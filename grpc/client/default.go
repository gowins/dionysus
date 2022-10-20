package client

import (
	"fmt"
	"sync"
	"time"

	"github.com/douyu/jupiter/pkg/client/grpc/balancer/p2c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

var mu sync.Mutex
var clients sync.Map

// DefaultTimeout 默认的连接超时时间
var DefaultTimeout = 3 * time.Second

// DefaultMaxRecvMsgSize maximum message that client can receive
// (4 MB).
var DefaultMaxRecvMsgSize = 1024 * 1024 * 4

// DefaultMaxSendMsgSize maximum message that client can send
// (4 MB).
var DefaultMaxSendMsgSize = 1024 * 1024 * 4

var clientParameters = keepalive.ClientParameters{
	Time:                10 * time.Second, // send pings every 10 seconds if there is no activity
	Timeout:             2 * time.Second,  // wait 2 second for ping ack before considering the connection dead
	PermitWithoutStream: true,             // send pings even without active streams
}

var defaultDialOpts = []grpc.DialOption{
	grpc.WithTransportCredentials(insecure.NewCredentials()),
	grpc.WithBlock(),
	grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, p2c.Name)), //nolint:staticcheck
	grpc.WithKeepaliveParams(clientParameters),
	grpc.WithDefaultCallOptions(
		grpc.MaxCallRecvMsgSize(DefaultMaxRecvMsgSize),
		grpc.MaxCallSendMsgSize(DefaultMaxSendMsgSize)),
}

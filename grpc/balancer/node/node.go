package node

import (
	"google.golang.org/grpc/balancer"
)

const Key = "node"

type Node struct {
	Conn      balancer.SubConn
	LoadCount int64
}

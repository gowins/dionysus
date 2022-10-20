package grpc

import (
	"google.golang.org/grpc"
)

type EndPoints []string
type ServerOption = grpc.ServerOption

// type Server = server.GrpcServer
type UnaryServerInterceptor = grpc.UnaryServerInterceptor
type StreamServerInterceptor = grpc.StreamServerInterceptor
type EmptyServerOption = grpc.EmptyServerOption
type ClientConn = grpc.ClientConn
type UnaryInvoker = grpc.UnaryInvoker
type CallOption = grpc.CallOption
type StreamDesc = grpc.StreamDesc
type Streamer = grpc.Streamer
type ClientStream = grpc.ClientStream

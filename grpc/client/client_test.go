package client

import (
	"net"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
)

func TestNew(t *testing.T) {
	_, err := Get(uuid.New().String())
	assert.NotNil(t, err)

	_, err = Get("test")
	assert.NotNil(t, err)

	ts, err := net.Listen("tcp", ":0")
	assert.Nil(t, err)
	srv := grpc.NewServer()
	defer srv.Stop()
	go srv.Serve(ts)

	addr := strings.ReplaceAll(ts.Addr().String(), "[::]", "localhost")
	conn, err1 := Get(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	assert.Nil(t, err1)
	assert.NotNil(t, conn)

	conn1, err2 := Get(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	assert.Nil(t, err2)

	// 相同name获取的对象是一样的
	assert.Equal(t, conn, conn1)

	assert.Equal(t, conn.GetState(), connectivity.Ready)
	assert.Nil(t, conn.Close())
	assert.Equal(t, conn.GetState(), connectivity.Shutdown)
	assert.Equal(t, conn1.GetState(), connectivity.Shutdown)
}

func TestNewClient(t *testing.T) {
	ts, err := net.Listen("tcp", ":0")
	assert.Nil(t, err)
	srv := grpc.NewServer()
	defer srv.Stop()
	go srv.Serve(ts)

	addr := strings.ReplaceAll(ts.Addr().String(), "[::]", "localhost")
	conn, err1 := New(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	assert.Nil(t, err1)
	assert.NotNil(t, conn)

	conn1, err2 := New(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	assert.Nil(t, err2)
	// new出来的是不同的对象, 新的对象
	assert.NotEqual(t, conn, conn1)

	assert.Equal(t, conn.GetState(), connectivity.Ready)
	assert.Nil(t, conn.Close())
	assert.Equal(t, conn.GetState(), connectivity.Shutdown)
	assert.Equal(t, conn1.GetState(), connectivity.Ready)
	assert.Nil(t, conn1.Close())
	assert.Equal(t, conn1.GetState(), connectivity.Shutdown)
}

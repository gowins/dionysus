package resolver

import (
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/resolver"

	"github.com/gowins/dionysus/grpc/registry"
	"github.com/gowins/dionysus/grpc/registry/etcdv3"
)

const rawUrl = "etcdv3://etcd.sf:2379"

var discover bool

func TestMain(m *testing.M) {
	_, discover = os.LookupEnv("GO_TEST_DISCOVER")
	os.Exit(m.Run())
}

func TestDiscovBuilder(t *testing.T) {
	if !discover {
		return
	}
	var name = "service_name"

	var b discovBuilder

	// 未初始化registry
	cc := new(mockedClientConn)
	res, err := b.Build(resolver.Target{URL: url.URL{Path: "service_name", Host: "sf.etcd", Scheme: "DiscovScheme"}}, cc, resolver.BuildOptions{})
	assert.NotNil(t, err)
	assert.Nil(t, res)

	assert.Nil(t, registry.Init(rawUrl))
	services, err := registry.Default.GetService(name)

	if err != nil {
		if !strings.Contains(err.Error(), "not found") {
			assert.NotNil(t, err)
		}
	}

	// 清空所有的服务
	for i := range services {
		assert.Nil(t, registry.Default.Deregister(services[i]))
	}

	// 初始化registry, 服务不存在
	cc = new(mockedClientConn)
	res, err = b.Build(resolver.Target{URL: url.URL{Path: name, Host: etcdv3.Name, Scheme: DiscovScheme}}, cc, resolver.BuildOptions{})
	assert.NotNil(t, err)
	assert.Nil(t, res)
	assert.NotNil(t, cc.state)
	assert.Equal(t, len(cc.state.Addresses), 0)

	// 注册服务信息
	assert.Nil(t, registry.Default.Register(&registry.Service{
		Name: name,
		Nodes: []*registry.Node{{
			Id:      uuid.New().String(),
			Address: "localhost",
			Port:    8087,
		}},
	}))

	// 检测服务的数量, 当前应该是一个
	services, err = registry.Default.GetService(name)
	assert.Nil(t, err)
	assert.Equal(t, len(services), 1)
	assert.Equal(t, len(services[0].Nodes), 1)

	cc = new(mockedClientConn)
	res, err = b.Build(resolver.Target{URL: url.URL{Path: name, Host: etcdv3.Name, Scheme: DiscovScheme}}, cc, resolver.BuildOptions{})
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.NotNil(t, cc.state)
	assert.Equal(t, len(cc.state.Addresses), 1*Replica)

	// 更新服务
	assert.Nil(t, registry.Default.Register(&registry.Service{
		Name: name,
		Nodes: []*registry.Node{{
			Id:      uuid.New().String(),
			Address: "localhost",
			Port:    8087,
		}},
	}))

	// 更新服务, 当前服务数量应该是2
	services, err = registry.Default.GetService(name)
	assert.Nil(t, err)
	assert.Equal(t, len(services), 1)
	assert.Equal(t, len(services[0].Nodes), 2)

	assert.NotNil(t, cc.state)
	// 服务自动更新
	// 当前的子链接地址应该是2*Replica个
	assert.Equal(t, len(cc.state.Addresses), 2*Replica)
}

func TestGetAddrs(t *testing.T) {
	if !discover {
		return
	}
	srvFoo := "service_foo"
	srvBar := "service_bar"
	// 初始化注册中心
	assert.Nil(t, registry.Init(rawUrl))
	// 清空 注册中心 srv
	clearRegSrv(t, srvFoo)
	clearRegSrv(t, srvBar)

	// 注册srvFoo
	assert.Nil(t, registry.Default.Register(&registry.Service{
		Name: srvFoo,
		Nodes: []*registry.Node{{
			Id:      uuid.New().String(),
			Address: "localhost",
			Port:    8087,
		}},
	}))

	//// 注册srvBar
	assert.Nil(t, registry.Default.Register(&registry.Service{
		Name: srvBar,
		Nodes: []*registry.Node{{
			Id:      uuid.New().String(),
			Address: "localhost",
			Port:    8088,
		}},
	}))

	var b discovBuilder

	ccFoo := new(mockedClientConn)

	res, err := b.Build(resolver.Target{URL: url.URL{Path: srvFoo, Host: etcdv3.Name, Scheme: DiscovScheme}}, ccFoo, resolver.BuildOptions{})
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.NotNil(t, ccFoo.state)

	for _, address := range ccFoo.state.Addresses {
		assert.Equal(t, address.ServerName, srvFoo, "should only have srvFoo ClientConn")
	}

	ccBar := new(mockedClientConn)
	res, err = b.Build(resolver.Target{URL: url.URL{Path: srvBar, Host: etcdv3.Name, Scheme: DiscovScheme}}, ccBar, resolver.BuildOptions{})
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.NotNil(t, ccBar.state)
	for _, address := range ccBar.state.Addresses {
		assert.Equal(t, address.ServerName, srvBar, "should only have srvBar ClientConn")
	}

}

// clearRegSrv 清空所有服务
func clearRegSrv(t *testing.T, srvFoo string) {
	services, err := registry.Default.GetService(srvFoo)

	if err != nil {
		if !strings.Contains(err.Error(), "not found") {
			assert.NotNil(t, err)
		}
	}

	// 清空所有的服务
	for i := range services {
		assert.Nil(t, registry.Default.Deregister(services[i]))
	}
}

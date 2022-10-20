package registry_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gowins/dionysus/grpc/registry"
	_ "github.com/gowins/dionysus/grpc/registry/etcdv3"
)

type mockRegistry struct{ name string }

func (m mockRegistry) Init(opts ...registry.Option) error {
	panic("implement me")
}

func (m mockRegistry) Register(service *registry.Service, option ...registry.RegisterOption) error {
	panic("implement me")
}

func (m mockRegistry) Deregister(service *registry.Service) error {
	panic("implement me")
}

func (m mockRegistry) GetService(s string) ([]*registry.Service, error) {
	panic("implement me")
}

func (m mockRegistry) ListServices() ([]*registry.Service, error) {
	panic("implement me")
}

func (m mockRegistry) Watch(option ...registry.WatchOption) (registry.Watcher, error) {
	panic("implement me")
}

func (m mockRegistry) String() string {
	return m.name
}

func TestInit(t *testing.T) {
	assert.NotNil(t, registry.Init(""))
	assert.NotNil(t, registry.Init("///"))
	assert.NotNil(t, registry.Init("/etcd.wpt"))
	assert.NotNil(t, registry.Init("etcdv3://"))
	assert.NotNil(t, registry.Init("etcdv3://123.456"))
	// assert.Nil(t, registry.Init("etcdv3://etcd.wpt:2379"))
}

func TestRegister(t *testing.T) {
	r := &mockRegistry{"1"}
	assert.Nil(t, registry.Register(r.String(), r))
	assert.NotNil(t, registry.Get(r.String()))
	assert.NotNil(t, registry.Register(r.String(), r))
	assert.NotNil(t, registry.Get(r.String()))

	r = &mockRegistry{"3"}
	assert.NotNil(t, registry.Register("", r))
	assert.Nil(t, registry.Get(""))

	r = &mockRegistry{"4"}
	assert.NotNil(t, registry.Register(r.String(), nil))
	assert.Nil(t, registry.Get(r.String()))
}

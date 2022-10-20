package etcdv3

import (
	"fmt"

	"github.com/gowins/dionysus/grpc/registry"
)

var (
	prefix = "/micro-registry"
	Name   = "etcdv3"
)

func init() {
	if err := registry.Register(Name, NewRegistry()); err != nil {
		panic(fmt.Errorf("registry %s error", Name))
	}
}

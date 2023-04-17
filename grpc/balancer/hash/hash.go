package hash

import (
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/metadata"
	"hash/fnv"
	"strings"
	"sync/atomic"
)

const Name = "hash"
const BalancerHashKey = "hash_balancer"

var logger = grpclog.Component("hashBalancer")

// newBuilder creates a new hash balancer builder.
func newBuilder() balancer.Builder {
	return base.NewBalancerBuilder(Name, &hashPickerBuilder{}, base.Config{HealthCheck: true})
}

func init() {
	balancer.Register(newBuilder())
}

type hashPickerBuilder struct{}

func (*hashPickerBuilder) Build(info base.PickerBuildInfo) balancer.Picker {
	logger.Infof("hashPicker: Build called with info: %v", info)
	if len(info.ReadySCs) == 0 {
		return base.NewErrPicker(balancer.ErrNoSubConnAvailable)
	}
	scs := make([]balancer.SubConn, 0, len(info.ReadySCs))
	for sc := range info.ReadySCs {
		scs = append(scs, sc)
	}
	return &hashPicker{
		subConns: scs,
		next:     uint32(1),
	}
}

type hashPicker struct {
	// subConns is the snapshot of the hash balancer when this picker was
	// created. The slice is immutable. Each Get() will do a hash
	// selection from it and return the selected SubConn.
	subConns []balancer.SubConn
	next     uint32
}

func (p *hashPicker) Pick(pickInfo balancer.PickInfo) (balancer.PickResult, error) {
	var sc balancer.SubConn
	subConnsLen := uint32(len(p.subConns))
	md, ok := metadata.FromOutgoingContext(pickInfo.Ctx)
	hashStrings := md.Get(BalancerHashKey)
	if ok && len(hashStrings) != 0 && hashStrings[0] != "" {
		hashString := strings.Join(hashStrings, ",")
		var hash32 = fnv.New32a()
		_, err := hash32.Write([]byte(hashString))
		if err == nil {
			sc = p.subConns[hash32.Sum32()%subConnsLen]
			return balancer.PickResult{SubConn: sc}, nil
		}
	}

	nextIndex := atomic.AddUint32(&p.next, 1)
	sc = p.subConns[nextIndex%subConnsLen]
	return balancer.PickResult{SubConn: sc}, nil
}

package resolver

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/gowins/dionysus/grpc/registry"
	xlog "github.com/gowins/dionysus/log"
	"github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/resolver"
	"google.golang.org/grpc/serviceconfig"
)

type mockWatcher struct{}

func (m *mockWatcher) Next() (*registry.Result, error) {
	return nil, nil
}
func (m *mockWatcher) Stop() {}

func TestBaseResolver(t *testing.T) {
	var br baseResolver
	br.r = &mockWatcher{}
	br.Close()
	br.ResolveNow(resolver.ResolveNowOptions{})
}

type mockedClientConn struct{ state resolver.State }

func (m *mockedClientConn) UpdateState(state resolver.State) error {
	m.state = state
	return nil
}
func (m *mockedClientConn) ReportError(err error)                   {}
func (m *mockedClientConn) NewAddress(addresses []resolver.Address) {}
func (m *mockedClientConn) NewServiceConfig(serviceConfig string)   {}
func (m *mockedClientConn) ParseServiceConfig(serviceConfigJSON string) *serviceconfig.ParseResult {
	return nil
}

func TestReshuffle(t *testing.T) {
	tests := []int{
		subsetSize / 2,
		subsetSize,
		subsetSize * 2,
	}

	for _, test := range tests {
		var addrs []resolver.Address
		for i := 0; i < test; i++ {
			addrs = append(addrs, resolver.Address{Addr: strconv.Itoa(i)})
		}

		assert.Equal(t, len(addrs), len(reshuffle(addrs)))
		assert.NotEqual(t, fmt.Sprintf("%v", addrs), fmt.Sprintf("%v", reshuffle(addrs)))
	}
}

func TestBuildDirectTarget(t *testing.T) {
	target := BuildDirectTarget([]string{"localhost:123", "localhost:456"})
	assert.Equal(t, "direct:///localhost:123,localhost:456", target)
}

func TestBuildDiscovTarget(t *testing.T) {
	target := BuildDiscovTarget([]string{"localhost:123", "localhost:456"}, "foo")
	assert.Equal(t, "discov://localhost:123,localhost:456/foo", target)
}

func TestGetServiceUniqueId(t *testing.T) {
	assert.Equal(t, "serviceName-1027", getServiceUniqueId("serviceName", 1027))
}

//go:norace
func TestSetLog(t *testing.T) {
	convey.Convey("resolver setup log", t, func() {
		xlog.Setup(xlog.SetProjectName("test"))
		SetLog(xlog.GetLogger())
	})
}

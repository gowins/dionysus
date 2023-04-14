package resolver

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/resolver"
)

var subsetSize = 32

func TestDirectBuilder_Build(t *testing.T) {
	tests := []int{
		0,
		1,
		2,
		subsetSize / 2,
		subsetSize,
		subsetSize * 2,
	}

	for _, test := range tests {
		test := test
		t.Run(strconv.Itoa(test), func(t *testing.T) {
			var servers []string
			for i := 0; i < test; i++ {
				servers = append(servers, fmt.Sprintf("localhost:%d", i))
			}

			var b directBuilder
			cc := new(mockedClientConn)
			_, err := b.Build(resolver.Target{URL: url.URL{Scheme: DirectScheme, Path: strings.Join(servers, ",")}}, cc, resolver.BuildOptions{})
			if test == 0 {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}

			m := make(map[resolver.Address]struct{})
			for _, each := range cc.state.Addresses {
				m[each] = struct{}{}
			}

			assert.Equal(t, test*Replica, len(m))
		})
	}
}

func TestBuilder_Scheme(t *testing.T) {
	assert.Equal(t, resolver.Get(DirectScheme).Scheme(), DirectScheme)
	assert.Equal(t, resolver.Get(DiscovScheme).Scheme(), DiscovScheme)
}

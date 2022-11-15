package redis

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNewClient(t *testing.T) {
	Convey("new cleint", t, func() {
		Convey("con not connect", func() {
			client, err := NewClient(&Rdconfig{
				MaxRetries:  0,
				Addr:        "127.0.0.1:56379",
				DB:          0,
				Password:    "",
				PoolSize:    80,
				IdleTimeout: 300,
			})
			So(err, ShouldNotBeNil)
			t.Log(err.Error())
			So(client, ShouldBeNil)
		})
	})
}

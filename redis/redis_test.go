package redis

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNewClient(t *testing.T) {
	Convey("new cleint", t, func() {
		Convey("can connect", func() {
			client, err := NewClient(&Rdconfig{
				MaxRetries:  0,
				Addr:        "r-bp1y6i1udbvmzn8o8upd.redis.rds.aliyuncs.com:6379",
				DB:          0,
				Password:    "kuh6cvk6ZMT*wmp5ufh",
				PoolSize:    80,
				IdleTimeout: 300,
			})
			So(err, ShouldBeNil)
			So(client, ShouldNotBeNil)
		})

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

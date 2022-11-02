package registry

import (
	"crypto/tls"
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
)

func TestOption(t *testing.T) {
	convey.Convey("option", t, func() {
		convey.Convey("options", func() {
			addrOpt := Addrs("127.0.0.1")
			timeoutOpt := Timeout(time.Second)
			secureOpt := Secure(true)
			tlsOpt := TLSConfig(&tls.Config{})
			ttlOpt := TTL(time.Second)
			o := &Options{}
			for _, opt := range []Option{addrOpt, timeoutOpt, secureOpt, tlsOpt, ttlOpt} {
				opt(o)
			}
			convey.So(o.TLSConfig, convey.ShouldNotBeNil)
			convey.So(o.Timeout, convey.ShouldEqual, time.Second)
			convey.So(o.Addrs, convey.ShouldNotBeNil)
			convey.So(o.Secure, convey.ShouldBeTrue)
		})
		convey.Convey("register option", func() {
			registerOpt := RegisterTTL(time.Second)
			o := &RegisterOptions{}
			for _, opt := range []RegisterOption{registerOpt} {
				opt(o)
			}
			convey.So(o.TTL, convey.ShouldEqual, time.Second)
		})
		convey.Convey("watch option", func() {
			watchOpt := WatchService("watch")
			o := &WatchOptions{}
			for _, opt := range []WatchOption{watchOpt} {
				opt(o)
			}
			convey.So(o.Service, convey.ShouldEqual, "watch")
		})
	})
}

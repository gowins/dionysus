package rotatelogs

import (
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
)

// TestNew test new Rotatelog
func TestNew(t *testing.T) {
	convey.Convey("new rotatelogs", t, func() {
		convey.Convey("invalid pattern", func() {
			_, err := New(
				"%",
				WithRotationCount(10),
				WithLinkName(""),
				WithMaxAge(time.Hour),
				WithRotationTime(time.Hour),
			)
			convey.So(err, convey.ShouldNotBeNil)
		})
		convey.Convey("MaxAge and RotationCount", func() {
			_, err := New(
				"/etc/logs",
				WithRotationCount(10),
				WithLinkName(""),
				WithMaxAge(time.Hour),
				WithRotationTime(time.Hour),
			)
			convey.So(err, convey.ShouldNotBeNil)
		})
		convey.Convey("success", func() {
			_, err := New(
				"/etc/log",
				WithRotationCount(0),
				WithLinkName(""),
				WithMaxAge(time.Hour),
				WithRotationTime(time.Hour),
			)
			convey.So(err, convey.ShouldBeNil)
		})
	})
}

// TestGenFileName test genarate file name
func TestGenFileName(t *testing.T) {
	convey.Convey("generate filename", t, func() {
		r, _ := New(
			"/etc/log.%Y%m%d%H%M",
			WithRotationCount(0),
			WithLinkName(""),
			WithMaxAge(time.Hour),
			WithRotationTime(time.Hour),
		)
		r.clock = clockFn(timeFix)
		convey.So(r.genFilename(), convey.ShouldEqual, "/etc/log.202210211600")
		r.clock = clockFn(timeFixUTC)
		convey.So(r.genFilename(), convey.ShouldEqual, "/etc/log.202210210800")
	})
}

// timeFix 2022-10-21 16:26:50
func timeFix() time.Time {
	loc, _ := time.LoadLocation("Asia/Shanghai")
	return time.Unix(1666340810, 0).In(loc)
}

// timeFixUTC 2022-10-21 16:26:50
func timeFixUTC() time.Time {
	return time.Unix(1666340810, 0).In(time.UTC)
}

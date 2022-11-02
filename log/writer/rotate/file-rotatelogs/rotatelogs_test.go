package rotatelogs

import (
	"os"
	"path/filepath"
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
				"/tmp/logs",
				WithRotationCount(10),
				WithLinkName(""),
				WithMaxAge(time.Hour),
				WithRotationTime(time.Hour),
			)
			convey.So(err, convey.ShouldNotBeNil)
		})
		convey.Convey("MaxAge and RotationCount is zero", func() {
			_, err := New(
				tmpDir("logs"),
				WithRotationCount(0),
				WithLinkName(""),
				WithMaxAge(0*time.Hour),
				WithRotationTime(time.Hour),
			)
			convey.So(err, convey.ShouldBeNil)
		})
		convey.Convey("MaxAge and RotationTime negative", func() {
			_, err := New(
				tmpDir("logs"),
				WithRotationCount(10),
				WithLinkName(""),
				WithMaxAge(-1*time.Hour),
				WithRotationTime(-1*time.Hour),
			)
			convey.So(err, convey.ShouldBeNil)
		})
		convey.Convey("success", func() {
			loc, _ := time.LoadLocation("Asia/Shanghai")
			_, err := New(
				tmpDir("log"),
				WithRotationCount(0),
				WithLinkName(""),
				WithMaxAge(time.Hour),
				WithRotationTime(time.Hour),
				WithLocation(loc),
			)
			convey.So(err, convey.ShouldBeNil)
		})
	})
}

// TestWrite test rotatelogs Write
func TestWrite(t *testing.T) {
	convey.Convey("test write", t, func() {
		r, _ := New(
			tmpDir("log.%Y%m%d%H%M"),
			WithLinkName(""),
			WithMaxAge(time.Hour),
			WithRotationTime(time.Hour),
			WithClock(clockFn(timeFix)),
		)
		_, err := r.Write([]byte("ffff"))
		convey.So(err, convey.ShouldBeNil)
		defer os.Remove(r.curFn)
		convey.So(r.curFn, convey.ShouldEqual, tmpDir("log.202210211600"))
		err = r.Close()
		convey.So(err, convey.ShouldBeNil)
	})
}

// TestGetWriter_nolock test getWriter_nolock
func TestGetWriter_nolock(t *testing.T) {
	convey.Convey("test get writer", t, func() {
		convey.Convey("get wirter error", func() {
			r, err := New(
				filepath.Join("vv", tmpDir("log.%Y%m%d%H%M")),
				WithLinkName(""),
				WithMaxAge(time.Hour),
				WithRotationTime(time.Hour),
				WithClock(clockFn(timeFix)),
			)
			convey.So(err, convey.ShouldBeNil)
			_, err = r.getWriter_nolock(false, false)
			convey.So(err, convey.ShouldNotBeNil)
		})
		convey.Convey("get wirter error maxAge", func() {
			r, err := New(
				tmpDir("log.%Y%m%d%H%M"),
				WithLinkName(""),
				WithMaxAge(time.Hour),
				WithRotationTime(time.Hour),
				WithClock(clockFn(timeFix)),
			)
			convey.So(err, convey.ShouldBeNil)
			r.maxAge = -1 * time.Hour
			r.rotationCount = 0
			_, err = r.getWriter_nolock(true, false)
			convey.So(err, convey.ShouldNotBeNil)
			defer os.Remove(r.curFn)
		})
		convey.Convey("rotate", func() {
			r, err := New(
				tmpDir("log.%Y%m%d%H%M"),
				WithLinkName("/tmp/testlog"),
				WithMaxAge(time.Hour),
				WithRotationTime(time.Hour),
				WithClock(clockFn(timeFix)),
			)
			convey.So(err, convey.ShouldBeNil)
			err = r.Rotate()
			convey.So(err, convey.ShouldBeNil)
			os.Remove(r.linkName)
			os.Remove(r.curFn)
		})
	})
}

// TestGenFileName test genarate file name
func TestGenFileName(t *testing.T) {
	convey.Convey("generate filename", t, func() {
		r, _ := New(
			tmpDir("log.%Y%m%d%H%M"),
			WithRotationCount(0),
			WithLinkName(""),
			WithMaxAge(time.Hour),
			WithRotationTime(time.Hour),
			WithClock(clockFn(timeFix)),
		)
		convey.So(r.genFilename(), convey.ShouldEqual, tmpDir("log.202210211600"))
		r.clock = clockFn(timeFixUTC)
		convey.So(r.genFilename(), convey.ShouldEqual, tmpDir("log.202210210800"))
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

func tmpDir(p string) string {
	return filepath.Join(os.TempDir(), p)
}

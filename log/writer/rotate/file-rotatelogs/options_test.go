package rotatelogs

import (
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
)

func TestOptions(t *testing.T) {
	convey.Convey("with options", t, func() {
		convey.So(func() {
			WithClock(clockFn(time.Now))
			WithLocation(time.Local)
			WithLinkName("linkName")
			WithMaxAge(time.Hour)
			WithRotationTime(time.Hour)
			WithRotationCount(100)
		}, convey.ShouldNotPanic)
	})
}

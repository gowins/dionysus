package rotatelogs

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestClock(t *testing.T) {
	convey.Convey("test Clock interface", t, func() {
		convey.So(func() {
			var _ Clock = Local
			var _ Clock = UTC
		}, convey.ShouldNotPanic)
	})
}

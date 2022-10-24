package log

import (
	"runtime/debug"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

type nopWriter struct{}

func (nw *nopWriter) Write(_ []byte) (int, error) {
	return 0, nil
}

// TestSetUp ...
func TestSetUp(t *testing.T) {
	convey.Convey("log set up", t, func() {
		convey.So(func() {
			WithField("stacktrace", string(debug.Stack())).
				Errorf("[error] Panic occurred in start process: %#v", "testing")
		}, convey.ShouldPanic)
		w := &nopWriter{}
		convey.So(func() {
			projectName = "testing"
			Setup(WithWriter(w))
			WithField("stacktrace", string(debug.Stack())).
				Errorf("[error] Panic occurred in start process: %#v", "testing")
		}, convey.ShouldNotPanic)
	})
}

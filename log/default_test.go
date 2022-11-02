package log

import (
	"io"
	"runtime/debug"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

// TestSetUp ...
func TestSetUp(t *testing.T) {
	convey.Convey("log set up", t, func() {
		convey.So(func() {
			WithField("stacktrace", string(debug.Stack())).
				Errorf("[error] Panic occurred in start process: %#v", "testing")
		}, convey.ShouldPanic)
		convey.So(func() {
			projectName = "testing"
			Setup(WithWriter(io.Discard), WithOnFatal(&MockCheckWriteHook{}))
			WithField("stacktrace", string(debug.Stack())).
				Errorf("[error] Panic occurred in start process: %#v", "testing")
		}, convey.ShouldNotPanic)
		projectName = "testing"
		Setup(WithWriter(io.Discard))
		WithField("stacktrace", string(debug.Stack())).
			Errorf("[error] Panic occurred in start process: %#v", "testing")
	})
}

// TestGetLogger ...
func TestGetLogger(t *testing.T) {
	convey.Convey("get logger", t, func() {
		Setup(WithWriter(io.Discard))
		logger := GetLogger()
		logger.Info("something")
	})
}

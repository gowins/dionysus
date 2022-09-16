package shutdown

import (
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestRun(t *testing.T) {
	Convey("Test run panic", t, func() {
		fc := make(chan struct{})
		So(func() { NotifyAfterFinish(fc, func() { panic("want panic") }) }, ShouldNotPanic)
		So(<-fc, ShouldResemble, struct{}{})

		So(func() { NotifyAfterFinish(fc, func() {}) }, ShouldNotPanic)
		So(<-fc, ShouldResemble, struct{}{})
	})
}

func TestShutdown(t *testing.T) {
	var wantCode int
	sysExit = func(code int) {
		wantCode = code
	}

	Convey("Test shutdown success", t, func() {
		fc := make(chan struct{})
		go func() { fc <- struct{}{} }()
		So(func() { WaitingForNotifies(fc, func() {}) }, ShouldNotPanic)
		So(wantCode, ShouldEqual, 0)

		go func() { quit <- os.Signal(nil) }()
		So(func() { WaitingForNotifies(fc, func() {}) }, ShouldNotPanic)
		So(wantCode, ShouldEqual, 0)

		Convey("shutdown panic", func() {
			fc := make(chan struct{})
			go func() { quit <- os.Signal(nil) }()
			So(func() { WaitingForNotifies(fc, func() { panic("want panic") }) }, ShouldNotPanic)
			So(wantCode, ShouldEqual, 3)
			wantCode = 0

		})

	})
}

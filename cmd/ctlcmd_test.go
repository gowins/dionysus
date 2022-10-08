package cmd

import (
	"context"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestRegShutdown(t *testing.T) {
	c := NewCtlCommand()
	Convey("", t, func() {
		So(c.RegRunFunc(func(ctx context.Context) {
			time.Sleep(1000000)
		}),
			ShouldBeNil)
		So(c.RegShutdownFunc(nil), ShouldBeError)
		f := func(ctx context.Context) {}
		So(c.RegShutdownFunc(f), ShouldBeNil)
		So(c.shutdownFunc, ShouldEqual, f)
	})
}

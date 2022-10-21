package cmd

import (
	"context"
	"testing"
	"time"

	"github.com/spf13/pflag"

	. "github.com/smartystreets/goconvey/convey"
)

func TestRegShutdown(t *testing.T) {
	c := NewCtlCommand()
	c.GetShutdownFunc()
	c.GetCmd()
	err := c.RegPreRunFunc("test", func() error {
		return nil
	})
	if err != nil {
		t.Errorf("want error nil get error %v", err)
		return
	}
	err = c.RegPostRunFunc("test2", func() error {
		return nil
	})
	if err != nil {
		t.Errorf("want error nil get error %v", err)
		return
	}
	c.RegFlagSet(&pflag.FlagSet{})
	c.Flags()
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

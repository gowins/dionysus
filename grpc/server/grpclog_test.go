package server

import (
	"io"
	"testing"

	xlog "github.com/gowins/dionysus/log"
	"github.com/smartystreets/goconvey/convey"
)

func TestLoggerWrapper(t *testing.T) {
	convey.Convey("grpclog Loggerv2", t, func() {
		convey.So(func() {
			SetLog(xlog.GetLogger())
		}, convey.ShouldPanic)
		convey.So(func() {
			lw := &loggerWrapper{logger: xlog.GetLogger()}
			lw.Info("logger wrapper")
		}, convey.ShouldPanic)
		convey.So(func() {
			defer func() {
				_ = recover()
			}()
			xlog.Setup(xlog.SetProjectName("test"), xlog.WithWriter(io.Discard), xlog.WithOnFatal(&xlog.MockCheckWriteHook{}))
			lw := &loggerWrapper{logger: xlog.GetLogger()}
			lw.Info("Info")
			lw.Infoln("Infoln")
			lw.Infof("Infof")
			lw.Warning("Warning")
			lw.Warningln("Warningln")
			lw.Warningf("Warningf")
			lw.Error("Error")
			lw.Errorln("Errorln")
			lw.Errorf("Errorf")
			lw.V(1)
			lw.Fatal("Fatal")
			lw.Fatalln("Fatalln")
		}, convey.ShouldNotPanic)
	})
}

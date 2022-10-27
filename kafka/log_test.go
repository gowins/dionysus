package kafka

import (
	"testing"

	logger "github.com/gowins/dionysus/log"
	. "github.com/smartystreets/goconvey/convey"
)

//go:norace
func TestSetLog(t *testing.T) {
	Convey("Before set", t, func() {
		So(kLogger.Logger, ShouldResemble, log)
		So(kErrorLogger.Logger, ShouldResemble, log)
		oldLogger := log

		Convey("set a new logger", func() {
			newLogger, err := logger.New(logger.ZapLogger)
			So(err, ShouldBeNil)
			SetLog(newLogger)

			So(oldLogger, ShouldNotResemble, newLogger)
			So(oldLogger, ShouldNotResemble, log)

			So(kLogger.Logger, ShouldResemble, log)
			So(kErrorLogger.Logger, ShouldResemble, log)
			kLogger.Printf("TestPrint %s", "arg")
			kErrorLogger.Printf("TestPrint %s", "arg")
		})

	})
}

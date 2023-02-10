package kafka

import (
	"testing"

	logger "github.com/gowins/dionysus/log"
	. "github.com/smartystreets/goconvey/convey"
)

//go:norace
func TestSetLog(t *testing.T) {
	Convey("Before set", t, func() {
		So(KLogger.Logger, ShouldResemble, log)
		So(KErrorLogger.Logger, ShouldResemble, log)
		// oldLogger := log

		Convey("set a new logger", func() {
			newLogger, err := logger.New(logger.ZapLogger)
			So(err, ShouldBeNil)
			SetLog(newLogger)

			// So(oldLogger, ShouldNotResemble, newLogger)
			// So(oldLogger, ShouldNotResemble, log)

			So(KLogger.Logger, ShouldResemble, log)
			So(KErrorLogger.Logger, ShouldResemble, log)
			KLogger.Printf("TestPrint %s", "arg")
			KErrorLogger.Printf("TestPrint %s", "arg")
		})

	})
}

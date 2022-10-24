package log

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

// TestNewEncoderConfig test NewEncoderConfig
func TestNewEncoderConfig(t *testing.T) {
	convey.Convey("encoder config", t, func() {
		cfg := NewEncoderConfig()
		convey.So(cfg.MessageKey, convey.ShouldEqual, defaultMessageKey)
		convey.So(cfg.LevelKey, convey.ShouldEqual, defaultLevelKey)
		convey.So(cfg.EncodeLevel, convey.ShouldEqual, LowercaseLevelEncoder)
		convey.So(cfg.TimeKey, convey.ShouldEqual, defaultTimeKey)
		convey.So(cfg.EncodeTime, convey.ShouldEqual, RFC3339MilliTimeEncoder)
		convey.So(cfg.CallerKey, convey.ShouldEqual, defaultCallerKey)
		convey.So(cfg.EncodeCaller, convey.ShouldEqual, FullCallerEncoder)
		convey.So(cfg.StacktraceKey, convey.ShouldEqual, defaultStacktraceKey)
	})
}

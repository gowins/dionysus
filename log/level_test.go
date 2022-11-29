package log

import (
	"io"
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/assert"
)

func TestLevel(t *testing.T) {

	Convey("Converting log levels to string", t, func() {
		So(DebugLevel.String(), ShouldEqual, "debug")
		So(TraceLevel.String(), ShouldEqual, "debug")
		So(InfoLevel.String(), ShouldEqual, "info")
		So(WarnLevel.String(), ShouldEqual, "warn")
		So(ErrorLevel.String(), ShouldEqual, "error")
		So(PanicLevel.String(), ShouldEqual, "panic")
		So(FatalLevel.String(), ShouldEqual, "fatal")
		So(NoticeLevel.String(), ShouldEqual, "notice")
		So(Level(10).String(), ShouldEqual, "Level(10)")
	})

	Convey("The given level is at or above this level", t, func() {
		So(DebugLevel.Enabled(FatalLevel), ShouldBeTrue)
		So(FatalLevel.Enabled(FatalLevel), ShouldBeTrue)
	})

	Convey("The given level is not at or above this level", t, func() {
		So(ErrorLevel.Enabled(InfoLevel), ShouldBeFalse)
		So(PanicLevel.Enabled(WarnLevel), ShouldBeFalse)
	})

	Convey("AllLevels", t, func() {
		So(len(AllLevels()), ShouldEqual, 8)
	})

	EncoderConfig := NewEncoderConfig()
	EncoderConfig.EncodeCaller = ShortCallerEncoder

	opts := []Option{
		WithEncoderCfg(EncoderConfig),
		WithLevelEnabler(DebugLevel),
		AddCaller(),
		AddCallerSkip(1),
		AddStacktrace(ErrorLevel),
		WithOnFatal(&MockCheckWriteHook{}),
	}

	// Debug
	log := captureOutput(func() {
		l, err := New(ZapLogger, append(opts, WithLevelEnabler(DebugLevel), WithWriter(os.Stdout))...)
		if err != nil {
			t.Fatal(err)
		}

		l.Debug("Debug")
		l.Info("Info")
	})
	assert.Contains(t, log, "Debug", "Info")

	log = captureOutput(func() {
		l, err := New(ZapLogger, append(opts, WithLevelEnabler(DebugLevel), WithWriter(os.Stdout))...)
		if err != nil {
			t.Fatal(err)
		}

		l.Debugf("Debug:%d", 123)
		l.Infof("Info:%d", 123)
	})
	assert.Contains(t, log, "Debug:123")
	assert.Contains(t, log, "Info:123")

	// Trace
	log = captureOutput(func() {
		l, err := New(ZapLogger, append(opts, WithLevelEnabler(DebugLevel), WithWriter(os.Stdout))...)
		if err != nil {
			t.Fatal(err)
		}

		l.Debug("Debug")
		l.Trace("Trace")
		l.Info("Info")
		l.Notice("Notice")
	})
	assert.Contains(t, log, "Debug")
	assert.Contains(t, log, "Trace")
	assert.Contains(t, log, "Info")
	assert.Contains(t, log, "Notice")

	log = captureOutput(func() {
		l, err := New(ZapLogger, append(opts, WithLevelEnabler(DebugLevel), WithWriter(os.Stdout))...)
		if err != nil {
			t.Fatal(err)
		}

		l.Debugf("Debug:%d", 123)
		l.Tracef("Trace:%d", 123)
		l.Infof("Info:%d", 123)
	})
	assert.Contains(t, log, "Debug:123")
	assert.Contains(t, log, "Trace:123")
	assert.Contains(t, log, "Info:123")

	// Info
	log = captureOutput(func() {
		l, err := New(ZapLogger, append(opts, WithLevelEnabler(InfoLevel), WithWriter(os.Stdout))...)
		if err != nil {
			t.Fatal(err)
		}

		l.Debug("Debug")
		l.Trace("Trace")
		l.Info("Info")
	})
	assert.NotContains(t, log, "Debug")
	assert.NotContains(t, log, "Trace")
	assert.Contains(t, log, "Info")

	log = captureOutput(func() {
		l, err := New(ZapLogger, append(opts, WithLevelEnabler(InfoLevel), WithWriter(os.Stdout))...)
		if err != nil {
			t.Fatal(err)
		}

		l.Debugf("Debug:%d", 123)
		l.Tracef("Trace:%d", 123)
		l.Infof("Info:%d", 123)
	})
	assert.NotContains(t, log, "Debug:123")
	assert.NotContains(t, log, "Trace:123")
	assert.Contains(t, log, "Info:123")

	// Warn
	log = captureOutput(func() {
		l, err := New(ZapLogger, append(opts, WithLevelEnabler(WarnLevel), WithWriter(os.Stdout))...)
		if err != nil {
			t.Fatal(err)
		}

		l.Debug("Debug")
		l.Info("Info")
		l.Warn("Warn")
	})
	assert.NotContains(t, log, "Debug")
	assert.NotContains(t, log, "Info")
	assert.Contains(t, log, "Warn")

	log = captureOutput(func() {
		l, err := New(ZapLogger, append(opts, WithLevelEnabler(WarnLevel), WithWriter(os.Stdout))...)
		if err != nil {
			t.Fatal(err)
		}

		l.Debugf("Debug:%d", 123)
		l.Infof("Info:%d", 123)
		l.Warnf("Warn:%d", 123)
	})
	assert.NotContains(t, log, "Debug:123")
	assert.NotContains(t, log, "Info:123")
	assert.Contains(t, log, "Warn:123")

	// Error
	log = captureOutput(func() {
		l, err := New(ZapLogger, append(opts, WithLevelEnabler(ErrorLevel), WithWriter(os.Stdout))...)
		if err != nil {
			t.Fatal(err)
		}

		l.Debug("Debug")
		l.Info("Info")
		l.Warn("Warn")
		l.Error("Error")
	})
	assert.NotContains(t, log, "Debug")
	assert.NotContains(t, log, "Info")
	assert.NotContains(t, log, "Warn")
	assert.Contains(t, log, "Error")

	log = captureOutput(func() {
		l, err := New(ZapLogger, append(opts, WithLevelEnabler(ErrorLevel), WithWriter(os.Stdout))...)
		if err != nil {
			t.Fatal(err)
		}

		l.Debugf("Debug:%d", 123)
		l.Infof("Info:%d", 123)
		l.Warnf("Warn:%d", 123)
		l.Errorf("Error:%d", 123)
	})
	assert.NotContains(t, log, "Debug:123")
	assert.NotContains(t, log, "Info:123")
	assert.NotContains(t, log, "Warn:123")
	assert.Contains(t, log, "Error:123")

	// Panic
	log = captureOutput(func() {
		l, err := New(ZapLogger, append(opts, WithLevelEnabler(PanicLevel), WithWriter(os.Stdout))...)
		if err != nil {
			t.Fatal(err)
		}

		l.Debug("Debug")
		l.Info("Info")
		l.Warn("Warn")
		l.Error("Error")
		l.Panic("Panic")
	})
	assert.NotContains(t, log, "Debug")
	assert.NotContains(t, log, "Info")
	assert.NotContains(t, log, "Warn")
	assert.NotContains(t, log, "Error")
	assert.NotContains(t, log, "Panic")

	log = captureOutput(func() {
		l, err := New(ZapLogger, append(opts, WithLevelEnabler(PanicLevel), WithWriter(os.Stdout))...)
		if err != nil {
			t.Fatal(err)
		}

		l.Debugf("Debug:%d", 123)
		l.Infof("Info:%d", 123)
		l.Warnf("Warn:%d", 123)
		l.Errorf("Error:%d", 123)
		l.Panicf("Panic:%d", 123)
	})
	assert.NotContains(t, log, "Debug:123")
	assert.NotContains(t, log, "Info:123")
	assert.NotContains(t, log, "Warn:123")
	assert.NotContains(t, log, "Error:123")
	assert.NotContains(t, log, "Panic:123")

	log = captureOutput(func() {
		l, err := New(ZapLogger, append(opts, WithLevelEnabler(ErrorLevel), WithWriter(os.Stdout))...)
		if err != nil {
			t.Fatal(err)
		}

		l.Debug("Debug")
		l.Info("Info")
		l.Warn("Warn")
		l.Error("Error")
		l.Panic("Panic")
	})
	assert.NotContains(t, log, "Debug")
	assert.NotContains(t, log, "Info")
	assert.NotContains(t, log, "Warn")
	assert.Contains(t, log, "Error")
	assert.Contains(t, log, "Panic")

	log = captureOutput(func() {
		l, err := New(ZapLogger, append(opts, WithLevelEnabler(ErrorLevel), WithWriter(os.Stdout))...)
		if err != nil {
			t.Fatal(err)
		}

		l.Debugf("Debug:%d", 123)
		l.Infof("Info:%d", 123)
		l.Warnf("Warn:%d", 123)
		l.Errorf("Error:%d", 123)
		l.Panicf("Panic:%d", 123)
	})
	assert.NotContains(t, log, "Debug:123")
	assert.NotContains(t, log, "Info:123")
	assert.NotContains(t, log, "Warn:123")
	assert.Contains(t, log, "Error:123")
	assert.Contains(t, log, "Panic:123")

	// Fatal
	log = captureOutput(func() {
		defer func() {
			_ = recover()
		}()
		l, err := New(ZapLogger, append(opts, WithLevelEnabler(FatalLevel), WithWriter(os.Stdout))...)
		if err != nil {
			t.Fatal(err)
		}

		l.Debug("Debug")
		l.Info("Info")
		l.Warn("Warn")
		l.Error("Error")
		l.Panic("Panic")
		l.Fatal("Fatal")
	})
	assert.NotContains(t, log, "Debug")
	assert.NotContains(t, log, "Info")
	assert.NotContains(t, log, "Warn")
	assert.NotContains(t, log, "Error")
	assert.NotContains(t, log, "Panic")
	assert.Contains(t, log, "Fatal")

	log = captureOutput(func() {
		defer func() {
			_ = recover()
		}()
		l, err := New(ZapLogger, append(opts, WithLevelEnabler(FatalLevel), WithWriter(os.Stdout))...)
		if err != nil {
			t.Fatal(err)
		}

		l.Debugf("Debug:%d", 123)
		l.Infof("Info:%d", 123)
		l.Warnf("Warn:%d", 123)
		l.Errorf("Error:%d", 123)
		l.Panicf("Panic:%d", 123)
		l.Noticef("Notice:%d", 123)
		l.Fatalf("Fatal:%d", 123)
	})
	assert.NotContains(t, log, "Debug:123")
	assert.NotContains(t, log, "Info:123")
	assert.NotContains(t, log, "Warn:123")
	assert.NotContains(t, log, "Error:123")
	assert.NotContains(t, log, "Panic:123")
	assert.NotContains(t, log, "Notice:123")
	assert.Contains(t, log, "Fatal:123")

	log = captureOutput(func() {
		defer func() {
			_ = recover()
		}()
		l, err := New(ZapLogger, append(opts, WithLevelEnabler(ErrorLevel), WithWriter(os.Stdout))...)
		if err != nil {
			t.Fatal(err)
		}

		l.Debugf("Debug")
		l.Infof("Info")
		l.Warnf("Warn")
		l.Errorf("Error")
		l.Panicf("Panic")
		l.Fatalf("Fatal")
	})
	assert.NotContains(t, log, "Debug")
	assert.NotContains(t, log, "Info")
	assert.NotContains(t, log, "Warn")
	assert.Contains(t, log, "Error")
	assert.Contains(t, log, "Panic")
	assert.Contains(t, log, "Fatal")

	log = captureOutput(func() {
		defer func() {
			_ = recover()
		}()
		l, err := New(ZapLogger, append(opts, WithLevelEnabler(ErrorLevel), WithWriter(os.Stdout))...)
		if err != nil {
			t.Fatal(err)
		}

		l.Debugf("Debug:%d", 123)
		l.Infof("Info:%d", 123)
		l.Warnf("Warn:%d", 123)
		l.Errorf("Error:%d", 123)
		l.Panicf("Panic:%d", 123)
		l.Fatalf("Fatal:%d", 123)
	})
	assert.NotContains(t, log, "Debug:123")
	assert.NotContains(t, log, "Info:123")
	assert.NotContains(t, log, "Warn:123")
	assert.Contains(t, log, "Error:123")
	assert.Contains(t, log, "Panic:123")
	assert.Contains(t, log, "Fatal:123")

	// ------------------------------------------------------------
	log = captureOutput(func() {
		l, err := New(ZapLogger, append(opts, WithLevelEnabler(DebugLevel), WithWriter(io.Discard), ErrorOutput(os.Stdout))...)
		if err != nil {
			t.Fatal(err)
		}

		l.Info("Info")
	})
	assert.NotContains(t, log, "Info")

	log = captureOutput(func() {
		l, err := New(ZapLogger, append(opts, WithLevelEnabler(DebugLevel), WithWriter(io.Discard), ErrorOutput(os.Stdout))...)
		if err != nil {
			t.Fatal(err)
		}

		l.Infof("Info:%d", 123)
	})
	assert.NotContains(t, log, "Info:123")

	log = captureOutput(func() {
		l, err := New(ZapLogger, append(opts, WithLevelEnabler(InfoLevel), WithWriter(io.Discard), ErrorOutput(os.Stdout))...)
		if err != nil {
			t.Fatal(err)
		}

		l.Debug("Debug")
		l.Info("Info")
	})
	assert.NotContains(t, log, "Debug")
	assert.NotContains(t, log, "Info")

	log = captureOutput(func() {
		l, err := New(ZapLogger, append(opts, WithLevelEnabler(InfoLevel), WithWriter(io.Discard), ErrorOutput(os.Stdout))...)
		if err != nil {
			t.Fatal(err)
		}

		l.Debugf("Debug:%d", 123)
		l.Infof("Info:%d", 123)
	})
	assert.NotContains(t, log, "Debug:123")
	assert.NotContains(t, log, "Info:123")

	log = captureOutput(func() {
		l, err := New(ZapLogger, append(opts, WithLevelEnabler(WarnLevel), WithWriter(io.Discard), ErrorOutput(os.Stdout))...)
		if err != nil {
			t.Fatal(err)
		}

		l.Debug("Debug")
		l.Info("Info")
		l.Warn("Warn")
	})
	assert.NotContains(t, log, "Debug")
	assert.NotContains(t, log, "Info")
	assert.NotContains(t, log, "Warn")

	log = captureOutput(func() {
		l, err := New(ZapLogger, append(opts, WithLevelEnabler(WarnLevel), WithWriter(io.Discard), ErrorOutput(os.Stdout))...)
		if err != nil {
			t.Fatal(err)
		}

		l.Debugf("Debug:%d", 123)
		l.Infof("Info:%d", 123)
		l.Warnf("Warn:%d", 123)
	})
	assert.NotContains(t, log, "Debug:123")
	assert.NotContains(t, log, "Info:123")
	assert.NotContains(t, log, "Warn:123")

	log = captureOutput(func() {
		l, err := New(ZapLogger, append(opts, WithLevelEnabler(ErrorLevel), WithWriter(io.Discard), ErrorOutput(os.Stdout))...)
		if err != nil {
			t.Fatal(err)
		}

		l.Debug("Debug")
		l.Info("Info")
		l.Warn("Warn")
		l.Error("Error")
	})
	assert.NotContains(t, log, "Debug")
	assert.NotContains(t, log, "Info")
	assert.NotContains(t, log, "Warn")
	assert.Contains(t, log, "Error")

	log = captureOutput(func() {
		l, err := New(ZapLogger, append(opts, WithLevelEnabler(ErrorLevel), WithWriter(io.Discard), ErrorOutput(os.Stdout))...)
		if err != nil {
			t.Fatal(err)
		}

		l.Debugf("Debug:%d", 123)
		l.Infof("Info:%d", 123)
		l.Warnf("Warn:%d", 123)
		l.Errorf("Error:%d", 123)
	})
	assert.NotContains(t, log, "Debug:123")
	assert.NotContains(t, log, "Info:123")
	assert.NotContains(t, log, "Warn:123")
	assert.Contains(t, log, "Error:123")

	log = captureOutput(func() {
		l, err := New(ZapLogger, append(opts, WithLevelEnabler(PanicLevel), WithWriter(io.Discard), ErrorOutput(os.Stdout))...)
		if err != nil {
			t.Fatal(err)
		}

		l.Debug("Debug")
		l.Info("Info")
		l.Warn("Warn")
		l.Error("Error")
		l.Panic("Panic")
	})
	assert.NotContains(t, log, "Debug")
	assert.NotContains(t, log, "Info")
	assert.NotContains(t, log, "Warn")
	assert.NotContains(t, log, "Error")
	assert.NotContains(t, log, "Panic")

	log = captureOutput(func() {
		l, err := New(ZapLogger, append(opts, WithLevelEnabler(PanicLevel), WithWriter(io.Discard), ErrorOutput(os.Stdout))...)
		if err != nil {
			t.Fatal(err)
		}

		l.Debugf("Debug:%d", 123)
		l.Infof("Info:%d", 123)
		l.Warnf("Warn:%d", 123)
		l.Errorf("Error:%d", 123)
		l.Panicf("Panic:%d", 123)
	})
	assert.NotContains(t, log, "Debug:123")
	assert.NotContains(t, log, "Info:123")
	assert.NotContains(t, log, "Warn:123")
	assert.NotContains(t, log, "Error:123")
	assert.NotContains(t, log, "Panic:123")

	log = captureOutput(func() {
		l, err := New(ZapLogger, append(opts, WithLevelEnabler(ErrorLevel), WithWriter(io.Discard), ErrorOutput(os.Stdout))...)
		if err != nil {
			t.Fatal(err)
		}

		l.Debug("Debug")
		l.Info("Info")
		l.Warn("Warn")
		l.Error("Error")
		l.Panic("Panic")
	})
	assert.NotContains(t, log, "Debug")
	assert.NotContains(t, log, "Info")
	assert.NotContains(t, log, "Warn")
	assert.Contains(t, log, "Error")
	assert.Contains(t, log, "Panic")

	log = captureOutput(func() {
		l, err := New(ZapLogger, append(opts, WithLevelEnabler(ErrorLevel), WithWriter(io.Discard), ErrorOutput(os.Stdout))...)
		if err != nil {
			t.Fatal(err)
		}

		l.Debugf("Debug:%d", 123)
		l.Infof("Info:%d", 123)
		l.Warnf("Warn:%d", 123)
		l.Errorf("Error:%d", 123)
		l.Panicf("Panic:%d", 123)
	})
	assert.NotContains(t, log, "Debug:123")
	assert.NotContains(t, log, "Info:123")
	assert.NotContains(t, log, "Warn:123")
	assert.Contains(t, log, "Error:123")
	assert.Contains(t, log, "Panic:123")

	log = captureOutput(func() {
		defer func() {
			_ = recover()
		}()
		l, err := New(ZapLogger, append(opts, WithLevelEnabler(FatalLevel), WithWriter(io.Discard), ErrorOutput(os.Stdout))...)
		if err != nil {
			t.Fatal(err)
		}

		l.Debug("Debug")
		l.Info("Info")
		l.Warn("Warn")
		l.Error("Error")
		l.Panic("Panic")
		l.Fatal("Fatal")
	})
	assert.NotContains(t, log, "Debug")
	assert.NotContains(t, log, "Info")
	assert.NotContains(t, log, "Warn")
	assert.NotContains(t, log, "Error")
	assert.NotContains(t, log, "Panic")
	assert.Contains(t, log, "Fatal")

	log = captureOutput(func() {
		defer func() {
			_ = recover()
		}()
		l, err := New(ZapLogger, append(opts, WithLevelEnabler(FatalLevel), WithWriter(io.Discard), ErrorOutput(os.Stdout))...)
		if err != nil {
			t.Fatal(err)
		}

		l.Debugf("Debug:%d", 123)
		l.Infof("Info:%d", 123)
		l.Warnf("Warn:%d", 123)
		l.Errorf("Error:%d", 123)
		l.Panicf("Panic:%d", 123)
		l.Fatalf("Fatal:%d", 123)
	})
	assert.NotContains(t, log, "Debug:123")
	assert.NotContains(t, log, "Info:123")
	assert.NotContains(t, log, "Warn:123")
	assert.NotContains(t, log, "Error:123")
	assert.NotContains(t, log, "Panic:123")
	assert.Contains(t, log, "Fatal:123")

	log = captureOutput(func() {
		defer func() {
			_ = recover()
		}()
		l, err := New(ZapLogger, append(opts, WithLevelEnabler(ErrorLevel), WithWriter(io.Discard), ErrorOutput(os.Stdout))...)
		if err != nil {
			t.Fatal(err)
		}

		l.Debug("Debug")
		l.Info("Info")
		l.Warn("Warn")
		l.Error("Error")
		l.Panic("Panic")
		l.Fatal("Fatal")
	})
	assert.NotContains(t, log, "Debug")
	assert.NotContains(t, log, "Info")
	assert.NotContains(t, log, "Warn")
	assert.Contains(t, log, "Error")
	assert.Contains(t, log, "Panic")
	assert.Contains(t, log, "Fatal")

	log = captureOutput(func() {
		defer func() {
			_ = recover()
		}()
		l, err := New(ZapLogger, append(opts, WithLevelEnabler(ErrorLevel), WithWriter(io.Discard), ErrorOutput(os.Stdout))...)
		if err != nil {
			t.Fatal(err)
		}

		l.Debugf("Debug:%d", 123)
		l.Infof("Info:%d", 123)
		l.Warnf("Warn:%d", 123)
		l.Errorf("Error:%d", 123)
		l.Panicf("Panic:%d", 123)

		l.WithField("xx", "oo").Fatalf("Hello %s", "World")
	})
	assert.NotContains(t, log, "Debug:123")
	assert.NotContains(t, log, "Info:123")
	assert.NotContains(t, log, "Warn:123")
	assert.Contains(t, log, "Error:123")
	assert.Contains(t, log, "Panic:123")
	assert.Contains(t, log, "xx", "oo", "Hello World")
}

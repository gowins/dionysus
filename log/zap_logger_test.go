package log

import (
	"bytes"
	"io"
	"os"
	"sync"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestLogger(t *testing.T) {
	opts := []ZapOption{
		WithZapEncoder(zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())),
		WithZapWriter(os.Stdout),
		WithZapLevelEnabler(DebugLevel),
		WithZapFields([]zapcore.Field{{Key: "wang", Type: zapcore.ReflectType, Interface: "666"}}),
		WithZapOptions([]zap.Option{zap.ErrorOutput(zapcore.AddSync(os.Stderr))}),
		WithZapOnFatal(&MockCheckWriteHook{}),
	}
	zLogger := func(zOpts ...ZapOption) Logger {
		opts := newZapOption()
		for _, opt := range zOpts {
			opt.apply(opts)
		}

		core := zapcore.NewCore(opts.encoder, opts.writer, opts.levelEnabler)

		l := zap.New(core, opts.zOpts...).With(opts.fields...)

		return &zLogger{L: l}
	}(opts...)

	Convey("Logging priority", t, func() {
		defer func() {
			_ = recover()
		}()
		zLogger.Debug("Debug:Hello World")
		zLogger.Debugf("Debugf:%s,Hello World \n", "wong")
		zLogger.Trace("Trace:Hello World")
		zLogger.Tracef("Tracef:%s,Hello World \n", "wong")
		zLogger.Info("Info:Hello World")
		zLogger.Infof("Infof:%s,Hello World \n", "wong")
		zLogger.Warn("Warn:Hello World")
		zLogger.Warnf("Warnf:%s,Hello World \n", "wong")
		zLogger.Error("Error:Hello World")
		zLogger.Errorf("Errorf:%s,Hello World \n", "wong")
		zLogger.Panic("Panic:Hello World")
		zLogger.Panicf("Panicf:%s,Hello World \n", "wong")
		zLogger.Fatal("Fatal:Hello World")
		zLogger.Fatalf("Fatalf:%s,Hello World \n", "wong")
		zLogger.WithField("field1", "xxx").Info("WithField")
		zLogger.WithField("", "").Info("WithField")
		zLogger.WithFields(map[string]interface{}{"field2": "xxx"}).Info("WithFields")
		zLogger.WithFields(nil).Info("WithFields")
	})

	Convey("levelParse", t, func() {
		So(zLogger.SetLogLevel(DebugLevel), ShouldBeNil)
		So(zLogger.SetLogLevel(TraceLevel), ShouldBeNil)
		So(zLogger.SetLogLevel(InfoLevel), ShouldBeNil)
		So(zLogger.SetLogLevel(WarnLevel), ShouldBeNil)
		So(zLogger.SetLogLevel(ErrorLevel), ShouldBeNil)
		So(zLogger.SetLogLevel(PanicLevel), ShouldBeNil)
		So(zLogger.SetLogLevel(FatalLevel), ShouldBeNil)
		So(zLogger.SetLogLevel(10), ShouldNotBeNil)

		So(zLogger.SetLogLevel(DebugLevel), ShouldBeNil)
		zLogger.Debug("debug")
		So(zLogger.SetLogLevel(ErrorLevel), ShouldBeNil)
		zLogger.Debug("no output")
		So(zLogger.SetLogLevel(DebugLevel), ShouldBeNil)
		zLogger.Debug("debug")
	})
}

func captureOutput(f func()) string {
	reader, writer, err := os.Pipe()
	if err != nil {
		panic(err)
	}
	stdout := os.Stdout
	stderr := os.Stderr
	defer func() {
		os.Stdout = stdout
		os.Stderr = stderr
	}()
	os.Stdout = writer
	os.Stderr = writer
	out := make(chan string)
	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func() {
		var buf bytes.Buffer
		wg.Done()
		_, _ = io.Copy(&buf, reader)
		out <- buf.String()
	}()
	wg.Wait()
	f()
	_ = os.Stderr.Sync()
	_ = os.Stdout.Close()
	return <-out
}

package log

import (
	"context"
	"fmt"
	oteltrace "go.opentelemetry.io/otel/trace"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// MockCheckWriteHook should be careful for test, should be recover
type MockCheckWriteHook struct{}

func (mcw *MockCheckWriteHook) OnWrite(e *zapcore.CheckedEntry, _ []zapcore.Field) {
	switch e.Level {
	case zapcore.FatalLevel:
		panic(e.Message)
	default:
	}

}

type zLogger struct {
	L            *zap.Logger
	errL         *zap.Logger
	levelEnabler Level
}

func newZapLogger(lopts ...Option) (Logger, error) {
	opts := Options{
		AddStack:     ErrorLevel,
		LevelEnabler: DebugLevel,
	}

	for _, opt := range lopts {
		opt.apply(&opts)
	}

	zOpts := make([]ZapOption, 0)

	optFunc := []func(opts Options) (ZapOption, error){
		withZapEncoder,
		withZapOptions,
		withZapWriter,
		withZapErrWriter,
		withZapFields,
		withZapLevelEnabler,
	}

	for _, fn := range optFunc {
		opt, err := fn(opts)
		if err != nil {
			return nil, err
		}
		if opt == nil {
			continue
		}
		zOpts = append(zOpts, opt)
	}

	zl := &zLogger{levelEnabler: opts.LevelEnabler}
	return zl.Init(zOpts...)
}

func (zl *zLogger) Init(zOpts ...ZapOption) (Logger, error) {
	opts := newZapOption()
	for _, opt := range zOpts {
		opt.apply(opts)
	}

	core := zapcore.NewCore(opts.encoder, opts.writer, opts.levelEnabler)
	zl.L = zap.New(core, opts.zOpts...).With(opts.fields...)

	if opts.errWriter != nil {
		errCore := zapcore.NewCore(opts.encoder, opts.errWriter, opts.levelEnabler)
		zl.errL = zap.New(errCore, opts.zOpts...).With(opts.fields...)
	}

	return zl, nil
}

func (zl *zLogger) Debug(args ...interface{}) {
	if !zl.levelEnabler.Enabled(DebugLevel) {
		return
	}
	zl.L.Debug(fmt.Sprint(args...))
}

func (zl *zLogger) Debugf(format string, args ...interface{}) {
	if !zl.levelEnabler.Enabled(DebugLevel) {
		return
	}
	zl.L.Debug(fmt.Sprintf(format, args...))
}

func (zl *zLogger) Info(args ...interface{}) {
	if !zl.levelEnabler.Enabled(InfoLevel) {
		return
	}
	zl.L.Info(fmt.Sprint(args...))
}

func (zl *zLogger) Infof(format string, args ...interface{}) {
	if !zl.levelEnabler.Enabled(InfoLevel) {
		return
	}
	zl.L.Info(fmt.Sprintf(format, args...))
}

func (zl *zLogger) Warn(args ...interface{}) {
	if !zl.levelEnabler.Enabled(WarnLevel) {
		return
	}
	zl.L.Warn(fmt.Sprint(args...))
}

func (zl *zLogger) Warnf(format string, args ...interface{}) {
	if !zl.levelEnabler.Enabled(WarnLevel) {
		return
	}
	zl.L.Warn(fmt.Sprintf(format, args...))
}

func (zl *zLogger) Error(args ...interface{}) {
	if !zl.levelEnabler.Enabled(ErrorLevel) {
		return
	}
	if zl.errL != nil {
		zl.errL.Error(fmt.Sprint(args...))
		return
	}
	zl.L.Error(fmt.Sprint(args...))
}

func (zl *zLogger) Errorf(format string, args ...interface{}) {
	if !zl.levelEnabler.Enabled(ErrorLevel) {
		return
	}
	if zl.errL != nil {
		zl.errL.Error(fmt.Sprintf(format, args...))
		return
	}
	zl.L.Error(fmt.Sprintf(format, args...))
}

func (zl *zLogger) Trace(args ...interface{}) {
	if !zl.levelEnabler.Enabled(DebugLevel) {
		return
	}
	zl.L.Debug(fmt.Sprint(args...))
}

func (zl *zLogger) Tracef(format string, args ...interface{}) {
	if !zl.levelEnabler.Enabled(DebugLevel) {
		return
	}
	zl.L.Debug(fmt.Sprintf(format, args...))
}

func (zl *zLogger) Panic(args ...interface{}) {
	if !zl.levelEnabler.Enabled(PanicLevel) {
		return
	}
	if zl.errL != nil {
		zl.errL.Error(fmt.Sprint(args...))
		return
	}
	zl.L.Error(fmt.Sprint(args...))
}

func (zl *zLogger) Panicf(format string, args ...interface{}) {
	if !zl.levelEnabler.Enabled(PanicLevel) {
		return
	}
	if zl.errL != nil {
		zl.errL.Error(fmt.Sprintf(format, args...))
		return
	}
	zl.L.Error(fmt.Sprintf(format, args...))
}

func (zl *zLogger) Fatal(args ...interface{}) {
	if !zl.levelEnabler.Enabled(FatalLevel) {
		return
	}
	if zl.errL != nil {
		zl.errL.Fatal(fmt.Sprint(args...))
		return
	}
	zl.L.Fatal(fmt.Sprint(args...))
}

func (zl *zLogger) Fatalf(format string, args ...interface{}) {
	if !zl.levelEnabler.Enabled(FatalLevel) {
		return
	}
	if zl.errL != nil {
		zl.errL.Fatal(fmt.Sprintf(format, args...))
		return
	}
	zl.L.Fatal(fmt.Sprintf(format, args...))
}

func (zl *zLogger) Notice(args ...interface{}) {
	if !zl.levelEnabler.Enabled(NoticeLevel) {
		return
	}
	zl.L.With(zap.String("alert", "notice")).Info(fmt.Sprint(args...))
}

func (zl *zLogger) Noticef(format string, args ...interface{}) {
	if !zl.levelEnabler.Enabled(NoticeLevel) {
		return
	}
	zl.L.With(zap.String("alert", "notice")).Info(fmt.Sprintf(format, args...))
}

func (zl *zLogger) WithField(key string, value interface{}) Logger {
	if key == "" {
		return zl
	}
	return zl.clone(zapcore.Field{Key: key, Type: zapcore.ReflectType, Interface: value})
}

func (zl *zLogger) WithFields(fields map[string]interface{}) Logger {
	if len(fields) == 0 {
		return zl
	}

	var zfields []zapcore.Field
	for k, v := range fields {
		zfields = append(zfields, zapcore.Field{Key: k, Type: zapcore.ReflectType, Interface: v})
	}

	return zl.clone(zfields...)
}

func (zl *zLogger) WithTraceId(ctx context.Context) Logger {
	traceid := oteltrace.SpanContextFromContext(ctx).TraceID().String()
	if traceid == "" {
		return zl
	}
	return zl.clone(zapcore.Field{Key: "traceId", Type: zapcore.ReflectType, Interface: traceid})
}

func (zl *zLogger) SetLogLevel(lv Level) error {
	if _, err := zapLevelParse(lv); err != nil {
		return err
	}
	zl.levelEnabler = lv
	return nil
}

func (zl *zLogger) LogLevel() Level {
	return zl.levelEnabler
}

func (zl *zLogger) clone(fields ...zapcore.Field) *zLogger {
	if zl.errL != nil {
		return &zLogger{
			L:            zl.L.With(fields...),
			errL:         zl.errL.With(fields...),
			levelEnabler: zl.levelEnabler,
		}
	}

	return &zLogger{L: zl.L.With(fields...), levelEnabler: zl.levelEnabler}
}

func zapLevelParse(lv Level) (zapcore.Level, error) {
	l := zapcore.DebugLevel
	switch lv {
	case DebugLevel:
		l = zapcore.DebugLevel
	case InfoLevel:
		l = zapcore.InfoLevel
	case WarnLevel:
		l = zapcore.WarnLevel
	case ErrorLevel:
		l = zapcore.ErrorLevel
	case PanicLevel:
		l = zapcore.PanicLevel
	case FatalLevel:
		l = zapcore.FatalLevel
	default:
		return l, fmt.Errorf("invalid level: %v ", l)
	}

	return l, nil
}

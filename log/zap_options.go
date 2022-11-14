package log

import (
	"fmt"
	"io"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type zapOptions struct {
	encoder      zapcore.Encoder
	zOpts        []zap.Option
	fields       []zap.Field
	levelEnabler zapcore.LevelEnabler
	writer       zapcore.WriteSyncer
	errWriter    zapcore.WriteSyncer
	onFatal      zapcore.CheckWriteHook
}

func newZapOption() *zapOptions {
	return &zapOptions{
		encoder:      zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		levelEnabler: zapcore.DebugLevel,
		writer:       zapcore.AddSync(os.Stdout),
	}
}

type ZapOption interface {
	apply(*zapOptions)
}

type zapOptionFunc func(*zapOptions)

func (f zapOptionFunc) apply(log *zapOptions) {
	f(log)
}

func WithZapOnFatal(onFatal zapcore.CheckWriteHook) ZapOption {
	return zapOptionFunc(func(zo *zapOptions) {
		if onFatal == nil {
			return
		}
		if zo.zOpts == nil {
			zo.zOpts = []zap.Option{zap.WithFatalHook(onFatal)}
		} else {
			zo.zOpts = append(zo.zOpts, zap.WithFatalHook(onFatal))
		}
	})
}

func WithZapEncoder(encoder zapcore.Encoder) ZapOption {
	return zapOptionFunc(func(opts *zapOptions) {
		opts.encoder = encoder
	})
}

func WithZapOptions(zopts []zap.Option) ZapOption {
	return zapOptionFunc(func(opts *zapOptions) {
		opts.zOpts = zopts
	})
}

func WithZapFields(fields []zap.Field) ZapOption {
	return zapOptionFunc(func(opts *zapOptions) {
		opts.fields = fields
	})
}

func WithZapLevelEnabler(lv Level) ZapOption {
	return zapOptionFunc(func(opts *zapOptions) {
		le, _ := zapLevelParse(lv)
		opts.levelEnabler = le
	})
}

func WithZapWriter(w io.Writer) ZapOption {
	return zapOptionFunc(func(opts *zapOptions) {
		opts.writer = zapcore.AddSync(w)
	})
}

func WithZapErrWriter(w io.Writer) ZapOption {
	return zapOptionFunc(func(opts *zapOptions) {
		opts.errWriter = zapcore.AddSync(w)
	})
}

// -------------------------------------------------------------------------------------------------
func withZapOptions(opts Options) (ZapOption, error) {
	zapOpts := make([]zap.Option, 0)

	if opts.CallerSkip >= 0 {
		zapOpts = append(zapOpts, zap.AddCallerSkip(opts.CallerSkip))
	}

	if lvl, err := zapLevelParse(opts.AddStack); err != nil {
		zapOpts = append(zapOpts, zap.AddStacktrace(zapcore.ErrorLevel))
	} else {
		zapOpts = append(zapOpts, zap.AddStacktrace(lvl))
	}

	if opts.AddCaller {
		zapOpts = append(zapOpts, zap.AddCaller())
	}

	if hook, ok := opts.OnFatal.(zapcore.CheckWriteHook); ok {
		zapOpts = append(zapOpts, zap.WithFatalHook(hook))
	}

	return WithZapOptions(zapOpts), nil
}

func withZapEncoder(opts Options) (ZapOption, error) {

	cfg := zapcore.EncoderConfig{}

	cfg.MessageKey = opts.EncoderCfg.MessageKey
	if cfg.MessageKey == "" {
		cfg.MessageKey = defaultMessageKey
	}
	cfg.LevelKey = opts.EncoderCfg.LevelKey

	switch opts.EncoderCfg.EncodeLevel {
	case LowercaseLevelEncoder:
		cfg.EncodeLevel = zapcore.LowercaseLevelEncoder
	case CapitalLevelEncoder:
		cfg.EncodeLevel = zapcore.CapitalLevelEncoder
	default:
		return nil, fmt.Errorf("invaild EncodeLevel: %v ", opts.EncoderCfg.EncodeLevel)
	}

	cfg.TimeKey = opts.EncoderCfg.TimeKey

	switch opts.EncoderCfg.EncodeTime {
	case RFC3339TimeEncoder:
		cfg.EncodeTime = zapcore.RFC3339TimeEncoder
	case RFC3339MilliTimeEncoder:
		cfg.EncodeTime = rfc3339MilliTimeEncoder
	case RFC3339NanoTimeEncoder:
		cfg.EncodeTime = zapcore.RFC3339NanoTimeEncoder
	default:
		return nil, fmt.Errorf("invaild EncodeTime: %v ", opts.EncoderCfg.EncodeTime)
	}

	cfg.CallerKey = opts.EncoderCfg.CallerKey

	switch opts.EncoderCfg.EncodeCaller {
	case FullCallerEncoder:
		cfg.EncodeCaller = zapcore.FullCallerEncoder
	case ShortCallerEncoder:
		cfg.EncodeCaller = zapcore.ShortCallerEncoder
	default:
		return nil, fmt.Errorf("invaild EncodeCaller: %v ", opts.EncoderCfg.EncodeCaller)
	}

	cfg.StacktraceKey = opts.EncoderCfg.StacktraceKey

	switch opts.Encoder {
	case JSONEncoder:
		return WithZapEncoder(zapcore.NewJSONEncoder(cfg)), nil
	case ConsoleEncoder:
		return WithZapEncoder(zapcore.NewConsoleEncoder(cfg)), nil
	default:
		return nil, fmt.Errorf("invaild Encoder: %v ", opts.EncoderCfg.EncodeTime)
	}
}

func rfc3339MilliTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {

	type appendTimeEncoder interface {
		AppendTimeLayout(time.Time, string)
	}

	if enc, ok := enc.(appendTimeEncoder); ok {
		enc.AppendTimeLayout(t, "2006-01-02T15:04:05.000Z07:00")
		return
	}

	enc.AppendString(t.Format("2006-01-02T15:04:05.000Z07:00"))
}

func withZapWriter(opts Options) (ZapOption, error) {
	if opts.Writer != nil {
		return WithZapWriter(opts.Writer), nil
	}
	return nil, nil
}

func withZapErrWriter(opts Options) (ZapOption, error) {
	if opts.ErrWriter != nil {
		return WithZapErrWriter(opts.ErrWriter), nil
	}
	return nil, nil
}

func withZapFields(opts Options) (ZapOption, error) {
	if opts.Fields != nil {
		var zfields []zapcore.Field
		for k, v := range opts.Fields {
			zfields = append(zfields, zapcore.Field{Key: k, Type: zapcore.ReflectType, Interface: v})
		}
		return WithZapFields(zfields), nil
	}
	return nil, nil
}

func withZapLevelEnabler(opts Options) (ZapOption, error) {
	return WithZapLevelEnabler(opts.LevelEnabler), nil
}

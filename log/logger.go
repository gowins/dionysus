package log

import (
	"fmt"
)

type Logger interface {
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})

	Info(args ...interface{})
	Infof(format string, args ...interface{})

	Warn(args ...interface{})
	Warnf(format string, args ...interface{})

	Error(args ...interface{})
	Errorf(format string, args ...interface{})

	Trace(args ...interface{})
	Tracef(format string, args ...interface{})

	Panic(args ...interface{})
	Panicf(format string, args ...interface{})

	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})

	Notice(args ...interface{})
	Noticef(format string, args ...interface{})

	WithField(key string, value interface{}) Logger
	WithFields(fields map[string]interface{}) Logger

	SetLogLevel(level Level) error

	LogLevel() Level
}

func New(Type LoggerType, opts ...Option) (Logger, error) {
	switch Type {
	case ZapLogger:
		return newZapLogger(opts...)
	default:
		return nil, fmt.Errorf("invaild LoggerType:%v ", Type)
	}
}

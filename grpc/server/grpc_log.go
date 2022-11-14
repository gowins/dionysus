package server

import (
	logger "github.com/gowins/dionysus/log"
	"google.golang.org/grpc/grpclog"
)

var _ grpclog.LoggerV2 = (*loggerWrapper)(nil)

// loggerWrapper wraps logger.Logger into a LoggerV2.
type loggerWrapper struct {
	logger logger.Logger
}

// Info logs to INFO log
func (l *loggerWrapper) Info(args ...interface{}) {
	l.logger.Info(args...)
}

// Infoln logs to INFO log
func (l *loggerWrapper) Infoln(args ...interface{}) {
	l.logger.Info(args...)
}

// Infof logs to INFO log
func (l *loggerWrapper) Infof(format string, args ...interface{}) {
	l.logger.Infof(format, args...)
}

// Warning logs to WARNING log
func (l *loggerWrapper) Warning(args ...interface{}) {
	l.logger.Warn(args...)
}

// Warning logs to WARNING log
func (l *loggerWrapper) Warningln(args ...interface{}) {
	l.logger.Warn(args...)
}

// Warning logs to WARNING log
func (l *loggerWrapper) Warningf(format string, args ...interface{}) {
	l.logger.Warnf(format, args...)
}

// Error logs to ERROR log
func (l *loggerWrapper) Error(args ...interface{}) {
	l.logger.Error(args...)
}

// Errorln Error  logs to ERROR log
func (l *loggerWrapper) Errorln(args ...interface{}) {
	l.logger.Error(args...)
}

// Errorf logs to ERROR log
func (l *loggerWrapper) Errorf(format string, args ...interface{}) {
	l.logger.Errorf(format, args...)
}

// Fatal logs to ERROR log
func (l *loggerWrapper) Fatal(args ...interface{}) {
	l.logger.Fatal(args...)
}

// Fatalln logs to ERROR log
func (l *loggerWrapper) Fatalln(args ...interface{}) {
	l.logger.Fatal(args...)
}

// Fatalf Error logs to ERROR log
func (l *loggerWrapper) Fatalf(format string, args ...interface{}) {
	l.logger.Fatalf(format, args...)
}

// V v returns true for all verbose level.
func (l *loggerWrapper) V(v int) bool {
	return l.logger.LogLevel() == logger.DebugLevel
}

package opentelemetry

import (
	logger "github.com/gowins/dionysus/log"
	"testing"
)

func TestSetLog(t *testing.T) {
	if log.LogLevel() != logger.DebugLevel {
		t.Errorf("want log level debug, got %v", log.LogLevel())
		return
	}
	newLogger, _ := logger.New(logger.ZapLogger, logger.WithLevelEnabler(logger.ErrorLevel))
	SetLog(newLogger)
	if log.LogLevel() != logger.ErrorLevel {
		t.Errorf("want log level error, got %v", log.LogLevel())
		return
	}
}

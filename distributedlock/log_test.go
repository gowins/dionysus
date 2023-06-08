package distributedlock

import (
	logger "github.com/gowins/dionysus/log"
	"testing"
)

func TestSetLog(t *testing.T) {
	l, _ := logger.New(logger.ZapLogger, logger.WithLevelEnabler(logger.FatalLevel))
	SetLog(l)
	if log.LogLevel().String() != logger.FatalLevel.String() {
		t.Errorf("want get level fatal, get %v", log.LogLevel().String())
	}
}

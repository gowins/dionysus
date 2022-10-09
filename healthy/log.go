package healthy

import (
	logger "github.com/gowins/dionysus/log"
)

var (
	defaultLogFields = map[string]interface{}{"pkg": "Health", "type": "internal"}
	log              = initLog()
)

func initLog() logger.Logger {
	l, _ := logger.New(logger.ZapLogger)
	return l.WithFields(defaultLogFields)
}

func SetLog(hyxLog logger.Logger) {
	log = hyxLog.WithFields(defaultLogFields)
	log.Debug("hystrix log is set")
}

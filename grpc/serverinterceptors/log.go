package serverinterceptors

import (
	logger "github.com/gowins/dionysus/log"
)

var (
	defaultLogFields = map[string]interface{}{"pkg": "server-interceptors", "type": "internal"}
	log              = initLog()
)

func initLog() logger.Logger {
	l, _ := logger.New(logger.ZapLogger)
	return l.WithFields(defaultLogFields)
}

func SetLog(log logger.Logger) {
	log = log.WithFields(defaultLogFields)
	log.Debug("grpcServer log is set")
}

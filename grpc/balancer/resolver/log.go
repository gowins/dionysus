package resolver

import logger "github.com/gowins/dionysus/log"

var (
	defaultLogFields = map[string]interface{}{"pkg": "grpcClientResolver", "type": "internal"}
	log              = initLog()
)

func initLog() logger.Logger {
	l, _ := logger.New(logger.ZapLogger)
	return l.WithFields(defaultLogFields)
}

func SetLog(l logger.Logger) {
	log = l.WithFields(defaultLogFields)
	log.Debug("grpcClientResolver log is set")
}

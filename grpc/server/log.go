package server

import (
	logger "github.com/gowins/dionysus/log"
	"google.golang.org/grpc/grpclog"
)

var (
	defaultLogFields = map[string]interface{}{"pkg": "grpcServer", "type": "internal"}
	log              = initLog()
)

func initLog() logger.Logger {
	l, _ := logger.New(logger.ZapLogger)
	return l.WithFields(defaultLogFields)
}

func SetLog(l logger.Logger) {
	log = l.WithFields(defaultLogFields)

	grpclog.SetLoggerV2(&loggerWrapper{logger: log})
	log.Debug("grpcServer log is set")
}

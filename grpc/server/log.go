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
	os := []logger.Option{ // 根据实际需求添加option
		logger.WithLevelEnabler(logger.DebugLevel),
		logger.WithEncoderCfg(logger.NewEncoderConfig()),
		logger.AddCallerSkip(1),
		logger.AddCaller(),
	}
	l, _ := logger.New(logger.ZapLogger, os...)
	return l.WithFields(defaultLogFields)
}

func SetLog(l logger.Logger) {
	log = l.WithFields(defaultLogFields)

	grpclog.SetLoggerV2(&loggerWrapper{logger: log})
	log.Debug("grpcServer log is set")
}

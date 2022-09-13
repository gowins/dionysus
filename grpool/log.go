package grpool

import "go.uber.org/zap"

type Logger interface {
	Debugf(template string, args ...interface{})
	Infof(template string, args ...interface{})
	Warnf(template string, args ...interface{})
	Errorf(template string, args ...interface{})
}

var log Logger = zap.NewExample().Sugar().With("pkg", "grpool", "type", "internal")

func SetLog(grpoolLog Logger) {
	log = grpoolLog
	log.Infof("grpool log is set")
}

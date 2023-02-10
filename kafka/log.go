package kafka

import logger "github.com/gowins/dionysus/log"

var (
	defaultLogFields = map[string]interface{}{"pkg": "kafka", "type": "internal"}
	log              = initLog()
	KLogger          = &debugLogger{log}
	KErrorLogger     = &errorLogger{log}
)

func initLog() logger.Logger {
	l, _ := logger.New(logger.ZapLogger)
	return l.WithFields(defaultLogFields)
}

func SetLog(l logger.Logger) {
	log = l.WithFields(defaultLogFields)
	log.Debug("kafka log is set")
	KLogger = &debugLogger{log}
	KErrorLogger = &errorLogger{log}
}

type debugLogger struct {
	logger.Logger
}

func (l *debugLogger) Printf(f string, args ...interface{}) {
	l.Logger.Debugf("[kafka Debug]"+f, args...)
}

type errorLogger struct {
	logger.Logger
}

func (l *errorLogger) Printf(f string, args ...interface{}) {
	l.Logger.Errorf("[kafka Error]"+f, args...)
}

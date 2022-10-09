package websocket

import (
	logger "github.com/gowins/dionysus/log"
)

var (
	defaultLogFields = map[string]interface{}{"pkg": "websocket", "type": "internal"}
	wLog             = initLog()
)

func initLog() logger.Logger {
	l, _ := logger.New(logger.ZapLogger)
	return l.WithFields(defaultLogFields)
}

type Lwriter struct {
}

func (writer Lwriter) Write(p []byte) (n int, err error) {
	wLog.Errorf(string(p))
	return len(p), nil
}

func SetLog(log logger.Logger) {
	wLog = log.WithFields(defaultLogFields)
	wLog.Debug("websocket log is set")
}

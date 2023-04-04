package ginx

import logger "github.com/gowins/dionysus/log"

var (
	defaultLogFields = map[string]interface{}{"pkg": "gin", "type": "internal"}
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
	log.Debug("gin log is set")
}

type panicWrite struct {
}

func (pw *panicWrite) Write(p []byte) (n int, err error) {
	log.Errorf("dionysus gin panic: %v", string(p))
	return len(p), nil
}

type errorWrite struct {
}

func (ew *errorWrite) Write(p []byte) (n int, err error) {
	log.Errorf("dionysus gin error: %v", string(p))
	return len(p), nil
}

type normalWrite struct {
}

func (nw *normalWrite) Write(p []byte) (n int, err error) {
	log.Infof("dionysus gin: %v", string(p))
	return len(p), nil
}

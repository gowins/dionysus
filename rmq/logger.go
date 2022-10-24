package rmq

import (
	"github.com/apache/rocketmq-client-go/v2/rlog"
	"github.com/gowins/dionysus/log"
)

type rlogger struct {
	dl log.Logger
}

func (r rlogger) Debug(msg string, fields map[string]interface{}) {
	r.dl.WithFields(fields).Debug(msg)
}

func (r rlogger) Info(msg string, fields map[string]interface{}) {
	r.dl.WithFields(fields).Info(msg)
}

func (r rlogger) Warning(msg string, fields map[string]interface{}) {
	r.dl.WithFields(fields).Warn(msg)
}

func (r rlogger) Error(msg string, fields map[string]interface{}) {
	r.dl.WithFields(fields).Error(msg)
}

func (r rlogger) Fatal(msg string, fields map[string]interface{}) {
	r.dl.WithFields(fields).Fatal(msg)
}

func (r rlogger) Level(level string) {
}

func (r rlogger) OutputPath(path string) (err error) {
	return nil
}

func SetupLogger(level log.Level, pname string) {
	dl, _ := log.New(log.ZapLogger, []log.Option{ // 根据实际需求添加option
		log.WithLevelEnabler(level),
		log.WithEncoderCfg(log.NewEncoderConfig()),
		log.AddCallerSkip(3),
		log.AddCaller(),
		log.Fields(map[string]interface{}{
			"app": pname,
		}),
	}...)
	rlog.SetLogLevel("info")
	rlog.SetLogger(&rlogger{
		dl: dl,
	})
}

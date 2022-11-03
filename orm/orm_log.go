package orm

import (
	"os"

	"github.com/gowins/dionysus/log"
	"gorm.io/gorm/logger"
)

// loggerWiter implement gorm logger.Writer
type loggerWiter struct {
	log.Logger
}

// NewWriter return loggerWiter
func NewWriter(opts ...log.Option) (logger.Writer, error) {
	defaultOpts := []log.Option{
		log.WithLevelEnabler(log.InfoLevel),
		log.WithWriter(os.Stdout),
		log.WithEncoderCfg(log.NewEncoderConfig()),
		log.AddCallerSkip(1),
		log.AddCaller(),
	}
	logger, err := log.New(log.ZapLogger, append(defaultOpts, opts...)...)
	if err != nil {
		return nil, err
	}
	return &loggerWiter{Logger: logger}, nil
}

// Printf implement gorm logger.Writer interface Printf method
func (lw *loggerWiter) Printf(msg string, args ...any) {
	lw.Infof(msg, args...)
}

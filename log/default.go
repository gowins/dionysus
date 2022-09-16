package log

var (
	defaultLogger Logger
	projectName   string
	Info          func(args ...interface{})
	Warn          func(args ...interface{})
	Error         func(args ...interface{})
	Fatal         func(args ...interface{})
	Debug         func(args ...interface{})
	Infof         func(format string, args ...interface{})
	Warnf         func(format string, args ...interface{})
	Errorf        func(format string, args ...interface{})
	Fatalf        func(format string, args ...interface{})
	Debugf        func(format string, args ...interface{})

	WithField  func(key string, value interface{}) Logger
	WithFields func(fields map[string]interface{}) Logger
)

func Setup(opts ...Option) {
	os := []Option{ // 根据实际需求添加option
		WithLevelEnabler(DebugLevel),
		WithEncoderCfg(NewEncoderConfig()),
		AddCallerSkip(1),
		AddCaller(),
	}

	//todo: opts隐藏与否
	defaultLogger, _ = New(ZapLogger, append(os, opts...)...)
	if projectName != "" {
		defaultLogger = defaultLogger.WithField("app", projectName)
	}
	Debug = defaultLogger.Debug
	Debugf = defaultLogger.Debugf
	Info = defaultLogger.Info
	Infof = defaultLogger.Infof
	Warn = defaultLogger.Warn
	Warnf = defaultLogger.Warnf
	Error = defaultLogger.Error
	Errorf = defaultLogger.Errorf
	Fatal = defaultLogger.Fatal
	Fatalf = defaultLogger.Fatalf
	WithField = defaultLogger.WithField
	WithFields = defaultLogger.WithFields

}

func SetProjectName(pname string) Option {
	return Fields(map[string]interface{}{
		"app": pname,
	})
}

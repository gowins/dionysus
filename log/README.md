# 日志组件说明
日志组件分为 Logger 跟 Writer 两块，对第三方库进行二次封装提供统一接口。

## logger包说明
- logger包是一组实现Logger接口的组件。

```go
type Logger interface {
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Warn(args ...interface{})
	Warnf(format string, args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Trace(args ...interface{})
	Tracef(format string, args ...interface{})
	Panic(args ...interface{})
	Panicf(format string, args ...interface{})
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})

	WithField(key string, value interface{}) Logger
	WithFields(fields map[string]interface{}) Logger

	SetLogLevel(level Level) error
}
```

## writer包说明
- writer包是一组实现io.Writer接口的组件。

```go
type Writer interface {
	Write(p []byte) (n int, err error)
}
```

## Option
```go
    // log options    

	// EncoderConfig 说明
	EncoderConfig := EncoderConfig{}
	EncoderConfig.MessageKey = "msg" // MessageKey

	EncoderConfig.LevelKey = "level"                  // LevelKey
	EncoderConfig.EncodeLevel = LowercaseLevelEncoder // 具体请看 LevelEncoder

	EncoderConfig.TimeKey = "@timestamp"               // TimeKey
	EncoderConfig.EncodeTime = RFC3339MilliTimeEncoder // 具体请看 TimeEncoder，默认RFC3339TimeEncoder

	// caller
	// note: CallerKey不为空 且 设置AddCaller 才会输出
	EncoderConfig.CallerKey = "caller"             // CallerKey
	EncoderConfig.EncodeCaller = FullCallerEncoder // CallerEncoder，默认FullCallerEncoder
	AddCaller()                                    // 显示日志调用者的文件名跟行号，默认不打开
	AddCallerSkip(1)                               // 跳过多少级caller 

	// Stacktrace
	// note: StacktraceKey不为空 且 达到Stacktrace设置的级别 才会输出
	EncoderConfig.StacktraceKey = "detail" // StacktraceKey ，配合 AddStacktrace 使用
	AddStacktrace(ErrorLevel)              // 默认ErrorLevel

	WithEncoderCfg(EncoderConfig) // 也可以使用默认配置 WithEncoderCfg(NewEncoderConfig())
	WithEncoder(JSONEncoder)      // 具体请看 Encoder，默认JSONEncoder

	WithLevelEnabler(DebugLevel)                 // 可选，设置日志输出级别，默认DebugLevel
	WithWriter(os.Stdout)                        // 可选，设置日志的wirter
	Fields(map[string]interface{}{"wpt": "yes"}) // 可选，增加字段到日志输出
	ErrorOutput(os.Stderr)                       // 可选，ErrorLevel及之后的日志级别将往这里输出
```

## Example
```go
// Logger
opts := []Option{ // 根据实际需求添加option
    WithLevelEnabler(DebugLevel),
    WithEncoderCfg(NewEncoderConfig()), 
}

l, err := New(ZapLogger, opts...)
l.Info("Hello World!")

// Writer
wopts := []rotate.Option{ // 根据实际需求添加option
    rotate.WithLogDir("/data/log"),
    rotate.WithLogSubDir("info"),
}
w := NewWriter(wopts...)
w.Write([]byte("Hello World!"))
```

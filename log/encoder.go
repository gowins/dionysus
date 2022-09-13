package log

// LowercaseLevelEncoder serializes a Level to a lowercase string. For example,
// InfoLevel is serialized to "info".
// CapitalLevelEncoder serializes a Level to an all-caps string. For example,
// InfoLevel is serialized to "INFO".
type LevelEncoder int

const (
	LowercaseLevelEncoder LevelEncoder = iota
	CapitalLevelEncoder
)

// FullCallerEncoder serializes a caller in /full/path/to/package/file:line
// format.
// ShortCallerEncoder serializes a caller in package/file:line format, trimming
// all but the final directory from the full path.
type CallerEncoder int

const (
	FullCallerEncoder CallerEncoder = iota
	ShortCallerEncoder
)

// A TimeEncoder serializes a time.Time to a primitive type.
type TimeEncoder int

const (
	RFC3339TimeEncoder TimeEncoder = iota
	RFC3339MilliTimeEncoder
	RFC3339NanoTimeEncoder
)

type Encoder int

const (
	JSONEncoder Encoder = iota
	ConsoleEncoder
)

type EncoderConfig struct {
	MessageKey    string
	LevelKey      string
	EncodeLevel   LevelEncoder
	TimeKey       string
	EncodeTime    TimeEncoder
	CallerKey     string
	EncodeCaller  CallerEncoder
	StacktraceKey string
}

const (
	defaultMessageKey    = "msg"
	defaultLevelKey      = "level"
	defaultTimeKey       = "@timestamp"
	defaultCallerKey     = "caller"
	defaultStacktraceKey = "detail"
)

func NewEncoderConfig() EncoderConfig {
	return EncoderConfig{
		MessageKey:    defaultMessageKey,
		LevelKey:      defaultLevelKey,
		EncodeLevel:   LowercaseLevelEncoder,
		TimeKey:       defaultTimeKey,
		EncodeTime:    RFC3339MilliTimeEncoder,
		CallerKey:     defaultCallerKey,
		EncodeCaller:  FullCallerEncoder,
		StacktraceKey: defaultStacktraceKey,
	}
}

package kafka

import (
	"time"

	"github.com/segmentio/kafka-go"
)

const defaultCommitInterval = time.Millisecond * 100

// default settings
var (
	defaultWriterConfig = kafka.WriterConfig{
		Async:        false,
		BatchSize:    100,
		BatchBytes:   2 << 19, // 1MB
		BatchTimeout: time.Millisecond * 100,
		Logger:       nil,
		ErrorLogger:  nil,
	}

	defaultReaderConfig = kafka.ReaderConfig{
		MinBytes:    2 << 9,  // 1K
		MaxBytes:    2 << 23, // 16M
		StartOffset: kafka.LastOffset,
		Logger:      nil,
		ErrorLogger: nil,
	}
)

type (
	ReaderOption func(*kafka.ReaderConfig)
	WriterOption func(*kafka.WriterConfig)
)

func WriterWithAsync() WriterOption {
	return func(config *kafka.WriterConfig) {
		config.Async = true
	}
}

func ReaderWithAsync() ReaderOption {
	return func(config *kafka.ReaderConfig) {
		config.CommitInterval = defaultCommitInterval
	}
}

func ReaderWithDialer(dialer *kafka.Dialer) ReaderOption {
	return func(config *kafka.ReaderConfig) {
		config.Dialer = dialer
	}
}

func ReaderWithOffset(offset int64) ReaderOption {
	return func(config *kafka.ReaderConfig) {
		config.StartOffset = offset
	}
}

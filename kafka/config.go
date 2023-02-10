package kafka

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/plain"
)

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
		StartOffset: kafka.FirstOffset,
		Logger:      nil,
		ErrorLogger: nil,
	}
)

type (
	ReaderOption func(*kafka.ReaderConfig)
	WriterOption func(*kafka.WriterConfig)
	DialerOption func(*kafka.Dialer)
)

func WriterWithAsync() WriterOption {
	return func(config *kafka.WriterConfig) {
		config.Async = true
	}
}

func WriterWithDialer(dialer *kafka.Dialer) WriterOption {
	return func(config *kafka.WriterConfig) {
		config.Dialer = dialer
	}
}

func WriterWithBatchSize(batchSize int) WriterOption {
	return func(config *kafka.WriterConfig) {
		config.BatchSize = batchSize
	}
}

func WriterWithBatchBytes(batchBytes int) WriterOption {
	return func(config *kafka.WriterConfig) {
		config.BatchBytes = batchBytes
	}
}

func WriterWithBatchTimeout(batchTimeout time.Duration) WriterOption {
	return func(config *kafka.WriterConfig) {
		config.BatchTimeout = batchTimeout
	}
}

func WriterWithLogger(debugLogger, errorLogger kafka.Logger) WriterOption {
	return func(config *kafka.WriterConfig) {
		if debugLogger != nil {
			config.Logger = debugLogger
		}
		if errorLogger != nil {
			config.ErrorLogger = errorLogger
		}
	}
}

// ReaderWithCommitInterval be careful commitInterval may cause msg repeated consume when server error exit.
func ReaderWithCommitInterval(commitInterval time.Duration) ReaderOption {
	return func(config *kafka.ReaderConfig) {
		config.CommitInterval = commitInterval
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

func ReaderWithMinBytes(minBytes int) ReaderOption {
	return func(config *kafka.ReaderConfig) {
		config.MinBytes = minBytes
	}
}

func ReaderWithMaxBytes(maxBytes int) ReaderOption {
	return func(config *kafka.ReaderConfig) {
		config.MaxBytes = maxBytes
	}
}

func DialerWithCaCertTls(file string) DialerOption {
	return func(dialer *kafka.Dialer) {
		certBytes, err := ioutil.ReadFile(file)
		if err != nil {
			log.Panicf("kafka client read cert file failed %v", err)
		}
		clientCertPool := x509.NewCertPool()
		ok := clientCertPool.AppendCertsFromPEM(certBytes)
		if !ok {
			log.Panicf("kafka client failed to parse root certificate")
		}
		dialer.TLS = &tls.Config{RootCAs: clientCertPool, InsecureSkipVerify: true} // nolint:gosec
	}
}

func DialerWithUserAndPasswd(user, passwd string) DialerOption {
	return func(dialer *kafka.Dialer) {
		mechanism := plain.Mechanism{
			Username: user,
			Password: passwd,
		}
		dialer.SASLMechanism = mechanism
	}
}

func GetKafkaDialer(options ...DialerOption) *kafka.Dialer {
	kafkaDialer := &kafka.Dialer{
		Timeout:   10 * time.Second,
		DualStack: true,
	}
	for _, opt := range options {
		opt(kafkaDialer)
	}
	return kafkaDialer
}

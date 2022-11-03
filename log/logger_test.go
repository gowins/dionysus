package log

import (
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewLogger(t *testing.T) {
	l, _ := New(ZapLogger)
	l.Info("INFO")

	EncoderConfig := NewEncoderConfig()
	EncoderConfig.EncodeCaller = ShortCallerEncoder

	opts := []Option{
		WithEncoderCfg(EncoderConfig),
		WithLevelEnabler(DebugLevel),
		AddCaller(),
		AddCallerSkip(1),
		AddStacktrace(ErrorLevel),
	}

	l, err := New(ZapLogger, opts...)
	if err != nil {
		t.Fatal(err)
	}

	l.Info("Hello World!")
	l.Error("Hello World!")

	_, err = New(7777)
	assert.EqualError(t, err, "invaild LoggerType:7777 ")
}

func TestErrWriter(t *testing.T) {

	EncoderConfig := NewEncoderConfig()
	EncoderConfig.EncodeCaller = ShortCallerEncoder

	opts := []Option{
		WithEncoderCfg(EncoderConfig),
		WithLevelEnabler(DebugLevel),
		AddCaller(),
		AddCallerSkip(1),
		AddStacktrace(ErrorLevel),
		WithOnFatal(&MockCheckWriteHook{}),
	}

	log := captureOutput(func() {
		l, err := New(ZapLogger, append(opts, WithWriter(os.Stdout), ErrorOutput(io.Discard))...)
		if err != nil {
			t.Fatal(err)
		}

		l.Info("Info")
		l.Error("Error")
	})
	assert.Contains(t, log, "Info")
	assert.NotContains(t, log, "Error")

	log = captureOutput(func() {
		l, err := New(ZapLogger, append(opts, ErrorOutput(os.Stdout), WithWriter(io.Discard))...)
		if err != nil {
			t.Fatal(err)
		}

		l.Info("Info")
		l.Error("Error")
	})
	assert.NotContains(t, log, "Info")
	assert.Contains(t, log, "Error")
}

func TestFields(t *testing.T) {

	EncoderConfig := NewEncoderConfig()
	EncoderConfig.EncodeCaller = ShortCallerEncoder

	opts := []Option{
		WithEncoderCfg(EncoderConfig),
		WithLevelEnabler(DebugLevel),
		AddCaller(),
		AddCallerSkip(1),
		AddStacktrace(ErrorLevel),
		WithOnFatal(&MockCheckWriteHook{}),
	}

	log := captureOutput(func() {

		l, err := New(ZapLogger, append(opts, WithWriter(os.Stdout))...)
		if err != nil {
			t.Fatal(err)
		}

		l.WithField("hh", "xx").Info("Info")
	})
	assert.Contains(t, log, "Info", "hh", "xx")

	log = captureOutput(func() {

		opts = append(
			opts,
			WithWriter(os.Stdout),
			Fields(map[string]interface{}{"oo": "qq"}),
		)

		l, err := New(ZapLogger, opts...)
		if err != nil {
			t.Fatal(err)
		}

		l.WithFields(map[string]interface{}{"hh": "xx"}).Info("Info")
	})
	assert.Contains(t, log, "Info", "hh", "xx", "oo", "qq")
}

func TestConsoleEncoder(t *testing.T) {

	EncoderConfig := EncoderConfig{}
	EncoderConfig.MessageKey = "msg"

	opts := []Option{
		WithLevelEnabler(DebugLevel),
		WithEncoderCfg(EncoderConfig),
		WithEncoder(ConsoleEncoder),
	}

	l, err := New(ZapLogger, opts...)
	if err != nil {
		t.Fatal(err)
	}

	l.Info("tracing")
}

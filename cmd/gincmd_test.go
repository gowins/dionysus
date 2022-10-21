package cmd

import (
	"github.com/spf13/pflag"
	"testing"
	"time"
)

func TestNewGinCommand(t *testing.T) {
	gcmd := NewGinCommand()
	gcmd.RegFlagSet(&pflag.FlagSet{})
	gcmd.Flags()
	err := gcmd.RegPostRunFunc("test", func() error {
		return nil
	})
	if err != nil {
		t.Errorf("want error nil get error %v", err)
		return
	}
	err = gcmd.RegPreRunFunc("test1", func() error {
		return nil
	})
	if err != nil {
		t.Errorf("want error nil get error %v", err)
		return
	}
	gcmd.GetCmd()
	gcmd.GetShutdownFunc()
	go func() {
		gcmd.startServer()
	}()
	time.Sleep(time.Second * 5)
	gcmd.stopServer()
}

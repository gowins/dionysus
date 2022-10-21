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
	gcmd.RegPostRunFunc("test", func() error {
		return nil
	})
	gcmd.RegPreRunFunc("test1", func() error {
		return nil
	})
	gcmd.GetCmd()
	gcmd.GetShutdownFunc()
	go func() {
		gcmd.startServer()
	}()
	time.Sleep(time.Second * 5)
	gcmd.stopServer()
}

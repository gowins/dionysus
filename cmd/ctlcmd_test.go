package cmd

import (
	"fmt"
	"github.com/spf13/pflag"
	"testing"
)

func TestNewCtlCommand(t *testing.T) {
	ctl := NewCtlCommand()
	err := ctl.RegShutdownFunc(nil)
	if err == nil {
		t.Errorf("want error not nil")
		return
	}
	err = ctl.RegRunFunc(nil)
	if err == nil {
		t.Errorf("want error not nil")
		return
	}
	_ = ctl.RegShutdownFunc(func() {
		fmt.Printf("this is shutdown\n")
	})

	_ = ctl.RegRunFunc(func() error {
		fmt.Printf("this is run func\n")
		return nil
	})
	ctl.RegFlagSet(&pflag.FlagSet{})
	ctl.Flags()
	if ctl.GetCmd().Use != CtlUse {
		t.Errorf("want get cmd use %v get %v", CtlUse, ctl.GetCmd().Use)
		return
	}
	if ctl.GetShutdownFunc() == nil {
		t.Errorf("want GetShutdownFunc not get nil")
		return
	}
	err = ctl.cmd.RunE(nil, nil)
	if err != nil {
		t.Errorf("want error nil get err %v", err)
	}
}

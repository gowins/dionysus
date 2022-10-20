//go:build !test

package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"time"

	"github.com/gowins/dionysus"
)

type testCmd struct {
	cmd      *cobra.Command
	stopChan chan struct{}
}

func (tc *testCmd) GetCmd() *cobra.Command {
	return tc.cmd
}

func (tc *testCmd) GetShutdownFunc() func() {
	return func() {
		fmt.Printf("this is testCmd shutdown func\n")
		tc.stopChan <- struct{}{}
	}
}

func (tc *testCmd) RegFlagSet(set *pflag.FlagSet) {
}

func (tc *testCmd) Flags() *pflag.FlagSet {
	return nil
}

func main() {
	tc := &testCmd{
		cmd:      &cobra.Command{Use: "testCmd", Short: "just for test"},
		stopChan: make(chan struct{}),
	}
	tc.cmd.RunE = func(cmd *cobra.Command, args []string) error {
		timer1 := time.NewTicker(5 * time.Second)
		for {
			select {
			case <-timer1.C:
				fmt.Printf("this is RunE %v\n", time.Now().String())
			case <-tc.stopChan:
				fmt.Printf("this is stopChan %v\n", time.Now().String())
				return nil
			}
		}
	}
	d := dionysus.NewDio()
	_ = d.PreRunRegWithPriority("userPre2", 2, func() error {
		fmt.Printf("this is userPre2\n")
		return nil
	})
	_ = d.PreRunRegWithPriority("userPre1", 1, func() error {
		fmt.Printf("this is userPre1\n")
		return nil
	})
	_ = d.PostRunRegWithPriority("userPost2", 2, func() error {
		fmt.Printf("this is userPost2\n")
		return nil
	})
	_ = d.PostRunRegWithPriority("userPost1", 1, func() error {
		fmt.Printf("this is userPost1\n")
		return nil
	})
	_ = d.PreRunStepsAppend("userPreA1", func() error {
		fmt.Printf("this is userPreA1\n")
		return nil
	})
	_ = d.PreRunStepsAppend("userPreA2", func() error {
		fmt.Printf("this is userPreA2\n")
		return nil
	})
	_ = d.PostRunStepsAppend("userPostA1", func() error {
		fmt.Printf("this is userPostA1\n")
		return nil
	})
	_ = d.PostRunStepsAppend("userPostA2", func() error {
		fmt.Printf("this is userPostA2\n")
		return nil
	})
	if err := d.DioStart("testcmd", tc); err != nil {
		fmt.Printf("DioStart err is %v\n", err)
	}
}

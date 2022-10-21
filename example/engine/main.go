//go:build !test

package main

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

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
	//register cmd
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

	//dio init
	d := dionysus.NewDio()
	//todo  对齐sysPriority 10个定义、后续append
	_ = d.PreRunRegWithPriority("userPre1", 102, func() error {
		fmt.Printf("this is userPre1\n")
		return nil
	})
	_ = d.PreRunRegWithPriority("userPre2", 102, func() error {
		fmt.Printf("this is userPre2\n")
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

	if err := d.DioStart("testcmd", tc); err != nil {
		fmt.Printf("DioStart err is %v\n", err)
	}
}

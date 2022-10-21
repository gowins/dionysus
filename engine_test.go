package dionysus

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"testing"
	"time"
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

func TestDio_DioStart(t *testing.T) {
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
	d := NewDio()
	d.cmd.SetArgs([]string{"testCmd"})
	gotPreStrings := []string{}
	wantPreStrings := []string{"this is userPre1\n", "this is userPre2\n", "this is userPreA1\n", "this is userPreA2\n"}
	gotPostStrings := []string{}
	wantpostStrings := []string{"this is userPost1\n", "this is userPost2\n", "this is userPostA1\n", "this is userPostA2\n"}
	_ = d.PreRunRegWithPriority("userPre2", 102, func() error {
		fmt.Printf("this is userPre2\n")
		gotPreStrings = append(gotPreStrings, "this is userPre2\n")
		return nil
	})
	_ = d.PreRunRegWithPriority("userPre1", 101, func() error {
		fmt.Printf("this is userPre1\n")
		gotPreStrings = append(gotPreStrings, "this is userPre1\n")
		return nil
	})
	_ = d.PostRunRegWithPriority("userPost2", 102, func() error {
		fmt.Printf("this is userPost2\n")
		gotPostStrings = append(gotPostStrings, "this is userPost2\n")
		return nil
	})
	_ = d.PostRunRegWithPriority("userPost1", 101, func() error {
		fmt.Printf("this is userPost1\n")
		gotPostStrings = append(gotPostStrings, "this is userPost1\n")
		return nil
	})
	_ = d.PreRunStepsAppend("userPreA1", func() error {
		gotPreStrings = append(gotPreStrings, "this is userPreA1\n")
		fmt.Printf("this is userPreA1\n")
		return nil
	})
	_ = d.PreRunStepsAppend("userPreA2", func() error {
		gotPreStrings = append(gotPreStrings, "this is userPreA2\n")
		fmt.Printf("this is userPreA2\n")
		return nil
	})
	_ = d.PostRunStepsAppend("userPostA1", func() error {
		fmt.Printf("this is userPostA1\n")
		gotPostStrings = append(gotPostStrings, "this is userPostA1\n")
		return nil
	})
	_ = d.PostRunStepsAppend("userPostA2", func() error {
		fmt.Printf("this is userPostA2\n")
		gotPostStrings = append(gotPostStrings, "this is userPostA2\n")
		return nil
	})
	go func() {
		if err := d.DioStart("testcmd", tc); err != nil {
			fmt.Printf("DioStart err is %v\n", err)
		}
	}()
	time.Sleep(10 * time.Second)
	fn := tc.GetShutdownFunc()
	fn()
	time.Sleep(5 * time.Second)
	for index, str := range wantpostStrings {
		if str != gotPostStrings[index] {
			t.Errorf("get %v, want %v\n", gotPostStrings[index], str)
			return
		}
	}
	for index, str := range wantPreStrings {
		if str != gotPreStrings[index] {
			t.Errorf("get %v, want %v\n", gotPreStrings[index], str)
			return
		}
	}
}

func TestDio_DioStartRun(t *testing.T) {
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
	d := NewDio()
	d.cmd.SetArgs([]string{"testCmd"})
	go func() {
		if err := d.DioStart("testcmd", tc); err != nil {
			fmt.Printf("DioStart err is %v\n", err)
		}
	}()
	time.Sleep(10 * time.Second)
	fn := tc.GetShutdownFunc()
	fn()
}

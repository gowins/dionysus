package dionysus

import (
	"fmt"
	"io"
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/gowins/dionysus/cmd"
	"github.com/gowins/dionysus/log"
	"github.com/gowins/dionysus/step"
	"github.com/smartystreets/goconvey/convey"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
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

//go:norace
func TestDioStartCmd(t *testing.T) {
	convey.Convey("commander", t, func() {
		tc := &testCmd{
			// cmd:
			cmd:      nil,
			stopChan: make(chan struct{}),
		}
		dio := NewDio()
		err := dio.DioStart("test_nil_cmd", tc)
		if err == nil {
			t.Errorf("expetced error is not nil, got %v", err)
		}
		convey.So(err, convey.ShouldNotBeNil)
	})
}

func TestWrapCobraRun(t *testing.T) {
	log.Setup(log.SetProjectName("testing"), log.WithWriter(io.Discard))
	convey.Convey("test wrap cobra run", t, func() {
		convey.Convey("wrap panic", func() {
			w1 := wrapCobrCmdRun(func(_ *cobra.Command, _ []string) error {
				s := make([]int, 1)
				fmt.Println(s[1])
				return nil
			}, func() {})
			cmd := &cobra.Command{Use: "testCobraRun", Short: "wrapp cobra run"}
			ars := make([]string, 0)
			err := w1(cmd, ars)
			convey.So(err, convey.ShouldBeNil)
		})
		convey.Convey("wrap return error", func() {
			w1 := wrapCobrCmdRun(func(_ *cobra.Command, _ []string) error {
				return fmt.Errorf("just test wrap return error")
			}, func() {})
			cmd := &cobra.Command{Use: "testCobraRun", Short: "wrapp cobra run"}
			ars := make([]string, 0)
			err := w1(cmd, ars)
			convey.So(err, convey.ShouldBeNil)
		})
	})
}

// TestWaitForNotify test
func TestWaitingForNotifies(t *testing.T) {
	log.Setup(log.SetProjectName("testing"), log.WithWriter(io.Discard))
	convey.Convey("waiting for notifies", t, func() {
		convey.Convey("quit SIGQUIT signal", func() {
			convey.So(func() {
				finishChan := make(chan struct{})
				quit := make(chan os.Signal, 1)
				quit <- syscall.SIGQUIT
				waitingForNotifies(quit, finishChan, func() {})
			}, convey.ShouldNotPanic)
		})
		convey.Convey("quit SIGQUIT signal nil shout down", func() {
			convey.So(func() {
				finishChan := make(chan struct{})
				quit := make(chan os.Signal, 1)
				quit <- syscall.SIGQUIT
				waitingForNotifies(quit, finishChan, nil)
			}, convey.ShouldNotPanic)
		})
	})
}

//go:norace
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
	wantPreStrings := []string{"this is PreRun1", "this is PreRun2", "this is PreRun3", "this is PreRun4", "this is PreRun5",
		"this is PreRun6", "this is PreRun7", "this is PreRun8", "this is PreRun9", "this is PreRun10", "this is userPreA1\n", "this is userPreA2\n"}
	gotPostStrings := []string{}
	wantpostStrings := []string{"this is PostRun1", "this is PostRun2", "this is PostRun3", "this is PostRun4", "this is PostRun5",
		"this is PostRun6", "this is PostRun7", "this is PostRun8", "this is PostRun9", "this is PostRun10",
		"this is userPostA1\n", "this is userPostA2\n"}
	d.RegUserFirstPostRunStep(step.InstanceStep{StepName: "PostRun1", Func: func() error {
		fmt.Printf("this is PostRun1\n")
		gotPostStrings = append(gotPostStrings, "this is PostRun1")
		return nil
	}})
	d.RegUserSecondPostRunStep(step.InstanceStep{"PostRun2", func() error {
		fmt.Printf("this is PostRun2\n")
		gotPostStrings = append(gotPostStrings, "this is PostRun2")
		return nil
	}})
	d.RegUserThirdPostRunStep(step.InstanceStep{"PostRun3", func() error {
		fmt.Printf("this is PostRun3\n")
		gotPostStrings = append(gotPostStrings, "this is PostRun3")
		return nil
	}})
	d.RegUserFourthPostRunStep(step.InstanceStep{"PostRun4", func() error {
		fmt.Printf("this is PostRun4\n")
		gotPostStrings = append(gotPostStrings, "this is PostRun4")
		return nil
	}})
	d.RegUserFifthPostRunStep(step.InstanceStep{"PostRun5", func() error {
		fmt.Printf("this is PostRun5\n")
		gotPostStrings = append(gotPostStrings, "this is PostRun5")
		return nil
	}})
	d.RegUserSixthPostRunStep(step.InstanceStep{"PostRun6", func() error {
		fmt.Printf("this is PostRun6\n")
		gotPostStrings = append(gotPostStrings, "this is PostRun6")
		return nil
	}})
	d.RegUserSeventhPostRunStep(step.InstanceStep{"PostRun7", func() error {
		fmt.Printf("this is PostRun7\n")
		gotPostStrings = append(gotPostStrings, "this is PostRun7")
		return nil
	}})
	d.RegUserEighthPostRunStep(step.InstanceStep{"PostRun8", func() error {
		fmt.Printf("this is PostRun8\n")
		gotPostStrings = append(gotPostStrings, "this is PostRun8")
		return nil
	}})
	d.RegUserNinethPostRunStep(step.InstanceStep{"PostRun9", func() error {
		fmt.Printf("this is PostRun9\n")
		gotPostStrings = append(gotPostStrings, "this is PostRun9")
		return nil
	}})
	d.RegUserTenthPostRunStep(step.InstanceStep{"PostRun10", func() error {
		fmt.Printf("this is PostRun10\n")
		gotPostStrings = append(gotPostStrings, "this is PostRun10")
		return nil
	}})
	d.RegUserFirstPreRunStep(step.InstanceStep{"PreRun1", func() error {
		fmt.Printf("this is PreRun1\n")
		gotPreStrings = append(gotPreStrings, "this is PreRun1")
		return nil
	}})
	d.RegUserSecondPreRunStep(step.InstanceStep{"PreRun2", func() error {
		fmt.Printf("this is PreRun2\n")
		gotPreStrings = append(gotPreStrings, "this is PreRun2")
		return nil
	}})
	d.RegUserThirdPreRunStep(step.InstanceStep{"PreRun3", func() error {
		fmt.Printf("this is PreRun3\n")
		gotPreStrings = append(gotPreStrings, "this is PreRun3")
		return nil
	}})
	d.RegUserFourthPreRunStep(step.InstanceStep{"PreRun4", func() error {
		fmt.Printf("this is PreRun4\n")
		gotPreStrings = append(gotPreStrings, "this is PreRun4")
		return nil
	}})
	d.RegUserFifthPreRunStep(step.InstanceStep{"PreRun5", func() error {
		fmt.Printf("this is PreRun5\n")
		gotPreStrings = append(gotPreStrings, "this is PreRun5")
		return nil
	}})
	d.RegUserSixthPreRunStep(step.InstanceStep{"PreRun6", func() error {
		fmt.Printf("this is PreRun6\n")
		gotPreStrings = append(gotPreStrings, "this is PreRun6")
		return nil
	}})
	d.RegUserSeventhPreRunStep(step.InstanceStep{"PreRun7", func() error {
		fmt.Printf("this is PreRun7\n")
		gotPreStrings = append(gotPreStrings, "this is PreRun7")
		return nil
	}})
	d.RegUserEighthPreRunStep(step.InstanceStep{"PreRun8", func() error {
		fmt.Printf("this is PreRun8\n")
		gotPreStrings = append(gotPreStrings, "this is PreRun8")
		return nil
	}})
	d.RegUserNinethPreRunStep(step.InstanceStep{"PreRun9", func() error {
		fmt.Printf("this is PreRun9\n")
		gotPreStrings = append(gotPreStrings, "this is PreRun9")
		return nil
	}})
	d.RegUserTenthPreRunStep(step.InstanceStep{"PreRun10", func() error {
		fmt.Printf("this is PreRun10\n")
		gotPreStrings = append(gotPreStrings, "this is PreRun10")
		return nil
	}})
	_ = d.PreRunStepsAppend(step.InstanceStep{"userPreA1", func() error {
		gotPreStrings = append(gotPreStrings, "this is userPreA1\n")
		fmt.Printf("this is userPreA1\n")
		return nil
	}})
	_ = d.PreRunStepsAppend(step.InstanceStep{"userPreA2", func() error {
		gotPreStrings = append(gotPreStrings, "this is userPreA2\n")
		fmt.Printf("this is userPreA2\n")
		return nil
	}})
	_ = d.PostRunStepsAppend(step.InstanceStep{"userPostA1", func() error {
		fmt.Printf("this is userPostA1\n")
		gotPostStrings = append(gotPostStrings, "this is userPostA1\n")
		return nil
	}})
	_ = d.PostRunStepsAppend(step.InstanceStep{"userPostA2", func() error {
		fmt.Printf("this is userPostA2\n")
		gotPostStrings = append(gotPostStrings, "this is userPostA2\n")
		return nil
	}})
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

//go:norace
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

func TestAddHealthCmd(t *testing.T) {
	d := NewDio()
	d.addHealthCmd(cmd.GrpcUse)
	d.addHealthCmd(cmd.GinUse)
	d.addHealthCmd(cmd.CtlUse)
}

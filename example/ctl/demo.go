package main

import (
	"context"
	"fmt"
	"github.com/gowins/dionysus/step"
	"time"

	"github.com/gowins/dionysus"
	"github.com/gowins/dionysus/cmd"
)

var (
	green  = string([]byte{27, 91, 51, 50, 109})
	white  = string([]byte{27, 91, 51, 55, 109})
	yellow = string([]byte{27, 91, 51, 51, 109})
	red    = string([]byte{27, 91, 51, 49, 109})
	blue   = string([]byte{27, 91, 51, 52, 109})
	// magenta = string([]byte{27, 91, 51, 53, 109})
	// cyan    = string([]byte{27, 91, 51, 54, 109})
	// color   = []string{green, white, yellow, red, blue, magenta, cyan}
)

func main() {
	d := dionysus.NewDio()
	postSteps := []step.InstanceStep{
		{
			StepName: "PostPrint1", Func: func() error {
			fmt.Println(green, "=========== post 1 =========== ", white)
			return nil
		},
		},
		{
			StepName: "PostPrint2", Func: func() error {
			fmt.Println(green, "=========== post 2 =========== ", white)
			return nil
		},
		},
	}
	preSteps := []step.InstanceStep{
		{
			StepName: "PrePrint1", Func: func() error {
			fmt.Println(green, "=========== pre 1 =========== ", white)
			return nil
		},
		},
		{
			StepName: "PrePrint2", Func: func() error {
			fmt.Println(green, "=========== pre 2 =========== ", white)
			return nil
		},
		},
	}
	// PreRun exec before server start
	_ = d.PreRunStepsAppend(preSteps...)

	ctlCmd := cmd.NewCtlCommand()
	_ = ctlCmd.RegRunFunc(func() error {
		timer1 := time.NewTicker(5 * time.Second)
		for {
			select {
			case <-timer1.C:
				fmt.Printf("this is RunE %v\n", time.Now().String())
			case <-ctlCmd.Ctx.Done():
				fmt.Printf("this is stopChan %v\n", time.Now().String())
				return nil
			}
		}
	})

	ctx, cancel := context.WithCancel(ctlCmd.Ctx)
	ctlCmd.Ctx = ctx
	_ = ctlCmd.RegShutdownFunc(func() {
		// ctl需要自己定义Shutdown的函数
		cancel()
	})

	// PostRun exec after server stop
	_ = d.PostRunStepsAppend(postSteps...)

	if err := d.DioStart("ctldemo", ctlCmd); err != nil {
		fmt.Printf("dio start error %v\n", err)
	}
}

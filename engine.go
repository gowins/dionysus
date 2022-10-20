package dionysus

import (
	"fmt"
	logger "log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gowins/dionysus/cmd"
	//logger "github.com/gowins/dionysus/log"
	"github.com/gowins/dionysus/step"
	"github.com/spf13/cobra"
)

type Dio struct {
	cmd                *cobra.Command
	persistentPreRunE  *step.Steps
	persistentPostRunE *step.Steps
}

var defaultDio = NewDio()

func NewDio() *Dio {
	d := &Dio{
		cmd:                &cobra.Command{Use: "root", Short: "just for root"},
		persistentPreRunE:  step.New(),
		persistentPostRunE: step.New(),
	}
	return d
}

// DioStart be care cmds should not use PersistentXXXRunX，this is use by Dio root cmd
func (d *Dio) DioStart(projectName string, cmds ...cmd.Commander) error {
	if projectName == "" {
		return fmt.Errorf("projectName can not be nil")
	}

	// global pre run function
	d.cmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		//logger.Setup()
		// TODO register logger, conf, tracing, metric
		//fmt.Printf("11111")
		d.persistentPreRunE.RegActionSteps("logger", 1, func() error {
			fmt.Printf("init logger here")
			return nil
		})
		d.persistentPreRunE.RegActionSteps("conf", 2, func() error {
			fmt.Printf("init conf here")
			return nil
		})
		d.persistentPreRunE.RegActionSteps("tracing", 3, func() error {
			fmt.Printf("init tracing here")
			return nil
		})
		d.persistentPreRunE.RegActionSteps("metric", 4, func() error {
			fmt.Printf("init metric here")
			return nil
		})
		return d.persistentPreRunE.Run()
	}

	// global post run function
	d.cmd.PersistentPostRunE = func(cmd *cobra.Command, args []string) error {
		// TODO shutdown tracing, metric
		d.persistentPostRunE.RegActionSteps("tracing", 1, func() error {
			fmt.Printf("shutdown tracing here")
			return nil
		})
		d.persistentPostRunE.RegActionSteps("metric", 2, func() error {
			fmt.Printf("shutdown metric here")
			return nil
		})
		return d.persistentPostRunE.Run()
	}

	// append other cmds
	for _, c := range cmds {
		originCmd := c.GetCmd()
		if originCmd.RunE != nil {
			originCmd.RunE = WrapCobrCmdRunE(originCmd.RunE, c.GetShutdownFunc())
		} else if originCmd.Run != nil {
			originCmd.Run = WrapCobrCmdRun(originCmd.Run, c.GetShutdownFunc())
		}
		d.cmd.AddCommand(originCmd)
	}

	// start
	if err := d.cmd.Execute(); err != nil {
		panic(err)
	}
	return nil
}

// PreRunStepsAppend append step will exec after step with priority which define by func PreRunRegWithPriority
func (d *Dio) PreRunStepsAppend(value string, fn func() error) error {
	return d.persistentPreRunE.ActionStepsAppend(value, fn)
}

// PostRunStepsAppend append step will exec after step with priority which define by func PostRunRegWithPriority
func (d *Dio) PostRunStepsAppend(value string, fn func() error) error {
	return d.persistentPostRunE.ActionStepsAppend(value, fn)
}

func (d *Dio) PreRunRegWithPriority(value string, priority int, fn func() error) error {
	if priority > 10000 {
		return fmt.Errorf("priority can not bigger than 10000")
	}
	// priority < 100 is reserve for system steps
	return d.persistentPreRunE.RegActionStepsE(value, priority+100, fn)
}

func (d *Dio) PostRunRegWithPriority(value string, priority int, fn func() error) error {
	if priority > 10000 {
		return fmt.Errorf("priority can not bigger than 10000")
	}
	// priority < 100 is reserve for system steps
	return d.persistentPostRunE.RegActionStepsE(value, priority+100, fn)
}

// Deprecated:: Use DioStart
func Start(project string, cmds ...cmd.Commander) {
	if err := defaultDio.DioStart(project, cmds...); err != nil {
		panic(err)
	}
}

type CobraRunE func(cmd *cobra.Command, args []string) error

func WrapCobrCmdRunE(cobraRunE CobraRunE, shutdownFunc func()) CobraRunE {
	finishChan := make(chan struct{})
	return func(cmd *cobra.Command, args []string) error {
		go func() {
			defer func() {
				if r := recover(); r != nil {
					logger.Printf("[error] Panic occurred in start process: %#v\n", r)
				}
				finishChan <- struct{}{}
			}()
			err := cobraRunE(cmd, args)
			if err != nil {
				logger.Printf("cobra cmd rune error %v\n", err)
			}
		}()
		// TODO health check start
		WaitingForNotifies(finishChan, shutdownFunc)
		return nil
	}
}

type CobraRun func(cmd *cobra.Command, args []string)

func WrapCobrCmdRun(cobraRun CobraRun, shutdownFunc func()) CobraRun {
	finishChan := make(chan struct{})
	return func(cmd *cobra.Command, args []string) {
		go func() {
			defer func() {
				if r := recover(); r != nil {
					logger.Printf("[error] Panic occurred in start process: %#v\n", r)
				}
				finishChan <- struct{}{}
			}()
			cobraRun(cmd, args)
		}()
		// TODO health check start
		WaitingForNotifies(finishChan, shutdownFunc)
	}
}

func WaitingForNotifies(finishChan <-chan struct{}, shutdownFunc func()) {
	quit := make(chan os.Signal)
	signal.Ignore(syscall.SIGHUP)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer func() {
		if r := recover(); r != nil {
			logger.Printf("[error] Panic occurred in shutdown process: %s\n", r)
			os.Exit(3)
		}
	}()

	select {
	case <-quit:
		logger.Printf("[info] Shuting down ...\n")
		if shutdownFunc == nil {
			fmt.Printf("shutdownFunc is nil")
		} else {
			shutdownFunc()
		}
	case <-finishChan:
		logger.Printf("[Dio] Finish.\n")
	}

	logger.Printf("[Dio] Exited.\n")
}

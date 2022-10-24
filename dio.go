package dionysus

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gowins/dionysus/cmd"
	"github.com/gowins/dionysus/log"
	"github.com/gowins/dionysus/step"
	"github.com/spf13/cobra"
)

type Dio struct {
	cmd                *cobra.Command
	persistentPreRunE  *step.Steps
	persistentPostRunE *step.Steps
}

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
		d.persistentPreRunE.RegSysFirstSteps(step.InstanceStep{
			StepName: "logger", Func: func() error {
				log.Setup(log.SetProjectName(projectName), log.WithWriter(os.Stdout))
				log.Info("add init logger here")
				return nil
			},
		})
		d.persistentPreRunE.RegSysSecondSteps(step.InstanceStep{
			StepName: "conf", Func: func() error {
				log.Info("add init conf here")
				return nil
			}})
		d.persistentPreRunE.RegSysThirdSteps(step.InstanceStep{
			StepName: "tracing", Func: func() error {
				log.Info("add init tracing here")
				return nil
			}})
		d.persistentPreRunE.RegSysFourthSteps(step.InstanceStep{
			StepName: "metric", Func: func() error {
				log.Info("add init metric here")
				return nil
			}})
		return d.persistentPreRunE.Run()
	}

	// global post run function
	d.cmd.PersistentPostRunE = func(cmd *cobra.Command, args []string) error {
		d.persistentPostRunE.RegSysFirstSteps(step.InstanceStep{
			StepName: "tracing", Func: func() error {
				log.Info("add shutdown tracing here")
				return nil
			}})
		d.persistentPostRunE.RegSysSecondSteps(step.InstanceStep{
			StepName: "metric", Func: func() error {
				log.Info("add shutdown metric here")
				return nil
			}})
		return d.persistentPostRunE.Run()
	}

	// append other cmds
	for _, c := range cmds {
		originCmd := c.GetCmd()
		if originCmd == nil {
			log.Info("cmd can not be nil")
			return fmt.Errorf("cmd can not be nil")
		}
		originCmd.RunE = wrapCobrCmdRun(originCmd.RunE, c.GetShutdownFunc())
		d.cmd.AddCommand(originCmd)
	}

	return d.cmd.Execute()
}

// Deprecated:: Use DioStart
var defaultDio = NewDio()

// Deprecated:: Use DioStart
func Start(project string, cmds ...cmd.Commander) {
	if err := defaultDio.DioStart(project, cmds...); err != nil {
		panic(err)
	}
}

type CobraRun func(cmd *cobra.Command, args []string) error

func wrapCobrCmdRun(cobraRun CobraRun, shutdownFunc func()) CobraRun {
	finishChan := make(chan struct{})
	return func(cmd *cobra.Command, args []string) error {
		go func() {
			defer func() {
				if r := recover(); r != nil {
					log.Info("[error] Panic occurred in start process: %#v\n", r)
				}
				finishChan <- struct{}{}
			}()
			err := cobraRun(cmd, args)
			if err != nil {
				log.Info("cobra cmd rune error %v\n", err)
			}
		}()
		// TODO health check start
		waitingForNotifies(finishChan, shutdownFunc)
		return nil
	}
}

// WaitingForNotifies todo shutdown 重复
func waitingForNotifies(finishChan <-chan struct{}, shutdownFunc func()) {
	quit := make(chan os.Signal)
	signal.Ignore(syscall.SIGHUP)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer func() {
		if r := recover(); r != nil {
			log.Infof("[error] Panic occurred in shutdown process: %s\n", r)
			os.Exit(3)
		}
	}()

	select {
	case <-quit:
		log.Info("[info] Shuting down ...\n")
		if shutdownFunc == nil {
			log.Info("shutdownFunc is nil")
		} else {
			shutdownFunc()
		}
	case <-finishChan:
		log.Info("[Dio] Finish.\n")
	}

	log.Info("[Dio] Exited.\n")
}

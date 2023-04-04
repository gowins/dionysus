package cmd

import (
	"context"
	"errors"
	"github.com/gowins/dionysus/healthy"
	"github.com/gowins/dionysus/log"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const CtlUse = "ctl"

type ctl struct {
	cmd           *cobra.Command
	Ctx           context.Context
	health        *healthy.Health
	runFunc       func() error
	shutdownSteps []StopStep
}

func NewCtlCommand() *ctl {
	return &ctl{
		cmd:           &cobra.Command{Use: CtlUse, Short: "Run as ctl mod"},
		health:        healthy.New(),
		Ctx:           context.TODO(),
		shutdownSteps: []StopStep{},
	}
}

func (c *ctl) RegShutdownFunc(stopSteps ...StopStep) {
	c.shutdownSteps = append(c.shutdownSteps, stopSteps...)
}

func (c *ctl) RegRunFunc(f func() error) error {
	if f == nil {
		return errors.New("Registering nil func ")
	}
	c.runFunc = func() error {
		defer func() {
			if r := recover(); r != nil {
				log.Errorf("dionysus ctl panic %v", r)
			}
		}()
		//run health
		if err := c.health.FileObserve(healthy.CheckInterval); err != nil {
			log.Errorf("health check error %v\n", err)
			return err
		}
		return f()
	}
	return nil
}

func (c *ctl) RegFlagSet(f *pflag.FlagSet) {
	c.cmd.Flags().AddFlagSet(f)
}

func (c *ctl) Flags() *pflag.FlagSet {
	return c.cmd.Flags()
}

func (c *ctl) GetCmd() *cobra.Command {
	c.cmd.RunE = func(cmd *cobra.Command, args []string) error {
		return c.runFunc()
	}
	return c.cmd
}

func (c *ctl) GetShutdownFunc() StopFunc {
	return func() {
		for _, stopSteps := range c.shutdownSteps {
			log.Infof("run shutdown %v", stopSteps.StepName)
			stopSteps.StopFn()
		}
	}
}

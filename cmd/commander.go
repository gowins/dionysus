package cmd

import (
	"github.com/spf13/cobra"
)

type Priority = int

type StopFunc func()

type StopStep struct {
	StopFn   StopFunc
	StepName string
}

type Commander interface {
	GetCmd() *cobra.Command
	GetShutdownFunc() StopFunc
	RegShutdownFunc(stopSteps ...StopStep)
}

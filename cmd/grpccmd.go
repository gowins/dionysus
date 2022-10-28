package cmd

import (
	"github.com/gowins/dionysus/log"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const (
	GrpcUse = "grpc"
)

var defaultGrpcAddr = ":1234"

type grpcCommand struct {
	cmd           *cobra.Command
	shutdownSteps []StopStep
}

func NewGrpcCommand() *grpcCommand {
	return &grpcCommand{
		cmd:           &cobra.Command{Use: GrpcUse, Short: "Run as grpc server"},
		shutdownSteps: []StopStep{},
	}
}

func (g *grpcCommand) GetCmd() *cobra.Command {
	g.cmd.RunE = func(cmd *cobra.Command, args []string) error {
		return nil
	}
	return g.cmd
}

func (g *grpcCommand) RegShutdownFunc(stopSteps ...StopStep) {
	g.shutdownSteps = append(g.shutdownSteps, stopSteps...)
}

func (g *grpcCommand) GetShutdownFunc() StopFunc {
	return func() {
		for _, stopSteps := range g.shutdownSteps {
			log.Infof("run shutdown %v", stopSteps.StepName)
			stopSteps.StopFn()
		}
		//grpcServer.GracefulStop()
	}
}

func (g *grpcCommand) RegFlagSet(set *pflag.FlagSet) {
}
func (g *grpcCommand) Flags() *pflag.FlagSet {
	return &pflag.FlagSet{}
}

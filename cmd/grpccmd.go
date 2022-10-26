package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const (
	GrpcUse = "grpc"
)

var defaultGrpcAddr = ":1234"

type grpcCommand struct {
	cmd *cobra.Command
}

func NewGrpcCommand() *grpcCommand {
	return &grpcCommand{
		cmd: &cobra.Command{Use: GrpcUse, Short: "Run as grpc server"},
	}
}

func (g *grpcCommand) GetCmd() *cobra.Command {
	g.cmd.RunE = func(cmd *cobra.Command, args []string) error {
		return nil
	}
	return g.cmd
}

func (g *grpcCommand) GetShutdownFunc() func() {
	return func() {
	}
}

func (g *grpcCommand) RegFlagSet(set *pflag.FlagSet) {
}
func (g *grpcCommand) Flags() *pflag.FlagSet {
	return &pflag.FlagSet{}
}

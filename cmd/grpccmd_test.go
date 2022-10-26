package cmd

import (
	"github.com/spf13/pflag"
	"testing"
)

func TestNewGrpcCommand(t *testing.T) {
	grpcCmd := NewGrpcCommand()
	grpcCmd.RegFlagSet(&pflag.FlagSet{})
	grpcCmd.Flags()
	shutdownFunc := grpcCmd.GetShutdownFunc()
	shutdownFunc()
	cmd := grpcCmd.GetCmd()
	cmd.RunE(nil, nil)
}

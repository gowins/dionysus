package cmd

import (
	"fmt"
	"github.com/spf13/pflag"
	"testing"
)

func TestNewGrpcCommand(t *testing.T) {
	grpcCmd := NewGrpcCommand()
	grpcCmd.RegFlagSet(&pflag.FlagSet{})
	grpcCmd.Flags()
	grpcCmd.RegShutdownFunc(StopStep{
		StepName: "grpc stop",
		StopFn: func() {
			fmt.Printf("this is grpc stop")
		},
	})
	shutdownFunc := grpcCmd.GetShutdownFunc()
	shutdownFunc()
	cmd := grpcCmd.GetCmd()
	cmd.RunE(nil, nil)
}

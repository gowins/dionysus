package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type Priority = int

type Commander interface {
	GetCmd() *cobra.Command
	GetShutdownFunc() func()

	RegFlagSet(set *pflag.FlagSet)
	Flags() *pflag.FlagSet
}

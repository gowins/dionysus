package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type Priority = int

type Commander interface {
	GetCmd() *cobra.Command

	RegFlagSet(set *pflag.FlagSet)
	Flags() *pflag.FlagSet

	RegPreRunFunc(value string, f func() error) error
	RegPostRunFunc(value string, f func() error) error
}

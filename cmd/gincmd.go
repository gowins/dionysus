package cmd

import (
	"context"
	"net/http"
	"os"
	"sync"

	"github.com/gowins/dionysus/ginx"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const (
	WebServerAddr = "GAPI_ADDR"
	addrFlagName  = "addr"
)

var (
	defaultWebServerAddr = ":8080"
)

type ginCommand struct {
	ginx.ZeroGinRouter

	cmd    *cobra.Command
	server *http.Server
	addr   string

	once sync.Once
}

func NewGinCommand() *ginCommand {
	return &ginCommand{
		ZeroGinRouter: ginx.NewZeroGinRouter(),
		cmd:           &cobra.Command{Use: "gin", Short: "Run as go-zero server"},
		server:        &http.Server{},
	}
}

func (t *ginCommand) GetCmd() *cobra.Command {
	t.once.Do(func() {
		if envAddr := os.Getenv(WebServerAddr); envAddr != "" {
			defaultWebServerAddr = envAddr
		}
		t.cmd.Flags().StringVarP(&t.addr, addrFlagName, "a", defaultWebServerAddr, "the http server address")

		t.server.Handler = t.Handler()
	})

	t.cmd.RunE = func(cmd *cobra.Command, args []string) error {
		t.server.Addr = t.addr
		return t.server.ListenAndServe()
	}
	return t.cmd
}

func (t *ginCommand) Start() {
	if err := t.cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func (t *ginCommand) RegFlagSet(set *pflag.FlagSet) {
	t.cmd.Flags().AddFlagSet(set)
}

func (t *ginCommand) Flags() *pflag.FlagSet {
	return t.cmd.Flags()
}

func (t *ginCommand) RegPreRunFunc(value string, f func() error) error {
	t.cmd.PreRunE = func(cmd *cobra.Command, args []string) error {
		return f()
	}
	return nil
}

func (t *ginCommand) RegPostRunFunc(value string, f func() error) error {
	t.cmd.PostRunE = func(cmd *cobra.Command, args []string) error {
		if err := f(); err != nil {
			return err
		}
		if err := t.server.Shutdown(context.TODO()); err != nil {
			return err
		}
		return nil
	}
	return nil
}

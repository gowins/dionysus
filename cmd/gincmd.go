package cmd

import (
	"context"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/gowins/dionysus/ginx"
	"github.com/gowins/dionysus/shutdown"
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
	cmd        *cobra.Command
	server     *http.Server
	addr       string
	once       sync.Once
	finishChan chan struct{}

	preRun  []func() error
	postRun []func() error
}

func NewGinCommand() *ginCommand {
	return &ginCommand{
		ZeroGinRouter: ginx.NewZeroGinRouter(),
		cmd:           &cobra.Command{Use: "gin", Short: "Run as go-zero server"},
		server:        &http.Server{},
		finishChan:    make(chan struct{}),
	}
}

func (t *ginCommand) GetCmd() *cobra.Command {
	t.once.Do(func() {
		if envAddr := os.Getenv(WebServerAddr); envAddr != "" {
			defaultWebServerAddr = envAddr
		}
		t.cmd.Flags().StringVarP(&t.server.Addr, addrFlagName, "a", defaultWebServerAddr, "the http server address")
	})

	t.cmd.PreRunE = func(cmd *cobra.Command, args []string) error {
		t.server.Handler = t.Handler()
		for _, v := range t.preRun {
			if err := v(); err != nil {
				return err
			}
		}
		return nil
	}

	t.cmd.RunE = func(cmd *cobra.Command, args []string) error {
		shutdown.NotifyAfterFinish(t.finishChan, t.startServer)
		return nil
	}

	t.cmd.PostRunE = func(cmd *cobra.Command, args []string) error {
		shutdown.WaitingForNotifies(t.finishChan, t.stopServer)
		for _, v := range t.postRun {
			if err := v(); err != nil {
				return err
			}
		}
		return nil
	}

	return t.cmd
}

func (t *ginCommand) RegFlagSet(set *pflag.FlagSet) {
	t.cmd.Flags().AddFlagSet(set)
}

func (t *ginCommand) Flags() *pflag.FlagSet {
	return t.cmd.Flags()
}

func (t *ginCommand) RegPreRunFunc(value string, f func() error) error {
	t.preRun = append(t.preRun, f)
	return nil
}

func (t *ginCommand) RegPostRunFunc(value string, f func() error) error {
	t.postRun = append(t.postRun, f)
	return nil
}

func (g *ginCommand) startServer() {
	log.Printf("[Dio] Engine setting with address %v", g.server.Addr)
	if err := g.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Printf("listen: %s\n", err)
		os.Exit(1)
	}
}

func (g *ginCommand) stopServer() {
	log.Println("[info] Server exiting")
	if err := g.server.Shutdown(context.TODO()); err != nil {
		log.Println("[error] Server forced to shutdown:", err)
		os.Exit(1)
	}
}

package cmd

import (
	"errors"
	"fmt"

	"github.com/gowins/dionysus/shutdown"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// The lifecycle of ctlcmd and four user functions: preFunc runFunc shutdownFunc and postFunc
//
// prerun stage                     run stage                       post run stage
// +-----------------+              +-------------------+        +--------------------------------------------------+
// |                 |              |                   |        |                                                  |
// | +-------------+ |              |  +--------------+ |        |  +-----------------+                             |
// | | parse flag  | |              |  |              +------------>+ stuck at select +-----------------+           |
// | +-----+-------+ |              |  |              | |        |  +--------+--------+                 |           |
// |       |         |       +-------->+ go (runFunc) | |        |           +                          +           |
// |       |         |       |      |  |              | |        |        os.Signal                user finish      |
// |  flags|in ctx   |       |      |  |              | |        |           +                          |           |
// |       |         |       |      |  +--------------+ |        |           |                          |           |
// |       |         |       |      |                   |        |           v                          v           |
// |       v         |       |      |                   |        |  +--------+--------+  succeed  +-----+--------+  |
// | +-----+-------+ |       |      |                   |        |  |  shutdownFunc   +---------->+  postFunc    |  |
// | |  preFunc    | |       |      |                   |        |  +--------+--------+           +-----+--------+  |
// | +-----+-------+ |       |      |                   |        |           |                          |           |
// |       |         |       |      |                   |        |           |                          |code=0     |
// |       v         |       |      |                   |        |           |                          v           |
// | +-----+-------+ |       |      |                   |        |           |  timeout code=1    +-----+--------+  |
// | | go (healthy)+---------+      |                   |        |           +------------------->+ os.Exit(code)|  |
// | +-------------+ |              |                   |        |                                +--------------+  |
// |                 |              |                   |        |                                                  |
// +-----------------+              +-------------------+        +--------------------------------------------------+

// CtxKey fix golint stage goanalysis_metalinter error
// SA1029: should not use built-in type string as key for value; define your own type to avoid collisio
type CtxKey string

// var runFunc, shutdownFunc = func(ctx context.Context) {}, func(ctx context.Context) {}
// var userFlagSet = &pflag.FlagSet{}

type ctl struct {
	cmd *cobra.Command

	runFunc, shutdownFunc func() error
	preFunc, postFunc     []func() error
}

func NewCtlCommand() *ctl {
	return &ctl{
		cmd: &cobra.Command{Use: "ctl", Short: "Run as ctl mod"},
	}
}

func (c *ctl) RegPreRunFunc(value string, f func() error) error {
	if f == nil {
		return errors.New("Registering nil func ")
	}
	c.preFunc = append(c.preFunc, f)
	return nil
}

func (c *ctl) RegRunFunc(f func() error) error {
	if f == nil {
		return errors.New("Registering nil func ")
	}
	c.runFunc = f
	return nil
}

func (c *ctl) RegShutdownFunc(f func() error) error {
	if f == nil {
		return errors.New("Registering nil func ")
	}
	c.shutdownFunc = f
	return nil
}

func (c *ctl) RegPostRunFunc(value string, f func() error) error {
	if f == nil {
		return errors.New("Registering nil func ")
	}
	c.postFunc = append(c.postFunc, f)
	return nil
}

func (c *ctl) RegFlagSet(f *pflag.FlagSet) {
	c.cmd.Flags().AddFlagSet(f)
}

func (c *ctl) Flags() *pflag.FlagSet {
	return c.cmd.Flags()
}

func (c *ctl) GetCmd() *cobra.Command {
	finishChan := make(chan struct{})

	c.cmd.PreRunE = func(cmd *cobra.Command, args []string) (err error) {
		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("Recovered from prerun. Err:%v ", r)
			}
		}()

		if len(c.preFunc) > 0 {
			for _, f := range c.preFunc {
				if err = f(); err != nil {
					return err
				}
			}
		}

		return nil
	}

	c.cmd.Run = func(cmd *cobra.Command, args []string) {
		shutdown.NotifyAfterFinish(finishChan, func() {
			_ = c.runFunc()
		})

	}

	c.cmd.PostRunE = func(cmd *cobra.Command, args []string) error {
		if len(c.postFunc) > 0 {
			for _, f := range c.postFunc {
				if err := f(); err != nil {
					return err
				}
			}
		}
		shutdown.WaitingForNotifies(finishChan, func() {
			_ = c.shutdownFunc()
		})
		return nil
	}

	return c.cmd
}

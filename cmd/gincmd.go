package cmd

import (
	"context"
	"github.com/gowins/dionysus/opentelemetry"
	"net/http"
	"os"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gowins/dionysus/ginx"
	"github.com/gowins/dionysus/healthy"
	"github.com/gowins/dionysus/log"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const (
	WebServerAddr = "GAPI_ADDR"
	closePath     = "/close"
	openPath      = "/open"
	GinUse        = "gin"
)

var (
	defaultWebServerAddr = ":8080"
	livenessStatus       = true
	readinessStatus      = true
	startupStatus        = true
)

type ginCommand struct {
	ginx.ZeroGinRouter
	cmd           *cobra.Command
	server        *http.Server
	addr          string
	once          sync.Once
	shutdownSteps []StopStep
}

func NewGinCommand(opts ...ginx.GinOption) *ginCommand {
	return &ginCommand{
		ZeroGinRouter: ginx.NewZeroGinRouter(opts...),
		cmd:           &cobra.Command{Use: GinUse, Short: "Run as go-zero server"},
		server:        &http.Server{},
		shutdownSteps: []StopStep{},
	}
}

func (t *ginCommand) GetCmd() *cobra.Command {
	t.once.Do(func() {
		if envAddr := os.Getenv(WebServerAddr); envAddr != "" {
			defaultWebServerAddr = envAddr
		}
	})

	t.server.Addr = defaultWebServerAddr
	t.cmd.RunE = func(cmd *cobra.Command, args []string) error {
		t.registerHealth()
		t.startServer()
		return nil
	}

	return t.cmd
}

func (t *ginCommand) RegShutdownFunc(stopSteps ...StopStep) {
	t.shutdownSteps = append(t.shutdownSteps, stopSteps...)
}

func (t *ginCommand) RegFlagSet(set *pflag.FlagSet) {
	t.cmd.Flags().AddFlagSet(set)
}

func (t *ginCommand) Flags() *pflag.FlagSet {
	return t.cmd.Flags()
}

func (g *ginCommand) startServer() {
	log.Infof("[Dio] Engine setting with address %v", g.server.Addr)
	g.server.Handler = g.Handler()
	if opentelemetry.TracerIsEnable() {
		log.Infof("[Dio] Engine use opentelemetry trace")
		g.server.Handler = opentelemetry.InitHttpHandler(g.server.Handler)

	}

	if err := g.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Infof("listen: %s\n", err)
		os.Exit(1)
	}
}

func (g *ginCommand) registerHealth() {
	healthxGroup := g.Group(healthy.HealthGroupPath)
	healthxGroup.Handle(http.MethodGet, healthy.HealthLivenessPath, func(c *gin.Context) ginx.Render {
		if err := healthy.CheckerFuncRun(healthy.HealthLiveness); err != nil {
			c.JSON(http.StatusInternalServerError, "checker error "+err.Error())
			return nil
		}
		if !livenessStatus {
			c.JSON(http.StatusInternalServerError, "livenesss is closed")
			return nil
		}
		return ginx.Success("ok")
	})
	healthxGroup.Handle(http.MethodPost, healthy.HealthLivenessPath+closePath, func(c *gin.Context) ginx.Render {
		livenessStatus = false
		return ginx.Success("ok")
	})
	healthxGroup.Handle(http.MethodPost, healthy.HealthLivenessPath+openPath, func(c *gin.Context) ginx.Render {
		livenessStatus = true
		return ginx.Success("ok")
	})
	healthxGroup.Handle(http.MethodGet, healthy.HealthReadinessPath, func(c *gin.Context) ginx.Render {
		if err := healthy.CheckerFuncRun(healthy.HealthReadiness); err != nil {
			c.JSON(http.StatusInternalServerError, "checker error "+err.Error())
			return nil
		}
		if !readinessStatus {
			c.JSON(http.StatusInternalServerError, "readiness is closed")
			return nil
		}
		return ginx.Success("ok")
	})
	healthxGroup.Handle(http.MethodPost, healthy.HealthReadinessPath+closePath, func(c *gin.Context) ginx.Render {
		readinessStatus = false
		return ginx.Success("ok")
	})
	healthxGroup.Handle(http.MethodPost, healthy.HealthReadinessPath+openPath, func(c *gin.Context) ginx.Render {
		readinessStatus = true
		return ginx.Success("ok")
	})
	healthxGroup.Handle(http.MethodGet, healthy.HealthStartupPath, func(c *gin.Context) ginx.Render {
		if err := healthy.CheckerFuncRun(healthy.HealthStartup); err != nil {
			c.JSON(http.StatusInternalServerError, "checker error "+err.Error())
			return nil
		}
		if !startupStatus {
			c.JSON(http.StatusInternalServerError, "startup is closed")
			return nil
		}
		return ginx.Success("ok")
	})
	healthxGroup.Handle(http.MethodPost, healthy.HealthStartupPath+closePath, func(c *gin.Context) ginx.Render {
		startupStatus = false
		return ginx.Success("ok")
	})
	healthxGroup.Handle(http.MethodPost, healthy.HealthStartupPath+openPath, func(c *gin.Context) ginx.Render {
		startupStatus = true
		return ginx.Success("ok")
	})
}

func (g *ginCommand) stopServer() {
	log.Infof("[info] Server exiting")
	if err := g.server.Shutdown(context.TODO()); err != nil {
		log.Infof("[error] Server forced to shutdown:", err)
		os.Exit(1)
	}
}

func (g *ginCommand) GetShutdownFunc() StopFunc {
	return func() {
		for _, stopSteps := range g.shutdownSteps {
			log.Infof("run shutdown %v", stopSteps.StepName)
			stopSteps.StopFn()
		}
		g.stopServer()
	}
}

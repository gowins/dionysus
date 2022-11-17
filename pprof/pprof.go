package pprof

import (
	"github.com/gin-gonic/gin"
	"github.com/gowins/dionysus/ginx"
	"log"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"syscall"
)

var pprofAddr = ":9092"

func Setup() {
	go func() {
		pprofSignal := make(chan os.Signal, 1)
		signal.Notify(pprofSignal, syscall.SIGUSR2)
		select {
		case <-pprofSignal:
			log.Printf("start pprof")
			setup()
		}
	}()
}

func setup() {
	r := ginx.NewZeroGinRouter()
	pprofGroup := r.Group("/debug/pprof")
	pprofGroup.Handle(http.MethodGet, "/", func(c *gin.Context) ginx.Render {
		pprof.Index(c.Writer, c.Request)
		return nil
	})
	pprofGroup.Handle(http.MethodGet, "/cmdline", func(c *gin.Context) ginx.Render {
		pprof.Cmdline(c.Writer, c.Request)
		return nil
	})
	pprofGroup.Handle(http.MethodGet, "/profile", func(c *gin.Context) ginx.Render {
		pprof.Profile(c.Writer, c.Request)
		return nil
	})
	pprofGroup.Handle(http.MethodGet, "/symbol", func(c *gin.Context) ginx.Render {
		pprof.Symbol(c.Writer, c.Request)
		return nil
	})
	pprofGroup.Handle(http.MethodGet, "/trace", func(c *gin.Context) ginx.Render {
		pprof.Trace(c.Writer, c.Request)
		return nil
	})
	pprofGroup.Handle(http.MethodGet, "/allocs", func(c *gin.Context) ginx.Render {
		pprof.Handler("allocs").ServeHTTP(c.Writer, c.Request)
		return nil
	})
	pprofGroup.Handle(http.MethodGet, "/block", func(c *gin.Context) ginx.Render {
		pprof.Handler("block").ServeHTTP(c.Writer, c.Request)
		return nil
	})
	pprofGroup.Handle(http.MethodGet, "/goroutine", func(c *gin.Context) ginx.Render {
		pprof.Handler("goroutine").ServeHTTP(c.Writer, c.Request)
		return nil
	})
	pprofGroup.Handle(http.MethodGet, "/heap", func(c *gin.Context) ginx.Render {
		pprof.Handler("heap").ServeHTTP(c.Writer, c.Request)
		return nil
	})
	pprofGroup.Handle(http.MethodGet, "/mutex", func(c *gin.Context) ginx.Render {
		pprof.Handler("mutex").ServeHTTP(c.Writer, c.Request)
		return nil
	})
	pprofGroup.Handle(http.MethodGet, "/threadcreate", func(c *gin.Context) ginx.Render {
		pprof.Handler("threadcreate").ServeHTTP(c.Writer, c.Request)
		return nil
	})
	go func() {
		_ = http.ListenAndServe(pprofAddr, r.Handler())
	}()
}

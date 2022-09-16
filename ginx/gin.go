package ginx

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ZeroGinRouter interface {
	GinRouters
	Group(path string, handler ...GinHandler) ZeroGinRouter
}

type GinRouter struct {
	ginGroup

	engine *gin.Engine
	server *http.Server
	root   bool
}

func NewZeroGinRouter() *GinRouter {
	g := gin.New()
	g.Use(gin.Recovery()) // 默认注册recovery
	r := &GinRouter{
		ginGroup: ginGroup{g: &g.RouterGroup},
		engine:   g,
		server: &http.Server{
			Handler: g.Handler(),
		},
		root: true,
	}
	return r
}

func (r *GinRouter) Group(path string, handler ...GinHandler) ZeroGinRouter {
	g := r.engine.Group(path, buildGinHandler(handler...)...)
	return &GinRouter{
		ginGroup: ginGroup{
			g: g,
		},
		engine: r.engine,
	}
}

// Run 启动
func (r *GinRouter) Run(addr string) error {
	r.server.Addr = addr
	if err := r.server.ListenAndServe(); err != nil {
		return err
	}
	return nil
}

// Shutdown 停止
func (r *GinRouter) Shutdown() error {
	if err := r.server.Shutdown(context.Background()); err != nil {
		return err
	}
	return nil
}

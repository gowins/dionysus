package ginx

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type ZeroGinRouter interface {
	GinRouters
	Group(path string, handler ...GinHandler) ZeroGinRouter
	Handler() http.Handler
}

type ginRouter struct {
	ginGroup

	engine *gin.Engine
}

func NewZeroGinRouter() ZeroGinRouter {
	g := gin.New()
	g.Use(gin.Recovery()) // 默认注册recovery
	r := &ginRouter{
		ginGroup: ginGroup{g: &g.RouterGroup},
		engine:   g,
	}
	return r
}

func (r *ginRouter) Group(path string, handler ...GinHandler) ZeroGinRouter {
	g := r.engine.Group(path, buildGinHandler(handler...)...)
	return &ginRouter{
		ginGroup: ginGroup{
			g: g,
		},
		engine: r.engine,
	}
}

func (r *ginRouter) Handler() http.Handler {
	return r.engine.Handler()
}

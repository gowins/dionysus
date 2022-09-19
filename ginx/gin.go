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
	group  *gin.RouterGroup
}

func NewZeroGinRouter() ZeroGinRouter {
	g := gin.New()
	g.Use(gin.Recovery()) // 默认注册recovery
	r := &ginRouter{
		ginGroup: ginGroup{g: &g.RouterGroup},
		engine:   g,
		group:    &g.RouterGroup,
	}
	return r
}

func (r *ginRouter) Group(path string, handler ...GinHandler) ZeroGinRouter {
	g := r.group.Group(path, buildGinHandler(handler...)...)
	return &ginRouter{
		ginGroup: ginGroup{
			g: g,
		},
		group:  g,
		engine: r.engine,
	}
}

func (r *ginRouter) Handler() http.Handler {
	return r.engine.Handler()
}

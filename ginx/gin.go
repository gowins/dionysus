package ginx

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ZeroGinRouter interface {
	GinRouters
	Group(path string, handlers ...gin.HandlerFunc) ZeroGinRouter
	Handler() http.Handler
}

type ginRouter struct {
	ginGroup
	engine *gin.Engine
	group  *gin.RouterGroup
}

func NewZeroGinRouter(opts ...GinOption) ZeroGinRouter {
	g := gin.New()
	// set properties of gin.Engine
	for _, opt := range opts {
		opt(g)
	}
	g.Use(gin.Recovery()) // 默认注册recovery
	r := &ginRouter{
		ginGroup: ginGroup{g: &g.RouterGroup},
		engine:   g,
		group:    &g.RouterGroup,
	}
	return r
}

func (r *ginRouter) Group(path string, handlers ...gin.HandlerFunc) ZeroGinRouter {
	g := r.group.Group(path, handlers...)
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

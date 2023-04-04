package ginx

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ZeroGinRouter route manage
type ZeroGinRouter interface {
	GinRouters
	// Group create and return new router group
	Group(path string) ZeroGinRouter
	// Handler return http.Handler
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
	gin.DefaultWriter = &normalWrite{}
	g.Use(gin.RecoveryWithWriter(&panicWrite{})) // 默认注册recovery
	r := &ginRouter{
		ginGroup: ginGroup{g: &g.RouterGroup},
		engine:   g,
		group:    &g.RouterGroup,
	}
	return r
}

// Group create and return new router group
func (r *ginRouter) Group(path string) ZeroGinRouter {
	g := r.group.Group(path)
	return &ginRouter{
		ginGroup: ginGroup{
			g: g,
		},
		group:  g,
		engine: r.engine,
	}
}

// Handler return http.Handler
func (r *ginRouter) Handler() http.Handler {
	return r.engine.Handler()
}

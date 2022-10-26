package ginx

import (
	"github.com/gin-gonic/gin"
)

// GinHandler gin.HandleFunc with Render interface
type GinHandler func(c *gin.Context) Render

// GinRouters register route
type GinRouters interface {
	// Use resgister middleware
	Use(handler ...gin.HandlerFunc) GinRouters
	// Handle use GinHandler to register route
	Handle(method, path string, handler GinHandler) GinRouters
}

type ginGroup struct {
	g gin.IRoutes
}

func (r *ginGroup) Use(handler ...gin.HandlerFunc) GinRouters {
	g := r.g.Use(handler...)
	return &ginGroup{g: g}
}

func (r *ginGroup) Handle(method, path string, handler GinHandler) GinRouters {
	g := r.g.Handle(method, path, buildGinHandler(handler))
	return &ginGroup{g: g}
}

func buildGinHandler(handler GinHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		if r := handler(c); r != nil {
			c.Render(c.Writer.Status(), r)
		}
	}
}

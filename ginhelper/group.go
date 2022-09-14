package ginhelper

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
)

type GinHandler func(c *gin.Context) render.Render

type GinRouters interface {
	Use(handler ...gin.HandlerFunc) GinRouters
	Handle(method, path string, handler ...GinHandler) GinRouters
}

type ginGroup struct {
	g gin.IRoutes
}

func (r *ginGroup) Use(handler ...gin.HandlerFunc) GinRouters {
	g := r.g.Use(handler...)
	return &ginGroup{g: g}
}

func (r *ginGroup) Handle(method, path string, handler ...GinHandler) GinRouters {
	g := r.g.Handle(method, path, buildGinHandler(handler...)...)
	return &ginGroup{g: g}
}

func buildGinHandler(handler ...GinHandler) []gin.HandlerFunc {
	gh := make([]gin.HandlerFunc, len(handler))

	for k, v := range handler {
		gh[k] = func(c *gin.Context) {
			c.Render(http.StatusOK, v(c))
		}
	}

	return gh
}

package ginx

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type GinHandler func(c *gin.Context) Render

type Handler func(c *gin.Context) error

type GinRouters interface {
	Use(handler ...gin.HandlerFunc) GinRouters
	Handle(method, path string, handler ...GinHandler) GinRouters
	HandleE(method, path string, handler ...Handler) GinRouters
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

func (r *ginGroup) HandleE(method, path string, handler ...Handler) GinRouters {
	g := r.g.Handle(method, path, buildGinHandlerE(handler...)...)
	return &ginGroup{g: g}
}

func buildGinHandlerE(handler ...Handler) []gin.HandlerFunc {
	gh := make([]gin.HandlerFunc, len(handler))

	for k, v := range handler {
		gh[k] = func(c *gin.Context) {
			if err := v(c); err != nil {
				c.Render(http.StatusOK, Error(err))
			}
		}
	}

	return gh
}

func buildGinHandler(handler ...GinHandler) []gin.HandlerFunc {
	gh := make([]gin.HandlerFunc, len(handler))

	for k, v := range handler {
		gh[k] = func(c *gin.Context) {
			if r := v(c); r != nil {
				c.Render(http.StatusOK, r)
			}
		}
	}

	return gh
}

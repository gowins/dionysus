package ginhelper

import (
	"net/http"
	"time"

	"github.com/gin-contrib/timeout"
	"github.com/gin-gonic/gin"
)

func TimeoutMiddleware(millSecond int) gin.HandlerFunc {
	return timeout.New(
		timeout.WithTimeout(time.Millisecond*time.Duration(millSecond)),
		timeout.WithHandler(func(c *gin.Context) { c.Next() }),
		timeout.WithResponse(func(c *gin.Context) { c.String(http.StatusGatewayTimeout, "timeout") }),
	)
}

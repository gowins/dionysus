package ginx

import (
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

func LimiterMiddleware(limit int) gin.HandlerFunc {
	limiter := rate.NewLimiter(rate.Every(time.Second), limit)
	return func(c *gin.Context) {
		if !limiter.Allow() {
			c.JSON(200, Response{
				Code: SpeedLimit,
				Msg:  CodeMsgMap[SpeedLimit],
			})
			c.Abort()
		}
		c.Next()
	}
}

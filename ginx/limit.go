package ginx

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
	"net/http"
)

func LimiterMiddleware(limit int) gin.HandlerFunc {
	limiter := rate.NewLimiter(rate.Limit(limit), limit)
	return func(c *gin.Context) {
		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, Response{
				Code: ginxLimitingCode,
				Msg:  ginxLimitingMsg,
			})
			c.Abort()
		}
		c.Next()
	}
}

package ginx

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

func LimiterMiddleware(limit int) gin.HandlerFunc {
	limiter := rate.NewLimiter(rate.Every(time.Second), limit)
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

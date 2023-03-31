package ginx

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

func LimiterMiddleware(limit int) gin.HandlerFunc {
	limiter := rate.NewLimiter(rate.Limit(limit), limit)
	return func(c *gin.Context) {
		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, Response{
				Code: ginxLimitingCode,
				Msg:  ginxLimitingMsg,
			})
			log.Errorf("too many requests, limit open, route %v", c.Request.URL.String())
			c.Abort()
		}
		c.Next()
	}
}

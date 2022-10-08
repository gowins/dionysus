package ginx

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-contrib/timeout"
	"github.com/gin-gonic/gin"
)

func TimeoutMiddleware() gin.HandlerFunc {
	var (
		bufPool        = timeout.BufferPool{}
		defaultTimeout = 10
	)
	return func(c *gin.Context) {
		var (
			requestTimeOut = defaultTimeout
		)
		if t := c.Request.Header.Get("Request_Timeout"); t != "" {
			if s, err := strconv.Atoi(t); err == nil && s < defaultTimeout && s > 0 {
				requestTimeOut = s
			}
		}

		timeoutSecond := time.Duration(requestTimeOut) * time.Second

		finish := make(chan struct{}, 1)
		panicChan := make(chan interface{}, 1)

		w := c.Writer
		buffer := bufPool.Get()
		tw := NewWriter(w, buffer)
		c.Writer = tw
		buffer.Reset()

		go func() {
			defer func() {
				if p := recover(); p != nil {
					panicChan <- p
				}
			}()
			c.Next()
			finish <- struct{}{}
		}()

		select {
		case p := <-panicChan:
			tw.FreeBuffer()
			c.Writer = w
			panic(p)

		case <-finish:
			c.Next()
			tw.mu.Lock()
			defer tw.mu.Unlock()
			dst := tw.ResponseWriter.Header()
			for k, vv := range tw.Header() {
				dst[k] = vv
			}
			tw.ResponseWriter.WriteHeader(tw.code)
			if _, err := tw.ResponseWriter.Write(buffer.Bytes()); err != nil {
				panic(err)
			}
			tw.FreeBuffer()
			bufPool.Put(buffer)

		case <-time.After(timeoutSecond):
			c.Abort()
			tw.mu.Lock()
			defer tw.mu.Unlock()
			tw.timeout = true
			tw.FreeBuffer()
			bufPool.Put(buffer)

			c.Writer = w
			c.String(http.StatusGatewayTimeout, "timeout")
			c.Writer = tw
		}
	}
}

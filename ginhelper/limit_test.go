package ginhelper

import (
	"encoding/json"
	"net/http"
	"sync"
	"testing"

	"github.com/gin-gonic/gin/render"

	"github.com/gin-gonic/gin"

	. "github.com/smartystreets/goconvey/convey"
)

func TestLimiterMiddleware(t *testing.T) {
	Convey("限流", t, func() {
		limit := 10

		r := NewZeroGinRouter()
		r.Use(LimiterMiddleware(limit))
		r.Handle(http.MethodGet, "/limiter", func(c *gin.Context) render.Render {
			return Success(struct{}{})
		})

		wg := sync.WaitGroup{}
		success := 0
		for i := 0; i < limit*2; i++ {
			wg.Add(1)
			go func(i int) {
				w := testHttpRequest("GET", "/limiter", nil, r)
				res := &Response{}
				_ = json.Unmarshal(w.Body.Bytes(), res)
				t.Logf("goroutine:%d response:%d", i, res.Code)
				if res.Code == 200 {
					success++
				}
				wg.Done()
			}(i)
		}
		wg.Wait()
		So(success, ShouldEqual, limit)
	})
}

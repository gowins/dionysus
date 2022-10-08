package ginx

import (
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	. "github.com/smartystreets/goconvey/convey"
)

func TestTimeout(t *testing.T) {
	Convey("超时", t, func() {
		r := NewZeroGinRouter()
		r.Use(TimeoutMiddleware())
		r.Handle(http.MethodGet, "timeout", func(c *gin.Context) Render {
			time.Sleep(2 * time.Second)
			return Success(struct{}{})
		})

		r.Handle(http.MethodGet, "no-timeout", func(c *gin.Context) Render {
			time.Sleep(500 * time.Millisecond)
			return Success(struct{}{})
		})

		res := testHttpRequest("GET", "/timeout", nil, r)
		So(res.Code, ShouldEqual, 504)

		res = testHttpRequest("GET", "/no-timeout", nil, r)
		So(res.Code, ShouldEqual, 200)
	})
}

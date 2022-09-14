package ginhelper

import (
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin/render"

	"github.com/gin-gonic/gin"
	. "github.com/smartystreets/goconvey/convey"
)

func TestTimeout(t *testing.T) {
	Convey("超时", t, func() {
		r := NewZeroGinRouter()
		r.Use(TimeoutMiddleware(1000))
		r.Handle(http.MethodGet, "timeout", func(c *gin.Context) render.Render {
			time.Sleep(2 * time.Second)
			return Success(struct{}{})
		})

		res := testHttpRequest("GET", "/timeout", nil, r)
		So(res.Code, ShouldEqual, 504)
	})
}

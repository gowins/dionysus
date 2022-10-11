package ginx

import (
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	. "github.com/smartystreets/goconvey/convey"
)

func TestTimeout(t *testing.T) {
	if testing.Short() {
		Convey("超时", t, func() {
			r := NewZeroGinRouter()
			r.Use(TimeoutMiddleware(1000))
			r.Handle(http.MethodGet, "timeout", func(_ *gin.Context) Render {
				time.Sleep(2 * time.Second)
				return Success(struct{}{})
			})

			res := testHttpRequest("GET", "/timeout", nil, r)
			So(res.Code, ShouldEqual, 504)
		})
	}
}

func TestNormal(t *testing.T) {
	if testing.Short() {
		Convey("未超时", t, func() {
			r := NewZeroGinRouter()
			r.Use(TimeoutMiddleware(1000))
			r.Handle(http.MethodGet, "timeout", func(_ *gin.Context) Render {
				return Success(struct{}{})
			})

			res := testHttpRequest("GET", "/timeout", nil, r)
			So(res.Code, ShouldEqual, 200)
		})
	}
}

package ginx

import (
	"bytes"
	"encoding/json"
	"github.com/gin-contrib/timeout"
	"github.com/gin-gonic/gin"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/zeromicro/go-zero/rest"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	c = rest.RestConf{
		Host:    "0.0.0.0",
		Port:    9999,
		Timeout: 1000,
	}
)

func TestNewGinRouter(t *testing.T) {
	Convey("new router", t, func() {
		r := NewZeroGinRouter()
		r.Handle("GET", "/test-router", func(ctx *gin.Context) Render {
			return Success(struct{}{})
		})

		res := testHttpRequest("GET", "/", nil, r)
		So(res.Code, ShouldEqual, 404)

		res = testHttpRequest("GET", "/test-router", nil, r)
		So(res.Code, ShouldEqual, 200)

		res = testHttpRequest("POST", "/test-router", nil, r)
		So(res.Code, ShouldEqual, 404)
	})
}

func TestNewGinRouter2(t *testing.T) {
	Convey("集成go-zero", t, func() {
		r := NewZeroGinRouter()

		r.Handle("GET", "/test-router", func(ctx *gin.Context) Render {
			return Success(struct{}{})
		})

		res := testHttpRequest("GET", "/test-router", nil, r)
		So(res.Code, ShouldEqual, 200)
	})
}

func TestAddRouter(t *testing.T) {
	Convey("增加路由", t, func() {
		r := NewZeroGinRouter()
		ag := r.Group("/admin/v1")
		ug := ag.Group("user")
		ug.Handle(http.MethodGet, "add", func(ctx *gin.Context) Render {
			return Success(struct{}{})
		})
		ug.Handle(http.MethodGet, "delete", func(ctx *gin.Context) Render {
			return Success(struct{}{})
		})
		ug.Handle(http.MethodGet, "update", func(ctx *gin.Context) Render {
			return Success(struct{}{})
		})
		ug.Handle(http.MethodGet, "list", func(ctx *gin.Context) Render {
			return Success(struct{}{})
		})
		ug.Handle(http.MethodGet, "detail", func(ctx *gin.Context) Render {
			return Success(struct{}{})
		})

		r.Use(gin.Recovery())
		r.Use(timeout.New())
	})
}

func testHttpRequest(method, path string, body interface{}, r ZeroGinRouter) *httptest.ResponseRecorder {
	bs, _ := json.Marshal(body)
	req, _ := http.NewRequest(method, path, bytes.NewReader(bs))
	req.Header.Set("Request_Timeout", "1")
	writer := httptest.NewRecorder()
	r.(*ginRouter).engine.ServeHTTP(writer, req)
	return writer
}

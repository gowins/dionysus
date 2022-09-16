package main

import (
	"fmt"
	"github.com/gowins/dionysus/ginx"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
)

func main() {
	r := ginx.NewZeroGinRouter()
	r.Use(gin.Logger())
	r.Handle(http.MethodGet, "test", func(c *gin.Context) render.Render {
		return ginx.Success(time.Now().Unix())
	})
	ag := r.Group("admin/v1")
	ag.Handle(http.MethodPost, "get", func(c *gin.Context) render.Render {
		var tt = struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		}{}
		if err := c.ShouldBind(&tt); err != nil {
			return ginx.Error(err)
		}
		fmt.Printf("%v\n", tt)
		return ginx.Success(time.Now().Unix())
	})

	ag.Handle(http.MethodPost, "panic", func(c *gin.Context) render.Render {
		panic("hehe")
	})

	ag.Handle(http.MethodPost, "error", func(c *gin.Context) render.Render {
		return ginx.Error(ginx.NewGinError(350001, "请重新登陆"))
	})

	defer r.Shutdown()
	r.Run(":9999")
}

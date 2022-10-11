package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gowins/dionysus"
	"github.com/gowins/dionysus/cmd"
	"github.com/gowins/dionysus/ginx"
)

func main() {
	ginCmd := cmd.NewGinCommand()
	addRoute(ginCmd)
	dionysus.Start("msg-srv", ginCmd)
}

func addRoute(engine ginx.ZeroGinRouter) {
	engine.Use(gin.Logger())
	engine.Use(func(_ *gin.Context) {
		fmt.Println("root RouteGroup")
	})
	adminGroup := engine.Group("admin/v1", func(_ *gin.Context) {
		fmt.Println("this is admin group")
	})
	adminGroup.Handle(http.MethodGet, "user/get", userGet)
	adminGroup.Handle(http.MethodPost, "user/post", userPost)
	adminGroup.HandleE(http.MethodGet, "user/list", func(c *gin.Context) error {
		return ginx.NewGinError(-1, "testing error")
	})

	webGroup := engine.Group("web/v1", func(_ *gin.Context) {
		fmt.Println("this is web group")
	})
	webGroup.Handle(http.MethodGet, "user/get", userGet)
	webGroup.Handle(http.MethodPost, "user/post", userPost)
}

func userGet(c *gin.Context) ginx.Render {
	return ginx.Success("获取用户成功")
}

func userPost(c *gin.Context) ginx.Render {
	var user = struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
		Age  int    `json:"age"`
	}{}

	if err := c.ShouldBind(&user); err != nil {
		return ginx.Error(err)
	}

	return ginx.Success("修改用户成功")
}

func customError(c *gin.Context) ginx.Render {
	return ginx.Error(ginx.NewGinError(350001, "自定义错误"))
}

func customPanic(c *gin.Context) ginx.Render {
	panic("自定义panic")
}

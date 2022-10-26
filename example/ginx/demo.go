package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gowins/dionysus"
	"github.com/gowins/dionysus/cmd"
	"github.com/gowins/dionysus/ginx"
	"net/http"
	"time"
)

func main() {
	gcmd := cmd.NewGinCommand()
	gcmd.Handle(http.MethodGet, "/test", func(c *gin.Context) ginx.Render {
		return ginx.Success(time.Now().Unix())
	})

	//注意这时返回默认http code状态码200, 这里的code不是事http code
	gcmd.Handle(http.MethodGet, "/test/error", func(c *gin.Context) ginx.Render {
		return ginx.Error(ginx.NewGinError(500100, "内部错误"))
	})

	demogroup := gcmd.Group("/demogroup")
	demogroup.Use(demoMiddle, demoMiddle2)
	demogroup.Handle(http.MethodGet, "/demoroute", func(c *gin.Context) ginx.Render {
		//会拿到demoMiddle中间件塞入gin.Contex的值
		va, ok := c.Get("demoMiddle")
		if !ok {
			return ginx.Error(ginx.NewGinError(500100, "内部错误"))
		}
		return ginx.Success(va)
	})

	//会返回http code 504错误码
	demogroup.Handle(http.MethodGet, "/error", func(c *gin.Context) ginx.Render {
		// 设置http返回状态码
		c.Status(http.StatusGatewayTimeout)
		//c.JSON(http.StatusInternalServerError, "InternalServerError")
		return ginx.Error(ginx.NewGinError(500100, "内部错误"))
	})
	d := dionysus.NewDio()
	if err := d.DioStart("ginxdemo", gcmd); err != nil {
		fmt.Printf("dio start error %v\n", err)
	}
}

func demoMiddle(c *gin.Context) {
	demoValue := time.Now().String() + " demoMiddle"
	fmt.Printf("this demoMiddle1\n")
	c.Set("demoMiddle", demoValue)
	c.Next()
}

func demoMiddle2(c *gin.Context) {
	demoValue := time.Now().String() + " demoMiddle2"
	fmt.Printf("this demoMiddle2\n")
	c.Set("demoMiddle2", demoValue)
	c.Next()
}

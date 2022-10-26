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
	demogroup := gcmd.Group("/demogroup")
	demogroup.Use(demoMiddle)
	demogroup.Handle(http.MethodGet, "/demoroute", func(c *gin.Context) ginx.Render {
		//会拿到demoMiddle中间件塞入gin.Contex的值
		va, ok := c.Get("demoMiddle")
		if !ok {
			return ginx.Error(ginx.NewGinError(500100, "内部错误"))
		}
		return ginx.Success(va)
	})

	//注意这时的http code返回是200, 这里的code不是事http code
	gcmd.Handle(http.MethodGet, "/test/error", func(c *gin.Context) ginx.Render {
		return ginx.Error(ginx.NewGinError(500100, "内部错误"))
	})

	//会返回http code 500错误码
	demogroup.Handle(http.MethodGet, "/error", func(c *gin.Context) ginx.Render {
		//c.Status(http.StatusInternalServerError)
		c.JSON(http.StatusInternalServerError, "InternalServerError")
		return nil
	})
	d := dionysus.NewDio()
	if err := d.DioStart("ginxdemo", gcmd); err != nil {
		fmt.Printf("dio start error %v\n", err)
	}
}

func demoMiddle(c *gin.Context) {
	demoValue := time.Now().String() + " demoMiddle"
	c.Set("demoMiddle", demoValue)
	c.Next()
}

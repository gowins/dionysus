package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gowins/dionysus"
	"github.com/gowins/dionysus/cmd"
	"github.com/gowins/dionysus/ginx"
	otm "github.com/gowins/dionysus/opentelemetry"
	"github.com/gowins/dionysus/step"
	oteltrace "go.opentelemetry.io/otel/trace"
	"log"
	"net/http"
	"time"
)

func main() {
	// 创建dionysus框架
	d := dionysus.NewDio()
	//创建gin的子cmd
	gcmd := cmd.NewGinCommand()
	// code必须大于100000
	if err := ginx.SetDefaultErrorCode(110000); err != nil {
		log.Fatalf("set default error code failed %v", err)
	}
	gcmd.RegShutdownFunc(cmd.StopStep{
		StopFn: func() {
			otm.Stop()
			fmt.Printf("this gcmd stop\n")
		},
		StepName: "gcmdstop",
	})
	d.PreRunStepsAppend(step.InstanceStep{
		StepName: "init trace",
		Func: func() error {
			otm.Setup(otm.WithServiceInfo(&otm.ServiceInfo{
				Name:      "testGin217",
				Namespace: "testGinNamespace217",
				Version:   "testGinVersion217",
			}), otm.WithTraceExporter(&otm.Exporter{
				ExporterEndpoint: otm.DefaultStdout,
				Insecure:         false,
				Creds:            otm.DefaultCred,
			}))
			return nil
		},
	})
	//定义路由/test和相应的handler函数
	gcmd.Handle(http.MethodGet, "/test", func(c *gin.Context) ginx.Render {
		return ginx.Success(time.Now().Unix())
	})

	//定义路由/test/error，注意这时返回默认http code状态码200, 这里的code不是事http code
	gcmd.Handle(http.MethodGet, "/test/error", func(c *gin.Context) ginx.Render {
		return ginx.Error(ginx.NewGinError(500100, "内部错误"))
	})

	//创建路由组/demogroup
	demogroup := gcmd.Group("/demogroup")
	//在路由组/demogroup应用中间件authMiddle，认证失败就返回错误
	//demogroup.Use(authMiddle)
	//在路由组/demogroup应用中间件demoMiddle，塞入请求时间
	demogroup.Use(demoMiddle)
	//在路由组/demogroup下组成路由/demogroup/demoroute将会执行上面注册的中间件authMiddle和demoMiddle
	demogroup.Handle(http.MethodGet, "/demoroute", func(c *gin.Context) ginx.Render {
		reqCtx := c.Request.Context()
		//fmt.Printf("Request header is %v\n", c.Request.Header)
		//ctx, span := tracer.Start(reqCtx, "aaa", opts...)
		span := oteltrace.SpanFromContext(reqCtx)
		span.AddEvent("qwqw")
		fmt.Printf("traced %v\n", span.SpanContext().TraceID().String())
		//c.Request = c.Request.WithContext()
		//defer span.End()
		//_, mspan := tracer.Start(ctx, "aaabbb")
		//mspan.AddEvent("aabbb")
		//mspan.End()
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
	//将子gin cmd注册到dionysus框架中，并启动程序
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

func authMiddle(c *gin.Context) {
	if authentication(c) {
		c.Next()
	} else {
		//如不想在执行其它中间件及handler函数，请调用c.AbortXXX
		c.AbortWithStatus(http.StatusForbidden)
	}
}

func authentication(c *gin.Context) bool {
	if c.GetHeader("identity") == "owner" {
		return true
	} else {
		return false
	}
}

/* 测试命令
curl -i 127.0.0.1:8080/demogroup/demoroute -H "identity: owner"
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Content-Length: 103

curl -i 127.0.0.1:8080/demogroup/error
HTTP/1.1 500 Internal Server Error
Content-Length: 0

curl -i 127.0.0.1:8080/test/error
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Content-Length: 46
{"code":500100,"msg":"内部错误","data":{}}
*/

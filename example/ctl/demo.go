package main

import (
	"context"
	"fmt"
	"github.com/gowins/dionysus/httpclient"
	otm "github.com/gowins/dionysus/opentelemetry"
	"github.com/gowins/dionysus/step"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	oteltrace "go.opentelemetry.io/otel/trace"
	"io"
	"net/http"
	"time"

	"github.com/gowins/dionysus"
	"github.com/gowins/dionysus/cmd"
)

var (
	green  = string([]byte{27, 91, 51, 50, 109})
	white  = string([]byte{27, 91, 51, 55, 109})
	yellow = string([]byte{27, 91, 51, 51, 109})
	red    = string([]byte{27, 91, 51, 49, 109})
	blue   = string([]byte{27, 91, 51, 52, 109})
	// magenta = string([]byte{27, 91, 51, 53, 109})
	// cyan    = string([]byte{27, 91, 51, 54, 109})
	// color   = []string{green, white, yellow, red, blue, magenta, cyan}
)

func main() {
	d := dionysus.NewDio()
	postSteps := []step.InstanceStep{
		{
			StepName: "PostPrint1", Func: func() error {
				fmt.Println(green, "=========== post 1 =========== ", white)
				otm.Stop()
				return nil
			},
		},
		{
			StepName: "PostPrint2", Func: func() error {
				fmt.Println(green, "=========== post 2 =========== ", white)
				return nil
			},
		},
	}
	preSteps := []step.InstanceStep{
		{
			StepName: "PrePrint1", Func: func() error {
				otm.Setup(otm.WithServiceInfo(&otm.ServiceInfo{
					Name:      "testCtl217",
					Namespace: "testCtlNamespace217",
					Version:   "testCtlVersion217",
				}))
				return nil
			},
		},
		{
			StepName: "PrePrint2", Func: func() error {
				fmt.Println(green, "=========== pre 2 =========== ", white)
				return nil
			},
		},
	}
	// PreRun exec before server start
	_ = d.PreRunStepsAppend(preSteps...)
	httpClient := httpclient.NewWithTracer()
	ctlCmd := cmd.NewCtlCommand()
	_ = ctlCmd.RegRunFunc(func() error {
		tracer := otel.GetTracerProvider().Tracer("hgsmytracer", oteltrace.WithInstrumentationVersion("v0.1.1"))
		timer1 := time.NewTicker(10 * time.Second)
		for {
			select {
			case <-timer1.C:
				opts := []oteltrace.SpanStartOption{
					oteltrace.WithAttributes(attribute.String("go.template1", "name")),
					oteltrace.WithSpanKind(oteltrace.SpanKindServer),
				}
				ctx, span := tracer.Start(context.Background(), "bbb", opts...)
				fmt.Printf("traced %v\n", span.SpanContext().TraceID().String())
				span.AddEvent("do http req")
				url := "http://127.0.0.1:8080/demogroup/demoroute"
				request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
				if err != nil {
					fmt.Printf("new request error %v\n", err)
					span.End()
					continue
				}
				rsp, err := httpClient.Do(request)
				if err != nil {
					fmt.Printf("time %v http error %v\n", time.Now().String(), err)
					span.End()
					continue
				}
				if rsp == nil {
					fmt.Printf("http error rsp is nil\n")
					span.End()
					continue
				}
				if rsp.StatusCode >= 300 {
					fmt.Printf("http error code %v\n", rsp.StatusCode)
				} else {
					fmt.Printf("http sucess code %v\n", rsp.StatusCode)
				}
				io.Copy(io.Discard, rsp.Body)
				rsp.Body.Close()
				span.End()
			case <-ctlCmd.Ctx.Done():
				fmt.Printf("this is stopChan %v\n", time.Now().String())
				return nil
			}
		}
	})

	ctx, cancel := context.WithCancel(ctlCmd.Ctx)
	ctlCmd.Ctx = ctx
	stopSteps := []cmd.StopStep{
		{
			StepName: "before stop",
			StopFn: func() {
				fmt.Printf("this is before stop\n")
			},
		},
		{
			StepName: "stop",
			StopFn: func() {
				fmt.Printf("this is stop\n")
				cancel()
			},
		},
	}
	ctlCmd.RegShutdownFunc(stopSteps...)

	// PostRun exec after server stop
	_ = d.PostRunStepsAppend(postSteps...)

	if err := d.DioStart("ctldemo", ctlCmd); err != nil {
		fmt.Printf("dio start error %v\n", err)
	}
}

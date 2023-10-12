package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gowins/dionysus"
	"github.com/gowins/dionysus/cmd"
	"github.com/gowins/dionysus/example/grpc/hw"
	"github.com/gowins/dionysus/ginx"
	"github.com/gowins/dionysus/grpc/client"
	"github.com/gowins/dionysus/httpclient"
	"github.com/gowins/dionysus/kafka"
	"github.com/gowins/dionysus/log"
	otm "github.com/gowins/dionysus/opentelemetry"
	"github.com/gowins/dionysus/orm"
	"github.com/gowins/dionysus/step"
	kafkago "github.com/segmentio/kafka-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	oteltrace "go.opentelemetry.io/otel/trace"
	"gorm.io/plugin/opentelemetry/tracing"
)

var httpClient httpclient.Client

// 在dio启动步骤中，增加opentelemetry的初始化及stop动作
func main() {
	dio := dionysus.NewDio()
	preSteps := []step.InstanceStep{
		{
			StepName: "initLog", Func: func() error {
				log.Setup()
				return nil
			},
		},
		{
			StepName: "initOtm", Func: func() error {
				otm.Setup(otm.WithServiceInfo(&otm.ServiceInfo{
					Name:      "dioName",
					Namespace: "dioNamespace",
					Version:   "dioVersion",
				}), otm.WithTraceExporter(&otm.Exporter{
					ExporterEndpoint: otm.DefaultStdout, //后续具体对接相应实例
					Insecure:         false,
					Creds:            otm.DefaultCred,
				}), otm.WithSampleRatio(1.0)) // 设置采样率，ParentBased采样策略，默认值为100%
				return nil
			},
		},
	}
	_ = dio.PreRunStepsAppend(preSteps...)
	postSteps := []step.InstanceStep{
		{
			StepName: "stopOtm", Func: func() error {
				//处理上报队列缓存中的trace，stop之前到的trace将会上报，而stop之后新产生的trace丢弃。
				otm.Stop()
				return nil
			},
		},
	}
	_ = dio.PostRunStepsAppend(postSteps...)
	//初始化gin, grpc, ctl等cmd
	gcmd := cmd.NewGinCommand()
	gcmd.Handle(http.MethodGet, "/otm/api/test", otmApiTest)
	// httpClient必须使用WithTracer版本
	httpClient = httpclient.NewWithTracer()
	if err := dio.DioStart("ginxdemo", gcmd); err != nil {
		log.Errorf("dio start error %v\n", err)
	}
}

//获取span有两种方法，第一种是创建新的span，第二种是context中获取到span

// 创建一个新的span并发起http请求
// spanHttpReq----------------------------------------------------------------spanHttpReq
//
//	spanHttpCli------------------------------------spanHttpCli
//	              spanHttpServer---spanHttpServer
func creatSpanDoHttpReq() {
	opts := []oteltrace.SpanStartOption{
		oteltrace.WithAttributes(attribute.String("httpCli", "name")),
		oteltrace.WithSpanKind(oteltrace.SpanKindClient),
	}
	ctx, span := otm.SpanStart(context.Background(), "httpReq", opts...)
	if span == nil {
		log.Errorf("trace is not init")
		return
	}
	// 如果创建了新的span，一定要end，不然会内存泄露，end调用后会将这次trace span放到上报队列中
	defer span.End()
	span.AddEvent("do http req")
	url := "http://127.0.0.1:8080/otm/api/test"
	//必须将上面生成的context放如req中, 才能将相关的trace传输到server端
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		fmt.Printf("new request error %v\n", err)
		return
	}
	// httpClient必须使用WithTracer版本
	rsp, err := httpClient.Do(request)
	if err != nil {
		fmt.Printf("time %v http error %v\n", time.Now().String(), err)
		return
	}
	if rsp == nil {
		fmt.Printf("http error rsp is nil\n")
		return
	}
	if rsp.StatusCode >= 300 {
		fmt.Printf("http error code %v\n", rsp.StatusCode)
	} else {
		log.Infof("http sucess code %v, time %v\n", rsp.StatusCode, time.Now().String())
	}
	io.Copy(io.Discard, rsp.Body)
	rsp.Body.Close()
}

// 第二种是从context中获取span，这种情况下不用end。
// spanA --- add event --- spanA
func otmApiTest(c *gin.Context) ginx.Render {
	//  这种没有新创建span，不要调用end
	span := oteltrace.SpanFromContext(c.Request.Context())
	span.AddEvent("this span father test")
	return ginx.Success("otm test")
}

// 对于grpc到client来说，创建连接时，必须使用NewConnWithTracer
func grpcCLient() {
	// 必须使用NewConnWithTracer
	grpcConn, _ := client.NewConnWithTracer("grpcServiceName")
	c := hw.NewGreeterClient(grpcConn)
	opts := []oteltrace.SpanStartOption{
		oteltrace.WithAttributes(attribute.String("GrpcCli", "name")),
		oteltrace.WithSpanKind(oteltrace.SpanKindClient),
	}
	ctx, span := otm.SpanStart(context.Background(), "grpcReq", opts...)
	// 如果创建了新的span，一定要end，不然会内存泄露，end调用后会将这次trace span放到上报队列中
	defer span.End()
	r, err := c.SayHello(ctx, &hw.HelloRequest{Name: "nameing"})
	if err != nil {
		log.Errorf("could not greet: %v", err)
		return
	}
	log.Infof("Greeting: %s", r.GetMessage())
}

func kafkaWithTrace(ctx context.Context) {
	kafkaWriter, _ := kafka.NewWriter([]string{"testAddr"}, "topic")
	// 如果ctx有span，那么spanDB会有old span为parent，反之则spanDB为root span
	ctxKafka, spanDB := otm.SpanStart(ctx, "kafkaWriter", oteltrace.WithSpanKind(oteltrace.SpanKindProducer),
		oteltrace.WithAttributes(attribute.String("messaging.system", "kafka")))
	if spanDB == nil {
		log.Errorf("trace is not init")
		return
	}
	// 新创建到span必须关闭
	defer spanDB.End()
	err := kafkaWriter.WriteMessages(ctxKafka, kafkago.Message{
		Key:   []byte("Key1"),
		Value: []byte("Value1"),
	})
	if err != nil {
		log.Errorf("kafka write message error %v", err)
		spanDB.SetStatus(codes.Error, fmt.Sprintf("kafka write message error %v", err))
		return
	}
	spanDB.AddEvent("kafkaWriter ok", oteltrace.WithAttributes(
		attribute.KeyValue{Key: "msg write", Value: attribute.StringValue("Key1")}))
}

func ormWithTrace(ctx context.Context) {
	ormDB := orm.GetDB("demoDB")
	if otm.TracerIsEnable() {
		err := ormDB.Use(tracing.NewPlugin(tracing.WithoutMetrics(), tracing.WithTracerProvider(otel.GetTracerProvider())))
		if err != nil {
			log.Fatalf("gorm db init tracer failed %s", err.Error())
		}
	}
}

func commonCtlTrace(ctx context.Context) {
	// 如果ctx有span，那么spanDB会有old span为parent，反之则spanDB为root span
	_, spanDB := otm.SpanStart(ctx, "commonCtl", oteltrace.WithSpanKind(oteltrace.SpanKindInternal))
	if spanDB == nil {
		log.Errorf("trace is not init")
		return
	}
	// 新创建到span必须关闭
	defer spanDB.End()
	log.Infof("commctl traceid id %v", otm.GetTraceId(ctx))
}

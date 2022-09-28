package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gowins/dionysus/rmq"

	"github.com/apache/rocketmq-client-go/v2/primitive"
)

func main() {
	client, err := rmq.NewProducer(&rmq.ProducerConfig{
		NameSrvAddr:    []string{"http://rmq-cn-zvp2ud6lc0e.cn-hangzhou.rmq.aliyuncs.com:8080"},
		UseCredentials: true,
		AccessKey:      "5qjJ0K0SdXIvuj2t",
		SecretKey:      "3yuDlpGFwL2itNkAx2Djxe+XcK1BDLC1VFqD6b7CV+A=",
		Group:          "hltv_g",
	})
	if err != nil {
		panic(err)
	}

	if err := client.Start(); err != nil {
		panic(err)
	}

	for i := 0; i < 1000; i++ {
		msg := &primitive.Message{
			Topic: "hltv",
			Body:  []byte(fmt.Sprintf("%d", time.Now().UnixNano())),
		}
		msg.WithTag("test")
		if mr, err := client.SendSync(context.Background(), msg); err != nil {
			log.Println(err.Error())
		} else {
			log.Printf("%v\n", mr)
		}
	}
}

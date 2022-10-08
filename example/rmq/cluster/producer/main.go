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
		SecretKey:      "lYlURcnL7z3q3GbO",
		Group:          "msg_g",
	})
	if err != nil {
		panic(err)
	}

	if err := client.Start(); err != nil {
		panic(err)
	}
	for i := 0; i < 1000; i++ {
		msg := &primitive.Message{
			Topic: "msg",
			Body:  []byte(fmt.Sprintf("%d", time.Now().UnixNano())),
		}
		msg.WithTag("tag_a")
		if mr, err := client.SendSync(context.Background(), msg); err != nil {
			log.Println(err.Error())
		} else {
			log.Printf("tag=tag_a count=%d result=%v\n", i, mr)
		}
	}

	for i := 0; i < 1000; i++ {
		msg := &primitive.Message{
			Topic: "msg",
			Body:  []byte(fmt.Sprintf("%d", time.Now().UnixNano())),
		}
		msg.WithTag("tag_b")
		if mr, err := client.SendSync(context.Background(), msg); err != nil {
			log.Println(err.Error())
		} else {
			log.Printf("tag=tag_b count=%d result=%v\n", i, mr)
		}
	}
}

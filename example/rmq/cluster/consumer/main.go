package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/gowins/dionysus/rmq"

	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
)

var (
	tag   = flag.String("tag", "tag_a", "tag")
	count = 0
)

func main() {
	flag.Parse()
	client, err := rmq.NewConsumer(&rmq.ConsumerConfig{
		NameSrvAddr:    []string{"http://rmq-cn-zvp2ud6lc0e.cn-hangzhou.rmq.aliyuncs.com:8080"},
		GroupName:      "msg_g",
		ConsumerModel:  1,
		UseCredentials: true,
		AccessKey:      "5qjJ0K0SdXIvuj2t",
		SecretKey:      "lYlURcnL7z3q3GbO",
		ConsumerOrder:  false,
	})
	if err != nil {
		panic(err)
	}
	if err := client.Subscribe("msg", consumer.MessageSelector{
		Type:       consumer.TAG,
		Expression: *tag,
	}, func(ctx context.Context, ext ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
		for _, m := range ext {
			count++
			fmt.Printf("receive msg topic %s tag %s broker %s queue %d offset %s body %s count %d\n",
				m.Topic, m.GetTags(), m.Queue.BrokerName, m.Queue.QueueId, m.OffsetMsgId, string(m.Body), count)
		}
		return consumer.ConsumeSuccess, nil
	}); err != nil {
		panic(err)
	}

	if err := client.Start(); err != nil {
		panic(err)
	}

	select {}
}

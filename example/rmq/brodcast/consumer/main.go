package main

import (
	"context"
	"fmt"

	"github.com/gowins/dionysus/rmq"

	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
)

func main() {
	client, err := rmq.NewConsumer(&rmq.ConsumerConfig{
		NameSrvAddr:    []string{"http://rmq-cn-zvp2ud6lc0e.cn-hangzhou.rmq.aliyuncs.com:8080"},
		GroupName:      "hltv_g",
		ConsumerModel:  0,
		UseCredentials: true,
		AccessKey:      "5qjJ0K0SdXIvuj2t",
		SecretKey:      "3yuDlpGFwL2itNkAx2Djxe+XcK1BDLC1VFqD6b7CV+A=",
		ConsumerOrder:  false,
	})
	if err != nil {
		panic(err)
	}
	if err := client.Subscribe("hltv", consumer.MessageSelector{}, func(ctx context.Context, ext ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
		for _, m := range ext {
			fmt.Printf("receive msg topic %s tag %s broker %s queue %d body %s\n", m.Topic, m.GetTags(), m.Queue.BrokerName, m.Queue.QueueId, string(m.Body))
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

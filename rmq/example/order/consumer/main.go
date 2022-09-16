package main

import (
	"context"
	"fmt"
	"github.com/gowins/dionysus/rmq"
	"sync/atomic"

	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
)

var (
	idx int64
)

func main() {
	client, err := rmq.NewConsumer(&rmq.ConsumerConfig{
		NameSrvAddr:    []string{"http://rmq-cn-zvp2ud6lc0e.cn-hangzhou.rmq.aliyuncs.com:8080"},
		GroupName:      "hltv_g",
		ConsumerModel:  1,
		UseCredentials: true,
		AccessKey:      "5qjJ0K0SdXIvuj2t",
		SecretKey:      "3yuDlpGFwL2itNkAx2Djxe+XcK1BDLC1VFqD6b7CV+A=",
		ConsumerOrder:  true,
	})
	if err != nil {
		panic(err)
	}
	if err := client.Subscribe("hltv", consumer.MessageSelector{}, func(ctx context.Context, ext ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
		for _, m := range ext {
			res := atomic.AddInt64(&idx, 1)
			fmt.Printf("receive msg idx %d topic %s tag %s shardingKey %s broker %s queue %d offset %d body %s\n",
				res, m.Topic, m.GetTags(), m.GetShardingKey(), m.Queue.BrokerName, m.Queue.QueueId, m.QueueOffset, string(m.Body))
		}
		return consumer.ConsumeSuccess, nil
	}); err != nil {
		panic(err)
	}

	if err := client.Start(); err != nil {
		panic(err)
	}
	defer func() {
		_ = client.Shutdown()
	}()
	select {}
}

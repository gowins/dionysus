package rmq

import (
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/rs/xid"
)

type ConsumerConfig struct {
	NameSrvAddr    []string `json:"name_srv_addr" yaml:"name_srv_url"`                // name srv地址
	GroupName      string   `json:"group_name,omitempty" yaml:"group_name"`           // 组名
	ConsumerModel  int      `json:"consumer_model,omitempty" yaml:"consumer_model"`   // 0-广播模式 1-集群模式
	UseCredentials bool     `json:"use-credentials,omitempty" yaml:"use_credentials"` // 是否需要身份验证
	AccessKey      string   `json:"access_key,omitempty" yaml:"access_key"`           // AccessKey 阿里云身份验证，在阿里云服务器管理控制台创建
	SecretKey      string   `json:"secret_key,omitempty" yaml:"secret_key"`           // SecretKey 阿里云身份验证，在阿里云服务器管理控制台创建
	ConsumerOrder  bool     `json:"consumer_order,omitempty" yaml:"consumer_order"`   // 是否顺序消费
}

func NewConsumer(c *ConsumerConfig) (rocketmq.PushConsumer, error) {
	var opts []consumer.Option
	opts = append(opts, consumer.WithNameServer(c.NameSrvAddr))
	if c.GroupName != "" {
		opts = append(opts, consumer.WithGroupName(c.GroupName))
	}
	if c.UseCredentials {
		opts = append(opts, consumer.WithCredentials(primitive.Credentials{
			AccessKey: c.AccessKey,
			SecretKey: c.SecretKey,
		}))
	}
	opts = append(opts, consumer.WithConsumerModel(consumer.MessageModel(c.ConsumerModel)))
	opts = append(opts, consumer.WithConsumerOrder(c.ConsumerOrder))
	opts = append(opts, consumer.WithInstance(xid.New().String()))

	client, err := consumer.NewPushConsumer(opts...)
	if err != nil {
		return nil, err
	}

	return client, nil
}

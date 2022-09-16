package rmq

import (
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
)

type ProducerConfig struct {
	NameSrvAddr    []string `json:"name_srv_addr" yaml:"name_srv_addr"`
	UseCredentials bool     `json:"use_credentials,omitempty" yaml:"use_credentials"`
	AccessKey      string   `json:"access_key,omitempty" yaml:"access_key"`
	SecretKey      string   `json:"secret_key,omitempty" yaml:"secret_key"`
	Group          string   `json:"group,omitempty" yaml:"group"`
}

func NewProducer(c *ProducerConfig) (rocketmq.Producer, error) {
	var opts []producer.Option
	opts = append(opts, producer.WithNameServer(c.NameSrvAddr))
	if c.UseCredentials {
		opts = append(opts, producer.WithCredentials(primitive.Credentials{
			AccessKey: c.AccessKey,
			SecretKey: c.SecretKey,
		}))
	}
	if c.Group != "" {
		opts = append(opts, producer.WithGroupName(c.Group))
	}

	client, err := rocketmq.NewProducer(opts...)
	if err != nil {
		return nil, err
	}
	if err = client.Start(); err != nil {
		return nil, err
	}
	return client, nil
}

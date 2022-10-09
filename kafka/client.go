package kafka

import (
	"fmt"

	kafka "github.com/segmentio/kafka-go"
)

func NewPartitionReader(brokers []string, topic string, partition int, options ...ReaderOption) (r *kafka.Reader, e error) {
	return newReader(brokers, topic, "", partition, options...)
}

func NewGroupReader(brokers []string, topic, groupID string, options ...ReaderOption) (r *kafka.Reader, e error) {
	if groupID == "" {
		return nil, fmt.Errorf("groupID can not be empty string ")
	}
	return newReader(brokers, topic, groupID, 0, options...)
}

func newReader(brokers []string, topic, groupID string, partition int, options ...ReaderOption) (r *kafka.Reader, e error) {
	defer func() {
		if rec := recover(); rec != nil {
			r, e = nil, fmt.Errorf("New Reader Error:%s ", rec)
			return
		}
	}()

	conf := defaultReaderConfig
	conf.Brokers = brokers
	conf.Topic = topic
	conf.Partition = partition
	conf.GroupID = groupID
	conf.Logger = kLogger
	conf.ErrorLogger = kErrorLogger

	if conf.GroupID != "" {
		conf.WatchPartitionChanges = true
	}

	for _, o := range options {
		o(&conf)
	}

	r, e = kafka.NewReader(conf), nil
	return
}

func NewPartitionWriter(brokers []string, topic string, partition int, options ...WriterOption) (w *kafka.Writer, e error) {
	p := &PartitionBalancer{partition}
	return newWriter(brokers, topic, p, options...)
}

func NewWriter(brokers []string, topic string, options ...WriterOption) (w *kafka.Writer, e error) {
	return newWriter(brokers, topic, nil, options...)
}

func newWriter(brokers []string, topic string, balancer kafka.Balancer, options ...WriterOption) (w *kafka.Writer, e error) {
	defer func() {
		if rec := recover(); rec != nil {
			w, e = nil, fmt.Errorf("New Writer Error:%s ", rec)
			return
		}
	}()

	conf := defaultWriterConfig
	conf.Brokers = brokers
	conf.Topic = topic
	conf.Logger = kLogger
	conf.ErrorLogger = kErrorLogger

	for _, o := range options {
		o(&conf)
	}

	w, e = kafka.NewWriter(conf), nil
	// there is a bug: https://github.com/segmentio/kafka-go/issues/524
	w.Balancer = balancer
	return
}

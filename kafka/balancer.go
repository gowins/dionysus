package kafka

import "github.com/segmentio/kafka-go"

// partitionBalancer
type PartitionBalancer struct {
	p int
}

func (p *PartitionBalancer) Balance(_ kafka.Message, _ ...int) (partition int) {
	return p.p
}

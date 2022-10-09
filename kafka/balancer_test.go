package kafka

import (
	"testing"

	"github.com/segmentio/kafka-go"
)

func TestPartitionBalancer(t *testing.T) {
	specificPartition := 2

	testCases := map[string]struct {
		Keys       [][]byte
		Partitions [][]int
		Partition  int
	}{
		"single message": {
			Keys: [][]byte{
				[]byte("key"),
			},
			Partitions: [][]int{
				{0, 1, 2},
			},
			Partition: specificPartition,
		},
		"multiple messages": {
			Keys: [][]byte{
				[]byte("a"),
				[]byte("ab"),
				[]byte("abc"),
				[]byte("abcd"),
			},
			Partitions: [][]int{
				{0, 1, 2},
				{0, 1, 2},
				{0, 1, 2},
				{0, 1, 2},
			},
			Partition: specificPartition,
		},
		"partition lost": {
			Keys: [][]byte{
				[]byte("hello world 1"),
				[]byte("hello world 2"),
				[]byte("hello world 3"),
			},
			Partitions: [][]int{
				{0, 1},
				{0, 1},
				{0, 1, 2},
			},
			Partition: specificPartition,
		},
	}

	for label, test := range testCases {
		t.Run(label, func(t *testing.T) {
			pb := &PartitionBalancer{specificPartition}

			var partition int
			for i, key := range test.Keys {
				msg := kafka.Message{Key: key}
				partition = pb.Balance(msg, test.Partitions[i]...)
			}

			if partition != test.Partition {
				t.Errorf("expected %v; got %v", test.Partition, partition)
			}
		})
	}
}

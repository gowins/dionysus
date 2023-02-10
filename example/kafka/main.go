package main

import (
	"context"
	"github.com/gowins/dionysus"
	"github.com/gowins/dionysus/cmd"
	"github.com/gowins/dionysus/kafka"
	"github.com/gowins/dionysus/log"
	"github.com/gowins/dionysus/step"
	kafkago "github.com/segmentio/kafka-go"
	"time"
)

func main() {
	d := dionysus.NewDio()
	var err error
	var kafkaReader *kafkago.Reader
	var kafkaWriter *kafkago.Writer
	stopChan := make(chan struct{})

	d.PreRunStepsAppend(step.InstanceStep{
		StepName: "init kafka",
		Func: func() error {
			kafkaReader, err = InitKafkaReader()
			if err != nil {
				return err
			}
			kafkaWriter, err = InitKafkaWriter()
			if err != nil {
				return err
			}
			return nil
		},
	})
	ctlCmd := cmd.NewCtlCommand()
	_ = ctlCmd.RegRunFunc(func() error {
		go KafkaReaderStart(kafkaReader)
		timer1 := time.NewTicker(5 * time.Second)
		for {
			select {
			case <-timer1.C:
				KafkaWriterMsg(kafkaWriter)
			case <-stopChan:
				log.Infof("this is stopChan %v\n", time.Now().String())
				return nil
			}
		}
	})
	ctlCmd.RegShutdownFunc(cmd.StopStep{
		StepName: "stop kafka",
		StopFn: func() {
			stopChan <- struct{}{}
			kafkaReader.Close()
			kafkaWriter.Close()
		},
	})
	d.DioStart("kafkademo", ctlCmd)
}

var testAddr = []string{
	"addr1",
	"addr2",
	"addr3"}

var caCert = "cacertfile"
var userName = "userName"
var password = "password"
var topic = "topic"
var groupId = "groupId"

func InitKafkaReader() (*kafkago.Reader, error) {
	kafkaDialer := kafka.GetKafkaDialer(kafka.DialerWithCaCertTls(caCert), kafka.DialerWithUserAndPasswd(userName, password))
	return kafka.NewGroupReader(testAddr, topic, groupId, kafka.ReaderWithDialer(kafkaDialer))
}

func KafkaReaderStart(kafkaReader *kafkago.Reader) {
	for {
		m, err := kafkaReader.ReadMessage(context.Background())
		if err != nil {
			return
		}
		log.Infof("\n\nmessage at topic/partition/offset %v/%v/%v: %s = %s\n\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))
	}
}

func InitKafkaWriter() (*kafkago.Writer, error) {
	kafkaDialer := kafka.GetKafkaDialer(kafka.DialerWithCaCertTls(caCert), kafka.DialerWithUserAndPasswd(userName, password))
	return kafka.NewWriter(testAddr, topic, kafka.WriterWithDialer(kafkaDialer))
}

func KafkaWriterMsg(kafkaWriter *kafkago.Writer) {
	timeNow := time.Now().String()
	err := kafkaWriter.WriteMessages(context.Background(),
		kafkago.Message{
			Key:   []byte("Key" + timeNow),
			Value: []byte("Value" + timeNow),
		},
	)
	if err != nil {
		log.Errorf("failed to write messages:", err)
	} else {
		log.Infof("write messages %v", timeNow)
	}
}

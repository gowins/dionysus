package kafka

import (
	"testing"
	"time"

	"github.com/segmentio/kafka-go"
	. "github.com/smartystreets/goconvey/convey"
)

func TestWriterOptions(t *testing.T) {
	Convey("Test writer options", t, func() {
		c := defaultWriterConfig

		Convey("Test using async option", func() {
			So(c.Async, ShouldBeFalse)
			WriterWithAsync()(&c)
			So(c.Async, ShouldBeTrue)
		})
		Convey("Test using dial option", func() {
			So(c.Dialer, ShouldBeNil)
			WriterWithDialer(GetKafkaDialer())(&c)
			So(c.Async, ShouldNotBeNil)
		})
		Convey("Test using batchSize option", func() {
			So(c.BatchSize, ShouldEqual, 100)
			WriterWithBatchSize(200)(&c)
			So(c.BatchSize, ShouldEqual, 200)
		})
		Convey("Test using batchBytes option", func() {
			So(c.BatchBytes, ShouldEqual, 2<<19)
			WriterWithBatchBytes(2 << 20)(&c)
			So(c.BatchBytes, ShouldEqual, 2<<20)
		})
		Convey("Test using batchTimeout option", func() {
			So(c.BatchTimeout, ShouldEqual, time.Millisecond*100)
			WriterWithBatchTimeout(time.Millisecond * 200)(&c)
			So(c.BatchTimeout, ShouldEqual, time.Millisecond*200)
		})
		Convey("Test using looger option", func() {
			So(c.Logger, ShouldBeNil)
			So(c.ErrorLogger, ShouldBeNil)
			WriterWithLogger(KLogger, KErrorLogger)(&c)
			So(c.Logger, ShouldNotBeNil)
			So(c.ErrorLogger, ShouldNotBeNil)
		})
	})
}

func TestReaderWithAsync(t *testing.T) {
	Convey("Test reader options", t, func() {
		c := defaultReaderConfig
		testCommitInterval := time.Millisecond * 100
		Convey("Test using async option", func() {
			So(c.CommitInterval, ShouldBeZeroValue)
			ReaderWithCommitInterval(testCommitInterval)(&c)
			So(c.CommitInterval, ShouldEqual, testCommitInterval)
		})
	})
}

func TestReaderWithDialer(t *testing.T) {
	Convey("Test reader options", t, func() {
		c := defaultReaderConfig

		dialer := &kafka.Dialer{}
		Convey("Test using dialer option", func() {
			So(c.Dialer, ShouldBeNil)
			ReaderWithDialer(dialer)(&c)
			So(c.Dialer, ShouldNotBeNil)
		})
	})
}

func TestReaderWithOffset(t *testing.T) {
	Convey("Test reader options", t, func() {
		c := defaultReaderConfig
		Convey("Test using dialer option", func() {
			So(c.StartOffset, ShouldEqual, kafka.FirstOffset)
			ReaderWithOffset(1)(&c)
			So(c.StartOffset, ShouldEqual, 1)
		})
	})
}

func TestGetKafkaDialer(t *testing.T) {
	Convey("Test kafka dialer", t, func() {
		kafkaDialer := GetKafkaDialer(DialerWithCaCertTls("./fake_cacert"), DialerWithUserAndPasswd("testuser", "testpasswd"))
		So(kafkaDialer, ShouldNotBeNil)
		So(kafkaDialer.TLS, ShouldNotBeNil)
		So(kafkaDialer.TLS.InsecureSkipVerify, ShouldBeTrue)
		So(kafkaDialer.TLS.RootCAs, ShouldNotBeNil)
		So(kafkaDialer.SASLMechanism.Name(), ShouldEqual, "PLAIN")
	})
}

func TestReaderWithMinBytes(t *testing.T) {
	Convey("Test reader minBytes options", t, func() {
		c := defaultReaderConfig
		Convey("Test using minBytes option", func() {
			So(c.MinBytes, ShouldEqual, 2<<9)
			ReaderWithMinBytes(3 << 9)(&c)
			So(c.MinBytes, ShouldEqual, 3<<9)
		})
	})
}

func TestReaderWithMaxBytes(t *testing.T) {
	Convey("Test reader maxBytes options", t, func() {
		c := defaultReaderConfig
		Convey("Test using maxBytes option", func() {
			So(c.MaxBytes, ShouldEqual, 2<<23)
			ReaderWithMaxBytes(3 << 23)(&c)
			So(c.MaxBytes, ShouldEqual, 3<<23)
		})
	})
}

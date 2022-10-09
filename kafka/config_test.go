package kafka

import (
	"testing"

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
	})
}

func TestReaderWithAsync(t *testing.T) {
	Convey("Test reader options", t, func() {
		c := defaultReaderConfig

		Convey("Test using async option", func() {
			So(c.CommitInterval, ShouldBeZeroValue)
			ReaderWithAsync()(&c)
			So(c.CommitInterval, ShouldEqual, defaultCommitInterval)
		})
	})
}

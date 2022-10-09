package kafka

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var testBrokerJson = "[\"addr1:9092\",\"addr2:9092\"]"
var testBrokerArr = []string{"addr1:9092", "addr2:9092"}

func TestGetReader(t *testing.T) {
	Convey("Test general reader", t, func() {
		r, err := newReader(nil, "none", "id", 1, ReaderWithAsync())
		So(r, ShouldBeNil)
		So(err, ShouldBeError)

		Convey("Put and try again", func() {

			Convey("Set groupID and partition sametime", func() {
				_, err := newReader(testBrokerArr, "none", "id", 1)
				So(err, ShouldNotBeNil)
			})

			Convey("Set group id", func() {
				_, err := newReader(testBrokerArr, "none", "id", 0)
				So(err, ShouldBeNil)
			})

			Convey("Set partition", func() {
				_, err := newReader(testBrokerArr, "none", "", 1)
				So(err, ShouldBeNil)
			})

			Convey("Get groupConsumer", func() {
				_, err := NewGroupReader(testBrokerArr, "topic", "")
				So(err, ShouldBeError)
				_, err = NewGroupReader(testBrokerArr, "topic", "testgroup")
				So(err, ShouldBeNil)
			})

			Convey("Get partition reader", func() {
				_, err := NewPartitionReader(testBrokerArr, "topic", 1)
				So(err, ShouldBeNil)
			})

		})
	})
}

func TestGetWriter(t *testing.T) {
	Convey("Test general writer", t, func() {
		w, err := newWriter(nil, "none", nil, WriterWithAsync())
		So(w, ShouldBeNil)
		So(err, ShouldBeError)

		Convey("Put and try again", func() {

			Convey("Set groupID and partition sametime", func() {
				_, err := newWriter(testBrokerArr, "none", nil)
				So(err, ShouldBeNil)
			})

			Convey("Get writer", func() {
				_, err := NewWriter(testBrokerArr, "topic")
				So(err, ShouldBeNil)
			})

			Convey("Get partition writer", func() {
				_, err := NewPartitionWriter(testBrokerArr, "topic", 1)
				So(err, ShouldBeNil)
			})

		})
	})
}

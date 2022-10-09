package algs

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestFirstNotEmpty(t *testing.T) {
	Convey("Test first func ", t, func() {
		target := "t"
		Convey("First not empty", func() {
			So(FirstNotEmpty(target, "t2", "t3", ""), ShouldEqual, target)
		})

		Convey("Second", func() {
			So(FirstNotEmpty("", target, "t2", "t3", ""), ShouldEqual, target)
		})

		Convey("All empty", func() {
			So(FirstNotEmpty("", "", ""), ShouldEqual, "")
		})
	})
}

func TestRand(t *testing.T) {
	Convey("Test rand", t, func() {
		So(RandStr(-1, false), ShouldEqual, "")
		So(RandStr(0, false), ShouldEqual, "")
	})
}

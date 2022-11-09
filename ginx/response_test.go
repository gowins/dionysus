package ginx

import (
	"fmt"
	"testing"

	"github.com/gin-gonic/gin/render"
	. "github.com/smartystreets/goconvey/convey"
)

func TestSetDefaultErrorCode(t *testing.T) {
	Convey("set error code", t, func() {
		err := SetDefaultErrorCode(1000)
		So(err, ShouldNotBeNil)
		err = SetDefaultErrorCode(110000)
		So(err, ShouldBeNil)
		testRender := Error(fmt.Errorf("testeror"))
		testJson, ok := testRender.(render.JSON)
		So(ok, ShouldBeTrue)
		data := testJson.Data
		testResponeData, ok := data.(Response)
		So(ok, ShouldBeTrue)
		So(testResponeData.Code, ShouldEqual, 110000)
		testRender = Error(NewGinError(110001, "error msg"))
		testJson, ok = testRender.(render.JSON)
		So(ok, ShouldBeTrue)
		data = testJson.Data
		testResponeData, ok = data.(Response)
		So(ok, ShouldBeTrue)
		So(testResponeData.Code, ShouldEqual, 110001)
		So(testResponeData.Msg, ShouldEqual, "error msg")
	})
}

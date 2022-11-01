package runenv

import (
	"os"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

// TestSetEnv test run enviroment
func TestSetEnv(t *testing.T) {
	convey.Convey("run enviroment", t, func() {
		convey.So(GetRunEnvKey(), convey.ShouldEqual, runEnvKey)
		os.Setenv(runEnvKey, "develop")
		convey.So(IsDev(), convey.ShouldBeTrue)
		convey.So(Not(Develop), convey.ShouldBeFalse)
		os.Setenv(runEnvKey, "test")
		convey.So(IsTest(), convey.ShouldBeTrue)
		os.Setenv(runEnvKey, "gray")
		convey.So(IsGray(), convey.ShouldBeTrue)
		os.Setenv(runEnvKey, "product")
		convey.So(IsProduct(), convey.ShouldBeTrue)
		convey.So(SetRunEnvKey("DIO_ENV"), convey.ShouldBeNil)
		convey.So(runEnvKey, convey.ShouldEqual, "DIO_ENV")
		os.Setenv(runEnvKey, "")
		convey.So(IsDev(), convey.ShouldBeTrue)
		convey.So(SetRunEnvKey(""), convey.ShouldNotBeNil)
	})
}

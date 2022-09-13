package httphystrix_test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/gowins/dionysus/httpclient"
	"github.com/gowins/dionysus/httpclient/httphystrix"
	"github.com/smartystreets/goconvey/convey"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func TestNew(t *testing.T) {
	const URL = "http://www.baidu.com"

	convey.Convey("Test HTTP Client With Hystrix", t, func() {
		convey.Convey("Case: Normal", func() {
			client := httpclient.New(httpclient.WithMiddleware(httphystrix.Middleware()))
			rsp, err := client.Get(URL, nil)
			convey.So(err, convey.ShouldBeNil)
			convey.So(rsp, convey.ShouldNotBeNil)
			convey.So(rsp.Body, convey.ShouldNotBeNil)
			convey.So(rsp.StatusCode, convey.ShouldEqual, http.StatusOK)
			b, err := ioutil.ReadAll(rsp.Body)
			m := make(map[string]string)
			err = json.Unmarshal(b, &m)
			convey.So(err, convey.ShouldBeNil)
			convey.So(m["version"], convey.ShouldNotBeNil)
		})

		convey.Convey("Case: Timeout", func() {
			client := httpclient.New(
				httpclient.WithMiddleware(httphystrix.Middleware(
					httphystrix.WithHystrixTimeout(time.Millisecond),
				)))
			_, err := client.Get(URL, nil)
			convey.So(err, convey.ShouldBeError)
			convey.So(strings.Contains(err.Error(), "timeout"), convey.ShouldBeTrue)
		})
	})
}

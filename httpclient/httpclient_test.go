package httpclient_test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"go.uber.org/atomic"

	"github.com/gowins/dionysus/httpclient"
	"github.com/smartystreets/goconvey/convey"
)

const URL = "http://www.baidu.com"

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func TestNew(t *testing.T) {
	convey.Convey("Test HTTP Client", t, func() {
		convey.Convey("Case: Normal", func() {
			client := httpclient.New()
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

		convey.Convey("Case: Not Found", func() {
			client := httpclient.New()
			rsp, err := client.Get(URL+"1", nil)
			convey.So(err, convey.ShouldBeNil)
			convey.So(rsp, convey.ShouldNotBeNil)
			convey.So(rsp.Body, convey.ShouldNotBeNil)
			convey.So(rsp.StatusCode, convey.ShouldEqual, http.StatusNotFound)
		})

		convey.Convey("Case: Timeout", func() {
			client := httpclient.New(httpclient.WithHTTPTimeout(time.Microsecond))
			rsp, err := client.Get(URL, nil)
			convey.So(err, convey.ShouldNotBeNil)
			convey.So(rsp, convey.ShouldBeNil)
			convey.ShouldContain(err.Error(), "Timeout")
		})

		convey.Convey("Case: Retry Count 3", func() {
			var c atomic.Int32
			client := httpclient.New(
				httpclient.WithHTTPTimeout(time.Microsecond),
				httpclient.WithRetryCount(3),
				httpclient.WithMiddleware(func(doFunc httpclient.DoFunc) httpclient.DoFunc {
					return func(request *http.Request, f func(*http.Response) error) error {
						c.Inc()
						return doFunc(request, f)
					}
				}),
			)

			rsp, err := client.Get(URL, nil)
			convey.So(c.Load(), convey.ShouldEqual, 3)
			convey.So(err, convey.ShouldNotBeNil)
			convey.So(rsp, convey.ShouldBeNil)
			convey.So(strings.Contains(err.Error(), "context deadline exceeded"), convey.ShouldBeTrue)
		})

		convey.Convey("Case: Retrier", func() {
			var c atomic.Int32
			client := httpclient.New(
				httpclient.WithHTTPTimeout(time.Microsecond),
				httpclient.WithRetrier(httpclient.NewRetrier(httpclient.NewConstantBackoff(10*time.Millisecond, 50*time.Millisecond))),
				httpclient.WithRetryCount(3),
				httpclient.WithMiddleware(func(doFunc httpclient.DoFunc) httpclient.DoFunc {
					return func(request *http.Request, f func(*http.Response) error) error {
						c.Inc()
						return doFunc(request, f)
					}
				}),
			)
			rsp, err := client.Get(URL, nil)
			convey.So(c.Load(), convey.ShouldEqual, 3)
			convey.So(err, convey.ShouldNotBeNil)
			convey.So(rsp, convey.ShouldBeNil)
			convey.ShouldContain(err.Error(), "Timeout")
		})
	})
}

func TestTimeout(t *testing.T) {
	convey.Convey("Test HTTP Client Timeout", t, func() {
		client := httpclient.New(
			httpclient.WithHTTPTimeout(time.Nanosecond),
		)

		_, err := client.Get(URL, nil)
		convey.So(err, convey.ShouldNotBeNil)

		client = client.Clone(httpclient.WithHTTPTimeout(time.Second * 5))
		_, err = client.Get(URL, nil)
		convey.So(err, convey.ShouldBeNil)
	})
}

func TestTransport(t *testing.T) {
	convey.Convey("Test HTTP Client Transport", t, func() {
		client := httpclient.New(
			httpclient.WithHTTPTimeout(time.Second),
			httpclient.WithTransport(http.DefaultTransport),
		)

		_, err := client.Get(URL, nil)
		convey.So(err, convey.ShouldBeNil)

		// 这个时候的transport使用的是transport Clone后的对象
		client = client.Clone(httpclient.WithHTTPTimeout(time.Second * 5))
		_, err = client.Get(URL, nil)
		convey.So(err, convey.ShouldBeNil)
	})
}

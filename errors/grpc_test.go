package errors

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestTypeError(t *testing.T) {
	convey.Convey("test custorm error type", t, func() {
		err := New("id", "detail", http.StatusOK)
		convey.So(err.Error(), convey.ShouldEqual, `{"id":"id","code":200,"detail":"detail","status":"OK"}`)
		v, ok := err.(*Error)
		convey.So(ok, convey.ShouldBeTrue)
		v.Code = http.StatusAccepted
		convey.So(err.Error(), convey.ShouldEqual, `{"id":"id","code":202,"detail":"detail","status":"OK"}`)
	})
}

func TestDefined(t *testing.T) {
	convey.Convey("defined", t, func() {
		convey.Convey("bad request", func() {
			err := BadRequest("1", "bad request")
			v, ok := err.(*Error)
			convey.So(ok, convey.ShouldBeTrue)
			convey.So(v.Code, convey.ShouldEqual, http.StatusBadRequest)
		})
		convey.Convey("unauthorized", func() {
			err := Unauthorized("1", "unauthorize request")
			v, ok := err.(*Error)
			convey.So(ok, convey.ShouldBeTrue)
			convey.So(v.Code, convey.ShouldEqual, http.StatusUnauthorized)
		})
		convey.Convey("forbidden", func() {
			err := Forbidden("1", "forbidden")
			v, ok := err.(*Error)
			convey.So(ok, convey.ShouldBeTrue)
			convey.So(v.Code, convey.ShouldEqual, http.StatusForbidden)
		})
		convey.Convey("not found", func() {
			err := NotFound("1", "not found")
			v, ok := err.(*Error)
			convey.So(ok, convey.ShouldBeTrue)
			convey.So(v.Code, convey.ShouldEqual, http.StatusNotFound)
		})
		convey.Convey("method not allowded", func() {
			err := MethodNotAllowed("1", "method not allowed")
			v, ok := err.(*Error)
			convey.So(ok, convey.ShouldBeTrue)
			convey.So(v.Code, convey.ShouldEqual, http.StatusMethodNotAllowed)
		})
		convey.Convey("timeout", func() {
			err := Timeout("1", "timeout")
			v, ok := err.(*Error)
			convey.So(ok, convey.ShouldBeTrue)
			convey.So(v.Code, convey.ShouldEqual, http.StatusRequestTimeout)
		})
		convey.Convey("conflict", func() {
			err := Conflict("1", "conflict")
			v, ok := err.(*Error)
			convey.So(ok, convey.ShouldBeTrue)
			convey.So(v.Code, convey.ShouldEqual, http.StatusConflict)
		})
		convey.Convey("internal server error", func() {
			err := InternalServerError("1", "internal server error")
			v, ok := err.(*Error)
			convey.So(ok, convey.ShouldBeTrue)
			convey.So(v.Code, convey.ShouldEqual, http.StatusInternalServerError)
		})
	})
}

func TestParse(t *testing.T) {
	convey.Convey("parse from json", t, func() {
		j := `{"id":"id","code":200,"detail":"detail","status":"OK"}`
		e := Parse(j)
		convey.So(e.Code, convey.ShouldEqual, http.StatusOK)
		j = `"id":"id","code":200,"detail":"detail","status":"OK"}`
		e1 := Parse(j)
		convey.So(e1.Code, convey.ShouldEqual, 0)
		convey.So(e1.Detail, convey.ShouldEqual, j)
	})
}

func TestIgnorableErr(t *testing.T) {
	convey.Convey("ignorable error", t, func() {
		convey.Convey("wrapper error", func() {
			convey.So(WrapIgnorableError(nil), convey.ShouldBeNil)
			e := WrapIgnorableError(fmt.Errorf("test ignorable"))
			err, ok := e.(*IgnorableError)
			convey.So(ok, convey.ShouldBeTrue)
			convey.So(err.Err, convey.ShouldEqual, "test ignorable")
			convey.So(err.Error(), convey.ShouldEqual, `{"error":"test ignorable","ignorable":true}`)
		})
		convey.Convey("unwrap ignorable error", func() {
			errStr := ""
			ok, _ := UnwrapIgnorableError(errStr)
			convey.So(ok, convey.ShouldBeFalse)
			errStr = `{"error":"test ignorable","ignorable":true}`
			ok, e := UnwrapIgnorableError(errStr)
			convey.So(ok, convey.ShouldBeTrue)
			convey.So(e, convey.ShouldEqual, "test ignorable")
			errStr = `{"error":"test ignorable","ignorable":true`
			ok, e = UnwrapIgnorableError(errStr)
			convey.So(ok, convey.ShouldBeFalse)
			convey.So(e, convey.ShouldEqual, errStr)
		})
	})
}

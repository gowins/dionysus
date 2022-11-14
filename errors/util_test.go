package errors

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/smartystreets/goconvey/convey"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestErr2GrpcCode(t *testing.T) {
	convey.Convey("error to grpc code", t, func() {
		sl := []struct {
			err      error
			expected codes.Code
		}{
			{err: nil, expected: codes.OK},
			{err: io.EOF, expected: codes.OutOfRange},
			{err: io.ErrClosedPipe, expected: codes.FailedPrecondition},
			{err: io.ErrNoProgress, expected: codes.FailedPrecondition},
			{err: io.ErrShortBuffer, expected: codes.FailedPrecondition},
			{err: io.ErrShortWrite, expected: codes.FailedPrecondition},
			{err: io.ErrUnexpectedEOF, expected: codes.FailedPrecondition},
			{err: os.ErrInvalid, expected: codes.InvalidArgument},
			{err: context.Canceled, expected: codes.Canceled},
			{err: context.DeadlineExceeded, expected: codes.DeadlineExceeded},
			{err: os.ErrPermission, expected: codes.PermissionDenied},
			{err: os.ErrExist, expected: codes.AlreadyExists},
			{err: os.ErrNotExist, expected: codes.NotFound},
			{err: fmt.Errorf("error"), expected: codes.Unknown},
		}
		for _, s := range sl {
			convey.So(ConvertGrpcCode(s.err), convey.ShouldEqual, s.expected)
		}
	})
}

func TestHttp2GrpcCode(t *testing.T) {
	convey.Convey("http to grpc code", t, func() {
		sl := []struct {
			err      *Error
			expected codes.Code
		}{
			{err: &Error{Code: http.StatusOK}, expected: codes.OK},
			{err: &Error{Code: http.StatusBadRequest}, expected: codes.InvalidArgument},
			{err: &Error{Code: http.StatusRequestTimeout}, expected: codes.DeadlineExceeded},
			{err: &Error{Code: http.StatusNotFound}, expected: codes.NotFound},
			{err: &Error{Code: http.StatusConflict}, expected: codes.AlreadyExists},
			{err: &Error{Code: http.StatusForbidden}, expected: codes.PermissionDenied},
			{err: &Error{Code: http.StatusUnauthorized}, expected: codes.Unauthenticated},
			{err: &Error{Code: http.StatusPreconditionFailed}, expected: codes.FailedPrecondition},
			{err: &Error{Code: http.StatusNotImplemented}, expected: codes.Unimplemented},
			{err: &Error{Code: http.StatusInternalServerError}, expected: codes.Internal},
			{err: &Error{Code: http.StatusServiceUnavailable}, expected: codes.Unavailable},
			{err: &Error{Code: 1}, expected: codes.Unknown},
		}
		for _, s := range sl {
			convey.So(Http2GrpcCode(s.err), convey.ShouldEqual, s.expected)
		}
	})
}

func TestGrpcAcceptable(t *testing.T) {
	convey.Convey("grpc acceptable", t, func() {
		convey.So(GrpcAcceptable(nil), convey.ShouldBeTrue)
		convey.So(GrpcAcceptable(fmt.Errorf("not grpc internal status error")), convey.ShouldBeTrue)
		sl := []struct {
			code    codes.Code
			msg     string
			expectd bool
		}{
			{code: codes.PermissionDenied, msg: "PermissionDenied", expectd: true},
			{code: codes.DeadlineExceeded, msg: "DeadlineExceeded", expectd: false},
			{code: codes.Internal, msg: "Internal", expectd: false},
			{code: codes.Unavailable, msg: "Unavailable", expectd: false},
			{code: codes.DataLoss, msg: "DataLoss", expectd: false},
		}
		for _, s := range sl {
			sta := status.New(s.code, s.msg)
			convey.So(GrpcAcceptable(sta.Err()), convey.ShouldEqual, s.expectd)
		}
	})
}

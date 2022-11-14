// Package errors provides a way to return detailed information
// for an RPC request error. The error is normally JSON encoded.
package errors

import (
	"encoding/json"
	"fmt"
	"net/http"
)

var _ error = (*Error)(nil)

// Error implements the error interface.
type Error struct {
	Id     string `json:"id"`
	Code   int32  `json:"code"`
	Detail string `json:"detail"`
	Status string `json:"status"`
}

func (e *Error) Error() string {
	b, _ := json.Marshal(e)
	return string(b)
}

// New generates a custom error.
func New(id, detail string, code int32) error {
	return &Error{
		Id:     id,
		Code:   code,
		Detail: detail,
		Status: http.StatusText(int(code)),
	}
}

// Parse tries to parse a JSON string into an error. If that
// fails, it will set the given string as the error detail.
func Parse(err string) *Error {
	e := new(Error)
	err1 := json.Unmarshal([]byte(err), e)
	if err1 != nil {
		e.Detail = err
	}
	return e
}

// BadRequest generates a 400 error.
func BadRequest(id, format string, a ...interface{}) error {
	return &Error{
		Id:     id,
		Code:   http.StatusBadRequest,
		Detail: fmt.Sprintf(format, a...),
		Status: http.StatusText(http.StatusBadRequest),
	}
}

// Unauthorized generates a 401 error.
func Unauthorized(id, format string, a ...interface{}) error {
	return &Error{
		Id:     id,
		Code:   http.StatusUnauthorized,
		Detail: fmt.Sprintf(format, a...),
		Status: http.StatusText(http.StatusUnauthorized),
	}
}

// Forbidden generates a 403 error.
func Forbidden(id, format string, a ...interface{}) error {
	return &Error{
		Id:     id,
		Code:   http.StatusForbidden,
		Detail: fmt.Sprintf(format, a...),
		Status: http.StatusText(http.StatusForbidden),
	}
}

// NotFound generates a 404 error.
func NotFound(id, format string, a ...interface{}) error {
	return &Error{
		Id:     id,
		Code:   http.StatusNotFound,
		Detail: fmt.Sprintf(format, a...),
		Status: http.StatusText(http.StatusNotFound),
	}
}

// MethodNotAllowed generates a 405 error.
func MethodNotAllowed(id, format string, a ...interface{}) error {
	return &Error{
		Id:     id,
		Code:   http.StatusMethodNotAllowed,
		Detail: fmt.Sprintf(format, a...),
		Status: http.StatusText(http.StatusMethodNotAllowed),
	}
}

// Timeout generates a 408 error.
func Timeout(id, format string, a ...interface{}) error {
	return &Error{
		Id:     id,
		Code:   http.StatusRequestTimeout,
		Detail: fmt.Sprintf(format, a...),
		Status: http.StatusText(http.StatusRequestTimeout),
	}
}

// Conflict generates a 409 error.
func Conflict(id, format string, a ...interface{}) error {
	return &Error{
		Id:     id,
		Code:   http.StatusConflict,
		Detail: fmt.Sprintf(format, a...),
		Status: http.StatusText(http.StatusConflict),
	}
}

// InternalServerError generates a 500 error.
func InternalServerError(id, format string, a ...interface{}) error {
	return &Error{
		Id:     id,
		Code:   http.StatusInternalServerError,
		Detail: fmt.Sprintf(format, a...),
		Status: http.StatusText(http.StatusInternalServerError),
	}
}

var _ error = (*IgnorableError)(nil)

type IgnorableError struct {
	Err       string `json:"error"`
	Ignorable bool   `json:"ignorable"`
}

func (e *IgnorableError) Error() string {
	b, _ := json.Marshal(e)
	return string(b)
}

// UnwrapIgnorableError tries to parse a JSON string into an error. If that
// fails, it will set the given string as the error detail.
func UnwrapIgnorableError(err string) (bool, string) {
	if err == "" {
		return false, err
	}

	igerr := new(IgnorableError)
	uerr := json.Unmarshal([]byte(err), igerr)
	if uerr != nil {
		return false, err
	}

	return igerr.Ignorable, igerr.Err
}

// IgnoreError generates a -1 error.
func WrapIgnorableError(err error) error {
	if err == nil {
		return nil
	}

	return &IgnorableError{
		Err:       err.Error(),
		Ignorable: true,
	}
}

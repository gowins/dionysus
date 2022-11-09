package ginx

import (
	"fmt"
	"github.com/gin-gonic/gin/render"
)

var (
	ginxOK               = 10000
	ginxOKMsg            = "success"
	ginxDefaultError     = 10001
	ginxDefaultErrorMsg  = "default gin error"
	ginxLimitingCode     = 10002
	ginxLimitingMsg      = "too many requests"
	businessMinErrorCode = 100000
)

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func ResponseData(code int, msg string, data interface{}) Render {
	return render.JSON{Data: Response{
		Code: code,
		Msg:  msg,
		Data: data,
	}}
}

func Success(data interface{}) Render {
	return ResponseData(ginxOK, ginxOKMsg, data)
}

func Error(err error) Render {
	if ge, ok := err.(GinError); ok {
		return ResponseData(ge.Code, ge.Error(), struct{}{})
	}
	return ResponseData(ginxDefaultError, ginxDefaultErrorMsg+": "+err.Error(), struct{}{})
}

func SetDefaultErrorCode(code int) error {
	if code < businessMinErrorCode {
		return fmt.Errorf("business error code should >= %v", businessMinErrorCode)
	}
	ginxDefaultError = code
	return nil
}

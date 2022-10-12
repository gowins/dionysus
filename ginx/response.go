package ginx

import (
	"github.com/gin-gonic/gin/render"
)

const (
	OK             = 0
	ParamError     = 100
	DefaultError   = 400
	SpeedLimit     = 429
	ServerInterval = 500
)

var (
	CodeMsgMap = map[int]string{
		OK:             "请求成功",
		ParamError:     "参数错误",
		DefaultError:   "请求失败",
		SpeedLimit:     "服务器正忙，请稍后再试",
		ServerInterval: "服务器错误",
	}
)

var (
	GinErrorParams = NewGinError(ParamError, CodeMsgMap[ParamError])
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
	return ResponseData(OK, CodeMsgMap[OK], data)
}

func Error(err error) Render {
	if ge, ok := err.(GinError); ok {
		return ResponseData(ge.Code, ge.Error(), struct{}{})
	}
	return ResponseData(DefaultError, CodeMsgMap[DefaultError], struct{}{})
}

package ginhelper

import (
	"github.com/gin-gonic/gin/render"
)

const (
	OK             = 200
	SpeedLimit     = -2
	ServerInterval = 500
	DefaultError   = 400
	ParamError     = 100
)

var (
	CodeMsgMap = map[int]string{
		SpeedLimit:     "服务器正忙，请稍后再试",
		ServerInterval: "服务器错误，请联系客服",
		OK:             "请求成功",
		DefaultError:   "请求失败",
		ParamError:     "参数错误",
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

func ResponseData(code int, msg string, data interface{}) render.Render {
	return render.JSON{Data: Response{
		Code: code,
		Msg:  msg,
		Data: data,
	}}
}

func Success(data interface{}) render.Render {
	return ResponseData(OK, CodeMsgMap[OK], data)
}

func Error(err error) render.Render {
	if ge, ok := err.(GinError); ok {
		return ResponseData(ge.Code, ge.Error(), struct{}{})
	}
	return ResponseData(DefaultError, CodeMsgMap[DefaultError], struct{}{})
}

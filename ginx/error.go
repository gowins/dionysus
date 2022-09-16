package ginx

type GinError struct {
	Code int
	Msg  string
}

func NewGinError(code int, msg string) GinError {
	return GinError{
		Code: code,
		Msg:  msg,
	}
}

func (e GinError) Error() string {
	return e.Msg
}

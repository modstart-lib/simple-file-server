package result

import (
	"simple-file-server/lib/defs"
)

func Generate(code int, msg string, data interface{}) defs.Response {
	return defs.Response{
		Code: code,
		Msg:  msg,
		Data: data,
	}
}

func GenerateSuccess(msg string) defs.Response {
	return Generate(0, msg, nil)
}

func GenerateSuccessWithData(msg string, data interface{}) defs.Response {
	return Generate(0, msg, data)
}

func GenerateSuccessData(data interface{}) defs.Response {
	return Generate(0, "ok", data)
}

func GenerateError(msg string) defs.Response {
	return Generate(-1, msg, nil)
}

package console

import (
	"encoding/json"
	"fmt"
	"simple-file-server/lib/defs"
)

func Generate(code int, msg string, data interface{}) {
	res := defs.Response{
		Code: code,
		Msg:  msg,
		Data: data,
	}
	resJson, _ := json.MarshalIndent(res, "", "    ")
	fmt.Println(string(resJson))
}

func GenerateSuccess(msg string) {
	Generate(0, msg, nil)
}

func GenerateSuccessData(data interface{}) {
	Generate(0, "ok", data)
}

func GenerateError(msg string) {
	Generate(-1, msg, nil)
}

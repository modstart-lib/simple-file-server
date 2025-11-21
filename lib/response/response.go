package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"simple-file-server/lib/defs"
)

func Generate(ctx *gin.Context, code int, msg string, data interface{}) {
	if data == nil {
		data = gin.H{}
	}
	res := defs.Response{
		Code: code,
		Msg:  msg,
		Data: data,
	}
	ctx.Header("Transfer-Encoding", "identity")
	ctx.JSON(http.StatusOK, res)
	ctx.Abort()
}

func GenerateSuccess(ctx *gin.Context, msg string) {
	Generate(ctx, 0, msg, nil)
}

func GenerateSuccessWithData(ctx *gin.Context, msg string, data interface{}) {
	Generate(ctx, 0, msg, data)
}

func GenerateSuccessData(ctx *gin.Context, data interface{}) {
	Generate(ctx, 0, "ok", data)
}

func GenerateError(ctx *gin.Context, msg string) {
	Generate(ctx, -1, msg, nil)
}

package misc

import "github.com/kataras/iris/context"

const (
	StandardTimeFormatPattern = "2006-01-02 15:04:05"
)

const (
	CodeSuccess = 1001
	CodeFailure = 1002
)

var ResponseCode = map[int]string{
	CodeSuccess: "成功",
	CodeFailure: "失败，原因：%s",
}

var JsonOption = context.JSON{}
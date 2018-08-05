package misc

const (
	StandardTimeFormatPattern = "2006-01-02 15:04:05"
)

const (
	CodeSuccess       = 1001
	CodeFailure       = 1002
	CodeParamMissing  = 1003
	CodeParamInvalid  = 1004
	CodeDuplicateData = 1005
	CodeTryAgainLater = 0
)

var ResponseCode = map[int]string{
	CodeSuccess:       "成功",
	CodeFailure:       "失败，原因：%s",
	CodeParamMissing:  "缺少参数：%s",
	CodeParamInvalid:  "非法参数：%s",
	CodeDuplicateData: "%s已存在",
	CodeTryAgainLater: "系统繁忙，请稍候再试",
}

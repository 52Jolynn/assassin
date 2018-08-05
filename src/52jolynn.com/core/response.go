package core

import (
	"52jolynn.com/misc"
	"fmt"
)

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func CreateResponse(code int, msgArg ...interface{}) *Response {
	return CreateResponseWithData(code, nil, msgArg...)
}

func CreateResponseWithData(code int, data interface{}, msgArg ...interface{}) *Response {
	if msgArg != nil && len(msgArg) > 0 {
		return &Response{Code: code, Msg: fmt.Sprintf(misc.ResponseCode[code], msgArg...), Data: data}
	}
	return &Response{Code: code, Msg: misc.ResponseCode[code], Data: data}
}

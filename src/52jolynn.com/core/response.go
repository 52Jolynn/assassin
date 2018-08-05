package core

type Response struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

type ResponseWithData struct {
	Response
	Data interface{} `json:"data"`
}

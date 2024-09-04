package bm

import "github.com/gin-gonic/gin"

type Response struct {
	Code     int         `json:"code"`
	Data     interface{} `json:"data,omitempty"`
	Msg      string      `json:"msg"`
	Callback string      `json:"callback"`
}

func (response *Response) Suc(data interface{}, msg string) *Response {
	if msg == "" {
		msg = "操作成功"
	}

	response.Code = 200
	response.Data = data
	response.Msg = msg

	return response
}

func (response *Response) FailBackend(msg string) *Response {
	response.Code = 500
	response.Msg = msg

	return response
}

func (response *Response) FailFront(msg string) *Response {
	response.Code = 400
	response.Msg = msg

	return response
}

func (response *Response) Send(c *gin.Context) {
	if response.Code == 0 {
		c.JSON(500, response)
	}
	c.JSON(200, response)
}

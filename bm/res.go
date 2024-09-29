package bm

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type Res struct {
	Code     int         `json:"code"`
	Data     interface{} `json:"data,omitempty"`
	Msg      string      `json:"msg"`
	Callback string      `json:"callback"`

	filePath string
	hopeName string
}

func NewRes() *Res {
	return &Res{}
}

func (response *Res) SucJson(data interface{}, msg ...any) *Res {

	if msg == nil {
		msg = append(msg, "操作成功")
	}

	response.Code = 200
	response.Data = data
	response.Msg = fmt.Sprint(msg...)

	return response
}

func (response *Res) SucFile(filePath string, fileName string, msg ...any) *Res {
	if msg == nil {
		msg = append(msg, "下载成功")
	}

	response.Code = 200
	response.filePath = filePath
	response.hopeName = fileName
	response.Msg = fmt.Sprint(msg...)

	return response
}

func (response *Res) SucList(data ResList, msg ...any) *Res {
	if msg == nil {
		msg = append(msg, "下载成功")
	}

	response.Code = 200
	response.Data = data
	response.Msg = fmt.Sprint(msg...)

	return response
}

func (response *Res) FailBackend(msg ...any) *Res {
	response.Code = 500
	response.Msg = fmt.Sprint(msg...)

	return response
}

func (response *Res) FailFront(msg ...any) *Res {
	response.Code = 400
	response.Msg = fmt.Sprint(msg...)

	return response
}

func (response *Res) Send(c *gin.Context) {
	if response.Code == 0 {
		c.JSON(500, response)
	}
	if response.filePath != "" {
		if response.hopeName == "" {
			response.hopeName = response.filePath
		}
		// 设置下载的标头
		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Transfer-Encoding", "binary")
		c.Header("Content-Disposition", "attachment; filename="+response.hopeName)
		c.Header("Content-Type", "application/octet-stream")
		c.Header("Content-Length", "0")
		// 将文件作为响应发送
		c.File(response.filePath)
	} else {
		c.JSON(200, response)
	}
}

func (response *Res) SendAbort(c *gin.Context) {
	response.Send(c)
	c.Abort()
}

type ResList struct {
	PageSize int         `json:"page_size"`
	Current  int         `json:"current"`
	Data     interface{} `json:"data"`
	Total    int64       `json:"total"`
}

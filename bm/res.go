package bm

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Res struct {
	Code     int         `json:"code"`
	Data     interface{} `json:"data,omitempty"`
	Msg      string      `json:"msg"`
	Callback string      `json:"callback"`

	filePath string              `json:"-"`
	hopeName string              `json:"-"`
	w        http.ResponseWriter `json:"-"`
}

type ResList struct {
	Pagination
	Data  interface{} `json:"data"`
	Total int64       `json:"total"`
}

const (
	ContentTypeJSON            = "application/json"
	ContentTypeOctetStream     = "application/octet-stream"
	ContentDescription         = "Content-Description"
	ContentTransferEncoding    = "Content-Transfer-Encoding"
	ContentDisposition         = "Content-Disposition"
	DefaultSuccessMessage      = "操作成功"
	DefaultDownloadMessage     = "下载成功"
	DefaultFailBackendMessage  = "网络错误"
	DefaultFailFrontendMessage = "客户端错误"
)

func NewRes(w http.ResponseWriter) *Res {
	return &Res{w: w}
}

func (r *Res) SucJson(data interface{}, msg ...interface{}) *Res {
	r.Code = http.StatusOK
	r.Data = data
	r.Msg = formatMessage(msg, DefaultSuccessMessage)
	return r
}

func (r *Res) SucFile(filePath, hopeName string, msg ...interface{}) *Res {
	r.Code = http.StatusOK
	r.filePath = filePath
	r.hopeName = hopeName
	r.Msg = formatMessage(msg, DefaultDownloadMessage)
	return r
}

func (r *Res) SucList(data ResList, msg ...interface{}) *Res {
	r.Code = http.StatusOK
	r.Data = data
	r.Msg = formatMessage(msg, DefaultSuccessMessage)
	return r
}

func (r *Res) FailBackend(msg ...interface{}) *Res {
	r.Code = http.StatusInternalServerError
	r.Msg = formatMessage(msg, DefaultFailBackendMessage)
	return r
}

func (r *Res) FailFront(msg ...interface{}) *Res {
	r.Code = http.StatusBadRequest
	r.Msg = formatMessage(msg, DefaultFailFrontendMessage)
	return r
}
func (r *Res) Send() {
	if r.Code == 0 {
		r.sendError(http.StatusInternalServerError, "Internal Server Error")
		return
	}

	if r.filePath != "" {
		r.sendFile()
	} else {
		r.sendJSON()
	}
}

func (r *Res) sendError(status int, message string) {
	http.Error(r.w, message, status)
}

func (r *Res) sendFile() {
	if r.hopeName == "" {
		r.hopeName = r.filePath
	}

	r.w.Header().Set(ContentDescription, "File Transfer")
	r.w.Header().Set(ContentTransferEncoding, "binary")
	r.w.Header().Set(ContentDisposition, "attachment; filename="+r.hopeName)
	r.w.Header().Set("Content-Type", ContentTypeOctetStream)

	http.ServeFile(r.w, nil, r.filePath)
}

func (r *Res) sendJSON() {
	r.w.Header().Set("Content-Type", ContentTypeJSON)
	r.w.WriteHeader(200)
	json.NewEncoder(r.w).Encode(r)
}

func formatMessage(msg []interface{}, defaultMsg string) string {
	if len(msg) == 0 {
		return defaultMsg
	}
	return fmt.Sprint(msg...)
}

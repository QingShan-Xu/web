package bm

type Pagination struct {
	PageSize int         `form:"page_size,default=10" json:"page_size,default=10"`
	Current  int         `form:"current,default=1" json:"current,default=1"`
	Data     interface{} `json:"data"`
	Total    int         `json:"total"`
}

package bm

type Pagination struct {
	PageSize int `bind:"page_size"`
	Current  int `bind:"current"`
}

// SortBy 为 ORDER 值的下标 +1
//
// Sort 为排序方式 ascend/descend, 代指数据库中的 ASC/DESC, 代指 升序/降序
type Order struct {
	Sort []Sort `form:"sort" json:"sort"`
}

type Sort struct {
	SortBy int    `form:"sort_by" json:"sort_by" bindind:"min=1"`
	Sort   string `form:"sort" json:"sort" bindind:"containsrune=descend,ascend"`
}

var SortMapOrder = map[string]string{
	"ascend":  "ASC",
	"descend": "DESC",
}

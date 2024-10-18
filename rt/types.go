package rt

import (
	"net/http"
)

// 路由
type Router struct {
	Name          string                            // 路由名称
	Path          string                            // 路由路径
	Middlewares   []func(http.Handler) http.Handler // 中间件
	Children      []Router                          // 子路由
	Method        string                            // 请求方法
	Handle        func()                            // 处理函数
	Bind          interface{}                       // 参数绑定 tag有 uri, query, form, json
	MODEL         interface{}                       // 数据库模型
	NoAutoMigrate bool                              // 不自动迁移该模型
	SCOPES        []Scope
	WHERE         [][]string

	CREATE_ONE map[string]string // 创建 map[orm字段名]bind字段名
	UPDATE_ONE map[string]string // 更新 map[orm字段名]bind字段名, nil 不更新
	DELETE_ONE bool              // 删除
	GET_ONE    bool              // 获取
	GET_LIST   bool              // 获取列表

	completePath string // 完整路径，从根节点到当前节点。
	completeName string // 完整名称，从根节点到当前节点。
	completeInfo string // 格式化的路由树结构信息。
}

// type BindData struct {
// }

// func (bindData *BindData) ReaderFile(key string) (multipart.File, *multipart.FileHeader, error)
// func (bindData *BindData) Reader(key string) (interface{}, error)

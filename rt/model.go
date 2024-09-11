package rt

import (
	"github.com/QingShan-Xu/xjh/bm"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// METHOD 是一个包含 HTTP 方法常量的结构体，用于在代码中统一引用 HTTP 方法名。
//
// 定义了常见的四种 HTTP 方法：GET、POST、PUT、DELETE，
var METHOD = struct {
	GET    string
	POST   string
	PUT    string
	DELETE string
}{
	GET:    "GET",
	POST:   "POST",
	PUT:    "PUT",
	DELETE: "DELETE",
}

var TYPE = struct {
	GET_LIST   string
	GET_ONE    string
	UPDATE_ONE string
	// UPDATE_LIST string
	CREATE_ONE string
	// CREATE_LIST string
	DELETE_ONE string
	// DELETE_LIST string
}{
	GET_LIST:   "GET_LIST",
	GET_ONE:    "GET_ONE",
	UPDATE_ONE: "UPDATE_ONE",
	// UPDATE_LIST: "UPDATE_LIST",
	CREATE_ONE: "CREATE_ONE",
	// CREATE_LIST: "CREATE_LIST",
	DELETE_ONE: "DELETE_ONE",
	// DELETE_LIST: "DELETE_LIST",
}

type Router struct {
	// 通用 路由组参数

	// 路由名称
	Name string
	// 路由路径 层级嵌套, 不需要手动添加 "/"
	Path string
	// 中间件 符合gin的中间件
	Middlewares []gin.HandlerFunc
	// 子路由
	Children []Router

	// 子路由
	Method  string      // 请求方法 rt.Method.GET
	Handler Handler     // 处理函数
	Bind    interface{} // 请求参数 struct{ Name string `json:"name" form:"name"` }

	// 数据库: 链式条件
	MODEL interface{}
	WHERE map[string]string

	Type string

	BeforeFinisher func(bind interface{}) interface{}
}

type Handler func(C *gin.Context, TX *gorm.DB, bind interface{}) (res *bm.Res)

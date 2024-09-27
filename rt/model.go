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
	Method          string                                                      // 请求方法 rt.Method.GET
	Handler         func(C *gin.Context, TX *gorm.DB, bind interface{}) *bm.Res // 处理函数
	OriginalHandler gin.HandlerFunc                                             // 原始gin处理函数
	Bind            interface{}                                                 // 请求参数 struct{ Name string `json:"name" form:"name"` }

	// 数据库: 链式条件
	NoAutoMigrate bool // 不自动迁移, 默认是自动迁移
	MODEL         interface{}
	// key: 数据库字段与条件, value: bind字段
	//
	// 例:
	//	map[string]string{
	//		"name", "Name"
	//		"age <> ?", "AgeNum"
	//	}
	WHERE map[string]string
	// 每一项为 数据库 字段
	//
	// 例:
	//	[]string{
	//		"pet.weight",	// 下标为 0 故 Sort[number][SortBy]中值为 1
	//		"age",		// 下标为 1 故 Sort[number][SortBy]中值为 2
	//	}
	ORDER []string
	// key: orm字段, value: bind字段
	//
	// 例:
	//	[]string{
	//		"Name",
	//		"Age",
	//		"Pet.Name",
	//	}
	SELECT []string
	// orm字段
	//
	// 例:
	//	[]string{
	//		"Pet.Shop",
	//		"Work",
	//	}
	PRELOAD []string
	// JOIN语句
	//
	// 例:
	//	[]string{
	//		"LEFT JOIN pet ON pet.user_id = user.id",
	//		"work",
	//	}
	JOINS []string

	Type string
	// key: orm字段, value: bind字段
	//
	// 例:
	//	[]string{
	//		"Name": "Name",
	//		"Work": "Work",
	//	}
	UPDATE map[string]string
	// key: orm字段, value: bind字段
	//
	// 例:
	//	[]string{
	//		"Name": "Name",
	//		"Work": "Work",
	//	}
	CREATE map[string]string
}

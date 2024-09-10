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

// Finisher 是一个包含常量的结构体，用于定义数据库操作的结尾句类型。
var Finisher = struct {
	//	示例:
	//
	//		tx := ctx.MustGet("reqTX_").(*gorm.DB) // tx := 在 Router 中定义的一系列 gorm 的 Chain Methods
	//		result := tx.Find(&data)
	//		if result.Error != nil { // 后端错误
	//			new(bm.Res).FailBackend(result.Error).Send(ctx)
	//			ctx.Abort()
	//			return
	//		}
	//		if result.RowsAffected == 0 { // 前端错误
	//			new(bm.Res).FailFront("数据不存在").Send(ctx)
	//			ctx.Abort()
	//			return
	//		}
	//		new(bm.Res).SucJson(data).Send(ctx) // 返回数据
	//		ctx.Abort() // 阻止后续 Handler, 即阻止 Router 中的 自定义Handler
	First string
	//	示例:
	//
	//		tx := ctx.MustGet("reqTX_").(*gorm.DB)  // tx := 在 Router 中定义的一系列 gorm 的 Chain Methods
	//		bind := ctx.MustGet("reqBind_") //  // bind := 在 Router 中定义的 Bind
	//		result := tx.Create(bind)
	//		if result.Error != nil { // 后端错误
	//			new(bm.Res).FailBackend(result.Error).Send(ctx)
	//			ctx.Abort()
	//			return
	//		}
	//		new(bm.Res).SucJson(bind).Send(ctx) // 返回数据
	//		ctx.Abort() // 阻止后续 Handler, 即阻止 Router 中的 自定义Handler
	Create string

	Update string

	Delete string
}{
	First:  "First",
	Create: "Create",
	Update: "Update",
	Delete: "Delete",
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

	// 数据库: 结尾句
	Finisher       string
	BeforeFinisher func(bind interface{}) interface{}
}

type Handler func(C *gin.Context, TX *gorm.DB, bind interface{}) (res *bm.Res)

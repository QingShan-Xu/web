package rt

import (
	"log"

	"github.com/QingShan-Xu/xjh/rt/internal/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// register 函数用于递归地注册路由，支持分组路由和单个路由的注册。
// 如果 regRouter 表示一个路由组，则会递归地为每个子路由进行注册。
// 否则，直接调用 registerRouter 函数来注册单个路由。
//
// 参数:
//   - pGroupRouter (*gin.RouterGroup): 当前路由组，用于注册新的路由或子路由组。
//   - regRouter (*Router): 包含路由的路径、方法、中间件及子路由信息的结构体。
//
// 示例:
//
//	register(routerGroup, &Router{
//	    Path: "/api",
//	    Children: []Router{
//	        {
//	            Method: "GET",
//	            Path:   "/items",
//	            Handler: func(ctx *gin.Context, tx *gorm.DB) Response {
//	                // 处理请求的代码
//	            },
//	        },
//	    },
//	})
func register(pGroupRouter *gin.RouterGroup, regRouter *Router) {
	// 检查 regRouter 是否表示一个路由组。
	isGRoup := isGroupRouter(regRouter)

	if isGRoup {
		// 如果是路由组，则创建一个新的子路由组。
		groupRouter := pGroupRouter.Group(regRouter.Path, regRouter.Middlewares...)

		// 如果该路由组没有子路由，则直接返回。
		if len(regRouter.Children) == 0 {
			return
		}

		// 递归地为每个子路由注册。
		for _, child := range regRouter.Children {
			register(groupRouter, &child)
		}
	} else {
		// 如果不是路由组，直接注册单个路由。
		registerRouter(pGroupRouter, regRouter)
	}
}

// registerRouter 函数用于在 Gin 框架中注册一个路由。
// 它通过传入的 regRouter 参数指定的路径和方法将中间件和处理程序函数绑定到指定的路由组 pGroupRouter 上。
//
// 参数:
//   - pGroupRouter (*gin.RouterGroup): 路由组，用于组织一组相关路由。
//   - regRouter (*Router): 包含路由的路径、方法、中间件和处理程序等信息的结构体。
//
// 示例:
//
//	registerRouter(routerGroup, &Router{
//	    Method: "GET",
//	    Path:   "/example",
//	    Handler: func(ctx *gin.Context, tx *gorm.DB) Response {
//	        // 处理请求的代码
//	    },
//	})
//
// 函数逻辑:
//  1. 判断 regRouter.Name 是否为空，如果不为空则使用 regRouter.Name 作为路由名称，否则使用 regRouter.Path。
//  2. 如果 pGroupRouter、regRouter 为空，或 regRouter.Method 为空，则记录警告日志并跳过路由注册。
//  3. 根据 regRouter 的配置添加相应的中间件函数到 middlewareFuncs 列表。
//     - 如果 regRouter.Bind 不为空，添加请求绑定中间件。
//     - 如果 regRouter.WHERE 不为空，添加数据库预处理中间件。
//     - 如果 regRouter.Finisher 不为空且 regRouter.Handler 为空，添加请求终结器中间件。
//     - 添加用户定义的其他中间件。
//  4. 如果 regRouter.Handler 不为空，添加处理程序函数，该函数在数据库事务上下文中处理请求并发送响应。
//  5. 最终通过 pGroupRouter.Handle 方法将路由与所有中间件和处理程序函数绑定到一起，并注册到 Gin 路由器中。
//
// 注意事项:
//   - 如果 regRouter.Finisher 和 regRouter.Handler 同时不为空，将触发致命日志并终止程序，
//     因为这两者不能同时使用。
func registerRouter(pGroupRouter *gin.RouterGroup, regRouter *Router) {
	var name string

	// 判断路由名称，如果 regRouter.Name 不为空，则使用它作为路由名称，否则使用路径作为名称。
	if regRouter.Name != "" {
		name = regRouter.Name
	} else {
		name = regRouter.Path
	}

	// 检查路由组、路由结构体和方法是否有效，如果无效则记录日志并跳过路由注册。
	if pGroupRouter == nil || regRouter == nil || regRouter.Method == "" {
		log.Printf("%s: 没有 Path 或 Method, 已跳过路由注册", name)
		return
	}

	// 初始化中间件函数列表，用于存储所有需要应用的中间件。
	middlewareFuncs := make([]gin.HandlerFunc, 0)

	// 如果有绑定参数，添加请求绑定中间件。
	if regRouter.Bind != nil {
		middlewareFuncs = append(middlewareFuncs, middleware.ReqBindMiddleware(regRouter.Bind, name))
	}

	// 如果有 WHERE 条件，添加数据库预处理中间件。
	if regRouter.WHERE != nil {
		middlewareFuncs = append(middlewareFuncs, middleware.ReqPreDBMiddleware(
			regRouter.WHERE,
			regRouter.Bind,
			regRouter.MODEL,
			name,
		))
	}

	// 如果有终结器且没有处理程序，添加请求终结器中间件。
	if regRouter.Finisher != "" {
		if regRouter.Handler != nil {
			log.Fatalf("%s: 不能同时使用 Handler 和 Finisher", name)
		}

		middlewareFuncs = append(middlewareFuncs, middleware.ReqFinisherMiddleware(
			regRouter.Finisher,
			regRouter.MODEL,
			name,
		))
	}

	// 添加用户定义的额外中间件。
	if len(regRouter.Middlewares) > 0 {
		middlewareFuncs = append(middlewareFuncs, regRouter.Middlewares...)
	}

	// 如果定义了处理程序函数，将其添加到中间件函数列表的最后。
	if regRouter.Handler != nil {
		middlewareFuncs = append(middlewareFuncs, func(ctx *gin.Context) {
			// 从上下文中获取数据库事务对象。
			TX := ctx.MustGet("reqTX_").(*gorm.DB)
			// 执行处理程序并获取响应。
			res := regRouter.Handler(ctx, TX)
			// 如果响应已经写入，直接返回。
			if ctx.Writer.Written() {
				return
			}
			// 发送响应。
			res.Send(ctx)
		})
	}

	// 最终，将所有中间件和处理程序绑定到指定的路由路径和方法上，并注册到路由组中。
	pGroupRouter.Handle(
		regRouter.Method,
		regRouter.Path,
		middlewareFuncs...,
	)
}

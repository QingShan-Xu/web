// Package rt 提供了路由器的定义和注册功能。
package rt

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Router 定义了路由器结构体。
type Router struct {
	Name          string                            // 路由名称
	Path          string                            // 路由路径
	Method        string                            // 请求方法
	Handler       http.HandlerFunc                  // 处理函数
	Bind          interface{}                       // 请求参数绑定结构体
	Model         interface{}                       // 数据库模型
	NoAutoMigrate bool                              // 是否自动迁移模型
	Middlewares   []func(http.Handler) http.Handler // 中间件列表
	Children      []Router                          // 子路由列表

	Scopes []Scope
	Where  [][]string

	CreateFields map[string]string // 创建操作字段映射
	UpdateFields map[string]string // 更新操作字段映射
	Delete       bool              // 是否为删除操作
	GetOne       bool              // 是否获取单个记录
	GetList      bool              // 是否获取列表

	completePath string // 完整路径
	completeName string // 完整名称
	completeInfo string // 路由信息
}

// Register 函数注册路由并返回 chi.Router。
// rootRouter: 根路由器。
// 返回生成的 chi.Router 或错误信息。
func Register(rootRouter *Router) (chi.Router, error) {
	if rootRouter == nil {
		return nil, fmt.Errorf("root router cannot be nil")
	}
	if !isGroup(*rootRouter) {
		return nil, fmt.Errorf("root router must be a group router")
	}

	// 初始化完整路径和名称。
	initCompletePathAndName(rootRouter)

	chiRouter := chi.NewRouter()
	if err := generateChiRouter(rootRouter, chiRouter); err != nil {
		return nil, fmt.Errorf("error generating router: %w", err)
	}

	// 自动迁移数据库模型。
	if err := generateDBModel(rootRouter); err != nil {
		return nil, err
	}

	// 生成查询条件。
	if err := generateQuery(rootRouter); err != nil {
		return nil, err
	}

	// 显示完整的路由信息。
	displayCompleteInfo(rootRouter)

	return chiRouter, nil
}

// generateChiRouter 递归生成 chi.Router。
// currentRouter: 当前处理的路由器。
// parentChiRouter: 父级 chi.Router。
func generateChiRouter(currentRouter *Router, parentChiRouter chi.Router) error {
	if currentRouter == nil {
		return fmt.Errorf("router cannot be nil")
	}

	if isGroup(*currentRouter) {
		// 为当前组定义一个新的子路由。
		parentChiRouter.Route(currentRouter.Path, func(subRouter chi.Router) {
			// 应用中间件。
			subRouter.Use(currentRouter.Middlewares...)

			// 递归处理子路由。
			for i := range currentRouter.Children {
				child := &currentRouter.Children[i]
				if err := generateChiRouter(child, subRouter); err != nil {
					log.Printf("Error processing child router %s(%s): %v", child.completePath, child.completeName, err)
					return
				}
			}
		})
	} else {
		// 如果未指定处理函数，默认使用 ServeHTTP 方法。
		if currentRouter.Handler == nil {
			currentRouter.Handler = currentRouter.ServeHTTP
		}
		parentChiRouter.Method(currentRouter.Method, currentRouter.Path, currentRouter.Handler)
	}

	return nil
}

// ServeHTTP 实现 http.Handler 接口。
// w: HTTP 响应写入器。
// req: HTTP 请求。
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	handler := &Handler{Router: r}
	handler.ServeHTTP(w, req)
}

// isGroup 检查路由器是否为组路由。
// router: 路由器。
// 返回是否为组路由的布尔值。
func isGroup(router Router) bool {
	return len(router.Children) > 0
}

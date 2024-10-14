package rt

import (
	"fmt"
	"log"
	"reflect"

	"github.com/go-chi/chi/v5"
)

// Register 函数根据给定的 Router 结构初始化一个 chi.Router。
// 确保提供的路由器是一个组路由器，并生成一个 chi.Router。
// 参数:
// - rootRouter: 用于生成 chi.Router 的 Router 配置。
// 返回:
// - (chi.Router, error): 配置好的 chi.Router 和可能出现的错误。
func Register(rootRouter *Router) (chi.Router, error) {
	if rootRouter == nil {
		return nil, fmt.Errorf("根路由器不能为空")
	}
	if !isGroup(*rootRouter) {
		return nil, fmt.Errorf("根路由器必须是一个组路由器")
	}

	// 初始化根节点及其子节点的完整路径和名称。
	initCompletePathAndName(rootRouter)

	chiRouter := chi.NewRouter()
	if err := generateChiRouter(rootRouter, rootRouter, chiRouter); err != nil {
		return nil, fmt.Errorf("生成路由时出错: %w", err)
	}

	if err := genDBModel(rootRouter); err != nil {
		return nil, err
	}

	if err := genQuery(rootRouter); err != nil {
		return nil, err
	}

	// 输出给定路由节点的完整树形结构信息。
	displayCompleteInfo(rootRouter)

	return chiRouter, nil
}

func genQuery(currentRouter *Router) error {
	if isGroup(*currentRouter) {
		for i := range currentRouter.Children {
			child := &currentRouter.Children[i]
			if err := genQuery(child); err != nil {
				return err
			}
		}
	}

	if currentRouter.MODEL == nil {
		return nil
	}

	query := NewQuery()

	if currentRouter.WHERE != nil {
		for _, where := range currentRouter.WHERE {
			scope, err := query.WHERE(where)
			if err != nil {
				return err
			}
			currentRouter.SCOPES = append(currentRouter.SCOPES, scope)
		}
	}

	return nil
}

func genDBModel(currentRouter *Router) error {
	if isGroup(*currentRouter) {
		for i := range currentRouter.Children {
			child := &currentRouter.Children[i]
			if err := genDBModel(child); err != nil {
				return err
			}
		}
	}

	if currentRouter.MODEL == nil || currentRouter.NoAutoMigrate {
		return nil
	}
	// 迁移
	if err := DB.AutoMigrate(reflect.New(reflect.Indirect(reflect.ValueOf(currentRouter.MODEL)).Type()).Interface()); err != nil {
		return fmt.Errorf("%s (%s): gorm AutoMigrate err: %v", currentRouter.Path, currentRouter.Name, err)
	}
	return nil
}

// generateChiRouter 递归地根据 Router 配置构建 chi 路由。
// 参数:
// - currentRouter: 当前处理的 Router 节点。
// - parentChiRouter: 将向其添加路由的 chi.Router 实例。
func generateChiRouter(rootRouter *Router, currentRouter *Router, parentChiRouter chi.Router) error {
	if currentRouter == nil {
		return fmt.Errorf("注册路由不能为空")
	}

	if isGroup(*currentRouter) {
		// 为当前组定义一个新的子路由。
		parentChiRouter.Route(currentRouter.Path, func(subRouter chi.Router) {

			// 应用中间件
			for _, mw := range currentRouter.Middlewares {
				subRouter.Use(mw)
			}

			for i := range currentRouter.Children {
				child := &currentRouter.Children[i]
				if err := generateChiRouter(rootRouter, child, subRouter); err != nil {
					// 记录错误，但继续处理其他子路由
					log.Printf("处理子路由 %s(%s) 时出错: %v", child.completePath, child.completeName, err)
				}
			}
		})
	} else {
		parentChiRouter.Method(currentRouter.Method, currentRouter.Path, currentRouter)
	}

	return nil
}

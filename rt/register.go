package rt

import (
	"log"

	"github.com/QingShan-Xu/xjh/cf"
	"github.com/gin-gonic/gin"
)

func register(pGroupRouter *gin.RouterGroup, regRouter *Router) {
	isGRoup := len(regRouter.Children) > 0

	if isGRoup {
		groupRouter := pGroupRouter.Group(regRouter.Path, regRouter.Middlewares...)

		if len(regRouter.Children) == 0 {
			return
		}

		// 递归地为每个子路由注册。
		for _, child := range regRouter.Children {
			register(groupRouter, &child)
		}
	} else {
		registerRouter(pGroupRouter, regRouter)
	}
}

func registerRouter(pGroupRouter *gin.RouterGroup, regRouter *Router) {
	var name string

	if regRouter.Name != "" {
		name = regRouter.Name
	} else {
		name = regRouter.Path
	}

	if regRouter.MODEL != nil && !regRouter.NoAutoMigrate {
		if err := cf.ORMDB.AutoMigrate(regRouter.MODEL); err != nil {
			log.Fatalf("%s: 数据库自动迁移失败: %v", name, err)
		}
	}

	// 检查路由组、路由结构体和方法是否有效，如果无效则记录日志并跳过路由注册。
	if pGroupRouter == nil || regRouter == nil || regRouter.Method == "" {
		log.Printf("%s: 没有 Method, 已跳过路由注册", name)
		return
	}

	pGroupRouter.Handle(
		regRouter.Method,
		regRouter.Path,
		append(regRouter.Middlewares, genHandler(regRouter, name))...,
	)
}

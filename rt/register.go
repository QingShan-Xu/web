package rt

import "github.com/gin-gonic/gin"

func register(pGroupRouter *gin.RouterGroup, regRouter *Router) {
	isGRoup := IsGroupRouter(regRouter)
	if isGRoup {
		groupRouter := pGroupRouter.Group(regRouter.Path)
		if len(regRouter.Middlewares) > 0 {
			groupRouter.Use(regRouter.Middlewares...)
		}
		if len(regRouter.Children) == 0 {
			return
		}
		for _, child := range regRouter.Children {
			register(groupRouter, &child)
		}
	} else {
		registerRouter(pGroupRouter, regRouter)
	}
}

func registerRouter(pGroupRouter *gin.RouterGroup, regRouter *Router) {
	if pGroupRouter == nil || regRouter == nil || regRouter.Path == "" {
		return
	}
	var name string
	if regRouter.Name != "" {
		name = regRouter.Name
	} else {
		name = regRouter.Path
	}

	router := pGroupRouter.Handle(regRouter.Method, regRouter.Path)

	if regRouter.Bind != nil {
		bindVal := getInstanceVal(regRouter.Bind)
		router.Use(bindReqMiddlewares(bindVal))
	}
}

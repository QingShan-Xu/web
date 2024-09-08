package rt

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

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
	pGroupRouter.Handle(
		regRouter.Method,
		regRouter.Path,
		reqBindMiddleware(regRouter, name),
		reqPreDBMiddleware(regRouter, name),
		func(ctx *gin.Context) {
			bind := ctx.MustGet("reqBind_")
			fmt.Printf("%+v", bind)
			ctx.JSON(200, gin.H{"data": bind})
		},
	)

	// if regRouter.Bind != nil {
	// 	bindVal := getInstanceVal(regRouter.Bind)
	// 	router.Use(bindReqMiddlewares(bindVal, name))
	// }
}

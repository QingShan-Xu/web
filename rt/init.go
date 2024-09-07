package rt

import (
	"github.com/QingShan-Xu/xjh/cf"
)

func Init(route *Router) {
	rootGroup := cf.GinGroup.Group("")
	register(rootGroup, route)
}

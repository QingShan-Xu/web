package rt

import (
	"github.com/QingShan-Xu/xjh/cf"
)

func Init(route *Router) {
	register(cf.GinGroup, route)
}

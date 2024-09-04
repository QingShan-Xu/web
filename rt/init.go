package rt

import (
	"github.com/QingShan-Xu/xjh/cf"
)

func Init(route *Route) {
	rootGroup := cf.GinGroup.Group("")
	register(rootGroup, route)
}

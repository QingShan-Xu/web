package rt

import (
	"gitee.com/be_clear/xjh/cf"
)

func Init(route *Route) {
	rootGroup := cf.GinGroup.Group("")
	register(rootGroup, route)
}

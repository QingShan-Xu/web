package rt

func RegisterRouter(rt Router) {
	if rt.isGroup() {
		rt.genGroupRouter(make([]string, 0))
	} else {
		rt.genGroupRouter(make([]string, 0))
	}
}

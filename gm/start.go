package gm

import (
	"github.com/QingShan-Xu/web/rt"
)

func Start(
	cfg Cfg,
	router rt.Router,
) {
	cfg.init()
	startRouter(cfg, router)
	startDB()
}

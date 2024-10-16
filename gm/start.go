package gm

import (
	"github.com/QingShan-Xu/web/db"
	"github.com/QingShan-Xu/web/rt"
)

func Start(
	cfg Cfg,
	router *rt.Router,
) {
	cfg.init()
	db.DB.Register()
	startRouter(cfg, router)
}

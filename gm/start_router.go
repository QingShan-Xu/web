package gm

import (
	"log"

	"github.com/QingShan-Xu/web/rt"
	"github.com/go-chi/chi/v5"
)

func routerRegister(router *rt.Router) *chi.Mux {

	// 注册
	r, err := rt.Register(router)
	if err != nil {
		log.Fatalf("%v", err)
	}

	return r
}

package gm

import (
	"log"
	"net/http"

	"github.com/QingShan-Xu/web/rt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func routerRegister(cfg Cfg, router *rt.Router) *chi.Mux {

	// 中间件
	router.Middlewares = append(
		router.Middlewares,
		[]func(http.Handler) http.Handler{
			middleware.Logger,
			middleware.CleanPath,
		}...,
	)

	// 注册
	r, err := rt.Register(router)
	if err != nil {
		log.Fatalf("%v", err)
	}

	return r
}

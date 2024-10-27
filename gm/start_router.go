package gm

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/QingShan-Xu/web/rt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/spf13/viper"
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

	// pingRouter
	ping_ := viper.Get("App.Ping").(string)
	ping, err := strconv.ParseBool(ping_)
	if err != nil {
		log.Fatalf("%s: App.Ping is not bool", strings.Join(cfg.FilePath, "/")+cfg.FileName+cfg.FileType)
	}
	if ping {
		// router.Children = append(router.Children, rt.PingRouter)
	}

	// 注册
	r, err := rt.Register(router)
	if err != nil {
		log.Fatalf("%v", err)
	}

	return r
}

package gm

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/QingShan-Xu/web/rt"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/spf13/viper"
)

func startRouter(cfg Cfg, router *rt.Router) {

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
	// 监听端口
	port, ok := viper.Get("App.Port").(string)
	if !ok {
		log.Fatalf("%s: App.Port is not string", strings.Join(cfg.FilePath, "/")+cfg.FileName+cfg.FileType)
	}
	if !strings.Contains(port, ":") {
		port = ":" + port
	}

	// 启动
	fmt.Printf("Server started at %s\n", port)
	if err := http.ListenAndServe(port, r); err != nil {
		log.Fatalf("Serve err %s", err.Error())
	}
}

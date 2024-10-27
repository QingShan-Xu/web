package gm

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/QingShan-Xu/web/db"
	"github.com/QingShan-Xu/web/rt"
	"github.com/spf13/viper"
)

func Start(
	cfg Cfg,
	router *rt.Router,
	brforeStart func(),
) {
	cfg.init()
	db.DB.Register()
	r := routerRegister(cfg, router)

	brforeStart()
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

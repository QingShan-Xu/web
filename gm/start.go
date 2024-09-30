package gm

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/QingShan-Xu/web/rt"
	"github.com/spf13/viper"
)

func Start(
	cfg Cfg,
	router rt.Router,
) {
	cfg.init()
	startRouter(cfg, router)
}

func startRouter(cfg Cfg, router rt.Router) {
	r := router.Register()
	port, ok := viper.Get("Serve.Port").(string)
	if !ok {
		log.Fatalf("%s: Serve.Port is not string", strings.Join(cfg.FilePath, "/")+cfg.FileName+cfg.FileType)
	}
	if !strings.Contains(port, ":") {
		port = ":" + port
	}

	fmt.Printf("Server started at %s\n", port)
	if err := http.ListenAndServe(port, r); err != nil {
		log.Fatalf("Serve err %s", err.Error())
	}
}

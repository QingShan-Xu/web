package gm

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/QingShan-Xu/web/db"
	"github.com/QingShan-Xu/web/rt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/docgen"
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

	initDoc(r)

	// 启动
	fmt.Printf("Server started at %s\n", port)
	if err := http.ListenAndServe(port, r); err != nil {
		log.Fatalf("Serve err %s", err.Error())
	}
}

func initDoc(r *chi.Mux) {
	relativePath := viper.GetString("Doc.RelativePath")
	if relativePath == "" {
		return
	}
	workDir, _ := os.Getwd()
	filePath := filepath.Join(workDir, relativePath)
	file, err := os.Create(filePath)
	if err != nil {
		log.Fatalf("API文档初始化失败: %v", err)
	}
	defer file.Close()
	content := docgen.MarkdownRoutesDoc(r, docgen.MarkdownOpts{
		ProjectPath: "项目接口文档",
		Intro:       "由 chi/docgo 自动生成, 请勿修改",
	})

	_, err = file.WriteString(content)
	if err != nil {
		log.Fatalf("API文档初始化失败: %v", err)
	}
}

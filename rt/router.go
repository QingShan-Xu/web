package rt

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

type (
	// 中间件
	Middlewares func(next http.Handler) http.Handler

	// 路由
	Router struct {
		Name        string        // 路由名称
		Path        string        // 路由路径
		Middlewares []Middlewares // 中间件
		Children    []Router      // 子路由
		Method      string        // 请求方法
		Handler     http.Handler  // 处理函数
		Bind        interface{}   // 参数绑定
	}
)

// 请求方法
var METHOD = struct {
	GET     string
	POST    string
	HEAD    string
	PUT     string
	PATCH   string
	DELETE  string
	OPTIONS string
	TRACE   string
	CONNECT string
}{
	GET:     "GET",
	POST:    "POST",
	HEAD:    "HEAD",
	PUT:     "PUT",
	PATCH:   "PATCH",
	DELETE:  "DELETE",
	OPTIONS: "OPTIONS",
	TRACE:   "TRACE",
	CONNECT: "CONNECT",
}

// func (curRT Router) genPath(pPath string) string {
// 	if curRT.Path == "" {
// 		return pPath
// 	}
// 	pathSlice := strings.Split(pPath, "/")
// 	pathSlice = append(pathSlice, curRT.Path)
// 	return strings.Join(pathSlice, "/")
// }

func (rt Router) isGroup() bool {
	return rt.Children != nil
}

func (curRT Router) genGroupRouter(nameSlice []string) chi.Router {
	nameSlice = append(nameSlice, curRT.Name)
	if err := curRT.checkGroupRouter(); err != nil {
		log.Fatalf("%s: %e", strings.Join(nameSlice, "."), err)
	}
	curCR := chi.NewRouter()

	curCR.Group(func(r chi.Router) {
		for _, child := range curRT.Children {
			if child.isGroup() {
				r.Mount(child.Path, child.genGroupRouter(nameSlice))
			} else {
				r.Mount(child.Path, child.genRouterItem(nameSlice))
			}
		}
	})
	return curCR
}

func (curRT Router) genRouterItem(nameSlice []string) chi.Router {
	nameSlice = append(nameSlice, curRT.Name)
	if err := curRT.checkRouter(); err != nil {
		log.Fatalf("%s: %e", strings.Join(nameSlice, "."), err)
	}
	curCR := chi.NewRouter()
	curCR.Method(curRT.Method, "/", curRT.Handler)
	return curCR
}

func (rt Router) checkGroupRouter() error {
	if rt.Method != "" {
		return fmt.Errorf("Method is not allowed on group router")
	}
	if rt.Handler != nil {
		return fmt.Errorf("Handler is not allowed on group router")
	}
	if rt.Bind != nil {
		return fmt.Errorf("Bind is not allowed on group router")
	}
	if rt.Name == "" {
		return fmt.Errorf("Name is required on group router")
	}
	return nil
}

func (rt Router) checkRouter() error {
	if rt.Name == "" {
		return fmt.Errorf("Name is required on router")
	}
	if rt.Method == "" {
		return fmt.Errorf("Method is required on router")
	}
	return nil
}

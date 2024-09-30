package rt

import (
	"fmt"
	"log"
	"net/http"
	"path"
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
		Handler     Handler       // 处理函数
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

func (curRT Router) Register() chi.Router {
	if !curRT.isGroup() {
		log.Fatalf("root router must be group router")
	}
	r := chi.NewRouter()
	curRT.genGroupRouter(make([]*Router, 0), r, "", true)
	return r
}

func (curRT *Router) genInfo(rtList []*Router) string {
	var sb strings.Builder
	pathSlice := make([]string, 0)
	nameSlice := make([]string, 0)
	for _, rt := range rtList {
		if rt.Path != "/" {
			pathSlice = append(pathSlice, rt.Path)
		}
		if rt.Name != "" {
			nameSlice = append(nameSlice, rt.Name)
		}
	}
	fullPath := path.Join(pathSlice...)
	name := strings.Join(nameSlice, ".")

	if fullPath == "" {
		sb.WriteString("/")
	} else {
		sb.WriteString(fullPath)
	}
	if curRT.Method != "" {
		sb.WriteString(" ")
		sb.WriteString(curRT.Method)
	}
	if name != "" {
		sb.WriteString(" (")
		sb.WriteString(name)
		sb.WriteString(")")
	}
	return sb.String()
}

func (curRT *Router) genRegisterInfo(rtList []*Router) string {
	if len(rtList) == 0 {
		return ""
	}
	return rtList[len(rtList)-1].genInfo([]*Router{rtList[len(rtList)-1]})
}

func getTreeSymbol(isLast bool) string {
	if isLast {
		return "└── "
	}
	return "├── "
}

func getNextPrefix(isLast bool) string {
	if isLast {
		return "    "
	}
	return "│   "
}

func (rt Router) isGroup() bool {
	return rt.Children != nil
}

func (curRT *Router) genGroupRouter(rtList []*Router, pCR chi.Router, prefix string, isLast bool) {
	rtList = append(rtList, curRT)
	if err := curRT.checkGroupRouter(); err != nil {
		log.Fatalf("\n%s\nerr: %s", curRT.genInfo(rtList), err.Error())
	}

	rtInfo := curRT.genRegisterInfo(rtList)
	if rtInfo != "" {
		fmt.Println(prefix + getTreeSymbol(isLast) + rtInfo)
	}

	newPrefix := prefix + getNextPrefix(isLast)

	pCR.Route(curRT.Path, func(curCR chi.Router) {
		for _, middlewares := range curRT.Middlewares {
			curCR.Use(middlewares)
		}
		for i, chlRT := range curRT.Children {
			if chlRT.isGroup() {
				chlRT.genGroupRouter(rtList, curCR, newPrefix, i == len(curRT.Children)-1)
			} else {
				chlRT.genRouterItem(rtList, curCR, newPrefix, i == len(curRT.Children)-1)
			}
		}
	})
}

func (curRT *Router) genRouterItem(rtList []*Router, pCR chi.Router, prefix string, isLast bool) {
	rtList = append(rtList, curRT)
	if err := curRT.checkRouter(); err != nil {
		log.Fatalf("\n%s\nerr: %s", curRT.genInfo(rtList), err.Error())
	}

	rtInfo := curRT.genRegisterInfo(rtList)
	if rtInfo != "" {
		fmt.Println(prefix + getTreeSymbol(isLast) + rtInfo)
	}
	pCR.Method(curRT.Method, curRT.Path, curRT.Handler)
}

func (rt *Router) checkGroupRouter() error {
	if rt.Method != "" {
		return fmt.Errorf("method is not allowed on group router")
	}
	if rt.Handler != nil {
		return fmt.Errorf("handler is not allowed on group router")
	}
	if rt.Bind != nil {
		return fmt.Errorf("bind is not allowed on group router")
	}
	if !strings.Contains(rt.Path, "/") {
		rt.Path = "/" + rt.Path
	}
	return nil
}

func (rt *Router) checkRouter() error {
	if rt.Method == "" {
		return fmt.Errorf("method is required on router")
	}
	if !strings.Contains(rt.Path, "/") {
		rt.Path = "/" + rt.Path
	}
	return nil
}

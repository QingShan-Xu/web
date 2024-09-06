package rt

import (
	"github.com/QingShan-Xu/xjh/bm"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type HendlerReq struct {
	C  *gin.Context
	TX *gorm.DB
	RT *Route
}
type Handler func(req *HendlerReq) *bm.Response

var TYPE = struct {
	GET_LIST    string
	GET_ONE     string
	UPDATE_ONE  string
	UPDATE_LIST string
	CREATE_ONE  string
	CREATE_LIST string
	DELETE      string
	DELETE_LIST string
}{
	GET_LIST:    "GET_LIST",
	GET_ONE:     "GET_ONE",
	UPDATE_ONE:  "UPDATE",
	UPDATE_LIST: "UPDATE_LIST",
	CREATE_ONE:  "CREATE",
	CREATE_LIST: "CREATE_LIST",
	DELETE:      "DELETE",
	DELETE_LIST: "DELETE_LIST",
}

var METHOD = struct {
	GET    string
	POST   string
	PUT    string
	DELETE string
}{
	GET:    "GET",
	POST:   "POST",
	PUT:    "PUT",
	DELETE: "DELETE",
}

type Route struct {
	// 通用 路由组
	Name        string            // 路由名称
	Path        string            // 路由路径
	Middlewares []gin.HandlerFunc // 中间件
	Children    []Route           // 子路由
	NoToken     bool

	// 子路由
	Type    string          // 简写方法 rt.Type.GET_LIST
	Method  string          // 请求方法 rt.Method.GET
	Handler gin.HandlerFunc // 处理函数
	Bind    interface{}     // 请求参数 struct{ Name string `json:"name" form:"name"` }

	// 数据库: 链式条件
	MODEL       interface{}
	SELECT      []string          // []string{"name"}
	OMIT        []string          // []string{"name"}
	DISTINCT    []string          // []string{"name"}
	WHERE       map[string]string // map[string]string{"name = ?": "Name"}
	NOT         map[string]string // map[string]string{"name = ?": "Name"}
	OR          map[string]string // map[string]string{"name = ?": "Name"}
	HAVING      map[string]string // map[string]string{"name = ?": "Name"}
	MAP_COLUMNS map[string]string // map[string]string{"name": "姓名"}
	RAW         map[string]string // map[string]string{"SELECT id, name, age FROM users WHERE name = ?": "Name"}
	ORDER       map[string]string // true升序, false降序
	LIMIT       int
	OFFSET      int
	JOINS       string
	INNER_JOINS string
	PRELOAD     []string
	TABLE       string
	GROUP       string
	CLAUSES     clause.OnConflict
}

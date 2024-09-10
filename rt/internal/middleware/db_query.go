package middleware

import (
	"log"
	"reflect"

	"github.com/QingShan-Xu/xjh/cf"
	"github.com/QingShan-Xu/xjh/rt/internal/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Query struct {
	Type  string // 类型
	Query string // 语句
	Data  interface{}
}

func ReqPreDBMiddleware(
	WHERE map[string]string,
	Bind interface{},
	MODEL interface{},
	name string,
) gin.HandlerFunc {
	var QueryList []Query

	if MODEL == nil {
		log.Fatalf("%s: MODEL 不能为空", name)
	}

	newMODEL := reflect.New(utils.GetInstanceVal(MODEL).Type()).Interface()

	bindVal := utils.GetInstanceVal(Bind)
	if bindVal.Type().Kind() != reflect.Struct {
		log.Fatalf("%s: Bind 必须为结构体", name)
	}

	bindFieldNames := make(map[string]struct{})
	t := bindVal.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		bindFieldNames[field.Name] = struct{}{}
	}

	if len(WHERE) > 0 {
		for query, data := range WHERE {
			if data == "" {
				log.Fatalf("%s: WHERE 条件值不能为空", name)
			}
			if _, exists := bindFieldNames[data]; !exists {
				log.Fatalf("%s: WHERE 引用值: %s 不在 Bind 或 InitValue 中", name, data)
			}
			QueryList = append(QueryList, Query{
				Type:  "WHERE",
				Query: query,
				Data:  data,
			})
		}
	}

	return func(ctx *gin.Context) {
		reqBind := struct2map(ctx.MustGet("reqBind_"))
		var needQuery []Query

		db := cf.ORMDB.Model(newMODEL)

		for _, query := range QueryList {
			if queryData, ok := reqBind[query.Data.(string)]; ok {
				switch qd := queryData.(type) {
				case string:
					if qd != "" {
						needQuery = append(needQuery, Query{
							Type:  query.Type,
							Query: query.Query,
							Data:  qd,
						})
					}
				case int:
					if qd != 0 {
						needQuery = append(needQuery, Query{
							Type:  query.Type,
							Query: query.Query,
							Data:  qd,
						})
					}
				}
			}
		}

		for _, query := range needQuery {
			if query.Type == "WHERE" {
				db.Where(query.Query, query.Data)
			}
		}

		ctx.Set("reqModel_", newMODEL)
		ctx.Set("reqTX_", db.Session(&gorm.Session{}))
		ctx.Next()
	}
}

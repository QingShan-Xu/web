package middleware

import (
	"fmt"
	"reflect"

	"github.com/QingShan-Xu/xjh/cf"
	"github.com/QingShan-Xu/xjh/rt/internal/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func ReqPreDBMiddleware(
	WHERE map[string]string,
	ORDER map[string]string,
	SELECT map[string]string,
	PRELOAD []string,
	JOINS []string,

	Bind interface{},
	TYPE string,
	MODEL interface{},
	name string,
) gin.HandlerFunc {

	return func(ctx *gin.Context) {
		newMODEL := reflect.New(utils.GetInstanceVal(MODEL).Type()).Interface()
		db := cf.ORMDB.Session(&gorm.Session{}).Model(newMODEL)

		reqBindMap := utils.MapFlatten(utils.Struct2map(ctx.MustGet("reqBind_")))

		if len(WHERE) > 0 {
			for query, data := range WHERE {
				if bindData, ok := reqBindMap[data]; ok {
					db.Where(query, bindData)
				}
			}
		}

		if len(ORDER) > 0 {
			for data, query := range ORDER {
				bindData, bindDataOk := reqBindMap[data]
				bindQuery, bindQueryOk := reqBindMap[query]
				if bindDataOk && bindQueryOk {
					db.Order(fmt.Sprintf("%s %s", bindData, bindQuery))
				}
			}
		}

		if len(SELECT) > 0 {
			for query := range SELECT {
				db.Select(query)
			}
		}

		if len(PRELOAD) > 0 {
			for _, query := range PRELOAD {
				db.Preload(query)
			}
		}

		if len(JOINS) > 0 {
			for _, query := range JOINS {
				db.Joins(query)
			}
		}

		ctx.Set("reqModel_", newMODEL)
		ctx.Set("reqTX_", db.Session(&gorm.Session{}))
		ctx.Next()
	}
}

package middleware

import (
	"fmt"
	"log"
	"reflect"

	"github.com/QingShan-Xu/xjh/cf"
	"github.com/QingShan-Xu/xjh/rt/internal/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func ReqPreDBMiddleware(
	WHERE map[string]string,
	ORDER map[string]string,

	Bind interface{},
	TYPE string,
	MODEL interface{},
	name string,
) gin.HandlerFunc {

	if MODEL == nil {
		log.Fatalf("%s: MODEL 不能为空", name)
	}

	newMODEL := reflect.New(utils.GetInstanceVal(MODEL).Type()).Interface()

	bindVal := utils.GetInstanceVal(Bind)
	if bindVal.Type().Kind() != reflect.Struct {
		log.Fatalf("%s: Bind 必须为结构体", name)
	}
	fmt.Printf("%+v", Bind)
	bindFieldNames := utils.MapFlatten(utils.Struct2map(Bind, true))

	if len(WHERE) > 0 {
		for query, data := range WHERE {
			if data == "" || query == "" {
				log.Fatalf("%s , %s: WHERE 条件值或语句不能为空", query, name)
			}
			if _, exists := bindFieldNames[data]; !exists {
				log.Fatalf("%s: WHERE 引用值: %s 不在 Bind 中", name, data)
			}
		}
	}

	if len(ORDER) > 0 {
		for data, query := range ORDER {
			if _, exists := bindFieldNames[data]; !exists {
				log.Fatalf("%s: ORDER 引用值: %s 不在 Bind 中", name, data)
			}
			if _, exists := bindFieldNames[query]; !exists {
				log.Fatalf("%s: ORDER 引用值: %s 不在 Bind 中", name, data)
			}
		}
	}

	if TYPE == "GET_LIST" {
		dataPage, existsPage := bindFieldNames["Pagination.PageSize"]
		if !existsPage {
			log.Fatalf("%s: GET_LIST 引用值: %s 不在 Bind 中", name, "Pagination.PageSize")
		}
		if _, ok := dataPage.(int); !ok {
			log.Fatalf("%s: GET_LIST 引用值: %s 不是 int 类型", name, "Pagination.PageSize")
		}

		dataCur, existsCur := bindFieldNames["Pagination.Current"]
		if !existsCur {
			log.Fatalf("%s: GET_LIST 引用值: %s 不在 Bind 中", name, "Pagination.Current")
		}
		if _, ok := dataCur.(int); !ok {
			log.Fatalf("%s: GET_LIST 引用值: %s 不是 int 类型", name, "Pagination.Current")
		}
	}

	/* 	if TYPE == "CREATE_ONE" || TYPE == "UPDATE_ONE" || TYPE == "DELETE_ONE" {
		log.Fatal("还没做")
	} */

	return func(ctx *gin.Context) {

		reqBindMap := utils.MapFlatten(utils.Struct2map(ctx.MustGet("reqBind_"), false))

		db := cf.ORMDB.Model(newMODEL)

		if TYPE == "GET_LIST" {
			pageSize := reqBindMap["Pagination.PageSize"].(int)
			current := reqBindMap["Pagination.Current"].(int)
			db.Limit(pageSize).Offset((current - 1) * pageSize)
		}

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

		ctx.Set("reqModel_", newMODEL)
		ctx.Set("reqTX_", db.Session(&gorm.Session{}))
		ctx.Next()
	}
}

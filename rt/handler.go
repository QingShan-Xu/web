package rt

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"

	"github.com/QingShan-Xu/xjh/bm"
	"github.com/QingShan-Xu/xjh/cf"
	"github.com/QingShan-Xu/xjh/rt/internal/class"
	"github.com/QingShan-Xu/xjh/rt/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

func genHandler(
	router *Router,
	name string,
) gin.HandlerFunc {

	dynamicBindStruct := class.DynamicStruct{Value: reflect.ValueOf(router.Bind)}

	if err := check(router, dynamicBindStruct); err != nil {
		log.Fatalf("%s: %v", name, err)
	}

	return handler(router)
}

func handler(router *Router) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 绑定请求值
		var res *bm.Res
		var bindData interface{} = struct{}{}
		if router.Bind != nil {
			bindVal := reflect.ValueOf(router.Bind)
			bindDataVal := reflect.New(reflect.TypeOf(router.Bind))

			// 赋初值
			for i := 0; i < bindVal.NumField(); i++ {
				fieldValue := bindVal.Field(i)
				if fieldValue.IsValid() && bindDataVal.Field(i).CanSet() {
					bindDataVal.Field(i).Set(fieldValue)
				}
			}

			bindData = bindDataVal.Interface()

			// uri
			if strings.Contains(ctx.FullPath(), ":") {
				ctx.ShouldBindUri(bindData)
			}

			// json, form, query
			if err := ctx.ShouldBind(bindData); err != nil {
				// binding 错误处理
				new(bm.Res).FailFront(err2Str(err)).SendAbort(ctx)
				return
			}
		}

		dynamicBindStruct := class.DynamicStruct{Value: reflect.ValueOf(bindData)}

		modelVal := reflect.New(reflect.TypeOf(router.MODEL))
		if modelVal.Kind() == reflect.Pointer {
			modelVal = modelVal.Elem()
		}

		db := cf.ORMDB
		model := modelVal.Interface()

		if model != nil {
			db.Model(model)
		}

		if len(router.WHERE) > 0 {
			for query, data := range router.WHERE {
				if bindData, err := dynamicBindStruct.GetField(data); err != nil {
					res.FailBackend(err).SendAbort(ctx)
					return
				} else {
					db.Where(query, bindData)
				}
			}
		}

		if len(router.ORDER) > 0 {
			for data, query := range router.ORDER {

				bindData, bindDataErr := dynamicBindStruct.GetField(data)
				if bindDataErr != nil {
					res.FailBackend(bindDataErr).SendAbort(ctx)
					return
				}
				bindQuery, bindQueryErr := dynamicBindStruct.GetField(query)
				if bindQueryErr != nil {
					res.FailBackend(bindQueryErr).SendAbort(ctx)
					return
				}
				db.Order(fmt.Sprintf("%s %s", bindData, bindQuery))
			}
		}

		if len(router.SELECT) > 0 {
			for query := range router.SELECT {
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

		ctx.Set("reqModel_", model)
		ctx.Set("reqTX_", db.Session(&gorm.Session{}))
		ctx.Next()

	}
}

func check(
	router *Router,
	dynamicBindStruct class.DynamicStruct,
) error {
	if router.Bind != nil && reflect.TypeOf(router.Bind).Kind() != reflect.Struct {
		return fmt.Errorf("router.Bind 必须是 [struct]")
	}

	if router.MODEL == nil &&
		(len(router.WHERE) != 0 ||
			len(router.ORDER) != 0 ||
			len(router.SELECT) != 0 ||
			len(router.PRELOAD) != 0 ||
			len(router.JOINS) != 0) {
		return fmt.Errorf("当使用数据库字段时 [MODEL] 不能为空")
	}

	if len(router.WHERE) > 0 {
		for query, data := range router.WHERE {
			if data == "" || query == "" {
				return fmt.Errorf("WHERE 条件值或语句 [%s] 不能为空", query)
			}
			if _, err := dynamicBindStruct.GetField(data); err != nil {
				return fmt.Errorf("WHERE 条件值或语句 %e", err)
			}
		}
	}

	if len(router.ORDER) > 0 {
		for data, query := range router.ORDER {
			if _, err := dynamicBindStruct.GetField(data); err != nil {
				return fmt.Errorf("ORDER 语句 %e", err)
			}
			if _, err := dynamicBindStruct.GetField(query); err != nil {
				return fmt.Errorf("ORDER 条件值 %e", err)
			}
		}
	}

	if len(router.SELECT) > 0 {
		for query, data := range router.SELECT {
			if data == "" || query == "" {
				log.Fatalf("SELECT 条件值或语句 [%s] 不能为空", query)
			}
			if _, err := dynamicBindStruct.GetField(data); err != nil {
				return fmt.Errorf("SELECT 条件值或语句 %e", err)
			}
		}
	}

	if len(router.PRELOAD) > 0 {
		for _, data := range router.PRELOAD {
			if data == "" {
				log.Fatalf("PRELOAD 条件值不能为空")
			}
		}
	}

	if len(router.JOINS) > 0 {
		for _, data := range router.JOINS {
			if data == "" {
				log.Fatalf("JOINS 条件值不能为空")
			}
		}
	}

	if router.Type == TYPE.GET_LIST {
		if _, err := dynamicBindStruct.GetField("PageSize"); err != nil {
			return fmt.Errorf("GET_LIST PageSize %e", err)
		}

		if _, err := dynamicBindStruct.GetField("Current"); err != nil {
			return fmt.Errorf("GET_LIST Current %e", err)
		}
	}

	return nil
}

var type2Front = map[string]string{
	"int64":   "数字",
	"int32":   "数字",
	"int":     "数字",
	"float64": "数字",
	"float32": "数字",
	"string":  "字符串",
	"bool":    "布尔",
}

// 转换错误为字符串
func err2Str(err error) string {
	if err != nil {
		switch e := err.(type) {
		case *strconv.NumError:
			if e.Func == "ParseBool" {
				return fmt.Sprintf(" %s--该值应为[布尔]类型", e.Num)
			}
			if e.Func == "ParseInt" {
				return fmt.Sprintf(" %s--该值应为[数字]类型", e.Num)
			}
		case *json.UnmarshalTypeError:
			a := e.Type.Name()
			eTpe, ok := type2Front[a]
			if !ok {
				eTpe = a
			}
			return fmt.Sprintf("%s: 应为 %s 类型", e.Field, eTpe)
		case validator.ValidationErrors:
			var errStrSlice []string
			for _, validatoE := range e.Translate(cf.Trans) {
				errStrSlice = append(errStrSlice, utils.ToSnakeCase(validatoE))
			}
			return strings.Join(errStrSlice, ",")
		default:
			err := e.Error()
			if err == "EOF" {
				return "不得为空"
			} else {
				return err
			}
		}
	}
	return ""
}

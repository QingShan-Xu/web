package rt

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
	"unicode"

	"github.com/QingShan-Xu/xjh/bm"
	"github.com/QingShan-Xu/xjh/cf"
	"github.com/QingShan-Xu/xjh/rt/internal/class"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

func genHandler(
	router *Router,
	name string,
) gin.HandlerFunc {

	if err := check(router); err != nil {
		log.Fatalf("%s: %v", name, err)
	}

	return handler(router)
}

func handler(router *Router) gin.HandlerFunc {
	if router.OriginalHandler != nil {
		return router.OriginalHandler
	}

	return func(ctx *gin.Context) {
		// 绑定请求值
		res := &bm.Res{}
		var bindData interface{} = struct{}{}
		if router.Bind != nil {
			bindVal := reflect.ValueOf(router.Bind)
			bindDataValPtr := reflect.New(reflect.TypeOf(router.Bind))
			bindDataVal := bindDataValPtr.Elem()

			// 赋初值
			for i := 0; i < bindVal.NumField(); i++ {
				fieldValue := bindVal.Field(i)
				if fieldValue.IsValid() && bindDataVal.Field(i).CanSet() {
					bindDataVal.Field(i).Set(fieldValue)
				}
			}

			bindData = bindDataValPtr.Interface()
			// uri
			if strings.Contains(ctx.FullPath(), ":") {
				ctx.ShouldBindUri(bindData)
			}

			// json, form, query``
			if err := ctx.ShouldBind(bindData); err != nil {
				// binding 错误处理
				new(bm.Res).FailFront(err2Str(err)).SendAbort(ctx)
				return
			}
		}
		fmt.Printf("%+v", bindData)
		dynamicBindStruct := class.DynamicStruct{Value: reflect.ValueOf(bindData)}

		var modelVal reflect.Value
		var model interface{}
		if router.MODEL != nil {
			_modelTpe := reflect.TypeOf(router.MODEL)
			if _modelTpe.Kind() == reflect.Pointer {
				_modelTpe = _modelTpe.Elem()
			}
			modelVal = reflect.New(_modelTpe)
			model = modelVal.Interface()
		}

		db := cf.ORMDB.Session(&gorm.Session{})

		if model != nil {
			db = db.Model(model)
		}

		if len(router.WHERE) > 0 {
			for query, data := range router.WHERE {
				if bindData, err := dynamicBindStruct.GetField(data); err != nil {
					res.FailBackend(err).SendAbort(ctx)
					return
				} else if bindData == nil {
					continue
				} else {
					db = db.Where(query, bindData)
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
				if bindData == nil || bindQuery == nil {
					continue
				}
				db = db.Order(fmt.Sprintf("%s %s", bindData, bindQuery))
			}
		}

		if len(router.SELECT) > 0 {
			for query := range router.SELECT {
				db = db.Select(query)
			}
		}

		if len(router.PRELOAD) > 0 {
			for _, query := range router.PRELOAD {
				db = db.Preload(query)
			}
		}

		if len(router.JOINS) > 0 {
			for _, query := range router.JOINS {
				db = db.Joins(query)
			}
		}

		if router.Type == "" && router.Handler != nil {
			res := router.Handler(ctx, db.Session(&gorm.Session{}), bindData)
			if ctx.Writer.Written() {
				return
			}
			res.Send(ctx)
			return
		}

		if router.Type == TYPE.GET_ONE {
			result := db.First(model)

			if result.RowsAffected == 0 {
				res.FailFront("数据不存在").SendAbort(ctx)
				return
			}

			if result.Error != nil {
				res.FailBackend("查询失败").SendAbort(ctx)
				return
			}

			res.SucJson(model).SendAbort(ctx)
			return
		}

		if router.Type == TYPE.GET_LIST {
			modelListVal := reflect.SliceOf(modelVal.Elem().Type())
			modelList := reflect.New(modelListVal).Interface()
			var total int64

			if err := db.Count(&total).Error; err != nil {
				res.FailBackend("获取失败").SendAbort(ctx)
				return
			}

			pageSize, pageSizeErr := dynamicBindStruct.GetField("PageSize")
			if pageSizeErr != nil {
				res.FailBackend("处理分页失败").SendAbort(ctx)
				return
			}

			current, currentErr := dynamicBindStruct.GetField("Current")
			if currentErr != nil {
				res.FailBackend("处理分页失败").SendAbort(ctx)
				return
			}

			if err := db.Limit(pageSize.(int)).Offset((current.(int) - 1) * pageSize.(int)).Find(modelList).Error; err != nil {
				res.FailBackend(err).SendAbort(ctx)
				return
			}

			res.SucJson(bm.ResList{
				Data:     modelList,
				Total:    total,
				PageSize: pageSize.(int),
				Current:  current.(int),
			}).SendAbort(ctx)
			return
		}

		if router.Type == TYPE.CREATE_ONE {
			result := db.Create(bindData)
			if result.Error != nil {
				res.FailBackend(result.Error).SendAbort(ctx)
				return
			}

			res.SucJson(bindData).SendAbort(ctx)
			return
		}

		if router.Type == TYPE.UPDATE_ONE {
			result := db.First(model)

			if result.RowsAffected == 0 {
				res.FailFront("数据不存在").SendAbort(ctx)
				return
			}

			if result.Error != nil {
				res.FailBackend("查询失败").SendAbort(ctx)
				return
			}

			if ok := reflect.DeepEqual(router.Bind, router.MODEL); ok {
				result = db.Updates(bindData)
			} else {
				bind := make(map[string]interface{}, 0)
				for k, v := range router.SELECT {
					if data, err := dynamicBindStruct.GetField(v); err != nil {
						res.FailBackend(fmt.Errorf("请求值 %s 缺失", v))
					} else {
						bind[k] = data
					}
				}
				result = db.Updates(bind)
			}

			if result.Error != nil {
				res.FailBackend("更新失败").SendAbort(ctx)
				return
			}

			res.SucJson(nil).SendAbort(ctx)
			return
		}

		if router.Type == TYPE.DELETE_ONE {
			result := db.First(model)

			if result.RowsAffected == 0 {
				res.FailFront("数据不存在").SendAbort(ctx)
				return
			}

			if result.Error != nil {
				res.FailBackend("查询失败").SendAbort(ctx)
				return
			}

			if err := db.Delete(model).Error; err != nil {
				res.FailBackend("删除失败").SendAbort(ctx)
				return
			}

			res.SucJson(true).SendAbort(ctx)
			return
		}

	}
}

func check(router *Router) error {
	if router.OriginalHandler != nil && router.Handler != nil {
		return fmt.Errorf("router.OriginalHandler 和 router.Handler 不能同时存在")
	}

	if router.OriginalHandler != nil && router.Handler == nil {
		return nil
	}

	if router.Type != "" && router.Handler != nil {
		return fmt.Errorf("router.Type 不能和 router.Handler 同时存在")
	}

	if router.Bind != nil && reflect.TypeOf(router.Bind).Kind() != reflect.Struct {
		return fmt.Errorf("router.Bind 必须是 [struct]")
	}

	if router.MODEL == nil &&
		(router.Type != "" ||
			len(router.WHERE) != 0 ||
			len(router.ORDER) != 0 ||
			len(router.SELECT) != 0 ||
			len(router.PRELOAD) != 0 ||
			len(router.JOINS) != 0) {
		return fmt.Errorf("当使用 router.Type 或 数据库字段时 [MODEL] 不能为空")
	}

	dynamicBindStruct := class.DynamicStruct{Value: reflect.ValueOf(router.Bind)}

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

	if router.Type != "" && router.Type != TYPE.CREATE_ONE && router.Type != TYPE.GET_ONE && router.Type != TYPE.GET_LIST && router.Type != TYPE.UPDATE_ONE && router.Type != TYPE.DELETE_ONE {
		return fmt.Errorf("不识别的 Type [%s]", router.Type)
	}

	if router.Type == TYPE.GET_LIST {
		if _, err := dynamicBindStruct.GetField("PageSize"); err != nil {
			return fmt.Errorf("GET_LIST PageSize %e", err)
		}

		if _, err := dynamicBindStruct.GetField("Current"); err != nil {
			return fmt.Errorf("GET_LIST Current %e", err)
		}
	}

	if router.Type == TYPE.CREATE_ONE {
		if ok := reflect.DeepEqual(router.Bind, router.MODEL); !ok {
			return fmt.Errorf("CREATE_ONE BIND 绑定对象与 MODEL 不一致")
		}
	}

	if router.Type == TYPE.UPDATE_ONE {
		if len(router.WHERE) == 0 {
			return fmt.Errorf("UPDATE_ONE 时 WHERE 不能为空")
		}
		if ok := reflect.DeepEqual(router.Bind, router.MODEL); !ok && len(router.SELECT) == 0 {
			return fmt.Errorf("UPDATE_ONE BIND 绑定对象与 MODEL 不一致时 SELECT 不得为空")
		}
	}

	if router.Type == TYPE.DELETE_ONE && len(router.WHERE) == 0 {
		return fmt.Errorf("DELETE_ONE 时 WHERE 不能为空")
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
				errStrSlice = append(errStrSlice, toSnakeCase(validatoE))
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

func toSnakeCase(str string) string {
	var result []rune
	for i, r := range str {
		if unicode.IsUpper(r) {
			// 如果不是字符串开始并且前一个字符不是下划线
			if i > 0 && !(unicode.IsUpper(rune(str[i-1]))) {
				result = append(result, '_') // 添加下划线
			}
			result = append(result, unicode.ToLower(r)) // 将大写字母转换为小写并添加到结果中
		} else {
			result = append(result, r) // 如果是小写字母或其他字符，直接添加到结果中
		}
	}
	return string(result) // 将 rune 数组转换为字符串并返回
}

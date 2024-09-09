package rt

import (
	"log"
	"reflect"

	"github.com/QingShan-Xu/xjh/bm"
	"github.com/QingShan-Xu/xjh/cf"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// // 绑定数据中间件
// var type2Front = map[string]string{
// 	"int64":   "数字",
// 	"int32":   "数字",
// 	"int":     "数字",
// 	"float64": "数字",
// 	"float32": "数字",
// 	"string":  "字符串",
// 	"bool":    "布尔",
// }

// type fieldTags struct {
// 	form  []reflect.StructField
// 	json  []reflect.StructField
// 	param []reflect.StructField
// 	uri   []reflect.StructField
// }

// func reqBindMiddleware(Bind interface{}, name string) gin.HandlerFunc {
// 	bindVal := getInstanceVal(Bind)
// 	bindTpe := bindVal.Type()

// 	if bindTpe == nil {
// 		log.Fatalf("%s: Bind 类型不能为 nil", name)
// 	}

// 	if bindVal.Kind() != reflect.Struct {
// 		log.Fatalf("%s: Bind 必须是 struct 类型", name)
// 	}

// 	var fieldsTags fieldTags

// 	for i := 0; i < bindTpe.NumField(); i++ {
// 		fieldTpe := bindTpe.Field(i)

// 		bindingTag := fieldTpe.Tag.Get("binding")

// 		if formTag := fieldTpe.Tag.Get("form"); formTag != "" {
// 			newField := fieldTpe
// 			if bindingTag != "" {
// 				newField.Tag = reflect.StructTag(fmt.Sprintf(
// 					`binding:"%s" form:"%s"`, bindingTag, formTag))
// 			} else {
// 				newField.Tag = reflect.StructTag(fmt.Sprintf(
// 					`form:"%s"`, formTag))
// 			}
// 			fieldsTags.form = append(fieldsTags.form, newField)
// 		}
// 		if jsonTag := fieldTpe.Tag.Get("json"); jsonTag != "" {
// 			newField := fieldTpe
// 			if bindingTag != "" {
// 				newField.Tag = reflect.StructTag(fmt.Sprintf(
// 					`binding:"%s" json:"%s"`, bindingTag, jsonTag))
// 			} else {
// 				newField.Tag = reflect.StructTag(fmt.Sprintf(
// 					`json:"%s"`,
// 					jsonTag,
// 				))
// 			}
// 			fieldsTags.json = append(fieldsTags.json, newField)
// 		}
// 		if paramTag := fieldTpe.Tag.Get("param"); paramTag != "" {
// 			newField := fieldTpe
// 			if bindingTag != "" {
// 				newField.Tag = reflect.StructTag(fmt.Sprintf(
// 					`binding:"%s" form:"%s"`, bindingTag, paramTag))
// 			} else {
// 				newField.Tag = reflect.StructTag(fmt.Sprintf(
// 					`form:"%s"`,
// 					paramTag,
// 				))
// 			}
// 			fieldsTags.param = append(fieldsTags.param, newField)
// 		}
// 		if uriTag := fieldTpe.Tag.Get("uri"); uriTag != "" {
// 			newField := fieldTpe
// 			if bindingTag != "" {
// 				newField.Tag = reflect.StructTag(fmt.Sprintf(
// 					`binding:"%s" uri:"%s"`, bindingTag, uriTag))
// 			} else {
// 				newField.Tag = reflect.StructTag(fmt.Sprintf(
// 					`uri:"%s"`,
// 					uriTag,
// 				))
// 			}
// 			fieldsTags.uri = append(fieldsTags.uri, newField)
// 		}
// 	}

// 	if len(fieldsTags.form) > 0 && len(fieldsTags.json) > 0 {
// 		log.Fatalf("%s: Bind 不能同时使用 form 和 json 标签", name)
// 	}

// 	return func(c *gin.Context) {
// 		var formBind interface{}
// 		var jsonBind interface{}
// 		var paramBind interface{}
// 		var uriBind interface{}

// 		if len(fieldsTags.json) > 0 {
// 			jsonBind = reflect.New(reflect.StructOf(fieldsTags.json)).Interface()
// 			err := c.ShouldBindBodyWith(jsonBind, binding.JSON)
// 			errStr := err2Str(err)
// 			if errStr != "" {
// 				new(bm.Res).FailFront(fmt.Sprintf("[JSON] %s", errStr)).Send(c)
// 				c.Abort()
// 				return
// 			}
// 		}

// 		if len(fieldsTags.form) > 0 {
// 			formBind = reflect.New(reflect.StructOf(fieldsTags.form)).Interface()
// 			err := c.ShouldBindWith(formBind, binding.Form)
// 			errStr := err2Str(err)
// 			if errStr != "" {
// 				new(bm.Res).FailFront(fmt.Sprintf("[FORM] %s", errStr)).Send(c)
// 				c.Abort()
// 				return
// 			}
// 		}

// 		if len(fieldsTags.param) > 0 {
// 			paramBind = reflect.New(reflect.StructOf(fieldsTags.param)).Interface()
// 			err := c.ShouldBindWith(paramBind, binding.Query)
// 			errStr := err2Str(err)
// 			if errStr != "" {
// 				new(bm.Res).FailFront(fmt.Sprintf("[PARAM] %s", errStr)).Send(c)
// 				c.Abort()
// 				return
// 			}
// 		}

// 		if len(fieldsTags.uri) > 0 {
// 			uriBind = reflect.New(reflect.StructOf(fieldsTags.uri)).Interface()
// 			err := c.ShouldBindUri(uriBind)
// 			errStr := err2Str(err)
// 			if errStr != "" {
// 				new(bm.Res).FailFront(fmt.Sprintf("[URI] %s", errStr)).Send(c)
// 				c.Abort()
// 				return
// 			}
// 		}

// 		bindDataVal := reflect.New(bindVal.Type())

// 		formData := struct2map(formBind)
// 		jsonData := struct2map(jsonBind)
// 		paramData := struct2map(paramBind)
// 		uriData := struct2map(uriBind)

// 		for i := 0; i < bindTpe.NumField(); i++ {
// 			fieldName := bindTpe.Field(i).Name
// 			bindDataField := bindDataVal.Elem().FieldByName(fieldName)

// 			if !bindDataField.CanSet() {
// 				continue
// 			}

// 			if data, ok := formData[fieldName]; ok {
// 				bindDataField.Set(reflect.ValueOf(data))
// 			}
// 			if data, ok := jsonData[fieldName]; ok {
// 				bindDataField.Set(reflect.ValueOf(data))
// 			}
// 			if data, ok := paramData[fieldName]; ok {
// 				bindDataField.Set(reflect.ValueOf(data))
// 			}
// 			if data, ok := uriData[fieldName]; ok {
// 				bindDataField.Set(reflect.ValueOf(data))
// 			}
// 		}

// 		bindData := bindDataVal.Interface()

// 		c.Set("reqBind_", bindData)
// 		c.Next()
// 	}
// }

// func err2Str(err error) string {
// 	if err != nil {
// 		switch e := err.(type) {
// 		case *strconv.NumError:
// 			if e.Func == "ParseBool" {
// 				return fmt.Sprintf(" %s--该值应为[布尔]类型", e.Num)
// 			}
// 			if e.Func == "ParseInt" {
// 				return fmt.Sprintf(" %s--该值应为[数字]类型", e.Num)
// 			}
// 		case *json.UnmarshalTypeError:
// 			a := e.Type.Name()
// 			eTpe, ok := type2Front[a]
// 			if !ok {
// 				eTpe = a
// 			}
// 			return fmt.Sprintf("%s: 应为 %s 类型", e.Field, eTpe)
// 		case validator.ValidationErrors:
// 			var errStrSlice []string
// 			for _, validatoE := range e.Translate(cf.Trans) {
// 				errStrSlice = append(errStrSlice, toSnakeCase(validatoE))
// 			}
// 			return strings.Join(errStrSlice, ",")
// 		default:
// 			err := e.Error()
// 			if err == "EOF" {
// 				return "不得为空"
// 			} else {
// 				return err
// 			}
// 		}
// 	}
// 	return ""
// }

// 操作中间件
func reqPreDBMiddleware(regRouter *Router, name string) gin.HandlerFunc {
	var QueryList []bm.Query

	if regRouter.MODEL == nil {
		log.Fatalf("%s: MODEL 不能为空", name)
	}

	bindVal := getInstanceVal(regRouter.Bind)
	if bindVal.Type().Kind() != reflect.Struct {
		log.Fatalf("%s: Bind 必须为结构体", name)
	}

	bindFileNames := make([]string, 0)
	for i := 0; i < bindVal.Type().NumField(); i++ {
		field := bindVal.Type().Field(i)
		bindFileNames = append(bindFileNames, field.Name)
	}

	if regRouter.WHERE != nil {
		for query, data := range regRouter.WHERE {
			if data == "" {
				log.Fatalf("%s: WHERE 条件值不能为空", name)
			}
			if !isIncludes(bindFileNames, data) {
				log.Fatalf("%s: WHERE 引用值: %s 不在 Bind 或 InitValue 中", name, data)
			}
			currentQuery := bm.Query{
				Type:  "WHERE",
				Query: query,
				Data:  data,
			}
			QueryList = append(QueryList, currentQuery)
		}
	}

	return func(ctx *gin.Context) {
		reqBind := struct2map(ctx.MustGet("reqBind_"))
		needQuery := make([]bm.Query, 0)

		db := cf.ORMDB

		db = db.Model(regRouter.MODEL)
		for _, query := range QueryList {
			queryData, ok := reqBind[query.Data.(string)]
			if !ok {
				continue
			}
			switch qd := queryData.(type) {
			case string:
				if qd != "" {
					newQuery := query
					newQuery.Data = qd
					needQuery = append(needQuery, newQuery)
				}
			case int:
				if qd != 0 {
					newQuery := query
					newQuery.Data = qd
					needQuery = append(needQuery, newQuery)
				}
			}
		}

		for _, query := range needQuery {
			if query.Type == "WHERE" {
				db.Where(query.Query, query.Data)
			}
		}

		ctx.Set("reqTX_", db.Session(&gorm.Session{}))

		ctx.Next()
	}
}

func reqFinisherMiddleware(regRouter *Router, name string) gin.HandlerFunc {
	if regRouter.Finisher == "" {
		return func(ctx *gin.Context) { ctx.Next() }
	}

	if regRouter.Finisher != "" && regRouter.Handler != nil {
		log.Fatalf("%s: 不能同时使用 Handler 和 Finisher", name)
	}

	if regRouter.Finisher == Finisher.First {
		return func(ctx *gin.Context) {
			tx := ctx.MustGet("reqTX_").(*gorm.DB)
			data := reflect.New(getInstanceVal(regRouter.MODEL).Type()).Interface()

			result := tx.Find(&data)
			if result.Error != nil {
				new(bm.Res).FailBackend(result.Error).Send(ctx)
				ctx.Abort()
				return
			}

			if result.RowsAffected == 0 {
				new(bm.Res).FailFront("数据不存在").Send(ctx)
				ctx.Abort()
				return
			}

			new(bm.Res).SucJson(data).Send(ctx)
			ctx.Abort()
		}
	}

	return func(ctx *gin.Context) {}
}

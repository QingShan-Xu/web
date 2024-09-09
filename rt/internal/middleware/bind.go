package middleware

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"

	"github.com/QingShan-Xu/xjh/bm"
	"github.com/QingShan-Xu/xjh/cf"
	"github.com/QingShan-Xu/xjh/rt/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type FieldTags struct {
	form  []reflect.StructField
	json  []reflect.StructField
	param []reflect.StructField
	uri   []reflect.StructField
}

func ReqBindMiddleware(Bind interface{}, name string) gin.HandlerFunc {
	bindVal := utils.GetInstanceVal(Bind)
	bindTpe := bindVal.Type()

	if bindTpe == nil {
		log.Fatalf("%s: Bind 类型不能为 nil", name)
	}

	if bindVal.Kind() != reflect.Struct {
		log.Fatalf("%s: Bind 必须是 struct 类型", name)
	}

	var fieldsTags FieldTags

	// 遍历结构体的字段，分类处理不同类型的标签
	for i := 0; i < bindTpe.NumField(); i++ {
		fieldTpe := bindTpe.Field(i)
		processFieldTags(fieldTpe, &fieldsTags)
	}

	// 检查是否有冲突的标签
	if len(fieldsTags.form) > 0 && len(fieldsTags.json) > 0 {
		log.Fatalf("%s: Bind 不能同时使用 form 和 json 标签", name)
	}

	return func(c *gin.Context) {
		bindDataVal := reflect.New(bindVal.Type())

		// 统一绑定不同类型的数据
		if err := bindData(c, &fieldsTags, bindDataVal); err != nil {
			new(bm.Res).FailFront(err.Error()).Send(c)
			c.Abort()
			return
		}

		// 将绑定的数据存储到上下文中
		c.Set("reqBind_", bindDataVal.Interface())
		c.Next()
	}
}

// 处理字段标签，分类存储到 FieldTags 结构体中
func processFieldTags(fieldTpe reflect.StructField, fieldsTags *FieldTags) {
	bindingTag := fieldTpe.Tag.Get("binding")

	if tag := fieldTpe.Tag.Get("form"); tag != "" {
		fieldsTags.form = append(fieldsTags.form, updateFieldTag(fieldTpe, "form", tag, bindingTag))
	}
	if tag := fieldTpe.Tag.Get("json"); tag != "" {
		fieldsTags.json = append(fieldsTags.json, updateFieldTag(fieldTpe, "json", tag, bindingTag))
	}
	if tag := fieldTpe.Tag.Get("param"); tag != "" {
		fieldsTags.param = append(fieldsTags.param, updateFieldTag(fieldTpe, "param", tag, bindingTag))
	}
	if tag := fieldTpe.Tag.Get("uri"); tag != "" {
		fieldsTags.uri = append(fieldsTags.uri, updateFieldTag(fieldTpe, "uri", tag, bindingTag))
	}
}

// 更新结构体字段的标签
func updateFieldTag(fieldTpe reflect.StructField, tagKey, tagValue, bindingTag string) reflect.StructField {
	newField := fieldTpe
	if bindingTag != "" {
		newField.Tag = reflect.StructTag(fmt.Sprintf(`binding:"%s" %s:"%s"`, bindingTag, tagKey, tagValue))
	} else {
		newField.Tag = reflect.StructTag(fmt.Sprintf(`%s:"%s"`, tagKey, tagValue))
	}
	return newField
}

// 统一绑定不同类型的数据
func bindData(c *gin.Context, fieldsTags *FieldTags, bindDataVal reflect.Value) error {
	// 定义并初始化绑定数据
	var binds = []struct {
		data  interface{}
		kind  string
		bType binding.Binding
		tags  []reflect.StructField
	}{
		{kind: "Form", bType: binding.Form, tags: fieldsTags.form},
		{kind: "Param", bType: binding.Query, tags: fieldsTags.param},
	}

	// 处理 Json 绑定
	if len(fieldsTags.json) > 0 {
		jsonBind := reflect.New(reflect.StructOf(fieldsTags.json)).Interface()
		if err := c.ShouldBindBodyWith(jsonBind, binding.JSON); err != nil {
			return fmt.Errorf("[Json] %s", err2Str(err))
		}
		mergeBindData(jsonBind, bindDataVal)
	}

	// 处理 URI 绑定
	if len(fieldsTags.uri) > 0 {
		uriBind := reflect.New(reflect.StructOf(fieldsTags.uri)).Interface()
		if err := c.ShouldBindUri(uriBind); err != nil {
			return fmt.Errorf("[URI] %s", err2Str(err))
		}
		mergeBindData(uriBind, bindDataVal)
	}

	// 处理 JSON、FORM、PARAM 绑定
	for _, b := range binds {
		if len(b.tags) > 0 {
			b.data = reflect.New(reflect.StructOf(b.tags)).Interface()
			if err := c.ShouldBindWith(b.data, b.bType); err != nil {
				return fmt.Errorf("[%s] %s", b.kind, err2Str(err))
			}
			mergeBindData(b.data, bindDataVal)
		}
	}

	return nil
}

// 将绑定的数据合并到最终的结构体中
func mergeBindData(srcData interface{}, dstDataVal reflect.Value) {
	srcMap := struct2map(srcData)
	for key, val := range srcMap {
		dstField := dstDataVal.Elem().FieldByName(key)
		if dstField.CanSet() {
			dstField.Set(reflect.ValueOf(val))
		}
	}
}

// 将结构体数据转为 map
func struct2map(s interface{}) map[string]interface{} {
	data := make(map[string]interface{})
	val := reflect.ValueOf(s).Elem()
	for i := 0; i < val.NumField(); i++ {
		data[val.Type().Field(i).Name] = val.Field(i).Interface()
	}
	return data
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

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
	"github.com/go-playground/validator/v10"
)

func ReqBindMiddleware(Bind interface{}, name string) gin.HandlerFunc {
	bindVal := utils.GetInstanceVal(Bind)

	if bindVal.Kind() != reflect.Struct {
		log.Fatalf("%s: Bind 必须是 struct 类型", name)
	}

	return func(c *gin.Context) {
		bindData := reflect.New(bindVal.Type()).Interface()
		bindDataVal := utils.GetInstanceVal(bindData)
		for i := 0; i < bindVal.NumField(); i++ {
			fieldValue := bindVal.Field(i)
			bindDataVal.Field(i).Set(fieldValue)
		}

		if strings.Contains(c.FullPath(), ":") {
			c.ShouldBindUri(bindData)
		}

		if err := c.ShouldBind(bindData); err != nil {
			new(bm.Res).FailFront(err2Str(err)).Send(c)
			c.Abort()
			return
		}

		// 将绑定的数据存储到上下文中
		c.Set("reqBind_", bindData)
		c.Next()
	}
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

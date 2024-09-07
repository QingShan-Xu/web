package rt

import (
	"log"
	"reflect"

	"github.com/gin-gonic/gin"
)

func bindReqMiddlewares(bindVal reflect.Value, name string) gin.HandlerFunc {

	if bindVal.Type() == nil {
		return func(c *gin.Context) {
			c.Next()
		}
	}

	if bindVal.Kind() != reflect.Struct {
		log.Fatalf("%s Bind字段必须是 struct 类型", name)
	}

	return func(c *gin.Context) {
		newBind := reflect.New(bindVal.Type()).Interface()
	}
}

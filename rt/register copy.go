package rt

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/QingShan-Xu/xjh/bm"
	"github.com/QingShan-Xu/xjh/cf"
	"github.com/QingShan-Xu/xjh/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

func reg1ister(groupRouter *gin.RouterGroup, route *Router) {
	isGRoup := IsGroupRouter(*route)
	if isGRoup {
		registerGroup(groupRouter, route)
	} else {
		registerRouter(groupRouter, route)
	}
}

func registerGroup(groupRouter *gin.RouterGroup, route *Router) {
	newGroupRouter := groupRouter.Group(route.Path)
	newGroupRouter.Use(route.Middlewares...)
	for _, child := range route.Children {
		if route.NoToken {
			child.NoToken = route.NoToken
		}
		register(newGroupRouter, &child)
	}
}

func registe1rRouter(groupRouter *gin.RouterGroup, route *Router) {
	handls := []gin.HandlerFunc{}
	if !route.NoToken {
		handls = append(handls, handlToken())
	}

	if route.Bind != nil {
		handls = append(handls, handleBindAll(route))
	}

	handls = append(handls, handleTX(route))

	if route.Type != "" && route.MODEL != nil {
		handls = append(handls, handleTYPE(route))
	}

	router := groupRouter.Handle(
		route.Method,
		route.Path,
		append(handls, route.Handler)...,
	)
	router.Use(route.Middlewares...)
}

var type2Front = map[string]string{
	"int64":   "数字",
	"int32":   "数字",
	"int":     "数字",
	"float64": "数字",
	"float32": "数字",
	"string":  "数字",
}

func handlToken() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var res bm.Response

		tokenStr := ctx.GetHeader("Authorization")
		if tokenStr == "" {
			res.Callback = "/"
			res.FailFront("请先登录后操作").Send(ctx)
			ctx.Abort()
			return
		}

		user_id, err := utils.Utils.Token.ParseToken(tokenStr)
		if err != nil {
			res.Callback = "/"
			res.FailFront(err.Error()).Send(ctx)
			ctx.Abort()
			return
		}

		ctx.Set("user_id", user_id)
		ctx.Next()
	}
}

func handleBindAll(route *Router) gin.HandlerFunc {
	bindVal, err := createDataInstance(route.Bind, true)
	if err != nil {
		log.Fatal(err)
	}

	return func(ctx *gin.Context) {
		newBindVal, _ := createDataInstance(bindVal.Interface(), true)
		route.Bind = newBindVal.Interface()

		var err error

		if strings.Contains(route.Path, ":") {
			err = ctx.ShouldBindUri(route.Bind)
		}

		if !strings.Contains(route.Path, ":") && route.Method == METHOD.GET {
			err = ctx.ShouldBindQuery(route.Bind)
		}

		if !strings.Contains(route.Path, ":") && route.Method != METHOD.GET {
			err = ctx.ShouldBind(route.Bind)
		}

		if err != nil {
			switch e := err.(type) {
			case *json.UnmarshalTypeError:
				a := e.Type.Name()
				eTpe, ok := type2Front[a]
				if !ok {
					eTpe = a
				}
				new(bm.Response).FailFront(fmt.Sprintf("%s: 应为%s类型", e.Field, eTpe)).Send(ctx)
				ctx.Abort()
				return
			case validator.ValidationErrors:
				// 处理验证错误
				errStr := ""
				for _, validatoE := range e.Translate(cf.Trans) {
					errStr += utils.ToSnakeCase(validatoE) + ";"
				}
				new(bm.Response).FailFront(errStr).Send(ctx)
				ctx.Abort()
				return
			default:
				err := e.Error()
				if err == "EOF" {
					new(bm.Response).FailFront("不得为空").Send(ctx)
				} else {
					new(bm.Response).FailFront(err).Send(ctx)
				}
				ctx.Abort()
				return
			}
		}
		ctx.Set("reqBind", route.Bind)
		ctx.Next()
	}
}

func handleTX(route *Router) gin.HandlerFunc {
	var noVModelVal reflect.Value
	var err error

	if route.MODEL != nil {
		noVModelVal, err = createDataInstance(route.MODEL, false)
		if err != nil {
			log.Fatal(err)
		}

		if !noVModelVal.CanInterface() {
			log.Fatalf("%s :%s 初始化DB失败", route.Name, route.Path)
		}

		if err := cf.ORMDB.AutoMigrate(noVModelVal.Interface()); err != nil {
			log.Fatalf("%s :%s 创建表失败 %v", route.Name, route.Path, err)
		}
	}

	return func(c *gin.Context) {
		tx := cf.ORMDB
		if route.MODEL != nil {
			newModelVal, _ := createDataInstance(noVModelVal.Interface(), false)
			tx = tx.Model(newModelVal.Interface())
			genDB(tx, route)
		}
		c.Set("reqTX", tx.Session(&gorm.Session{}))
		c.Next()
	}
}
func handleTYPE(route *Router) gin.HandlerFunc {
	return func(c *gin.Context) {
		tx := c.MustGet("reqTX").(*gorm.DB)
		res := genType(tx, route)
		if res != nil {
			res.Send(c)
			c.Abort()
			return
		}
		c.Next()
	}
}

func createDataInstance(dataType interface{}, copyValues bool) (reflect.Value, error) {

	if dataType == nil {
		return reflect.Value{}, fmt.Errorf("dataType cannot be nil")
	}

	dataReflectType := reflect.TypeOf(dataType)
	dataReflectValue := reflect.ValueOf(dataType)

	if dataReflectType.Kind() == reflect.Pointer {
		dataReflectType = dataReflectType.Elem()
		dataReflectValue = dataReflectValue.Elem()
	}

	if dataReflectType.Kind() != reflect.Struct && dataReflectType.Kind() != reflect.Slice {
		return reflect.Value{}, fmt.Errorf("无法创建实例：不支持的类型 %s", dataReflectType.Kind())
	}

	newInstance := reflect.New(dataReflectType)

	if copyValues {
		newInstance.Elem().Set(dataReflectValue)
	}

	return newInstance, nil
}

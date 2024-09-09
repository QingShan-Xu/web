package middleware

import (
	"log"
	"reflect"

	"github.com/QingShan-Xu/xjh/bm"
	"github.com/QingShan-Xu/xjh/rt/internal/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func ReqFinisherMiddleware(
	Finisher string,
	MODEL interface{},
	name string,
) gin.HandlerFunc {

	if MODEL == nil {
		log.Fatalf("%s: MODEL 在使用 Finisher 时不能为空", name)
	}

	if Finisher == "First" {
		return func(ctx *gin.Context) {
			tx := ctx.MustGet("reqTX_").(*gorm.DB)
			data := reflect.New(utils.GetInstanceVal(MODEL).Type()).Interface()

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

	log.Fatalf("%s: %s 该方法还未在 Finisher 中实现", name, Finisher)
	return func(ctx *gin.Context) { ctx.Next() }
}

package middleware

import (
	"log"
	"reflect"

	"github.com/QingShan-Xu/xjh/bm"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func ReqFinisherMiddleware(
	BeforeFinisher func(bind interface{}) interface{},
	Finisher string,
	name string,
) gin.HandlerFunc {

	if Finisher == "First" {
		return func(ctx *gin.Context) {
			tx := ctx.MustGet("reqTX_").(*gorm.DB)
			model := ctx.MustGet("reqModel_")

			result := tx.Find(model)
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

			new(bm.Res).SucJson(model).Send(ctx)
			ctx.Abort()
		}
	}

	if Finisher == "Create" {
		return func(ctx *gin.Context) {
			var bind interface{}
			reqBind := ctx.MustGet("reqBind_")
			reqModel := ctx.MustGet("reqModel_")

			tx := ctx.MustGet("reqTX_").(*gorm.DB)
			if BeforeFinisher != nil {
				beforeBind := BeforeFinisher(reqBind)

				if reflect.TypeOf(beforeBind).Kind() != reflect.Pointer {
					beforeBind = &beforeBind
				}
				if reflect.TypeOf(beforeBind).Kind() != reflect.TypeOf(reqModel).Kind() {
					new(bm.Res).FailBackend(name, ": 创建对象 与 数据结构不一致").Send(ctx)
					ctx.Abort()
					return
				}
			} else {
				bind = reqBind
			}

			result := tx.Create(bind)
			if result.Error != nil {
				new(bm.Res).FailBackend(result.Error).Send(ctx)
				ctx.Abort()
				return
			}

			new(bm.Res).SucJson(bind).Send(ctx)
			ctx.Abort()
		}
	}

	if Finisher == "Update" {
		return func(ctx *gin.Context) {
			var bind interface{}
			reqBind := ctx.MustGet("reqBind_")
			tx := ctx.MustGet("reqTX_").(*gorm.DB)
			if BeforeFinisher != nil {
				beforBind := BeforeFinisher(reqBind)
				bind = beforBind
			} else {
				bind = reqBind
			}
			// mysql不支持Returning
			result := tx.Clauses(clause.Returning{}).Updates(bind)
			if result.Error != nil {
				new(bm.Res).FailBackend(result.Error).Send(ctx)
				ctx.Abort()
				return
			}

			new(bm.Res).SucJson(nil).Send(ctx)
			ctx.Abort()
		}
	}

	if Finisher == "Delete" {
		return func(ctx *gin.Context) {
			tx := ctx.MustGet("reqTX_").(*gorm.DB)
			reqModel := ctx.MustGet("reqModel_")

			result := tx.Delete(reqModel)
			if result.Error != nil {
				new(bm.Res).FailBackend(result.Error).Send(ctx)
				ctx.Abort()
				return
			}

			new(bm.Res).SucJson(nil).Send(ctx)
			ctx.Abort()
		}
	}

	log.Fatalf("%s: %s 该方法还未在 Finisher 中实现", name, Finisher)
	return func(ctx *gin.Context) { ctx.Next() }
}

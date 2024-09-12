package middleware

import (
	"log"
	"reflect"

	"github.com/QingShan-Xu/xjh/bm"
	"github.com/QingShan-Xu/xjh/rt/internal/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func ReqTypeMiddleware(
	BeforeInset func(bind interface{}) interface{},
	Type string,
	name string,
) gin.HandlerFunc {

	if Type == "GET_ONE" {
		return func(ctx *gin.Context) {
			tx := ctx.MustGet("reqTX_").(*gorm.DB)
			model := ctx.MustGet("reqModel_")

			result := tx.First(model)
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

	if Type == "GET_LIST" {
		return func(ctx *gin.Context) {
			tx := ctx.MustGet("reqTX_").(*gorm.DB)
			model := ctx.MustGet("reqModel_")
			bind := ctx.MustGet("reqBind_")
			bindMap := utils.Struct2map(bind, true)
			modelList := reflect.New(reflect.SliceOf(reflect.TypeOf(model))).Interface()
			var total int64

			if err := tx.Count(&total).Error; err != nil {
				new(bm.Res).FailBackend(err).Send(ctx)
				ctx.Abort()
				return
			}

			pageSize := bindMap["Pagination"].(map[string]interface{})["PageSize"].(int)
			current := bindMap["Pagination"].(map[string]interface{})["Current"].(int)

			if err := tx.Limit(pageSize).Offset((current - 1) * pageSize).Find(modelList).Error; err != nil {
				new(bm.Res).FailBackend(err).Send(ctx)
				ctx.Abort()
				return
			}

			new(bm.Res).SucJson(bm.ResList{
				Data:     modelList,
				Total:    total,
				PageSize: pageSize,
				Current:  current,
			}).Send(ctx)
			ctx.Abort()
		}
	}

	if Type == "CREATE_ONE" {
		return func(ctx *gin.Context) {
			var bind interface{}
			reqBind := ctx.MustGet("reqBind_")
			reqModel := ctx.MustGet("reqModel_")

			tx := ctx.MustGet("reqTX_").(*gorm.DB)
			if BeforeInset != nil {
				beforeBind := BeforeInset(reqBind)

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

	if Type == "UPDATE_ONE" {
		return func(ctx *gin.Context) {
			var bind interface{}
			reqBind := ctx.MustGet("reqBind_")
			tx := ctx.MustGet("reqTX_").(*gorm.DB)
			if BeforeInset != nil {
				beforBind := BeforeInset(reqBind)
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

	if Type == "DELETE_ONE" {
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

	log.Fatalf("%s: %s 该方法还未在 Type 中实现", name, Type)
	return func(ctx *gin.Context) { ctx.Next() }
}

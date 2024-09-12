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
			tx := ctx.MustGet("reqTX_").(*gorm.DB)
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
			tx := ctx.MustGet("reqTX_").(*gorm.DB)

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

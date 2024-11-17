// Package rt 提供了请求处理的核心逻辑。
package rt

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/QingShan-Xu/web/bm"
	"github.com/QingShan-Xu/web/db"
	"github.com/QingShan-Xu/web/ds"
	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

// handler 处理请求的核心逻辑。
type handler struct {
	Router *Router
}

// serveHTTP 实现 http.Handler 接口。
// w: HTTP 响应写入器。
// r: HTTP 请求。
func (h *handler) serveHTTP(w http.ResponseWriter, r *http.Request) *bm.Res {
	var bindData interface{}
	var err error
	response := bm.NewRes(w)

	if h.Router.Bind != nil {
		// 数据绑定和验证。
		binder := newBinder()
		bindData, err = binder.bindAndValidate(h.Router.Bind, r)
		if err != nil {
			return response.FailFront(err)
		}
	}

	var bindReader ds.FieldReader
	if bindData != nil {
		// 使用 ds 包解析绑定数据。
		bindReader, err = ds.NewStructReader(bindData)
		if err != nil {
			return response.FailBackend(err)
		}
	}

	fmt.Printf("%+v", bindData)
	currentDB := db.DB.GORM.Session(&gorm.Session{})

	// 应用查询范围（Scopes）。
	var scopes []func(db *gorm.DB) *gorm.DB
	for _, scope := range h.Router.Scopes {
		scopes = append(scopes, scope(bindReader))
	}
	currentDB = currentDB.Scopes(scopes...)

	// 检查是否同时设置了多个 Finisher 方法。
	finisherMethodCount := 0
	if h.Router.CreateOne != nil {
		finisherMethodCount++
	}
	if h.Router.UpdateOne != nil {
		finisherMethodCount++
	}
	if h.Router.DeleteOne {
		finisherMethodCount++
	}
	if h.Router.GetOne {
		finisherMethodCount++
	}
	if h.Router.GetList {
		finisherMethodCount++
	}

	if finisherMethodCount > 1 {
		fmt.Println("Cannot use multiple finisher methods simultaneously")
		return response.FailBackend("Cannot use multiple finisher methods simultaneously")

	}

	if finisherMethodCount > 0 && h.Router.Model == nil {
		fmt.Println("Router.Model cannot be nil when using finisher methods")
		return response.FailBackend("Router.Model cannot be nil when using finisher methods")
	}

	// 处理各个 Finisher 方法。
	switch {
	case h.Router.CreateOne != nil:
		// 处理创建操作。
		finisherParams, err := h.genCreateParams(bindReader)
		if err != nil {
			return response.FailFront(err)

		}
		if err := currentDB.Create(finisherParams).Error; err != nil {
			return response.FailFront(err)
		}
		return response.SucJson(finisherParams)

	case h.Router.UpdateOne != nil:
		// 处理更新操作。
		newModel := reflect.New(reflect.TypeOf(h.Router.Model)).Interface()
		if err := currentDB.First(newModel).Error; err != nil {
			return response.FailFront("No corresponding data")
		}

		finisherParams, err := h.genUpdateParams(bindReader, newModel)
		if err != nil {
			return response.FailFront(err)
		}
		if err := db.DB.GORM.Save(finisherParams).Error; err != nil {
			return response.FailFront(err)
		}
		return response.SucJson(finisherParams)

	case h.Router.DeleteOne:
		// 处理删除操作。
		newModel := reflect.New(reflect.TypeOf(h.Router.Model)).Interface()
		if err := currentDB.First(newModel).Error; err != nil {
			return response.FailFront("No corresponding data")

		}
		if err := currentDB.Delete(newModel).Error; err != nil {
			return response.FailFront(err)
		}
		return response.SucJson(newModel)

	case h.Router.GetOne:
		// 处理获取单个记录操作。
		newModel := reflect.New(reflect.TypeOf(h.Router.Model)).Interface()
		if err := currentDB.First(newModel).Error; err != nil {
			return response.FailFront("No corresponding data")

		}
		return response.SucJson(newModel)

	case h.Router.GetList:
		// 处理获取列表操作。
		pagination := bm.Pagination{
			PageSize: 10,
			Current:  1,
		}

		if bindReader != nil {
			// 获取分页参数。
			if pageSizeField, err := bindReader.GetField("PageSize"); err == nil {
				if pageSizeValue, ok := pageSizeField.SafeInt(); ok && pageSizeValue > 0 && pageSizeValue <= 100 {
					pagination.PageSize = pageSizeValue
				}
			}
			if currentField, err := bindReader.GetField("Current"); err == nil {
				if currentValue, ok := currentField.SafeInt(); ok && currentValue > 0 {
					pagination.Current = currentValue
				}
			}
		}

		var total int64
		if err := currentDB.Count(&total).Error; err != nil {
			return response.FailBackend("Query failed")

		}
		if total == 0 {
			return response.SucList(bm.ResList{
				Pagination: pagination,
				Data:       []interface{}{},
				Total:      total,
			})

		}

		currentDB = currentDB.Scopes(PaginationScope(pagination))
		newModelSlice := reflect.New(reflect.SliceOf(reflect.TypeOf(h.Router.Model))).Interface()
		if err := currentDB.Find(newModelSlice).Error; err != nil {
			return response.FailBackend("Query failed")

		}

		return response.SucList(bm.ResList{
			Pagination: pagination,
			Data:       newModelSlice,
			Total:      total,
		})
	}

	if h.Router.Handler != nil {
		currentDB.Transaction(func(tx *gorm.DB) error {
			response = h.Router.Handler(HandlerParams{
				W:          w,
				R:          r,
				Res:        bm.NewRes(w),
				Tx:         tx,
				BindReader: NewBinderReader(bindReader),
			})

			if response.Code != 200 {
				return fmt.Errorf(response.Msg)
			}

			return nil
		})
	}

	return response
}

// genCreateParams 生成创建操作的参数。
// bindReader: 绑定数据的结构体读取器。
// 返回生成的模型实例或错误信息。
func (h *handler) genCreateParams(bindReader ds.FieldReader) (interface{}, error) {
	newModel := reflect.New(reflect.TypeOf(h.Router.Model)).Interface()
	modelReader, err := ds.NewStructReader(newModel)
	if err != nil {
		return nil, err
	}

	modelMap := make(map[string]interface{})
	for modelField, bindField := range h.Router.CreateOne {
		if _, err := modelReader.GetField(modelField); err != nil {
			return nil, fmt.Errorf("model lacks field '%s'", modelField)
		}
		bindValueField, err := bindReader.GetField(bindField)
		if err != nil {
			return nil, fmt.Errorf("bind data lacks field '%s'", bindField)
		}
		bindValue := bindValueField.Interface()
		if bindValue == nil {
			continue
		}
		modelMap[modelField] = bindValue
	}

	decoderConfig := &mapstructure.DecoderConfig{
		Squash:               true,
		WeaklyTypedInput:     true,
		TagName:              "bind",
		IgnoreUntaggedFields: true,
		Result:               newModel,
	}
	decoder, _ := mapstructure.NewDecoder(decoderConfig)
	if err := decoder.Decode(modelMap); err != nil {
		return nil, err
	}

	return newModel, nil
}

// genUpdateParams 生成更新操作的参数。
// bindReader: 绑定数据的结构体读取器。
// model: 当前数据库中已有的模型实例。
// 返回更新后的模型实例或错误信息。
func (h *handler) genUpdateParams(bindReader ds.FieldReader, model interface{}) (interface{}, error) {
	modelReader, err := ds.NewStructReader(model)
	if err != nil {
		return nil, err
	}

	modelMap := make(map[string]interface{})
	for modelField, bindField := range h.Router.UpdateOne {
		if _, err := modelReader.GetField(modelField); err != nil {
			return nil, fmt.Errorf("model lacks field '%s'", modelField)
		}
		bindValueField, err := bindReader.GetField(bindField)
		if err != nil {
			return nil, fmt.Errorf("bind data lacks field '%s'", bindField)
		}
		bindValue := bindValueField.Interface()
		if bindValue == nil {
			continue
		}
		modelMap[modelField] = bindValue
	}

	decoderConfig := &mapstructure.DecoderConfig{
		Squash:               true,
		WeaklyTypedInput:     true,
		TagName:              "bind",
		IgnoreUntaggedFields: true,
		Result:               model,
		ZeroFields:           true,
	}
	decoder, _ := mapstructure.NewDecoder(decoderConfig)
	if err := decoder.Decode(modelMap); err != nil {
		return nil, err
	}

	return model, nil
}

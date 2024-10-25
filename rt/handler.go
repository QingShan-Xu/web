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

// Handler 处理请求的核心逻辑。
type Handler struct {
	Router *Router
}

// ServeHTTP 实现 http.Handler 接口。
// w: HTTP 响应写入器。
// r: HTTP 请求。
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var bindData interface{}
	var err error
	response := bm.NewRes(w)

	if h.Router.Bind != nil {
		// 数据绑定和验证。
		binder := NewBinder()
		bindData, err = binder.BindAndValidate(h.Router.Bind, r)
		if err != nil {
			response.FailFront(err.Error()).Send()
			return
		}
	}

	fmt.Printf("%+v\n", bindData)

	// 使用 ds 包解析绑定数据。
	bindReader, err := ds.NewStructReader(bindData)
	if err != nil {
		response.FailBackend(err).Send()
		return
	}

	currentDB := db.DB.GORM.Session(&gorm.Session{})

	// 应用查询范围（Scopes）。
	var scopes []func(db *gorm.DB) *gorm.DB
	for _, scope := range h.Router.Scopes {
		scopes = append(scopes, scope(bindReader))
	}
	currentDB = currentDB.Scopes(scopes...)

	// 检查是否同时设置了多个 Finisher 方法。
	finisherMethodCount := 0
	if h.Router.CreateFields != nil {
		finisherMethodCount++
	}
	if h.Router.UpdateFields != nil {
		finisherMethodCount++
	}
	if h.Router.Delete {
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
		response.FailBackend("Cannot use multiple finisher methods simultaneously").Send()
		return
	}

	if finisherMethodCount > 0 && h.Router.Model == nil {
		fmt.Println("Router.Model cannot be nil when using finisher methods")
		response.FailBackend("Router.Model cannot be nil when using finisher methods").Send()
		return
	}

	// 处理各个 Finisher 方法。
	switch {
	case h.Router.CreateFields != nil:
		// 处理创建操作。
		finisherParams, err := h.genCreateParams(bindReader)
		if err != nil {
			response.FailFront(err.Error()).Send()
			return
		}
		if err := currentDB.Create(finisherParams).Error; err != nil {
			response.FailFront(err.Error()).Send()
			return
		}
		response.SucJson(finisherParams).Send()
		return

	case h.Router.UpdateFields != nil:
		// 处理更新操作。
		newModel := reflect.New(reflect.TypeOf(h.Router.Model)).Interface()
		if err := currentDB.First(newModel).Error; err != nil {
			response.FailFront("No corresponding data").Send()
			return
		}

		finisherParams, err := h.genUpdateParams(bindReader, newModel)
		if err != nil {
			response.FailFront(err.Error()).Send()
			return
		}
		if err := currentDB.Updates(finisherParams).Error; err != nil {
			response.FailFront(err.Error()).Send()
			return
		}
		response.SucJson(finisherParams).Send()
		return

	case h.Router.Delete:
		// 处理删除操作。
		newModel := reflect.New(reflect.TypeOf(h.Router.Model)).Interface()
		if err := currentDB.First(newModel).Error; err != nil {
			response.FailFront("No corresponding data").Send()
			return
		}
		if err := currentDB.Delete(newModel).Error; err != nil {
			response.FailFront(err.Error()).Send()
			return
		}
		response.SucJson(newModel).Send()
		return

	case h.Router.GetOne:
		// 处理获取单个记录操作。
		newModel := reflect.New(reflect.TypeOf(h.Router.Model)).Interface()
		if err := currentDB.First(newModel).Error; err != nil {
			response.FailFront("No corresponding data").Send()
			return
		}
		response.SucJson(newModel).Send()
		return

	case h.Router.GetList:
		// 处理获取列表操作。
		pagination := bm.Pagination{
			PageSize: 10,
			Current:  1,
		}

		if bindReader != nil {
			// 获取分页参数。
			if pageSizeField, err := bindReader.GetFieldByName("PageSize"); err == nil {
				if pageSizeValue, ok := pageSizeField.Value.Interface().(int); ok && pageSizeValue > 0 && pageSizeValue <= 100 {
					pagination.PageSize = pageSizeValue
				}
			}
			if currentField, err := bindReader.GetFieldByName("Current"); err == nil {
				if currentValue, ok := currentField.Value.Interface().(int); ok && currentValue > 0 {
					pagination.Current = currentValue
				}
			}
		}

		var total int64
		if err := currentDB.Count(&total).Error; err != nil {
			response.FailBackend("Query failed").Send()
			return
		}
		if total == 0 {
			response.SucList(bm.ResList{
				Pagination: pagination,
				Data:       []interface{}{},
				Total:      total,
			}).Send()
			return
		}

		currentDB = currentDB.Scopes(PaginationScope(pagination))
		newModelSlice := reflect.New(reflect.SliceOf(reflect.TypeOf(h.Router.Model))).Interface()
		if err := currentDB.Find(newModelSlice).Error; err != nil {
			response.FailBackend("Query failed").Send()
			return
		}

		response.SucList(bm.ResList{
			Pagination: pagination,
			Data:       newModelSlice,
			Total:      total,
		}).Send()
		return
	}
}

// genCreateParams 生成创建操作的参数。
// bindReader: 绑定数据的结构体读取器。
// 返回生成的模型实例或错误信息。
func (h *Handler) genCreateParams(bindReader *ds.StructReader) (interface{}, error) {
	newModel := reflect.New(reflect.TypeOf(h.Router.Model)).Interface()
	modelReader, err := ds.NewStructReader(newModel)
	if err != nil {
		return nil, err
	}

	modelMap := make(map[string]interface{})
	for modelField, bindField := range h.Router.CreateFields {
		if _, err := modelReader.GetFieldByName(modelField); err != nil {
			return nil, fmt.Errorf("Model lacks field '%s'", modelField)
		}
		bindValueField, err := bindReader.GetFieldByName(bindField)
		if err != nil {
			return nil, fmt.Errorf("Bind data lacks field '%s'", bindField)
		}
		bindValue := bindValueField.Value.Interface()
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
func (h *Handler) genUpdateParams(bindReader *ds.StructReader, model interface{}) (interface{}, error) {
	modelReader, err := ds.NewStructReader(model)
	if err != nil {
		return nil, err
	}

	modelMap := make(map[string]interface{})
	for modelField, bindField := range h.Router.UpdateFields {
		if _, err := modelReader.GetFieldByName(modelField); err != nil {
			return nil, fmt.Errorf("Model lacks field '%s'", modelField)
		}
		bindValueField, err := bindReader.GetFieldByName(bindField)
		if err != nil {
			return nil, fmt.Errorf("Bind data lacks field '%s'", bindField)
		}
		bindValue := bindValueField.Value.Interface()
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
	}
	decoder, _ := mapstructure.NewDecoder(decoderConfig)
	if err := decoder.Decode(modelMap); err != nil {
		return nil, err
	}

	return model, nil
}

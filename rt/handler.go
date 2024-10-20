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

func (curRT *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var bindData interface{}
	var err error
	response := bm.NewRes(w)

	if curRT.Bind != nil {
		dataBinder := NewDataBinder()
		bindData, err = dataBinder.BindData(curRT, r)
		if err != nil {
			response.FailFront(err.Error()).Send()
			return
		}
	}

	fmt.Println(fmt.Sprintf("%+v", bindData))

	bindReader, err := ds.NewStructReader(bindData)

	if err != nil {
		response.FailBackend(err).Send()
	}

	currentDB := db.DB.GORM.Session(&gorm.Session{})

	scopes := []func(db *gorm.DB) *gorm.DB{}
	for _, scope := range curRT.SCOPES {
		scopes = append(scopes, scope(bindReader))
	}

	currentDB = currentDB.Scopes(scopes...)

	if err != nil {
		response.FailBackend().Send()
		return
	}

	if curRT.CREATE_ONE != nil && curRT.UPDATE_ONE != nil && curRT.DELETE_ONE && curRT.GET_ONE && curRT.GET_LIST {
		fmt.Println("不支持同时使用 Finisher Method")
		return
	}
	if (curRT.CREATE_ONE != nil || curRT.UPDATE_ONE != nil || curRT.DELETE_ONE && curRT.GET_ONE && curRT.GET_LIST) && curRT.MODEL == nil {
		fmt.Println("使用 Finisher Method 时, rt.MODEL 不能为空")
		response.FailBackend().Send()
		return
	}

	if curRT.CREATE_ONE != nil {
		finisherParams, err := curRT.genCreateOneParams(bindReader)
		if err != nil {
			response.FailFront(err).Send()
			return
		}
		if err := currentDB.Create(finisherParams).Error; err != nil {
			response.FailFront(err).Send()
			return
		}
		response.SucJson(finisherParams).Send()
		return
	}

	if curRT.UPDATE_ONE != nil {
		newOrm := reflect.New(reflect.TypeOf(curRT.MODEL)).Interface()
		firstQ := currentDB.First(newOrm)
		if firstQ.RowsAffected == 0 {
			response.FailFront("没有对应数据").Send()
			return
		}
		if firstQ.Error != nil {
			response.FailFront(firstQ.Error.Error()).Send()
			return
		}

		finisherParams, err := curRT.genUpdateOneParams(bindReader, newOrm)
		if err != nil {
			response.FailFront(err.Error()).Send()
			return
		}
		if currentDB.Updates(finisherParams).Error != nil {
			response.FailFront(err).Send()
			return
		}
		response.SucJson(finisherParams).Send()
		return
	}

	if curRT.DELETE_ONE {
		newOrm := reflect.New(reflect.TypeOf(curRT.MODEL)).Interface()
		firstQ := currentDB.First(newOrm)
		if firstQ.RowsAffected == 0 {
			response.FailFront("没有对应数据").Send()
			return
		}
		if firstQ.Error != nil {
			response.FailFront(firstQ.Error.Error()).Send()
			return
		}

		if err := currentDB.Delete(newOrm).Error; err != nil {
			response.FailFront(err.Error()).Send()
			return
		}
		response.SucJson(newOrm).Send()
		return
	}

	if curRT.GET_ONE {
		newOrm := reflect.New(reflect.TypeOf(curRT.MODEL)).Interface()
		firstQ := currentDB.First(newOrm)
		if firstQ.RowsAffected == 0 {
			response.FailFront("没有对应数据").Send()
			return
		}
		if firstQ.Error != nil {
			response.FailFront(firstQ.Error.Error()).Send()
			return
		}
		response.SucJson(newOrm).Send()
		return
	}

	if curRT.GET_LIST {
		pagination := bm.Pagination{
			PageSize: 10,
			Current:  1,
		}
		pageSizeField, _ := bindReader.GetFieldByName("PageSize")
		pageSizeValue := pageSizeField.GetValue()
		currentField, _ := bindReader.GetFieldByName("Current")
		currentValue := currentField.GetValue()

		if pageSizeValue != 0 && pageSizeValue.(int) <= 100 {
			pagination.PageSize = pageSizeValue.(int)
		}
		pagination.Current = currentValue.(int)

		var total int64
		if err := currentDB.Count(&total).Error; err != nil {
			response.FailBackend("查询失败").Send()
		}
		if total == 0 {
			response.SucList(bm.ResList{
				Pagination: pagination,
				Data:       []interface{}{},
				Total:      total,
			})
			return
		}

		currentDB.Scopes(PaginationScope(pagination))
		newOrm := reflect.New(reflect.SliceOf(reflect.TypeOf(curRT.MODEL))).Interface()
		if err := currentDB.Find(newOrm).Error; err != nil {
			response.FailBackend("查询失败").Send()
			return
		}

		response.SucList(bm.ResList{
			Pagination: pagination,
			Data:       newOrm,
			Total:      total,
		}).Send()
		return
	}
}

func (curRT *Router) genCreateOneParams(bindReader *ds.StructReader) (interface{}, error) {
	newOrm := reflect.New(reflect.TypeOf(curRT.MODEL)).Interface()
	ormReader, err := ds.NewStructReader(newOrm)
	if err != nil {
		return nil, err
	}

	ormMap := map[string]interface{}{}
	for ormName, bindName := range curRT.CREATE_ONE {
		if _, err := ormReader.GetFieldByName(ormName); err != nil {
			return nil, fmt.Errorf("请检查ORM中有无 %s 字段", ormName)
		}
		bindField, err := bindReader.GetFieldByName(bindName)
		if err != nil {
			return nil, fmt.Errorf("请求值 %s 缺失", bindName)
		}
		bindValue := bindField.Value.Interface()
		if bindValue == nil {
			continue
		}
		ormMap[ormName] = bindValue
	}

	decoder, _ := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Squash:               true,
		WeaklyTypedInput:     true,
		TagName:              "bind",
		IgnoreUntaggedFields: true,
		Result:               newOrm,
	})
	if err := decoder.Decode(ormMap); err != nil {
		return nil, err
	}

	return newOrm, nil
}

func (curRT *Router) genUpdateOneParams(bindReader *ds.StructReader, ormModel interface{}) (interface{}, error) {

	ormReader, err := ds.NewStructReader(ormModel)
	if err != nil {
		return nil, err
	}

	ormMap := map[string]interface{}{}
	for ormName, bindName := range curRT.UPDATE_ONE {
		if _, err := ormReader.GetFieldByName(ormName); err != nil {
			return nil, fmt.Errorf("rt.MODEL中无 %s 缺失", ormName)
		}
		bindField, err := bindReader.GetFieldByName(bindName)
		if err != nil {
			return nil, fmt.Errorf("rt.Bind中无 %s 字段", bindName)
		}
		bindValue := bindField.Value.Interface()
		if bindValue == nil {
			continue
		}
		ormMap[ormName] = bindValue
	}

	decoder, _ := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Squash:               true,
		WeaklyTypedInput:     true,
		TagName:              "bind",
		IgnoreUntaggedFields: true,
		Result:               ormModel,
	})
	if err := decoder.Decode(ormMap); err != nil {
		return nil, err
	}

	fmt.Printf("%+v", ormModel)
	return ormModel, nil
}

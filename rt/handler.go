package rt

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/QingShan-Xu/web/bm"
	"github.com/QingShan-Xu/web/db"
	"github.com/QingShan-Xu/web/ds"
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

	fmt.Printf("%+v", bindData)

	bindReader := ds.NewReader(bindData)

	currentDB := db.DB.GORM.Session(&gorm.Session{})

	scopes := []func(db *gorm.DB) *gorm.DB{}
	for _, scope := range curRT.SCOPES {
		scopes = append(scopes, scope(bindReader))
	}

	currentDB.Scopes(scopes...)

	if curRT.CREATE_ONE != nil && curRT.UPDATE_ONE != nil && curRT.DELETE_ONE && curRT.GET_ONE && curRT.GET_LIST {
		fmt.Println("不支持同时使用 Finisher Method")
		response.FailBackend().Send()
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
			response.FailFront(err.Error()).Send()
			return
		}
		if err := currentDB.Create(finisherParams).Error; err != nil {
			response.FailFront(err.Error()).Send()
			return
		}
		response.SucJson(finisherParams)
	}

	if curRT.UPDATE_ONE != nil {
		finisherParams, err := curRT.genUpdateOneParams(bindReader)
		if err != nil {
			response.FailFront(err.Error()).Send()
			return
		}
		if err := currentDB.Updates(finisherParams).Error; err != nil {
			response.FailFront(err.Error()).Send()
			return
		}
	}

	if curRT.DELETE_ONE {
		newOrm := reflect.New(reflect.TypeOf(curRT.MODEL)).Interface()
		if err := currentDB.Delete(newOrm).Error; err != nil {
			response.FailFront(err.Error()).Send()
			return
		}
	}

	if curRT.GET_ONE {
		newOrm := reflect.New(reflect.TypeOf(curRT.MODEL)).Interface()
		firstQ := currentDB.First(newOrm)
		if firstQ.Error != nil {
			response.FailFront(firstQ.Error.Error()).Send()
			return
		}
		if firstQ.RowsAffected == 0 {
			response.FailFront("没有对应数据").Send()
			return
		}
	}

	if curRT.GET_LIST {
		pageSize := 10
		if bindReader.HasField("PageSize") {
			pageSize = bindReader.GetField("PageSize").Int()
		}
		current := 1
		if bindReader.HasField("current") {
			current = bindReader.GetField("Current").Int()
		}

		var total int64
		if err := currentDB.Model(curRT.MODEL).Count(&total); err != nil {
			response.FailBackend("查询失败").Send()
		}
		if total == 0 {
			response.SucList(bm.ResList{
				Current:  current,
				Data:     []interface{}{},
				PageSize: pageSize,
				Total:    total,
			})
			return
		}

		currentDB.Scopes(PaginationScope(pageSize, current))
		newOrm := reflect.New(reflect.SliceOf(reflect.TypeOf(curRT.MODEL))).Interface()
		if err := currentDB.Find(newOrm).Error; err != nil {
			response.FailBackend("查询失败").Send()
			return
		}

		response.SucList(bm.ResList{
			Current:  current,
			Data:     newOrm,
			PageSize: pageSize,
			Total:    total,
		})
		return
	}
}

func (curRT *Router) genCreateOneParams(bindReader ds.Reader) (interface{}, error) {
	newOrm := reflect.New(reflect.TypeOf(curRT.MODEL)).Interface()
	ormReader := ds.NewReader(newOrm)
	for ormName, bindName := range curRT.CREATE_ONE {
		bindField := bindReader.GetField(bindName)
		if bindField == nil {
			return nil, fmt.Errorf("请求值 %s 缺失", bindName)
		}
		if !ormReader.HasField(ormName) {
			return nil, fmt.Errorf("后台值 %s 缺失", ormName)
		}
		if err := ormReader.SetValue(ormName, bindField.Interface()); err != nil {
			return nil, err
		}
	}
	fmt.Printf("%+v", newOrm)
	return newOrm, nil
}

func (curRT *Router) genUpdateOneParams(bindReader ds.Reader) (interface{}, error) {
	newOrm := reflect.New(reflect.TypeOf(curRT.MODEL)).Interface()
	ormReader := ds.NewReader(newOrm)
	for ormName, bindName := range curRT.UPDATE_ONE {
		bindField := bindReader.GetField(bindName)
		if bindField == nil {
			return nil, fmt.Errorf("请求值 %s 缺失", bindName)
		}
		if !ormReader.HasField(ormName) {
			return nil, fmt.Errorf("后台值 %s 缺失", ormName)
		}
		if err := ormReader.SetValue(ormName, bindField.Interface()); err != nil {
			return nil, err
		}
	}
	fmt.Printf("%+v", newOrm)
	return newOrm, nil
}

// func (curRT *Router) Handler(w http.ResponseWriter, r *http.Request) {

// // WHERE语句
// if curRT.WHERE != nil && curRT.MODEL == nil {
// 	log.Fatalf("%s: MODEL required when WHERE not nil", curRT.Path)
// }
// for query, data := range curRT.WHERE {
// 	db = db.Where(query, data)
// }

// }

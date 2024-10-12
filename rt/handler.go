package rt

import (
	"encoding/json"
	"io"
	"net/http"
	"reflect"
	"strconv"

	"github.com/QingShan-Xu/web/ds"
)

func (curRT *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if curRT.Bind != nil {
		dataBinder := NewDataBinder()
		dataBinder.BindData(curRT, r)
	}
	// if err := h(); err != nil {
	// 	// handle returned error here.
	// 	w.WriteHeader(503)
	// 	w.Write([]byte("bad"))
	// }
}

// JSONBindingStrategy 实现了JSON数据绑定
type JSONBindingStrategy struct{}

// QueryBindingStrategy 实现了查询参数绑定
type QueryBindingStrategy struct{}

// FormBindingStrategy 实现了表单数据绑定
type FormBindingStrategy struct{}

// BindData 将HTTP请求数据绑定到指定结构体
//
// 参数:
//   - curRT: 当前路由信息
//   - r: HTTP请求
//
// 返回:
//   - ds.Reader: 绑定数据的Reader
//   - error: 绑定过程中的错误
func (binder *DataBinder) BindData(curRT *Router, r *http.Request) (ds.Reader, error) {
	bind := reflect.New(reflect.TypeOf(curRT.Bind)).Interface()
	reader := ds.NewReader(&bind)

	bindFieldMap := binder.getBindFieldMap(reader)

	requestData := make(map[string]interface{})

	// 根据Content-Type选择绑定策略
	contentType := r.Header.Get("Content-Type")
	if strategy, ok := binder.strategies[contentType]; ok {
		strategyData, err := strategy.Bind(r, bindFieldMap)
		if err != nil {
			return nil, err
		}
		for k, v := range strategyData {
			requestData[k] = v
		}
	}

	// 设置绑定的值
	for _, bindField := range bindFieldMap {
		if bindData, ok := requestData[bindField.BindName]; ok {
			binder.setValue(reader, bindField, bindData)
		}
	}

	return reader, nil
}

// getBindFieldMap 获取需要绑定的字段映射
func (binder *DataBinder) getBindFieldMap(reader ds.Reader) map[string]BindField {
	bindFieldMap := make(map[string]BindField)
	bindAllFields := reader.GetAllFields()

	for _, field := range bindAllFields {
		fieldTag := field.Tag()
		for _, bindFieldType := range []string{"json", "form", "uri", "query"} {
			if bindName, ok := fieldTag.Lookup(bindFieldType); ok {
				bindField := BindField{
					FieldName: field.Name(),
					FieldKind: field.Kind(),
					BindName:  bindName,
					BindType:  bindFieldType,
				}
				if bindTag, ok := fieldTag.Lookup("bind"); ok {
					bindField.BindTag = bindTag
				}
				bindFieldMap[bindField.BindName] = bindField
			}
		}
	}

	return bindFieldMap
}

// setValue 设置字段的值，处理特殊类型转换
func (db *DataBinder) setValue(reader ds.Reader, bindField BindField, bindData interface{}) error {
	if bindField.BindTag == "strint" {
		return db.setStrIntValue(reader, bindField, bindData)
	}
	reader.SetValue(bindField.FieldName, bindData)
	return nil
}

// setStrIntValue 处理字符串和整数之间的转换
func (db *DataBinder) setStrIntValue(reader ds.Reader, bindField BindField, bindData interface{}) error {
	dataKind := reflect.TypeOf(bindData).Kind()
	if dataKind == bindField.FieldKind {
		reader.SetValue(bindField.FieldName, bindData)
		return nil
	}

	if dataKind == reflect.String && bindField.FieldKind == reflect.Int {
		data, err := strconv.Atoi(bindData.(string))
		if err != nil {
			return err
		}
		reader.SetValue(bindField.FieldName, data)
	} else if dataKind == reflect.Int && bindField.FieldKind == reflect.String {
		data := strconv.Itoa(bindData.(int))
		reader.SetValue(bindField.FieldName, data)
	}

	return nil
}

// Bind 实现JSON绑定策略
func (j *JSONBindingStrategy) Bind(r *http.Request, bindFieldMap map[string]BindField) (map[string]interface{}, error) {
	requestData := make(map[string]interface{})
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	err = json.Unmarshal(body, &requestData)
	if err != nil {
		return nil, err
	}

	return requestData, nil
}

// Bind 实现查询参数绑定策略
func (q *QueryBindingStrategy) Bind(r *http.Request, bindFieldMap map[string]BindField) map[string]interface{} {
	requestData := make(map[string]interface{})
	queryParams := r.URL.Query()

	for key, values := range queryParams {
		if bindField, ok := bindFieldMap[key]; ok {
			if bindField.FieldKind == reflect.Array || bindField.FieldKind == reflect.Slice {
				requestData[bindField.BindName] = values
			} else {
				requestData[bindField.BindName] = values[0]
			}
		}
	}

	return requestData
}

// Bind 实现表单数据绑定策略
func (f *FormBindingStrategy) Bind(r *http.Request, bindFieldMap map[string]BindField) (map[string]interface{}, error) {
	requestData := make(map[string]interface{})
	err := r.ParseForm()
	if err != nil {
		return nil, err
	}

	for key, values := range r.Form {
		if bindField, ok := bindFieldMap[key]; ok {
			if bindField.FieldKind == reflect.Array || bindField.FieldKind == reflect.Slice {
				requestData[bindField.BindName] = values
			} else {
				requestData[bindField.BindName] = values[0]
			}
		}
	}

	return requestData, nil
}

// func (curRT *Router) Handler(w http.ResponseWriter, r *http.Request) {

// db := DB.Session(&gorm.Session{})
// if curRT.MODEL == nil && curRT.NoAutoMigrate {
// 	log.Fatalf("%s: MODEL required when NoAutoMigrate is true", curRT.Path)
// }
// // 迁移
// if curRT.NoAutoMigrate {
// 	if err := DB.AutoMigrate(curRT.MODEL); err != nil {
// 		log.Fatalf("%s: gorm AutoMigrate err: %v", curRT.Path, err)
// 	}
// }
// // WHERE语句
// if curRT.WHERE != nil && curRT.MODEL == nil {
// 	log.Fatalf("%s: MODEL required when WHERE not nil", curRT.Path)
// }
// for query, data := range curRT.WHERE {
// 	db = db.Where(query, data)
// }

// }

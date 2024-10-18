package rt

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"

	"github.com/QingShan-Xu/web/ds"
	"github.com/go-chi/chi/v5"
)

type (

	// 数据绑定接口
	BindingStrategy interface {
		bind(r *http.Request, reader ds.Reader, bindFieldSlice []bindField) error
	}

	// 表示一个需要绑定的字段
	bindField struct {
		FieldName string
		FieldKind reflect.Kind
		BindName  string
		BindType  string
		BindTag   string
	}

	dataBinder struct{}

	queryBindingStrategy struct{}
	uriBindingStrategy   struct{}
	jsonBindingStrategy  struct{}
	formBindingStrategy  struct{}
)

func NewDataBinder() *dataBinder {
	return &dataBinder{}
}

func (binder *dataBinder) BindData(curRT *Router, r *http.Request) (interface{}, error) {
	bindData := reflect.New(reflect.TypeOf(curRT.Bind)).Interface()

	reader := ds.NewReader(bindData)

	bindFieldSlice := binder.getBindFieldSlice(reader)

	uriBind := uriBindingStrategy{}
	if err := uriBind.bind(r, reader, bindFieldSlice); err != nil {
		return nil, err
	}

	queryBind := queryBindingStrategy{}

	if err := queryBind.bind(r, reader, bindFieldSlice); err != nil {
		return nil, err
	}

	contentType := r.Header.Get("Content-Type")

	if strings.HasPrefix(contentType, "application/json") {
		jsonBind := jsonBindingStrategy{}
		jsonBind.bind(r, &bindData)
	}

	if strings.HasPrefix(contentType, "application/x-www-form-urlencoded") {
		formBind := formBindingStrategy{}
		formBind.bind(r, reader, bindFieldSlice)
	}

	validateZhInfo := ValidateStruct(bindData)
	if validateZhInfo != nil {
		var values []string
		for k, v := range validateZhInfo {
			snakeK := ToSnakeCase(k)
			values = append(values, strings.ReplaceAll(fmt.Sprintf("%v", v), k, snakeK))
		}
		validateZhErr := fmt.Errorf("%s", strings.Join(values, ", "))
		return nil, validateZhErr
	}

	return bindData, nil
}

func (binder *dataBinder) getBindFieldSlice(reader ds.Reader) []bindField {
	bindFieldslice := []bindField{}
	bindAllFields := reader.GetAllFields()

	for _, field := range bindAllFields {
		fieldTag := field.Tag()
		for _, bindFieldType := range []string{"json", "form", "uri", "query"} {
			if bindName, ok := fieldTag.Lookup(bindFieldType); ok {
				bindField := bindField{
					FieldName: field.Name(),
					FieldKind: field.Kind(),
					BindName:  bindName,
					BindType:  bindFieldType,
				}
				if bindTag, ok := fieldTag.Lookup("bind"); ok {
					bindField.BindTag = bindTag
				}
				bindFieldslice = append(bindFieldslice, bindField)
			}
		}
	}

	return bindFieldslice
}

func (q queryBindingStrategy) bind(r *http.Request, reader ds.Reader, bindFieldSlice []bindField) error {

	queryParams := r.URL.Query()

	for _, bindField := range bindFieldSlice {
		if bindField.BindType != "query" {
			continue
		}

		values, ok := queryParams[bindField.BindName]
		if !ok {
			continue
		}

		if bindField.FieldKind == reflect.Array || bindField.FieldKind == reflect.Slice {
			if err := reader.SetValue(bindField.FieldName, values); err != nil {
				return err
			}
		} else {
			if err := reader.SetValue(bindField.FieldName, values[0]); err != nil {
				return err
			}
		}

	}

	return nil
}
func (q uriBindingStrategy) bind(r *http.Request, reader ds.Reader, bindFieldSlice []bindField) error {

	for _, bindField := range bindFieldSlice {
		if bindField.BindType != "uri" {
			continue
		}

		dateParam := chi.URLParam(r, bindField.BindName)
		if dateParam == "" {
			continue
		}

		reader.SetValue(bindField.FieldName, dateParam)
	}

	return nil
}

func (q formBindingStrategy) bind(r *http.Request, reader ds.Reader, bindFieldSlice []bindField) error {
	return nil
}

func (q jsonBindingStrategy) bind(r *http.Request, bindData interface{}) error {
	if reflect.TypeOf(bindData).Kind() != reflect.Pointer {
		return fmt.Errorf("bindData must be a pointer")
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	err = json.Unmarshal(body, bindData)
	if err != nil {
		return err
	}
	return nil
}

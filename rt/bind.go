package rt

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/mitchellh/mapstructure"
)

type (
	bindImpl struct {
	}
)

func NewDataBinder() *bindImpl {
	return &bindImpl{}
}

func (bind *bindImpl) BindData(curRT *Router, r *http.Request) (interface{}, error) {
	bindData := reflect.New(reflect.TypeOf(curRT.Bind)).Interface()

	// 绑定
	if err := bind.bind(r, bindData); err != nil {
		return nil, err
	}

	if err := bind.validate(bindData); err != nil {
		return nil, err
	}

	return bindData, nil
}

// 绑定
func (binding *bindImpl) bind(r *http.Request, bindData interface{}) error {
	decoder, _ := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Squash:               true,
		WeaklyTypedInput:     true,
		TagName:              "bind",
		IgnoreUntaggedFields: true,
		Result:               bindData,
	})

	routeContext := chi.RouteContext(r.Context())
	uriMap := map[string]string{}
	for i, key := range routeContext.URLParams.Keys {
		uriMap[key] = routeContext.URLParams.Values[i]
	}
	if err := decoder.Decode(uriMap); err != nil {
		return err
	}

	// query
	if r.Method == http.MethodGet {
		querMap := map[string]interface{}{}
		queryParams := r.URL.Query()
		for queryK, queryV := range queryParams {
			if len(queryV) > 1 {
				querMap[queryK] = queryV
			} else {
				querMap[queryK] = queryV[0]
			}
		}
		if err := decoder.Decode(querMap); err != nil {
			return err
		}
	}

	contentType := r.Header.Get("Content-Type")
	// json
	if strings.HasPrefix(contentType, "application/json") && r.Method != http.MethodGet {
		jsonMap := map[string]interface{}{}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			return err
		}
		defer r.Body.Close()
		if err = json.Unmarshal(body, &jsonMap); err != nil {
			return err
		}
		if err := decoder.Decode(jsonMap); err != nil {
			return err
		}
	}

	// Form data
	if strings.HasPrefix(contentType, "application/x-www-form-urlencoded") && r.Method != http.MethodGet {
		if err := r.ParseForm(); err != nil {
			return err
		}

		formMap := make(map[string]interface{})
		for formKey, formValues := range r.PostForm {
			if len(formValues) > 1 {
				formMap[formKey] = formValues
			} else {
				formMap[formKey] = formValues[0]
			}
		}
		if err := decoder.Decode(formMap); err != nil {
			return err
		}
	}

	return nil
}

// 校验
func (binding *bindImpl) validate(bindData interface{}) error {
	validateZhInfo := ValidateStruct(bindData)
	if validateZhInfo == nil {
		return nil
	}
	var values []string
	for k, v := range validateZhInfo {
		snakeK := ToSnakeCase(k)
		values = append(values, strings.ReplaceAll(fmt.Sprintf("%v", v), k, snakeK))
	}
	validateZhErr := fmt.Errorf("%s", strings.Join(values, ", "))
	return validateZhErr
}

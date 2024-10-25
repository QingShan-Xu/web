// Package rt 提供了路由处理相关的功能，包括数据绑定、验证等。
package rt

import (
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
	"reflect"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/mitchellh/mapstructure"
)

// Binder 实现了数据绑定和验证的功能。
type Binder struct{}

// NewBinder 创建一个新的数据绑定器。
func NewBinder() *Binder {
	return &Binder{}
}

// BindAndValidate 绑定请求数据并进行验证。
// routerBind: 路由绑定的结构体类型。
// r: HTTP 请求。
// 返回绑定并验证后的数据或错误信息。
func (b *Binder) BindAndValidate(routerBind interface{}, r *http.Request) (interface{}, error) {
	if routerBind == nil {
		return nil, fmt.Errorf("Router.Bind cannot be nil")
	}

	// 创建绑定数据的实例。
	bindType := reflect.TypeOf(routerBind)
	bindValue := reflect.New(bindType).Interface()

	// 执行数据绑定。
	if err := b.bindData(r, bindValue); err != nil {
		return nil, err
	}

	// 执行数据验证。
	if err := b.validateData(bindValue); err != nil {
		return nil, err
	}

	return bindValue, nil
}

// bindData 绑定请求数据到结构体。
// r: HTTP 请求。
// bindValue: 绑定数据的实例。
func (b *Binder) bindData(r *http.Request, bindValue interface{}) error {
	// 配置 mapstructure 解码器。
	decoderConfig := &mapstructure.DecoderConfig{
		Squash:               true,
		WeaklyTypedInput:     true,
		TagName:              "bind",
		IgnoreUntaggedFields: true,
		Result:               bindValue,
	}
	decoder, _ := mapstructure.NewDecoder(decoderConfig)

	// 解析 URI 参数。
	routeCtx := chi.RouteContext(r.Context())
	uriParams := make(map[string]string, len(routeCtx.URLParams.Keys))
	for i, key := range routeCtx.URLParams.Keys {
		uriParams[key] = routeCtx.URLParams.Values[i]
	}
	if err := decoder.Decode(uriParams); err != nil {
		return fmt.Errorf("failed to decode URI parameters: %w", err)
	}

	// 解析查询参数。
	queryParams := valuesToMap(r.URL.Query())
	if err := decoder.Decode(queryParams); err != nil {
		return fmt.Errorf("failed to decode query parameters: %w", err)
	}

	// 解析请求体数据（仅针对非 GET 请求）。
	if r.Method != http.MethodGet {
		contentType := r.Header.Get("Content-Type")
		mediaType, _, err := mime.ParseMediaType(contentType)
		if err != nil {
			return fmt.Errorf("invalid Content-Type header: %w", err)
		}

		switch mediaType {
		case "application/json":
			// 解析 JSON 数据。
			bodyMap := map[string]interface{}{}
			body, err := io.ReadAll(r.Body)
			if err != nil {
				return fmt.Errorf("failed to read request body: %w", err)
			}
			defer r.Body.Close()
			if err = json.Unmarshal(body, &bodyMap); err != nil {
				return fmt.Errorf("failed to unmarshal JSON body: %w", err)
			}
			if err := decoder.Decode(bodyMap); err != nil {
				return fmt.Errorf("failed to decode JSON body: %w", err)
			}
		case "application/x-www-form-urlencoded":
			// 解析表单数据。
			if err := r.ParseForm(); err != nil {
				return fmt.Errorf("failed to parse form data: %w", err)
			}
			formMap := valuesToMap(r.PostForm)
			if err := decoder.Decode(formMap); err != nil {
				return fmt.Errorf("failed to decode form data: %w", err)
			}
		}
	}

	return nil
}

// validateData 验证绑定的数据。
// bindValue: 绑定数据的实例。
func (b *Binder) validateData(bindValue interface{}) error {
	validationErrors := ValidateStruct(bindValue)
	if validationErrors == nil {
		return nil
	}
	var errorMessages []string
	for field, errMsg := range validationErrors {
		snakeField := ToSnakeCase(field)
		formattedMsg := strings.ReplaceAll(errMsg, field, snakeField)
		errorMessages = append(errorMessages, formattedMsg)
	}
	return fmt.Errorf(strings.Join(errorMessages, ", "))
}

// valuesToMap 将 url.Values 转换为 map[string]interface{}。
// values: URL 参数或表单参数。
// 返回转换后的 map。
func valuesToMap(values url.Values) map[string]interface{} {
	result := make(map[string]interface{}, len(values))
	for key, vals := range values {
		if len(vals) > 1 {
			result[key] = vals
		} else {
			result[key] = vals[0]
		}
	}
	return result
}

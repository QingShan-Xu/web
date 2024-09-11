package utils

import (
	"reflect"
	"unicode"
)

func GetInstanceVal(anyData interface{}) reflect.Value {
	instanceTpe := reflect.TypeOf(anyData)
	instanceVal := reflect.ValueOf(anyData)

	if instanceTpe.Kind() == reflect.Pointer {
		instanceVal = instanceVal.Elem()
	}

	return instanceVal
}

func ToSnakeCase(str string) string {
	var result []rune
	for i, r := range str {
		if unicode.IsUpper(r) {
			if i > 0 && (!unicode.IsUpper(r) || (i+1 < len(str) && unicode.IsLower(rune(str[i+1])))) {
				result = append(result, '_')
			}
			result = append(result, unicode.ToLower(r))
		} else {
			result = append(result, r)
		}
	}
	return string(result)
}

func IsIncludes[T comparable](slice []T, item T) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

func IsBasicType(tpe reflect.Type) bool {
	if tpe.Kind() == reflect.Pointer {
		tpe = tpe.Elem()
	}

	switch tpe.Kind() {
	case reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64,
		reflect.Complex64, reflect.Complex128,
		reflect.String:
		return true
	default:
		return false
	}
}

// Struct2map 将结构体转换为 map，并且可以选择是否保留零值字段（布尔值除外）
func Struct2map(s interface{}, keepZeroValues bool) map[string]interface{} {
	data := make(map[string]interface{})
	val := reflect.ValueOf(s)

	// 处理指针的情况
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	// 确保 val 是结构体
	if val.Kind() != reflect.Struct {
		return data
	}

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldName := val.Type().Field(i).Name
		fieldValue := field.Interface()

		// 检查字段是否是结构体
		if field.Kind() == reflect.Struct {
			// 递归处理嵌套结构体
			innerMap := Struct2map(fieldValue, keepZeroValues)
			if len(innerMap) > 0 {
				data[fieldName] = innerMap
			}
		} else if field.Kind() == reflect.Bool || keepZeroValues || !isZero(field) {
			// 保留布尔值或根据参数决定是否保留零值
			data[fieldName] = fieldValue
		}
	}
	return data
}

func isZero(val reflect.Value) bool {
	zero := reflect.Zero(val.Type()).Interface()
	return reflect.DeepEqual(val.Interface(), zero)
}

func MapFlatten(m map[string]interface{}) map[string]interface{} {
	flatMap := make(map[string]interface{})
	flatten("", m, flatMap)
	return flatMap
}

func flatten(prefix string, m map[string]interface{}, flatMap map[string]interface{}) {
	for k, v := range m {
		fullKey := k
		if prefix != "" {
			fullKey = prefix + "." + k
		}

		// 判断是否是嵌套的 map
		if nestedMap, ok := v.(map[string]interface{}); ok {
			flatten(fullKey, nestedMap, flatMap)
		} else {
			flatMap[fullKey] = v
		}
	}
}

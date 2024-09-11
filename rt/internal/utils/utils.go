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
			// 如果不是字符串开始并且前一个字符不是下划线
			if i > 0 && !(unicode.IsUpper(rune(str[i-1]))) {
				result = append(result, '_') // 添加下划线
			}
			result = append(result, unicode.ToLower(r)) // 将大写字母转换为小写并添加到结果中
		} else {
			result = append(result, r) // 如果是小写字母或其他字符，直接添加到结果中
		}
	}
	return string(result) // 将 rune 数组转换为字符串并返回
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

		// 只处理导出字段
		if field.CanInterface() {
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
	}
	return data
}

// 判断是否为零值
func isZero(v reflect.Value) bool {
	// 根据类型判断是否为零值
	switch v.Kind() {
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Complex64, reflect.Complex128:
		return v.Complex() == 0
	case reflect.Array, reflect.Slice:
		return v.Len() == 0
	case reflect.Map:
		return v.IsNil() || v.Len() == 0
	case reflect.String:
		return v.String() == ""
	}
	return false
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

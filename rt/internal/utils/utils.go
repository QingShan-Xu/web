package utils

import (
	"fmt"
	"reflect"
	"strings"
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
func Struct2map(s interface{}) map[string]interface{} {
	data := make(map[string]interface{})
	sVal := reflect.ValueOf(s)
	sTpe := reflect.TypeOf(s)
	fmt.Printf("%+v\n", s)

	// 处理指针的情况
	if sTpe.Kind() == reflect.Ptr {
		sVal = sVal.Elem()
		sTpe = sTpe.Elem()
	}

	// 确保 val 是结构体
	if sTpe.Kind() != reflect.Struct {
		return data
	}

	for i := 0; i < sVal.NumField(); i++ {
		fieldVal := sVal.Field(i)
		fieldKind := fieldVal.Kind()
		fieldName := fieldVal.Type().Name()

		if !fieldVal.CanInterface() {
			continue
		}
		fieldValue := fieldVal.Interface()

		if fieldKind == reflect.Pointer {
			fieldValue = fieldVal.Elem().Interface()
			fieldVal = fieldVal.Elem()
			fieldKind = fieldVal.Kind()
		}

		// 检查字段是否是结构体
		if fieldKind == reflect.Struct {
			innerMap := Struct2map(fieldValue)
			if len(innerMap) > 0 {
				data[fieldName] = innerMap
			}
		}

		if fieldValue != nil {
			data[fieldName] = fieldValue
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

type DynamicStruct struct {
	Value reflect.Value
}

// GetField 根据路径获取字段或键值
func (ds *DynamicStruct) GetField(path string) (interface{}, error) {
	parts := strings.Split(path, ".")
	val := ds.Value

	for _, part := range parts {
		switch val.Kind() {
		case reflect.Pointer:
			val = val.Elem()
		case reflect.Map:
			val = val.MapIndex(reflect.ValueOf(part))
		case reflect.Struct:
			val = getFieldByNameOrEmbedded(val, part)
		case reflect.Slice, reflect.Array:
			index := -1
			_, err := fmt.Sscanf(part, "%d", &index)
			if err != nil || index < 0 || index >= val.Len() {
				return nil, fmt.Errorf("path %s: index %s out of range", path, part)
			}
			val = val.Index(index)
		default:
			return nil, fmt.Errorf("path %s: unexpected type %s", path, val.Kind())
		}
		if !val.IsValid() {
			return nil, fmt.Errorf("path %s: field or key %s not found", path, part)
		}
	}

	if val.Kind() == reflect.Pointer {
		if val.IsNil() {
			return nil, nil
		}
		val = val.Elem()
	}

	return val.Interface(), nil
}

// getFieldByNameOrEmbedded 检查嵌入结构体的字段并访问
func getFieldByNameOrEmbedded(val reflect.Value, fieldName string) reflect.Value {
	field := val.FieldByName(fieldName)
	if field.IsValid() {
		return field
	}

	// 如果是嵌入结构体，遍历所有字段进行嵌入检查
	for i := 0; i < val.NumField(); i++ {
		fieldType := val.Type().Field(i)
		if fieldType.Anonymous { // 这是嵌入结构体
			embeddedField := val.Field(i).FieldByName(fieldName)
			if embeddedField.IsValid() {
				return embeddedField
			}
		}
	}
	return reflect.Value{}
}

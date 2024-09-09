package rt

import (
	"reflect"
	"unicode"
)

func getInstanceVal(anyData interface{}) reflect.Value {
	instanceTpe := reflect.TypeOf(anyData)
	instanceVal := reflect.ValueOf(anyData)

	if instanceTpe.Kind() == reflect.Pointer {
		instanceVal = instanceVal.Elem()
	}

	return instanceVal
}

// func IsEmptyStruct(s interface{}) bool {
// 	v := getInstanceVal(s)

// 	if v.Kind() != reflect.Struct {
// 		return false
// 	}

// 	for i := 0; i < v.NumField(); i++ {
// 		field := v.Field(i)

// 		if field.Kind() == reflect.Struct {
// 			return IsEmptyStruct(field.Interface())
// 		} else if !field.IsZero() {
// 			return false
// 		}
// 	}

// 	return true
// }

func isIncludes[T comparable](slice []T, item T) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

func toSnakeCase(str string) string {
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

func struct2map(s interface{}) (data map[string]interface{}) {
	data = make(map[string]interface{})
	if s == nil {
		return
	}
	sVal := getInstanceVal(s)
	sTpe := sVal.Type()
	if sTpe.Kind() != reflect.Struct {
		return
	}

	for i := 0; i < sVal.NumField(); i++ {
		fName := sTpe.Field(i).Name
		fValue := sVal.Field(i).Interface()
		data[fName] = fValue
	}
	return
}

func isGroupRouter(router *Router) bool {
	return len(router.Children) != 0
}

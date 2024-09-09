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

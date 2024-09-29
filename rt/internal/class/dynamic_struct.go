package class

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
)

type DynamicStruct struct {
	Value reflect.Value
	Ctx   *gin.Context
}

// GetField 根据路径获取字段或键值
func (ds *DynamicStruct) GetField(path string) (interface{}, error) {
	parts := strings.Split(path, ".")
	val := ds.Value

	if strings.HasPrefix(path, "@") && ds.Ctx != nil {
		ginData, isExist := ds.Ctx.Get(path[1:])
		if !isExist {
			return nil, nil
		}
		return ginData, nil
	}

	for i, part := range parts {
		val = dereferencePointer(val)

		if !val.IsValid() {
			break
		}

		// 检查是否是 $len 操作
		if part == "$len" {
			if i != len(parts)-1 {
				return nil, fmt.Errorf("$len 必须是路径的最后一部分")
			}
			switch val.Kind() {
			case reflect.Slice, reflect.Array, reflect.Map, reflect.String:
				return val.Len(), nil
			default:
				return nil, fmt.Errorf("无法获取类型 %s 的长度", val.Kind())
			}
		}

		// 处理 $int 操作
		if part == "$int" {
			switch val.Kind() {
			case reflect.Slice, reflect.Array:
				result := make([]interface{}, val.Len())
				for j := 0; j < val.Len(); j++ {
					item := val.Index(j)
					if i == len(parts)-1 {
						// 如果 $int 是最后一部分，直接返回元素
						result[j] = item.Interface()
					} else {
						// 如果 $int 后面还有路径，继续处理
						itemDS := &DynamicStruct{Value: item}
						fieldValue, err := itemDS.GetField(strings.Join(parts[i+1:], "."))
						if err != nil {
							return nil, fmt.Errorf("处理切片/数组的第 %d 个元素时出错: %v", j, err)
						}
						result[j] = fieldValue
					}
				}
				return result, nil
			default:
				return nil, fmt.Errorf("类型 %s 不是数组或切片", val.Kind())
			}
		}

		// 根据类型处理字段或键的访问
		switch val.Kind() {
		case reflect.Map:
			val = val.MapIndex(reflect.ValueOf(part))
		case reflect.Struct:
			val = getFieldByNameOrEmbedded(val, part)
		case reflect.Slice, reflect.Array:
			index, err := parseIndex(part, path)
			if err != nil {
				return nil, err
			}
			if index < 0 || index >= val.Len() {
				return nil, fmt.Errorf("路径 %s: 索引 %d 超出范围", path, index)
			}
			val = val.Index(index)
		default:
			return nil, fmt.Errorf("路径 %s: 类型 %s 无法处理", path, val.Kind())
		}

		if !val.IsValid() {
			return nil, fmt.Errorf("路径 %s: 字段 %s 不存在", path, part)
		}
	}

	// 如果是 nil 则返回 nil
	if val.Kind() == reflect.Invalid {
		return nil, nil
	}

	isNil := checkNil(val)
	if isNil {
		return nil, nil
	}

	return val.Interface(), nil
}

// dereferencePointer 解引用指针
func dereferencePointer(val reflect.Value) reflect.Value {
	if val.Kind() == reflect.Pointer && !val.IsNil() {
		return val.Elem()
	}
	return reflect.ValueOf(nil)
}

// parseIndex 解析字符串为索引
func parseIndex(part, path string) (int, error) {
	index := -1
	if _, err := fmt.Sscanf(part, "%d", &index); err != nil {
		return -1, fmt.Errorf("路径 %s: 无效的索引 %s", path, part)
	}
	return index, nil
}

// getFieldByNameOrEmbedded 检查嵌入结构体的字段并访问
func getFieldByNameOrEmbedded(val reflect.Value, fieldName string) reflect.Value {
	// 先查找直接字段
	if field := val.FieldByName(fieldName); field.IsValid() {
		return field
	}

	// 如果是嵌入结构体，递归查找嵌入字段
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

func checkNil(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Ptr, reflect.Slice, reflect.Map, reflect.Chan, reflect.Func, reflect.Interface:
		return v.IsNil()
	default:
		return false
	}
}

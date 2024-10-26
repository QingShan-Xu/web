// 字段相关的定义和方法

package ds

import (
	"fmt"
	"reflect"
)

// sourceField 包含了字段的元信息。
type sourceField struct {
	reflect.StructField
}

// NewFieldReader 创建一个新的字段读取器。
// field: 结构体字段。
// fieldValue: 字段的反射值。
// 返回 StructReader 实例或错误信息。
func NewFieldReader(field reflect.StructField, fieldValue reflect.Value) (*structReader, error) {
	fieldReader := &structReader{
		Name:        field.Name,
		FieldType:   fieldValue.Type(),
		FieldKind:   fieldValue.Kind(),
		FieldValue:  fieldValue,
		Fields:      []*structReader{}, // 初始化为空
		FieldSource: &sourceField{field},
		Tags:        make(map[string]Tag),
	}

	// 解析字段标签
	fieldReader.parseFieldTags()

	// 处理嵌入字段，但不递归解析
	if field.Anonymous && fieldValue.Kind() == reflect.Struct {
		fieldReader.IsEmbedded = true
		// 不递归解析嵌入字段的子字段
	}

	// 检查是否为数组或切片
	if fieldValue.Kind() == reflect.Slice || fieldValue.Kind() == reflect.Array {
		fieldReader.IsList = true
	}

	return fieldReader, nil
}

// mapArray 对数组或切片进行映射，提取指定的键。
// keys: 需要提取的键列表。
// 返回包含提取结果的新切片或错误信息。
func (r *structReader) mapArray(keys []string) ([]map[string]interface{}, error) {
	if !r.IsList {
		return nil, fmt.Errorf("MapArray: the field '%s' is not a list", r.Name)
	}

	var result []map[string]interface{}
	for i := 0; i < r.FieldValue.Len(); i++ {
		itemValue := r.FieldValue.Index(i)
		itemReader, err := generateStructReader(itemValue)
		if err != nil {
			return nil, err
		}
		itemMap := make(map[string]interface{})
		for _, key := range keys {
			fieldReader, err := itemReader.GetField(key)
			if err != nil {
				return nil, err
			}
			itemMap[key] = fieldReader.Interface()
		}
		result = append(result, itemMap)
	}

	return result, nil
}

// 核心的 StructReader 定义和方法
// package ds 提供了结构体数据的读取和操作功能。
package ds

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

type (
	FieldReader interface {
		len() (int, error)
		parseFieldTags()
		mapArray(keys []string) ([]map[string]interface{}, error)

		GetField(name string) (FieldReader, error)
		GetName() string
		GetTag() map[string]Tag
		SafePointerInt() (*int, bool)
		SafeInt() (int, bool)
		SafePointerInt8() (*int8, bool)
		SafeInt8() (int8, bool)
		SafePointerInt16() (*int16, bool)
		SafeInt16() (int16, bool)
		SafePointerInt32() (*int32, bool)
		SafeInt32() (int32, bool)
		SafePointerInt64() (*int64, bool)
		SafeInt64() (int64, bool)
		SafePointerUint() (*uint, bool)
		SafeUint() (uint, bool)
		SafePointerUint8() (*uint8, bool)
		SafeUint8() (uint8, bool)
		SafePointerUint16() (*uint16, bool)
		SafeUint16() (uint16, bool)
		SafePointerUint32() (*uint32, bool)
		SafeUint32() (uint32, bool)
		SafePointerUint64() (*uint64, bool)
		SafeUint64() (uint64, bool)
		SafePointerFloat32() (*float32, bool)
		SafeFloat32() (float32, bool)
		SafePointerFloat64() (*float64, bool)
		SafeFloat64() (float64, bool)
		SafePointerString() (*string, bool)
		SafeString() (string, bool)
		SafePointerBool() (*bool, bool)
		SafeBool() (bool, bool)
		SafePointerTime() (*time.Time, bool)
		SafeTime() (time.Time, bool)
		Interface() interface{}
	}

	// structReader 用于读取结构体的信息，支持嵌套、匿名字段和标签解析。
	structReader struct {
		Name       string          // 字段名称
		FieldType  reflect.Type    // 字段类型
		FieldKind  reflect.Kind    // 字段种类
		FieldValue reflect.Value   // 字段值
		Fields     []*structReader // 子字段列表

		IsEmbedded   bool // 是否为嵌入字段
		IsList       bool // 是否为数组或切片类型
		fieldsParsed bool

		FieldSource *sourceField   // 原始的结构体字段信息
		Tags        map[string]Tag // 字段标签
	}
)

// NewStructReader 创建一个新的 StructReader。
// data: 输入的结构体实例。
// 返回 StructReader 实例或错误信息。
func NewStructReader(data interface{}) (FieldReader, error) {
	if data == nil {
		return nil, fmt.Errorf("newStructReader: data cannot be nil")
	}

	value := reflect.ValueOf(data)
	if value.Kind() == reflect.Ptr {
		if value.IsNil() {
			return nil, fmt.Errorf("newStructReader: data cannot be nil pointer")
		}
		value = value.Elem()
	}

	if value.Kind() != reflect.Struct {
		return nil, fmt.Errorf("NewStructReader: expected struct type but got %s", value.Kind())
	}

	reader, err := generateStructReader(value)
	if err != nil {
		return nil, err
	}

	return reader, nil
}

// generateStructReader 生成 StructReader。
// value: 反射值。
// 返回 StructReader 实例或错误信息。
func generateStructReader(value reflect.Value) (FieldReader, error) {
	structType := value.Type()

	reader := &structReader{
		Name:       structType.Name(),
		FieldType:  structType,
		FieldKind:  value.Kind(),
		FieldValue: value,
		Fields:     []*structReader{}, // 初始化为空
		Tags:       make(map[string]Tag),
	}

	if value.Kind() != reflect.Struct {
		return reader, nil
	}

	for i := 0; i < value.NumField(); i++ {
		fieldValue := value.Field(i)
		field := structType.Field(i)

		// 跳过未导出的字段
		if !field.IsExported() {
			continue
		}

		// 仅解析当前字段，不递归嵌套字段
		fieldReader, err := NewFieldReader(field, fieldValue)
		if err != nil {
			return nil, err
		}

		reader.Fields = append(reader.Fields, fieldReader)
	}

	return reader, nil
}

// GetField 根据名称获取字段。
// name: 字段名称，支持嵌套（例如 "A.B.C"）、特殊方法（例如 "items._len"）、数组映射（例如 "items[key1,key2]"）。
// 返回 StructReader 实例或错误信息。
func (r *structReader) GetField(name string) (FieldReader, error) {
	if name == "" {
		return nil, fmt.Errorf("GetField: name cannot be empty")
	}

	// 处理特殊语法
	if strings.HasSuffix(name, "._len") {
		fieldName := strings.TrimSuffix(name, "._len")
		fieldReader, err := r.GetField(fieldName)
		if err != nil {
			return nil, err
		}
		length, err := fieldReader.len()
		if err != nil {
			return nil, err
		}

		// 返回长度值的 StructReader
		return &structReader{
			Name:       "_len",
			FieldType:  reflect.TypeOf(length),
			FieldKind:  reflect.Int,
			FieldValue: reflect.ValueOf(length),
		}, nil
	}

	// 处理数组映射语法，例如 "items[key1,key2]"
	if strings.Contains(name, "[") && strings.HasSuffix(name, "]") {
		index := strings.Index(name, "[")
		fieldName := name[:index]
		keysStr := name[index+1 : len(name)-1]
		keys := strings.Split(keysStr, ",")

		fieldReader, err := r.GetField(fieldName)
		if err != nil {
			return nil, err
		}

		mappedArray, err := fieldReader.mapArray(keys)
		if err != nil {
			return nil, err
		}

		// 返回映射后的数组的 StructReader
		return &structReader{
			Name:       fieldReader.GetName(),
			FieldType:  reflect.TypeOf(mappedArray),
			FieldKind:  reflect.Slice,
			FieldValue: reflect.ValueOf(mappedArray),
		}, nil
	}

	// 正常的字段访问，支持嵌套字段
	nameComponents := parseFieldName(name)
	currentReader := r

	for _, component := range nameComponents {
		found := false
		for _, field := range currentReader.Fields {
			if field.Name == component {
				// 如果字段是结构体且未解析子字段，则动态解析
				if (field.FieldKind == reflect.Struct || (field.IsEmbedded && field.FieldKind == reflect.Struct)) && !field.fieldsParsed {
					embeddedReader, err := generateStructReader(field.FieldValue)
					if err != nil {
						return nil, err
					}
					if sr, ok := embeddedReader.(*structReader); ok {
						field.Fields = sr.Fields
						// 不要覆盖 field.Tags
						field.fieldsParsed = true
					} else {
						return nil, fmt.Errorf("unexpected type for embeddedReader")
					}
				}

				currentReader = field
				found = true
				break
			}

			// 如果是嵌入字段，继续在其字段中查找
			if field.IsEmbedded {
				if field.FieldKind == reflect.Struct && !field.fieldsParsed {
					embeddedReader, err := generateStructReader(field.FieldValue)
					if err != nil {
						return nil, err
					}
					if sr, ok := embeddedReader.(*structReader); ok {
						field.Fields = sr.Fields
						// 不要覆盖 field.Tags
						field.fieldsParsed = true
					} else {
						return nil, fmt.Errorf("unexpected type for embeddedReader")
					}
				}

				embeddedField, err := field.GetField(component)
				if err == nil {
					currentReader = embeddedField.(*structReader)
					found = true
					break
				}
			}
		}
		if !found {
			return nil, fmt.Errorf("field '%s' not found", name)
		}
	}

	return currentReader, nil
}

// len 获取支持 len() 方法的值的长度。
// 返回长度或错误信息。
func (r *structReader) len() (int, error) {
	switch r.FieldKind {
	case reflect.Array, reflect.Slice, reflect.Map, reflect.Chan, reflect.String:
		return r.FieldValue.Len(), nil
	default:
		return 0, fmt.Errorf("len: unsupported kind %s", r.FieldKind)
	}
}

// len 获取支持 len() 方法的值的长度。
// 返回长度或错误信息。
func (r *structReader) GetName() string {
	return r.Name
}

// SafePointerInt 尝试从 reflect.Value 中获取 *int 类型的指针。
// 如果成功，返回指针和 true；否则，返回 nil 和 false。
func (r *structReader) SafePointerInt() (*int, bool) {
	if r.FieldValue.IsNil() {
		return nil, false
	}
	value, ok := r.SafeInt()
	if ok {
		return &value, true
	}
	return nil, false
}

// SafeInt 尝试从 reflect.Value 中获取 int 类型的值。
// 如果成功，返回整数值和 true；否则，返回 0 和 false。
func (r *structReader) SafeInt() (int, bool) {
	v := reflect.Indirect(r.FieldValue)
	if v.CanInt() {
		return int(v.Int()), true
	}
	return 0, false
}

// SafePointerInt8 尝试从 reflect.Value 中获取 *int8 类型的指针。
// 如果成功，返回指针和 true；否则，返回 nil 和 false。
func (r *structReader) SafePointerInt8() (*int8, bool) {
	if r.FieldValue.IsNil() {
		return nil, false
	}
	value, ok := r.SafeInt8()
	if ok {
		return &value, true
	}
	return nil, false
}

// SafeInt8 尝试从 reflect.Value 中获取 int8 类型的值。
// 如果成功，返回 int8 值和 true；否则，返回 0 和 false。
func (r *structReader) SafeInt8() (int8, bool) {
	value, ok := r.SafeInt()
	if ok {
		value := int8(value)
		return value, true
	}
	return 0, false
}

// SafePointerInt16 尝试从 reflect.Value 中获取 *int16 类型的指针。
// 如果成功，返回指针和 true；否则，返回 nil 和 false。
func (r *structReader) SafePointerInt16() (*int16, bool) {
	if r.FieldValue.IsNil() {
		return nil, false
	}
	value, ok := r.SafeInt16()
	if ok {
		return &value, true
	}
	return nil, false
}

// SafeInt16 尝试从 reflect.Value 中获取 int16 类型的值。
// 如果成功，返回 int16 值和 true；否则，返回 0 和 false。
func (r *structReader) SafeInt16() (int16, bool) {
	value, ok := r.SafeInt()
	if ok {
		value := int16(value)
		return value, true
	}
	return 0, false
}

// SafePointerInt32 尝试从 reflect.Value 中获取 *int32 类型的指针。
// 如果成功，返回指针和 true；否则，返回 nil 和 false。
func (r *structReader) SafePointerInt32() (*int32, bool) {
	if r.FieldValue.IsNil() {
		return nil, false
	}
	value, ok := r.SafeInt32()
	if ok {
		return &value, true
	}
	return nil, false
}

// SafeInt32 尝试从 reflect.Value 中获取 int32 类型的值。
// 如果成功，返回 int32 值和 true；否则，返回 0 和 false。
func (r *structReader) SafeInt32() (int32, bool) {
	value, ok := r.SafeInt()
	if ok {
		value := int32(value)
		return value, true
	}
	return 0, false
}

// SafePointerInt64 尝试从 reflect.Value 中获取 *int64 类型的指针。
// 如果成功，返回指针和 true；否则，返回 nil 和 false。
func (r *structReader) SafePointerInt64() (*int64, bool) {
	if r.FieldValue.IsNil() {
		return nil, false
	}
	value, ok := r.SafeInt64()
	if ok {
		return &value, true
	}
	return nil, false
}

// SafeInt64 尝试从 reflect.Value 中获取 int64 类型的值。
// 如果成功，返回 int64 值和 true；否则，返回 0 和 false。
func (r *structReader) SafeInt64() (int64, bool) {
	value, ok := r.SafeInt()
	if ok {
		value := int64(value)
		return value, true
	}
	return 0, false
}

// SafeUint 尝试从 reflect.Value 中获取 uint 类型的值。
// 如果成功，返回 uint 值和 true；否则，返回 0 和 false。
func (r *structReader) SafeUint() (uint, bool) {
	v := reflect.Indirect(r.FieldValue)
	if v.CanUint() {
		return uint(v.Uint()), true
	}
	return 0, false
}

// SafePointerUint 尝试从 reflect.Value 中获取 *uint 类型的指针。
// 如果成功，返回指针和 true；否则，返回 nil 和 false。
func (r *structReader) SafePointerUint() (*uint, bool) {
	if r.FieldValue.IsNil() {
		return nil, false
	}
	value, ok := r.SafeUint()
	if ok {
		return &value, true
	}
	return nil, false
}

// SafeUint8 尝试从 reflect.Value 中获取 uint8 类型的值。
// 如果成功，返回 uint8 值和 true；否则，返回 0 和 false。
func (r *structReader) SafeUint8() (uint8, bool) {
	value, ok := r.SafeUint()
	if ok {
		value := uint8(value)
		return value, true
	}
	return 0, false
}

// SafePointerUint8 尝试从 reflect.Value 中获取 *uint8 类型的指针。
// 如果成功，返回指针和 true；否则，返回 nil 和 false。
func (r *structReader) SafePointerUint8() (*uint8, bool) {
	if r.FieldValue.IsNil() {
		return nil, false
	}
	value, ok := r.SafeUint8()
	if ok {
		return &value, true
	}
	return nil, false
}

// SafeUint16 尝试从 reflect.Value 中获取 uint16 类型的值。
// 如果成功，返回 uint16 值和 true；否则，返回 0 和 false。
func (r *structReader) SafeUint16() (uint16, bool) {
	value, ok := r.SafeUint()
	if ok {
		value := uint16(value)
		return value, true
	}
	return 0, false
}

// SafePointerUint16 尝试从 reflect.Value 中获取 *uint16 类型的指针。
// 如果成功，返回指针和 true；否则，返回 nil 和 false。
func (r *structReader) SafePointerUint16() (*uint16, bool) {
	if r.FieldValue.IsNil() {
		return nil, false
	}
	value, ok := r.SafeUint16()
	if ok {
		return &value, true
	}
	return nil, false
}

// SafeUint32 尝试从 reflect.Value 中获取 uint32 类型的值。
// 如果成功，返回 uint32 值和 true；否则，返回 0 和 false。
func (r *structReader) SafeUint32() (uint32, bool) {
	value, ok := r.SafeUint()
	if ok {
		value := uint32(value)
		return value, true
	}
	return 0, false
}

// SafePointerUint32 尝试从 reflect.Value 中获取 *uint32 类型的指针。
// 如果成功，返回指针和 true；否则，返回 nil 和 false。
func (r *structReader) SafePointerUint32() (*uint32, bool) {
	if r.FieldValue.IsNil() {
		return nil, false
	}
	value, ok := r.SafeUint32()
	if ok {
		return &value, true
	}
	return nil, false
}

// SafeUint64 尝试从 reflect.Value 中获取 uint64 类型的值。
// 如果成功，返回 uint64 值和 true；否则，返回 0 和 false。
func (r *structReader) SafeUint64() (uint64, bool) {
	value, ok := r.SafeUint()
	if ok {
		value := uint64(value)
		return value, true
	}
	return 0, false
}

// SafePointerUint64 尝试从 reflect.Value 中获取 *uint64 类型的指针。
// 如果成功，返回指针和 true；否则，返回 nil 和 false。
func (r *structReader) SafePointerUint64() (*uint64, bool) {
	if r.FieldValue.IsNil() {
		return nil, false
	}
	value, ok := r.SafeUint64()
	if ok {
		return &value, true
	}
	return nil, false
}

// SafeFloat32 尝试从 reflect.Value 中获取 float32 类型的值。
// 如果成功，返回 float32 值和 true；否则，返回 0.0 和 false。
func (r *structReader) SafeFloat32() (float32, bool) {
	v := reflect.Indirect(r.FieldValue)
	if v.CanFloat() {
		return float32(v.Float()), true
	}
	return 0, false
}

// SafePointerFloat32 尝试从 reflect.Value 中获取 *float32 类型的指针。
// 如果成功，返回指针和 true；否则，返回 nil 和 false。
func (r *structReader) SafePointerFloat32() (*float32, bool) {
	if r.FieldValue.IsNil() {
		return nil, false
	}
	value, ok := r.SafeFloat32()
	if ok {
		return &value, true
	}
	return nil, false
}

// SafeFloat64 尝试从 reflect.Value 中获取 float64 类型的值。
// 如果成功，返回 float64 值和 true；否则，返回 0.0 和 false。
func (r *structReader) SafeFloat64() (float64, bool) {
	v := reflect.Indirect(r.FieldValue)
	if v.CanFloat() {
		return v.Float(), true
	}
	return 0, false
}

// SafePointerFloat64 尝试从 reflect.Value 中获取 *float64 类型的指针。
// 如果成功，返回指针和 true；否则，返回 nil 和 false。
func (r *structReader) SafePointerFloat64() (*float64, bool) {
	if r.FieldValue.IsNil() {
		return nil, false
	}
	value, ok := r.SafeFloat64()
	if ok {
		return &value, true
	}
	return nil, false
}

// SafeString 尝试从 reflect.Value 中获取 string 类型的值。
// 如果成功，返回字符串和 true；否则，返回空字符串和 false。
func (r *structReader) SafeString() (string, bool) {
	v := reflect.Indirect(r.FieldValue)
	if v.Kind() == reflect.String {
		return v.String(), true
	}
	return "", false
}

// SafePointerString 尝试从 reflect.Value 中获取 *string 类型的指针。
// 如果成功，返回指针和 true；否则，返回 nil 和 false。
func (r *structReader) SafePointerString() (*string, bool) {
	if r.FieldValue.IsNil() {
		return nil, false
	}
	value, ok := r.SafeString()
	if ok {
		return &value, true
	}
	return nil, false
}

// SafeBool 尝试从 reflect.Value 中获取 bool 类型的值。
// 如果成功，返回 bool 值和 true；否则，返回 false 和 false。
func (r *structReader) SafeBool() (bool, bool) {
	v := reflect.Indirect(r.FieldValue)
	if v.Kind() == reflect.Bool {
		return v.Bool(), true
	}
	return false, false
}

// SafePointerBool 尝试从 reflect.Value 中获取 *bool 类型的指针。
// 如果成功，返回指针和 true；否则，返回 nil 和 false。
func (r *structReader) SafePointerBool() (*bool, bool) {
	if r.FieldValue.IsNil() {
		return nil, false
	}
	value, ok := r.SafeBool()
	if ok {
		return &value, true
	}
	return nil, false
}

// SafeTime 尝试从 reflect.Value 中获取 time.Time 类型的值。
// 如果成功，返回 time.Time 值和 true；否则，返回 time.Time{} 和 false。
func (r *structReader) SafeTime() (time.Time, bool) {
	v := reflect.Indirect(r.FieldValue)
	value, ok := v.Interface().(time.Time)
	if ok {
		return value, true
	}
	return time.Time{}, false
}

// SafePointerTime 尝试从 reflect.Value 中获取 *time.Time 类型的指针。
// 如果成功，返回指针和 true；否则，返回 nil 和 false。
func (r *structReader) SafePointerTime() (*time.Time, bool) {
	if r.FieldValue.IsNil() {
		return nil, false
	}
	value, ok := r.SafeTime()
	if ok {
		return &value, true
	}
	return nil, false
}

func (r *structReader) GetTag() map[string]Tag {
	return r.Tags
}

func (r *structReader) Interface() interface{} {
	return r.FieldValue.Interface()
}

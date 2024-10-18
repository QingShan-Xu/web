package ds

import (
	"fmt"
	"reflect"
	"time"

	"github.com/QingShan-Xu/web/bm"
)

type (
	Reader interface {
		HasField(name string) bool
		GetField(name string) Field
		GetAllFields() []Field
		GetValue() interface{}
		SetValue(name string, newVal interface{}) error
	}

	Field interface {
		Name() string
		Tag() reflect.StructTag
		Kind() reflect.Kind
		PointerInt() *int
		Int() int
		PointerInt8() *int8
		Int8() int8
		PointerInt16() *int16
		Int16() int16
		PointerInt32() *int32
		Int32() int32
		PointerInt64() *int64
		Int64() int64
		PointerUint() *uint
		Uint() uint
		PointerUint8() *uint8
		Uint8() uint8
		PointerUint16() *uint16
		Uint16() uint16
		PointerUint32() *uint32
		Uint32() uint32
		PointerUint64() *uint64
		Uint64() uint64
		PointerFloat32() *float32
		Float32() float32
		PointerFloat64() *float64
		Float64() float64
		PointerString() *string
		String() string
		PointerBool() *bool
		Bool() bool
		PointerTime() *time.Time
		Time() time.Time
		Interface() interface{}
	}

	readImpl struct {
		fields map[string]fieldImpl
		value  interface{}
	}

	fieldImpl struct {
		field reflect.StructField
		value reflect.Value
	}
)

// 传入指针
func NewReader(value interface{}) Reader {
	if reflect.ValueOf(value).Kind() != reflect.Ptr {
		panic("value must be pointer")
	}

	fields := map[string]fieldImpl{}

	valueOf := reflect.Indirect(reflect.ValueOf(value))
	typeOf := valueOf.Type()

	if typeOf.Kind() == reflect.Struct {
		for i := 0; i < valueOf.NumField(); i++ {
			field := typeOf.Field(i)
			value := valueOf.Field(i)
			if field.Anonymous {
				for i := 0; i < value.NumField(); i++ {
					currentField := field.Type.Field(i)
					fields[currentField.Name] = fieldImpl{
						field: currentField,
						value: value.Field(i),
					}
				}
			} else {
				fields[field.Name] = fieldImpl{
					field: field,
					value: valueOf.Field(i),
				}
			}
		}
	}

	return readImpl{
		fields: fields,
		value:  value,
	}
}

func (r readImpl) HasField(name string) bool {
	_, ok := r.fields[name]
	return ok
}

func (r readImpl) GetField(name string) Field {
	if !r.HasField(name) {
		return nil
	}
	return r.fields[name]
}

func (r readImpl) GetAllFields() []Field {
	var fields []Field

	for _, field := range r.fields {
		fields = append(fields, field)
	}

	return fields
}

func (r readImpl) GetValue() interface{} {
	return r.value
}

func (r readImpl) SetValue(name string, newVal interface{}) error {
	if !r.HasField(name) {
		return fmt.Errorf("field %s does not exist", name)
	}

	field := r.fields[name]
	value := reflect.ValueOf(newVal)

	// 特殊处理 FlexibleInt
	if flexInt, ok := newVal.(bm.FlexibleInt); ok {
		value = reflect.ValueOf(flexInt.Int())
	}

	// 转换新值类型以匹配字段类型
	if value.Type() != field.value.Type() {
		if !value.Type().ConvertibleTo(field.value.Type()) {
			return fmt.Errorf("%s 类型错误: 需要 %s, 获得 %s", name, field.value.Type(), value.Type())
		}
		value = value.Convert(field.value.Type())
	}

	if field.value.Kind() == reflect.Ptr {
		if field.value.IsNil() {
			field.value.Set(reflect.New(field.value.Type().Elem()))
		}
		if value.Kind() != reflect.Ptr {
			// 如果新值不是指针，创建一个新的指针并设置其值
			ptr := reflect.New(field.value.Type().Elem())
			ptr.Elem().Set(value)
			value = ptr
		}
		field.value.Set(value)
	} else if field.value.CanSet() {
		field.value.Set(value)
	} else {
		return fmt.Errorf("field %s is not settable", name)
	}

	return nil
}

func (f fieldImpl) Name() string {
	return f.field.Name
}

func (f fieldImpl) Tag() reflect.StructTag {
	return f.field.Tag
}
func (f fieldImpl) Kind() reflect.Kind {
	return f.field.Type.Kind()
}

func (f fieldImpl) PointerInt() *int {
	if f.value.IsNil() {
		return nil
	}
	value := f.Int()
	return &value
}

func (f fieldImpl) Int() int {
	return int(reflect.Indirect(f.value).Int())
}

func (f fieldImpl) PointerInt8() *int8 {
	if f.value.IsNil() {
		return nil
	}
	value := f.Int8()
	return &value
}

func (f fieldImpl) Int8() int8 {
	return int8(reflect.Indirect(f.value).Int())
}

func (f fieldImpl) PointerInt16() *int16 {
	if f.value.IsNil() {
		return nil
	}
	value := f.Int16()
	return &value
}

func (f fieldImpl) Int16() int16 {
	return int16(reflect.Indirect(f.value).Int())
}

func (f fieldImpl) PointerInt32() *int32 {
	if f.value.IsNil() {
		return nil
	}
	value := f.Int32()
	return &value
}

func (f fieldImpl) Int32() int32 {
	return int32(reflect.Indirect(f.value).Int())
}

func (f fieldImpl) PointerInt64() *int64 {
	if f.value.IsNil() {
		return nil
	}
	value := f.Int64()
	return &value
}

func (f fieldImpl) Int64() int64 {
	return reflect.Indirect(f.value).Int()
}

func (f fieldImpl) PointerUint() *uint {
	if f.value.IsNil() {
		return nil
	}
	value := f.Uint()
	return &value
}

func (f fieldImpl) Uint() uint {
	return uint(reflect.Indirect(f.value).Uint())
}

func (f fieldImpl) PointerUint8() *uint8 {
	if f.value.IsNil() {
		return nil
	}
	value := f.Uint8()
	return &value
}

func (f fieldImpl) Uint8() uint8 {
	return uint8(reflect.Indirect(f.value).Uint())
}

func (f fieldImpl) PointerUint16() *uint16 {
	if f.value.IsNil() {
		return nil
	}
	value := f.Uint16()
	return &value
}

func (f fieldImpl) Uint16() uint16 {
	return uint16(reflect.Indirect(f.value).Uint())
}

func (f fieldImpl) PointerUint32() *uint32 {
	if f.value.IsNil() {
		return nil
	}
	value := f.Uint32()
	return &value
}

func (f fieldImpl) Uint32() uint32 {
	return uint32(reflect.Indirect(f.value).Uint())
}

func (f fieldImpl) PointerUint64() *uint64 {
	if f.value.IsNil() {
		return nil
	}
	value := f.Uint64()
	return &value
}

func (f fieldImpl) Uint64() uint64 {
	return reflect.Indirect(f.value).Uint()
}

func (f fieldImpl) PointerFloat32() *float32 {
	if f.value.IsNil() {
		return nil
	}
	value := f.Float32()
	return &value
}

func (f fieldImpl) Float32() float32 {
	return float32(reflect.Indirect(f.value).Float())
}

func (f fieldImpl) PointerFloat64() *float64 {
	if f.value.IsNil() {
		return nil
	}
	value := f.Float64()
	return &value
}

func (f fieldImpl) Float64() float64 {
	return reflect.Indirect(f.value).Float()
}

func (f fieldImpl) PointerString() *string {
	if f.value.IsNil() {
		return nil
	}
	value := f.String()
	return &value
}

func (f fieldImpl) String() string {
	return reflect.Indirect(f.value).String()
}

func (f fieldImpl) PointerBool() *bool {
	if f.value.IsNil() {
		return nil
	}
	value := f.Bool()
	return &value
}

func (f fieldImpl) Bool() bool {
	return reflect.Indirect(f.value).Bool()
}

func (f fieldImpl) PointerTime() *time.Time {
	if f.value.IsNil() {
		return nil
	}
	value := f.Time()
	return &value
}

func (f fieldImpl) Time() time.Time {
	value, ok := reflect.Indirect(f.value).Interface().(time.Time)
	if !ok {
		panic(fmt.Sprintf(`field "%s" is not instance of time.Time`, f.field.Name))
	}

	return value
}

func (f fieldImpl) Interface() interface{} {
	return f.value.Interface()
}

package rt

import (
	"time"

	"github.com/QingShan-Xu/web/ds"
)

type (
	BindReader interface {
		GetTag(name string) (map[string]ds.Tag, bool)
		SafePointerInt(name string) (*int, bool)
		SafeInt(name string) (int, bool)
		SafePointerInt8(name string) (*int8, bool)
		SafeInt8(name string) (int8, bool)
		SafePointerInt16(name string) (*int16, bool)
		SafeInt16(name string) (int16, bool)
		SafePointerInt32(name string) (*int32, bool)
		SafeInt32(name string) (int32, bool)
		SafePointerInt64(name string) (*int64, bool)
		SafeInt64(name string) (int64, bool)
		SafePointerUint(name string) (*uint, bool)
		SafeUint(name string) (uint, bool)
		SafePointerUint8(name string) (*uint8, bool)
		SafeUint8(name string) (uint8, bool)
		SafePointerUint16(name string) (*uint16, bool)
		SafeUint16(name string) (uint16, bool)
		SafePointerUint32(name string) (*uint32, bool)
		SafeUint32(name string) (uint32, bool)
		SafePointerUint64(name string) (*uint64, bool)
		SafeUint64(name string) (uint64, bool)
		SafePointerFloat32(name string) (*float32, bool)
		SafeFloat32(name string) (float32, bool)
		SafePointerFloat64(name string) (*float64, bool)
		SafeFloat64(name string) (float64, bool)
		SafePointerString(name string) (*string, bool)
		SafeString(name string) (string, bool)
		SafePointerBool(name string) (*bool, bool)
		SafeBool(name string) (bool, bool)
		SafePointerTime(name string) (*time.Time, bool)
		SafeTime(name string) (time.Time, bool)
		Interface(name string) interface{}
	}

	bindReaderImpl struct {
		dsReader ds.FieldReader
	}
)

func NewBinderReader(dsReader ds.FieldReader) BindReader {
	return bindReaderImpl{
		dsReader: dsReader,
	}
}

// SafePointerInt 尝试从 reflect.Value 中获取 *int 类型的指针。
// 如果成功，返回指针和 true；否则，返回 nil 和 false。
func (r bindReaderImpl) SafePointerInt(name string) (*int, bool) {
	field, _ := r.dsReader.GetField(name)
	if field != nil {
		return field.SafePointerInt()
	}
	return nil, false
}

// SafeInt 尝试从 reflect.Value 中获取 int 类型的值。
// 如果成功，返回整数值和 true；否则，返回 0 和 false。
func (r bindReaderImpl) SafeInt(name string) (int, bool) {
	field, _ := r.dsReader.GetField(name)
	if field != nil {
		return field.SafeInt()
	}
	return 0, false
}

// SafePointerInt8 尝试从 reflect.Value 中获取 *int8 类型的指针。
// 如果成功，返回指针和 true；否则，返回 nil 和 false。
func (r bindReaderImpl) SafePointerInt8(name string) (*int8, bool) {
	field, _ := r.dsReader.GetField(name)
	if field != nil {
		return field.SafePointerInt8()
	}
	return nil, false
}

// SafeInt8 尝试从 reflect.Value 中获取 int8 类型的值。
// 如果成功，返回 int8 值和 true；否则，返回 0 和 false。
func (r bindReaderImpl) SafeInt8(name string) (int8, bool) {
	field, _ := r.dsReader.GetField(name)
	if field != nil {
		return field.SafeInt8()
	}
	return 0, false
}

// SafePointerInt16 尝试从 reflect.Value 中获取 *int16 类型的指针。
// 如果成功，返回指针和 true；否则，返回 nil 和 false。
func (r bindReaderImpl) SafePointerInt16(name string) (*int16, bool) {
	field, _ := r.dsReader.GetField(name)
	if field != nil {
		return field.SafePointerInt16()
	}
	return nil, false
}

// SafeInt16 尝试从 reflect.Value 中获取 int16 类型的值。
// 如果成功，返回 int16 值和 true；否则，返回 0 和 false。
func (r bindReaderImpl) SafeInt16(name string) (int16, bool) {
	field, _ := r.dsReader.GetField(name)
	if field != nil {
		return field.SafeInt16()
	}
	return 0, false
}

// SafePointerInt32 尝试从 reflect.Value 中获取 *int32 类型的指针。
// 如果成功，返回指针和 true；否则，返回 nil 和 false。
func (r bindReaderImpl) SafePointerInt32(name string) (*int32, bool) {
	field, _ := r.dsReader.GetField(name)
	if field != nil {
		return field.SafePointerInt32()
	}
	return nil, false
}

// SafeInt32 尝试从 reflect.Value 中获取 int32 类型的值。
// 如果成功，返回 int32 值和 true；否则，返回 0 和 false。
func (r bindReaderImpl) SafeInt32(name string) (int32, bool) {
	field, _ := r.dsReader.GetField(name)
	if field != nil {
		return field.SafeInt32()
	}
	return 0, false
}

// SafePointerInt64 尝试从 reflect.Value 中获取 *int64 类型的指针。
// 如果成功，返回指针和 true；否则，返回 nil 和 false。
func (r bindReaderImpl) SafePointerInt64(name string) (*int64, bool) {
	field, _ := r.dsReader.GetField(name)
	if field != nil {
		return field.SafePointerInt64()
	}
	return nil, false
}

// SafeInt64 尝试从 reflect.Value 中获取 int64 类型的值。
// 如果成功，返回 int64 值和 true；否则，返回 0 和 false。
func (r bindReaderImpl) SafeInt64(name string) (int64, bool) {
	field, _ := r.dsReader.GetField(name)
	if field != nil {
		return field.SafeInt64()
	}
	return 0, false
}

// SafeUint 尝试从 reflect.Value 中获取 uint 类型的值。
// 如果成功，返回 uint 值和 true；否则，返回 0 和 false。
func (r bindReaderImpl) SafeUint(name string) (uint, bool) {
	field, _ := r.dsReader.GetField(name)
	if field != nil {
		return field.SafeUint()
	}
	return 0, false
}

// SafePointerUint 尝试从 reflect.Value 中获取 *uint 类型的指针。
// 如果成功，返回指针和 true；否则，返回 nil 和 false。
func (r bindReaderImpl) SafePointerUint(name string) (*uint, bool) {
	field, _ := r.dsReader.GetField(name)
	if field != nil {
		return field.SafePointerUint()
	}
	return nil, false
}

// SafeUint8 尝试从 reflect.Value 中获取 uint8 类型的值。
// 如果成功，返回 uint8 值和 true；否则，返回 0 和 false。
func (r bindReaderImpl) SafeUint8(name string) (uint8, bool) {
	field, _ := r.dsReader.GetField(name)
	if field != nil {
		return field.SafeUint8()
	}
	return 0, false
}

// SafePointerUint8 尝试从 reflect.Value 中获取 *uint8 类型的指针。
// 如果成功，返回指针和 true；否则，返回 nil 和 false。
func (r bindReaderImpl) SafePointerUint8(name string) (*uint8, bool) {
	field, _ := r.dsReader.GetField(name)
	if field != nil {
		return field.SafePointerUint8()
	}
	return nil, false
}

// SafeUint16 尝试从 reflect.Value 中获取 uint16 类型的值。
// 如果成功，返回 uint16 值和 true；否则，返回 0 和 false。
func (r bindReaderImpl) SafeUint16(name string) (uint16, bool) {
	field, _ := r.dsReader.GetField(name)
	if field != nil {
		return field.SafeUint16()
	}
	return 0, false
}

// SafePointerUint16 尝试从 reflect.Value 中获取 *uint16 类型的指针。
// 如果成功，返回指针和 true；否则，返回 nil 和 false。
func (r bindReaderImpl) SafePointerUint16(name string) (*uint16, bool) {
	field, _ := r.dsReader.GetField(name)
	if field != nil {
		return field.SafePointerUint16()
	}
	return nil, false
}

// SafeUint32 尝试从 reflect.Value 中获取 uint32 类型的值。
// 如果成功，返回 uint32 值和 true；否则，返回 0 和 false。
func (r bindReaderImpl) SafeUint32(name string) (uint32, bool) {
	field, _ := r.dsReader.GetField(name)
	if field != nil {
		return field.SafeUint32()
	}
	return 0, false
}

// SafePointerUint32 尝试从 reflect.Value 中获取 *uint32 类型的指针。
// 如果成功，返回指针和 true；否则，返回 nil 和 false。
func (r bindReaderImpl) SafePointerUint32(name string) (*uint32, bool) {
	field, _ := r.dsReader.GetField(name)
	if field != nil {
		return field.SafePointerUint32()
	}
	return nil, false
}

// SafeUint64 尝试从 reflect.Value 中获取 uint64 类型的值。
// 如果成功，返回 uint64 值和 true；否则，返回 0 和 false。
func (r bindReaderImpl) SafeUint64(name string) (uint64, bool) {
	field, _ := r.dsReader.GetField(name)
	if field != nil {
		return field.SafeUint64()
	}
	return 0, false
}

// SafePointerUint64 尝试从 reflect.Value 中获取 *uint64 类型的指针。
// 如果成功，返回指针和 true；否则，返回 nil 和 false。
func (r bindReaderImpl) SafePointerUint64(name string) (*uint64, bool) {
	field, _ := r.dsReader.GetField(name)
	if field != nil {
		return field.SafePointerUint64()
	}
	return nil, false
}

// SafeFloat32 尝试从 reflect.Value 中获取 float32 类型的值。
// 如果成功，返回 float32 值和 true；否则，返回 0.0 和 false。
func (r bindReaderImpl) SafeFloat32(name string) (float32, bool) {
	field, _ := r.dsReader.GetField(name)
	if field != nil {
		return field.SafeFloat32()
	}
	return 0, false
}

// SafePointerFloat32 尝试从 reflect.Value 中获取 *float32 类型的指针。
// 如果成功，返回指针和 true；否则，返回 nil 和 false。
func (r bindReaderImpl) SafePointerFloat32(name string) (*float32, bool) {
	field, _ := r.dsReader.GetField(name)
	if field != nil {
		return field.SafePointerFloat32()
	}
	return nil, false
}

// SafeFloat64 尝试从 reflect.Value 中获取 float64 类型的值。
// 如果成功，返回 float64 值和 true；否则，返回 0.0 和 false。
func (r bindReaderImpl) SafeFloat64(name string) (float64, bool) {
	field, _ := r.dsReader.GetField(name)
	if field != nil {
		return field.SafeFloat64()
	}
	return 0, false
}

// SafePointerFloat64 尝试从 reflect.Value 中获取 *float64 类型的指针。
// 如果成功，返回指针和 true；否则，返回 nil 和 false。
func (r bindReaderImpl) SafePointerFloat64(name string) (*float64, bool) {
	field, _ := r.dsReader.GetField(name)
	if field != nil {
		return field.SafePointerFloat64()
	}
	return nil, false
}

// SafeString 尝试从 reflect.Value 中获取 string 类型的值。
// 如果成功，返回字符串和 true；否则，返回空字符串和 false。
func (r bindReaderImpl) SafeString(name string) (string, bool) {
	field, _ := r.dsReader.GetField(name)
	if field != nil {
		return field.SafeString()
	}
	return "", false
}

// SafePointerString 尝试从 reflect.Value 中获取 *string 类型的指针。
// 如果成功，返回指针和 true；否则，返回 nil 和 false。
func (r bindReaderImpl) SafePointerString(name string) (*string, bool) {
	field, _ := r.dsReader.GetField(name)
	if field != nil {
		return field.SafePointerString()
	}
	return nil, false
}

// SafeBool 尝试从 reflect.Value 中获取 bool 类型的值。
// 如果成功，返回 bool 值和 true；否则，返回 false 和 false。
func (r bindReaderImpl) SafeBool(name string) (bool, bool) {
	field, _ := r.dsReader.GetField(name)
	if field != nil {
		return field.SafeBool()
	}
	return false, false
}

// SafePointerBool 尝试从 reflect.Value 中获取 *bool 类型的指针。
// 如果成功，返回指针和 true；否则，返回 nil 和 false。
func (r bindReaderImpl) SafePointerBool(name string) (*bool, bool) {
	field, _ := r.dsReader.GetField(name)
	if field != nil {
		return field.SafePointerBool()
	}
	return nil, false
}

// SafeTime 尝试从 reflect.Value 中获取 time.Time 类型的值。
// 如果成功，返回 time.Time 值和 true；否则，返回 time.Time{} 和 false。
func (r bindReaderImpl) SafeTime(name string) (time.Time, bool) {
	field, _ := r.dsReader.GetField(name)
	if field != nil {
		return field.SafeTime()
	}
	return time.Time{}, false
}

// SafePointerTime 尝试从 reflect.Value 中获取 *time.Time 类型的指针。
// 如果成功，返回指针和 true；否则，返回 nil 和 false。
func (r bindReaderImpl) SafePointerTime(name string) (*time.Time, bool) {
	field, _ := r.dsReader.GetField(name)
	if field != nil {
		return field.SafePointerTime()
	}
	return nil, false
}

func (r bindReaderImpl) GetTag(name string) (map[string]ds.Tag, bool) {
	field, _ := r.dsReader.GetField(name)
	if field != nil {
		return field.GetTag(), true
	}
	return map[string]ds.Tag{}, false
}

func (r bindReaderImpl) Interface(name string) interface{} {
	field, _ := r.dsReader.GetField(name)
	if field != nil {
		return field.Interface()
	}
	return nil
}

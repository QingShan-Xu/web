package bm

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"gorm.io/gorm"
)

// FlexibleInt 是一个可以从字符串或数字JSON值解析的整数类型
type FlexibleInt int

// UnmarshalJSON 实现了json.Unmarshaler接口
func (fi *FlexibleInt) UnmarshalJSON(data []byte) error {
	// 首先尝试解析为字符串
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		// 如果成功解析为字符串，尝试将其转换为整数
		i, err := strconv.Atoi(s)
		if err != nil {
			return fmt.Errorf("FlexibleInt: failed to parse string as int: %v", err)
		}
		*fi = FlexibleInt(i)
		return nil
	}

	// 如果不是字符串，尝试直接解析为整数
	var i int
	if err := json.Unmarshal(data, &i); err != nil {
		return fmt.Errorf("FlexibleInt: failed to parse as int: %v", err)
	}
	*fi = FlexibleInt(i)
	return nil
}

// Int 返回底层的int值
func (fi FlexibleInt) Int() int {
	return int(fi)
}

type Model struct {
	ID        FlexibleInt    `gorm:"primarykey,type:integer(11)" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"delete_at"`
}

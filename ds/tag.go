// 标签解析相关的定义和方法

package ds

import (
	"reflect"
	"strconv"
	"strings"
)

// Tag 表示结构体字段的标签信息。
type Tag struct {
	Key     string   // 标签键
	Value   string   // 标签值
	Options []string // 标签选项
}

// parseFieldTags 解析字段的标签。
func (r *structReader) parseFieldTags() {
	if r.FieldSource == nil {
		return
	}

	r.Tags = make(map[string]Tag)

	// 获取所有标签键
	tagKeys := getTagKeys(r.FieldSource.Tag)

	// 解析每个标签
	for _, key := range tagKeys {
		if value, ok := r.FieldSource.Tag.Lookup(key); ok {
			// 去掉引号
			unquotedValue, err := strconv.Unquote(value)
			if err != nil {
				// 如果去引号失败，使用原始值
				unquotedValue = value
			}

			tagComponents := strings.Split(unquotedValue, ",")
			fieldTag := Tag{
				Key:     key,
				Value:   tagComponents[0],
				Options: tagComponents[1:],
			}
			r.Tags[key] = fieldTag
		}
	}
}

// getTagKeys 获取字段的所有标签键。
// tag: 字段的结构体标签。
// 返回标签键的切片。
func getTagKeys(tag reflect.StructTag) []string {
	var tagKeys []string
	tagStr := string(tag)
	tagPairs := strings.Split(tagStr, " ")
	for _, pair := range tagPairs {
		pair = strings.TrimSpace(pair)
		if pair == "" {
			continue
		}
		keyValue := strings.SplitN(pair, ":", 2)
		if len(keyValue) == 2 {
			tagKeys = append(tagKeys, keyValue[0])
		}
	}
	return tagKeys
}

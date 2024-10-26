// 工具函数

package ds

import "strings"

// parseFieldName 解析字段名称，支持嵌套字段。
// name: 字段名称字符串，例如 "A.B.C"。
// 返回字段名称切片。
func parseFieldName(name string) []string {
	return strings.Split(name, ".")
}

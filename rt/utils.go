// Package rt 提供了通用的工具函数。
package rt

import (
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var (
	commonInitialisms         = []string{"API", "ASCII", "CPU", "CSS", "DNS", "EOF", "GUID", "HTML", "HTTP", "HTTPS", "ID", "IP", "JSON", "LHS", "QPS", "RAM", "RHS", "RPC", "SLA", "SMTP", "SSH", "TLS", "TTL", "UID", "UI", "UUID", "URI", "URL", "UTF8", "VM", "XML", "XSRF", "XSS"}
	commonInitialismsReplacer *strings.Replacer
)

func init() {
	// 初始化首字母缩略词的替换器。
	initialismsForReplacer := make([]string, 0, len(commonInitialisms)*2)
	for _, initialism := range commonInitialisms {
		initialismsForReplacer = append(initialismsForReplacer, initialism, cases.Title(language.Und).String(initialism))
	}
	commonInitialismsReplacer = strings.NewReplacer(initialismsForReplacer...)
}

// ToSnakeCase 将字符串转换为蛇形命名（snake_case）。
// name: 输入字符串。
// 返回转换后的字符串。
func ToSnakeCase(name string) string {
	if name == "" {
		return ""
	}

	value := commonInitialismsReplacer.Replace(name)
	var buf strings.Builder
	buf.Grow(len(value) + 2)

	lastCase := false
	curCase := false

	for i, v := range value {
		if v >= 'A' && v <= 'Z' {
			curCase = true
		} else {
			curCase = false
		}

		if i > 0 && curCase && !lastCase {
			buf.WriteByte('_')
		}

		buf.WriteRune([]rune(strings.ToLower(string(v)))[0])
		lastCase = curCase
	}

	return buf.String()
}

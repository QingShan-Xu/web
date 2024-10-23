package rt

import (
	"fmt"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// 初始化根节点及其子节点的完整路径和名称。
func initCompletePathAndName(root *Router) {
	if root == nil {
		return
	}

	var infoBuilder strings.Builder

	// 初始化根节点的路径和名称。
	root.completePath = root.Path
	root.completeName = root.Name

	// 输出根节点信息，不带前缀符号。
	infoBuilder.WriteString(fmt.Sprintf("%s\n", root.Path))

	// 遍历子节点，设置完整路径和名称。
	for i, child := range root.Children {
		depthFirstProcess(&child, root.Path, root.Name, "", i == len(root.Children)-1, &infoBuilder)
	}

	// 设置根节点的完整信息。
	root.completeInfo = infoBuilder.String()
}

// 执行深度优先遍历，为每个节点设置完整的路径和名称。
func depthFirstProcess(current *Router, currentPath, currentName, prefix string, isLast bool, infoBuilder *strings.Builder) {
	// 计算并设置新的完整路径和名称。
	newPath := removeTrailingSlash(fmt.Sprintf("%s%s", strings.TrimRight(currentPath, "/"), current.Path))
	newName := fmt.Sprintf("%s%s", strings.TrimRight(currentName, "."), strings.TrimLeft(current.Name, "."))

	current.completePath = newPath
	current.completeName = newName

	// 格式化并记录当前节点信息。
	infoLine := fmt.Sprintf("%s%s %s (%s) %s\n", prefix, treeSymbol(isLast), newPath, newName, current.Method)
	infoBuilder.WriteString(infoLine)

	// 获取子节点的前缀。
	childPrefix := generateChildPrefix(prefix, isLast)

	// 递归处理子节点。
	for i := range current.Children {
		child := &current.Children[i]
		depthFirstProcess(child, newPath, newName, childPrefix, i == len(current.Children)-1, infoBuilder)
	}
}

// 移除路径末尾的斜杠。
func removeTrailingSlash(path string) string {
	return strings.TrimRight(path, "/")
}

// treeSymbol 返回合适的树形符号，依据节点位置。
func treeSymbol(isLast bool) string {
	if isLast {
		return "└──"
	}
	return "├──"
}

// 返回子节点的前缀，依据父节点位置。
func generateChildPrefix(prefix string, isLast bool) string {
	if isLast {
		return prefix + "    "
	}
	return prefix + "│   "
}

// 输出给定路由节点的完整树形结构信息。
func displayCompleteInfo(node *Router) {
	fmt.Println(node.completeInfo)
}

// 检查路由节点是否有子节点。
func isGroup(node Router) bool {
	return node.Children != nil
}

var (
	// https://github.com/golang/lint/blob/master/lint.go#L770
	commonInitialisms         = []string{"API", "ASCII", "CPU", "CSS", "DNS", "EOF", "GUID", "HTML", "HTTP", "HTTPS", "ID", "IP", "JSON", "LHS", "QPS", "RAM", "RHS", "RPC", "SLA", "SMTP", "SSH", "TLS", "TTL", "UID", "UI", "UUID", "URI", "URL", "UTF8", "VM", "XML", "XSRF", "XSS"}
	commonInitialismsReplacer *strings.Replacer
)

func init() {
	commonInitialismsForReplacer := make([]string, 0, len(commonInitialisms))
	for _, initialism := range commonInitialisms {
		commonInitialismsForReplacer = append(commonInitialismsForReplacer, initialism, cases.Title(language.Und).String(initialism))
	}
	commonInitialismsReplacer = strings.NewReplacer(commonInitialismsForReplacer...)
}

func ToSnakeCase(name string) string {
	if name == "" {
		return ""
	}

	var (
		value                          = commonInitialismsReplacer.Replace(name)
		buf                            strings.Builder
		lastCase, nextCase, nextNumber bool // upper case == true
		curCase                        = value[0] <= 'Z' && value[0] >= 'A'
	)

	for i, v := range value[:len(value)-1] {
		nextCase = value[i+1] <= 'Z' && value[i+1] >= 'A'
		nextNumber = value[i+1] >= '0' && value[i+1] <= '9'

		if curCase {
			if lastCase && (nextCase || nextNumber) {
				buf.WriteRune(v + 32)
			} else {
				if i > 0 && value[i-1] != '_' && value[i+1] != '_' {
					buf.WriteByte('_')
				}
				buf.WriteRune(v + 32)
			}
		} else {
			buf.WriteRune(v)
		}

		lastCase = curCase
		curCase = nextCase
	}

	if curCase {
		if !lastCase && len(value) > 1 {
			buf.WriteByte('_')
		}
		buf.WriteByte(value[len(value)-1] + 32)
	} else {
		buf.WriteByte(value[len(value)-1])
	}
	ret := buf.String()
	return ret
}

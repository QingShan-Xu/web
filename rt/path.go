// Package rt 包含了处理路由路径的函数。
package rt

import (
	"fmt"
	"strings"
)

// initCompletePathAndName 初始化完整路径和名称。
// root: 根路由器。
func initCompletePathAndName(root *Router) {
	if root == nil {
		return
	}

	var infoBuilder strings.Builder

	root.completePath = root.Path
	root.completeName = root.Name

	// 输出根节点信息。
	infoBuilder.WriteString(fmt.Sprintf("%s\n", root.Path))

	// 遍历子节点，设置完整路径和名称。
	for i, child := range root.Children {
		depthFirstProcess(&child, root.Path, root.Name, "", i == len(root.Children)-1, &infoBuilder)
	}

	// 设置根节点的完整信息。
	root.completeInfo = infoBuilder.String()
}

// depthFirstProcess 深度优先遍历路由树。
// current: 当前路由器。
// currentPath: 当前路径。
// currentName: 当前名称。
// prefix: 前缀字符串，用于显示树形结构。
// isLast: 是否为同级的最后一个节点。
// infoBuilder: 构建路由信息的字符串构建器。
func depthFirstProcess(current *Router, currentPath, currentName, prefix string, isLast bool, infoBuilder *strings.Builder) {
	// 计算并设置新的完整路径和名称。
	newPath := removeTrailingSlash(fmt.Sprintf("%s%s", strings.TrimRight(currentPath, "/"), current.Path))
	newName := strings.TrimRight(strings.TrimLeft(fmt.Sprintf("%s.%s", currentName, current.Name), "."), ".")

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

// removeTrailingSlash 移除路径末尾的斜杠。
// path: 路径字符串。
// 返回处理后的路径。
func removeTrailingSlash(path string) string {
	return strings.TrimRight(path, "/")
}

// treeSymbol 返回合适的树形符号，依据节点位置。
// isLast: 是否为同级的最后一个节点。
// 返回树形符号字符串。
func treeSymbol(isLast bool) string {
	if isLast {
		return "└──"
	}
	return "├──"
}

// generateChildPrefix 返回子节点的前缀，依据父节点位置。
// prefix: 当前前缀。
// isLast: 父节点是否为同级的最后一个节点。
// 返回新的前缀字符串。
func generateChildPrefix(prefix string, isLast bool) string {
	if isLast {
		return prefix + "    "
	}
	return prefix + "│   "
}

// displayCompleteInfo 输出给定路由节点的完整树形结构信息。
// node: 路由器节点。
func displayCompleteInfo(node *Router) {
	fmt.Println(node.completeInfo)
}

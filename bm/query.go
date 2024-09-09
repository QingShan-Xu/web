package bm

type Query struct {
	Type  string // 类型
	Query string // 语句
	Data  interface{}
}

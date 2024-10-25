// Package rt 提供了构建数据库查询的功能。
package rt

import (
	"fmt"
	"reflect"

	"github.com/QingShan-Xu/web/bm"
	"github.com/QingShan-Xu/web/ds"
	"gorm.io/gorm"
)

// Scope 定义了数据库查询范围（Scopes）。
type Scope func(reader *ds.StructReader) func(*gorm.DB) *gorm.DB

// Query 构建数据库查询。
type Query struct {
	bindReader *ds.StructReader
}

// NewQuery 创建新的查询构建器。
// bindReader: 绑定数据的结构体读取器。
// 返回 Query 实例。
func NewQuery(bindReader *ds.StructReader) *Query {
	return &Query{
		bindReader: bindReader,
	}
}

// Model 指定查询的模型。
// model: 数据库模型。
// 返回 Scope 函数或错误信息。
func (q *Query) Model(model interface{}) (Scope, error) {
	modelInstance := reflect.Indirect(reflect.ValueOf(model))
	if modelInstance.Kind() != reflect.Struct {
		return nil, fmt.Errorf("Router.Model must be a struct type")
	}

	modelPtr := reflect.New(modelInstance.Type()).Interface()

	scope := func(reader *ds.StructReader) func(db *gorm.DB) *gorm.DB {
		return func(db *gorm.DB) *gorm.DB {
			return db.Model(modelPtr)
		}
	}

	return scope, nil
}

// Where 添加查询条件。
// whereQuery: 查询条件数组。
// 返回 Scope 函数或错误信息。
func (q *Query) Where(whereQuery []string) (Scope, error) {
	if whereQuery == nil {
		return nil, fmt.Errorf("Where parameter cannot be nil")
	}

	if len(whereQuery) < 2 {
		return nil, fmt.Errorf("Where clause must contain at least two elements")
	}

	for _, fieldName := range whereQuery[1:] {
		if _, err := q.bindReader.GetFieldByName(fieldName); err != nil {
			return nil, err
		}
	}

	return func(reader *ds.StructReader) func(db *gorm.DB) *gorm.DB {
		fieldValues := []interface{}{}
		for _, fieldName := range whereQuery[1:] {
			fieldValue, _ := reader.GetFieldByName(fieldName)
			fieldValues = append(fieldValues, fieldValue.Value.Interface())
		}
		return func(db *gorm.DB) *gorm.DB {
			if containsNil(fieldValues) {
				return db
			}
			return db.Where(whereQuery[0], fieldValues...)
		}
	}, nil
}

// containsNil 检查切片中是否包含 nil 值。
// values: 接口切片。
// 返回是否包含 nil 的布尔值。
func containsNil(values []interface{}) bool {
	for _, v := range values {
		if v == nil {
			return true
		}
	}
	return false
}

// PaginationScope 实现分页。
// pagination: 分页参数。
// 返回 Scope 函数。
func PaginationScope(pagination bm.Pagination) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Limit(pagination.PageSize).Offset((pagination.Current - 1) * pagination.PageSize)
	}
}

// generateQuery 生成查询条件。
// currentRouter: 当前路由器。
// 返回错误信息（如果有）。
func generateQuery(currentRouter *Router) error {
	if isGroup(*currentRouter) {
		for i := range currentRouter.Children {
			child := &currentRouter.Children[i]
			if err := generateQuery(child); err != nil {
				return err
			}
		}
		return nil
	}

	if currentRouter.Bind == nil && (currentRouter.Where != nil || currentRouter.Model != nil) {
		return fmt.Errorf("router '%s' requires Bind when using WHERE or MODEL", currentRouter.completePath)
	}

	var bindReader *ds.StructReader
	var err error
	if currentRouter.Bind != nil {
		bindReader, err = ds.NewStructReader(currentRouter.Bind)
		if err != nil {
			return err
		}
	}

	query := NewQuery(bindReader)

	if currentRouter.Model != nil {
		scope, err := query.Model(currentRouter.Model)
		if err != nil {
			return err
		}
		currentRouter.Scopes = append(currentRouter.Scopes, scope)
	}

	if currentRouter.Where != nil {
		if currentRouter.Model == nil {
			return fmt.Errorf("router '%s' requires Model when using WHERE", currentRouter.completePath)
		}
		for _, where := range currentRouter.Where {
			scope, err := query.Where(where)
			if err != nil {
				return err
			}
			currentRouter.Scopes = append(currentRouter.Scopes, scope)
		}
	}

	return nil
}

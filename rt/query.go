// Package rt 提供了构建数据库查询的功能。
package rt

import (
	"fmt"
	"log"
	"reflect"

	"github.com/QingShan-Xu/web/bm"
	"github.com/QingShan-Xu/web/ds"
	"gorm.io/gorm"
)

// Scope 定义了数据库查询范围（Scopes）。
type Scope func(reader ds.FieldReader) func(*gorm.DB) *gorm.DB

// query 构建数据库查询。
type query struct {
	bindReader ds.FieldReader
}

// newQuery 创建新的查询构建器。
// bindReader: 绑定数据的结构体读取器。
// 返回 Query 实例。
func newQuery(bindReader ds.FieldReader) *query {
	return &query{
		bindReader: bindReader,
	}
}

// model 指定查询的模型。
// model: 数据库模型。
// 返回 Scope 函数或错误信息。
func (q *query) model(model interface{}) (Scope, error) {
	modelInstance := reflect.Indirect(reflect.ValueOf(model))
	if modelInstance.Kind() != reflect.Struct {
		return nil, fmt.Errorf("Router.Model must be a struct type")
	}

	modelPtr := reflect.New(modelInstance.Type()).Interface()

	scope := func(reader ds.FieldReader) func(db *gorm.DB) *gorm.DB {
		return func(db *gorm.DB) *gorm.DB {
			return db.Model(modelPtr)
		}
	}

	return scope, nil
}

// where 添加查询条件。
// whereQuery: 查询条件数组。
// 返回 Scope 函数或错误信息。
func (q *query) where(whereQuery []string) (Scope, error) {
	if whereQuery == nil {
		return nil, fmt.Errorf("where parameter cannot be nil")
	}

	if len(whereQuery) < 2 {
		return nil, fmt.Errorf("where clause must contain at least two elements")
	}

	for _, fieldName := range whereQuery[1:] {
		if _, err := q.bindReader.GetField(fieldName); err != nil {
			return nil, fmt.Errorf("%v in Bind but in Where", err)
		}
	}

	return func(reader ds.FieldReader) func(db *gorm.DB) *gorm.DB {
		fieldValues := []interface{}{}
		for _, fieldName := range whereQuery[1:] {
			fieldValue, _ := reader.GetField(fieldName)
			fieldValues = append(fieldValues, fieldValue.Interface())
		}
		return func(db *gorm.DB) *gorm.DB {
			if containsNil(fieldValues) {
				return db
			}
			return db.Where(whereQuery[0], fieldValues...)
		}
	}, nil
}

// preload 添加查询条件。
// preloadQuery: 查询条件数组。
// 返回 Scope 函数或错误信息。
func (q *query) preload(preloadQuery []string) (Scope, error) {
	if preloadQuery == nil {
		return nil, fmt.Errorf("preload parameter cannot be nil")
	}

	return func(reader ds.FieldReader) func(db *gorm.DB) *gorm.DB {
		return func(db *gorm.DB) *gorm.DB {
			if len(preloadQuery) > 1 {
				preloadArgs := []interface{}{}
				for _, args := range preloadQuery[1:] {
					preloadArgs = append(preloadArgs, args)
				}
				return db.Preload(preloadQuery[0], preloadArgs...)
			}
			return db.Preload(preloadQuery[0])
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
func generateQuery(currentRouter *Router) {
	if isGroup(*currentRouter) {
		for i := range currentRouter.Children {
			child := &currentRouter.Children[i]
			generateQuery(child)
		}
	}

	if currentRouter.Where != nil && currentRouter.Model == nil {
		log.Fatalf("router '%s' requires Model when using WHERE", currentRouter.completePath)
	}

	var bindReader ds.FieldReader
	var err error
	if currentRouter.Bind != nil {
		bindReader, err = ds.NewStructReader(currentRouter.Bind)
		if err != nil {
			log.Fatalf("%v", err)
		}
	}

	query := newQuery(bindReader)

	if currentRouter.Model != nil {
		scope, err := query.model(currentRouter.Model)
		if err != nil {
			log.Fatalf("%v", err)
		}
		currentRouter.Scopes = append(currentRouter.Scopes, scope)
	}

	if currentRouter.Where != nil {
		if currentRouter.Model == nil {
			log.Fatalf("router '%s' requires Model when using WHERE", currentRouter.completePath)
		}
		for _, where := range currentRouter.Where {
			scope, err := query.where(where)
			if err != nil {
				log.Fatalf("%s(%s) %v", currentRouter.completePath, currentRouter.completeName, err)
			}
			currentRouter.Scopes = append(currentRouter.Scopes, scope)
		}
	}
	if currentRouter.Preload != nil {
		if currentRouter.Model == nil {
			log.Fatalf("router '%s' requires Model when using Preload", currentRouter.completePath)
		}
		for _, preload := range currentRouter.Preload {
			scope, err := query.preload(preload)
			if err != nil {
				log.Fatalf("%s(%s) %v", currentRouter.completePath, currentRouter.completeName, err)
			}
			currentRouter.Scopes = append(currentRouter.Scopes, scope)
		}
	}
}

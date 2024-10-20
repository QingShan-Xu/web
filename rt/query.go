package rt

import (
	"fmt"
	"reflect"
	"slices"

	"github.com/QingShan-Xu/web/bm"
	"github.com/QingShan-Xu/web/ds"
	"gorm.io/gorm"
)

type (
	Scope func(reader *ds.StructReader) func(*gorm.DB) *gorm.DB
	Query struct {
		bindReader *ds.StructReader
	}
)

func NewQuery(bindReader *ds.StructReader) *Query {
	return &Query{
		bindReader: bindReader,
	}
}

func (q *Query) MODEL(model interface{}) (Scope, error) {
	modelInstance := reflect.Indirect(reflect.ValueOf(model))
	if modelInstance.Kind() != reflect.Struct {
		return nil, fmt.Errorf("Rt.MODEL 参数 必须为 struct 类型")
	}

	modelPtr := reflect.New(modelInstance.Type()).Interface()

	scope := func(reader *ds.StructReader) func(db *gorm.DB) *gorm.DB {
		return func(db *gorm.DB) *gorm.DB {
			return db.Model(modelPtr)
		}
	}

	return scope, nil
}

func (q *Query) WHERE(whereQuery []string) (Scope, error) {
	if whereQuery == nil {
		return nil, fmt.Errorf("Query.WHERE 参数 不得为 nil")
	}

	if len(whereQuery) < 2 {
		return nil, fmt.Errorf("where 语句 必须包含两个或两个以上元素")
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
			// 如果有空参数, 不创建scope
			if slices.Contains(fieldValues, nil) {
				return db
			}
			return db.Where(whereQuery[0], fieldValues...)
		}
	}, nil
}

var PaginationScope = func(pagination bm.Pagination) func(db *gorm.DB) *gorm.DB {

	return func(db *gorm.DB) *gorm.DB {
		return db.Limit(pagination.PageSize).Offset((pagination.Current - 1) * pagination.PageSize)
	}
}

package rt

import (
	"fmt"
	"slices"

	"github.com/QingShan-Xu/web/ds"
	"gorm.io/gorm"
)

type (
	Scope func(reader ds.Reader) func(*gorm.DB) *gorm.DB
	Query struct {
	}
)

func NewQuery() *Query {
	return &Query{}
}

func (q *Query) WHERE(whereQuery []string) (Scope, error) {
	if whereQuery == nil {
		return nil, fmt.Errorf("Query.WHERE 参数 不得为 nil")
	}

	if len(whereQuery) < 2 {
		return nil, fmt.Errorf("where 语句 必须包含两个或两个以上元素")
	}

	scope := func(reader ds.Reader) func(db *gorm.DB) *gorm.DB {
		return func(db *gorm.DB) *gorm.DB {
			fieldValues := []interface{}{}
			for _, fieldName := range whereQuery[1:] {
				fieldValues = append(fieldValues, reader.GetField(fieldName).Interface())
			}
			if slices.Contains(fieldValues, nil) {
				return db
			}
			return db.Where(whereQuery[0], fieldValues...)
		}
	}

	return scope, nil
}

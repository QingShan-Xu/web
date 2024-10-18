package rt

import (
	"gorm.io/gorm"
)

var METHOD = struct {
	GET, POST, HEAD, PUT, PATCH, DELETE, OPTIONS, TRACE, CONNECT string
}{
	GET:     "GET",
	POST:    "POST",
	HEAD:    "HEAD",
	PUT:     "PUT",
	PATCH:   "PATCH",
	DELETE:  "DELETE",
	OPTIONS: "OPTIONS",
	TRACE:   "TRACE",
	CONNECT: "CONNECT",
}

var PaginationScope = func(pageSize, current int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Limit(pageSize).Offset((current - 1) * pageSize)
	}
}

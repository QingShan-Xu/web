package xjh

import (
	"github.com/QingShan-Xu/xjh/rt"

	"gorm.io/gorm"
)

type ConfigOptions struct {
	DB     *gorm.DB
	Router rt.Route
}

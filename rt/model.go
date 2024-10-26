// Package rt 提供了数据库模型的自动迁移功能。
package rt

import (
	"fmt"

	"github.com/QingShan-Xu/web/db"
)

// generateDBModel 自动迁移数据库模型。
// currentRouter: 当前路由器。
// 返回错误信息（如果有）。
func generateDBModel(currentRouter *Router) error {
	if isGroup(*currentRouter) {
		for i := range currentRouter.Children {
			child := &currentRouter.Children[i]
			if err := generateDBModel(child); err != nil {
				return err
			}
		}
		return nil
	}

	if currentRouter.Model != nil && !currentRouter.NoAutoMigrate {
		if err := db.DB.GORM.AutoMigrate(&currentRouter.Model); err != nil {
			return fmt.Errorf("failed to auto-migrate model for router '%s': %w", currentRouter.completePath, err)
		}
	}

	return nil
}

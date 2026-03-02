package services

import (
	"fmt"

	"gorm.io/gorm"
)

func applyOrgDescendantScope(tx *gorm.DB, column, key string) *gorm.DB {
	condition := fmt.Sprintf("(%s = ? OR %s LIKE ?)", column, column)
	return tx.Where(condition, key, key+"-%")
}

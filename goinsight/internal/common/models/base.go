package models

import "database/sql/driver"

type Model struct {
	ID        uint64    `gorm:"primaryKey" json:"id"`
	CreatedAt LocalTime `gorm:"index:idx_created_at;autoCreateTime;comment:创建时间" json:"created_at"`
	UpdatedAt LocalTime `gorm:"index:idx_updated_at;autoUpdateTime;comment:更新时间" json:"updated_at"`
}

// 枚举类型
type EnumType string

func (rt *EnumType) Scan(value interface{}) error {
	*rt = EnumType(value.([]byte))
	return nil
}

func (rt EnumType) Value() (driver.Value, error) {
	return string(rt), nil
}

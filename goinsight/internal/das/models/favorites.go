package models

import (
	"goInsight/internal/common/models"
)

type InsightDASFavorites struct {
	*models.Model
	Username string `gorm:"type:varchar(128);not null;comment:执行的用户;uniqueIndex:uniq_username_title" json:"username"`
	Title    string `gorm:"type:varchar(256);not null;default:'';comment:标题;uniqueIndex:uniq_username_title" json:"title"`
	Sqltext  string `gorm:"type:text;null;default:null;comment:用户收藏的SQL" json:"sqltext"`
}

func (InsightDASFavorites) TableName() string {
	return "insight_das_favorites"
}

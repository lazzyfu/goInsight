package pagination

import (
	"gorm.io/gorm"
)

type Pagination struct {
	PageSize int  `form:"page_size" json:"page_size"`
	Page     int  `form:"page" json:"page"`
	IsPage   bool `form:"is_page" json:"is_page"`
}

func Pager(p *Pagination, queryTx *gorm.DB, list interface{}) (total int64) {
	if !p.IsPage {
		queryTx.Find(list)
		return
	}
	if p.PageSize < 1 {
		p.PageSize = 10
	}
	if p.Page < 1 {
		p.Page = 1
	}
	offset := p.PageSize * (p.Page - 1)
	// count
	queryTx.Count(&total)
	// 获取分页数据
	queryTx.Limit(p.PageSize).Offset(offset).Find(list)
	return
}

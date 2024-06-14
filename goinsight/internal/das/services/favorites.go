package services

import (
	"errors"
	"goInsight/global"
	"goInsight/internal/das/forms"
	"goInsight/internal/das/models"
	"goInsight/pkg/pagination"

	"github.com/gin-gonic/gin"
)

type GetFavoritesService struct {
	*forms.GetFavoritesForm
	C        *gin.Context
	Username string
}

func (s *GetFavoritesService) Run() (responseData *[]models.InsightDASFavorites, total int64, err error) {
	var list []models.InsightDASFavorites
	tx := global.App.DB.Model(&models.InsightDASFavorites{}).Where("username=?", s.Username).Order("updated_at desc")
	// 搜索schema
	if s.Search != "" {
		tx = tx.Where("`title` like ? or `sqltext` like ?", "%"+s.Search+"%", "%"+s.Search+"%")
	}
	total = pagination.Pager(&s.PaginationQ, tx, &list)
	return &list, total, nil
}

type CreateFavoritesService struct {
	*forms.CreateFavoritesForm
	C        *gin.Context
	Username string
}

func (s *CreateFavoritesService) Run() error {
	favorites := models.InsightDASFavorites{Title: s.Title, Sqltext: s.Sqltext, Username: s.Username}
	tx := global.App.DB.Model(&models.InsightDASFavorites{}).Create(&favorites)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

type UpdateFavoritesService struct {
	*forms.UpdateFavoritesForm
	C        *gin.Context
	ID       uint32
	Username string
}

func (s *UpdateFavoritesService) Run() error {
	tx := global.App.DB.Model(&models.InsightDASFavorites{}).Where("id=? and username=?", s.ID, s.Username)
	data := make(map[string]interface{})
	if s.Title != "" {
		data["Title"] = s.Title
	}
	if s.Sqltext != "" {
		data["Sqltext"] = s.Sqltext
	}
	tx.Updates(data)
	if tx.RowsAffected == 0 {
		return errors.New("更新失败,影响行数为0")
	}
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

type DeleteFavoritesService struct {
	C        *gin.Context
	ID       uint32
	Username string
}

func (s *DeleteFavoritesService) Run() error {
	tx := global.App.DB.Where("id=? and username=?", s.ID, s.Username).Delete(&models.InsightDASFavorites{})
	if tx.RowsAffected == 0 {
		return errors.New("删除失败,影响行数为0")
	}
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

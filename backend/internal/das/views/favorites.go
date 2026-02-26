package views

import (
	"github.com/lazzyfu/goinsight/pkg/response"

	"github.com/lazzyfu/goinsight/internal/das/forms"
	"github.com/lazzyfu/goinsight/internal/das/services"

	"github.com/gin-gonic/gin"
)

// 获取收藏的SQL
func GetFavoritesView(c *gin.Context) {
	username, ok := getUsername(c)
	if !ok {
		return
	}
	var form *forms.GetFavoritesForm = &forms.GetFavoritesForm{}
	if err := c.ShouldBind(&form); err == nil {
		service := services.GetFavoritesService{
			GetFavoritesForm: form,
			C:                c,
			Username:         username,
		}
		returnData, total, err := service.Run()
		if err != nil {
			response.Fail(c, err.Error())
		} else {
			response.PaginationSuccess(c, total, returnData)
		}
	} else {
		response.ValidateFail(c, err.Error())
	}
}

// 新建
func CreateFavoritesView(c *gin.Context) {
	username, ok := getUsername(c)
	if !ok {
		return
	}
	var form *forms.CreateFavoritesForm = &forms.CreateFavoritesForm{}
	if err := c.ShouldBind(&form); err == nil {
		service := services.CreateFavoritesService{
			CreateFavoritesForm: form,
			C:                   c,
			Username:            username,
		}
		err := service.Run()
		if err != nil {
			response.Fail(c, err.Error())
		} else {
			response.Success(c, nil, "success")
		}
	} else {
		response.ValidateFail(c, err.Error())
	}
}

// 更新
func UpdateFavoritesView(c *gin.Context) {
	username, ok := getUsername(c)
	if !ok {
		return
	}
	var form *forms.UpdateFavoritesForm = &forms.UpdateFavoritesForm{}
	id, ok := parseUint32Param(c, "id")
	if !ok {
		return
	}
	if err := c.ShouldBind(&form); err == nil {
		service := services.UpdateFavoritesService{
			UpdateFavoritesForm: form,
			C:                   c,
			ID:                  id,
			Username:            username,
		}
		err := service.Run()
		if err != nil {
			response.Fail(c, err.Error())
		} else {
			response.Success(c, nil, "success")
		}
	} else {
		response.ValidateFail(c, err.Error())
	}
}

// 删除
func DeleteFavoritesView(c *gin.Context) {
	username, ok := getUsername(c)
	if !ok {
		return
	}
	id, ok := parseUint32Param(c, "id")
	if !ok {
		return
	}
	service := services.DeleteFavoritesService{
		C:        c,
		ID:       id,
		Username: username,
	}
	err := service.Run()
	if err != nil {
		response.Fail(c, err.Error())
	} else {
		response.Success(c, nil, "success")
	}
}

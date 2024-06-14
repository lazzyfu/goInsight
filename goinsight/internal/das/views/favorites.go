/*
@Time    :   2023/06/09 14:26:16
@Author  :   xff
*/

package views

import (
	"goInsight/internal/das/forms"
	"goInsight/internal/das/services"
	"goInsight/pkg/response"
	"strconv"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// 获取收藏的SQL
func GetFavoritesView(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	username := claims["id"].(string)
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
	claims := jwt.ExtractClaims(c)
	username := claims["id"].(string)
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
	claims := jwt.ExtractClaims(c)
	username := claims["id"].(string)
	var form *forms.UpdateFavoritesForm = &forms.UpdateFavoritesForm{}
	id, _ := strconv.Atoi(c.Param("id"))
	if err := c.ShouldBind(&form); err == nil {
		service := services.UpdateFavoritesService{
			UpdateFavoritesForm: form,
			C:                   c,
			ID:                  uint32(id),
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
	claims := jwt.ExtractClaims(c)
	username := claims["id"].(string)
	id, _ := strconv.Atoi(c.Param("id"))
	service := services.DeleteFavoritesService{
		C:        c,
		ID:       uint32(id),
		Username: username,
	}
	err := service.Run()
	if err != nil {
		response.Fail(c, err.Error())
	} else {
		response.Success(c, nil, "success")
	}
}

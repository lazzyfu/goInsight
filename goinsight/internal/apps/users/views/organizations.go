/*
@Time    :   2023/09/04 17:09:11
@Author  :   xff
*/

package views

import (
	"goInsight/internal/apps/users/forms"
	"goInsight/internal/apps/users/services"
	"goInsight/internal/pkg/response"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

func GetOrganizationsView(c *gin.Context) {
	service := services.GetOrganizationsServices{C: c}
	returnData := service.Run()
	response.Success(c, returnData, "success")
}

func CreateRootOrganizationsView(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	username := claims["id"].(string)
	var form *forms.CreateRootOrganizationsForm = &forms.CreateRootOrganizationsForm{}
	if err := c.ShouldBind(&form); err == nil {
		service := services.CreateRootOrganizationsService{
			CreateRootOrganizationsForm: form,
			C:                           c,
			Username:                    username,
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

func CreateChildOrganizationsView(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	username := claims["id"].(string)
	var form *forms.CreateChildOrganizationsForm = &forms.CreateChildOrganizationsForm{}
	if err := c.ShouldBind(&form); err == nil {
		service := services.CreateChildOrganizationsService{
			CreateChildOrganizationsForm: form,
			C:                            c,
			Username:                     username,
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

func UpdateOrganizationsView(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	username := claims["id"].(string)
	var form *forms.UpdateOrganizationsForm = &forms.UpdateOrganizationsForm{}
	if err := c.ShouldBind(&form); err == nil {
		service := services.UpdateOrganizationsService{
			UpdateOrganizationsForm: form,
			C:                       c,
			Username:                username,
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

func DeleteOrganizationsView(c *gin.Context) {
	var form *forms.DeleteOrganizationsForm = &forms.DeleteOrganizationsForm{}
	if err := c.ShouldBind(&form); err == nil {
		service := services.DeleteOrganizationsService{
			DeleteOrganizationsForm: form,
			C:                       c,
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

func GetOrganizationsUsersView(c *gin.Context) {
	var form *forms.GetOrganizationsUsersForm = &forms.GetOrganizationsUsersForm{}
	if err := c.ShouldBind(&form); err == nil {
		service := services.GetOrganizationsUsersServices{
			GetOrganizationsUsersForm: form,
			C:                         c,
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

func BindOrganizationsUsersView(c *gin.Context) {
	// claims := jwt.ExtractClaims(c)
	// username := claims["id"].(string)
	var form *forms.BindOrganizationsUsersForm = &forms.BindOrganizationsUsersForm{}
	if err := c.ShouldBind(&form); err == nil {
		service := services.BindOrganizationsUsersService{
			BindOrganizationsUsersForm: form,
			C:                          c,
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

func DeleteOrganizationsUsersView(c *gin.Context) {
	var form *forms.DeleteOrganizationsUsersForm = &forms.DeleteOrganizationsUsersForm{}
	if err := c.ShouldBind(&form); err == nil {
		service := services.DeleteOrganizationsUsersService{
			DeleteOrganizationsUsersForm: form,
			C:                            c,
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

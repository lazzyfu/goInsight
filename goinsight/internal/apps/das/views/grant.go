package views

import (
	"goInsight/internal/apps/das/forms"
	"goInsight/internal/apps/das/services"
	"goInsight/internal/pkg/response"
	"strconv"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// 获取用户权限
func GetUserGrantsView(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	username := claims["id"].(string)
	var form *forms.UserGrantsForm = &forms.UserGrantsForm{}
	if err := c.ShouldBind(&form); err == nil {
		service := services.GetUserGrantsService{UserGrantsForm: form, C: c, Username: username}
		returnData, err := service.Run()
		if err != nil {
			response.Fail(c, err.Error())
		} else {
			response.Success(c, returnData, "success")
		}
	} else {
		response.ValidateFail(c, err.Error())
	}
}

// 获取管理后台权限
func AdminGetSchemasGrantView(c *gin.Context) {
	var form *forms.AdminSchemasGrantForm = &forms.AdminSchemasGrantForm{}
	if err := c.ShouldBind(&form); err == nil {
		service := services.AdminGetSchemasGrantService{AdminSchemasGrantForm: form, C: c}
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

func AdminGetInstancesListView(c *gin.Context) {
	var form *forms.AdminGetInstancesListForm = &forms.AdminGetInstancesListForm{}
	if err := c.ShouldBind(&form); err == nil {
		service := services.AdminGetInstancesListService{
			AdminGetInstancesListForm: form,
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

// 获取指定环境的Schemas
func AdminGetSchemasListView(c *gin.Context) {
	var form *forms.AdminGetSchemasListForm = &forms.AdminGetSchemasListForm{}
	if err := c.ShouldBind(&form); err == nil {
		service := services.AdminGetSchemasListService{
			AdminGetSchemasListForm: form,
			C:                       c,
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

func AdminGetTablesListView(c *gin.Context) {
	var form *forms.AdminGetTablesListForm = &forms.AdminGetTablesListForm{}
	if err := c.ShouldBind(&form); err == nil {
		service := services.AdminGetTablesListService{
			AdminGetTablesListForm: form,
			C:                      c,
		}
		returnData, err := service.Run()
		if err != nil {
			response.Fail(c, err.Error())
		} else {
			response.Success(c, returnData, "success")
		}
	} else {
		response.ValidateFail(c, err.Error())
	}
}

func AdminCreateSchemasGrantView(c *gin.Context) {
	var form *forms.AdminCreateSchemasGrantForm = &forms.AdminCreateSchemasGrantForm{}
	if err := c.ShouldBind(&form); err == nil {
		service := services.AdminCreateSchemasGrantService{
			AdminCreateSchemasGrantForm: form,
			C:                           c,
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

func AdminGetTablesGrantView(c *gin.Context) {
	var form *forms.AdminGetTablesGrantForm = &forms.AdminGetTablesGrantForm{}
	if err := c.ShouldBind(&form); err == nil {
		service := services.AdminGetTablesGrantService{AdminGetTablesGrantForm: form, C: c}
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

func AdminCreateTablesGrantView(c *gin.Context) {
	var form *forms.AdminCreateTablesGrantForm = &forms.AdminCreateTablesGrantForm{}
	if err := c.ShouldBind(&form); err == nil {
		service := services.AdminCreateTablesGrantService{
			AdminCreateTablesGrantForm: form,
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

func AdminDeleteSchemasGrantView(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	service := services.AdminDeleteSchemasGrantService{
		C:  c,
		ID: uint32(id),
	}
	err := service.Run()
	if err != nil {
		response.Fail(c, err.Error())
	} else {
		response.Success(c, nil, "success")
	}
}

func AdminDeleteTablesGrantView(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	service := services.AdminDeleteTablesGrantService{
		C:  c,
		ID: uint32(id),
	}
	err := service.Run()
	if err != nil {
		response.Fail(c, err.Error())
	} else {
		response.Success(c, nil, "success")
	}
}

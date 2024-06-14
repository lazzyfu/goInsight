package views

import (
	"goInsight/internal/inspect/forms"
	"goInsight/internal/inspect/services"
	"goInsight/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

func AdminGetInspectParamsView(c *gin.Context) {
	var form *forms.AdminInspectParamsForm = &forms.AdminInspectParamsForm{}
	if err := c.ShouldBind(&form); err == nil {
		service := services.AdminInspectParamsServices{
			AdminInspectParamsForm: form,
			C:                      c,
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

func AdminUpdateInspectParamsView(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var form *forms.AdminUpdateInspectParamsForm = &forms.AdminUpdateInspectParamsForm{}
	if err := c.ShouldBind(&form); err == nil {
		service := services.AdminUpdateInspectParamsService{
			AdminUpdateInspectParamsForm: form,
			C:                            c,
			ID:                           uint64(id),
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

package views

import (
	"strconv"

	"github.com/lazzyfu/goinsight/pkg/response"

	"github.com/lazzyfu/goinsight/internal/inspect/forms"
	"github.com/lazzyfu/goinsight/internal/inspect/services"

	"github.com/gin-gonic/gin"
)

func AdminGetGlobalInspectParamsView(c *gin.Context) {
	var form forms.AdminGlobalInspectParamsForm
	if err := c.ShouldBind(&form); err != nil {
		response.ValidateFail(c, err.Error())
		return
	}
	service := services.AdminGlobalInspectParamsServices{
		AdminGlobalInspectParamsForm: &form,
		C:                            c,
	}
	returnData, total, err := service.Run()
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.PaginationSuccess(c, total, returnData)
}

func AdminUpdateGlobalInspectParamsView(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var form forms.AdminUpdateGlobalInspectParamsForm
	if err := c.ShouldBind(&form); err != nil {
		response.ValidateFail(c, err.Error())
		return
	}
	service := services.AdminUpdateGlobalInspectParamsService{
		AdminUpdateGlobalInspectParamsForm: &form,
		C:                                  c,
		ID:                                 uint64(id),
	}
	err := service.Run()
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Success(c, nil, "success")
}

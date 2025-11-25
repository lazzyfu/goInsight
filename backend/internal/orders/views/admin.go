package views

import (
	"strconv"

	"github.com/lazzyfu/goinsight/pkg/response"

	"github.com/lazzyfu/goinsight/internal/orders/forms"
	"github.com/lazzyfu/goinsight/internal/orders/services"

	"github.com/gin-gonic/gin"
)

func AdminGetApprovalFlowsView(c *gin.Context) {
	var form forms.AdminGetApprovalFlowsForm
	if err := c.ShouldBind(&form); err != nil {
		response.ValidateFail(c, err.Error())
		return
	}
	service := services.AdminGetApprovalFlowsService{
		AdminGetApprovalFlowsForm: &form,
		C:                         c,
	}
	returnData, total, err := service.Run()
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.PaginationSuccess(c, total, returnData)
}

func AdminUpdateApprovalFlowsView(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var form forms.AdminUpdateApprovalFlowsForm

	if err := c.ShouldBind(&form); err != nil {
		response.ValidateFail(c, err.Error())
		return
	}
	service := services.AdminUpdateApprovalFlowsService{
		AdminUpdateApprovalFlowsForm: &form,
		C:                            c,
		ID:                           uint64(id),
	}
	returnData, total, err := service.Run()
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.PaginationSuccess(c, total, returnData)
}

func AdminCreateApprovalFlowsView(c *gin.Context) {
	var form forms.AdminCreateApprovalFlowsForm

	if err := c.ShouldBind(&form); err != nil {
		response.ValidateFail(c, err.Error())
		return
	}
	service := services.AdminCreateApprovalFlowsService{
		AdminCreateApprovalFlowsForm: &form,
		C:                            c,
	}
	err := service.Run()
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Success(c, nil, "success")
}

func AdminDeleteApprovalFlowsView(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	service := services.AdminDeleteApprovalFlowsService{
		C:  c,
		ID: uint64(id),
	}
	err := service.Run()
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.Success(c, nil, "success")
}

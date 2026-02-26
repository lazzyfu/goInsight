package views

import (
	"strconv"

	"github.com/lazzyfu/goinsight/middleware"
	"github.com/lazzyfu/goinsight/pkg/response"

	"github.com/gin-gonic/gin"
)

func getUsername(c *gin.Context) (string, bool) {
	username, ok := middleware.GetUserNameFromJWT(c)
	if !ok {
		response.Fail(c, "认证信息无效")
		return "", false
	}
	return username, true
}

func parseUint32Param(c *gin.Context, name string) (uint32, bool) {
	raw := c.Param(name)
	id, err := strconv.ParseUint(raw, 10, 32)
	if err != nil {
		response.ValidateFail(c, "非法参数: "+name)
		return 0, false
	}
	return uint32(id), true
}

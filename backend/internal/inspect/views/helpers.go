package views

import (
	"strconv"

	"github.com/lazzyfu/goinsight/pkg/response"

	"github.com/gin-gonic/gin"
)

func parseUint64Param(c *gin.Context, name string) (uint64, bool) {
	raw := c.Param(name)
	id, err := strconv.ParseUint(raw, 10, 64)
	if err != nil {
		response.ValidateFail(c, "非法参数: "+name)
		return 0, false
	}
	return id, true
}

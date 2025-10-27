package views

import (
	"net/http"

	"github.com/lazzyfu/goinsight/internal/global"

	"github.com/gin-gonic/gin"
)

// 返回自定义网站title
func GetAppTitleView(c *gin.Context) {
	title := "GoInsight"
	if len(global.App.Config.App.Title) > 0 {
		title = global.App.Config.App.Title
	}
	c.JSON(http.StatusOK, gin.H{"data": title})
}

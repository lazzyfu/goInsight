package bootstrap

import (
	"github.com/lazzyfu/goinsight/middleware"

	"github.com/lazzyfu/goinsight/internal/global"
)

func InitializeLog() {
	global.App.Log = middleware.InitLogger("app.log")
}

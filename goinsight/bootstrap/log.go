/*
@Time    :   2023/08/14 15:49:31
@Author  :   xff
*/

package bootstrap

import (
	"goInsight/global"
	"goInsight/middleware"
)

func InitializeLog() {
	global.App.Log = middleware.InitLogger("app.log")
}

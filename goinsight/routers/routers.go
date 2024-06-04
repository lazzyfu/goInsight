/*
@Time    :   2023/08/14 18:16:00
@Author  :   xff
*/

package routers

import (
	"goInsight/middleware"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
)

type Option func(*gin.Engine)

var options = []Option{}

// 注册app的路由配置
func Include(opts ...Option) {
	options = append(options, opts...)
}

// 初始化
func Init() *gin.Engine {
	r := gin.New()

	// 使用CORS中间件
	r.Use(cors.Default())

	// 使用requestid
	r.Use(requestid.New())

	// 初始化请求日志记录器
	requestLogger := middleware.InitLogger(time.Now().Format("2006-01-02") + "-request.log")
	r.Use(middleware.LoggerRequestToFile(requestLogger))

	for _, opt := range options {
		opt(r)
	}
	return r
}

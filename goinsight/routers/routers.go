/*
@Time    :   2023/08/14 18:16:00
@Author  :   xff
*/

package routers

import (
	"goInsight/middleware"

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
	// 将http请求记录到文件
	r.Use(middleware.LoggerRequestToFile())
	for _, opt := range options {
		opt(r)
	}
	return r
}

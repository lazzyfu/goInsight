/*
@Time    :   2023/08/14 15:51:08
@Author  :   lazzyfu
*/

package middleware

import (
	"fmt"
	"goInsight/global"
	"time"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func LoggerRequestToFile() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 开始时间
		startTime := time.Now()
		c.Next()
		endTime := time.Now()
		latencyTime := fmt.Sprintf("%dms", endTime.Sub(startTime).Milliseconds())

		//日志格式
		global.App.Log.WithFields(logrus.Fields{
			"status_code":       c.Writer.Status(),
			"latency_time":      latencyTime,
			"request_client_ip": c.ClientIP(),
			"request_method":    c.Request.Method,
			"request_uri":       c.Request.RequestURI,
			"request_ua":        c.Request.UserAgent(),
			"request_referer":   c.Request.Referer(),
			"request_id":        requestid.Get(c),
			"request_header":    c.Request.Header,
		}).Info()
	}
}

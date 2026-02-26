package middleware

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"time"

	"github.com/lazzyfu/goinsight/internal/global"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

// InitLogger 初始化日志记录器，支持日志文件轮转
func InitLogger(logFileName string) *logrus.Logger {
	// 创建日志文件夹
	logFilePath := global.App.Config.Log.RootDir
	if err := os.MkdirAll(logFilePath, 0777); err != nil {
		fmt.Println(err.Error())
	}

	// 实例化
	logger := logrus.New()

	// 使用lumberjack进行日志轮转
	logger.SetOutput(&lumberjack.Logger{
		Filename:   path.Join(logFilePath, logFileName),
		MaxSize:    100,  // 每个日志文件的最大尺寸（单位：MB）
		MaxBackups: 30,   // 保留的旧日志文件的最大数量
		MaxAge:     7,    // 保留旧日志文件的最大天数
		Compress:   true, // 是否压缩旧的日志文件
	})

	// 设置日志级别
	switch strings.ToLower(global.App.Config.Log.Level) {
	case "debug":
		logger.SetLevel(logrus.DebugLevel)
	case "info":
		logger.SetLevel(logrus.InfoLevel)
	case "warn":
		logger.SetLevel(logrus.WarnLevel)
	case "error":
		logger.SetLevel(logrus.ErrorLevel)
	default:
		logger.SetLevel(logrus.InfoLevel)
	}

	// 设置日志格式
	logger.Formatter = &logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	}
	return logger
}

func LoggerRequestToFile(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 开始时间
		startTime := time.Now()
		// 继续处理请求
		c.Next()
		// 结束时间
		endTime := time.Now()
		latencyTime := fmt.Sprintf("%dms", endTime.Sub(startTime).Milliseconds())

		safeURI := sanitizeRequestURI(c.Request.URL)
		safeHeaders := sanitizeRequestHeaders(c.Request.Header)

		//日志格式
		logger.WithFields(logrus.Fields{
			"status_code":       c.Writer.Status(),
			"latency_time":      latencyTime,
			"request_client_ip": c.ClientIP(),
			"request_method":    c.Request.Method,
			"request_uri":       safeURI,
			"request_ua":        c.Request.UserAgent(),
			"request_referer":   c.Request.Referer(),
			"request_id":        requestid.Get(c),
			"request_header":    safeHeaders,
		}).Info()
	}
}

var sensitiveHeaderKeys = map[string]struct{}{
	"authorization": {},
	"cookie":        {},
	"set-cookie":    {},
	"x-api-key":     {},
}

var sensitiveQueryKeys = map[string]struct{}{
	"token":        {},
	"password":     {},
	"old_password": {},
	"new_password": {},
	"otp_code":     {},
}

func sanitizeRequestHeaders(headers http.Header) map[string][]string {
	if headers == nil {
		return map[string][]string{}
	}
	masked := make(map[string][]string, len(headers))
	for k, values := range headers {
		if _, ok := sensitiveHeaderKeys[strings.ToLower(k)]; ok {
			masked[k] = []string{"***"}
			continue
		}
		copied := make([]string, len(values))
		copy(copied, values)
		masked[k] = copied
	}
	return masked
}

func sanitizeRequestURI(u *url.URL) string {
	if u == nil {
		return ""
	}
	query := u.Query()
	for key := range query {
		if _, ok := sensitiveQueryKeys[strings.ToLower(key)]; ok {
			query.Set(key, "***")
		}
	}
	safeURL := *u
	safeURL.RawQuery = query.Encode()
	return safeURL.RequestURI()
}

/*
@Time    :   2023/08/14 15:49:31
@Author  :   xff
*/

package bootstrap

import (
	"fmt"
	"goInsight/global"
	"os"
	"path"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

func InitializeLog() {
	now := time.Now()
	global.App.Log = initLog(now.Format("2006-01-02") + ".log")
}

func initLog(logFileName string) *logrus.Logger {
	logFilePath := global.App.Config.Log.RootDir
	if err := os.MkdirAll(logFilePath, 0777); err != nil {
		fmt.Println(err.Error())
	}
	// 日志文件
	fileName := path.Join(logFilePath, logFileName)
	if _, err := os.Stat(fileName); err != nil {
		if _, err := os.Create(fileName); err != nil {
			fmt.Println(err.Error())
		}
	}
	// 写入文件
	src, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		fmt.Println("err", err)
	}

	// 实例化
	logger := logrus.New()

	// 设置输出
	logger.Out = src

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
		TimestampFormat:   "2006-01-02 15:04:05",
		DisableHTMLEscape: true,
	}
	return logger
}

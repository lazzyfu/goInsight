package tidb

import (
	"fmt"
	"goInsight/pkg/utils"
	"strings"
	"time"

	"goInsight/internal/orders/api/base"
)

// TiDB DML
type ExecuteTiDBDML struct {
	*base.DBConfig
}

// 执行TiDB DML
func (e *ExecuteTiDBDML) Run() (data base.ReturnData, err error) {
	var executeLog []string

	// Function to log messages and publish
	logAndPublish := func(msg string) {
		timestamp := time.Now().Format("2006-01-02 15:04:05")
		formattedMsg := fmt.Sprintf("[%s] %s", timestamp, msg)
		executeLog = append(executeLog, formattedMsg)
		base.PublishMessageToChannel(e.OrderID, formattedMsg, "")
	}

	// Logging function for errors
	logErrorAndReturn := func(err error, errMsg string) (base.ReturnData, error) {
		logAndPublish(errMsg + err.Error())
		data.ExecuteLog = strings.Join(executeLog, "\n")
		return data, err
	}

	// CREATE A NEW DATABASE CONNECTION
	db, err := NewTiDBCnx(e.DBConfig)
	if err != nil {
		return logErrorAndReturn(base.SQLExecuteError{Err: err},
			fmt.Sprintf("访问数据库(%s:%d)失败，错误：", e.DBConfig.Hostname, e.DBConfig.Port))
	}
	defer db.Close()
	logAndPublish(fmt.Sprintf("访问数据库(%s:%d)成功", e.DBConfig.Hostname, e.DBConfig.Port))

	// GET CONNECTION ID
	connectionID, err := DaoGetTiDBConnectionID(db)
	if err != nil {
		return logErrorAndReturn(base.SQLExecuteError{Err: err}, "获取数据库Connection ID失败，错误：")
	}
	logAndPublish(fmt.Sprintf("数据库Connection ID：%d", connectionID))

	// SHOW PROCESS
	ch1 := make(chan int64)
	go DaoTiDBGetProcesslist(e.DBConfig, e.OrderID, connectionID, ch1)

	// 执行SQL
	startTime := time.Now()
	affectedRows, err := DaoTiDBExecute(db, e.SQL, ch1)
	if err != nil {
		return logErrorAndReturn(base.SQLExecuteError{Err: err}, "SQL执行失败，错误：")
	}
	endTime := time.Now()
	executeCostTime := utils.HumanfriendlyTimeUnit(endTime.Sub(startTime))
	logAndPublish(fmt.Sprintf("SQL执行成功，影响行数%d，执行耗时：%s", affectedRows, executeCostTime))

	// TiDB不支持生成回滚SQL
	logAndPublish("TiDB不支持生成回滚SQL")

	// 返回数据
	data.ExecuteLog = strings.Join(executeLog, "\n")
	data.AffectedRows = affectedRows
	data.ExecuteCostTime = executeCostTime
	return
}

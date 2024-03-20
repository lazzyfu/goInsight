/*
@Author  :   lazzyfu
@Desc    :   执行DML
*/

package api

import (
	"fmt"
	"goInsight/internal/pkg/utils"
	"strings"
	"time"
)

// MySQL DML
type ExecuteMySQLDML struct {
	*DBConfig
}

// 执行MySQL DML语句
func (e *ExecuteMySQLDML) Run() (data ReturnData, err error) {
	// Slice for executeLog
	var executeLog []string

	// Function to log messages and publish
	logAndPublish := func(msg string) {
		timestamp := time.Now().Format("2006-01-02 15:04:05")
		formattedMsg := fmt.Sprintf("[%s] %s", timestamp, msg)
		executeLog = append(executeLog, formattedMsg)
		PublishMsg(e.OrderID, formattedMsg, "")
	}

	// Logging function for errors
	logErrorAndReturn := func(err error, errMsg string) (ReturnData, error) {
		logAndPublish(errMsg + err.Error())
		data.ExecuteLog = strings.Join(executeLog, "\n")
		return data, err
	}

	// CREATE A NEW DATABASE CONNECTION
	db, err := NewMySQLCnx(e.DBConfig)
	if err != nil {
		return logErrorAndReturn(SQLExecuteError{Err: err}, fmt.Sprintf("访问数据库(%s:%d)失败，错误：%s", e.DBConfig.Hostname, e.DBConfig.Port, err.Error()))
	}
	defer db.Close()
	logAndPublish(fmt.Sprintf("访问数据库(%s:%d)成功", e.DBConfig.Hostname, e.DBConfig.Port))

	// GET CONNECTION ID
	connectionID, err := DaoGetConnectionID(db)
	if err != nil {
		return logErrorAndReturn(SQLExecuteError{Err: err}, "获取数据库Connection ID失败，错误：")
	}
	logAndPublish(fmt.Sprintf("数据库Connection ID：%d", connectionID))

	// SHOW PROCESS
	ch1 := make(chan int64)
	go DaoGetProcesslist(e.DBConfig, e.OrderID, connectionID, ch1)

	// 获取执行开始前的binlog position
	startFile, startPosition, err := DaoGetMySQLPos(db)
	if err != nil {
		return logErrorAndReturn(SQLExecuteError{Err: err}, "获取Start Binlog File和Position失败，错误：")
	}
	logAndPublish(fmt.Sprintf("Start Binlog File：%s，Position：%d", startFile, startPosition))

	// 执行SQL
	startTime := time.Now()
	affectedRows, err := DaoMySQLExecute(db, e.SQL, ch1)
	if err != nil {
		return logErrorAndReturn(SQLExecuteError{Err: err}, "SQL执行失败，错误：")
	}
	endTime := time.Now()
	executeCostTime := utils.HumanfriendlyTimeUnit(endTime.Sub(startTime))
	logAndPublish(fmt.Sprintf("SQL执行成功，影响行数%d，执行耗时：%s", affectedRows, executeCostTime))

	data.AffectedRows = affectedRows
	data.ExecuteCostTime = executeCostTime

	// 获取执行后的binlog position
	endFile, endPosition, err := DaoGetMySQLPos(db)
	if err != nil {
		return logErrorAndReturn(RollbackSQLError{Err: err}, "获取End Binlog File和Position失败，错误：")
	}
	logAndPublish(fmt.Sprintf("End Binlog File：%s，Position：%d", endFile, endPosition))

	var rollbackSQL, backupCostTime string
	// 影响行数大于0，才执行生成回滚SQL操作
	if affectedRows > 0 {
		// 生成回滚SQL
		logAndPublish("开始解析Binlog生成回滚SQL")
		startTime = time.Now()
		binlog := Binlog{DBConfig: e.DBConfig,
			ConnectionID:  connectionID,
			StartFile:     startFile,
			StartPosition: startPosition,
			EndFile:       endFile,
			EndPosition:   endPosition}
		rollbackSQL, err = binlog.Run()
		if err != nil {
			return logErrorAndReturn(RollbackSQLError{Err: err}, "生成回滚SQL失败，错误：")
		}

		endTime = time.Now()
		backupCostTime = utils.HumanfriendlyTimeUnit(endTime.Sub(startTime))
		logAndPublish(fmt.Sprintf("生成回滚SQL成功，耗时：%s", backupCostTime))
	}
	// 返回数据
	data.ExecuteLog = strings.Join(executeLog, "\n")
	data.RollbackSQL = rollbackSQL
	data.BackupCostTime = backupCostTime
	return
}

// TiDB DML
type ExecuteTiDBDML struct {
	*DBConfig
}

// 执行TiDB DML
func (e *ExecuteTiDBDML) Run() (data ReturnData, err error) {
	var executeLog []string

	// Function to log messages and publish
	logAndPublish := func(msg string) {
		timestamp := time.Now().Format("2006-01-02 15:04:05")
		formattedMsg := fmt.Sprintf("[%s] %s", timestamp, msg)
		executeLog = append(executeLog, formattedMsg)
		PublishMsg(e.OrderID, formattedMsg, "")
	}

	// Logging function for errors
	logErrorAndReturn := func(err error, errMsg string) (ReturnData, error) {
		logAndPublish(errMsg + err.Error())
		data.ExecuteLog = strings.Join(executeLog, "\n")
		return data, err
	}

	// CREATE A NEW DATABASE CONNECTION
	db, err := NewMySQLCnx(e.DBConfig)
	if err != nil {
		return logErrorAndReturn(SQLExecuteError{Err: err},
			fmt.Sprintf("访问数据库(%s:%d)失败，错误：", e.DBConfig.Hostname, e.DBConfig.Port))
	}
	defer db.Close()
	logAndPublish(fmt.Sprintf("访问数据库(%s:%d)成功", e.DBConfig.Hostname, e.DBConfig.Port))

	// GET CONNECTION ID
	connectionID, err := DaoGetConnectionID(db)
	if err != nil {
		return logErrorAndReturn(SQLExecuteError{Err: err}, "获取数据库Connection ID失败，错误：")
	}
	logAndPublish(fmt.Sprintf("数据库Connection ID：%d", connectionID))

	// SHOW PROCESS
	ch1 := make(chan int64)
	go DaoGetProcesslist(e.DBConfig, e.OrderID, connectionID, ch1)

	// 执行SQL
	startTime := time.Now()
	affectedRows, err := DaoMySQLExecute(db, e.SQL, ch1)
	if err != nil {
		return logErrorAndReturn(SQLExecuteError{Err: err}, "SQL执行失败，错误：")
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

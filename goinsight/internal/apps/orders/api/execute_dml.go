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

// MySQL
type ExecuteMySQLDML struct {
	*DBConfig
}

func (e *ExecuteMySQLDML) Run() (data ReturnData, err error) {
	var executeLog []string
	var msg string

	// Create a new database connection
	db, err := NewMySQLCnx(e.DBConfig)
	if err != nil {
		data.ExecuteLog = fmt.Sprintf("访问数据库(%s:%d)失败，错误：%s", e.DBConfig.Hostname, e.DBConfig.Port, err.Error())
		return data, err
	}
	defer db.Close()
	msg = fmt.Sprintf("[%s] 访问数据库(%s:%d)成功", time.Now().Format("2006-01-02 15:04:05"), e.DBConfig.Hostname, e.DBConfig.Port)
	executeLog = append(executeLog, msg)
	PublishMsg(e.OrderID, msg, "")

	// get connection id
	ConnectionID, err := DaoGetConnectionID(db)
	if err != nil {
		data.ExecuteLog = fmt.Sprintf("获取数据库Connection ID失败，错误：%s", err.Error())
		return data, err
	}
	msg = fmt.Sprintf("[%s] 数据库Connection ID：%d", time.Now().Format("2006-01-02 15:04:05"), ConnectionID)
	executeLog = append(executeLog, msg)
	PublishMsg(e.OrderID, msg, "")

	// show process
	ch1 := make(chan int64)
	go DaoGetProcesslist(e.DBConfig, e.OrderID, ConnectionID, ch1)

	// 获取position
	startFile, startPosition, err := DaoGetMySQLPos(db)
	if err != nil {
		data.ExecuteLog = fmt.Sprintf("获取Start Binlog File和Position失败，错误：%s", err.Error())
		return data, err
	}
	msg = fmt.Sprintf("[%s] Start Binlog File：%s，Position：%d", time.Now().Format("2006-01-02 15:04:05"), startFile, startPosition)
	executeLog = append(executeLog, msg)
	PublishMsg(e.OrderID, msg, "")

	// 开始执行SQL
	startTime := time.Now()
	affectedRows, err := DaoMySQLExecute(db, e.SQL, ch1)
	if err != nil {
		data.ExecuteLog = fmt.Sprintf("SQL执行失败，错误：%s", err.Error())
		return data, err
	}
	endTime := time.Now()
	executeCostTime := utils.HumanfriendlyTimeUnit(endTime.Sub(startTime))
	msg = fmt.Sprintf("[%s] SQL执行成功，影响行数%d，执行耗时：%s", time.Now().Format("2006-01-02 15:04:05"), affectedRows, executeCostTime)
	executeLog = append(executeLog, msg)
	PublishMsg(e.OrderID, msg, "")

	data.AffectedRows = affectedRows
	data.ExecuteCostTime = executeCostTime

	// 获取position
	endFile, endPosition, err := DaoGetMySQLPos(db)
	if err != nil {
		data.ExecuteLog = fmt.Sprintf("获取End Binlog File和Position失败，错误：%s", err.Error())
		return data, err
	}
	msg = fmt.Sprintf("[%s] End Binlog File：%s，Position：%d", time.Now().Format("2006-01-02 15:04:05"), endFile, endPosition)
	executeLog = append(executeLog, msg)
	PublishMsg(e.OrderID, msg, "")

	// 回滚SQL
	msg = fmt.Sprintf("[%s] 开始生成回滚SQL，解析Binlog", time.Now().Format("2006-01-02 15:04:05"))
	executeLog = append(executeLog, msg)
	PublishMsg(e.OrderID, msg, "")

	var rollbackSQL, backupCostTime string
	if affectedRows > 0 {
		// 获取回滚SQL
		startTime = time.Now()
		binlog := Binlog{DBConfig: e.DBConfig,
			ConnectionID:  ConnectionID,
			StartFile:     startFile,
			StartPosition: startPosition,
			EndFile:       endFile,
			EndPosition:   endPosition}
		rollbackSQL, err = binlog.Run()
		if err != nil {
			data.ExecuteLog = fmt.Sprintf("回滚SQL生成失败，错误：%s", err.Error())
			return data, err
		}

		endTime = time.Now()
		backupCostTime = utils.HumanfriendlyTimeUnit(endTime.Sub(startTime))

		msg = fmt.Sprintf("[%s] 回滚SQL生成成功，耗时：%s", time.Now().Format("2006-01-02 15:04:05"), backupCostTime)
		executeLog = append(executeLog, msg)
		PublishMsg(e.OrderID, msg, "")
	}
	// 返回数据
	data.ExecuteLog = strings.Join(executeLog, "\n")
	data.RollbackSQL = rollbackSQL
	data.BackupCostTime = backupCostTime
	return
}

// TiDB
type ExecuteTiDBDML struct {
	*DBConfig
}

func (e *ExecuteTiDBDML) Run() (data ReturnData, err error) {
	var executeLog []string
	var msg string

	// Create a new database connection
	db, err := NewMySQLCnx(e.DBConfig)
	if err != nil {
		data.ExecuteLog = fmt.Sprintf("访问数据库(%s:%d)失败，错误：%s", e.DBConfig.Hostname, e.DBConfig.Port, err.Error())
		return data, err
	}
	defer db.Close()
	msg = fmt.Sprintf("[%s] 访问数据库(%s:%d)成功", time.Now().Format("2006-01-02 15:04:05"), e.DBConfig.Hostname, e.DBConfig.Port)
	executeLog = append(executeLog, msg)
	PublishMsg(e.OrderID, msg, "")

	// get connection id
	ConnectionID, err := DaoGetConnectionID(db)
	if err != nil {
		data.ExecuteLog = fmt.Sprintf("获取数据库Connection ID失败，错误：%s", err.Error())
		return data, err
	}
	msg = fmt.Sprintf("[%s] 数据库Connection ID：%d", time.Now().Format("2006-01-02 15:04:05"), ConnectionID)
	executeLog = append(executeLog, msg)
	PublishMsg(e.OrderID, msg, "")

	// show process
	ch1 := make(chan int64)
	go DaoGetProcesslist(e.DBConfig, e.OrderID, ConnectionID, ch1)

	// 执行SQL
	startTime := time.Now()
	affectedRows, err := DaoMySQLExecute(db, e.SQL, ch1)
	if err != nil {
		data.ExecuteLog = fmt.Sprintf("SQL执行失败，错误：%s", err.Error())
		return data, err
	}
	endTime := time.Now()
	executeCostTime := utils.HumanfriendlyTimeUnit(endTime.Sub(startTime))

	msg = fmt.Sprintf("[%s] SQL执行成功，影响行数%d，执行耗时：%s", time.Now().Format("2006-01-02 15:04:05"), affectedRows, executeCostTime)
	executeLog = append(executeLog, msg)
	PublishMsg(e.OrderID, msg, "")

	// TiDB不支持生成回滚SQL
	msg = fmt.Sprintf("[%s] TiDB不支持生成回滚SQL", time.Now().Format("2006-01-02 15:04:05"))
	executeLog = append(executeLog, msg)
	PublishMsg(e.OrderID, msg, "")

	// 返回数据
	data.ExecuteLog = strings.Join(executeLog, "\n")
	data.AffectedRows = affectedRows
	data.ExecuteCostTime = executeCostTime
	return
}

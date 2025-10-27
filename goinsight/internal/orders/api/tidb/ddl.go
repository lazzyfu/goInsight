package tidb

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/lazzyfu/goinsight/pkg/parser"
	"github.com/lazzyfu/goinsight/pkg/utils"

	"github.com/lazzyfu/goinsight/internal/orders/api/base"
)

// 执行Online DDL语句
func ExecuteOnlineDDL(dc *base.DBConfig) (data base.ReturnData, err error) {
	var executeLog []string

	// Function to log messages and publish
	logAndPublish := func(msg string) {
		timestamp := time.Now().Format("2006-01-02 15:04:05")
		formattedMsg := fmt.Sprintf("[%s] %s", timestamp, msg)
		executeLog = append(executeLog, formattedMsg)
		base.PublishMessageToChannel(dc.OrderID, formattedMsg, "")
	}

	// Logging function for errors
	logErrorAndReturn := func(err error, errMsg string) (base.ReturnData, error) {
		logAndPublish(errMsg + err.Error())
		data.ExecuteLog = strings.Join(executeLog, "\n")
		return data, err
	}

	// CREATE A NEW DATABASE CONNECTION
	db, err := NewTiDBCnx(dc)
	if err != nil {
		return logErrorAndReturn(base.SQLExecuteError{Err: err}, fmt.Sprintf("访问数据库(%s:%d)失败，错误：%s", dc.Hostname, dc.Port, err.Error()))
	}
	defer db.Close()
	logAndPublish(fmt.Sprintf("访问数据库(%s:%d)成功", dc.Hostname, dc.Port))

	// GET CONNECTION ID
	connectionID, err := DaoGetTiDBConnectionID(db)
	if err != nil {
		return logErrorAndReturn(base.SQLExecuteError{Err: err}, "获取数据库Connection ID失败，错误：")
	}
	logAndPublish(fmt.Sprintf("数据库Connection ID：%d", connectionID))

	// SHOW PROCESS
	ch1 := make(chan int64)
	go DaoTiDBGetProcesslist(dc, dc.OrderID, connectionID, ch1)

	// 执行SQL
	startTime := time.Now()
	affectedRows, err := DaoTiDBExecute(db, dc.SQL, ch1)
	if err != nil {
		return logErrorAndReturn(base.SQLExecuteError{Err: err}, "SQL执行失败，错误：")
	}
	endTime := time.Now()
	executeCostTime := utils.HumanfriendlyTimeUnit(endTime.Sub(startTime))
	logAndPublish(fmt.Sprintf("SQL执行成功，影响行数%d，执行耗时：%s", affectedRows, executeCostTime))

	// 返回数据
	data.ExecuteLog = strings.Join(executeLog, "\n")
	data.AffectedRows = affectedRows
	data.ExecuteCostTime = executeCostTime
	return
}

// TiDB DDL
type ExecuteTiDBDDL struct {
	*base.DBConfig
}

func (e *ExecuteTiDBDDL) Run() (data base.ReturnData, err error) {
	// Create/Drop/Truncate/Rename/Alter
	sqlType, err := parser.GetSqlStatement(e.SQL)
	if err != nil {
		return
	}

	switch sqlType {
	case "CreateDatabase", "CreateTable", "CreateView":
		return ExecuteOnlineDDL(e.DBConfig)
	case "DropTable", "DropIndex":
		return ExecuteOnlineDDL(e.DBConfig)
	case "TruncateTable":
		return ExecuteOnlineDDL(e.DBConfig)
	case "RenameTable":
		return ExecuteOnlineDDL(e.DBConfig)
	case "CreateIndex":
		return ExecuteOnlineDDL(e.DBConfig)
	case "DropDatabase":
		return data, errors.New("【风险】禁止执行drop database操作")
	case "AlterTable":
		return ExecuteOnlineDDL(e.DBConfig)
	default:
		return data, fmt.Errorf("当前SQL未匹配到规则，执行失败，SQL类型为：%s", sqlType)
	}
}

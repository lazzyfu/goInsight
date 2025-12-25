package tidb

import (
	"errors"
	"fmt"
	"time"

	"github.com/lazzyfu/goinsight/pkg/parser"
	"github.com/lazzyfu/goinsight/pkg/utils"

	"github.com/lazzyfu/goinsight/internal/orders/api/base"
)

// 执行Online DDL语句
func (e *ExecuteTiDBDDL) ExecuteOnlineDDL() (data base.ReturnData, err error) {
	logger := base.NewExecuteLogger()
	publisher := base.NewRedisPublisher()

	log := func(msg string) {
		formatted := logger.Add(msg)
		publisher.Publish(e.OrderID, formatted, "")
	}

	// CREATE A NEW DATABASE CONNECTION
	db, err := NewTiDBCnx(e.DBConfig)
	if err != nil {
		log(fmt.Sprintf("访问数据库(%s:%d)失败，错误：%s", e.Hostname, e.Port, err.Error()))
		data.ExecuteLog = logger.String()
		return data, base.SQLExecuteError{Err: err}
	}
	defer db.Close()
	log(fmt.Sprintf("访问数据库(%s:%d)成功", e.Hostname, e.Port))

	// GET CONNECTION ID
	connectionID, err := DaoGetTiDBConnectionID(db)
	if err != nil {
		log(fmt.Sprintf("获取数据库Connection ID失败，错误：%s", err.Error()))
		data.ExecuteLog = logger.String()
		return data, base.SQLExecuteError{Err: err}
	}
	log(fmt.Sprintf("数据库Connection ID：%d", connectionID))

	// SHOW PROCESS
	ch1 := make(chan int64)
	go DaoTiDBGetProcesslist(e.DBConfig, e.OrderID, connectionID, ch1)

	// 执行SQL
	startTime := time.Now()
	affectedRows, err := DaoTiDBExecute(db, e.SQL, ch1)
	if err != nil {
		log(fmt.Sprintf("SQL执行失败，错误：%s", err.Error()))
		data.ExecuteLog = logger.String()
		return data, base.SQLExecuteError{Err: err}
	}
	endTime := time.Now()
	executeCostTime := utils.HumanfriendlyTimeUnit(endTime.Sub(startTime))
	log(fmt.Sprintf("SQL执行成功，影响行数%d，执行耗时：%s", affectedRows, executeCostTime))

	// 返回数据
	data.ExecuteLog = logger.String()
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
		return e.ExecuteOnlineDDL()
	case "DropTable", "DropIndex":
		return e.ExecuteOnlineDDL()
	case "TruncateTable":
		return e.ExecuteOnlineDDL()
	case "RenameTable":
		return e.ExecuteOnlineDDL()
	case "CreateIndex":
		return e.ExecuteOnlineDDL()
	case "DropDatabase":
		return data, errors.New("【风险】禁止执行drop database操作")
	case "AlterTable":
		return e.ExecuteOnlineDDL()
	default:
		return data, fmt.Errorf("当前SQL未匹配到规则，执行失败，SQL类型为：%s", sqlType)
	}
}

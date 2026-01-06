package tidb

import (
	"fmt"
	"strings"
	"time"

	"github.com/lazzyfu/goinsight/pkg/utils"

	"github.com/lazzyfu/goinsight/internal/orders/api/base"
)

// TiDB DML
type ExecuteTiDBDML struct {
	*base.DBConfig
}

// 执行TiDB DML
func (e *ExecuteTiDBDML) Run() (data base.ReturnData, err error) {
	logger := base.NewExecuteLogger()
	publisher := base.NewRedisPublisher()

	log := func(msg string) {
		cleaned := strings.TrimSpace(msg)
		formatted := logger.Add(cleaned) + "\n"
		publisher.Publish(e.OrderID, e.TaskID, utils.RenderLogStream, formatted)
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

	// // 转换SQL
	// rw, err := base.NewRewrite(e.SQL)
	// if err != nil {
	// 	return logErrorAndReturn(base.SQLExecuteError{Err: err}, "SQL转换失败，错误：")
	// }
	// err = rw.RewriteDML2Select()
	// if err != nil {
	// 	return logErrorAndReturn(base.SQLExecuteError{Err: err}, "SQL转换失败，错误：")
	// }
	// fmt.Println("sql wa::: ", rw.SQL)

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

	// TiDB不支持生成回滚SQL
	log("TiDB不支持生成回滚SQL")

	// 返回数据
	data.ExecuteLog = logger.String()
	data.AffectedRows = affectedRows
	data.ExecuteCostTime = executeCostTime
	return
}

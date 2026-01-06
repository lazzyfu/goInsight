package mysql

import (
	"fmt"
	"strings"
	"time"

	"github.com/lazzyfu/goinsight/pkg/utils"

	"github.com/lazzyfu/goinsight/internal/orders/api/base"
)

// MySQL DML
type ExecuteMySQLDML struct {
	*base.DBConfig
}

// 执行MySQL DML语句
func (e *ExecuteMySQLDML) Run() (data base.ReturnData, err error) {
	logger := base.NewExecuteLogger()
	publisher := base.NewRedisPublisher()

	log := func(msg string) {
		cleaned := strings.TrimSpace(msg)
		formatted := logger.Add(cleaned) + "\n"
		publisher.Publish(e.OrderID, e.TaskID, utils.RenderLogStream, formatted)
	}

	// CREATE A NEW DATABASE CONNECTION
	db, err := NewMySQLCnx(e.DBConfig)
	if err != nil {
		log(fmt.Sprintf("访问数据库(%s:%d)失败，错误：%s", e.DBConfig.Hostname, e.DBConfig.Port, err.Error()))
		data.ExecuteLog = logger.String()
		return data, base.SQLExecuteError{Err: err}
	}
	defer db.Close()
	log(fmt.Sprintf("访问数据库(%s:%d)成功", e.DBConfig.Hostname, e.DBConfig.Port))

	// GET CONNECTION ID
	connectionID, err := DaoMySQLGetConnectionID(db)
	if err != nil {
		log(fmt.Sprintf("获取数据库Connection ID失败，错误：%s", err.Error()))
		data.ExecuteLog = logger.String()
		return data, base.SQLExecuteError{Err: err}
	}
	log(fmt.Sprintf("数据库Connection ID：%d", connectionID))

	// SHOW PROCESS
	ch1 := make(chan int64)
	go DaoMySQLGetProcesslist(e.DBConfig, e.OrderID, connectionID, ch1)

	// 获取执行开始前的binlog position
	startFile, startPosition, err := DaoMySQLGetBinlogPos(db)
	if err != nil {
		log(fmt.Sprintf("获取Start Binlog File和Position失败，错误：%s", err.Error()))
		data.ExecuteLog = logger.String()
		return data, base.SQLExecuteError{Err: err}
	}
	log(fmt.Sprintf("Start Binlog File：%s，Position：%d", startFile, startPosition))

	// 执行SQL
	startTime := time.Now()
	affectedRows, err := DaoMySQLExecute(db, e.SQL, ch1)
	if err != nil {
		log(fmt.Sprintf("SQL执行失败，错误：%s", err.Error()))
		data.ExecuteLog = logger.String()
		return data, base.SQLExecuteError{Err: err}
	}
	endTime := time.Now()
	executeCostTime := utils.HumanfriendlyTimeUnit(endTime.Sub(startTime))
	log(fmt.Sprintf("SQL执行成功，影响行数%d，执行耗时：%s", affectedRows, executeCostTime))

	data.AffectedRows = affectedRows
	data.ExecuteCostTime = executeCostTime

	// 获取执行后的binlog position
	endFile, endPosition, err := DaoMySQLGetBinlogPos(db)
	if err != nil {
		log(fmt.Sprintf("获取End Binlog File和Position失败，错误：%s", err.Error()))
		data.ExecuteLog = logger.String()
		return data, base.SQLExecuteError{Err: err}
	}
	log(fmt.Sprintf("End Binlog File：%s，Position：%d", endFile, endPosition))

	var rollbackSQL, backupCostTime string
	// 影响行数大于0，才执行生成回滚SQL操作
	if affectedRows > 0 {
		// 生成回滚SQL
		log("开始解析Binlog生成回滚SQL")
		startTime = time.Now()
		binlog := Binlog{
			DBConfig:      e.DBConfig,
			ConnectionID:  connectionID,
			StartFile:     startFile,
			StartPosition: startPosition,
			EndFile:       endFile,
			EndPosition:   endPosition}
		rollbackSQL, err = binlog.Run()
		if err != nil {
			log(fmt.Sprintf("生成回滚SQL失败，错误：%s", err.Error()))
			data.ExecuteLog = logger.String()
			return data, base.SQLExecuteError{Err: err}
		}

		endTime = time.Now()
		backupCostTime = utils.HumanfriendlyTimeUnit(endTime.Sub(startTime))
		log(fmt.Sprintf("生成回滚SQL成功，耗时：%s", backupCostTime))
	}
	// 返回数据
	data.ExecuteLog = logger.String()
	data.RollbackSQL = rollbackSQL
	data.BackupCostTime = backupCostTime
	return
}

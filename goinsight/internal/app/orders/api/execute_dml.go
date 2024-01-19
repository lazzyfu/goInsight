/*
@Author  :   lazzyfu
@Desc    :   执行DML
*/

package api

import (
	"goInsight/internal/pkg/utils"
	"time"
)

// MySQL
type ExecuteMySQLDML struct {
	*DBConfig
}

func (e *ExecuteMySQLDML) Run() (data ReturnData, err error) {
	// Create a new database connection
	db, err := NewMySQLCnx(e.DBConfig)
	if err != nil {
		return data, err
	}
	defer db.Close()
	// get connection id
	ConnectionID, err := DaoGetConnectionID(db)
	if err != nil {
		return data, err
	}
	// show process
	ch1 := make(chan int64)
	go DaoGetProcesslist(e.DBConfig, e.OrderID, ConnectionID, ch1)
	// 获取position
	startFile, startPosition, err := DaoGetMySQLPos(db)
	if err != nil {
		return data, err
	}
	// 发送消息
	PublishMsg(e.OrderID, "开始执行SQL", "")
	// 开始执行SQL
	startTime := time.Now()
	affectedRows, err := DaoMySQLExecute(db, e.SQL, ch1)
	if err != nil {
		return data, err
	}
	data.AffectedRows = affectedRows

	endTime := time.Now()
	executeCostTime := utils.HumanfriendlyTimeUnit(endTime.Sub(startTime))
	data.ExecuteCostTime = executeCostTime
	// 获取position
	endFile, endPosition, err := DaoGetMySQLPos(db)
	if err != nil {
		return data, err
	}
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
			return data, err
		}
		endTime = time.Now()
		backupCostTime = utils.HumanfriendlyTimeUnit(endTime.Sub(startTime))
	}
	// 返回数据
	data.RollbackSQL = rollbackSQL
	data.BackupCostTime = backupCostTime
	return
}

// TiDB
type ExecuteTiDBDML struct {
	*DBConfig
}

func (e *ExecuteTiDBDML) Run() (data ReturnData, err error) {
	// Create a new database connection
	db, err := NewMySQLCnx(e.DBConfig)
	if err != nil {
		return data, err
	}
	defer db.Close()
	// get connection id
	ConnectionID, err := DaoGetConnectionID(db)
	if err != nil {
		return data, err
	}
	// show process
	ch1 := make(chan int64)
	go DaoGetProcesslist(e.DBConfig, e.OrderID, ConnectionID, ch1)

	// 执行SQL
	startTime := time.Now()
	affectedRows, err := DaoMySQLExecute(db, e.SQL, ch1)
	if err != nil {
		return data, err
	}
	endTime := time.Now()
	executeCostTime := utils.HumanfriendlyTimeUnit(endTime.Sub(startTime))

	var backupCostTime string
	// 返回数据
	data.RollbackSQL = "TiDB不支持生成回滚SQL"
	data.AffectedRows = affectedRows
	data.ExecuteCostTime = executeCostTime
	data.BackupCostTime = backupCostTime
	return
}

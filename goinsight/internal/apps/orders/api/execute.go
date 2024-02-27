/*
@Time    :   2023/11/03 15:15:13
@Author  :   lazzyfu
*/

package api

import (
	"errors"
	"fmt"
)

type ExportFile struct {
	FileName      string `json:"file_name"`
	FileSize      int64  `json:"file_size"`
	FilePath      string `json:"file_path"`
	ContentType   string `json:"content_type"`
	EncryptionKey string `json:"encryption_key"`
	ExportRows    int64  `json:"export_rows"`
	DownloadUrl   string `json:"download_url"`
}
type ReturnData struct {
	RollbackSQL     string `json:"rollback_sql"`
	AffectedRows    int64  `json:"affected_rows"`
	ExecuteCostTime string `json:"execute_cost_time"`
	BackupCostTime  string `json:"backup_cost_time"`
	ExecuteLog      string `json:"execute_log"`
	ExportFile
	Error string `json:"error"`
}

type DBConfig struct {
	Hostname         string
	Port             uint16
	Charaset         string
	UserName         string
	Password         string
	Schema           string
	DBType           string
	SQLType          string
	SQL              string
	OrderID          string
	TaskID           string
	ExportFileFormat string
}

type Executor interface {
	Run() (ReturnData, error)
}

type ExecuteSQLAPI struct {
	*DBConfig
	Executor
}

func (e *ExecuteSQLAPI) Run() (data ReturnData, err error) {
	data, err = e.Executor.Run()
	if err != nil {
		// 发送消息
		PublishMsg(e.OrderID, err.Error(), "")
		data.Error = err.Error()
	}
	return
}

func NewExecuteSQLAPI(config *DBConfig) *ExecuteSQLAPI {
	switch config.DBType {
	case "MySQL":
		return &ExecuteSQLAPI{config, &MySQLExecutor{config}}
	case "TiDB":
		return &ExecuteSQLAPI{config, &TiDBExecutor{config}}
	case "ClickHouse":
		// todo
		// return &ExecuteSQLAPI{config, &ExecuteClickHouseDML{config}}
	}
	return nil
}

// MySQL
type MySQLExecutor struct {
	*DBConfig
}

func (m *MySQLExecutor) Run() (data ReturnData, err error) {
	switch m.SQLType {
	case "DML":
		execute := ExecuteMySQLDML{m.DBConfig}
		return execute.Run()
	case "DDL":
		execute := ExecuteMySQLDDL{m.DBConfig}
		return execute.Run()
	case "EXPORT":
		execute := ExecuteExportMySQLData{m.DBConfig}
		return execute.Run()
	default:
		data.Error = fmt.Sprintf("不支持的SQL类型：%s", m.SQLType)
		err = errors.New(data.Error)
		return
	}
}

// TiDBExecutor 结构体用于执行 TiDB 数据库操作
type TiDBExecutor struct {
	// 继承自 DBConfig 结构体，包含数据库连接信息
	*DBConfig
}

// Run 方法根据 SQL 类型执行 TiDB 操作并返回结果和错误信息
func (m *TiDBExecutor) Run() (data ReturnData, err error) {
	switch m.SQLType {
	case "DML":
		execute := ExecuteTiDBDML{m.DBConfig}
		return execute.Run()
	case "DDL":
		execute := ExecuteTiDBDDL{m.DBConfig}
		return execute.Run()
	case "EXPORT":
		execute := ExecuteExportMySQLData{m.DBConfig}
		return execute.Run()
	default:
		data.Error = fmt.Sprintf("不支持的SQL类型：%s", m.SQLType)
		err = errors.New(data.Error)
		return
	}
}

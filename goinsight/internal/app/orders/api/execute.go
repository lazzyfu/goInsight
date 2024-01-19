/*
@Time    :   2023/11/03 15:15:13
@Author  :   lazzyfu
*/

package api

type ReturnData struct {
	RollbackSQL     string `json:"rollback_sql"`
	AffectedRows    int64  `json:"affected_rows"`
	ExecuteCostTime string `json:"execute_cost_time"`
	BackupCostTime  string `json:"backup_cost_time"`
	ExecuteLog      string `json:"execute_log"`
	Error           string `json:"error"`
}

type DBConfig struct {
	Hostname string
	Port     uint16
	Charaset string
	UserName string
	Password string
	Schema   string
	DBType   string
	SQLType  string
	SQL      string
	OrderID  string
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
	default:
		data.Error = "unknown SQLType"
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
	// 当 SQL 类型为 DML 时，使用 ExecuteTiDBDML 结构体执行操作
	case "DML":
		execute := ExecuteTiDBDML{m.DBConfig}
		return execute.Run()
	// 当 SQL 类型为 DDL 时，使用 ExecuteTiDBDDL 结构体执行操作
	case "DDL":
		execute := ExecuteTiDBDDL{m.DBConfig}
		return execute.Run()
	// 当SQL类型未知时，设置错误信息并返回
	default:
		data.Error = "unknown SQLType"
		return
	}
}

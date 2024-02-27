/*
@Author  :   lazzyfu
@Desc    :   执行DML
*/

package api

import (
	"context"
	"errors"
	"fmt"
	"goInsight/global"
	"goInsight/internal/pkg/parser"
	"goInsight/internal/pkg/utils"
	"regexp"
	"strings"
	"time"
)

// 执行Online DDL语句
func ExecuteOnlineDDL(dc *DBConfig) (data ReturnData, err error) {
	var executeLog []string
	var msg string

	// Create a new database connection
	db, err := NewMySQLCnx(dc)
	if err != nil {
		data.ExecuteLog = fmt.Sprintf("访问数据库(%s:%d)失败，错误：%s", dc.Hostname, dc.Port, err.Error())
		return data, err
	}
	defer db.Close()
	msg = fmt.Sprintf("[%s] 访问数据库(%s:%d)成功", time.Now().Format("2006-01-02 15:04:05"), dc.Hostname, dc.Port)
	executeLog = append(executeLog, msg)
	PublishMsg(dc.OrderID, msg, "")

	// get connection id
	ConnectionID, err := DaoGetConnectionID(db)
	if err != nil {
		data.ExecuteLog = fmt.Sprintf("获取数据库Connection ID失败，错误：%s", err.Error())
		return data, err
	}
	msg = fmt.Sprintf("[%s] 数据库Connection ID：%d", time.Now().Format("2006-01-02 15:04:05"), ConnectionID)
	executeLog = append(executeLog, msg)
	PublishMsg(dc.OrderID, msg, "")

	// show process
	ch1 := make(chan int64)
	go DaoGetProcesslist(dc, dc.OrderID, ConnectionID, ch1)

	// 执行SQL
	startTime := time.Now()
	affectedRows, err := DaoMySQLExecute(db, dc.SQL, ch1)
	if err != nil {
		data.ExecuteLog = fmt.Sprintf("SQL执行失败，错误：%s", err.Error())
		return data, err
	}
	endTime := time.Now()
	executeCostTime := utils.HumanfriendlyTimeUnit(endTime.Sub(startTime))

	msg = fmt.Sprintf("[%s] SQL执行成功，影响行数%d，执行耗时：%s", time.Now().Format("2006-01-02 15:04:05"), affectedRows, executeCostTime)
	executeLog = append(executeLog, msg)
	PublishMsg(dc.OrderID, msg, "")
	// 返回数据
	data.ExecuteLog = strings.Join(executeLog, "\n")
	data.AffectedRows = affectedRows
	data.ExecuteCostTime = executeCostTime
	return
}

// MySQL
type ExecuteMySQLDDL struct {
	*DBConfig
}

func (e *ExecuteMySQLDDL) ExecuteCommand(command string) (data []string, err error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ch := make(chan string)
	// 读取输出
	go func(ch <-chan string) {
		for {
			v, ok := <-ch
			if !ok {
				break
			}
			data = append(data, v)
			err = utils.Publish(context.Background(), e.OrderID, v, "ghost")
			if err != nil {
				global.App.Log.Error(err)
			}
		}
	}(ch)
	err = Command(ctx, ch, command)
	if err != nil {
		return
	}
	return
}

// 执行ghost封装后的DDL
func (e *ExecuteMySQLDDL) ExecuteDDLWithGhost(sql string) (data ReturnData, err error) {
	// 移除字符串前后的所有空白字符，包括空格、制表符、换行符等
	newSQL := strings.TrimSpace(sql)
	// 获取表名
	tableName, err := parser.GetTableNameFromAlterStatement(sql)
	if err != nil {
		return data, err
	}
	// 正则匹配
	syntax := `(?i)^ALTER(\s+)TABLE(\s+)([\S]*)(\s+)(ADD|CHANGE|RENAME|MODIFY|DROP|ENGINE|CONVERT)(\s*)([\S\s]*)`
	re, err := regexp.Compile(syntax)
	if err != nil {
		return data, err
	}
	match := re.FindStringSubmatch(newSQL)
	if len(match) < 5 {
		return data, errors.New("当前SQL正则匹配失败")
	}
	//  将反引号处理为空，将双引号处理成单引号
	vv := strings.Replace(strings.Replace(strings.Join(match[5:], ""), "`", "", -1), "\"", "'", -1)
	// 生成ghost命令
	ghostCMD := strings.Join(
		[]string{
			global.App.Config.Ghost.Path,
			strings.Join(global.App.Config.Ghost.Args, " "),
			fmt.Sprintf("--user=\"%s\" --password=\"%s\"", global.App.Config.RemoteDB.UserName, global.App.Config.RemoteDB.Password),
			fmt.Sprintf("--host=\"%s\" --port=%d", e.Hostname, e.Port),
			fmt.Sprintf("--database=%s --table=%s", e.Schema, tableName),
			fmt.Sprintf("--alter=\"%s\" --execute", vv),
		}, " ")
	startTime := time.Now()
	// 执行ghost命令
	executeLog, err := e.ExecuteCommand(ghostCMD)
	endTime := time.Now()
	executeCostTime := utils.HumanfriendlyTimeUnit(endTime.Sub(startTime))
	// 返回数据
	data.ExecuteLog = strings.Join(executeLog, "")
	data.ExecuteCostTime = executeCostTime
	if err != nil {
		data.Error = err.Error()
	}
	if err != nil {
		return
	}
	return
}

func (e *ExecuteMySQLDDL) Run() (data ReturnData, err error) {
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
		return data, errors.New("请更正语法为alter table ... rename")
	case "CreateIndex":
		return data, errors.New("请更正语法为alter table ... add")
	case "DropDatabase":
		return data, errors.New("【风险】禁止执行drop database操作")
	case "AlterTable":
		return e.ExecuteDDLWithGhost(e.SQL)
	default:
		return data, fmt.Errorf("当前SQL未匹配到规则，执行失败，SQL类型为：%s", sqlType)
	}
}

// TiDB
type ExecuteTiDBDDL struct {
	*DBConfig
}

func (e *ExecuteTiDBDDL) Run() (data ReturnData, err error) {
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

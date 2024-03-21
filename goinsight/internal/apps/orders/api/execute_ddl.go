/*
@Author  :   lazzyfu
@Desc    :   执行DDL
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

	// Function to log messages and publish
	logAndPublish := func(msg string) {
		timestamp := time.Now().Format("2006-01-02 15:04:05")
		formattedMsg := fmt.Sprintf("[%s] %s\n", timestamp, msg)
		executeLog = append(executeLog, formattedMsg)
		PublishMsg(dc.OrderID, formattedMsg, "ghost")
	}

	// Logging function for errors
	logErrorAndReturn := func(err error, errMsg string) (ReturnData, error) {
		logAndPublish(errMsg + err.Error())
		data.ExecuteLog = strings.Join(executeLog, "\n")
		return data, err
	}

	// CREATE A NEW DATABASE CONNECTION
	db, err := NewMySQLCnx(dc)
	if err != nil {
		return logErrorAndReturn(SQLExecuteError{Err: err}, fmt.Sprintf("访问数据库(%s:%d)失败，错误：%s", dc.Hostname, dc.Port, err.Error()))
	}
	defer db.Close()
	logAndPublish(fmt.Sprintf("访问数据库(%s:%d)成功", dc.Hostname, dc.Port))

	// GET CONNECTION ID
	connectionID, err := DaoGetConnectionID(db)
	if err != nil {
		return logErrorAndReturn(SQLExecuteError{Err: err}, "获取数据库Connection ID失败，错误：")
	}
	logAndPublish(fmt.Sprintf("数据库Connection ID：%d", connectionID))

	// SHOW PROCESS
	ch1 := make(chan int64)
	go DaoGetProcesslist(dc, dc.OrderID, connectionID, ch1)

	// 执行SQL
	startTime := time.Now()
	affectedRows, err := DaoMySQLExecute(db, dc.SQL, ch1)
	if err != nil {
		return logErrorAndReturn(SQLExecuteError{Err: err}, "SQL执行失败，错误：")
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

// MySQL DDL
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
	var executeLog []string

	// Function to log messages and publish
	logAndPublish := func(msg string) {
		timestamp := time.Now().Format("2006-01-02 15:04:05")
		formattedMsg := fmt.Sprintf("[%s] %s\n", timestamp, msg)
		executeLog = append(executeLog, formattedMsg)
		PublishMsg(e.OrderID, formattedMsg, "ghost")
	}

	// Logging function for errors
	logErrorAndReturn := func(err error, errMsg string) (ReturnData, error) {
		logAndPublish(errMsg + err.Error())
		data.ExecuteLog = strings.Join(executeLog, "\n")
		return data, err
	}

	// 移除字符串前后的所有空白字符，包括空格、制表符、换行符等
	newSQL := strings.TrimSpace(sql)
	logAndPublish("移除SQL语句前后的所有空白字符，包括空格、制表符、换行符等")

	// 获取表名
	tableName, err := parser.GetTableNameFromAlterStatement(sql)
	if err != nil {
		return logErrorAndReturn(SQLExecuteError{Err: err}, "解析SQL提取表名失败")
	}
	logAndPublish("从SQL语句中提取表名成功")

	// 正则匹配
	syntax := `(?i)^ALTER(\s+)TABLE(\s+)([\S]*)(\s+)(ADD|CHANGE|RENAME|MODIFY|DROP|ENGINE|CONVERT)(\s*)([\S\s]*)`
	re, err := regexp.Compile(syntax)
	if err != nil {
		return logErrorAndReturn(SQLExecuteError{Err: err}, "正则匹配SQL语句失败")
	}
	match := re.FindStringSubmatch(newSQL)
	if len(match) < 5 {
		return logErrorAndReturn(SQLExecuteError{Err: err}, "正则匹配SQL语句失败")
	}
	logAndPublish("正则匹配SQL语句成功")

	//  将反引号处理为空，将双引号处理成单引号
	vv := strings.Replace(strings.Replace(strings.Join(match[5:], ""), "`", "", -1), "\"", "'", -1)
	logAndPublish("将反引号处理为空，将双引号处理成单引号")

	// 生成ghost命令
	logAndPublish("生成gh-ost执行命令")
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
	// 打印命令，已掩码password
	re = regexp.MustCompile(`--password="([^"]*)"`)
	printGhostCMD := re.ReplaceAllString(ghostCMD, `--password="..."`)
	logAndPublish(fmt.Sprintf("执行gh-ost命令：%s", printGhostCMD))
	// 执行ghost命令
	log, err := e.ExecuteCommand(ghostCMD)
	executeLog = append(executeLog, log...)
	if err != nil {
		return logErrorAndReturn(SQLExecuteError{Err: err}, "执行失败，错误：")
	}
	logAndPublish("gh-ost命令执行成功")
	endTime := time.Now()
	executeCostTime := utils.HumanfriendlyTimeUnit(endTime.Sub(startTime))

	// 返回数据
	data.ExecuteLog = strings.Join(executeLog, "")
	data.ExecuteCostTime = executeCostTime
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
		return data, errors.New("请更正为alter table ... rename语法")
	case "CreateIndex":
		return data, errors.New("请更正为alter table ... add语法")
	case "DropDatabase":
		return data, errors.New("【风险】禁止执行drop database操作")
	case "AlterTable":
		return e.ExecuteDDLWithGhost(e.SQL)
	default:
		return data, fmt.Errorf("当前SQL未匹配到规则，执行失败，SQL类型为：%s", sqlType)
	}
}

// TiDB DDL
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

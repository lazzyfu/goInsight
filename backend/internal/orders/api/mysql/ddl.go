package mysql

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/lazzyfu/goinsight/internal/global"

	"github.com/lazzyfu/goinsight/pkg/parser"
	"github.com/lazzyfu/goinsight/pkg/utils"

	"github.com/lazzyfu/goinsight/internal/orders/api/base"
)

// MySQL DDL
type ExecuteMySQLDDL struct {
	*base.DBConfig
}

func (e *ExecuteMySQLDDL) ExecuteCommand(command string) (data []string, err error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch := make(chan string, 100)

	// 读取输出 & 推送
	done := make(chan struct{})
	go func() {
		defer close(done)
		for v := range ch {
			data = append(data, v)

			if perr := utils.Publish(context.Background(), e.OrderID, e.TaskID, utils.RenderGhostStream, v); perr != nil {
				global.App.Log.Error(perr)
			}
		}
	}()

	err = base.Command(ctx, ch, command)
	close(ch) // 明确关闭
	<-done    // 等待消费完成

	return
}

// 执行Online DDL语句
func (e *ExecuteMySQLDDL) ExecuteOnlineDDL() (data base.ReturnData, err error) {
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
		log(fmt.Sprintf("访问数据库(%s:%d)失败，错误：%s", e.Hostname, e.Port, err.Error()))
		data.ExecuteLog = logger.String()
		return data, base.SQLExecuteError{Err: err}
	}
	defer db.Close()

	log(fmt.Sprintf("访问数据库(%s:%d)成功", e.Hostname, e.Port))

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

	// 返回数据
	data.ExecuteLog = logger.String()
	data.AffectedRows = affectedRows
	data.ExecuteCostTime = executeCostTime
	return
}

// 执行ghost封装后的DDL
func (e *ExecuteMySQLDDL) ExecuteDDLWithGhost(sql string) (data base.ReturnData, err error) {
	logger := base.NewExecuteLogger()
	publisher := base.NewRedisPublisher()

	log := func(msg string) {
		cleaned := strings.TrimSpace(msg)
		formatted := logger.Add(cleaned) + "\n"
		publisher.Publish(e.OrderID, e.TaskID, utils.RenderGhostStream, formatted)
	}

	// 移除字符串前后的所有空白字符，包括空格、制表符、换行符等
	newSQL := strings.TrimSpace(sql)
	log("移除SQL语句前后的所有空白字符，包括空格、制表符、换行符等")

	// 获取表名
	tableName, err := parser.GetTableNameFromAlterStatement(sql)
	if err != nil {
		log(fmt.Sprintf("解析SQL提取表名失败，错误：%s", err.Error()))
		data.ExecuteLog = logger.String()
		return data, base.SQLExecuteError{Err: err}
	}
	log("从SQL语句中提取表名成功")

	// 正则匹配
	syntax := `(?i)^ALTER(\s+)TABLE(\s+)([\S]*)(\s+)(ADD|CHANGE|RENAME|MODIFY|DROP|ENGINE|CONVERT)(\s*)([\S\s]*)`
	re, err := regexp.Compile(syntax)
	if err != nil {
		log(fmt.Sprintf("正则匹配SQL语句失败，错误：%s", err.Error()))
		data.ExecuteLog = logger.String()
		return data, base.SQLExecuteError{Err: err}
	}
	match := re.FindStringSubmatch(newSQL)
	if len(match) < 5 {
		log("正则匹配SQL语句失败，未匹配到预期结果")
		data.ExecuteLog = logger.String()
		return data, base.SQLExecuteError{Err: errors.New("正则匹配SQL语句失败，未匹配到预期结果")}
	}
	log("正则匹配SQL语句成功")

	//  将反引号处理为空，将双引号处理成单引号
	vv := strings.ReplaceAll(strings.ReplaceAll(strings.Join(match[5:], ""), "`", ""), "\"", "'")
	log("将反引号处理为空，将双引号处理成单引号")

	// 生成ghost命令
	log("生成gh-ost执行命令")

	ghostCMDParts := []string{
		global.App.Config.Ghost.Path,
		strings.Join(global.App.Config.Ghost.Args, " "),
		fmt.Sprintf("--user=\"%s\" --password=\"%s\"", e.UserName, e.Password),
		fmt.Sprintf("--host=\"%s\" --port=%d", e.Hostname, e.Port),
		fmt.Sprintf("--database=%s --table=%s", e.Schema, tableName),
		fmt.Sprintf("--alter=\"%s\" --execute", vv),
	}

	if strings.Contains(e.Hostname, "rds.aliyuncs.com") {
		ghostCMDParts = append(ghostCMDParts, "-aliyun-rds=true")
		ghostCMDParts = append(ghostCMDParts, fmt.Sprintf("-assume-master-host=\"%s\"", e.Hostname))
	}

	ghostCMD := strings.Join(ghostCMDParts, " ")

	startTime := time.Now()
	// 打印命令，已掩码password
	re = regexp.MustCompile(`--password="([^"]*)"`)
	printGhostCMD := re.ReplaceAllString(ghostCMD, `--password="..."`)
	log(fmt.Sprintf("执行gh-ost命令：%s", printGhostCMD))
	// 执行ghost命令
	elog, eerr := e.ExecuteCommand(ghostCMD)
	for _, l := range elog {
		log(l)
	}
	if eerr != nil {
		log(fmt.Sprintf("gh-ost命令执行失败，错误：%s", eerr.Error()))
		data.ExecuteLog = logger.String()
		return data, base.SQLExecuteError{Err: eerr}
	}
	log("gh-ost命令执行成功")
	endTime := time.Now()
	executeCostTime := utils.HumanfriendlyTimeUnit(endTime.Sub(startTime))

	// 返回数据
	data.ExecuteLog = logger.String()
	data.ExecuteCostTime = executeCostTime
	return
}

func (e *ExecuteMySQLDDL) Run() (data base.ReturnData, err error) {
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

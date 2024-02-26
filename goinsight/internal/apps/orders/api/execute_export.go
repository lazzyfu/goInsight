/*
@Author  :   lazzyfu
@Desc    :
*/

package api

import (
	"database/sql"
	"fmt"
	"goInsight/global"
	"goInsight/internal/pkg/utils"
	"strconv"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
)

// MySQL
type ExecuteExportMySQLData struct {
	*DBConfig
}

// 通过游标读取数据流式写入到文件，避免海量数据打爆内存
func (e *ExecuteExportMySQLData) toExcel(db *sql.DB, query string) (file ExportFile, err error) {
	// Execute the SQL query
	rows, err := db.Query(query)
	if err != nil {
		global.App.Log.Error("execute sql query error", err.Error())
		return
	}
	defer rows.Close()
	// 创建Excel文件
	sheet := "Sheet1"
	f := excelize.NewFile()
	// 创建流式写入器
	sw, err := f.NewStreamWriter(sheet)
	if err != nil {
		global.App.Log.Error("new excel stream writer error", err.Error())
		return
	}
	// 获取列名并写入表头
	columns, err := rows.Columns()
	if err != nil {
		global.App.Log.Error("get columns error", err.Error())
		return
	}
	var columnsInterface []interface{}
	for _, value := range columns {
		columnsInterface = append(columnsInterface, value)
	}
	if err = sw.SetRow("A1", columnsInterface); err != nil {
		global.App.Log.Error("set row error for columns", err.Error())
		return
	}
	// excel起始行号
	var rowIndex int64 = 2
	// 读取数据
	vals := make([]interface{}, len(columns))
	for i := range columns {
		vals[i] = new(sql.RawBytes)
	}
	for rows.Next() {
		if err = rows.Scan(vals...); err != nil {
			return
		}
		// 处理数据
		vmap := make([]interface{}, len(vals))
		for i, c := range vals {
			// Type assertion and value conversion
			switch v := c.(type) {
			case *sql.RawBytes:
				if *v == nil {
					// Handle null value as nil
					vmap[i] = nil
				} else {
					// Convert RawBytes to string
					vmap[i] = string(*v)
				}
			}
		}
		// 写入Excel
		if err = sw.SetRow("A"+strconv.Itoa(int(rowIndex)), vmap); err != nil {
			global.App.Log.Error("Error writing row to excel", err.Error())
			return
		}
		rowIndex++
	}
	// Check for errors after closing the rows
	if err = rows.Close(); err != nil {
		global.App.Log.Error("Error closing rows", err.Error())
		return
	}
	// Check for any additional errors
	if err = rows.Err(); err != nil {
		global.App.Log.Error("Error reading rows", err.Error())
		return
	}
	// 调用Flush函数来结束流式写入过程
	if err = sw.Flush(); err != nil {
		global.App.Log.Error("Flush failed", err.Error())
		return
	}
	// 保存Excel文件
	fileName := fmt.Sprintf("%s.xlsx", e.TaskID)
	filePath := "./media/"
	if err = f.SaveAs(filePath + fileName); err != nil {
		global.App.Log.Error("Error saving excel file", err.Error())
		return
	}
	// 加密和压缩文件
	encryptFileName := fileName + ".zip"
	encryptFilePath := filePath + encryptFileName
	key, err := e.encryptAndZipFile(fileName, encryptFileName, filePath)
	if err != nil {
		global.App.Log.Error("Error encrypting and zipping file", err.Error())
		return
	}
	// 保存信息
	file.EncryptionKey = string(key)
	file.FileName = encryptFileName
	file.FilePath = encryptFilePath
	file.ContentType = "xlsx"
	file.FileSize, _ = utils.GetFileSize(file.FilePath)
	file.ExportRows = rowIndex - 2
	file.DownloadUrl = fmt.Sprintf("%s/orders/download/exportfile/%s", global.App.Config.Notify.NoticeURL, file.FileName)
	return
}

func (e *ExecuteExportMySQLData) encryptAndZipFile(inputFile, outputFile, filePath string) (key string, err error) {
	key = utils.GenerateRandomString(24)
	return key, utils.EncryptAndTarGzFiles(inputFile, outputFile, filePath, key)
}

func (e *ExecuteExportMySQLData) Run() (data ReturnData, err error) {
	var executeLog []string
	var msg string

	// Create a new database connection
	db, err := NewMySQLCnx(e.DBConfig)
	if err != nil {
		data.ExecuteLog = fmt.Sprintf("访问数据库(%s:%d)失败，错误：%s", e.DBConfig.Hostname, e.DBConfig.Port, err.Error())
		return data, err
	}
	defer db.Close()
	msg = fmt.Sprintf("[%s] 访问数据库(%s:%d)成功", time.Now().Format("2006-01-02 15:04:05"), e.DBConfig.Hostname, e.DBConfig.Port)
	executeLog = append(executeLog, msg)
	PublishMsg(e.OrderID, msg, "")

	// 执行
	startTime := time.Now()
	file, err := e.toExcel(db, e.SQL)
	if err != nil {
		data.ExecuteLog = fmt.Sprintf("SQL执行失败，错误：%s", err.Error())
		return data, err
	}
	endTime := time.Now()
	executeCostTime := utils.HumanfriendlyTimeUnit(endTime.Sub(startTime))

	msg = fmt.Sprintf("[%s] SQL执行成功，影响行数%d，执行耗时：%s", time.Now().Format("2006-01-02 15:04:05"), file.ExportRows, executeCostTime)
	executeLog = append(executeLog, msg)
	PublishMsg(e.OrderID, msg, "")

	// 存储返回数据
	data.ExecuteLog = strings.Join(executeLog, "\n")
	data.ExecuteCostTime = executeCostTime
	data.ExportFile = file
	return
}

package mysql

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/lazzyfu/goinsight/internal/global"

	"github.com/lazzyfu/goinsight/pkg/utils"

	"github.com/lazzyfu/goinsight/internal/orders/api/base"
	"github.com/lazzyfu/goinsight/internal/orders/api/file"

	"golang.org/x/sync/errgroup"
)

type ExecuteMySQLExportToFile struct {
	*base.DBConfig
}

var filePath string = "./media"

func (e *ExecuteMySQLExportToFile) processRowsAndExport(rows *sql.Rows, columns []string) (int64, string, string, error) {
	g := new(errgroup.Group)

	vals := make([]any, len(columns))
	for i := range columns {
		vals[i] = new(sql.RawBytes)
	}
	ch := make(chan []any, 10)

	// Determine file format and start export
	var fileName, fullFileName string
	switch e.ExportFileFormat {
	case "CSV":
		fileName = fmt.Sprintf("%s.csv", e.TaskID)
		fullFileName = fmt.Sprintf("%s/%s", filePath, fileName)
		g.Go(func() error {
			return file.ToCSV(columns, ch, fullFileName)
		})
	case "XLSX":
		fileName = fmt.Sprintf("%s.xlsx", e.TaskID)
		fullFileName = fmt.Sprintf("%s/%s", filePath, fileName)
		g.Go(func() error {
			return file.ToExcel(columns, ch, fullFileName)
		})
	}

	// Read and process rows
	rowCount, err := e.readAndProcessRows(rows, vals, ch)
	if err != nil {
		return 0, "", "", err
	}

	// Close channel and wait for export to finish
	close(ch)

	// Wait to complete
	if err := g.Wait(); err != nil {
		return 0, "", "", err
	}

	return rowCount, fileName, fullFileName, nil
}

func (e *ExecuteMySQLExportToFile) readAndProcessRows(rows *sql.Rows, vals []any, ch chan []any) (int64, error) {
	var rowCount int64
	for rows.Next() {
		if err := rows.Scan(vals...); err != nil {
			return 0, err
		}
		vmap := e.processRowData(vals)
		ch <- vmap
		rowCount++
	}
	return rowCount, nil
}

func (e *ExecuteMySQLExportToFile) processRowData(vals []any) []any {
	vmap := make([]any, len(vals))
	for i, c := range vals {
		switch v := c.(type) {
		case *sql.RawBytes:
			if *v == nil {
				vmap[i] = nil
			} else {
				vmap[i] = string(*v)
			}
		}
	}
	return vmap
}

func (e *ExecuteMySQLExportToFile) Run() (data base.ReturnData, err error) {
	var (
		executeLog []string
		startTime  = time.Now()
	)

	// Function to log messages and publish
	logAndPublish := func(msg string) {
		formattedMsg := fmt.Sprintf("[%s] %s", time.Now().Format("2006-01-02 15:04:05"), msg)
		executeLog = append(executeLog, formattedMsg)
		base.PublishMessageToChannel(e.OrderID, fmt.Sprintf("%s \n", formattedMsg), "")
	}

	// Logging function for errors
	logErrorAndReturn := func(err error, errMsg string) (base.ReturnData, error) {
		logAndPublish(fmt.Sprintf("%s: %s", errMsg, err.Error()))
		data.ExecuteLog = strings.Join(executeLog, "\n")
		return data, err
	}

	// Establish database connection
	db, err := NewMySQLCnx(e.DBConfig)
	if err != nil {
		return logErrorAndReturn(base.SQLExecuteError{Err: err}, fmt.Sprintf("Failed to access database (%s:%d)", e.DBConfig.Hostname, e.DBConfig.Port))
	}
	defer db.Close()
	logAndPublish(fmt.Sprintf("连接到数据库(%s:%d)", e.DBConfig.Hostname, e.DBConfig.Port))

	// Execute the SQL
	rows, err := db.Query(e.SQL)
	if err != nil {
		return logErrorAndReturn(base.SQLExecuteError{Err: err}, "Failed to execute SQL query")
	}
	defer rows.Close()
	logAndPublish("执行SQL语句")

	// Retrieve column names
	columns, err := rows.Columns()
	if err != nil {
		return logErrorAndReturn(base.SQLExecuteError{Err: err}, "Failed to retrieve columns")
	}
	logAndPublish("检索列名")

	// Process rows and export to file
	rowCount, fileName, fullFileName, err := e.processRowsAndExport(rows, columns)
	if err != nil {
		return logErrorAndReturn(err, "Failed to process rows and export")
	}
	logAndPublish("处理行数据")

	// 加密和压缩文件
	encryptFileName := fileName + ".zip"
	encryptFilePath := fmt.Sprintf("%s/%s", filePath, encryptFileName)
	key := utils.GenerateRandomString(24)
	_ = utils.EncryptAndTarGzFiles(fileName, encryptFileName, filePath, key)
	logAndPublish(fmt.Sprintf("加密和压缩文件%s -> %s", fullFileName, encryptFileName))

	// 删除原文件
	err = os.Remove(fullFileName)
	if err != nil {
		return logErrorAndReturn(err, "Error deleting original file")
	}
	logAndPublish(fmt.Sprintf("源文件%s删除成功", fullFileName))

	// 结束时间
	endTime := time.Now()
	executeCostTime := utils.HumanfriendlyTimeUnit(endTime.Sub(startTime))

	logAndPublish(fmt.Sprintf("执行成功，影响行数%d，执行耗时：%s", rowCount, executeCostTime))

	// Prepare export file metadata
	FileSize, _ := utils.GetFileSize(encryptFilePath)
	data.ExportFile = base.ExportFile{
		EncryptionKey: string(key),
		FileName:      encryptFileName,
		FilePath:      encryptFilePath,
		ContentType:   "xlsx",
		FileSize:      FileSize,
		ExportRows:    rowCount,
		DownloadUrl:   fmt.Sprintf("%s/orders/tasks/download/exportfile/%s", global.App.Config.Notify.NoticeURL, encryptFileName),
	}
	data.ExecuteLog = strings.Join(executeLog, "\n")
	data.ExecuteCostTime = executeCostTime

	// Return execution data
	return
}

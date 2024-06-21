package file

import (
	"database/sql"
	"goInsight/global"
	"strconv"

	"github.com/xuri/excelize/v2"
)

func ToExcel(columns []string, row <-chan []interface{}, fileName string) (err error) {
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

	// 从通道中读取数据行并写入 CSV 文件
	for v := range row {
		// 写入Excel
		if err = sw.SetRow("A"+strconv.Itoa(int(rowIndex)), v); err != nil {
			global.App.Log.Error("Error writing row to excel", err.Error())
			return
		}
		rowIndex++
	}

	// 调用Flush函数来结束流式写入过程
	if err = sw.Flush(); err != nil {
		global.App.Log.Error("Flush failed", err.Error())
		return
	}

	// 保存Excel文件
	if err = f.SaveAs(fileName); err != nil {
		global.App.Log.Error("Error saving excel file", err.Error())
		return
	}
	return
}

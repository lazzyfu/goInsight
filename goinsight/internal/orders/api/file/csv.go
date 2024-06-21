package file

import (
	"encoding/csv"
	"goInsight/global"
	"os"

	"github.com/spf13/cast"
)

func ToCSV(columns []string, row <-chan []interface{}, fileName string) (err error) {
	// 创建 CSV 文件
	csvFile, err := os.Create(fileName)
	if err != nil {
		global.App.Log.Error("create csv file error", err.Error())
		return
	}
	defer csvFile.Close()

	// 创建 CSV Writer
	csvWriter := csv.NewWriter(csvFile)
	defer csvWriter.Flush()

	// 写入表头
	if err = csvWriter.Write(columns); err != nil {
		global.App.Log.Error("write columns error", err.Error())
		return
	}

	// csv起始行号
	var rowIndex int64 = 2

	// 从通道中读取数据行并写入 CSV 文件
	for v := range row {
		// 按行写入到文件
		newV := cast.ToStringSlice(v)
		if err = csvWriter.Write(newV); err != nil {
			global.App.Log.Error("write row error", err.Error())
			return
		}

		// 将缓冲区中的数据写入文件
		csvWriter.Flush()

		// 检查错误，如果有错误则打印并退出
		if err = csvWriter.Error(); err != nil {
			global.App.Log.Error("csv writer error", err.Error())
			return
		}
		rowIndex++
	}
	return
}

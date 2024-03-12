package process

import (
	"fmt"
	"goInsight/global"
	"goInsight/internal/pkg/kv"
	"strings"

	"github.com/jinzhu/copier"
)

// getRowFormatMaxSize
func getRowFormatMaxSize(rowFormat string) int {
	rowFormatMap := map[string]int{
		"REDUNDANT":  8000,
		"DYNAMIC":    65535,
		"COMPRESSED": 15360,
		"COMPACT":    8126,
	}

	return rowFormatMap[strings.ToUpper(rowFormat)]
}

// RowSizeTooLarge
type PartSpecification struct {
	Column  string
	Tp      byte
	Elems   []string // Elems is the element list for enum and set type.
	Flen    int      // 字段长度
	Decimal int      // decimal字段专用,decimal(12,2)中的2
	Charset string   // 列字符集
}
type InnoDBRowSize struct {
	Table     string // 表名
	Engine    string // 表引擎
	Charset   string // 表字符集
	RowFormat string // 行格式
	ColsMaps  []PartSpecification
}

// https://dev.mysql.com/doc/refman/8.3/en/innodb-row-format.html
func (l *InnoDBRowSize) Check(kv *kv.KVCache) error {
	if l.Engine != "InnoDB" {
		return nil
	}
	// 判断行格式
	var rowFormat string = l.RowFormat
	if l.RowFormat == "DEFAULT" {
		rowFormat = kv.Get("innodbDefaultRowFormat").(string)
	}
	maxRowSize := getRowFormatMaxSize(rowFormat)

	// version
	versionIns := DbVersion{kv.Get("dbVersion").(string)}

	// 计算列长度
	var maxSumRowsLength int

	// 判断字符集，当列字符集为空，使用表的字符集
	for _, i := range l.ColsMaps {
		// &{{riskcontrol_derived_variable_conf1 utf8mb4 [{i_id 3 [] 11 -1 } {ch_code 15 [] 200 -1 }]}}
		// 处理字符集为空的情况
		if len(i.Charset) == 0 {
			i.Charset = l.Charset
		}

		var instDataBytes DataBytes
		err := copier.CopyWithOption(&instDataBytes, i, copier.Option{IgnoreEmpty: true, DeepCopy: true})
		if err != nil {
			return err
		}
		maxSumRowsLength += instDataBytes.Get(versionIns.Int())
	}
	// 判断是否触发了行大小限制
	msg := fmt.Sprintf("表`%s`触发了Row Size Limit，最大行大小为%d，当前为%d（表存储引擎为%s，行格式为%s）", l.Table, maxRowSize, maxSumRowsLength, l.Engine, rowFormat)
	global.App.Log.Info(msg)

	if maxSumRowsLength > maxRowSize {
		return fmt.Errorf(msg)
	}

	return nil
}

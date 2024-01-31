package process

import (
	"fmt"
	"goInsight/global"
	"goInsight/internal/pkg/kv"

	"github.com/jinzhu/copier"
)

// RowSizeTooLarge
type RowSizeTooLargePartSpecification struct {
	Column  string
	Tp      byte
	Elems   []string // Elems is the element list for enum and set type.
	Flen    int      // 字段长度
	Decimal int      // decimal字段专用,decimal(12,2)中的2
	Charset string   // 列字符集
}
type RowSizeTooLarge struct {
	Table                   string // 表名
	Charset                 string // 表字符集
	RowSizeTooLargeColsMaps []RowSizeTooLargePartSpecification
}

func (l *RowSizeTooLarge) Check(kv *kv.KVCache) error {
	maxRowSize := 65535
	versionIns := DbVersion{kv.Get("dbVersion").(string)}
	var maxSumLength int
	for _, i := range l.RowSizeTooLargeColsMaps {
		// &{{riskcontrol_derived_variable_conf1 utf8mb4 [{i_id 3 [] 11 -1 } {ch_code 15 [] 200 -1 }]}}
		// 判断字符集，当列字符集为空，使用表的字符集
		if len(i.Charset) == 0 {
			i.Charset = l.Charset
		}

		var instDataBytes DataBytes
		err := copier.CopyWithOption(&instDataBytes, i, copier.Option{IgnoreEmpty: true, DeepCopy: true})
		if err != nil {
			return err
		}
		maxSumLength += instDataBytes.Get(versionIns.Int())
	}
	global.App.Log.Debug(fmt.Sprintf("maxSumLength:%d, maxRowSize:%d", maxSumLength, maxRowSize))
	if maxSumLength > maxRowSize {
		return fmt.Errorf("表`%s`触发了Row Size Limit,The maximum row size is 65535，当前为%d字节，您可以将一些列更改为TEXT类型(参考:https://dev.mysql.com/doc/refman/5.7/en/column-count-limit.html)", l.Table, maxSumLength)
	}
	return nil
}

/*
@Desc    :   下面仅为粗略计算
MySQL紧凑行格式和动态行格式的表都有65535字节最大行限制（BLOB和TEXT仅对行大小贡献9到12字节，内容是分开存储的），但是紧凑行格式的表比较特殊，额外还有8126字节最大行限制（单列按最大不超过768字节计算，如果小于768，则以真实长度进行计算）。紧凑行格式最大行限制原理分析如下（下面不考虑列定义的NULL、NOT NULL占用字节的因素）：
当表为紧凑行格式时，对于可变长的列，当前大于等于768字节时，超出部分数据存储在页外。对于固定长度的列，当大于等于768字节时，固定长度的列会被编码为可变长度列。例如varchar(300)在utf8mb4字符集下，列字节长度为300*4=1200，大于768字节，当插入大于768字节的数据时，多余的数据会存储在单独的页中。因此就有2种计算方式：①768字节会作为8126字节最大行限制一部分进行计算（如果小于768，则以真实长度进行计算）。②1200字节会作为65535字节最大行限制一部分进行计算（以真实长度进行计算）。
1）如果先触发8126字节最大行限制，则抛出：ERROR 1118 (42000): Row size too large (> 8126).
测试用例：
CREATE TABLE tab1 (
	id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
	col1 varchar(192) NOT NULL default '',
	col2 varchar(192) NOT NULL default '',
	col3 varchar(192) NOT NULL default '',
	col4 varchar(192) NOT NULL default '',
	col5 varchar(192) NOT NULL default '',
	col6 varchar(192) NOT NULL default '',
	col7 varchar(192) NOT NULL default '',
	col8 varchar(192) NOT NULL default '',
	col9 varchar(192) NOT NULL default '',
	col10 varchar(1000) NOT NULL default '',
	col11 varchar(1000) NOT NULL default '',
	PRIMARY KEY (id)
) ENGINE = InnoDB row_format = compact DEFAULT CHARSET = utf8mb4;

（2）如果先触发65535字节最大行限制，则抛出：ERROR 1118 (42000): Row size too large.
测试用例：
CREATE TABLE tab1 (
	id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
	col1 varchar(10000) NOT NULL default '',
	col2 varchar(10000) NOT NULL default '',
	PRIMARY KEY (id)
) ENGINE = InnoDB row_format = compact DEFAULT CHARSET = utf8mb4;
*/

package process

import (
	"fmt"
	"strings"

	"github.com/lazzyfu/goinsight/pkg/kv"

	"github.com/jinzhu/copier"
	"github.com/pingcap/tidb/pkg/parser/mysql"
)

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
func (rs *InnoDBRowSize) Check(kv *kv.KVCache) error {
	if rs.Engine != "InnoDB" {
		return nil
	}

	// 获取指定行格式行大小
	var rowFormat string = rs.RowFormat
	if rowFormat == "DEFAULT" {
		rowFormat = strings.ToUpper(kv.Get("innodbDefaultRowFormat").(string))
	}

	// MySQL表的内部具有65,535字节的最大行大小限制
	maxRowSize := 65535

	// 每种行格式独有的行大小限制
	var specifyRowFormatSize int
	switch rowFormat {
	case "REDUNDANT":
		specifyRowFormatSize = 8000
	case "COMPACT":
		specifyRowFormatSize = 8126
	case "COMPRESSED":
		maxRowSize = 15360
	}

	// version
	versionIns := DbVersion{kv.Get("dbVersion").(string)}

	// 计算列长度
	var totalRowLength int
	var totalSpecifiedRowLength int

	// 判断字符集，当列字符集为空，使用表的字符集
	for _, colSpec := range rs.ColsMaps {
		// &{{riskcontrol_derived_variable_conf1 utf8mb4 [{i_id 3 [] 11 -1 } {ch_code 15 [] 200 -1 }]}}
		// 处理字符集为空的情况
		if len(colSpec.Charset) == 0 {
			colSpec.Charset = rs.Charset
		}

		var instDataBytes DataBytes
		err := copier.CopyWithOption(&instDataBytes, colSpec, copier.Option{IgnoreEmpty: true, DeepCopy: true})
		if err != nil {
			return err
		}
		// 每个列计算后的长度
		colLen := instDataBytes.Get(versionIns.Int())

		// 判断行格式
		if (rowFormat == "REDUNDANT" || rowFormat == "COMPACT") && (colSpec.Tp == mysql.TypeString || colSpec.Tp == mysql.TypeVarchar) {
			totalSpecifiedRowLength += min(colLen, 768)
		}
		totalRowLength += colLen
	}

	if (strings.ToUpper(rowFormat) == "REDUNDANT" || strings.ToUpper(rowFormat) == "COMPACT") && totalSpecifiedRowLength > specifyRowFormatSize {
		return fmt.Errorf("表`%s`触发了Row Size Limit（>%d），当前为%d（表存储引擎为%s，表行格式为%s）[请参考：https://dev.mysql.com/doc/refman/8.0/en/innodb-row-format.html]", rs.Table, specifyRowFormatSize, totalSpecifiedRowLength, rs.Engine, rowFormat)
	}

	if totalRowLength > maxRowSize {
		return fmt.Errorf("表`%s`触发了Row Size Limit（>%d），当前为%d（表存储引擎为%s, 表行格式为%s）[请参考：https://dev.mysql.com/doc/refman/8.0/en/innodb-row-format.html]", rs.Table, maxRowSize, totalRowLength, rs.Engine, rowFormat)
	}
	return nil
}

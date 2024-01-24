/*
@Time    :   2022/07/06 10:12:48
@Author  :   zongfei.fu
@Desc    :   None
*/

package process

import (
	"bytes"
	"fmt"
	"goInsight/internal/pkg/kv"
	"goInsight/internal/pkg/utils"
	"sqlSyntaxAudit/config"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/pingcap/tidb/parser/mysql"
)

// 检查列的属性
type ColOptions struct {
	Table           string      // 表名
	OldColumn       string      // 旧列, CHANGE [COLUMN] old_col_name new_col_name中的old_col_name
	Column          string      // 列名
	Tp              byte        // 列类型
	Flen            int         // 类型长度
	NotNullFlag     bool        // 列是否NOT NULL
	HasDefaultValue bool        // 列是否有默认值
	DefaultValue    interface{} // 列的默认值
	DefaultIsNull   bool        // 列默认值是否为NULL
	HasComment      bool        // 是否有注释
	AuditConfig     *config.AuditConfiguration
}

// 检查列名长度
func (c *ColOptions) CheckColumnLength() error {
	if utf8.RuneCountInString(c.Column) > c.AuditConfig.MAX_COLUMN_NAME_LENGTH {
		return fmt.Errorf("列`%s`命名长度超出限制，最大字符数为%d[表`%s`]", c.Column, c.AuditConfig.MAX_COLUMN_NAME_LENGTH, c.Table)
	}
	return nil
}

// 检查列名合法性
func (c *ColOptions) CheckColumnIdentifer() error {
	if c.AuditConfig.CHECK_IDENTIFIER {
		if ok := utils.IsMatchPattern(utils.NamePattern, c.Column); !ok {
			return fmt.Errorf("列`%s`命名不符合要求，仅允许匹配正则`%s`[表`%s`]", c.Column, utils.NamePattern, c.Table)
		}
	}
	return nil
}

// 检查列名是否为关键字
func (c *ColOptions) CheckColumnIdentiferKeyword() error {
	if c.AuditConfig.CHECK_IDENTIFER_KEYWORD {
		if _, ok := Keywords[strings.ToUpper(c.Column)]; ok {
			return fmt.Errorf("列`%s`命名不允许使用关键字[表`%s`]", c.Column, c.Table)
		}
	}
	return nil
}

// 检查列注释
func (c *ColOptions) CheckColumnComment() error {
	if c.AuditConfig.CHECK_COLUMN_COMMENT && !c.HasComment {
		return fmt.Errorf("列`%s`必须要有注释[表`%s`]", c.Column, c.Table)
	}
	return nil
}

// char建议转换为varchar
func (c *ColOptions) CheckColumnCharToVarchar() error {
	if c.AuditConfig.COLUMN_MAX_CHAR_LENGTH < c.Flen && c.Tp == mysql.TypeString {
		return fmt.Errorf("列`%s`推荐设置为varchar(%d)[表`%s`]", c.Column, c.Flen, c.Table)
	}
	return nil
}

// 最大允许定义的varchar长度
func (c *ColOptions) CheckColumnMaxVarcharLength() error {
	if c.AuditConfig.MAX_VARCHAR_LENGTH < c.Flen && c.Tp == mysql.TypeVarchar {
		return fmt.Errorf("列`%s`最大允许定义的varchar长度为%d，当前varchar长度为%d[表`%s`]", c.Column, c.AuditConfig.MAX_VARCHAR_LENGTH, c.Flen, c.Table)
	}
	return nil
}

// 将float/double转成int/bigint/decimal等
func (c *ColOptions) CheckColumnFloatDouble() error {
	if c.AuditConfig.CHECK_COLUMN_FLOAT_DOUBLE {
		if c.Tp == mysql.TypeFloat || c.Tp == mysql.TypeDouble {
			return fmt.Errorf("列`%s`的类型为float或double，建议转换为int/bigint/decimal类型[表`%s`]", c.Column, c.Table)
		}
	}
	return nil
}

// 列不允许定义的类型
func (c *ColOptions) CheckColumnNotAllowedType() error {
	if !c.AuditConfig.ENABLE_COLUMN_JSON_TYPE && c.Tp == mysql.TypeJSON {
		return fmt.Errorf("列`%s`不允许定义JSON类型[表`%s`]", c.Column, c.Table)
	}
	if !c.AuditConfig.ENABLE_COLUMN_BLOB_TYPE && (c.Tp == mysql.TypeTinyBlob || c.Tp == mysql.TypeMediumBlob || c.Tp == mysql.TypeLongBlob || c.Tp == mysql.TypeBlob) {
		return fmt.Errorf("列`%s`不允许定义BLOB/TEXT类型[表`%s`]", c.Table, c.Column)
	}
	if !c.AuditConfig.ENABLE_COLUMN_TIMESTAMP_TYPE && c.Tp == mysql.TypeTimestamp {
		return fmt.Errorf("列`%s`不允许定义TIMESTAMP类型[表`%s`]", c.Column, c.Table)
	}
	if !c.AuditConfig.ENABLE_COLUMN_BIT_TYPE && c.Tp == mysql.TypeBit {
		return fmt.Errorf("列`%s`不允许定义BIT类型[表`%s`]", c.Column, c.Table)
	}
	return nil
}

// 检查列not null
func (c *ColOptions) CheckColumnNotNull() error {
	if !c.AuditConfig.ENABLE_COLUMN_NOT_NULL {
		return nil
	}
	// 允许为NULL的类型
	allowNULLType := []byte{mysql.TypeBlob, mysql.TypeTinyBlob, mysql.TypeMediumBlob, mysql.TypeLongBlob, mysql.TypeJSON}
	// 是否允许时间类型设置为null
	if c.AuditConfig.ENABLE_COLUMN_TIME_NULL {
		allowNULLType = append(allowNULLType, []byte{mysql.TypeDatetime, mysql.TypeTimestamp, mysql.TypeDate, mysql.TypeYear}...)
	}
	// 列必须定义NOT NULL
	if !utils.IsByteContain(allowNULLType, c.Tp) && !c.NotNullFlag {
		return fmt.Errorf("列`%s`必须定义为`NOT NULL`[表`%s`]", c.Column, c.Table)
	}
	// 不合法的定义`NOT NULL DEFAULT NULL`
	if c.NotNullFlag && c.HasDefaultValue && c.DefaultIsNull {
		return fmt.Errorf("列`%s`不能定义`NOT NULL DEFAULT NULL`[表`%s`]", c.Column, c.Table)
	}
	return nil
}

// 检查列默认值
func (c *ColOptions) CheckColumnDefaultValue() error {
	// BLOB,TEXT,GEOMETRY,JSON类型不能设置默认值
	cannotSetDefaultValueType := []byte{mysql.TypeBlob, mysql.TypeTinyBlob, mysql.TypeMediumBlob, mysql.TypeLongBlob, mysql.TypeJSON, mysql.TypeGeometry}
	if utils.IsByteContain(cannotSetDefaultValueType, c.Tp) {
		if c.HasDefaultValue {
			return fmt.Errorf("列`%s`不能有一个默认值(BLOB/TEXT/GEOMETRY/JSON类型不能有一个默认值)[表`%s`]", c.Column, c.Table)
		}
	}
	// 列需要设置默认值
	if c.AuditConfig.CHECK_COLUMN_DEFAULT_VALUE && !c.HasDefaultValue && !utils.IsByteContain(cannotSetDefaultValueType, c.Tp) {
		return fmt.Errorf("列`%s`需要设置一个默认值[表`%s`]", c.Column, c.Table)
	}
	// 检查默认值(有默认值、且不为NULL)和数据类型是否匹配，Invalid default value
	if c.HasDefaultValue && !c.DefaultIsNull && !utils.IsByteContain(cannotSetDefaultValueType, c.Tp) {
		switch c.Tp {
		case mysql.TypeTiny, mysql.TypeShort, mysql.TypeInt24,
			mysql.TypeLong, mysql.TypeLonglong,
			mysql.TypeYear,
			mysql.TypeFloat, mysql.TypeDouble, mysql.TypeNewDecimal:
			// 验证string型默认值的合法性
			switch val := c.DefaultValue.(type) {
			case string:
				_, intErr := strconv.ParseInt(val, 10, 16)
				_, floatErr := strconv.ParseFloat(val, 64)
				if intErr != nil && floatErr != nil {
					return fmt.Errorf("列`%s`默认值和类型不匹配[表`%s`]", c.Column, c.Table)
				}
			}
		case mysql.TypeVarchar, mysql.TypeString:
			// 判断string型默认值的长度是否超过了定义的长度
			if utf8.RuneCountInString(c.DefaultValue.(string)) > c.Flen {
				return fmt.Errorf("列`%s`的默认值超过了字段类型定义的长度[表`%s`]", c.Column, c.Table)
			}
		}

	}
	// 有默认值，配置了无效的默认值，如default current_timestamp
	if c.HasDefaultValue && !(c.Tp == mysql.TypeTimestamp || c.Tp == mysql.TypeDatetime) && c.DefaultValue == "current_timestamp" {
		return fmt.Errorf("列`%s`配置了无效的默认值(default current_timestamp)[表`%s`]", c.Column, c.Table)
	}
	return nil
}

// CheckColsTypeChanged
func CheckColsTypeChanged(col ColOptions, vCol ColOptions, auditConfig *config.AuditConfiguration, kv *kv.KVCache, tp string, table string) error {
	// tp = modify or change
	var oriColumn string
	if tp == "modify" {
		oriColumn = col.Column
	}
	if tp == "change" {
		oriColumn = col.OldColumn
	}
	// 检查change的列是否进行列类型变更
	intMap := map[byte]string{
		mysql.TypeTiny:     "tinyint",
		mysql.TypeShort:    "smallint",
		mysql.TypeLong:     "int",
		mysql.TypeInt24:    "mediumint",
		mysql.TypeLonglong: "bigint",
	}
	intTp := []byte{mysql.TypeTiny, mysql.TypeShort, mysql.TypeLong, mysql.TypeInt24, mysql.TypeLonglong}

	stringMap := map[byte]string{
		mysql.TypeString:  "char",
		mysql.TypeVarchar: "varchar",
	}
	stringTp := []byte{mysql.TypeString, mysql.TypeVarchar}

	// DB版本
	dbVersionIns := DbVersion{Version: kv.Get("dbVersion").(string)}

	// 匿名函数
	funcCheck := func() error {
		// 开启了兼容变更模式，适用于tidb和mysql
		// 允许同一类型，不同长度变更且变更后的长度必须大于变更前的长度
		// 允许操作: tinyint-> int、int->bigint、char->varchar ...
		// 不允许操作：int -> tinyint、varchar -> char ...
		var tidbTips string
		if dbVersionIns.IsTiDB() {
			tidbTips = "(TiDB不支持当前数据类型变更)"
		}
		if utils.IsByteContain(intTp, col.Tp) && utils.IsByteContain(intTp, vCol.Tp) {
			if bytes.IndexByte(intTp, col.Tp) < bytes.IndexByte(intTp, vCol.Tp) {
				return fmt.Errorf("列`%s`不允许变更数据类型(%s -> %s)[表`%s`]%s", oriColumn, intMap[vCol.Tp], intMap[col.Tp], table, tidbTips)
			}
		} else if utils.IsByteContain(stringTp, col.Tp) && utils.IsByteContain(stringTp, vCol.Tp) {
			if bytes.IndexByte(stringTp, col.Tp) < bytes.IndexByte(stringTp, vCol.Tp) {
				return fmt.Errorf("列`%s`不允许变更数据类型(%s -> %s)[表`%s`]%s", oriColumn, stringMap[vCol.Tp], stringMap[col.Tp], table, tidbTips)
			}
		} else {
			if col.Tp != vCol.Tp {
				return fmt.Errorf("列`%s`不允许变更数据类型[表`%s`]", col.Column, table)
			}
		}
		return nil
	}

	if oriColumn == vCol.Column {
		if auditConfig.ENABLE_COLUMN_TYPE_CHANGE {
			// 允许列类型变更
			return nil
		}
		// 不允许列类型变更
		// 当ENABLE_COLUMN_TYPE_CHANGE = false时，ENABLE_COLUMN_TYPE_CHANGE_COMPATIBLE生效
		if auditConfig.ENABLE_COLUMN_TYPE_CHANGE_COMPATIBLE {
			// 启用兼容模式
			return funcCheck()
		}
		if !auditConfig.ENABLE_COLUMN_TYPE_CHANGE_COMPATIBLE {
			// 禁用兼容模式
			if col.Tp != vCol.Tp {
				return fmt.Errorf("列`%s`不允许变更数据类型[表`%s`]", col.Column, table)
			}
		}
	}
	return nil
}

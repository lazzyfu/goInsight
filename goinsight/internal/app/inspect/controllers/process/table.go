/*
@Time    :   2022/06/24 13:12:20
@Author  :   zongfei.fu
*/

package process

import (
	"fmt"
	"goInsight/internal/pkg/utils"
	"sqlSyntaxAudit/config"
	"strings"
	"unicode/utf8"
)

type TableOptions struct {
	Table         string // 表名
	Type          string // 类型:create or alter
	Engine        string // 引擎
	Charset       string // 表字符集
	Collate       string // 表排序规则
	AutoIncrement uint64 // 表自增
	RowFormat     string // 行格式
	IsPartition   bool   // 是否分区
	PartitionType string // 分区类型
	HasComment    bool   // 表示有注释，代表不了注释为空
	Comment       string // 表注释
	AuditConfig   *config.AuditConfiguration
}

// 检查表名长度
func (t *TableOptions) CheckTableLength() error {
	if utf8.RuneCountInString(t.Table) > t.AuditConfig.MAX_TABLE_NAME_LENGTH {
		return fmt.Errorf("表名`%s`命名长度超出限制，最大字符数%d", t.Table, t.AuditConfig.MAX_TABLE_NAME_LENGTH)
	}
	return nil
}

// 检查表名合法性
func (t *TableOptions) CheckTableIdentifer() error {
	if t.AuditConfig.CHECK_IDENTIFIER {
		if ok := utils.IsMatchPattern(utils.NamePattern, t.Table); !ok {
			return fmt.Errorf("表`%s`的命名不符合要求，仅允许匹配正则`%s`", t.Table, utils.NamePattern)
		}
	}
	return nil
}

// 检查表名是否为关键字
func (t *TableOptions) CheckTableIdentiferKeyword() error {
	if t.AuditConfig.CHECK_IDENTIFER_KEYWORD {
		if _, ok := Keywords[strings.ToUpper(t.Table)]; ok {
			return fmt.Errorf("表`%s`的命名不允许使用关键字", t.Table)
		}
	}
	return nil
}

// 检查存储引擎
func (t *TableOptions) CheckTableEngine() error {
	if t.Type == "create" && len(t.Engine) == 0 {
		return fmt.Errorf("表`%s`必须显式指定存储引擎，支持的存储引擎为%s", t.Table, t.AuditConfig.TABLE_SUPPORT_ENGINE)
	}
	if len(t.Engine) > 0 {
		if !utils.IsContain(t.AuditConfig.TABLE_SUPPORT_ENGINE, t.Engine) {
			return fmt.Errorf("表`%s`指定的存储引擎`%s`不符合要求，支持的存储引擎为`%s`", t.Table, t.Engine, t.AuditConfig.TABLE_SUPPORT_ENGINE)
		}
	}
	return nil
}

// 检查表分区
func (t *TableOptions) CheckTablePartition() error {
	if !t.AuditConfig.ENABLE_PARTITION_TABLE && t.IsPartition {
		return fmt.Errorf("表`%s`不允许定义分区表", t.Table)
	}
	if t.AuditConfig.ENABLE_PARTITION_TABLE && t.IsPartition {
		supportPartTypes := []string{"RANGE", "HASH", "LIST", "KEY", "SYSTEM_TIME"}
		if !utils.IsContain(supportPartTypes, t.PartitionType) {
			return fmt.Errorf("表`%s`不支持自定义的分区类型`%s`，支持的分区类型%s", t.Table, t.PartitionType, supportPartTypes)
		}
	}
	return nil
}

// 检查表注释
func (t *TableOptions) CheckTableComment() error {
	if t.AuditConfig.CHECK_TABLE_COMMENT {
		if t.Type == "create" && !t.HasComment {
			return fmt.Errorf("表`%s`必须要有注释", t.Table)
		}
		if t.HasComment && len(strings.TrimSpace(t.Comment)) == 0 {
			return fmt.Errorf("表`%s`的注释不能为空或空字符", t.Table)
		}
		if t.HasComment && utf8.RuneCountInString(strings.TrimSpace(t.Comment)) > t.AuditConfig.TABLE_COMMENT_LENGTH {
			return fmt.Errorf("表`%s`的注释长度超出限制，最大字符限制为%d", t.Table, t.AuditConfig.TABLE_COMMENT_LENGTH)
		}
	}
	return nil
}

// 检查表字符集
func (t *TableOptions) CheckTableCharset() error {
	// 获取支持的字符集
	var tblSupportCharset []string
	// 获取推荐的排序规则
	var tblRecommendCollation string
	for _, item := range t.AuditConfig.TABLE_SUPPORT_CHARSET {
		tblSupportCharset = append(tblSupportCharset, item["charset"])
		if len(t.Charset) > 0 && item["charset"] == t.Charset {
			tblRecommendCollation = item["recommend"]
		}
	}
	charset := Charset{
		SupportCharset:     tblSupportCharset,
		RecommendCollation: tblRecommendCollation,
		Table:              TableCharset{Table: t.Table, Charset: t.Charset, Collate: t.Collate},
	}
	// 表字符集检查
	if t.AuditConfig.CHECK_TABLE_CHARSET {
		switch t.Type {
		case "create":
			if err := charset.CheckTable(); err != nil {
				return err
			}
		case "alter":
			if len(t.Charset) > 0 || len(t.Collate) > 0 {
				if err := charset.CheckTable(); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// 建表时，自增初始值必须设置为1
func (t *TableOptions) CheckTableAutoIncrementInitValue() error {
	if t.AuditConfig.CHECK_TABLE_AUTOINCREMENT_INIT_VALUE {
		if t.AutoIncrement != 1 && t.Type == "create" {
			// create语句自增值需要设置为1
			return fmt.Errorf("表`%s`的自增初始值必须显式指定且设置为1【例如:ENGINE = InnoDB AUTO_INCREMENT=1】", t.Table)
		}
	}
	return nil
}

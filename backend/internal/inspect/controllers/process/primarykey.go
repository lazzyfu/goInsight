package process

import (
	"fmt"

	"github.com/lazzyfu/goinsight/internal/inspect/config"

	"github.com/pingcap/tidb/pkg/parser/mysql"
)

type PrimaryKey struct {
	Table            string // 表名
	Column           string // 列
	Tp               byte   // 类型
	Flag             uint   // flag
	HasNotNull       bool   // 是否not null
	HasAutoIncrement bool   // 是否自增
	InspectParams    *config.InspectParams
}

func (p *PrimaryKey) CheckBigint() error {
	if p.Tp != mysql.TypeLonglong && p.InspectParams.CHECK_PRIMARYKEY_USE_BIGINT {
		// 必须使用bigint类型
		return fmt.Errorf("主键`%s`必须使用bigint类型（表`%s`）", p.Column, p.Table)
	}
	return nil
}

func (p *PrimaryKey) CheckUnsigned() error {
	if !mysql.HasUnsignedFlag(p.Flag) && p.InspectParams.CHECK_PRIMARYKEY_USE_UNSIGNED {
		// bigint必须定义unsigned
		return fmt.Errorf("主键`%s`必须定义unsigned（表`%s`）", p.Column, p.Table)
	}
	return nil
}

func (p *PrimaryKey) CheckAutoIncrement() error {
	if !p.HasAutoIncrement && p.InspectParams.CHECK_PRIMARYKEY_USE_AUTO_INCREMENT {
		// 必须定义AUTO_INCREMENT
		return fmt.Errorf("主键`%s`必须定义auto_increment（表`%s`）", p.Column, p.Table)
	}
	return nil
}

func (p *PrimaryKey) CheckNotNull() error {
	if !p.HasNotNull {
		// 必须定义NOT NULL
		return fmt.Errorf("主键`%s`必须定义NOT NULL（表`%s`）", p.Column, p.Table)
	}
	return nil
}

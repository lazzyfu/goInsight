/*
@Time    :   2022/07/06 10:12:48
@Author  :   xff
@Desc    :   None
*/

package process

import (
	"errors"
	"goInsight/internal/inspect/controllers/dao"
	"goInsight/pkg/kv"
	"goInsight/pkg/utils"
	"strconv"
	"strings"

	"github.com/mitchellh/mapstructure"
)

// 对应explain的输出
type ExplainOutput struct {
	Table string `json:"Column:table"`
	// MySQL的Explain预估行数
	Rows int `json:"Column:rows"`

	// TiDB (v4.0及之后)的Explain预估行数存储在Count中
	EstRows interface{} `json:"Column:estRows"`
}

type Explain struct {
	DB  *dao.DB
	SQL string
	KV  *kv.KVCache
}

func (e Explain) ConvertToExplain() string {
	var explain []string
	explain = append(explain, "EXPLAIN ")
	explain = append(explain, e.SQL)
	return strings.Join(explain, "")
}

func (e *Explain) Get(EXPLAIN_RULE string) (int, error) {
	explainSQL := e.ConvertToExplain()
	if !strings.HasPrefix(explainSQL, "EXPLAIN") {
		return 0, errors.New("Explain语句未检测到以`EXPLAIN`开头，请联系管理员")
	}
	rows, err := e.DB.Query(explainSQL)
	if err != nil {
		return 0, err
	}
	// 赋值给结构体
	var data []ExplainOutput
	err = mapstructure.WeakDecode(rows, &data)
	if err != nil {
		// 处理多表删除未匹配到行返回rows=NULL的情况
		if strings.Contains(err.Error(), "strconv.ParseInt: parsing \"NULL\"") {
			return 0, nil
		}
		return 0, err
	}
	// 获取db版本
	dbVersionIns := DbVersion{e.KV.Get("dbVersion").(string)}

	var AffectedRows []int
	for _, item := range data {
		if dbVersionIns.IsTiDB() {
			// tidb的执行计划第一行可能是estRows=N/A
			floatEstRows, err := strconv.ParseFloat(item.EstRows.(string), 64)
			if err != nil {
				continue
			}
			// float64 -> int
			intEstRows := int(floatEstRows)
			AffectedRows = append(AffectedRows, intEstRows)
		} else {
			AffectedRows = append(AffectedRows, item.Rows)
		}
	}
	if EXPLAIN_RULE == "first" {
		return AffectedRows[0], nil
	}
	if EXPLAIN_RULE == "max" {
		return utils.MaxInt(AffectedRows), nil
	}
	return 0, nil
}

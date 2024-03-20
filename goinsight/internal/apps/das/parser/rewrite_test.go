/*
@Time    :   2023/06/21 10:45:58
@Author  :   zongfei.fu
@Desc    :   测试用例
*/

package parser

import (
	"fmt"
	"goInsight/global"
	"goInsight/internal/pkg/utils"
	"testing"

	"github.com/pingcap/tidb/pkg/parser/ast"
	"github.com/stretchr/testify/assert"
)

func init() {
	// 初始化配置
	global.App.Config.Das.DefaultReturnRows = 100
	global.App.Config.Das.MaxReturnRows = 100
	global.App.Config.Das.MaxExecutionTime = 600000

}

func parserStmtForRewrite(sql string) (*TiStmt, error) {
	// 检查语法是否有效
	TiStmt, warns, err := NewParse(sql, "", "")
	if len(warns) > 0 {
		return TiStmt, fmt.Errorf("Parse Warning: %s", utils.ErrsJoin("; ", warns))
	}
	if err != nil {
		return TiStmt, fmt.Errorf("SQL语法解析错误:%s", err.Error())
	}
	return TiStmt, nil
}

func TestRewrite(t *testing.T) {
	// 定义用例
	testCases := []struct {
		name      string
		RequestID string
		sql       string
		wantErr   error
		wantRes   string
	}{
		{
			name:      "SELECT",
			sql:       "select * from t1",
			RequestID: "6852dffa-2004-4629-8a4b-9c641b576b93",
			wantRes:   "SELECT /*+ max_execution_time(600000)*/ * FROM `t1` LIMIT 100 /* 6852dffa-2004-4629-8a4b-9c641b576b93 */",
		},
		{
			name:      "SELECT WITH HINT",
			sql:       "SELECT /*+ USE_INDEX(t1, idx1), HASH_AGG() */ count(*) FROM t t1, t t2 WHERE t1.a = t2.b;",
			RequestID: "6852dffa-2004-4629-8a4b-9c641b576b93",
			wantRes:   "SELECT /*+ max_execution_time(600000) USE_INDEX(`t1` `idx1`) HASH_AGG()*/ count(1) FROM (`t` AS `t1`) JOIN `t` AS `t2` WHERE `t1`.`a` = `t2`.`b` LIMIT 100 /* 6852dffa-2004-4629-8a4b-9c641b576b93 */",
		},
		{
			name:      "SELECT WITH OPTIONS",
			sql:       "select /*!40001 SQL_NO_CACHE */ * from t1",
			RequestID: "6852dffa-2004-4629-8a4b-9c641b576b93",
			wantRes:   "SELECT SQL_NO_CACHE /*+ max_execution_time(600000)*/ * FROM `t1` LIMIT 100 /* 6852dffa-2004-4629-8a4b-9c641b576b93 */",
		},
		{
			name:      "SELECT WITH LIMIT 10",
			sql:       "select * from t1 limit 10",
			RequestID: "6852dffa-2004-4629-8a4b-9c641b576b93",
			wantRes:   "SELECT /*+ max_execution_time(600000)*/ * FROM `t1` LIMIT 10 /* 6852dffa-2004-4629-8a4b-9c641b576b93 */",
		},
		{
			name:      "SELECT WITH LIMIT 1000",
			sql:       "select * from t1 limit 1000",
			RequestID: "6852dffa-2004-4629-8a4b-9c641b576b93",
			wantRes:   "SELECT /*+ max_execution_time(600000)*/ * FROM `t1` LIMIT 100 /* 6852dffa-2004-4629-8a4b-9c641b576b93 */",
		},
		{
			name:      "SELECT WITH LIMIT 10 OFFSET 10",
			sql:       "select * from t1 limit 10 offset 10",
			RequestID: "6852dffa-2004-4629-8a4b-9c641b576b93",
			wantRes:   "SELECT /*+ max_execution_time(600000)*/ * FROM `t1` LIMIT 10,10 /* 6852dffa-2004-4629-8a4b-9c641b576b93 */",
		},
		{
			name:      "SELECT WITH LIMIT 10 OFFSET 1000",
			sql:       "select * from t1 limit 10 offset 1000",
			RequestID: "6852dffa-2004-4629-8a4b-9c641b576b93",
			wantRes:   "SELECT /*+ max_execution_time(600000)*/ * FROM `t1` LIMIT 1000,10 /* 6852dffa-2004-4629-8a4b-9c641b576b93 */",
		},
		{
			name:      "SELECT WITH LIMIT 1000 OFFSET 1000",
			sql:       "select * from t1 limit 1000 offset 1000",
			RequestID: "6852dffa-2004-4629-8a4b-9c641b576b93",
			wantRes:   "SELECT /*+ max_execution_time(600000)*/ * FROM `t1` LIMIT 1000,100 /* 6852dffa-2004-4629-8a4b-9c641b576b93 */",
		},
		{
			name:      "SELECT WITH LIMIT 1000,1000",
			sql:       "select * from t1 limit 1000,1000",
			RequestID: "6852dffa-2004-4629-8a4b-9c641b576b93",
			wantRes:   "SELECT /*+ max_execution_time(600000)*/ * FROM `t1` LIMIT 1000,100 /* 6852dffa-2004-4629-8a4b-9c641b576b93 */",
		},
	}
	// 运行
	for _, c := range testCases {
		// 提取stmt
		TiStmt, err := parserStmtForRewrite(c.sql)
		if err != nil {
			t.Error(err)
		}
		t.Run(c.name, func(t *testing.T) {
			// 提取单条SQL的stmt
			var singleStmt ast.StmtNode = TiStmt.Stmts[0]
			rewriter := Rewrite{Stmt: singleStmt, RequestID: c.RequestID}
			restoreSQL := rewriter.Run()
			assert.Equal(t, c.wantRes, restoreSQL)
		})

	}
}

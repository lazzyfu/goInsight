/*
@Time    :   2023/06/21 10:45:58
@Desc    :   测试用例
*/

package parser

import (
	"fmt"
	"goInsight/pkg/utils"
	"testing"

	"github.com/pingcap/tidb/pkg/parser/ast"
	"github.com/stretchr/testify/assert"
)

func parserStmtForExtract(sql string) (*TiStmt, error) {
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

func TestExtract(t *testing.T) {
	// 定义用例
	testCases := []struct {
		name    string
		sql     string
		schema  string
		wantErr error
		wantRes []Table
	}{
		{
			name:    "SELECT 1",
			sql:     "select 1",
			schema:  "information_schema",
			wantRes: []Table{},
		},
		{
			name:   "简单查询 1",
			sql:    "select * from t1",
			schema: "information_schema",
			wantRes: []Table{
				{Schema: "information_schema", Table: "t1"},
			},
		},
		{
			name:   "简单查询 2",
			sql:    "select * from t1",
			schema: "",
			wantRes: []Table{
				{Schema: "", Table: "t1"},
			},
		},
		{
			name:   "简单查询 指定库名 1",
			sql:    "select * from sbtest.t1",
			schema: "",
			wantRes: []Table{
				{Schema: "sbtest", Table: "t1"},
			},
		},
		{
			name:   "简单查询 指定库名 2",
			sql:    "select * from sbtest.t1",
			schema: "information_schema",
			wantRes: []Table{
				{Schema: "sbtest", Table: "t1"},
			},
		},
		{
			name:   "复杂查询 1",
			sql:    "SELECT	d.CH_ID,IFNULL(w.I_WEIGHT, 5) AS I_WEIGHT FROM (SELECT DISTINCT CH_ID FROM b_setting WHERE IS_DROPPED = 0 AND IS_STOP = 0 AND I_QUOTA > 0) AS d LEFT JOIN b_weight AS w ON w.CH_ID = d.CH_ID AND w.I_TYPE = 8010",
			schema: "",
			wantRes: []Table{
				{Schema: "", Table: "b_setting"},
				{Schema: "", Table: "b_weight"},
			},
		},
		{
			name:   "简单UPDATE",
			sql:    "UPDATE b_record SET I_STATUS=2, I_STATUS_DETAIL=0 WHERE I_REF=123456789 and I_STATUS=1",
			schema: "",
			wantRes: []Table{
				{Schema: "", Table: "b_record"},
			},
		},
		{
			name:   "UPDATE JOIN",
			sql:    "UPDATE employees INNER JOIN merits ON employees.performance = merits.performance SET salary = salary + salary * percentage",
			schema: "",
			wantRes: []Table{
				{Schema: "", Table: "employees"},
				{Schema: "", Table: "merits"},
			},
		},
		{
			name:   "CREATE TABLE",
			sql:    "CREATE TABLE IF NOT EXISTS `slamonitor` (`id` bigint(20) not null auto_increment comment '自增主键',primary key (`id`)) ENGINE = InnoDB DEFAULT CHARSET = utf8 COMMENT 'slamonitor'",
			schema: "",
			wantRes: []Table{
				{Schema: "", Table: "slamonitor"},
			},
		},
		{
			name:   "UNION查询",
			sql:    "select id,name from t1 union select id,name from t2 union select id,name from (select * from t3 where id > 1) as xx",
			schema: "sbtest",
			wantRes: []Table{
				{Schema: "sbtest", Table: "t1"},
				{Schema: "sbtest", Table: "t2"},
				{Schema: "sbtest", Table: "t3"},
			},
		},
		{
			name:   "Hint查询",
			sql:    "SELECT /*!40001 SQL_NO_CACHE */ * FROM `film`",
			schema: "",
			wantRes: []Table{
				{Schema: "", Table: "film"},
			},
		},
		{
			name:   "非递归的CTE",
			sql:    "WITH xm_gl AS ( SELECT * FROM products WHERE pname IN ( '小米电视机', '格力空调' ) ) SELECT avg( price ) FROM xm_gl",
			schema: "",
			wantRes: []Table{
				{Schema: "", Table: "products"},
			},
		},
		{
			name:   "非递归的CTE(别名与表名相同)",
			sql:    "WITH products AS ( SELECT * FROM products WHERE pname IN ( select name FROM `order` where user_id=1 ) ) SELECT avg( price ) FROM products",
			schema: "",
			wantRes: []Table{
				{Schema: "", Table: "products"},
				{Schema: "", Table: "order"},
			},
		},
		{
			name:    "非递归的CTE(不包含任何表)",
			sql:     "WITH CTE AS (SELECT 1, 2) SELECT * FROM cte t1, cte t2",
			schema:  "",
			wantRes: []Table{},
		},
		{
			name:    "递归的CTE(不包含任何表)",
			sql:     "WITH RECURSIVE cte(a) AS (SELECT 1 UNION SELECT a+1 FROM cte WHERE a < 5) SELECT * FROM cte",
			schema:  "",
			wantRes: []Table{},
		},
		{
			name:   "递归的CTE",
			sql:    "WITH RECURSIVE cte ( node, path ) AS ( SELECT node, '1'  FROM bst WHERE parent IS NULL UNION ALL SELECT bst.node,  CONCAT ( cte.path, '-->', bst.node ) FROM cte JOIN bst ON cte.node = bst.parent) SELECT * FROM cte ORDER BY node",
			schema: "",
			wantRes: []Table{
				{Schema: "", Table: "bst"},
			},
		},
		{
			name:   "简单JOIN",
			sql:    "select t.table_schema,t.table_name,engine from information_schema.tables t inner join information_schema.columns c on t.table_schema=c.table_schema and t.table_name=c.table_name group by t.table_schema,t.table_name",
			schema: "",
			wantRes: []Table{
				{Schema: "information_schema", Table: "tables"},
				{Schema: "information_schema", Table: "columns"},
			},
		},
		{
			name:   "标量子查询",
			sql:    "select (select max(salary) from b where b.id=a.id) from a",
			schema: "",
			wantRes: []Table{
				{Schema: "", Table: "a"},
				{Schema: "", Table: "b"},
			},
		},
		{
			name:   "简单DELETE",
			sql:    "delete from sbtest.t1 where D_TIME >='2022-08-17 00:00:00' and D_TIME < '2022-08-18 00:00:00'",
			schema: "",
			wantRes: []Table{
				{Schema: "sbtest", Table: "t1"},
			},
		},
		{
			name:   "关联DELETE",
			sql:    "DELETE t1 FROM t1, t2 WHERE t1.id=t2.id",
			schema: "",
			wantRes: []Table{
				{Schema: "", Table: "t1"},
				{Schema: "", Table: "t2"},
			},
		},
		{
			name:   "INSERT INTO SELECT",
			sql:    "INSERT INTO T1 SELECT * FROM T2 WHERE id in (SELECT ID FROM T3)",
			schema: "",
			wantRes: []Table{
				{Schema: "", Table: "t1"},
				{Schema: "", Table: "t2"},
				{Schema: "", Table: "t3"},
			},
		},
		{
			name:   "RENAME TABLE",
			sql:    "RENAME TABLE t1 TO t2, t3 to t4",
			schema: "",
			wantRes: []Table{
				{Schema: "", Table: "t1"},
				{Schema: "", Table: "t2"},
				{Schema: "", Table: "t3"},
				{Schema: "", Table: "t4"},
			},
		},
		{
			name:   "SIMPLE CREATE VIEW",
			sql:    "CREATE VIEW v1 AS SELECT * FROM t1 WHERE c1 > 2",
			schema: "",
			wantRes: []Table{
				{Schema: "", Table: "v1"},
				{Schema: "", Table: "t1"},
			},
		},
		{
			name:   "CREATE VIEW WITH ALGORITHM AND DEFINER",
			sql:    "CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`localhost` SQL SECURITY DEFINER VIEW `v_F_players` AS select `PLAYERS`.`PLAYERNO`,`PLAYERS`.`NAME`,`PLAYERS`.`SEX`,`PLAYERS`.`PHONENO` from `PLAYERS` where (`PLAYERS`.`SEX` = 'F') WITH CASCADED CHECK OPTION",
			schema: "",
			wantRes: []Table{
				{Schema: "", Table: "v_f_players"},
				{Schema: "", Table: "players"},
			},
		},
		{
			name:   "SIMPLE DROP VIEW",
			sql:    "drop view db1.v1",
			schema: "",
			wantRes: []Table{
				{Schema: "db1", Table: "v1"},
			},
		},
		{
			name:   "REPLACE",
			sql:    "REPLACE INTO t1 (id, c1) VALUES(3, 99)",
			schema: "",
			wantRes: []Table{
				{Schema: "", Table: "t1"},
			},
		},
		{
			name:   "CREATE INDEX",
			sql:    "CREATE INDEX idx1 ON t1 ((col1 + col2))",
			schema: "",
			wantRes: []Table{
				{Schema: "", Table: "t1"},
			},
		},
		{
			name:   "DROP INDEX",
			sql:    "DROP INDEX `PRIMARY` ON t",
			schema: "",
			wantRes: []Table{
				{Schema: "", Table: "t"},
			},
		},
		{
			name:   "ALTER TABLE",
			sql:    "ALTER TABLE t1 RENAME COLUMN a TO b, RENAME COLUMN b TO a",
			schema: "",
			wantRes: []Table{
				{Schema: "", Table: "t1"},
			},
		},
		{
			name:   "TRUNCATE TABLE",
			sql:    "truncate table t1",
			schema: "",
			wantRes: []Table{
				{Schema: "", Table: "t1"},
			},
		},
	}
	// 运行
	for _, c := range testCases {
		// 提取stmt
		TiStmt, err := parserStmtForExtract(c.sql)
		if err != nil {
			t.Error(err)
		}
		t.Run(c.name, func(t *testing.T) {
			// 提取单条SQL的stmt
			var singleStmt ast.StmtNode = TiStmt.Stmts[0]
			extracter := Extracter{Stmt: singleStmt, Schema: c.schema}
			extractTables := extracter.Run()
			assert.Equal(t, c.wantRes, extractTables)
		})

	}
}

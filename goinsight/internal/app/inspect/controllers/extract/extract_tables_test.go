package extract

import (
	"errors"
	"sqlSyntaxAudit/config"
	"sqlSyntaxAudit/forms"
	"sqlSyntaxAudit/global"
	logger "sqlSyntaxAudit/middleware/log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	// 初始化配置
	global.App.AuditConfig = &config.AuditConfiguration{
		LogFilePath: "../../logs",
	}
	logger.Setup()
}

func TestChecker_Extract(t *testing.T) {
	tests := []struct {
		name    string
		form    forms.ExtractTablesForm
		wantErr error
		wantRes []ReturnData
	}{
		{
			name: "简单查询",
			form: forms.ExtractTablesForm{
				SqlText:   "select * from t1",
				RequestID: "78c25a06-3b34-4ecb-b9dd-7197078873c7",
			},
			wantRes: []ReturnData{
				{
					Tables: []string{"t1"},
					Type:   "SELECT",
					Query:  "select * from t1",
				},
			},
		},
		{
			name: "UNION查询",
			form: forms.ExtractTablesForm{
				SqlText:   "select id,name from t1 union select id,name from t2 union select id,name from (select * from t3 where id > 1) as xx",
				RequestID: "78c25a06-3b34-4ecb-b9dd-7197078873c7",
			},
			wantRes: []ReturnData{
				{
					Tables: []string{"t1", "t2", "t3"},
					Type:   "SELECT",
					Query:  "select id,name from t1 union select id,name from t2 union select id,name from (select * from t3 where id > 1) as xx",
				},
			},
		},
		{
			// TODO 待决策, 是否支持这种语句
			name: "TIDB不支持的hint(mysql hint)",
			form: forms.ExtractTablesForm{
				SqlText:   "SELECT /*+ NO_RANGE_OPTIMIZATION(t3 PRIMARY, f2_idx) */ f1 FROM t3 WHERE f1 > 30 AND f1 < 33;",
				RequestID: "78c25a06-3b34-4ecb-b9dd-7197078873c7",
			},
			wantErr: errors.New("Parse Warning: [parser:8061]Optimizer hint NO_RANGE_OPTIMIZATION is not supported by TiDB and is ignored"),
		},
		{
			name: "TIDB支持的hint",
			form: forms.ExtractTablesForm{
				SqlText:   "SELECT /*!40001 SQL_NO_CACHE */ * FROM `film`;",
				RequestID: "78c25a06-3b34-4ecb-b9dd-7197078873c7",
			},
			wantRes: []ReturnData{
				{
					// WITH子查询别名不应当放入Tables中
					Tables: []string{
						"film",
					},
					Type:  "SELECT",
					Query: "SELECT /*!40001 SQL_NO_CACHE */ * FROM `film`;",
				},
			},
		},
		{
			name: "非递归的CTE",
			form: forms.ExtractTablesForm{
				SqlText:   "WITH xm_gl AS ( SELECT * FROM products WHERE pname IN ( '小米电视机', '格力空调' ) ) SELECT avg( price ) FROM xm_gl;",
				RequestID: "78c25a06-3b34-4ecb-b9dd-7197078873c7",
			},
			wantRes: []ReturnData{
				{
					// WITH子查询别名不应当放入Tables中
					Tables: []string{
						"products",
					},
					Type:  "SELECT",
					Query: "WITH xm_gl AS ( SELECT * FROM products WHERE pname IN ( '小米电视机', '格力空调' ) ) SELECT avg( price ) FROM xm_gl;",
				},
			},
		},
		{
			name: "非递归的CTE(别名与表名相同)",
			form: forms.ExtractTablesForm{
				SqlText:   "WITH products AS ( SELECT * FROM products WHERE pname IN ( select name FROM `order` where user_id=1 ) ) SELECT avg( price ) FROM products;",
				RequestID: "78c25a06-3b34-4ecb-b9dd-7197078873c7",
			},
			wantRes: []ReturnData{
				{
					// WITH子查询别名不应当放入Tables中
					Tables: []string{
						"products",
						"order",
					},
					Type:  "SELECT",
					Query: "WITH products AS ( SELECT * FROM products WHERE pname IN ( select name FROM `order` where user_id=1 ) ) SELECT avg( price ) FROM products;",
				},
			},
		},
		{
			name: "非递归的CTE(不包含任何表)",
			form: forms.ExtractTablesForm{
				SqlText:   "WITH CTE AS (SELECT 1, 2) SELECT * FROM cte t1, cte t2;",
				RequestID: "78c25a06-3b34-4ecb-b9dd-7197078873c7",
			},
			wantRes: []ReturnData{
				{
					Tables: []string{},
					Type:   "SELECT",
					Query:  "WITH CTE AS (SELECT 1, 2) SELECT * FROM cte t1, cte t2;",
				},
			},
		},
		{
			name: "递归的CTE(不包含任何表)",
			form: forms.ExtractTablesForm{
				SqlText:   "WITH RECURSIVE cte(a) AS (SELECT 1 UNION SELECT a+1 FROM cte WHERE a < 5) SELECT * FROM cte;",
				RequestID: "78c25a06-3b34-4ecb-b9dd-7197078873c7",
			},
			wantRes: []ReturnData{
				{
					Tables: []string{},
					Type:   "SELECT",
					Query:  "WITH RECURSIVE cte(a) AS (SELECT 1 UNION SELECT a+1 FROM cte WHERE a < 5) SELECT * FROM cte;",
				},
			},
		},
		{
			name: "递归的CTE",
			form: forms.ExtractTablesForm{
				SqlText:   "WITH RECURSIVE cte ( node, path ) AS ( SELECT node, '1'  FROM bst WHERE parent IS NULL UNION ALL SELECT bst.node,  CONCAT ( cte.path, '-->', bst.node ) FROM cte JOIN bst ON cte.node = bst.parent) SELECT * FROM cte ORDER BY node;",
				RequestID: "78c25a06-3b34-4ecb-b9dd-7197078873c7",
			},
			wantRes: []ReturnData{
				{
					Tables: []string{
						"bst",
					},
					Type:  "SELECT",
					Query: "WITH RECURSIVE cte ( node, path ) AS ( SELECT node, '1'  FROM bst WHERE parent IS NULL UNION ALL SELECT bst.node,  CONCAT ( cte.path, '-->', bst.node ) FROM cte JOIN bst ON cte.node = bst.parent) SELECT * FROM cte ORDER BY node;",
				},
			},
		},
		{
			name: "SELECT未查询任何表",
			form: forms.ExtractTablesForm{
				SqlText:   "select 'hello';",
				RequestID: "78c25a06-3b34-4ecb-b9dd-7197078873c7",
			},
			wantRes: []ReturnData{
				{
					Tables: []string{},
					Type:   "SELECT",
					Query:  "select 'hello';",
				},
			},
		},
		{
			name: "简单JOIN",
			form: forms.ExtractTablesForm{
				SqlText:   "select t.table_schema,t.table_name,engine from information_schema.tables t inner join information_schema.columns c on t.table_schema=c.table_schema and t.table_name=c.table_name group by t.table_schema,t.table_name;",
				RequestID: "78c25a06-3b34-4ecb-b9dd-7197078873c7",
			},
			wantRes: []ReturnData{
				{
					Tables: []string{
						"tables",
						"columns",
					},
					Type:  "SELECT",
					Query: "select t.table_schema,t.table_name,engine from information_schema.tables t inner join information_schema.columns c on t.table_schema=c.table_schema and t.table_name=c.table_name group by t.table_schema,t.table_name;",
				},
			},
		},
		{
			name: "标量子查询",
			form: forms.ExtractTablesForm{
				SqlText:   "select (select max(salary) from b where b.id=a.id) from a;",
				RequestID: "78c25a06-3b34-4ecb-b9dd-7197078873c7",
			},
			wantRes: []ReturnData{
				{
					Tables: []string{
						"a",
						"b",
					},
					Type:  "SELECT",
					Query: "select (select max(salary) from b where b.id=a.id) from a;",
				},
			},
		},
		{
			name: "简单DELETE",
			form: forms.ExtractTablesForm{
				SqlText:   "delete from t1 where D_TIME >='2022-08-17 00:00:00' and D_TIME < '2022-08-18 00:00:00';",
				RequestID: "78c25a06-3b34-4ecb-b9dd-7197078873c7",
			},
			wantRes: []ReturnData{
				{
					Tables: []string{
						"t1",
					},
					Type:  "DELETE",
					Query: "delete from t1 where D_TIME >='2022-08-17 00:00:00' and D_TIME < '2022-08-18 00:00:00';",
				},
			},
		},
		{
			name: "关联DELETE",
			form: forms.ExtractTablesForm{
				SqlText:   "DELETE t1 FROM t1, t2 WHERE t1.id=t2.id",
				RequestID: "78c25a06-3b34-4ecb-b9dd-7197078873c7",
			},
			wantRes: []ReturnData{
				{
					Tables: []string{
						"t1",
						"t2",
					},
					Type:  "DELETE",
					Query: "DELETE t1 FROM t1, t2 WHERE t1.id=t2.id",
				},
			},
		},
		{
			name: "INSERT INTO SELECT",
			form: forms.ExtractTablesForm{
				SqlText:   "INSERT INTO T1 SELECT * FROM T2 WHERE id in (SELECT ID FROM T3)",
				RequestID: "78c25a06-3b34-4ecb-b9dd-7197078873c7",
			},
			wantRes: []ReturnData{
				{
					Tables: []string{
						"t1",
						"t2",
						"t3",
					},
					Type:  "INSERT",
					Query: "INSERT INTO T1 SELECT * FROM T2 WHERE id in (SELECT ID FROM T3)",
				},
			},
		},
		{
			name: "RENAME TABLE",
			form: forms.ExtractTablesForm{
				SqlText:   "RENAME TABLE t1 TO t2, t3 to t4",
				RequestID: "78c25a06-3b34-4ecb-b9dd-7197078873c7",
			},
			wantRes: []ReturnData{
				{
					Tables: []string{
						"t1",
						"t2",
						"t3",
						"t4",
					},
					Type:  "RENAME TABLE",
					Query: "RENAME TABLE t1 TO t2, t3 to t4",
				},
			},
		},
		{
			name: "SIMPLE CREATE VIEW",
			form: forms.ExtractTablesForm{
				SqlText:   "CREATE VIEW v1 AS SELECT * FROM t1 WHERE c1 > 2;",
				RequestID: "78c25a06-3b34-4ecb-b9dd-7197078873c7",
			},
			wantRes: []ReturnData{
				{
					Tables: []string{
						"v1",
						"t1",
					},
					Type:  "CREATE VIEW",
					Query: "CREATE VIEW v1 AS SELECT * FROM t1 WHERE c1 > 2;",
				},
			},
		},
		{
			name: "CREATE VIEW WITH ALGORITHM AND DEFINER",
			form: forms.ExtractTablesForm{
				SqlText:   "CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`localhost` SQL SECURITY DEFINER VIEW `v_F_players` AS select `PLAYERS`.`PLAYERNO`,`PLAYERS`.`NAME`,`PLAYERS`.`SEX`,`PLAYERS`.`PHONENO` from `PLAYERS` where (`PLAYERS`.`SEX` = 'F') WITH CASCADED CHECK OPTION;",
				RequestID: "78c25a06-3b34-4ecb-b9dd-7197078873c7",
			},
			wantRes: []ReturnData{
				{
					Tables: []string{
						"v_f_players",
						"players",
					},
					Type:  "CREATE VIEW",
					Query: "CREATE ALGORITHM=UNDEFINED DEFINER=`root`@`localhost` SQL SECURITY DEFINER VIEW `v_F_players` AS select `PLAYERS`.`PLAYERNO`,`PLAYERS`.`NAME`,`PLAYERS`.`SEX`,`PLAYERS`.`PHONENO` from `PLAYERS` where (`PLAYERS`.`SEX` = 'F') WITH CASCADED CHECK OPTION;",
				},
			},
		},
		{
			name: "SIMPLE DROP VIEW",
			form: forms.ExtractTablesForm{
				SqlText:   "drop view db1.v1;",
				RequestID: "78c25a06-3b34-4ecb-b9dd-7197078873c7",
			},
			wantRes: []ReturnData{
				{
					Tables: []string{
						"v1",
					},
					// 因为create table和create view已经区分出来了, 所以drop也应该区分
					Type:  "DROP VIEW",
					Query: "drop view db1.v1;",
				},
			},
		},
		{
			name: "REPLACE",
			form: forms.ExtractTablesForm{
				SqlText:   "REPLACE INTO t1 (id, c1) VALUES(3, 99);",
				RequestID: "78c25a06-3b34-4ecb-b9dd-7197078873c7",
			},
			wantRes: []ReturnData{
				{
					Tables: []string{
						"t1",
					},
					Type:  "REPLACE",
					Query: "REPLACE INTO t1 (id, c1) VALUES(3, 99);",
				},
			},
		},
		{
			name: "CREATE INDEX",
			form: forms.ExtractTablesForm{
				SqlText:   "CREATE INDEX idx1 ON t1 ((col1 + col2));",
				RequestID: "78c25a06-3b34-4ecb-b9dd-7197078873c7",
			},
			wantRes: []ReturnData{
				{
					Tables: []string{
						"t1",
					},
					Type:  "CREATE INDEX",
					Query: "CREATE INDEX idx1 ON t1 ((col1 + col2));",
				},
			},
		},
		{
			name: "DROP INDEX",
			form: forms.ExtractTablesForm{
				SqlText:   "DROP INDEX `PRIMARY` ON t;",
				RequestID: "78c25a06-3b34-4ecb-b9dd-7197078873c7",
			},
			wantRes: []ReturnData{
				{
					Tables: []string{
						"t",
					},
					Type:  "DROP INDEX",
					Query: "DROP INDEX `PRIMARY` ON t;",
				},
			},
		},
		{
			name: "ALTER TABLE",
			form: forms.ExtractTablesForm{
				SqlText:   "ALTER TABLE t1 RENAME COLUMN a TO b, RENAME COLUMN b TO a;",
				RequestID: "78c25a06-3b34-4ecb-b9dd-7197078873c7",
			},
			wantRes: []ReturnData{
				{
					Tables: []string{
						"t1",
					},
					Type:  "ALTER TABLE",
					Query: "ALTER TABLE t1 RENAME COLUMN a TO b, RENAME COLUMN b TO a;",
				},
			},
		},
		{
			name: "TRUNCATE TABLE",
			form: forms.ExtractTablesForm{
				SqlText:   "truncate table t1",
				RequestID: "78c25a06-3b34-4ecb-b9dd-7197078873c7",
			},
			wantRes: []ReturnData{
				{
					Tables: []string{
						"t1",
					},
					Type:  "TRUNCATE TABLE",
					Query: "truncate table t1",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Checker{
				Form: tt.form,
			}
			err, res := c.Extract(tt.form.RequestID)
			assert.Equal(t, tt.wantErr, err)
			if tt.wantErr != nil {
				// 预期会有错误返回，就不需要进一步校验res了
				return
			}
			assert.Equal(t, tt.wantRes, res)
		})
	}
}

func Test_removeElement(t *testing.T) {
	type args struct {
		org         []string
		toBeRemoved []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "simple",
			args: args{
				org:         []string{"t1", "t2"},
				toBeRemoved: []string{"t2"},
			},
			want: []string{"t1"},
		},
		{
			name: "toBeRemoved比org多",
			args: args{
				org:         []string{"t1", "t2"},
				toBeRemoved: []string{"t2", "t8", "t9"},
			},
			want: []string{"t1"},
		},
		{
			name: "org is empty",
			args: args{
				org:         []string{},
				toBeRemoved: []string{"t2", "t8", "t9"},
			},
			want: []string{},
		},
		{
			name: "toBeRemoved is empty",
			args: args{
				org:         []string{"t1", "t2"},
				toBeRemoved: []string{},
			},
			want: []string{"t1", "t2"},
		},
		{
			name: "org is nil",
			args: args{
				org:         nil,
				toBeRemoved: []string{"t2", "t8", "t9"},
			},
			want: nil,
		},
		{
			name: "toBeRemoved is nil",
			args: args{
				org:         []string{"t1", "t2"},
				toBeRemoved: nil,
			},
			want: []string{"t1", "t2"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, removeElement(tt.args.org, tt.args.toBeRemoved), "removeElement(%v, %v)", tt.args.org, tt.args.toBeRemoved)
		})
	}
}

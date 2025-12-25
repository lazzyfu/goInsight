package parser

import (
	"fmt"
	"testing"

	"github.com/lazzyfu/goinsight/pkg/utils"

	"github.com/pingcap/tidb/pkg/parser/ast"
	"github.com/stretchr/testify/assert"
)

func parserStmtForStatementType(sql string) (*TiStmt, error) {
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

func TestStatementType(t *testing.T) {
	// 定义用例
	testCases := []struct {
		name    string
		sql     string
		wantErr error
		wantRes string
	}{
		{
			name:    "SELECT",
			sql:     "select * from t1",
			wantRes: "SELECT",
		},
		{
			name:    "SELECT WITH CTE",
			sql:     "WITH xm_gl AS (SELECT * FROM products WHERE pname IN ('小米电视机', '格力空调')) SELECT avg(price) FROM xm_gl",
			wantRes: "SELECT",
		},
		{
			name:    "DELETE",
			sql:     "delete from t1 where id = 10",
			wantRes: "DELETE",
		},
		{
			name:    "INSERT",
			sql:     "insert into t1(id, name) values(1, 'zs')",
			wantRes: "INSERT",
		},
		{
			name:    "UPDATE",
			sql:     "UPDATE employees INNER JOIN merits ON employees.performance = merits.performance SET salary = salary + salary * percentage",
			wantRes: "UPDATE",
		},
		{
			name:    "ShowEngines",
			sql:     "show engines",
			wantRes: "ShowEngines",
		},
		{
			name:    "ShowDatabases",
			sql:     "show databases",
			wantRes: "ShowDatabases",
		},
		{
			name:    "ShowTables",
			sql:     "show tables",
			wantRes: "ShowTables",
		},
		{
			name:    "ShowTableStatus",
			sql:     "SHOW table status like 't1'",
			wantRes: "ShowTableStatus",
		},
		{
			name:    "ShowColumns",
			sql:     "show columns from v1",
			wantRes: "ShowColumns",
		},
		{
			name:    "ShowWarnings",
			sql:     "show warnings",
			wantRes: "ShowWarnings",
		},
		{
			name:    "ShowCharset",
			sql:     "SHOW CHARACTER SET",
			wantRes: "ShowCharset",
		},
		{
			name:    "ShowCollation",
			sql:     "SHOW COLLATION",
			wantRes: "ShowCollation",
		},
		{
			name:    "ShowVariables",
			sql:     "show global variables",
			wantRes: "ShowVariables",
		},
		{
			name:    "ShowStatus",
			sql:     "show global status",
			wantRes: "ShowStatus",
		},
		{
			name:    "ShowCreateTable",
			sql:     "SHOW CREATE TABLE t1",
			wantRes: "ShowCreateTable",
		},
		{
			name:    "ShowCreateView",
			sql:     "SHOW CREATE view v_t1",
			wantRes: "ShowCreateView",
		},
		{
			name:    "ShowCreateUser",
			sql:     "SHOW CREATE USER 'root'",
			wantRes: "ShowCreateUser",
		},
		{
			name:    "ShowCreateSequence",
			sql:     "SHOW CREATE SEQUENCE seq",
			wantRes: "ShowCreateSequence",
		},
		{
			name:    "ShowCreatePlacementPolicy",
			sql:     "SHOW CREATE PLACEMENT POLICY p1",
			wantRes: "ShowCreatePlacementPolicy",
		},
		{
			name:    "ShowGrants",
			sql:     "show grants",
			wantRes: "ShowGrants",
		},
		{
			name:    "ShowTriggers",
			sql:     "show triggers",
			wantRes: "ShowTriggers",
		},
		{
			name:    "ShowProcedureStatus",
			sql:     "show procedure status",
			wantRes: "ShowProcedureStatus",
		},
		{
			name:    "ShowIndex",
			sql:     "SHOW INDEXES FROM t1",
			wantRes: "ShowIndex",
		},
		{
			name:    "ShowProcessList",
			sql:     "show processlist",
			wantRes: "ShowProcessList",
		},
		{
			name:    "ShowCreateDatabase",
			sql:     "show create database test",
			wantRes: "ShowCreateDatabase",
		},
		{
			name:    "ShowConfig",
			sql:     "SHOW CONFIG WHERE type = 'tidb' AND name = 'advertise-address'",
			wantRes: "ShowConfig",
		},
		{
			name:    "ShowEvents",
			sql:     "show events",
			wantRes: "ShowEvents",
		},
		{
			name:    "ShowStatsExtended",
			sql:     "show STATS_EXTENDED",
			wantRes: "ShowStatsExtended",
		},
		{
			name:    "ShowStatsMeta",
			sql:     "show stats_meta",
			wantRes: "ShowStatsMeta",
		},
		{
			name:    "ShowStatsHistograms",
			sql:     "show stats_histograms",
			wantRes: "ShowStatsHistograms",
		},
		{
			name:    "ShowStatsTopN",
			sql:     "show STATS_TOPN",
			wantRes: "ShowStatsTopN",
		},

		{
			name:    "ShowStatsBuckets",
			sql:     "show STATS_BUCKETS",
			wantRes: "ShowStatsBuckets",
		},
		{
			name:    "ShowStatsHealthy",
			sql:     "show STATS_HEALTHY",
			wantRes: "ShowStatsHealthy",
		},
		{
			name:    "ShowStatsLocked",
			sql:     "show STATS_LOCKED",
			wantRes: "ShowStatsLocked",
		},
		{
			name:    "ShowHistogramsInFlight",
			sql:     "show HISTOGRAMS_IN_FLIGHT",
			wantRes: "ShowHistogramsInFlight",
		},
		{
			name:    "ShowColumnStatsUsage",
			sql:     "show COLUMN_STATS_USAGE",
			wantRes: "ShowColumnStatsUsage",
		},
		{
			name:    "ShowPlugins",
			sql:     "show PLUGINS",
			wantRes: "ShowPlugins",
		},
		{
			name:    "ShowProfile",
			sql:     "SHOW PROFILE",
			wantRes: "ShowProfile",
		},
		{
			name:    "ShowProfiles",
			sql:     "SHOW PROFILES",
			wantRes: "ShowProfiles",
		},
		{
			name:    "ShowMasterStatus",
			sql:     "SHOW master status",
			wantRes: "ShowMasterStatus",
		},
		{
			name:    "ShowPrivileges",
			sql:     "SHOW privileges",
			wantRes: "ShowPrivileges",
		},
		{
			name:    "ShowErrors",
			sql:     "SHOW errors",
			wantRes: "ShowErrors",
		},
		{
			name:    "ShowBindings",
			sql:     "SHOW bindings",
			wantRes: "ShowBindings",
		},
		{
			name:    "ShowBindingCacheStatus",
			sql:     "SHOW BINDING_CACHE STATUS",
			wantRes: "ShowBindingCacheStatus",
		},
		{
			name:    "ShowPumpStatus",
			sql:     "SHOW PUMP STATUS",
			wantRes: "ShowPumpStatus",
		},
		{
			name:    "ShowDrainerStatus",
			sql:     "SHOW DRAINER STATUS",
			wantRes: "ShowDrainerStatus",
		},
		{
			name:    "ShowOpenTables",
			sql:     "SHOW OPEN TABLES",
			wantRes: "ShowOpenTables",
		},
		{
			name:    "ShowAnalyzeStatus",
			sql:     "SHOW ANALYZE STATUS",
			wantRes: "ShowAnalyzeStatus",
		},
		{
			name:    "ShowRegions",
			sql:     "SHOW TABLE t1 REGIONS",
			wantRes: "ShowRegions",
		},
		{
			name:    "ShowBuiltins",
			sql:     "SHOW BUILTINS",
			wantRes: "ShowBuiltins",
		},
		{
			name:    "ShowTableNextRowId",
			sql:     "show table t next_row_id",
			wantRes: "ShowTableNextRowId",
		},
		{
			name:    "ShowBackups",
			sql:     "show BACKUPS",
			wantRes: "ShowBackups",
		},
		{
			name:    "ShowRestores",
			sql:     "show RESTORES",
			wantRes: "ShowRestores",
		},
		{
			name:    "ShowPlacement",
			sql:     "show PLACEMENT",
			wantRes: "ShowPlacement",
		},
		{
			name:    "ShowPlacementForDatabase",
			sql:     "show PLACEMENT FOR DATABASE sbtest",
			wantRes: "ShowPlacementForDatabase",
		},
		{
			name:    "ShowPlacementForTable",
			sql:     "show PLACEMENT FOR table sbtest1",
			wantRes: "ShowPlacementForTable",
		},
		{
			name:    "ShowPlacementForPartition",
			sql:     "show PLACEMENT FOR table t1 partition p1",
			wantRes: "ShowPlacementForPartition",
		},
		{
			name:    "ShowPlacementLabels",
			sql:     "SHOW PLACEMENT LABELS",
			wantRes: "ShowPlacementLabels",
		},
		{
			name:    "ShowSessionStates",
			sql:     "SHOW SESSION_STATES",
			wantRes: "ShowSessionStates",
		},
		{
			name:    "ExplainSelect",
			sql:     "explain select * from t1 where id = 100",
			wantRes: "ExplainSelect",
		},
		{
			name:    "ExplainDelete",
			sql:     "explain delete from t1",
			wantRes: "ExplainDelete",
		},
		{
			name:    "ExplainInsert",
			sql:     "explain insert into t1(id, name) values(1, 'zs')",
			wantRes: "ExplainInsert",
		},
		{
			name:    "ExplainUpdate",
			sql:     "explain update t1 set name='hello' where id = 100",
			wantRes: "ExplainUpdate",
		},
		{
			name:    "ExplainFor",
			sql:     "explain for connection 100",
			wantRes: "ExplainFor",
		},
		{
			name:    "Desc",
			sql:     "desc t1",
			wantRes: "Desc",
		},
		{
			name:    "CALL",
			sql:     "call showDbSize('sbtest1')",
			wantRes: "CALL",
		},
		{
			name:    "SET",
			sql:     "SET CHARACTER SET utf8mb4",
			wantRes: "SET",
		},
		{
			name:    "ALTER",
			sql:     "alter table t1 add col1 char(30) not null default null comment 'xxx'",
			wantRes: "ALTER",
		},
		{
			name:    "CREATE",
			sql:     "create database sbtest",
			wantRes: "CREATE",
		},
		{
			name:    "DROP",
			sql:     "drop database sbtest",
			wantRes: "DROP",
		},
		{
			name:    "RENAME",
			sql:     "RENAME table t1 to t2",
			wantRes: "RENAME",
		},
		{
			name:    "TRUNCATE",
			sql:     "TRUNCATE table t1",
			wantRes: "TRUNCATE",
		},
		{
			name:    "REPAIR",
			sql:     "ADMIN REPAIR TABLE tbl_name CREATE TABLE STATEMENT",
			wantRes: "REPAIR",
		},
		{
			name:    "LoadData",
			sql:     "LOAD DATA LOCAL INFILE '/mnt/evo970/data-sets/bikeshare-data/2017Q4-capitalbikeshare-tripdata.csv' INTO TABLE trips FIELDS TERMINATED BY x'2c' ENCLOSED BY b'100010' LINES TERMINATED BY '\r\n' IGNORE 1 LINES (duration, start_date, end_date, start_station_number, start_station, end_station_number, end_station, bike_number, member_type)",
			wantRes: "LoadData",
		},
		{
			name:    "SplitRegion",
			sql:     "SPLIT TABLE t BETWEEN (-9223372036854775808) AND (9223372036854775807) REGIONS 16",
			wantRes: "SplitRegion",
		},
		{
			name:    "Use",
			sql:     "USE mysql",
			wantRes: "Use",
		},
		{
			name:    "ShutDown",
			sql:     "SHUTDOWN",
			wantRes: "ShutDown",
		},
	}
	// 运行
	for _, c := range testCases {
		// 提取stmt
		TiStmt, err := parserStmtForStatementType(c.sql)
		if err != nil {
			t.Error(err)
		}
		t.Run(c.name, func(t *testing.T) {
			// 提取单条SQL的stmt
			var singleStmt ast.StmtNode = TiStmt.Stmts[0]
			st := StatementType{}
			statementType := st.Extract(singleStmt)
			assert.Equal(t, c.wantRes, statementType)
		})

	}
}

/*
@Time    :   2022/07/06 10:12:48
@Author  :   zongfei.fu
@Desc    :   None
*/

package config

type DisableTablesAudit struct {
	DB     string   // 库名
	Tables []string // 表名
	Reason string   // 原因
}

type InspectParams struct {
	// TABLE
	MAX_TABLE_NAME_LENGTH                int                 // 表名的长度
	CHECK_TABLE_COMMENT                  bool                // 检查表是否有注释
	TABLE_COMMENT_LENGTH                 int                 // 表的注释长度
	CHECK_IDENTIFIER                     bool                // 对象名必须使用字符串范围为正则[a-zA-Z0-9_]
	CHECK_IDENTIFER_KEYWORD              bool                // 对象名是否可以使用关键字
	CHECK_TABLE_CHARSET                  bool                // 是否检查表的字符集和排序规则
	TABLE_SUPPORT_CHARSET                []map[string]string // 表支持的字符集
	CHECK_TABLE_ENGINE                   bool                // 是否检查表的存储引擎
	TABLE_SUPPORT_ENGINE                 []string            // 表支持的存储引擎
	ENABLE_PARTITION_TABLE               bool                // 是否启用分区表
	CHECK_TABLE_PRIMARY_KEY              bool                // 检查表是否有主键
	TABLE_AT_LEAST_ONE_COLUMN            bool                // 表至少要有一列，语法默认支持
	CHECK_TABLE_AUDIT_TYPE_COLUMNS       bool                // 启用审计类型的字段(col1 datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP && col2 datetime DEFAULT CURRENT_TIMESTAMP)
	ENABLE_CREATE_TABLE_AS               bool                // 是否允许create table as语法
	ENABLE_CREATE_TABLE_LIKE             bool                // 是否允许create table like语法
	ENABLE_FOREIGN_KEY                   bool                // 是否启用外键
	CHECK_TABLE_AUTOINCREMENT_INIT_VALUE bool                // 检查建表是自增列初始值是否为1
	ENABLE_CREATE_VIEW                   bool                // 是否支持创建和使用视图
	// COLUMN
	MAX_COLUMN_NAME_LENGTH               int  // 列名的长度
	CHECK_COLUMN_CHARSET                 bool // 是否检查列的字符集
	CHECK_COLUMN_COMMENT                 bool // 是否检查列的注释
	COLUMN_MAX_CHAR_LENGTH               int  // char长度大于N的时候需要改为varchar
	MAX_VARCHAR_LENGTH                   int  // 最大允许定义的varchar长度
	ENABLE_COLUMN_BLOB_TYPE              bool // 是否允许列的类型为BLOB/TEXT
	ENABLE_COLUMN_JSON_TYPE              bool // 是否允许列的类型为JSON
	ENABLE_COLUMN_BIT_TYPE               bool // 是否允许列的类型为BIT
	ENABLE_COLUMN_TIMESTAMP_TYPE         bool // 是否允许列的类型为TIMESTAMP
	CHECK_PRIMARYKEY_USE_BIGINT          bool // 主键是否为bigint
	CHECK_PRIMARYKEY_USE_UNSIGNED        bool // 主键bigint是否为unsigned
	CHECK_PRIMARYKEY_USE_AUTO_INCREMENT  bool // 主键是否定义为自增
	ENABLE_COLUMN_NOT_NULL               bool // 是否允许列定义为NOT NULL
	ENABLE_COLUMN_TIME_NULL              bool // 是否允许时间类型设置为NULL
	CHECK_COLUMN_DEFAULT_VALUE           bool // 列必须要有默认值
	CHECK_COLUMN_FLOAT_DOUBLE            bool // 将float/double转成int/bigint/decimal等
	ENABLE_COLUMN_TYPE_CHANGE            bool // 是否允许变更列类型
	ENABLE_COLUMN_TYPE_CHANGE_COMPATIBLE bool // 允许tinyint-> int、int->bigint、char->varchar等
	ENABLE_COLUMN_CHANGE_COLUMN_NAME     bool // 是否允许CHANGE修改列名操作
	// INDEX
	CHECK_UNIQ_INDEX_PREFIX       bool   // 是否检查唯一索引前缀,如唯一索引必须以uniq_为前缀
	CHECK_SECONDARY_INDEX_PREFIX  bool   // 是否检查二级索引前缀,如普通索引必须以idx_为前缀
	CHECK_FULLTEXT_INDEX_PREFIX   bool   // 是否检查全文索引前缀,如全文索引必须以full_为前缀
	UNQI_INDEX_PREFIX             string // 定义唯一索引前缀，不区分大小写
	SECONDARY_INDEX_PREFIX        string // 定义二级索引前缀，不区分大小写
	FULLTEXT_INDEX_PREFIX         string // 定义全文索引前缀，不区分大小写
	SECONDARY_INDEX_MAX_KEY_PARTS int    // 组成二级索引的列数不能超过指定的个数,包括唯一索引
	PRIMARYKEY_MAX_KEY_PARTS      int    // 组成主键索引的列数不能超过指定的个数
	MAX_INDEX_KEYS                int    // 最多有N个索引,包括唯一索引/二级索引
	ENABLE_INDEX_RENAME           bool   // 是否允许rename索引名
	// ALTER
	ENABLE_DROP_COLS               bool // 是否允许DROP列
	ENABLE_DROP_INDEXES            bool // 是否允许DROP索引
	ENABLE_DROP_PRIMARYKEY         bool // 是否允许DROP主键
	ENABLE_DROP_TABLE              bool // 是否允许DROP TABLE
	ENABLE_TRUNCATE_TABLE          bool // 是否允许TRUNCATE TABLE
	ENABLE_RENAME_TABLE_NAME       bool // 是否允许rename表名
	ENABLE_MYSQL_MERGE_ALTER_TABLE bool // MySQL同一个表的多个ALTER是否合并为单条语句
	ENABLE_TIDB_MERGE_ALTER_TABLE  bool // TiDB同一个表的多个ALTER是否合并为单条语句
	// DML
	DML_MUST_HAVE_WHERE        bool   // DML语句必须有where条件
	DML_DISABLE_LIMIT          bool   // DML语句中不允许有LIMIT
	DML_DISABLE_ORDERBY        bool   // DML语句中不允许有orderby
	DML_DISABLE_SUBQUERY       bool   // DML语句不能有子查询
	CHECK_DML_JOIN_WITH_ON     bool   // DML的JOIN语句必须有ON语句
	EXPLAIN_RULE               string // explain判断受影响行数时使用的规则("first", "max")。 "first": 使用第一行的explain结果作为受影响行数, "max": 使用explain结果中的最大值作为受影响行数
	MAX_AFFECTED_ROWS          int    // 最大影响行数，默认100
	MAX_INSERT_ROWS            int    // 一次最多允许insert的行, eg: insert into tbl(col,...) values(row1), (row2)...
	DISABLE_REPLACE            bool   // 是否禁用replace语句
	DISABLE_INSERT_INTO_SELECT bool   // 是否禁用insert/replace into select语法
	DISABLE_ON_DUPLICATE       bool   // 是否禁止insert on duplicate语法
	// 禁止语法审核的表
	DISABLE_AUDIT_DML_TABLES []DisableTablesAudit // 禁止指定的表的DML语句进行审核
	DISABLE_AUDIT_DDL_TABLES []DisableTablesAudit // 禁止指定的表的DDL语句进行审核
}

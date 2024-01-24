/*
@Time    :   2023/08/14 15:53:53
@Author  :   lazzyfu
*/

package bootstrap

import (
	"context"
	"encoding/json"
	"fmt"
	"goInsight/global"
	commonModels "goInsight/internal/app/common/models"
	dasModels "goInsight/internal/app/das/models"
	inspectModels "goInsight/internal/app/inspect/models"
	ordersModels "goInsight/internal/app/orders/models"
	usersModels "goInsight/internal/app/users/models"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

// 初始化redis
func InitializeRedis() *redis.Client {
	var ctx = context.Background()
	config := global.App.Config.Redis
	if config.Host == "" {
		return nil
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Host, config.Port),
		Password: config.Password,
		DB:       config.DB,
		PoolSize: 512,
	})
	err := rdb.Set(ctx, "testkey", "value", 0).Err()
	if err != nil {
		global.App.Log.Error(err)
		panic(err)
	}
	return rdb
}

func InitializeDB() *gorm.DB {
	switch global.App.Config.Database.Driver {
	case "mysql":
		return initializeMySQLGorm()
	default:
		return initializeMySQLGorm()
	}
}

// 初始化MySQL
func initializeMySQLGorm() *gorm.DB {
	config := global.App.Config.Database
	if config.Database == "" {
		return nil
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local", config.UserName, config.Password, config.Host, config.Port, config.Database, config.Charset)
	if db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{}); err != nil {
		global.App.Log.Error("mysql connect failed, err:", err.Error())
		panic(err.Error())
	} else {
		if err := db.Use(
			dbresolver.Register(dbresolver.Config{ /* xxx */ }).
				SetConnMaxIdleTime(time.Duration(config.ConnMaxIdleTime) * time.Second).
				SetConnMaxLifetime(time.Duration(config.ConnMaxLifetime) * time.Second).
				SetMaxIdleConns(config.MaxIdleConns).
				SetMaxOpenConns(config.MaxOpenConns),
		); err != nil {
			global.App.Log.Error("mysql set conn params failed, err:", err.Error())
			panic(err.Error())
		}
		// 初始化表
		initializeTables(db)
		// 初始化用户允许的操作
		initializeAllowedOperations(db)
		// 初始化审核参数
		initializeInspectParams(db)
		// 初始化系统管理员
		initializeAdminUser(db)
		return db
	}
}

// 初始化系统表
func initializeTables(db *gorm.DB) {
	err := db.Set("gorm:table_options", "ENGINE=InnoDB CHARSET=utf8mb4").AutoMigrate(
		// 用户
		&usersModels.InsightUsers{},
		&usersModels.InsightRoles{},
		&usersModels.InsightOrganizations{},
		&usersModels.InsightOrganizationsUsers{},
		// common
		&commonModels.InsightDBEnvironments{},
		&commonModels.InsightDBConfig{},
		&commonModels.InsightDBSchemas{},
		// inspect
		&inspectModels.InsightInspectParams{},
		// das
		&dasModels.InsightDASUserSchemaPermissions{},
		&dasModels.InsightDASUserTablePermissions{},
		&dasModels.InsightDASRecords{},
		&dasModels.InsightDASAllowedOperations{},
		&dasModels.InsightDASFavorites{},
		// orders
		&ordersModels.InsightOrderRecords{},
		&ordersModels.InsightOrderTasks{},
		&ordersModels.InsightOrderOpLogs{},
		&ordersModels.InsightOrderMessages{},
	)
	if err != nil {
		global.App.Log.Fatal("migrate table failed", err.Error())
	}
}

// 初始化系统管理员
func initializeAdminUser(db *gorm.DB) {
	var user usersModels.InsightUsers
	_ = db.FirstOrCreate(&user, usersModels.InsightUsers{
		Username:    "admin",
		Password:    "$2a$10$36U.TwQGCRC8jlse3SgQY.sJweMofGXLtEejF.xBlzQN0iqauHIN.",
		Email:       "admin@example.com",
		NickName:    "管理员",
		Mobile:      "",
		AvatarFile:  "/static/avatar2.jpg",
		IsSuperuser: true,
		IsActive:    true,
		IsStaff:     false,
		IsTwoFA:     false,
	})
}

// 初始化用户允许的操作
func initializeAllowedOperations(db *gorm.DB) {
	var ops []map[string]interface{} = []map[string]interface{}{
		{"name": "SELECT", "is_enable": true, "remark": ""},
		{"name": "UNION", "is_enable": true, "remark": ""},
		{"name": "Use", "is_enable": true, "remark": ""},
		{"name": "Desc", "is_enable": true, "remark": ""},
		{"name": "ExplainSelect", "is_enable": true, "remark": ""},
		{"name": "ExplainDelete", "is_enable": true, "remark": ""},
		{"name": "ExplainInsert", "is_enable": true, "remark": ""},
		{"name": "ExplainUpdate", "is_enable": true, "remark": ""},
		{"name": "ExplainUnion", "is_enable": true, "remark": ""},
		{"name": "ExplainFor", "is_enable": true, "remark": "ExplainForStmt is a statement to provite information about how is SQL statement executeing in connection #ConnectionID"},
		{"name": "ShowEngines", "is_enable": false, "remark": ""},
		{"name": "ShowDatabases", "is_enable": false, "remark": ""},
		{"name": "ShowTables", "is_enable": false, "remark": ""},
		{"name": "ShowTableStatus", "is_enable": false, "remark": ""},
		{"name": "ShowColumns", "is_enable": false, "remark": ""},
		{"name": "ShowWarnings", "is_enable": false, "remark": ""},
		{"name": "ShowCharset", "is_enable": false, "remark": ""},
		{"name": "ShowVariables", "is_enable": false, "remark": ""},
		{"name": "ShowStatus", "is_enable": false, "remark": ""},
		{"name": "ShowCollation", "is_enable": false, "remark": ""},
		{"name": "ShowCreateTable", "is_enable": false, "remark": ""},
		{"name": "ShowCreateView", "is_enable": false, "remark": ""},
		{"name": "ShowCreateUser", "is_enable": false, "remark": ""},
		{"name": "ShowCreateSequence", "is_enable": false, "remark": ""},
		{"name": "ShowCreatePlacementPolicy", "is_enable": false, "remark": ""},
		{"name": "ShowGrants", "is_enable": false, "remark": ""},
		{"name": "ShowTriggers", "is_enable": false, "remark": ""},
		{"name": "ShowProcedureStatus", "is_enable": false, "remark": ""},
		{"name": "ShowIndex", "is_enable": false, "remark": ""},
		{"name": "ShowProcessList", "is_enable": false, "remark": ""},
		{"name": "ShowCreateDatabase", "is_enable": false, "remark": ""},
		{"name": "ShowConfig", "is_enable": false, "remark": ""},
		{"name": "ShowEvents", "is_enable": false, "remark": ""},
		{"name": "ShowStatsExtended", "is_enable": false, "remark": ""},
		{"name": "ShowStatsMeta", "is_enable": false, "remark": ""},
		{"name": "ShowStatsHistograms", "is_enable": false, "remark": ""},
		{"name": "ShowStatsTopN", "is_enable": false, "remark": ""},
		{"name": "ShowStatsBuckets", "is_enable": false, "remark": ""},
		{"name": "ShowStatsHealthy", "is_enable": false, "remark": ""},
		{"name": "ShowStatsLocked", "is_enable": false, "remark": ""},
		{"name": "ShowHistogramsInFlight", "is_enable": false, "remark": ""},
		{"name": "ShowColumnStatsUsage", "is_enable": false, "remark": ""},
		{"name": "ShowPlugins", "is_enable": false, "remark": ""},
		{"name": "ShowProfile", "is_enable": false, "remark": ""},
		{"name": "ShowProfiles", "is_enable": false, "remark": ""},
		{"name": "ShowMasterStatus", "is_enable": false, "remark": ""},
		{"name": "ShowPrivileges", "is_enable": false, "remark": ""},
		{"name": "ShowErrors", "is_enable": false, "remark": ""},
		{"name": "ShowBindings", "is_enable": false, "remark": ""},
		{"name": "ShowBindingCacheStatus", "is_enable": false, "remark": ""},
		{"name": "ShowPumpStatus", "is_enable": false, "remark": ""},
		{"name": "ShowDrainerStatus", "is_enable": false, "remark": ""},
		{"name": "ShowOpenTables", "is_enable": false, "remark": ""},
		{"name": "ShowAnalyzeStatus", "is_enable": false, "remark": ""},
		{"name": "ShowRegions", "is_enable": false, "remark": ""},
		{"name": "ShowBuiltins", "is_enable": false, "remark": ""},
		{"name": "ShowTableNextRowId", "is_enable": false, "remark": ""},
		{"name": "ShowBackups", "is_enable": false, "remark": ""},
		{"name": "ShowRestores", "is_enable": false, "remark": ""},
		{"name": "ShowPlacement", "is_enable": false, "remark": ""},
		{"name": "ShowPlacementForDatabase", "is_enable": false, "remark": ""},
		{"name": "ShowPlacementForTable", "is_enable": false, "remark": ""},
		{"name": "ShowPlacementForPartition", "is_enable": false, "remark": ""},
		{"name": "ShowPlacementLabels", "is_enable": false, "remark": ""},
		{"name": "ShowSessionStates", "is_enable": false, "remark": ""},
	}
	for _, i := range ops {
		var allowedOperations dasModels.InsightDASAllowedOperations
		_ = db.FirstOrCreate(&allowedOperations, dasModels.InsightDASAllowedOperations{
			Name:     i["name"].(string),
			IsEnable: i["is_enable"].(bool),
			Remark:   i["remark"].(string),
		})
	}
}

// 初始化审核参数
func initializeInspectParams(db *gorm.DB) {
	var params []map[string]interface{} = []map[string]interface{}{
		// TABLE
		{"key": "MAX_TABLE_NAME_LENGTH", "value": map[string]int{"value": 32}, "remark": "表名的长度"},
		{"key": "CHECK_TABLE_COMMENT", "value": map[string]bool{"value": true}, "remark": "检查表是否有注释"},
		{"key": "TABLE_COMMENT_LENGTH", "value": map[string]int{"value": 64}, "remark": "表注释的长度"},
		{"key": "CHECK_IDENTIFIER", "value": map[string]bool{"value": true}, "remark": "对象名必须使用字符串范围为正则[a-zA-Z0-9_]"},
		{"key": "CHECK_IDENTIFER_KEYWORD", "value": map[string]bool{"value": false}, "remark": "对象名是否可以使用关键字"},
		{"key": "CHECK_TABLE_CHARSET", "value": map[string]bool{"value": true}, "remark": "是否检查表的字符集和排序规则"},
		{"key": "TABLE_SUPPORT_CHARSET", "value": map[string][]map[string]string{"value": []map[string]string{
			{"charset": "utf8", "recommend": "utf8_general_ci"},
			{"charset": "utf8mb4", "recommend": "utf8mb4_general_ci"},
		}}, "remark": "表支持的字符集"},
		{"key": "CHECK_TABLE_ENGINE", "value": map[string]bool{"value": true}, "remark": "是否检查表的存储引擎"},
		{"key": "TABLE_SUPPORT_ENGINE", "value": []string{"InnoDB"}, "remark": "表支持的存储引擎"},
		{"key": "ENABLE_PARTITION_TABLE", "value": map[string]bool{"value": false}, "remark": "是否启用分区表"},
		{"key": "CHECK_TABLE_PRIMARY_KEY", "value": map[string]bool{"value": true}, "remark": "检查表是否有主键"},
		{"key": "TABLE_AT_LEAST_ONE_COLUMN", "value": map[string]bool{"value": true}, "remark": "表至少要有一列，语法默认支持"},
		{"key": "CHECK_TABLE_AUDIT_TYPE_COLUMNS", "value": map[string]bool{"value": true}, "remark": "启用审计类型的字段(col1 datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP && col2 datetime DEFAULT CURRENT_TIMESTAMP)"},
		{"key": "ENABLE_CREATE_TABLE_AS", "value": map[string]bool{"value": false}, "remark": "是否允许create table as语法"},
		{"key": "ENABLE_CREATE_TABLE_LIKE", "value": map[string]bool{"value": false}, "remark": "是否允许create table like语法"},
		{"key": "ENABLE_FOREIGN_KEY", "value": map[string]bool{"value": false}, "remark": "是否启用外键"},
		{"key": "CHECK_TABLE_AUTOINCREMENT_INIT_VALUE", "value": map[string]bool{"value": true}, "remark": "检查建表是自增列初始值是否为1"},
		{"key": "ENABLE_CREATE_VIEW", "value": map[string]bool{"value": true}, "remark": "是否支持创建和使用视图"},
		// COLUMN
		{"key": "MAX_COLUMN_NAME_LENGTH", "value": map[string]int{"value": 64}, "remark": "列名的长度"},
		{"key": "CHECK_COLUMN_CHARSET", "value": map[string]bool{"value": true}, "remark": "是否检查列的字符集"},
		{"key": "CHECK_COLUMN_COMMENT", "value": map[string]bool{"value": true}, "remark": "是否检查列的注释"},
		{"key": "COLUMN_MAX_CHAR_LENGTH", "value": map[string]int{"value": 64}, "remark": "char长度大于N的时候需要改为varchar"},
		{"key": "MAX_VARCHAR_LENGTH", "value": map[string]int{"value": 65535}, "remark": "最大允许定义的varchar长度"},
		{"key": "ENABLE_COLUMN_BLOB_TYPE", "value": map[string]bool{"value": true}, "remark": "是否允许列的类型为BLOB/TEXT"},
		{"key": "ENABLE_COLUMN_JSON_TYPE", "value": map[string]bool{"value": true}, "remark": "是否允许列的类型为JSON"},
		{"key": "ENABLE_COLUMN_BIT_TYPE", "value": map[string]bool{"value": true}, "remark": "是否允许列的类型为BIT"},
		{"key": "ENABLE_COLUMN_TIMESTAMP_TYPE", "value": map[string]bool{"value": false}, "remark": "是否允许列的类型为TIMESTAMP"},
		{"key": "CHECK_PRIMARYKEY_USE_BIGINT", "value": map[string]bool{"value": true}, "remark": "主键是否为bigint"},
		{"key": "CHECK_PRIMARYKEY_USE_UNSIGNED", "value": map[string]bool{"value": true}, "remark": "主键bigint是否为unsigned"},
		{"key": "CHECK_PRIMARYKEY_USE_AUTO_INCREMENT", "value": map[string]bool{"value": true}, "remark": "主键是否定义为自增"},
		{"key": "ENABLE_COLUMN_NOT_NULL", "value": map[string]bool{"value": true}, "remark": "是否允许列定义为NOT NULL"},
		{"key": "ENABLE_COLUMN_TIME_NULL", "value": map[string]bool{"value": true}, "remark": "是否允许时间类型设置为NULL"},
		{"key": "CHECK_COLUMN_DEFAULT_VALUE", "value": map[string]bool{"value": true}, "remark": "列必须要有默认值"},
		{"key": "CHECK_COLUMN_FLOAT_DOUBLE", "value": map[string]bool{"value": true}, "remark": "将float/double转成int/bigint/decimal等"},
		{"key": "ENABLE_COLUMN_TYPE_CHANGE", "value": map[string]bool{"value": false}, "remark": "是否允许变更列类型"},
		{"key": "ENABLE_COLUMN_TYPE_CHANGE_COMPATIBLE", "value": map[string]bool{"value": true}, "remark": "允许tinyint-> int、int->bigint、char->varchar等"},
		{"key": "ENABLE_COLUMN_CHANGE_COLUMN_NAME", "value": map[string]bool{"value": false}, "remark": "是否允许CHANGE修改列名操作"},
		// INDEX
		{"key": "CHECK_UNIQ_INDEX_PREFIX", "value": map[string]bool{"value": true}, "remark": "是否检查唯一索引前缀，如唯一索引必须以uniq_为前缀"},
		{"key": "CHECK_SECONDARY_INDEX_PREFIX", "value": map[string]bool{"value": true}, "remark": "是否检查二级索引前缀，如普通索引必须以idx_为前缀"},
		{"key": "CHECK_FULLTEXT_INDEX_PREFIX", "value": map[string]bool{"value": true}, "remark": "是否检查全文索引前缀，如全文索引必须以full_为前缀"},
		{"key": "UNQI_INDEX_PREFIX", "value": map[string]string{"value": "UNIQ_"}, "remark": "定义唯一索引前缀，不区分大小写"},
		{"key": "SECONDARY_INDEX_PREFIX", "value": map[string]string{"value": "IDX_"}, "remark": "定义二级索引前缀，不区分大小写"},
		{"key": "FULLTEXT_INDEX_PREFIX", "value": map[string]string{"value": "FULL_"}, "remark": "定义全文索引前缀，不区分大小写"},
		{"key": "SECONDARY_INDEX_MAX_KEY_PARTS", "value": map[string]int{"value": 8}, "remark": "组成二级索引的列数不能超过指定的个数,包括唯一索引"},
		{"key": "PRIMARYKEY_MAX_KEY_PARTS", "value": map[string]int{"value": 1}, "remark": "组成主键索引的列数不能超过指定的个数"},
		{"key": "MAX_INDEX_KEYS", "value": map[string]int{"value": 12}, "remark": "最多有N个索引，包括唯一索引/二级索引"},
		{"key": "ENABLE_INDEX_RENAME", "value": map[string]bool{"value": false}, "remark": "是否允许rename索引名"},
		// ALTER
		{"key": "ENABLE_DROP_COLS", "value": map[string]bool{"value": true}, "remark": "是否允许DROP列"},
		{"key": "ENABLE_DROP_INDEXES", "value": map[string]bool{"value": true}, "remark": "是否允许DROP索引"},
		{"key": "ENABLE_DROP_PRIMARYKEY", "value": map[string]bool{"value": false}, "remark": "是否允许DROP主键"},
		{"key": "ENABLE_DROP_TABLE", "value": map[string]bool{"value": true}, "remark": "是否允许DROP TABLE"},
		{"key": "ENABLE_TRUNCATE_TABLE", "value": map[string]bool{"value": true}, "remark": "是否允许TRUNCATE TABLE"},
		{"key": "ENABLE_RENAME_TABLE_NAME", "value": map[string]bool{"value": false}, "remark": "是否允许rename表名"},
		{"key": "ENABLE_MYSQL_MERGE_ALTER_TABLE", "value": map[string]bool{"value": true}, "remark": "MySQL同一个表的多个ALTER是否合并为单条语句"},
		{"key": "ENABLE_TIDB_MERGE_ALTER_TABLE", "value": map[string]bool{"value": false}, "remark": "TiDB同一个表的多个ALTER是否合并为单条语句"},
		// DML
		{"key": "DML_MUST_HAVE_WHERE", "value": map[string]bool{"value": true}, "remark": "DML语句必须有where条件"},
		{"key": "DML_DISABLE_LIMIT", "value": map[string]bool{"value": true}, "remark": "DML语句中不允许有LIMIT"},
		{"key": "DML_DISABLE_ORDERBY", "value": map[string]bool{"value": true}, "remark": "DML语句中不允许有orderby"},
		{"key": "DML_DISABLE_SUBQUERY", "value": map[string]bool{"value": true}, "remark": "DML语句不能有子查询"},
		{"key": "CHECK_DML_JOIN_WITH_ON", "value": map[string]bool{"value": true}, "remark": "DML的JOIN语句必须有ON语句"},
		{"key": "EXPLAIN_RULE", "value": map[string]string{"value": "first"}, "remark": "explain判断受影响行数时使用的规则('first', 'max')。 'first': 使用第一行的explain结果作为受影响行数, 'max': 使用explain结果中的最大值作为受影响行数"},
		{"key": "MAX_AFFECTED_ROWS", "value": map[string]int{"value": 100}, "remark": "最大影响行数，默认100"},
		{"key": "MAX_INSERT_ROWS", "value": map[string]int{"value": 100}, "remark": " 一次最多允许insert的行, eg: insert into tbl(col,...) values(row1), (row2)..."},
		{"key": "DISABLE_REPLACE", "value": map[string]bool{"value": true}, "remark": "是否禁用replace语句"},
		{"key": "DISABLE_INSERT_INTO_SELECT", "value": map[string]bool{"value": true}, "remark": "是否禁用insert/replace into select语法"},
		{"key": "DISABLE_ON_DUPLICATE", "value": map[string]bool{"value": true}, "remark": "是否禁止insert on duplicate语法"},
		// 禁止语法审核的表
		{"key": "DISABLE_AUDIT_DML_TABLES", "value": map[string]interface{}{"value": []map[string]interface{}{
			{"DB": "d1", "Tables": []string{"t1", "t2"}, "Reason": "研发禁止审核和提交"},
			{"DB": "d2", "Tables": []string{"t1", "t2"}, "Reason": "研发禁止审核和提交"},
		}}, "remark": "禁止指定的表的DML语句进行审核"},
		{"key": "DISABLE_AUDIT_DDL_TABLES", "value": map[string]interface{}{"value": []map[string]interface{}{
			{"DB": "d1", "Tables": []string{"t1", "t2"}, "Reason": "研发禁止审核和提交"},
			{"DB": "d2", "Tables": []string{"t1", "t2"}, "Reason": "研发禁止审核和提交"},
		}}, "remark": "禁止指定的表的DDL语句进行审核"},
	}
	for _, i := range params {
		var inspectParams inspectModels.InsightInspectParams
		jsonParams, err := json.Marshal(i["value"])
		if err != nil {
			global.App.Log.Error(err)
			panic(err)
		}
		_ = db.FirstOrCreate(&inspectParams, inspectModels.InsightInspectParams{
			Key:    i["key"].(string),
			Value:  jsonParams,
			Remark: i["remark"].(string),
		})
	}
}

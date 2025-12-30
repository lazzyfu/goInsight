package bootstrap

import (
	"context"
	"fmt"
	"time"

	"github.com/lazzyfu/goinsight/internal/global"

	commonModels "github.com/lazzyfu/goinsight/internal/common/models"
	dasModels "github.com/lazzyfu/goinsight/internal/das/models"
	inspectModels "github.com/lazzyfu/goinsight/internal/inspect/models"
	ordersModels "github.com/lazzyfu/goinsight/internal/orders/models"
	usersModels "github.com/lazzyfu/goinsight/internal/users/models"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

// 初始化redis
func InitializeRedis() *redis.Client {
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

	// 健康检查PING
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		global.App.Log.Error("redis connect failed, err:", err.Error())
		panic(fmt.Sprintf("redis connect failed, err: %s", err.Error()))
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
		initializeGlobalInspectParams(db)
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
		&usersModels.InsightOrgs{},
		&usersModels.InsightOrgUsers{},
		// common
		&commonModels.InsightInstanceEnvironments{},
		&commonModels.InsightInstances{},
		&commonModels.InsightInstanceSchemas{},
		// inspect
		&inspectModels.InsightInspectGlobalParams{},
		&inspectModels.InsightInspectInstanceParams{},
		// das
		&dasModels.InsightDasSchemaPerms{},
		&dasModels.InsightDasTablePerms{},
		&dasModels.InsightDASRecords{},
		&dasModels.InsightDASOperations{},
		&dasModels.InsightDASFavorites{},
		// orders
		&ordersModels.InsightOrderRecords{},
		&ordersModels.InsightOrderTasks{},
		&ordersModels.InsightOrderMessages{},
		&ordersModels.InsightApprovalFlows{},
		&ordersModels.InsightApprovalRecords{},
		&ordersModels.InsightApprovalFlowUsers{},
		&ordersModels.InsightOrderLogs{},
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
	var ops []map[string]any = []map[string]any{
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
		var allowedOperations dasModels.InsightDASOperations
		_ = db.FirstOrCreate(&allowedOperations, dasModels.InsightDASOperations{
			Name:     i["name"].(string),
			IsEnable: i["is_enable"].(bool),
			Remark:   i["remark"].(string),
		})
	}
}

// 初始化审核参数
func initializeGlobalInspectParams(db *gorm.DB) {
	var params []map[string]any = []map[string]any{
		// TABLE
		{"title": "检查表名长度", "Key": "MAX_TABLE_NAME_LENGTH", "Value": "32", "Type": "number"},
		{"title": "检查表是否有注释", "Key": "CHECK_TABLE_COMMENT", "Value": "true", "Type": "boolean"},
		{"title": "检查表注释的长度", "Key": "TABLE_COMMENT_LENGTH", "Value": "64", "Type": "number"},
		{"title": "对象名必须使用字符串范围为正则[a-zA-Z0-9_]", "Key": "CHECK_IDENTIFIER", "Value": "true", "Type": "boolean"},
		{"title": "对象名是否可以使用关键字", "Key": "CHECK_IDENTIFER_KEYWORD", "Value": "false", "Type": "boolean"},
		{"title": "是否检查表的字符集和排序规则", "Key": "CHECK_TABLE_CHARSET", "Value": "true", "Type": "boolean"},
		{"title": "表支持的字符集", "Key": "TABLE_SUPPORT_CHARSET", "Value": "utf8,utf8_general_ci;utf8mb4,utf8mb4_general_ci", "Type": "string"},
		{"title": "是否检查表的存储引擎", "Key": "CHECK_TABLE_ENGINE", "Value": "true", "Type": "boolean"},
		{"title": "表支持的存储引擎", "Key": "TABLE_SUPPORT_ENGINE", "Value": "InnoDB,MyISAM", "Type": "string"},
		{"title": "是否启用分区表", "Key": "ENABLE_PARTITION_TABLE", "Value": "false", "Type": "boolean"},
		{"title": "检查表是否有主键", "Key": "CHECK_TABLE_PRIMARY_KEY", "Value": "true", "Type": "boolean"},
		{"title": "表至少要有一列，语法默认支持", "Key": "TABLE_AT_LEAST_ONE_COLUMN", "Value": "true", "Type": "boolean"},
		{"title": "启用审计类型的字段(col1 datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP && col2 datetime DEFAULT CURRENT_TIMESTAMP)", "Key": "CHECK_TABLE_AUDIT_TYPE_COLUMNS", "Value": "true", "Type": "boolean"},
		{"title": "是否允许create table as语法", "Key": "ENABLE_CREATE_TABLE_AS", "Value": "false", "Type": "boolean"},
		{"title": "是否允许create table like语法", "Key": "ENABLE_CREATE_TABLE_LIKE", "Value": "false", "Type": "boolean"},
		{"title": "是否启用外键", "Key": "ENABLE_FOREIGN_KEY", "Value": "false", "Type": "boolean"},
		{"title": "检查建表是自增列初始值是否为1", "Key": "CHECK_TABLE_AUTOINCREMENT_INIT_VALUE", "Value": "true", "Type": "boolean"},
		{"title": "是否支持创建和使用视图", "Key": "ENABLE_CREATE_VIEW", "Value": "true", "Type": "boolean"},
		{"title": "InnoDB表支持的行格式", "Key": "INNODB_ROW_FORMAT", "Value": "DYNAMIC,COMPACT,REDUNDANT", "Type": "string"},
		// COLUMN
		{"title": "列名的长度", "Key": "MAX_COLUMN_NAME_LENGTH", "Value": "64", "Type": "number"},
		{"title": "是否检查列的字符集", "Key": "CHECK_COLUMN_CHARSET", "Value": "true", "Type": "boolean"},
		{"title": "是否检查列的注释", "Key": "CHECK_COLUMN_COMMENT", "Value": "true", "Type": "boolean"},
		{"title": "char长度大于N的时候需要改为varchar", "Key": "COLUMN_MAX_CHAR_LENGTH", "Value": "64", "Type": "number"},
		{"title": "最大允许定义的varchar长度", "Key": "MAX_VARCHAR_LENGTH", "Value": "16383", "Type": "number"},
		{"title": "是否允许列的类型为BLOB/TEXT", "Key": "ENABLE_COLUMN_BLOB_TYPE", "Value": "true", "Type": "boolean"},
		{"title": "是否允许列的类型为JSON", "Key": "ENABLE_COLUMN_JSON_TYPE", "Value": "true", "Type": "boolean"},
		{"title": "是否允许列的类型为BIT", "Key": "ENABLE_COLUMN_BIT_TYPE", "Value": "true", "Type": "boolean"},
		{"title": "是否允许列的类型为TIMESTAMP", "Key": "ENABLE_COLUMN_TIMESTAMP_TYPE", "Value": "false", "Type": "boolean"},
		{"title": "主键是否为bigint", "Key": "CHECK_PRIMARYKEY_USE_BIGINT", "Value": "true", "Type": "boolean"},
		{"title": "主键bigint是否为unsigned", "Key": "CHECK_PRIMARYKEY_USE_UNSIGNED", "Value": "true", "Type": "boolean"},
		{"title": "主键是否定义为自增", "Key": "CHECK_PRIMARYKEY_USE_AUTO_INCREMENT", "Value": "true", "Type": "boolean"},
		{"title": "是否允许列定义为NOT NULL", "Key": "ENABLE_COLUMN_NOT_NULL", "Value": "true", "Type": "boolean"},
		{"title": "是否允许时间类型设置为NULL", "Key": "ENABLE_COLUMN_TIME_NULL", "Value": "true", "Type": "boolean"},
		{"title": "列必须要有默认值", "Key": "CHECK_COLUMN_DEFAULT_VALUE", "Value": "true", "Type": "boolean"},
		{"title": "将float/double转成int/bigint/decimal等", "Key": "CHECK_COLUMN_FLOAT_DOUBLE", "Value": "true", "Type": "boolean"},
		{"title": "是否允许变更列类型", "Key": "ENABLE_COLUMN_TYPE_CHANGE", "Value": "false", "Type": "boolean"},
		{"title": "是否允许通过兼容的类型变更列类型", "Key": "ENABLE_COLUMN_TYPE_CHANGE_COMPATIBLE", "Value": "true", "Type": "boolean"},
		{"title": "是否允许CHANGE修改列名操作", "Key": "ENABLE_COLUMN_CHANGE_COLUMN_NAME", "Value": "false", "Type": "boolean"},
		// INDEX
		{"title": "是否检查唯一索引前缀，如唯一索引必须以uniq_为前缀", "Key": "CHECK_UNIQ_INDEX_PREFIX", "Value": "true", "Type": "boolean"},
		{"title": "是否检查二级索引前缀，如普通索引必须以idx_为前缀", "Key": "CHECK_SECONDARY_INDEX_PREFIX", "Value": "true", "Type": "boolean"},
		{"title": "是否检查全文索引前缀，如全文索引必须以full_为前缀", "Key": "CHECK_FULLTEXT_INDEX_PREFIX", "Value": "true", "Type": "boolean"},
		{"title": "定义唯一索引前缀，不区分大小写", "Key": "UNQI_INDEX_PREFIX", "Value": "UNIQ_", "Type": "string"},
		{"title": "定义二级索引前缀，不区分大小写", "Key": "SECONDARY_INDEX_PREFIX", "Value": "IDX_", "Type": "string"},
		{"title": "定义全文索引前缀，不区分大小写", "Key": "FULLTEXT_INDEX_PREFIX", "Value": "FULL_", "Type": "string"},
		{"title": "组成二级索引的列数不能超过指定的个数,包括唯一索引", "Key": "SECONDARY_INDEX_MAX_KEY_PARTS", "Value": "8", "Type": "number"},
		{"title": "组成主键索引的列数不能超过指定的个数", "Key": "PRIMARYKEY_MAX_KEY_PARTS", "Value": "1", "Type": "number"},
		{"title": "最多有N个索引，包括唯一索引/二级索引", "Key": "MAX_INDEX_KEYS", "Value": "12", "Type": "number"},
		{"title": "是否允许rename索引名", "Key": "ENABLE_INDEX_RENAME", "Value": "false", "Type": "boolean"},
		{"title": "是否允许冗余索引", "Key": "ENABLE_REDUNDANT_INDEX", "Value": "false", "Type": "boolean"},
		// ALTER
		{"title": "是否允许DROP列", "Key": "ENABLE_DROP_COLS", "Value": "true", "Type": "boolean"},
		{"title": "是否允许DROP索引", "Key": "ENABLE_DROP_INDEXES", "Value": "true", "Type": "boolean"},
		{"title": "是否允许DROP主键", "Key": "ENABLE_DROP_PRIMARYKEY", "Value": "false", "Type": "boolean"},
		{"title": "是否允许DROP TABLE", "Key": "ENABLE_DROP_TABLE", "Value": "true", "Type": "boolean"},
		{"title": "是否允许TRUNCATE TABLE", "Key": "ENABLE_TRUNCATE_TABLE", "Value": "true", "Type": "boolean"},
		{"title": "是否允许rename表名", "Key": "ENABLE_RENAME_TABLE_NAME", "Value": "false", "Type": "boolean"},
		{"title": "MySQL同一个表的多个ALTER是否合并为单条语句", "Key": "ENABLE_MYSQL_MERGE_ALTER_TABLE", "Value": "true", "Type": "boolean"},
		{"title": "TiDB同一个表的多个ALTER是否合并为单条语句", "Key": "ENABLE_TIDB_MERGE_ALTER_TABLE", "Value": "false", "Type": "boolean"},
		// DML
		{"title": "DML语句必须有where条件", "Key": "DML_MUST_HAVE_WHERE", "Value": "true", "Type": "boolean"},
		{"title": "DML语句中不允许有LIMIT", "Key": "DML_DISABLE_LIMIT", "Value": "true", "Type": "boolean"},
		{"title": "DML语句中不允许有orderby", "Key": "DML_DISABLE_ORDERBY", "Value": "true", "Type": "boolean"},
		{"title": "DML语句不能有子查询", "Key": "DML_DISABLE_SUBQUERY", "Value": "true", "Type": "boolean"},
		{"title": "DML的JOIN语句必须有ON语句", "Key": "CHECK_DML_JOIN_WITH_ON", "Value": "true", "Type": "boolean"},
		{"title": "explain判断受影响行数时使用的规则('first', 'max')。 'first': 使用第一行的explain结果作为受影响行数, 'max': 使用explain结果中的最大值作为受影响行数", "Key": "EXPLAIN_RULE", "Value": "first", "Type": "string"},
		{"title": "最大影响行数，默认100", "Key": "MAX_AFFECTED_ROWS", "Value": "100", "Type": "number"},
		{"title": "一次最多允许insert的行, eg: insert into tbl(col,...) values(row1), (row2)...", "Key": "MAX_INSERT_ROWS", "Value": "100", "Type": "number"},
		{"title": "是否禁用replace语句", "Key": "DISABLE_REPLACE", "Value": "true", "Type": "boolean"},
		{"title": "是否禁用insert/replace into select语法", "Key": "DISABLE_INSERT_INTO_SELECT", "Value": "true", "Type": "boolean"},
		{"title": "是否禁止insert on duplicate语法", "Key": "DISABLE_ON_DUPLICATE", "Value": "true", "Type": "boolean"},
	}

	for _, param := range params {
		var inspectGlobalParam inspectModels.InsightInspectGlobalParams
		err := db.Where("`key` = ?", param["Key"].(string)).First(&inspectGlobalParam).Error
		if err == gorm.ErrRecordNotFound {
			newParam := inspectModels.InsightInspectGlobalParams{
				Title: param["title"].(string),
				Key:   param["Key"].(string),
				Value: param["Value"].(string),
				Type:  commonModels.EnumType(param["Type"].(string)),
			}
			_ = db.Create(&newParam)
		}
	}
}

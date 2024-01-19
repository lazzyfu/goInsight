/*
@Time    :   2023/08/14 15:53:53
@Author  :   lazzyfu
*/

package bootstrap

import (
	"context"
	"fmt"
	"goInsight/global"
	commonModels "goInsight/internal/app/common/models"
	dasModels "goInsight/internal/app/das/models"
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
		IsSuperuser: true,
		IsActive:    true,
		IsStaff:     false,
		IsTwoFA:     false,
	})
}

// 初始化用户允许的操作
func initializeAllowedOperations(db *gorm.DB) {
	var ops []map[string]string = []map[string]string{
		{"name": "SELECT", "remark": ""},
		{"name": "UNION", "remark": ""},
		{"name": "Use", "remark": ""},
		{"name": "Desc", "remark": ""},
		{"name": "ExplainSelect", "remark": ""},
		{"name": "ExplainDelete", "remark": ""},
		{"name": "ExplainInsert", "remark": ""},
		{"name": "ExplainUpdate", "remark": ""},
		{"name": "ExplainUnion", "remark": ""},
		{"name": "ExplainFor", "remark": "ExplainForStmt is a statement to provite information about how is SQL statement executeing in connection #ConnectionID"},
		{"name": "ShowEngines", "remark": ""},
		{"name": "ShowDatabases", "remark": ""},
		{"name": "ShowTables", "remark": ""},
		{"name": "ShowTableStatus", "remark": ""},
		{"name": "ShowColumns", "remark": ""},
		{"name": "ShowWarnings", "remark": ""},
		{"name": "ShowCharset", "remark": ""},
		{"name": "ShowVariables", "remark": ""},
		{"name": "ShowStatus", "remark": ""},
		{"name": "ShowCollation", "remark": ""},
		{"name": "ShowCreateTable", "remark": ""},
		{"name": "ShowCreateView", "remark": ""},
		{"name": "ShowCreateUser", "remark": ""},
		{"name": "ShowCreateSequence", "remark": ""},
		{"name": "ShowCreatePlacementPolicy", "remark": ""},
		{"name": "ShowGrants", "remark": ""},
		{"name": "ShowTriggers", "remark": ""},
		{"name": "ShowProcedureStatus", "remark": ""},
		{"name": "ShowIndex", "remark": ""},
		{"name": "ShowProcessList", "remark": ""},
		{"name": "ShowCreateDatabase", "remark": ""},
		{"name": "ShowConfig", "remark": ""},
		{"name": "ShowEvents", "remark": ""},
		{"name": "ShowStatsExtended", "remark": ""},
		{"name": "ShowStatsMeta", "remark": ""},
		{"name": "ShowStatsHistograms", "remark": ""},
		{"name": "ShowStatsTopN", "remark": ""},
		{"name": "ShowStatsBuckets", "remark": ""},
		{"name": "ShowStatsHealthy", "remark": ""},
		{"name": "ShowStatsLocked", "remark": ""},
		{"name": "ShowHistogramsInFlight", "remark": ""},
		{"name": "ShowColumnStatsUsage", "remark": ""},
		{"name": "ShowPlugins", "remark": ""},
		{"name": "ShowProfile", "remark": ""},
		{"name": "ShowProfiles", "remark": ""},
		{"name": "ShowMasterStatus", "remark": ""},
		{"name": "ShowPrivileges", "remark": ""},
		{"name": "ShowErrors", "remark": ""},
		{"name": "ShowBindings", "remark": ""},
		{"name": "ShowBindingCacheStatus", "remark": ""},
		{"name": "ShowPumpStatus", "remark": ""},
		{"name": "ShowDrainerStatus", "remark": ""},
		{"name": "ShowOpenTables", "remark": ""},
		{"name": "ShowAnalyzeStatus", "remark": ""},
		{"name": "ShowRegions", "remark": ""},
		{"name": "ShowBuiltins", "remark": ""},
		{"name": "ShowTableNextRowId", "remark": ""},
		{"name": "ShowBackups", "remark": ""},
		{"name": "ShowRestores", "remark": ""},
		{"name": "ShowPlacement", "remark": ""},
		{"name": "ShowPlacementForDatabase", "remark": ""},
		{"name": "ShowPlacementForTable", "remark": ""},
		{"name": "ShowPlacementForPartition", "remark": ""},
		{"name": "ShowPlacementLabels", "remark": ""},
		{"name": "ShowSessionStates", "remark": ""},
	}
	for _, i := range ops {
		var allowedOperations dasModels.InsightDASAllowedOperations
		_ = db.FirstOrCreate(&allowedOperations, dasModels.InsightDASAllowedOperations{Name: i["name"], Remark: i["remark"]})
	}
}

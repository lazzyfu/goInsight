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

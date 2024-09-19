package dao

import (
	"errors"
	"fmt"
	"goInsight/internal/inspect/controllers/parser"
	"goInsight/pkg/kv"
	"goInsight/pkg/utils"
	"strconv"
	"strings"

	mysqlapi "github.com/go-sql-driver/mysql"
)

// ShowCreateTable retrieves the structure of the specified table.
func ShowCreateTable(table string, db *DB, kv *kv.KVCache) (data interface{}, err error) {
	// Return table structure from cache if available
	data = kv.Get(table)
	if data != nil {
		return data, nil
	}
	query := fmt.Sprintf("SHOW CREATE TABLE `%s`", table)
	result, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	var createStatement string
	for _, sql := range *result {
		// Table
		if _, ok := sql["Create Table"]; ok {
			createStatement = sql["Create Table"].(string)
		}
		// View
		if _, ok := sql["Create View"]; ok {
			createStatement = sql["Create View"].(string)
		}
	}

	var warns []error
	data, warns, err = parser.NewParse(createStatement, "", "")
	if len(warns) > 0 {
		return nil, fmt.Errorf("解析警告: %s", utils.ErrsJoin("; ", warns))
	}
	if err != nil {
		return nil, fmt.Errorf("SQL语法解析错误: %s", err.Error())
	}
	kv.Put(table, data)
	return data, nil
}

// CheckIfTableExists checks if the specified table exists in the current database.
func CheckIfTableExists(table string, db *DB) (string, error) {
	err := db.Execute(fmt.Sprintf("DESC `%s`", table))
	if me, ok := err.(*mysqlapi.MySQLError); ok {
		if me.Number == 1146 {
			// Table does not exist
			return fmt.Sprintf("表或视图`%s`不存在", table), err
		} else if me.Number == 1045 {
			return fmt.Sprintf("无法访问目标数据库 %s:%d, %s", db.Host, db.Port, err.Error()), err
		}
	}
	return fmt.Sprintf("表或视图`%s`已存在", table), nil
}

// CheckIfDatabaseExists checks if the specified database exists.
func CheckIfDatabaseExists(database string, db *DB) (string, error) {
	// Query the information_schema.schemata to check if the database exists
	result, err := db.Query(fmt.Sprintf("SELECT COUNT(*) as count FROM information_schema.schemata WHERE schema_name='%s'", database))
	if err != nil {
		return fmt.Sprintf("执行SQL失败, 主机: %s:%d, 错误: %s", db.Host, db.Port, err.Error()), err
	}
	var count int
	for _, row := range *result {
		count, _ = strconv.Atoi(row["count"].(string))
		break
	}
	if count == 0 {
		// Database does not exist
		return fmt.Sprintf("数据库`%s`不存在", database), errors.New("error")
	}
	// Database exists
	return fmt.Sprintf("数据库`%s`已存在", database), nil
}

// CheckIfTableExistsCrossDB checks if the specified table exists across databases.
func CheckIfTableExistsCrossDB(table string, db *DB) (string, error) {
	// Check if the table exists using information_schema.tables, suitable for cross-database checks
	result, err := db.Query(fmt.Sprintf("SELECT COUNT(*) as count FROM information_schema.tables WHERE table_name='%s'", table))
	if err != nil {
		return fmt.Sprintf("执行SQL失败, 主机: %s:%d, 错误: %s", db.Host, db.Port, err.Error()), err
	}
	var count int
	for _, row := range *result {
		count, _ = strconv.Atoi(row["count"].(string))
		break
	}
	if count == 0 {
		// Table does not exist
		return fmt.Sprintf("表或视图`%s`不存在", table), errors.New("error")
	}
	// Table exists
	return fmt.Sprintf("表或视图`%s`已存在", table), nil
}

// GetDBVars retrieves database variables.
func GetDBVars(db *DB) (map[string]string, error) {
	result, err := db.Query(`SHOW VARIABLES WHERE Variable_name IN ('innodb_large_prefix','version','character_set_database','innodb_default_row_format')`)
	if err != nil {
		return nil, err
	}

	var data map[string]string = map[string]string{
		"dbVersion":              "",
		"dbCharset":              "utf8",
		"largePrefix":            "OFF",
		"innodbDefaultRowFormat": "dynamic",
	}

	// [map[Value:utf8 Variable_name:character_set_database] map[Value:5.7.35-log Variable_name:version]]
	for _, row := range *result {
		variableName, ok := row["Variable_name"].(string)
		if !ok {
			return nil, fmt.Errorf("Variable_name类型意外")
		}

		value, ok := row["Value"].(string)
		if !ok {
			return nil, fmt.Errorf("行中Value类型意外")
		}

		switch variableName {
		case "version":
			data["dbVersion"] = value
		case "character_set_database":
			data["dbCharset"] = value
		case "innodb_large_prefix":
			switch value {
			case "0":
				data["largePrefix"] = "OFF"
			case "1":
				data["largePrefix"] = "ON"
			default:
				data["largePrefix"] = strings.ToUpper(value)
			}
		case "innodb_default_row_format":
			data["innodbDefaultRowFormat"] = value
		}
	}
	return data, nil
}

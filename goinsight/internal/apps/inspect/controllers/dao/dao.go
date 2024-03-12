/*
@Time    :   2022/06/24 10:18:49
@Author  :   zongfei.fu
@Desc    :   获取目标数据库元信息
*/

package dao

import (
	"errors"
	"fmt"
	"goInsight/internal/apps/inspect/controllers/parser"
	"goInsight/internal/pkg/kv"
	"goInsight/internal/pkg/utils"
	"strconv"
	"strings"

	mysqlapi "github.com/go-sql-driver/mysql"
)

// ShowCreateTable
func ShowCreateTable(table string, db *DB, kv *kv.KVCache) (data interface{}, err error) {
	// 返回表结构
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
		// 表
		if _, ok := sql["Create Table"]; ok {
			createStatement = sql["Create Table"].(string)
		}
		// 视图
		if _, ok := sql["Create View"]; ok {
			createStatement = sql["Create View"].(string)
		}
	}

	var warns []error
	data, warns, err = parser.NewParse(createStatement, "", "")
	if len(warns) > 0 {
		return nil, fmt.Errorf("Parse Warning: %s", utils.ErrsJoin("; ", warns))
	}
	if err != nil {
		return nil, fmt.Errorf("SQL语法解析错误：%s", err.Error())
	}
	kv.Put(table, data)
	return data, nil
}

// descTable
func DescTable(table string, db *DB) (error, string) {
	// 检查表是否存在，适用于确认当前实例当前库的表
	err := db.Execute(fmt.Sprintf("desc `%s`", table))
	if me, ok := err.(*mysqlapi.MySQLError); ok {
		if me.Number == 1146 {
			// 表不存在
			return err, fmt.Sprintf("表或视图`%s`不存在", table)
		} else if me.Number == 1045 {
			return err, fmt.Sprintf("访问目标数据库%s:%d失败,%s", db.Host, db.Port, err.Error())
		}
	}
	return nil, fmt.Sprintf("表或视图`%s`已经存在", table)
}

// verifyTable
func VerifyTable(table string, db *DB) (error, string) {
	// 通过information_schema.tables检查表是否存在，适用于确认当前实例跨库的表
	result, err := db.Query(fmt.Sprintf("select count(*) as count from information_schema.tables where table_name='%s'", table))
	if err != nil {
		return err, fmt.Sprintf("执行SQL失败,主机:%s:%d,错误:%s", db.Host, db.Port, err.Error())
	}
	var count int
	for _, row := range *result {
		count, _ = strconv.Atoi(row["count"].(string))
		break
	}
	if count == 0 {
		// 表不存在
		return errors.New("error"), fmt.Sprintf("表或视图`%s`不存在", table)
	}
	// 表存在
	return nil, fmt.Sprintf("表或视图`%s`已经存在", table)
}

// 获取DB变量
func GetDBVars(db *DB) (map[string]string, error) {
	result, err := db.Query("show variables where Variable_name in  ('innodb_large_prefix','version','character_set_database','innodb_default_row_format')")
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
		if strings.EqualFold(row["Variable_name"].(string), "version") {
			data["dbVersion"] = row["Value"].(string)
		}
		if strings.EqualFold(row["Variable_name"].(string), "character_set_database") {
			data["dbCharset"] = row["Value"].(string)
		}
		if strings.EqualFold(row["Variable_name"].(string), "innodb_large_prefix") {
			switch row["Value"].(string) {
			case "0":
				data["largePrefix"] = "OFF"
			case "1":
				data["largePrefix"] = "ON"
			default:
				data["largePrefix"] = strings.ToUpper(row["Value"].(string))
			}
		}
		if strings.EqualFold(row["Variable_name"].(string), "innodb_default_row_format") {
			data["innodbDefaultRowFormat"] = row["Value"].(string)
		}
	}
	return data, nil
}

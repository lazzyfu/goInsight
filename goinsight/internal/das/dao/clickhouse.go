/*
@Time    :   2023/04/11 11:27:36
@Author  :   xff
@Desc    :   如果clickhouse版本大于22.3，请使用高版本的clickhouse-go
			 参考：https://github.com/ClickHouse/clickhouse-go/blob/v2.8.3/resources/meta.yml
*/

package dao

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/google/uuid"
	"github.com/spf13/cast"
)

type ClickhouseDB struct {
	User     string
	Password string
	Host     string
	Port     int
	Database string
	Settings map[string]interface{}
	Ctx      context.Context
}

func (c *ClickhouseDB) connect() (driver.Conn, error) {
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{fmt.Sprintf("%s:%d", c.Host, c.Port)},
		Auth: clickhouse.Auth{
			Database: c.Database,
			Username: c.User,
			Password: c.Password,
		},
		Compression: &clickhouse.Compression{
			Method: clickhouse.CompressionLZ4,
		},
		DialTimeout:  3 * time.Second,
		MaxOpenConns: 1,
		MaxIdleConns: 1,
		Settings:     c.Settings,
	})
	if err != nil {
		return nil, err
	}
	if err := conn.Ping(c.Ctx); err != nil {
		// if exception, ok := err.(*clickhouse.Exception); ok {
		// 	fmt.Printf("Exception [%d] %s \n%s\n", exception.Code, exception.Message, exception.StackTrace)
		// }
		return nil, err
	}
	return conn, nil
}

func (c *ClickhouseDB) Query(query string) (*[]string, *[]map[string]interface{}, error) {
	// 连接到db
	conn, err := c.connect()
	if err != nil {
		return nil, nil, err
	}
	defer conn.Close()
	// 执行查询
	rows, err := conn.Query(c.Ctx, query)
	if err != nil {
		return nil, nil, err
	}
	// 获取列名
	columns := rows.Columns()
	columnTypes := rows.ColumnTypes()

	vals := make([]interface{}, len(columns))
	for i := range columnTypes {
		vals[i] = reflect.New(columnTypes[i].ScanType()).Interface()
	}

	resultSlice := make([]map[string]interface{}, 0)
	for rows.Next() {
		if err := rows.Scan(vals...); err != nil {
			return nil, nil, err
		}
		vmap := make(map[string]interface{}, len(vals))
		for i, v := range vals {
			// 解决select 1 as id,2 as id,3 as id列名重复的问题，自动加别名
			if _, ok := vmap[columns[i]]; ok {
				columns[i] = fmt.Sprintf("%s_%s[别名]", columns[i], uuid.New())
			}
			var newV string
			// 类型断言
			switch v := v.(type) {
			case *string:
				newV = *v
			case *time.Time:
				newV = fmt.Sprintf("%v", v.Format("2006-01-02 15:04:05"))
			default:
				newV = cast.ToString(v)
			}
			vmap[columns[i]] = newV
		}
		resultSlice = append(resultSlice, vmap)
	}
	return &columns, &resultSlice, nil
}

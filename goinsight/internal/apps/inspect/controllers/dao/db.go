/*
@Time    :   2022/06/23 16:37:41
@Author  :   zongfei.fu
@Desc    :   操作目标审核数据库
*/

package dao

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/go-sql-driver/mysql"
)

// DB Struct
type DB struct {
	User     string
	Password string
	Host     string
	Port     int
	Database string
}

// Open connection for db
func (d *DB) Open() (*sql.DB, error) {
	config := mysql.Config{
		User:                 d.User,
		Passwd:               d.Password,
		Addr:                 fmt.Sprintf("%s:%d", d.Host, d.Port),
		Net:                  "tcp",
		DBName:               d.Database,
		AllowNativePasswords: true,
		Timeout:              3000 * time.Millisecond,
		ReadTimeout:          3000 * time.Millisecond,
		WriteTimeout:         3000 * time.Millisecond,
	}

	DSN := config.FormatDSN()
	db, err := sql.Open("mysql", DSN)
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(0)
	return db, err
}

// Executes a query without returning any rows.
func (d *DB) Execute(query string) error {
	db, err := d.Open()
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(query)
	return err
}

// Query
func (d *DB) Query(query string) (*[]map[string]interface{}, error) {
	db, err := d.Open()
	if err != nil {
		return nil, err
	}
	defer db.Close()
	// 执行查询
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	// 获取列名
	columns, error := rows.Columns()
	if error != nil {
		return nil, error
	}
	// Make a slice
	vals := make([]interface{}, len(columns))
	for i := range columns {
		vals[i] = new(sql.RawBytes)
	}
	// Fetch rows
	result := make([]map[string]interface{}, 0)
	for rows.Next() {
		if err := rows.Scan(vals...); err != nil {
			return nil, err
		}
		// var value string
		vmap := make(map[string]interface{}, len(vals))
		for i, c := range vals {
			// 类型断言
			switch v := c.(type) {
			case *sql.RawBytes:
				if *v == nil {
					// nil在前端解析的是null，符合预期，也可以直接return nil
					vmap[columns[i]] = "NULL"
				} else {
					vmap[columns[i]] = string(*v)
				}
			}
		}
		result = append(result, vmap)
	}
	if err = rows.Close(); err != nil {
		return nil, err
	}

	// Rows.Err will report the last error encountered by Rows.Scan.
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return &result, nil
}

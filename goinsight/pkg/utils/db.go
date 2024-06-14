package utils

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
)

// DB Struct
type DB struct {
	User     string
	Password string
	Host     string
	Port     int
	Database string
	Params   map[string]string
	Ctx      context.Context
}

// Open connection for db
func (d *DB) Open() (*sql.DB, error) {
	// user:1234.com@tcp(127.0.0.1:3306)/testdb
	config := mysql.Config{
		User:                 d.User,
		Passwd:               d.Password,
		Addr:                 fmt.Sprintf("%s:%d", d.Host, d.Port),
		Net:                  "tcp",
		DBName:               d.Database,
		AllowNativePasswords: true,
		Params:               d.Params,
		Timeout:              3 * time.Second, // Dial timeout
	}

	DSN := config.FormatDSN()
	db, err := sql.Open("mysql", DSN)
	db.SetMaxOpenConns(1)
	return db, err
}

// Query
func (d *DB) Query(query string) (*[]map[string]interface{}, error) {
	// 连接到db
	db, err := d.Open()
	if err != nil {
		return nil, err
	}
	defer db.Close()
	// 执行查询
	rows, err := db.QueryContext(d.Ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	// 获取列名
	columns, error := rows.Columns()
	if error != nil {
		return nil, err
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
			// 解决select 1 as id,2 as id,3 as id列名重复的问题，自动加别名
			if _, ok := vmap[columns[i]]; ok {
				columns[i] = fmt.Sprintf("%s_%s[别名]", columns[i], uuid.New())
			}
			// 类型断言
			switch v := c.(type) {
			case *sql.RawBytes:
				if *v == nil {
					// nil在前端解析的是null，符合预期，直接return nil
					vmap[columns[i]] = nil
				} else {
					vmap[columns[i]] = string(*v)
				}
			}
		}
		result = append(result, vmap)
	}
	// If the database is being written to ensure to check for Close
	// errors that may be returned from the driver. The query may
	// encounter an auto-commit error and be forced to rollback changes.
	if err = rows.Close(); err != nil {
		return nil, err
	}

	// Rows.Err will report the last error encountered by Rows.Scan.
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return &result, nil
}

// Executes a query without returning any rows.
func (d *DB) Execute(query string) error {
	db, err := d.Open()
	if err != nil {
		return err
	}
	defer db.Close()
	_, err = db.ExecContext(d.Ctx, query)
	return err
}

package api

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"goInsight/global"
	"goInsight/internal/pkg/utils"
	"strconv"
	"time"

	"github.com/go-sql-driver/mysql"
)

func NewMySQLCnx(cfg *DBConfig) (*sql.DB, error) {
	// Configure database connection
	config := mysql.Config{
		User:                 cfg.UserName,
		Passwd:               cfg.Password,
		Addr:                 fmt.Sprintf("%s:%d", cfg.Hostname, cfg.Port),
		Net:                  "tcp",
		DBName:               cfg.Schema,
		AllowNativePasswords: true,
		Params:               map[string]string{"sql_mode": "STRICT_TRANS_TABLES"},
		Timeout:              5 * time.Second, // Dial timeout
	}
	// Create a DSN string from the configuration
	DSN := config.FormatDSN()
	// Open a connection to the database
	db, err := sql.Open("mysql", DSN)
	if err != nil {
		return nil, err // Handle open error
	}
	// Limit the maximum number of open connections to 1
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(0)
	// Return the database connection and nil error
	return db, nil
}

// Query executes a query that returns rows, typically a SELECT. The args are for any placeholder parameters in the query
func DaoMySQLQuery(db *sql.DB, query string) (*[]map[string]interface{}, error) {
	// Execute the SQL query
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	// Get column names
	columns, error := rows.Columns()
	if error != nil {
		return nil, err
	}
	// Prepare data slices
	vals := make([]interface{}, len(columns))
	for i := range columns {
		vals[i] = new(sql.RawBytes)
	}
	// Collect results
	result := make([]map[string]interface{}, 0)
	for rows.Next() {
		err := rows.Scan(vals...) // Scan each row into slices
		if err != nil {
			return nil, err
		}
		// Map for each row
		vmap := make(map[string]interface{}, len(vals))
		for i, c := range vals {
			// Type assertion and value conversion
			switch v := c.(type) {
			case *sql.RawBytes:
				if *v == nil {
					// Handle null value as nil
					vmap[columns[i]] = nil
				} else {
					// Convert RawBytes to string
					vmap[columns[i]] = string(*v)
				}
			}
		}
		// Append each row map to the results
		result = append(result, vmap)
	}
	// Check for errors after closing the rows
	if err = rows.Close(); err != nil {
		return nil, err
	}
	// Check for any additional errors
	if err = rows.Err(); err != nil {
		return nil, err
	}
	// Return pointer to slice of result maps
	return &result, nil
}

// Exec executes a query without returning any rows. The args are for any placeholder parameters in the query
func DaoMySQLExecute(db *sql.DB, sql string, ch chan<- int64) (int64, error) {
	// Send a signal to the channel to indicate that the function has started
	ch <- 1
	// Execute the SQL statement
	result, err := db.Exec(sql)
	if err != nil {
		// Return an error if the statement failed
		return 0, err
	}
	// Close the channel to indicate that the function has completed
	close(ch)
	// Return the number of rows affected by the statement
	return result.RowsAffected()
}

// 获取MySQL类数据库的连接ID
func DaoGetConnectionID(db *sql.DB) (connectionID int64, err error) {
	// Define the SQL query to retrieve the connection ID
	sql := "SELECT CONNECTION_ID() as CONNECTION_ID"
	// Execute the query and capture any errors
	data, err := DaoMySQLQuery(db, sql)
	if err != nil {
		return connectionID, err
	}
	// Check if no data was retrieved
	if len(*data) == 0 {
		return connectionID, errors.New("Failed to get connection ID: no valid row found")
	}
	// Expect only one row to be returned
	row := (*data)[0]
	// Extract the "CONNECTION_ID" value from the first row
	connectionIDStr := row["CONNECTION_ID"].(string)
	// Attempt to parse the string value as a 64-bit integer
	connectionID, err = strconv.ParseInt(connectionIDStr, 10, 64)
	if err != nil {
		return connectionID, errors.New("Failed to get connection ID: parsing error")
	}
	// Return the retrieved connection ID and nil error
	return
}

// 获取MySQL类数据库的Processlist
func DaoGetProcesslist(dbconfig *DBConfig, order_id string, connection_id int64, ch <-chan int64) {
	// Create a new database connection
	db, err := NewMySQLCnx(dbconfig)
	if err != nil {
		global.App.Log.Error(err.Error())
		return
	}
	defer db.Close() // Close the database connection when the function exits
	// Construct the SQL query
	sql := fmt.Sprintf("SELECT * FROM INFORMATION_SCHEMA.PROCESSLIST WHERE ID=%d", connection_id)
	// Loop
	for {
		exitFlag := false
		select {
		case _, ok := <-ch: // Receive a signal from the channel
			if !ok { // If the channel is closed
				exitFlag = true // Set the exit flag to true
			}
		case <-time.After(500 * time.Millisecond): // Wait for 500 milliseconds
			// This allows the function to check for new data periodically
		}
		if exitFlag { // If the exit flag is set
			break // Exit the loop
		}
		// Execute the SQL query
		data, err := DaoMySQLQuery(db, sql)
		if err != nil {
			global.App.Log.Error("Failed to get processlist: ", err.Error())
			break // Exit the loop
		}
		if len(*data) == 0 { // If no rows are returned
			global.App.Log.Error("Failed to get processlist: no valid row found")
			break // Exit the loop
		}
		// Expect only one row of data
		row := (*data)[0]
		// Publish the process information
		PublishMsg(order_id, row, "processlist")
	}
}

// 获取MySQL数据库的Position
func DaoGetMySQLPos(db *sql.DB) (file string, position int64, err error) {
	// Query the database for the master status
	data, err := DaoMySQLQuery(db, "SHOW MASTER STATUS")
	if err != nil {
		return file, position, err
	}
	// Check if data is empty
	if len(*data) == 0 {
		return file, position, errors.New("Failed to get MySQL position: no valid row found，请检查MySQL是否开启了binlog")
	}
	// Expect to return one row of data
	row := (*data)[0]
	// Parse the file and position from the row data
	file = row["File"].(string)
	position, err = strconv.ParseInt(row["Position"].(string), 10, 64)
	if err != nil {
		return file, position, errors.New("Failed to get MySQL position: position parsing error")
	}
	return file, position, nil
}

func PublishMsg(channel string, data interface{}, renderType string) {
	// 发送消息
	err := utils.Publish(context.Background(), channel, data, renderType)
	if err != nil {
		global.App.Log.Error(err)
	}
}

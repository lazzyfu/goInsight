/*
@Author  :   xff
@Desc    :   生成MySQL回滚语句
*/

package api

import (
	"context"
	"database/sql/driver"
	"fmt"
	"goInsight/pkg/parser"
	"goInsight/pkg/utils"
	"strings"
	"time"

	"github.com/pingcap/tidb/pkg/parser/ast"

	"github.com/go-mysql-org/go-mysql/mysql"
	"github.com/go-mysql-org/go-mysql/replication"
)

type Binlog struct {
	*DBConfig
	ConnectionID  int64
	StartFile     string
	StartPosition int64
	EndFile       string
	EndPosition   int64
}

func (b *Binlog) parserTableStmt(table string) (*ast.CreateTableStmt, error) {
	// 初始化DB
	db, err := NewMySQLCnx(b.DBConfig)
	if err != nil {
		return nil, err
	}
	defer db.Close()
	// 查看表结构
	data, err := DaoMySQLQuery(db, fmt.Sprintf("show create table %s", table))
	if err != nil {
		return nil, err
	}
	var tableStruct string
	for _, row := range *data {
		// Expect to return one row of data
		tableStruct = row["Create Table"].(string)
		break
	}
	// 解析表结构
	stmt, err := parser.NewParseOneStmt(tableStruct, "", "")
	if err != nil {
		return nil, err
	}
	switch s := stmt.(type) {
	case *ast.CreateTableStmt:
		return s, nil
	}
	return nil, nil
}

func (b *Binlog) extractPK(stmt *ast.CreateTableStmt) (bool, []string) {
	var keys []string
	// 从I_ID bigint unsigned NOT NULL AUTO_INCREMENT primary key COMMENT '自增ID' 提取主键
	for _, col := range stmt.Cols {
		for _, opt := range col.Options {
			if opt.Tp == ast.ColumnOptionPrimaryKey {
				keys = append(keys, col.Name.Name.O)
			}
		}
	}
	// 从PRIMARY KEY (I_ID) 提取主键
	for _, cons := range stmt.Constraints {
		if cons.Tp == ast.ConstraintPrimaryKey {
			for _, col := range cons.Keys {
				if !utils.IsContain(keys, col.Column.Name.O) {
					keys = append(keys, col.Column.Name.O)
				}
			}
		}
	}
	return len(keys) > 0, keys
}

func (b *Binlog) Run() (string, error) {
	cfg := replication.BinlogSyncerConfig{
		ServerID:   20231108 + uint32(uint32(time.Now().Unix())%10000),
		Flavor:     "mysql",
		Host:       b.Hostname,
		Port:       b.Port,
		User:       b.UserName,
		Password:   b.Password,
		UseDecimal: true,
	}
	syncer := replication.NewBinlogSyncer(cfg)
	defer syncer.Close()
	// 定义开始结束的pos
	startPosition := mysql.Position{Name: b.StartFile, Pos: uint32(b.StartPosition)}
	stopPosition := mysql.Position{Name: b.EndFile, Pos: uint32(b.EndPosition)}
	// 开启同步
	streamer, err := syncer.StartSync(startPosition)
	if err != nil {
		return "", err
	}
	// 声明当前的pos
	currentPosition := startPosition
	// 获取当前事件的thread id，用来和执行SQL的thread id进行比较
	var currentThreadID uint32
	// 回滚SQL
	var rbsqls []string
	// 循环
	for {
		e, err := streamer.GetEvent(context.Background())
		if err != nil {
			return "", err
		}

		if e.Header.LogPos > 0 {
			currentPosition.Pos = e.Header.LogPos
		}

		if e.Header.EventType == replication.ROTATE_EVENT {
			if event, ok := e.Event.(*replication.RotateEvent); ok {
				currentPosition = mysql.Position{Name: string(event.NextLogName),
					Pos: uint32(event.Position)}
			}
		}

		if currentPosition.Compare(startPosition) == -1 {
			continue
		}
		// 如果当前pos大于停止的pos，退出
		if currentPosition.Compare(stopPosition) > -1 {
			break
		}
		// 事件类型判断
		switch e.Header.EventType {
		case replication.QUERY_EVENT:
			if event, ok := e.Event.(*replication.QueryEvent); ok {
				currentThreadID = event.SlaveProxyID
			}
		case replication.WRITE_ROWS_EVENTv1, replication.WRITE_ROWS_EVENTv2:
			if event, ok := e.Event.(*replication.RowsEvent); ok {
				// 获取表的stmt
				tableName := fmt.Sprintf("`%s`.`%s`", event.Table.Schema, event.Table.Table)
				stmt, err := b.parserTableStmt(tableName)
				if err != nil {
					return "", err
				}
				// 解析回滚语句
				if b.ConnectionID == int64(currentThreadID) {
					sql, err := b.generateDeleteSql(event, stmt)
					if err != nil {
						return "", err
					}
					rbsqls = append(rbsqls, sql)
				}
			}
		case replication.DELETE_ROWS_EVENTv1, replication.DELETE_ROWS_EVENTv2:
			if event, ok := e.Event.(*replication.RowsEvent); ok {
				// 获取表的stmt
				tableName := fmt.Sprintf("`%s`.`%s`", event.Table.Schema, event.Table.Table)
				stmt, err := b.parserTableStmt(tableName)
				if err != nil {
					return "", err
				}
				if b.ConnectionID == int64(currentThreadID) {
					sql, err := b.generateInsertSql(event, stmt)
					if err != nil {
						return "", err
					}
					rbsqls = append(rbsqls, sql)
				}
			}
		case replication.UPDATE_ROWS_EVENTv1, replication.UPDATE_ROWS_EVENTv2:
			if event, ok := e.Event.(*replication.RowsEvent); ok {
				// 获取表的stmt
				tableName := fmt.Sprintf("`%s`.`%s`", event.Table.Schema, event.Table.Table)
				stmt, err := b.parserTableStmt(tableName)
				if err != nil {
					return "", err
				}
				if b.ConnectionID == int64(currentThreadID) {
					sql, err := b.generateUpdateSql(event, stmt)
					if err != nil {
						return "", err
					}
					rbsqls = append(rbsqls, sql)
				}
			}
		}
	}
	return strings.Join(rbsqls, ";\r\n"), nil
}

func (b *Binlog) generateUpdateSql(e *replication.RowsEvent, stmt *ast.CreateTableStmt) (string, error) {
	template := "UPDATE `%s`.`%s` SET %s WHERE"
	hasPrimaryKey, PrimaryKeys := b.extractPK(stmt)

	var (
		oldValues []driver.Value
		newValues []driver.Value
		newSql    string
	)

	var rbsqls []string
	var sets []string
	var sql string

	// e.Rows:  [[2 dasdas6 MySQL 2] [2 dasdas7 MySQL 2] [3 dasdas6 MySQL 3] [3 dasdas7 MySQL 3]]
	for i, rows := range e.Rows {
		// rows：[2 dasdas6 MySQL 2]
		var columns []string
		if i%2 == 0 {
			// old values
			for i, d := range rows {
				col := stmt.Cols[i]
				// 跳过计算列
				if isGenerated(col.Options) {
					continue
				}
				if isUnsigned(stmt.Cols[i].Tp.GetFlag()) {
					d = processValue(d, stmt.Cols[i].Tp.GetType())
				}
				sets = append(sets, fmt.Sprintf(" `%s`=?", col.Name.Name.O))
				newValues = append(newValues, d)
			}
			sql = fmt.Sprintf(template, e.Table.Schema, e.Table.Table, strings.Join(sets, ","))
		} else {
			// new values
			for i, d := range rows {
				col := stmt.Cols[i]
				// 跳过计算列
				if isGenerated(col.Options) {
					continue
				}
				if hasPrimaryKey {
					if utils.IsContain(PrimaryKeys, col.Name.Name.O) {
						if isUnsigned(stmt.Cols[i].Tp.GetFlag()) {
							d = processValue(d, stmt.Cols[i].Tp.GetType())
						}
						oldValues = append(oldValues, d)
						if d == nil {
							columns = append(columns,
								fmt.Sprintf(" `%s` IS ?", col.Name.Name.O))
						} else {
							columns = append(columns,
								fmt.Sprintf(" `%s`=?", col.Name.Name.O))
						}
					}
				} else {
					if isUnsigned(stmt.Cols[i].Tp.GetFlag()) {
						d = processValue(d, stmt.Cols[i].Tp.GetType())
					}
					oldValues = append(oldValues, d)
					if d == nil {
						columns = append(columns,
							fmt.Sprintf(" `%s` IS ?", col.Name.Name.O))
					} else {
						columns = append(columns,
							fmt.Sprintf(" `%s`=?", col.Name.Name.O))
					}
				}
				// 重置
				sets = []string{}
			}
			newSql = strings.Join([]string{sql, strings.Join(columns, " AND")}, "")
			newValues = append(newValues, oldValues...)
			r, err := interpolateParams(newSql, newValues, true)
			if err != nil {
				return "", err
			}
			rbsqls = append(rbsqls, string(r))
			oldValues = nil
			newValues = nil
		}
	}
	return strings.Join(rbsqls, ";\r\n"), nil
}

func (b *Binlog) generateInsertSql(e *replication.RowsEvent, stmt *ast.CreateTableStmt) (string, error) {
	var columns []string
	template := "INSERT INTO `%s`.`%s`(%s) VALUES(%s)"
	for _, col := range stmt.Cols {
		// 跳过计算列
		if isGenerated(col.Options) {
			continue
		}
		columns = append(columns, fmt.Sprintf("`%s`", col.Name.Name.O))
	}
	paramValues := strings.TrimRight(strings.Repeat("?,", len(columns)), ",")
	sql := fmt.Sprintf(template, e.Table.Schema, e.Table.Table,
		strings.Join(columns, ","), paramValues)

	var rbsqls []string
	for _, rows := range e.Rows {
		var vv []driver.Value
		for i, d := range rows {
			col := stmt.Cols[i]
			// 跳过计算列
			if isGenerated(col.Options) {
				continue
			}
			if isUnsigned(col.Tp.GetFlag()) {
				d = processValue(d, col.Tp.GetType())
			}
			vv = append(vv, d)
		}
		r, err := interpolateParams(sql, vv, false)
		if err != nil {
			return "", err
		}
		rbsqls = append(rbsqls, string(r))
	}
	return strings.Join(rbsqls, ";\r\n"), nil
}

func (b *Binlog) generateDeleteSql(e *replication.RowsEvent, stmt *ast.CreateTableStmt) (string, error) {
	template := "DELETE FROM `%s`.`%s` WHERE "
	sql := fmt.Sprintf(template, e.Table.Schema, e.Table.Table)
	hasPrimaryKey, PrimaryKeys := b.extractPK(stmt)

	var rbsqls []string
	// e.Rows为insert into xx values(),(),();
	for _, rows := range e.Rows {
		var vv []driver.Value
		var columns []string
		// 判断是否有主键并提取主键
		if hasPrimaryKey {
			for i, d := range rows {
				col := stmt.Cols[i]
				// 跳过计算列
				if isGenerated(col.Options) {
					continue
				}
				if utils.IsContain(PrimaryKeys, col.Name.Name.O) {
					vv = append(vv, d)
					columns = append(columns, fmt.Sprintf("`%s`=?", col.Name.Name.O))
				}
			}
		} else {
			for i, d := range rows {
				col := stmt.Cols[i]
				// 跳过计算列
				if isGenerated(col.Options) {
					continue
				}
				vv = append(vv, d)
				if d == nil {
					columns = append(columns, fmt.Sprintf("`%s` IS ?", col.Name.Name.O))
				} else {
					columns = append(columns, fmt.Sprintf("`%s`=?", col.Name.Name.O))
				}
			}
		}

		newSql := strings.Join([]string{sql, strings.Join(columns, " AND ")}, "")
		r, err := interpolateParams(newSql, vv, false)
		if err != nil {
			return "", err
		}
		rbsqls = append(rbsqls, string(r))
	}
	return strings.Join(rbsqls, ";\r\n"), nil
}

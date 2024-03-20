package api

import "fmt"

// SQL执行失败
type SQLExecuteError struct {
	Err error
}

func (e SQLExecuteError) Error() string {
	return fmt.Sprintf("SQL execution failed，Error: %v", e.Err)
}

// 生成回滚SQL失败
type RollbackSQLError struct {
	Err error
}

func (e RollbackSQLError) Error() string {
	return fmt.Sprintf("Failed to generate rollback SQL, Error: %v", e.Err)
}

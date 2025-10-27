package base

// SQL执行失败
type SQLExecuteError struct {
	Err error
}

func (e SQLExecuteError) Error() string {
	return e.Err.Error()
}

// 生成回滚SQL失败
type RollbackSQLError struct {
	Err error
}

func (e RollbackSQLError) Error() string {
	return e.Err.Error()
}

package base

import (
	"fmt"
	"strings"
	"time"
)

type ExecuteLogger struct {
	logs []string
}

func NewExecuteLogger() *ExecuteLogger {
	return &ExecuteLogger{
		logs: make([]string, 0, 16),
	}
}

func (l *ExecuteLogger) Add(msg string) string {
	formatted := fmt.Sprintf(
		"[%s] %s",
		time.Now().Format("2006-01-02 15:04:05"),
		msg,
	)
	l.logs = append(l.logs, formatted)
	return formatted
}

func (l *ExecuteLogger) String() string {
	return strings.Join(l.logs, "\n")
}

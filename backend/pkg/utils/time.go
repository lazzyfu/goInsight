package utils

import (
	"time"

	"github.com/dustin/go-humanize"
)

// 时间格式化
func HumanfriendlyTimeUnit(d time.Duration) string {
	relTime := humanize.RelTimeMagnitude{D: d, Format: "now", DivBy: time.Second}
	return relTime.D.String()
}

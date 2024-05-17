/*
@Time    :   2022/07/06 10:12:48
@Author  :   xff
@Desc    :   None
*/

package process

import (
	"fmt"
	"strconv"
	"strings"
)

// 版本格式：major.minor.patch
type DbVersion struct {
	Version string
}

func (d *DbVersion) Format() []string {
	// 5.7.25-TiDB-v5.0.4
	// 5.7.35-log
	// 8.0.26
	//
	tokens := strings.Split(d.Version, "-")
	if len(tokens) == 0 {
		return []string{"1", "0", "00"}
	}
	versionSeg := strings.Split(tokens[0], ".")
	if len(versionSeg) != 3 {
		return []string{"1", "0", "00"}
	}
	return versionSeg
}

func (d *DbVersion) Int() int {
	// 5.7.35-log 返回 5735
	tokens := d.Format()
	versionStr := fmt.Sprintf("%s%02s%02s", tokens[0], tokens[1], tokens[2])
	version, err := strconv.Atoi(versionStr)
	if err != nil {
		return 10000
	}
	return version
}

func (d *DbVersion) IsTiDB() bool {
	// 判断是否为tidb
	return strings.Contains(strings.ToLower(d.Version), "tidb")
}

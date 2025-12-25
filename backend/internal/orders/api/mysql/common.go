package mysql

import (
	"database/sql/driver"
	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/pingcap/tidb/pkg/parser/ast"
	"github.com/pingcap/tidb/pkg/parser/mysql"
	"github.com/shopspring/decimal"
)

// 判断是否为计算列
func isGenerated(opts []*ast.ColumnOption) bool {
	for _, co := range opts {
		// 跳过计算列
		if co.Tp == ast.ColumnOptionGenerated {
			return true
		}
	}
	return false
}

func isUnsigned(flag uint) bool {
	return mysql.HasUnsignedFlag(flag)
}

func interpolateParams(query string, args []driver.Value, hexBlob bool) ([]byte, error) {
	// Number of ? should be same to len(args)
	if strings.Count(query, "?") != len(args) {
		return nil, fmt.Errorf("生成回滚sql的参数个数不匹配：查询：%s，需要参数数量：%d，提供参数：%d", query, strings.Count(query, "?"), len(args))
	}

	var buf []byte

	argPos := 0

	for i := 0; i < len(query); i++ {
		q := strings.IndexByte(query[i:], '?')
		if q == -1 {
			buf = append(buf, query[i:]...)
			break
		}
		buf = append(buf, query[i:i+q]...)
		i += q

		arg := args[argPos]
		argPos++

		if arg == nil {
			buf = append(buf, "NULL"...)
			continue
		}

		// log.Info(arg)
		// log.Infof("%T", arg)

		switch v := arg.(type) {
		case int8:
			buf = strconv.AppendInt(buf, int64(v), 10)
		case int16:
			buf = strconv.AppendInt(buf, int64(v), 10)
		case int32:
			buf = strconv.AppendInt(buf, int64(v), 10)
		case int64:
			buf = strconv.AppendInt(buf, v, 10)
		case uint64:
			buf = strconv.AppendUint(buf, uint64(v), 10)
		case int:
			buf = strconv.AppendInt(buf, int64(v), 10)
		case decimal.Decimal:
			buf = append(buf, v.String()...)
		case float32:
			buf = strconv.AppendFloat(buf, float64(v), 'g', -1, 32)
		case float64:
			buf = strconv.AppendFloat(buf, v, 'g', -1, 64)
		case bool:
			if v {
				buf = append(buf, '1')
			} else {
				buf = append(buf, '0')
			}
		case time.Time:
			if v.IsZero() {
				buf = append(buf, "'0000-00-00'"...)
			} else {
				v := v.In(time.UTC)
				v = v.Add(time.Nanosecond * 500) // To round under microsecond
				year := v.Year()
				year100 := year / 100
				year1 := year % 100
				month := v.Month()
				day := v.Day()
				hour := v.Hour()
				minute := v.Minute()
				second := v.Second()
				micro := v.Nanosecond() / 1000

				buf = append(buf, []byte{
					'\'',
					digits10[year100], digits01[year100],
					digits10[year1], digits01[year1],
					'-',
					digits10[month], digits01[month],
					'-',
					digits10[day], digits01[day],
					' ',
					digits10[hour], digits01[hour],
					':',
					digits10[minute], digits01[minute],
					':',
					digits10[second], digits01[second],
				}...)

				if micro != 0 {
					micro10000 := micro / 10000
					micro100 := micro / 100 % 100
					micro1 := micro % 100
					buf = append(buf, []byte{
						'.',
						digits10[micro10000], digits01[micro10000],
						digits10[micro100], digits01[micro100],
						digits10[micro1], digits01[micro1],
					}...)
				}
				buf = append(buf, '\'')
			}
		case string:
			if hexBlob {
				if utf8.ValidString(v) {
					buf = append(buf, '\'')
					buf = escapeBytesBackslash(buf, []byte(v))
				} else {
					buf = append(buf, 'X')
					buf = append(buf, '\'')
					b := hex.EncodeToString([]byte(v))
					buf = append(buf, b...)
				}
			} else {
				buf = append(buf, '\'')
				buf = escapeBytesBackslash(buf, []byte(v))
			}

			buf = append(buf, '\'')
		case []byte:
			if v == nil {
				buf = append(buf, "NULL"...)
			} else {
				// buf = append(buf, "_binary'"...)
				if hexBlob {
					if utf8.Valid(v) {
						buf = append(buf, '\'')
						buf = escapeBytesBackslash(buf, v)
					} else {
						buf = append(buf, 'X')
						buf = append(buf, '\'')
						b := hex.EncodeToString(v)
						buf = append(buf, b...)
					}
				} else {
					buf = append(buf, '\'')
					buf = escapeBytesBackslash(buf, v)
				}
				buf = append(buf, '\'')
			}
		default:
			return nil, errors.New("driver: skip fast-path; continue as if unimplemented")
		}

		// 4 << 20 , 4MB
		if len(buf)+4 > 4<<20 {
			return nil, errors.New("driver: skip fast-path; continue as if unimplemented")
		}
	}
	if argPos != len(args) {
		return nil, errors.New("driver: skip fast-path; continue as if unimplemented")
	}
	return buf, nil
}

const digits01 = "0123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789"
const digits10 = "0000000000111111111122222222223333333333444444444455555555556666666666777777777788888888889999999999"

func escapeBytesBackslash(buf, v []byte) []byte {
	pos := len(buf)
	buf = reserveBuffer(buf, len(v)*2)

	for _, c := range v {
		switch c {
		case '\x00':
			buf[pos] = '\\'
			buf[pos+1] = '0'
			pos += 2
		case '\n':
			buf[pos] = '\\'
			buf[pos+1] = 'n'
			pos += 2
		case '\r':
			buf[pos] = '\\'
			buf[pos+1] = 'r'
			pos += 2
		case '\x1a':
			buf[pos] = '\\'
			buf[pos+1] = 'Z'
			pos += 2
		case '\'':
			buf[pos] = '\\'
			buf[pos+1] = '\''
			pos += 2
		case '"':
			buf[pos] = '\\'
			buf[pos+1] = '"'
			pos += 2
		case '\\':
			buf[pos] = '\\'
			buf[pos+1] = '\\'
			pos += 2
		default:
			buf[pos] = c
			pos++
		}
	}

	return buf[:pos]
}

func reserveBuffer(buf []byte, appendSize int) []byte {
	newSize := len(buf) + appendSize
	if cap(buf) < newSize {
		// Grow buffer exponentially
		newBuf := make([]byte, len(buf)*2+appendSize)
		copy(newBuf, buf)
		buf = newBuf
	}
	return buf[:newSize]
}

// processValue 处理无符号值(unsigned)
func processValue(value driver.Value, tp byte) driver.Value {
	if value == nil {
		return value
	}

	switch v := value.(type) {
	case int8:
		if v >= 0 {
			return value
		}
		return int64(1<<8 + int64(v))
	case int16:
		if v >= 0 {
			return value
		}
		return int64(1<<16 + int64(v))
	case int32:
		if v >= 0 {
			return value
		}
		if tp == mysql.TypeInt24 {
			return int64(1<<24 + int64(v))
		}
		return int64(1<<32 + int64(v))
	case int64:
		if v >= 0 {
			return value
		}
		return math.MaxUint64 - uint64(abs(v)) + 1
	// case int:
	// case float32:
	// case float64:

	default:
		// log.Error("解析错误")
		// log.Errorf("%T", v)
		return value
	}
}

func abs(n int64) int64 {
	y := n >> 63
	return (n ^ y) - y
}

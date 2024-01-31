package process

import (
	"strings"

	"github.com/pingcap/tidb/parser/mysql"
)

var charSets = map[string]int{
	"armscii8": 1,
	"ascii":    1,
	"big5":     2,
	"binary":   1,
	"cp1250":   1,
	"cp1251":   1,
	"cp1256":   1,
	"cp1257":   1,
	"cp850":    1,
	"cp852":    1,
	"cp866":    1,
	"cp932":    2,
	"dec8":     1,
	"eucjpms":  3,
	"euckr":    2,
	"gb18030":  4,
	"gb2312":   2,
	"gbk":      2,
	"geostd8":  1,
	"greek":    1,
	"hebrew":   1,
	"hp8":      1,
	"keybcs2":  1,
	"koi8r":    1,
	"koi8u":    1,
	"latin1":   1,
	"latin2":   1,
	"latin5":   1,
	"latin7":   1,
	"macce":    1,
	"macroman": 1,
	"sjis":     2,
	"swe7":     1,
	"tis620":   1,
	"ucs2":     2,
	"ujis":     3,
	"utf16":    4,
	"utf16le":  4,
	"utf32":    4,
	"utf8":     3,
	"utf8mb4":  4,
}

// getDataBytes 计算数据类型字节数
// 计算方法参考:https://dev.mysql.com/doc/refman/8.0/en/storage-requirements.html
type DataBytes struct {
	Column  string
	Tp      byte
	Elems   []string // Elems is the element list for enum and set type.
	Ilen    int      // key `idx_name`(name(32))中的32
	Flen    int      // 字段长度
	Decimal int      // decimal字段专用,decimal(12,2)中的2
	Charset string   // 列字符集
}

func (l *DataBytes) Get(dbVersion int) int {
	// Flen: 字段定义的长度
	// Ilen: key `idx_name`(name(32))中的32
	// DLen: 创建索引时指定字段的实际长度，如果Ilen有值，长度为Ilen，否则为Flen; 如果类型为enum或set时，长度为len(Elems)

	var Dlen int
	switch l.Tp {
	case mysql.TypeEnum, mysql.TypeSet:
		Dlen = len(l.Elems)
	default:
		Dlen = l.Flen
		if l.Ilen > 0 {
			Dlen = l.Ilen
		}
	}

	switch l.Tp {
	case mysql.TypeString, mysql.TypeVarchar, mysql.TypeEnum, mysql.TypeSet:
		return StringStorageReq(l.Tp, l.Charset, Dlen)
	case mysql.TypeTiny, mysql.TypeShort, mysql.TypeInt24, mysql.TypeLonglong, mysql.TypeLong, mysql.TypeFloat, mysql.TypeDouble, mysql.TypeNewDecimal, mysql.TypeBit:
		return numericStorageReq(l.Tp, l.Flen, Dlen, l.Decimal)
	case mysql.TypeDate, mysql.TypeDuration, mysql.TypeDatetime, mysql.TypeYear, mysql.TypeNewDate, mysql.TypeTimestamp:
		return timeStorageReq(l.Tp, Dlen, dbVersion)
	default:
		return -1
	}
}

// 字符串类字节长度计算
func StringStorageReq(dataType byte, charset string, Dlen int) int {
	// get bytes per character, default 1
	bytesPerChar := 1
	if _, ok := charSets[strings.ToLower(charset)]; ok {
		// 根据字符集返回字符对应的字节长度
		bytesPerChar = charSets[strings.ToLower(charset)]
	}

	switch dataType {
	case mysql.TypeString:
		// char or binary
		if Dlen > 255 {
			Dlen = 255
		}
		if strings.ToLower(charset) == "binary" {
			// M bytes, 0 <= M <= 255
			return Dlen
		}
		return Dlen * bytesPerChar
	case mysql.TypeVarchar:
		// varchar or varbinary
		// L + 1 bytes if column values require 0 − 255 bytes, L + 2 bytes if values may require more than 255 bytes
		if Dlen*bytesPerChar <= 255 {
			return Dlen*bytesPerChar + 1
		}
		return Dlen*bytesPerChar + 2
	case mysql.TypeEnum:
		// 1 or 2 bytes, depending on the number of enumeration values (65,535 values maximum)
		return Dlen/(2^15) + 1
	case mysql.TypeSet:
		// 1, 2, 3, 4, or 8 bytes, depending on the number of set members (64 members maximum)
		return Dlen/8 + 1
	default:
		return 0
	}
}

// 数值类字节长度计算
func numericStorageReq(dataType byte, Flen int, Dlen int, Decimal int) int {
	switch dataType {
	case mysql.TypeTiny:
		// tinyint
		return 1
	case mysql.TypeShort:
		// smallint
		return 2
	case mysql.TypeInt24:
		// mediumint
		return 3
	case mysql.TypeLong:
		// int
		return 4
	case mysql.TypeLonglong, mysql.TypeDouble:
		// bigint
		return 8
	case mysql.TypeFloat:
		return 4
	case mysql.TypeNewDecimal:
		// 计算decimal
		if Flen == -1 {
			// decimal() 直接返回
			return 4
		}
		leftover := func(leftover int) int {
			if leftover > 0 && leftover <= 2 {
				return 1
			} else if leftover > 2 && leftover <= 4 {
				return 2
			} else if leftover > 4 && leftover <= 6 {
				return 3
			} else if leftover > 6 && leftover <= 8 {
				return 4
			} else {
				return 4
			}
		}
		integer := Flen/9*4 + leftover(Flen%9)
		fractional := Decimal/9*4 + leftover(Decimal%9)
		return integer + fractional
	case mysql.TypeBit:
		// approximately (M+7)/8 bytes
		if Flen == -1 {
			return 1
		}
		return (Dlen + 7) / 8
	default:
		return -1
	}
}

// 时间类型长度计算
func timeStorageReq(dataType byte, Dlen int, dbVersion int) int {
	extr := func(length int) int {
		if length > 0 && length <= 2 {
			return 1
		} else if length > 2 && length <= 4 {
			return 2
		} else if length > 4 && length <= 6 || length > 6 {
			return 3
		}
		return 0
	}
	switch dataType {
	case mysql.TypeYear:
		// year
		return 1
	case mysql.TypeDate:
		// date
		return 3
	case mysql.TypeDuration:
		// time
		if dbVersion < 50604 {
			return 3
		}
		// 3 bytes + fractional seconds storage
		return 3 + extr(Dlen)
	case mysql.TypeDatetime:
		// time
		if dbVersion < 50604 {
			return 8
		}
		// 5 bytes + fractional seconds storage
		return 5 + extr(Dlen)
	case mysql.TypeTimestamp:
		// timestamp
		if dbVersion < 50604 {
			return 4
		}
		// 4 bytes + fractional seconds storage
		return 4 + extr(Dlen)
	default:
		return 8
	}
}

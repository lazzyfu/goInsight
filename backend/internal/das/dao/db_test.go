package dao

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 定义测试用例的表
var testCaseTables map[string]string = map[string]string{
	"testcase_bit": `CREATE TABLE testcase_bit (
						col_bit bit(2)
					) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4`,
	"testcase_int": `CREATE TABLE IF NOT EXISTS testcase_int (
						col_tinyint tinyint,
						col_tinyint_unsigned tinyint unsigned not null,
						col_smallint smallint,
						col_smallint_unsigned smallint unsigned not null,
						col_mediumint mediumint,
						col_mediumint_unsigned mediumint unsigned not null,
						col_int int,
						col_int_unsigned int unsigned not null,
						col_bigint bigint,
						col_bigint_unsigned bigint unsigned not null
					) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4`,
	"testcase_decimal": `CREATE TABLE IF NOT EXISTS testcase_decimal (
							col_decimal decimal(10,2),
							col_float float,
							col_double double
						) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4`,
	"testcase_date": `CREATE TABLE IF NOT EXISTS testcase_date (
						col_year year,
						col_date date,
						col_time time,
						col_datetime datetime,
						col_timestamp timestamp
						) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4`,
	"testcase_char": `CREATE TABLE IF NOT EXISTS testcase_char (
						col_char char(30),
						col_char_not_null char(30) not null default '',
						col_varchar varchar(30),
						col_varchar_not_null varchar(30) not null default ''
					) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4`,
	"testcase_byte": `CREATE TABLE IF NOT EXISTS testcase_byte (
						col_binary binary(3),
						col_varbinary varbinary(3),
						col_tinyblob tinyblob,
						col_blob blob,
						col_mediumblob mediumblob,
						col_longblob longblob
					) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4`,
	"testcase_text": `CREATE TABLE IF NOT EXISTS testcase_text (
						col_tinytext tinytext,
						col_mediumtext mediumtext,
						col_text text,
						col_longtext longtext
					) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4`,
	"testcase_enum_set": `CREATE TABLE IF NOT EXISTS testcase_enum_set (
							col_enum enum('x-small', 'small', 'medium', 'large', 'x-large'),
							col_set set('a', 'b', 'c', 'd')
						) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4`,
	"testcase_spatial": `CREATE TABLE IF NOT EXISTS testcase_spatial (
							col_geometry geometry
						) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4`,
	"testcase_json": `CREATE TABLE IF NOT EXISTS testcase_json (
						col_json json
					) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4`,
}

// 定义测试用例的数据，每个表的用例数据一组
var testCaseData [][]string = [][]string{
	{
		`insert into testcase_bit values(b'10')`,
		`insert into testcase_bit values(null)`,
		`insert into testcase_bit values('')`,
	},
	{
		`insert into testcase_int values (127, 255, 32767, 65535, 8388607, 16777215, 2147483647, 4294967295, 9223372036854775807, 18446744073709551615)`,
		`insert into testcase_int values (-12, 255, -3276, 65535, -838860, 16777215, -214748364, 4294967295, -922337203685477580, 18446744073709551615)`,
		`insert into testcase_int values (null, 0, null, 0, null, 0, null, 0, null, 0)`,
	},
	{
		`insert into testcase_decimal values (100.22, 100.00, 100)`,
		`insert into testcase_decimal values (-10.22, -100.00, -100)`,
		`insert into testcase_decimal values (null, null, null)`,
	},
	{
		`insert into testcase_date values ("2023", "2023-06-06", "12:00:00", "2023-06-06 12:00:00", "2023-06-06 12:00:00")`,
		`insert into testcase_date values ("0000", "0000-00-00", "00:00:00", "2023-06-06 12:00:00", "2023-06-06 12:00:00")`,
		`insert into testcase_date values (null, null, null, null, "2023-06-06 12:00:00")`,
		`insert into testcase_date values (0, 0, 0, 0, 0)`,
	},
	{
		`insert into testcase_char values ("char", "char_not_null", "varchar", "varchar_not_null")`,
		`insert into testcase_char values (null, "", null, "")`,
	},
	{
		`insert into testcase_byte values ("a", "ab", "tinyblob", "blob", "mediumblob", "longblob")`,
		`insert into testcase_byte values (null, null, null, null, null, null)`,
	},
	{
		`insert into testcase_text values ("tinytext", "mediumtext", "text", "longtext")`,
		`insert into testcase_text values (null, null, null, null)`,
	},
	{
		`insert into testcase_enum_set values ("x-small", null)`,
		`insert into testcase_enum_set values (null, ('a,d'))`,
	},
	{
		`insert into testcase_spatial values (ST_GeomFromText('point(108.9498710632 34.2588125935)'))`,
		`insert into testcase_spatial values (null)`,
	},
	{
		`insert into testcase_json values ('{"key1": "value1", "key2": "value2"}')`,
		`insert into testcase_json values (null)`,
	},
}

// 测试用例使用的数据库，建议和das放到一个DB实例
var (
	testCaseUser     string = "das_rw"
	testCasePassword string = "das@1234.Com"
	testCaseHost     string = "127.0.0.1"
	testCasePort     int    = 3306
	testCaseDatabase string = "das"
)

func init() {
	db := DB{
		User:     testCaseUser,
		Password: testCasePassword,
		Host:     testCaseHost,
		Port:     testCasePort,
		Database: testCaseDatabase,
		Params:   map[string]string{"sql_mode": "\"\""},
		Ctx:      context.Background(),
	}
	// 初始化测试用例的表
	for key, value := range testCaseTables {
		fmt.Println("=== 初始化用例表: ", key)
		// 删除表
		if err := db.Execute(fmt.Sprintf("drop table IF EXISTS `%s`", key)); err != nil {
			panic(err)
		}
		// 新建表
		if err := db.Execute(value); err != nil {
			panic(err)
		}
	}
	// 初始化测试用例的数据
	for _, i := range testCaseData {
		for _, j := range i {
			if strings.HasPrefix(j, "insert") {
				fmt.Println("=== 初始化用例数据: ", j)
				if err := db.Execute(j); err != nil {
					panic(err)
				}
			}
		}
	}
}

func TestDb(t *testing.T) {
	// 定义用例
	testCases := []struct {
		name    string
		sql     string
		wantRes []map[string]interface{}
	}{
		{
			name: "TEST NullInt64",
			sql:  "select cast(null as unsigned) as NullInt64",
			wantRes: []map[string]interface{}{
				{"NullInt64": interface{}(nil)},
			},
		},
		{
			name: "TEST BOOL",
			sql:  "SELECT IF(0, 'true', 'false') as bool_false, IF(1, 'true', 'false') as bool_true",
			wantRes: []map[string]interface{}{
				{"bool_false": "false", "bool_true": "true"},
			},
		},
		{
			name: "TEST BIT",
			sql:  "select col_bit from testcase_bit",
			wantRes: []map[string]interface{}{
				{"col_bit": "\x02"},
				{"col_bit": interface{}(nil)},
				{"col_bit": "\x00"},
			},
		},
		{
			name: "TEST HEX BIT",
			sql:  "select hex(col_bit) as hex_col_bit from testcase_bit",
			wantRes: []map[string]interface{}{
				{"hex_col_bit": "2"},
				{"hex_col_bit": interface{}(nil)},
				{"hex_col_bit": "0"},
			},
		},
		{
			name: "TEST TINYINT",
			sql:  "select col_tinyint from testcase_int",
			wantRes: []map[string]interface{}{
				{"col_tinyint": "127"},
				{"col_tinyint": "-12"},
				{"col_tinyint": interface{}(nil)},
			},
		},
		{
			name: "TEST TINYINT UNSIGNED",
			sql:  "select col_tinyint_unsigned  from testcase_int",
			wantRes: []map[string]interface{}{
				{"col_tinyint_unsigned": "255"},
				{"col_tinyint_unsigned": "255"},
				{"col_tinyint_unsigned": "0"},
			},
		},
		{
			name: "TEST SMALLINT",
			sql:  "select col_smallint  from testcase_int",
			wantRes: []map[string]interface{}{
				{"col_smallint": "32767"},
				{"col_smallint": "-3276"},
				{"col_smallint": interface{}(nil)},
			},
		},
		{
			name: "TEST SMALLINT UNSIGNED",
			sql:  "select col_smallint_unsigned  from testcase_int",
			wantRes: []map[string]interface{}{
				{"col_smallint_unsigned": "65535"},
				{"col_smallint_unsigned": "65535"},
				{"col_smallint_unsigned": "0"},
			},
		},
		{
			name: "TEST MEDIUMINT",
			sql:  "select col_mediumint from testcase_int",
			wantRes: []map[string]interface{}{
				{"col_mediumint": "8388607"},
				{"col_mediumint": "-838860"},
				{"col_mediumint": interface{}(nil)},
			},
		},
		{
			name: "TEST MEDIUMINT UNSIGNED",
			sql:  "select col_mediumint_unsigned from testcase_int",
			wantRes: []map[string]interface{}{
				{"col_mediumint_unsigned": "16777215"},
				{"col_mediumint_unsigned": "16777215"},
				{"col_mediumint_unsigned": "0"},
			},
		},
		{
			name: "TEST INT",
			sql:  "select col_int from testcase_int",
			wantRes: []map[string]interface{}{
				{"col_int": "2147483647"},
				{"col_int": "-214748364"},
				{"col_int": interface{}(nil)},
			},
		},
		{
			name: "TEST INT UNSIGNED",
			sql:  "select col_int_unsigned from testcase_int",
			wantRes: []map[string]interface{}{
				{"col_int_unsigned": "4294967295"},
				{"col_int_unsigned": "4294967295"},
				{"col_int_unsigned": "0"},
			},
		},
		{
			name: "TEST BIGINT",
			sql:  "select col_bigint from testcase_int",
			wantRes: []map[string]interface{}{
				{"col_bigint": "9223372036854775807"},
				{"col_bigint": "-922337203685477580"},
				{"col_bigint": interface{}(nil)},
			},
		},
		{
			name: "TEST BIGINT UNSIGNED",
			sql:  "select col_bigint_unsigned from testcase_int",
			wantRes: []map[string]interface{}{
				{"col_bigint_unsigned": "18446744073709551615"},
				{"col_bigint_unsigned": "18446744073709551615"},
				{"col_bigint_unsigned": "0"},
			},
		},
		{
			name: "TEST DECIMAL",
			sql:  "select * from testcase_decimal",
			wantRes: []map[string]interface{}{
				{"col_decimal": "100.22", "col_double": "100", "col_float": "100"},
				{"col_decimal": "-10.22", "col_double": "-100", "col_float": "-100"},
				{"col_decimal": interface{}(nil), "col_double": interface{}(nil), "col_float": interface{}(nil)},
			},
		},
		{
			name: "TEST DATE",
			sql:  "select * from testcase_date",
			wantRes: []map[string]interface{}{
				{"col_date": "2023-06-06", "col_datetime": "2023-06-06 12:00:00", "col_time": "12:00:00", "col_timestamp": "2023-06-06 12:00:00", "col_year": "2023"},
				{"col_date": "0000-00-00", "col_datetime": "2023-06-06 12:00:00", "col_time": "00:00:00", "col_timestamp": "2023-06-06 12:00:00", "col_year": "0000"},
				{"col_date": interface{}(nil), "col_datetime": interface{}(nil), "col_time": interface{}(nil), "col_timestamp": "2023-06-06 12:00:00", "col_year": interface{}(nil)},
				{"col_date": "0000-00-00", "col_datetime": "0000-00-00 00:00:00", "col_time": "00:00:00", "col_timestamp": "0000-00-00 00:00:00", "col_year": "0000"},
			},
		},
		{
			name: "TEST CHAR",
			sql:  "select * from testcase_char",
			wantRes: []map[string]interface{}{
				{"col_char": "char", "col_char_not_null": "char_not_null", "col_varchar": "varchar", "col_varchar_not_null": "varchar_not_null"},
				{"col_char": interface{}(nil), "col_char_not_null": "", "col_varchar": interface{}(nil), "col_varchar_not_null": ""},
			},
		},
		{
			name: "TEST BYTE",
			sql:  "select * from testcase_byte",
			wantRes: []map[string]interface{}{
				{"col_binary": "a\x00\x00", "col_blob": "blob", "col_longblob": "longblob", "col_mediumblob": "mediumblob", "col_tinyblob": "tinyblob", "col_varbinary": "ab"},
				{"col_binary": interface{}(nil), "col_blob": interface{}(nil), "col_longblob": interface{}(nil), "col_mediumblob": interface{}(nil), "col_tinyblob": interface{}(nil), "col_varbinary": interface{}(nil)},
			},
		},
		{
			name: "TEST TEXT",
			sql:  "select * from testcase_text",
			wantRes: []map[string]interface{}{
				{"col_longtext": "longtext", "col_mediumtext": "mediumtext", "col_text": "text", "col_tinytext": "tinytext"},
				{"col_longtext": interface{}(nil), "col_mediumtext": interface{}(nil), "col_text": interface{}(nil), "col_tinytext": interface{}(nil)},
			},
		},
		{
			name: "TEST ENUM SET",
			sql:  "select * from testcase_enum_set",
			wantRes: []map[string]interface{}{
				{"col_enum": "x-small", "col_set": interface{}(nil)},
				{"col_enum": interface{}(nil), "col_set": "a,d"},
			},
		},
		{
			name: "TEST SPATIAL",
			sql:  "select ST_AsText(col_geometry) as col_geometry from testcase_spatial",
			wantRes: []map[string]interface{}{
				{"col_geometry": "POINT(108.9498710632 34.2588125935)"},
				{"col_geometry": interface{}(nil)},
			},
		},
		{
			name: "TEST JSON",
			sql:  "select * from testcase_json",
			wantRes: []map[string]interface{}{
				{"col_json": "{\"key1\": \"value1\", \"key2\": \"value2\"}"},
				{"col_json": interface{}(nil)},
			},
		},
	}
	// 初始化连接
	db := DB{
		User:     testCaseUser,
		Password: testCasePassword,
		Host:     testCaseHost,
		Port:     testCasePort,
		Database: testCaseDatabase,
		Ctx:      context.Background(),
	}
	// 运行
	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			_, data, err := db.Query(c.sql)
			if err != nil {
				t.Fatal(err)
			}
			assert.EqualValues(t, c.wantRes, *data)
		})
	}
}

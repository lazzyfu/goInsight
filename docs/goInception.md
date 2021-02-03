- [goInception](#goinception)
- [下载安装包](#下载安装包)
- [解压](#解压)
- [编辑配置文件](#编辑配置文件)
- [配置到supervisor](#配置到supervisor)
- [启动goInception服务](#启动goinception服务)
- [重启django服务](#重启django服务)

#### goInception
goInception地址：

https://github.com/hanchuanchuan/goInception/tags

> goInception是一个非常棒的开源项目，YaSQL目前仅用到了goInception的语法审核功能

#### 下载安装包
下载最新的Releases版本即可

`wget https://github.com/hanchuanchuan/goInception/releases/download/v1.2.3/goInception-linux-amd64-v1.2.3.tar.gz`

#### 解压
```bash
mkdir /usr/local/goInception
tar -zxf goInception-linux-amd64-v1.2.3.tar.gz -C /usr/local/goInception
cd /usr/local/goInception
```

#### 编辑配置文件
vim config/config.toml
>请根据自己业务的实际情况进行修改，参考goInception官方文档调整参数即可
```editorconfig
# IP地址
host = "127.0.0.1"

# 端口
port = 4000

# TiDB数据库目录
path = "/tmp/tidb"

# 忽略终端连接断开信号
ignore_sighup = true

[log]
# 日志级别: debug, info, warn, error, fatal.
level = "info"

# 日志格式, one of json, text, console.
format = "text"

# 禁用时间戳输出
disable-timestamp = false

# 日志文件
[log.file]
# 日志文件名
filename = ""

# 日志文件的最大上限(MB)
max-size = 300

# Max日志文件的保存天数，默认值 `0`，即不清理
max-days = 0

# 要保留的最大旧日志文件数，默认值 `0`，即不清理
max-backups = 0

# 日志轮询，默认值 `true`，即开启
log-rotate = true

# 审核规则
[inc]
# 自增列
## 当建表时自增列的类型不为int或者bigint时报错
check_autoincrement_datatype = true
## 当建表时自增列的值指定的不为1，则报错
check_autoincrement_init_value = true
## 自增列必须要为无符号型
enable_autoincrement_unsigned = true

# 主键
## 是否强制主键列必须是int
enable_pk_columns_only_int = true
## 建表时，如果没有主键，则报错
check_primary_key = true


# 注释
## 建表时，列没有注释时报错
check_column_comment = true
## 建表时，表没有注释时报错
check_table_comment = true

# 默认值
## 检查在建表、修改列、新增列时，新的列属性是不是要有默认值
check_column_default_value = false
## 允许blob/text/json类型置为null
enable_blob_not_null = false
## 允许列类型为BLOB/TEXT
enable_blob_type = true

# 关键字
## 不检查关键字了，历史原因
enable_identifer_keyword = true

# JSON
## 允许列类型为JSON
enable_json_type = true

# 检查字段类型变更
## 不允许变更字段类型，兼容tidb同步
check_column_type_change = true

# LIMIT
## 允许在DML语句中使用了LIMIT时
check_dml_limit = true

# 不能有order by 语句
check_dml_orderby = true

# WHERE
## 在DML语句中没有WHERE条件时，是不是要报错
check_dml_where = true

# 限制一条insert values的总行数
max_insert_rows = 100

# 必须指定插入列表，也就是要写入哪些列
check_insert_field = true

# 当update/delete预估受影响行数超出设置值时警告
max_update_rows = 10000

# 不允许select *
enable_select_star = true

# float/double
## 开启时,当使用 float/double 类型时提示转成 decimal 类型
check_float_double = true

# 最大char长度,当超出时警告转换为varchar类型
max_char_length = 30

# 检查表名/索引名前缀
check_index_prefix = true
index_prefix = "IDX_"
uniq_index_prefix = "uniq_"
table_prefix = ""

# 一个索引最多可指定的列数
max_key_parts = 5

# 单表允许的最大索引数
max_keys = 12

# 在多个改同一个表的语句出现是，报错，提示合成一个
# 当检测数据库类型为tidb时，动态修改
merge_alter_table = true

# 用以指定建表时必须创建的列。多个列时以逗号分隔
must_have_columns = 'D_CREATED_AT,D_UPDATED_AT'

# 字符集设置
enable_set_charset = true
enable_set_collation = true
support_charset = "utf8,utf8mb4"
## 允许列自己设置字符集
enable_column_charset = true

# DROP
## 不允许drop database
enable_drop_database = false
## 是否允许删除表
enable_drop_table = true

# 枚举类型
## 是不是支持enum,set,bit数据类型
enable_enum_set_bit = false

# 是不是支持外键
enable_foreign_key = false

# 创建索引时是否允许空索引名
enable_null_index_name = true

# 创建或者新增列时是否允许列为NULL
enable_nullable = false

# 允许change列
enable_change_column = true

# ENGINE
support_engine = "innodb"

# DISABLE TIMESTAMP类型
enable_timestamp_type = false

# explain判断受影响行数时使用的规则, 默认值"first"
# 可选值: "first", "max"
#      "first":    使用第一行的explain结果作为受影响行数
#      "max":      使用explain结果中的最大值作为受影响行数
explain_rule = "first"

# 返回的信息使用语言
lang = "zh-CN"

# 全量日志
general_log = false

[osc]
osc_on = false

[ghost]
ghost_on = false

# 自定义审核级别用以实现指定审核结果的错误级别
# inception show levels
[inc_level]
er_alter_table_once = 1
er_auto_incr_id_warning = 1
er_autoinc_unsigned = 1
er_blob_cant_have_default = 1
er_cant_change_column = 1
er_cant_change_column_position = 1
er_cant_set_charset = 1
er_cant_set_collation = 1
er_cant_set_engine = 1
er_change_column_type = 1
er_change_too_much_rows = 1
er_char_to_varchar_len = 1
er_charset_on_column = 1
er_column_have_no_comment = 1
er_datetime_default = 1
er_foreign_key = 2
er_ident_use_keyword = 1
er_implicit_type_conversion = 1
er_inc_init_err = 1
er_index_name_idx_prefix = 1
er_index_name_uniq_prefix = 1
er_insert_too_much_rows = 1
er_invalid_data_type = 1
er_invalid_ident = 1
er_join_no_on_condition = 1
er_json_type_support = 2
er_must_have_columns = 1
er_columns_must_have_index = 1
er_columns_must_have_index_type_err = 1
er_no_where_condition = 1
er_not_allowed_nullable = 1
er_ordery_by_rand = 1
er_partition_not_allowed = 1
er_pk_cols_not_int = 1
er_pk_too_many_parts = 1
er_select_only_star = 1
er_set_data_type_int_bigint = 2
er_table_charset_must_null = 1
er_table_charset_must_utf8 = 1
er_table_must_have_comment = 1
er_table_must_have_pk = 1
er_table_prefix = 1
er_text_not_nullable_error = 1
er_timestamp_default = 1
er_too_many_key_parts = 1
er_too_many_keys = 1
er_too_much_auto_datetime_cols = 2
er_too_much_auto_timestamp_cols = 2
er_udpate_too_much_rows = 1
er_use_enum = 1
er_use_text_or_blob = 2
er_with_default_add_column = 1
er_with_insert_field = 1
er_with_limit_condition = 1
er_with_orderby_condition = 1
er_wrong_and_expr = 1
```

#### 配置到supervisor
vim /etc/supervisord.d/yasql.conf
```editorconfig
[program:goInception]
user=root
autorestart=true
autostart=true
directory=/data/www/yasql/yasql
command=/usr/local/goInception/goInception --config=/usr/local/goInception/config/config.toml
redirect_stderr=true
stdout_logfile=/data/www/yasql/yasql/logs/goinception.log
```

#### 启动goInception服务
编辑config.py文件，配置goInception的地址和端口
```python
# GoInception
INCEPTION = {
    'host': '127.0.0.1',
    'port': 4000
}
```

启动goInception服务
```bash
supervisorctl update
supervisorctl start goInception
```

#### 重启django服务
`supervisorctl restart yasql-server`
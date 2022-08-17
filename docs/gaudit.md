- [介绍](#介绍)
- [下载安装包](#下载安装包)
- [解压](#解压)
- [创建配置文件](#创建配置文件)
- [修改系统配置](#修改系统配置)
- [自定义审核参数](#自定义审核参数)
- [配置supervisor](#配置supervisor)
- [授权](#授权)
- [YaSQL配置gAudit地址](#yasql配置gaudit地址)
- [启动gAudit服务](#启动gaudit服务)
- [重启django服务](#重启django服务)

### 介绍
> gAudit是一个SQL语法审核工具，支持MySQL/TiDB，通过解析SQL语法树实现语法规则审核。

https://github.com/lazzyfu/gAudit

### 下载安装包
> 下载最新的Releases版本即可，请选择tar.gz结尾的

`wget https://github.com/lazzyfu/gAudit/releases/download/v1.0.1/gAudit-linux-v1.0.1.tar.gz`

### 解压
```bash
mkdir /usr/local/gAudit
tar -jxf gAudit-linux-v1.0.1.tar.gz -C /usr/local/gAudit
```

### 创建配置文件
> 根据您自己的需求调整审核参数

模板文件: template/config.json
格式: json
审核参数: 请参考[审核参数](parameters.md)进行自定义增加或调整
文件：`/usr/local/gAudit/config.json`


### 修改系统配置
`vim /usr/local/gAudit/config.json`
```json
"ListenAddress": "127.0.0.1:8082",
"LogFilePath": "./logs",
"LogLevel": "debug",
```
### 自定义审核参数
> 有时需要临时放开某个审核规则，不希望每次修改配置文件然后去重启gAudit服务

**审核参数生效优先级**
`自定义传递参数（custom_audit_parameters）` > `template/config.json` > `config/config.go`

- **config/config.go**
系统内置审核参数，优先级最低，可以被`template/config.json`覆盖

- **template/config.json**
自定义的模板参数文件，您启动时加载的配置文件，可以被`自定义传递参数（custom_audit_parameters）`覆盖

- **custom_audit_parameters**
POST请求时自定义传参，优先级最高。支持一次传递多个审核参数（系统参数除外）

**custom_audit_parameters使用方法**
> 请参考[审核参数](parameters.md), 参数名不区分大小写
- 不允许使用关键字
- 必须要有审核字段
```
"custom_audit_parameters": {"check_identifer_keyword": true,"check_table_audit_type_columns": true},
```

### 配置supervisor
vim /etc/supervisord.d/yasql.conf
```editorconfig
[program:gAudit]
user=root
autorestart=true
autostart=true
directory=/data/www/yasql/yasql
command=/usr/local/gAudit/gAudit -config /usr/local/gAudit/config.json
redirect_stderr=true
stdout_logfile=/data/www/yasql/yasql/logs/gAudit.log
```

### 授权
> 语法审核用到的权限：SELECT, INSERT, UPDATE, DELETE

凡是需要进行语法审核的业务库，均需要进行创建账号和授权

```sql
create user 'gaudit_rw'@'%' identified by '1234.com';
GRANT SELECT, INSERT, UPDATE, DELETE ON *.* TO 'gaudit_rw'@'%'
```

由于工单系统会执行工单、进行备份和查询等功能，因此需要其他的权限

```sql 
GRANT all ON *.* TO 'gaudit_rw'@'%'
```

### YaSQL配置gAudit地址
编辑config.py文件，配置gAudit的地址和端口

vim yasql/config.py
```python
# 指定审核用户
REOMOTE_USER = {
    'user': 'gaudit_rw',
    'password': '1234.com'
}

# gAudit
GAUDIT_API = "http://127.0.0.1:8082/api/v1/audit"
```


### 启动gAudit服务
```bash
supervisorctl update
supervisorctl start gAudit
```

### 重启django服务
`supervisorctl restart yasql-server`
# goInsight

![](https://img.shields.io/github/languages/top/lazzyfu/goInsight)
![](https://img.shields.io/static/v1?label=License&message=MIT&color=green&?style=flat-square)
<img alt="Github Stars" src="https://img.shields.io/github/stars/lazzyfu/goInsight?logo=github">

goInsight是一个集数据库工单和数据检索的平台，提供语法审核、工单流、库表数据查询能力，旨在规范流程、降低风险、解放DBA生产力。

## 文档
:point_right: [goInsight文档](https://github.com/lazzyfu/goInsight/wiki)

## 功能
- **SQL审计** — 支持对工单SQL进行语法审核，内置自研语法审核器[gAudit](https://github.com/lazzyfu/gAudit)，支持多达68种语法审核规则。支持动态调整审核参数，除全局审核参数外，可在DB实例级别定义自定义的审核参数，实现每个DB不同的审核规则。
- **数据检索** - 内置库表数据检索引擎，支持为用户配置库表级别数据只读查询权限。支持多达65种只读查询语句类型，如SELECT、EXPLAIN、SHOW等，支持查询记录审计，支持MySQL（含分支、云RDS）、TiDB、ClickHouse。
- **组织管理** - 支持用户组织层级管理，可以为每个DB实例绑定组织，实现用户只访问当前组织的数据库，支持用户角色管理。
- **其他功能**
  - 支持双因素身份认证登录、支持日志审计。
  - 支持提交DDL/DML工单提交、审批、执行、钩子等功能。
  - 支持钉钉/企业微信机器人/邮件推送消息。
  
## 限制 
  - **工单** - 目前仅支持MySQL5.6+、TiDB数据库提交DML/DDL语句，单个工单最大限制2048条SQL语句。
  - **查询** - 目前仅支持MySQL、TiDB、ClickHouse数据库，默认返回100条记录（可调整），仅支持只读操作，变更请使用工单。


## 安装

### 下载二进制包

下载前请更新下载链接中的`版本号`为最新的`发行版本`。

```
wget https://github.com/lazzyfu/goInsight/releases/download/v1.0.0/goInsight-linux-amd64-v1.0.2.tar.gz
```

### 安装gh-ost工具

`MySQL DDL`改表工单自动集成`gh-ost`工具，请提前安装[gh-ost](https://github.com/github/gh-ost)。

### 启动服务

请确保已正确配置config.yaml。

```
./goInsight -config config.yaml
```

### 访问Web界面

现在您可以访问goInsight的Web地址：<http://ip:port>，请确保防火墙放通了您启动服务时指定的port端口。

请输入系统默认的管理员账号密码，管理员默认未开启双因素身份认证，登录后请及时修改密码。

- 默认管理员账号：admin
- 默认管理员密码：1234.Com!


## 联系作者

E-mail: `1126227133@qq.com`

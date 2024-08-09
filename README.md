## GoInsight
![](https://img.shields.io/static/v1?label=License&message=MIT&color=green&?style=flat-square)
[![downloads](https://img.shields.io/github/downloads/lazzyfu/goInsight/total.svg)](https://github.com/lazzyfu/goInsight/releases)
[![release](https://img.shields.io/github/v/release/lazzyfu/goInsight.svg)](https://github.com/lazzyfu/goInsight/releases)
<img alt="Github Stars" src="https://img.shields.io/github/stars/lazzyfu/goInsight?logo=github">

`GoInsight`是一个基于`Gin`和`Vue`开发的企业级前后端分离的WEB端平台，支持MySQL/TiDB/ClickHouse数据库，提供**工单管理**和**数据查询**功能。

## 文档
- [功能预览](https://github.com/lazzyfu/goInsight/wiki/Feature-Preview)
- [使用文档](https://github.com/lazzyfu/goInsight/wiki)

## 主要功能
- **工单管理（管理工单的生命周期）**
  - 支持提交DDL/DML/数据导出（CSV/XLSX格式）工单
  - 支持工单审批、执行、回滚、HOOK等功能
  - 支持钉钉/企业微信机器人/邮件推送工单消息

- **SQL审计（规则拦得住）**
  - 内置[gAudit](https://github.com/lazzyfu/gAudit)语法审核器，支持多达70种语法审核规则，可联系开发者定制更多规则
  - 除全局审核参数外，可在DB实例级别自定义审核参数（优先级最高）

- **数据查询**
  - 强大的库表数据检索引擎，支持多大65种只读查询语句类型（SELECT、CTE、EXPLAIN、SHOW等）
  - 支持用户库表级权限管理、支持数据字典、查询记录审计等功能

- **组织管理**
  - 支持用户组织和角色管理
  - 支持DB实例绑定组织，实现用户访问权限控制

- **自定义工单环境**
  - 支持定义多个工单环境（如研发环境、预发布环境、测试环境、生产环境等）
  - 支持不同环境绑定不同数据库实例，同一个工单可以在不同环境HOOK，避免重复提交

- **其他功能** 
  - 支持管理员统一管理数据库实例
  - 支持双因素身份认证登录
  - 支持日志审计

## 支持数据库
| 数据库                                          | 版本 | 功能                                                                                                                                                                                      |
| ----------------------------------------------- | -------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| MySQL数据库（含分支、云厂商数据库、AWS Aurora） | >= 5.7   | &#9745; DML工单，支持生成DML回滚语句 <br> &#9745; DDL工单，支持ONLINE DDL（集成gh-ost工具） <br> &#9745; 导出工单，支持加密导出XLSX/CSV格式文件 <br> &#9745; 数据查询，支持库表级权限控制 |
| TiDB                                            | >= 4.0   | &#9745; DML工单，暂时不支持生成DML回滚语句 <br> &#9745; DDL工单，TiDB原生支持ONLINE DDL <br> &#9745; 导出工单，支持加密导出XLSX/CSV格式文件 <br> &#9745; 数据查询，支持库表级权限控制     |
| Clickhouse                                      | >= 18.1  | &#9745; 数据查询，支持库表级权限控制                                                                                                                                                      |


## 快速部署
> [使用Supervisor管理服务](https://github.com/lazzyfu/goInsight/wiki/Service-Deployment)

**1、下载二进制包**

```
# 下载前请更新下载链接中的`版本号`为最新的`发行版本`
wget https://github.com/lazzyfu/goInsight/releases/download/v1.3.5/goinsight-linux-amd64-v1.3.5.tar.gz
```

**2、配置config.yaml**

配置文件相对简单，请参考: [如何配置config.yaml](https://github.com/lazzyfu/goInsight/wiki/Config)

**3、启动服务**

```
./goInsight -config config.yaml
```

**4、访问服务**
服务启动后，请通过 <http://ip:port> 访问，登录后请及时修改密码。生产使用时，强烈建议您在该服务前加一层HTTPS代理。
 - 默认管理员账号：`admin`
 - 默认管理员密码：`1234.Com!`

## 联系作者
- E-mail: [1126227133@qq.com](mailto:1126227133@qq.com)
- 微信号: Lazzyfu
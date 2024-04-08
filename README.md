## 基于Gin + Vue + Ant Design前后端分离的工单系统
![](https://img.shields.io/static/v1?label=License&message=MIT&color=green&?style=flat-square)
[![downloads](https://img.shields.io/github/downloads/lazzyfu/goInsight/total.svg)](https://github.com/lazzyfu/goInsight/releases)
[![release](https://img.shields.io/github/v/release/lazzyfu/goInsight.svg)](https://github.com/lazzyfu/goInsight/releases)
<img alt="Github Stars" src="https://img.shields.io/github/stars/lazzyfu/goInsight?logo=github">

goInsight是集`数据库工单`和`数据查询`的平台，提供`语法审核`、`工单审批执行`、`数据查询`能力，支持`MySQL/TiDB/ClickHouse`。旨在规范变更流程、降低数据库管理员、研发测试等人员使用成本。支持本地快速部署，使用简单，注重用户隐私保护。

## 文档
- [功能预览](https://github.com/lazzyfu/goInsight/wiki/Feature-Preview)
- [使用文档](https://github.com/lazzyfu/goInsight/wiki)

## 主要功能
- **SQL审计** - 支持对工单SQL进行语法审核，内置自研语法审核器[gAudit](https://github.com/lazzyfu/gAudit)，支持多达68种语法审核规则。支持动态调整审核参数，除全局审核参数外，可在DB实例级别定义自定义的审核参数，实现每个DB不同的审核规则。
- **数据查询** - 内置强大的库表数据检索引擎，可为用户配置库表级别数据查询权限。支持多达65种只读查询语句类型，如SELECT、CTE、EXPLAIN、SHOW等，支持查询记录审计。
- **工单流** - 支持DDL、DML、数据导出（CSV/XLSX格式）工单，支持审批、执行、Hook功能，单个工单最大支持提交2048条SQL语句。支持钉钉/企业微信机器人/邮件推送消息。
- **组织管理** - 支持用户组织层级管理，可以为每个DB实例绑定组织，实现用户只访问当前组织的数据库，支持用户角色管理。
- **自定义工单环境** - 支持自定义工单环境，可根据自己场景定义多个工单环境并绑定不同数据库实例，如研发环境、预发布环境、测试环境、生产环境等。 
- **安全** - 支持双因素身份认证登录、支持日志审计。

## 支持数据库
| 数据库                               | 版本   | 功能                                                                                                                          |
| ------------------------------------ | ------ | ----------------------------------------------------------------------------------------------------------------------------- |
| MySQL/华为云RDS/阿里云RDS/AWS Aurora等 | >= 5.7  | &#9745; DDL/DML工单 <br> &#9745; 导出工单 <br> &#9745; 数据查询 <br> &#9745; DML回滚 <br> &#9745; Online DDL（集成gh-ost工具）     |
| TiDB                                 | >= 4.0  | &#9745; DDL/DML工单 <br> &#9745; 导出工单 <br> &#9745; 数据查询 <br> &#9744; DML回滚 <br> &#9745; Online DDL（TiDB原生支持） |
| Clickhouse                           | >= 18.1 | &#9744; DDL/DML工单 <br> &#9744; 导出工单 <br> &#9745; 数据查询                                              |

## 快速部署
> [!NOTE]
> [建议使用Supervisor管理服务](https://github.com/lazzyfu/goInsight/wiki/Service-Deployment)

#### 1.下载二进制包
> 下载前请更新下载链接中的`版本号`为最新的`发行版本`

```
wget https://github.com/lazzyfu/goInsight/releases/download/v1.3.1/goinsight-linux-amd64-v1.3.1.tar.gz
```

#### 2.启动服务

启动服务前，请确保已正确配置`config.yaml`。

```
./goInsight -config config.yaml
```

#### 3.访问服务

现在您可以通过<http://ip:port>访问服务（请确保防火墙放通了您启动服务时指定的port端口）。然后输入系统默认的管理员账号密码，管理员默认未开启双因素身份认证，登录后请及时修改密码。

```
默认管理员账号：admin
默认管理员密码：1234.Com!
```

## 联系作者

E-mail: `1126227133@qq.com`
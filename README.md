## GoInsight
![](https://img.shields.io/static/v1?label=License&message=MIT&color=green&?style=flat-square)
[![downloads](https://img.shields.io/github/downloads/lazzyfu/goInsight/total.svg)](https://github.com/lazzyfu/goInsight/releases)
[![release](https://img.shields.io/github/v/release/lazzyfu/goInsight.svg)](https://github.com/lazzyfu/goInsight/releases)
<img alt="Github Stars" src="https://img.shields.io/github/stars/lazzyfu/goInsight?logo=github">

`GoInsight`是一个基于`Gin`和`Vue`开发的前后端分离Web平台，支持MySQL/TiDB/ClickHouse数据库，提供**工单管理**和**数据查询**功能。

## 文档
- [功能预览](https://github.com/lazzyfu/goInsight/wiki/Feature-Preview)
- [使用文档](https://github.com/lazzyfu/goInsight/wiki)

## 主要功能
- **工单管理**
  - 提交DDL、DML、数据导出（CSV/XLSX格式）工单
  - 支持审批、执行、回滚、HOOK功能
  - 单个工单最大支持提交2048条SQL语句
  - 支持钉钉/企业微信机器人/邮件推送消息

- **SQL审计**
  - 内置[gAudit](https://github.com/lazzyfu/gAudit)语法审核器
  - 支持多达70种语法审核规则
  - 除全局审核参数外，可在DB实例级别自定义审核参数

- **数据查询**
  - 强大的库表数据检索引擎
  - 支持 65 种只读查询语句类型（SELECT、CTE、EXPLAIN、SHOW 等）
  - 查询记录审计

- **组织管理**
  - 用户组织和角色管理
  - 可DB 实例绑定组织，实现用户访问权限控制

- **自定义工单环境**
  - 支持定义多个工单环境（如研发、预发布、测试、生产环境等）
  - 不同环境绑定不同数据库实例

- **数据库实例管理** 
  - 统一管理数据库实例

- **安全**
  - 双因素身份认证登录
  - 日志审计

## 支持数据库
| 数据库                                 | 版本    | 功能                                                                                               |
| -------------------------------------- | ------- | -------------------------------------------------------------------------------------------------- |
| MySQL/云RDS/Aurora等 | >= 5.7  | ✅ DDL/DML工单 <br> ✅ 导出工单 <br> ✅ 数据查询 <br> ✅ DML回滚 <br> ✅ Online DDL（集成 gh-ost 工具） |
| TiDB                                   | >= 4.0  | ✅ DDL/DML工单 <br> ✅ 导出工单 <br> ✅ 数据查询 <br> ❌ DML回滚 <br> ✅ Online DDL（TiDB 原生支持）    |
| ClickHouse                             | >= 18.1 | ❌ DDL/DML工单 <br> ❌ 导出工单 <br> ✅ 数据查询                                                      |


## 快速部署
> [使用Supervisor管理服务](https://github.com/lazzyfu/goInsight/wiki/Service-Deployment)

**1.下载二进制包**

```
# 下载前请更新下载链接中的`版本号`为最新的`发行版本`
wget https://github.com/lazzyfu/goInsight/releases/download/v1.3.3/goinsight-linux-amd64-v1.3.3.tar.gz
```

**2.启动服务**

```
# 确保已正确配置 config.yaml
./goInsight -config config.yaml
```

**3.访问服务**
 - 通过 <http://ip:port> 访问
 - 默认管理员账号：`admin`
 - 默认管理员密码：`1234.Com!`
 - 登录后请及时修改密码

## 联系作者

E-mail: [1126227133@qq.com](mailto:1126227133@qq.com)
## GoInsight

![](https://img.shields.io/static/v1?label=License&message=MIT&color=green&?style=flat-square)
[![downloads](https://img.shields.io/github/downloads/lazzyfu/goInsight/total.svg)](https://github.com/lazzyfu/goInsight/releases)
[![release](https://img.shields.io/github/v/release/lazzyfu/goInsight.svg)](https://github.com/lazzyfu/goInsight/releases)
<img alt="Github Stars" src="https://img.shields.io/github/stars/lazzyfu/goInsight?logo=github">

`GoInsight` 是一个基于 `Gin` 和 `Vue3` 开发的企业级前后端分离的 WEB 平台，支持 MySQL/TiDB 数据库，提供 **工单管理** 和 **数据查询** 功能。

## 文档

- [功能预览](https://github.com/lazzyfu/goInsight/wiki/Feature-Preview)
- [使用文档](https://github.com/lazzyfu/goInsight/wiki)

## 主要功能

- **工单管理**
  - 支持提交 DDL/DML/数据导出（CSV/XLSX 格式）工单
  - 支持工单审批、执行、回滚、HOOK 等功能
  - 支持钉钉 / 企业微信机器人 / 邮件推送工单消息

- **审批流**
  - 支持会签、或签
  - 自定义审批流程和审批层级

- **SQL 审计**
  - 内置语法审核器，支持多达 68 种语法审核规则
  - 支持全局审核参数、实例自定义审核参数，优先级：实例级别 > 全局级别

- **数据查询**
  - 强大的库表数据检索引擎，支持多达 65 种只读查询语句类型（SELECT、CTE、EXPLAIN、SHOW 等）
  - 支持用户库表级权限管理、支持数据字典、查询记录审计等功能

- **组织管理**
  - 支持用户组织和角色管理
  - 支持 DB 实例绑定组织，实现用户访问权限控制

- **自定义工单环境**
  - 支持定义多个工单环境（如研发、预发布、测试、生产等）
  - 支持为不同环境配置独立的审批与参数

- **安全**
  - 支持双因素身份认证登录
  - 支持日志审计

## 支持数据库

| 数据库 | 版本 | 功能概览 |
| --- | --- | --- |
| MySQL 数据库（含分支、云厂商数据库、AWS Aurora） | >= 5.7 | ✓ DML 工单（含回滚语句）；✓ ONLINE DDL（gh-ost）；✓ 导出 XLSX/CSV（加密可选）；✓ 库表级查询权限 |
| TiDB | >= 4.0 | ✓ DML 工单（暂不生成回滚）；✓ 原生 ONLINE DDL；✓ 导出 XLSX/CSV（加密可选）；✓ 库表级查询权限 |

## 快速部署
>
> [使用Supervisor管理服务](https://github.com/lazzyfu/goInsight/wiki/Service-Deployment)

### 1、下载二进制包

```
# 下载前请更新下载链接中的`版本号`为最新的`发行版本`
wget https://github.com/lazzyfu/goInsight/releases/download/v2.0.0/goinsight-linux-amd64-v2.0.0.tar.gz
```

### 2、配置 config.yaml

配置文件相对简单，请参考: [如何配置config.yaml](https://github.com/lazzyfu/goInsight/wiki/Config)

### 3、启动服务

```
./goinsight-linux-amd64 -config config.yaml
```

### 4、访问服务

服务启动后，请通过 <http://ip:port> 访问，首次登录后请立即修改默认密码。生产环境强烈建议在服务前加一层 HTTPS 反向代理。

- 默认管理员账号：`admin`
- 默认管理员密码：`1234.Com!`，请登录后立刻更改并妥善保管。

## 联系开发者

- E-mail: [1126227133@qq.com](mailto:1126227133@qq.com)
- 微信号: Lazzyfu

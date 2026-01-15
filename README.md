# 🚀 GoInsight

<p align="center">
  <img src="https://img.shields.io/static/v1?label=License&message=MIT&color=2ea44f&style=flat-square" />
  <a href="https://github.com/lazzyfu/goInsight/releases">
    <img src="https://img.shields.io/github/v/release/lazzyfu/goInsight.svg?style=flat-square&color=orange" />
  </a>
  <a href="https://github.com/lazzyfu/goInsight/releases">
    <img src="https://img.shields.io/github/downloads/lazzyfu/goInsight/total.svg?style=flat-square&color=blue" />
  </a>
  <img src="https://img.shields.io/github/stars/lazzyfu/goInsight?style=flat-square&logo=github" />
</p>

<p align="center">
  <b> 企业级数据库工单与数据查询平台 </b>
</p>

<p align="center">
  <img src="https://img.shields.io/badge/Backend-Gin-00ADD8?style=for-the-badge&logo=go" />
  <img src="https://img.shields.io/badge/Frontend-Vue3-4FC08D?style=for-the-badge&logo=vue.js" />
  <img src="https://img.shields.io/badge/Database-MySQL%20%7C%20TiDB-4479A1?style=for-the-badge&logo=mysql" />
</p>

## 📖 项目简介

**GoInsight** 是一个基于 **Gin + Vue3** 的开发企业级数据库WEB管理平台，提供**数据库工单**和**数据查询**能力，解决DBA与开发人员之间协作低效、操作高危、审计缺失等痛点。

## 🗺️ 文档与预览

* 📌 [功能预览](https://github.com/lazzyfu/goInsight/wiki/Feature-Preview)

## 🛠️ 核心功能模块

| 功能模块       | 功能说明                                                                        |
| :------------- | :------------------------------------------------------------------------------ |
| **📝 工单管理** | 支持DDL/DML/数据导出（CSV / XLSX）工单；支持审批、执行、回滚及一键复制工单 |
| **⚖️ 审批流**   | 支持会签/或签，可按组织架构自定义审批流程与层级                          |
| **🔍 SQL 审计** | 内置高性能语法审核器，支持68种审核规则；支持全局及实例级参数配置         |
| **🔓 数据查询** | 支持库表级访问权限控制与查询审计           |
| **👥 权限组织** | 基于组织架构的用户角色管理；DB实例按组织绑定，实现访问隔离               |
| **🌐 环境隔离** | 支持自定义多个工单环境（测试、预发、生产），确保变更发布流程规范化              |
| **🔔 消息通知** | 集成钉钉/企业微信/邮件消息推送                      |
| **🛡️ 其它** | 支持双因素认证                       |

## 🗄️ 支持数据库

| 数据库引擎         | 版本要求 | 核心能力说明                                                        |
| :----------------- | :------- | :------------------------------------------------------------------ |
| **MySQL / 云服务** | ≥ 5.7    | DML工单（含回滚）；Online DDL (gh-ost)；数据加密导出；数据查询库表级权限控制 |
| **TiDB**           | ≥ 4.0    | DML工单（不支持生成回滚SQL）；原生Online DDL；数据加密导出；数据查询库表级权限控制     |

## 🚀 快速部署

### 1. 下载二进制包

前往 [Releases](https://github.com/lazzyfu/goInsight/releases) 页面获取最新的安装包。

### 2. 启动服务

```bash
./goinsight-linux-amd64 -config config.yaml
```

### 3. 访问系统

在浏览器打开：<http://your-ip:port>

* 默认管理员： admin

* 默认密码： 1234.Com!

⚠️ 安全建议： 首次登录请立即修改密码。生产环境强烈建议在 GoInsight 前置部署 Nginx 并配置 HTTPS 反向代理。

## 🤝 联系与支持

如果您在使用过程中发现 Bug 或有功能建议，欢迎提交 Issue 或 Pull Request。

* 📧 Email: <1126227133@qq.com>

* 💬 微信: Lazzyfu

# goInsight

![](https://img.shields.io/static/v1?label=Go&message=1.20&color=green&?style=flat-square)
![](https://img.shields.io/static/v1?label=Vue&message=AntDesignVue&color=green&?style=flat-square)
![](https://img.shields.io/static/v1?label=License&message=MIT&color=green&?style=flat-square)

:two_hearts:goInsight是一个集数据库工单提交、审核、执行和数据查询的平台，支持MySQL/TiDB。


:point_right:[功能预览](https://github.com/lazzyfu/goInsight/wiki/预览)

### 核心功能

- 支持用户管理功能、角色、组织层级功能、双因素身份认证登录
- 支持提交DDL/DML工单、语法审核、审批流、基于组织级别的权限控制，工单执行、回滚SQL
- 支持数据查询功能、用户库表级别数据查询权限、支持黑名单、支持审计
- 支持钉钉/企业微信机器人/邮件推送消息

### 安装说明

#### 一、下载二进制包

> 下载前请更新为最新的发行版本

```
wget https://github.com/lazzyfu/goInsight/releases/download/v1.0.0/goInsight-linux-amd64-v1.0.0.tar.gz
```

#### 二、集成语法审核工具

请下载安装：<https://github.com/lazzyfu/gAudit>

#### 三、安装gh-ost工具

请下载安装：<https://github.com/github/gh-ost>

#### 四、启动服务

> 请先调整配置文件，否则可能因访问不到数据库而无法启动

```
./goInsight -config config.yaml
```

#### 五、访问Web

访问地址：<http://ip:port>

请输入系统默认的管理员账号密码，管理员默认未开启双因素身份认证，登录后请及时修改密码。

- 默认管理员账号：admin
- 默认管理员密码：1234.Com!


### 联系作者

E-mail: `1126227133@qq.com`

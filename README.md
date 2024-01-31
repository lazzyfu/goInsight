# goInsight

![](https://img.shields.io/static/v1?label=Go&message=1.20&color=green&?style=flat-square)
![](https://img.shields.io/static/v1?label=Vue&message=AntDesignVue&color=green&?style=flat-square)
![](https://img.shields.io/static/v1?label=License&message=MIT&color=green&?style=flat-square)

:two_hearts:goInsight是一个集语法审核、工单提交/审批/执行和数据查询的平台，旨在规范上线流程、减少风险、解放DBA生产力。


:point_right:[点击查看预览图和文档](https://github.com/lazzyfu/goInsight/wiki)

### 核心功能

- 支持用户管理、角色管理、组织层级管理、双因素身份认证登录等功能
- 内置语法审核功能，支持多达68种语法审核规则，支持全局、DB实例级别定义审核参数，可为不同DB实例配置不同的审核规则
- 支持DDL/DML工单、审批流、权限控制
- 支持工单执行、执行结果、执行日志、DML生成回滚SQL语句
- 支持数据查询功能、支持为用户配置库表级别数据查询权限、支持黑名单、支持查询审计
- 支持钉钉/企业微信机器人/邮件推送消息
- 支持MySQL5.6+/TiDB4+

### 安装说明

#### 1.下载二进制包

> 下载前请更新下载链接中的版本号为最新的发行版本

```
wget https://github.com/lazzyfu/goInsight/releases/download/v1.0.0/goInsight-linux-amd64-v1.0.2.tar.gz
```


#### 2.安装gh-ost工具
> MySQL DDL工单自动集成gh-ost工具，请提前安装gh-ost

请下载安装：<https://github.com/github/gh-ost>

#### 3.启动服务

> 请先更改配置文件，否则可能因访问不到数据库等情况而启动失败

```
./goInsight -config config.yaml
```

#### 4.访问Web界面

现在您可以访问goInsight的地址：<http://ip:port>，请确保防火墙放通了您启动服务时指定的port端口。

请输入系统默认的管理员账号密码，管理员默认未开启双因素身份认证，登录后请及时修改密码。

- 默认管理员账号：admin
- 默认管理员密码：1234.Com!


### 联系作者

E-mail: `1126227133@qq.com`

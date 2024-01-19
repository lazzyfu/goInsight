# goInsight
![](https://img.shields.io/static/v1?label=Go&message=1.20&color=green&?style=flat-square)
![](https://img.shields.io/static/v1?label=Vue&message=AntDesignVue&color=green&?style=flat-square)
![License](https://img.shields.io/github/license/lazzyfu/goInsight?style=flat-square)

goInsight是一个集数据库工单提交、审核、执行和数据查询的平台，旨在规范上线流程，解放DBA生产力。

### 安装说明
#### 下载二进制包
```
```

#### 启动服务
> 请先调整配置文件，否则可能因访问不到数据库无法启动
```
./goInsight -config config.yaml
```

#### 访问WEB界面
访问http://ip:port，输入系统默认的管理员账号密码，管理员默认未开启双因素身份认证，登录后请及时修改密码。
- 默认管理员账号：admin
- 默认管理员密码：1234.Com!

#### 集成语法审核工具
请下载安装：https://github.com/lazzyfu/gAudit

#### 安装gh-ost工具
请下载安装：https://github.com/github/gh-ost

### 支持功能
* 支持用户管理功能，双因素身份认证、定义任意层级的组织功能
* 支持工单审批流，基于组织级别的权限控制
* 支持工单执行，工单执行记录、工单消息推送，支持钉钉/企业微信机器人/邮件推送
* 支持数据查询功能，支持用户定义表级别数据查询权限
* 支持语法审核，集成gAudit语法审核
* 支持提交DML/DDL工单；支持MySQL5.7+/TiDB Version 4+
  * MySQL DDL自动集成gh-ost工具
  * TiDB ONLINE DDL
  * 支持MySQL DML生成回滚SQL
  * 【限制】暂不支持TiDB DML语句生成回滚SQL

### 功能预览
#### 用户
![](https://github.com/lazzyfu/goInsight/blob/master/docs/pics/user.png)

#### 数据查询
![](https://github.com/lazzyfu/goInsight/blob/master/docs/pics/das.png)

#### 提交DML/DDL工单
![](https://github.com/lazzyfu/goInsight/blob/master/docs/pics/dml.png)
![](https://github.com/lazzyfu/goInsight/blob/master/docs/pics/ddl.png)

#### 工单列表和详情
![](https://github.com/lazzyfu/goInsight/blob/master/docs/pics/orders.png)
![](https://github.com/lazzyfu/goInsight/blob/master/docs/pics/order_detail.png)

#### 工单消息通知
![](https://github.com/lazzyfu/goInsight/blob/master/docs/pics/msg_commit.png)
![](https://github.com/lazzyfu/goInsight/blob/master/docs/pics/msg_audit.png)

### 联系
E-mail: `1126227133@qq.com`
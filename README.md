# AuditSQL

AuditSQL是web版本的MySQL数据库审核平台，旨在降低DBA的运维成本.

## 开发组件
```text
 Python 3.6+
 Django 2.0+
 Celery 4.2.0+
 Django-channel
 Paramiko
 Pymysql
 Python-mysql-replication
```

  
## 功能
```text
* 自定义环境
  - 可以根据自己环境需求，定义多个环境，比如：测试环境、预发布环境、生产环境等

* 工单
  - 支持DML、DDL、运维工单、数据导出
  - 支持钩子、上线版本号等功能
 
* 语法
  - 提供SQL语法审核、格式化、高亮、注释、补全等功能
 
* 审核
  - 通过Inception进行语法规则检测(目前inception仅做此用途)
  - 审核人(部门Leader)审核工单，是否通过
  - 执行人(DBA)执行工单，反馈工单状态

* 执行
  - DDL语句支持gh-ost改表，抛弃pt-osc
  - DML语句实现事务级别的执行保证
  - 支持一键「全部执行」和有选择的「单条执行」
 
* 进度展示
  - 前台实时显示DML、DDL、数据导出的进度（websocket）

* 推送通知
  - 支持钉钉
  - 支持邮件
  - all

* SQL查询
  - 支持SQL查询功能，库授权，limit等功能
 
* Xterm
  - 支持Xterm接口，可以接入redis、mongodb等命令行接口

* 其他功能
   - 支持LDAP或本地手动创建用户授权登陆功能，支持修改头像，密码等功能
   - 支持自动检测已配置数据库实例的死锁检测功能，并提供推送
```

## 设计流程
![](https://github.com/lazzyfu/AuditSQL/blob/master/media/png/design.png)

## 效果展示
https://github.com/lazzyfu/AuditSQL/wiki/show

## 文档（使用和安装）
https://github.com/lazzyfu/AuditSQL/wiki

## 联系方式：
QQ群号：710797678

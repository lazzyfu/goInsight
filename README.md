# AuditSQL

AuditSQL是web版本的MySQL数据库审核平台，旨在降低DBA的运维成本.

- __开发组件__
   - Python 3.6+
   - Django 2.0+
   - Celery 4.2.0+
   - Django-channel
   - Paramiko
   - Pymysql
   - Python-mysql-replication

- __功能__
   - 支持自定义环境（比如：测试环境、预发布环境、生产环境）
   - 支持DML、DDL、运维工单
   - 提供SQL语法审核、格式化、高亮、注释、补全等功能
   - 并提供上线版本号支持
   - 支持DDL、DML工单一键自动执行（DML事务保证）
   - 支持gh-ost、原生alter改表（抛弃了pt-osc）
   - 提供工单执行时，实时进度显示和执行日志预览功能
   - 支持工单流 支持工单提交、审核、执行、反馈、关闭、钩子等操作
   - 支持钩子功能，可勾到定义的环境（不需要重复提交工单）
   - 支持SQL查询，查询库授权、查询日志审核等功能
   - 支持多种推送方式，邮件、钉钉（后台支持一键开关）
   - 支持LDAP或本地手动创建用户授权登陆功能，支持修改头像，密码等功能
   - 提供xterm集成，支持绑定redis、mongodb等查询接口
   - 支持自动检测已配置数据库实例的死锁检测功能，并提供推送

## 设计流程

![](https://github.com/lazzyfu/AuditSQL/blob/master/media/png/design.png

## 效果展示
https://github.com/lazzyfu/AuditSQL/wiki/show

## 文档（使用和安装）
https://github.com/lazzyfu/AuditSQL/wiki

## 迭代周期
https://github.com/lazzyfu/AuditSQL/blob/master/upgrade.txt

## 联系方式（扫码加入）：

![](https://github.com/lazzyfu/AuditSQL/blob/master/media/png/ql.png)

## 打赏（感谢支持）：

![](https://github.com/lazzyfu/AuditSQL/blob/master/media/png/ds.png)
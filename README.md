# AuditSQL

AuditSQL是基于Inception开发的web版本的MySQL数据库审核平台，旨在降低DBA的运维成本，解放大家的生产力。

请大家放心使用，若是遇到使用上的问题或bug，请提出Issues，我会及时关注并给出解决办法。若是觉得ok，请给颗**Star**，谢谢。

本系统经过生产业务验证，稳定，可长期使用，支持rds，谢谢。

## 开发组件

* Python 3.6+
* Django 2.0+
* Celery 4.2.0+
* Django-channels
* AdminLTE
* Paramiko

## 功能
- 支持自定义级联环境（比如：测试环境--> 预发布环境 -->生产环境）
- 支持DDL、DML的SQL语法审核、提示、美化、高亮、注释识别、补全等功能
- 支持DDL、DML上线工单，并提供上线版本号支持
- 支持DDL、DML工单一键自动执行（可选单条执行或全部执行）
- 支持gh-ost、pt-osc、原生alter改表
- 支持工单流审核操作，提供：工单提交、工单执行、工单审核等权限
- 支持钩子功能，DDL只需提交一次，逐级环境勾取即可
- 支持SQL查询，查询库授权、查询日志审核等功能
- 提供部分工单的回滚功能
- 提供工单执行时，实时进度显示和执行日志预览功能
- 支持多种推送方式，邮件、钉钉（后台支持一键开关）
- 支持LDAP或本地手动创建用户授权登陆功能，支持修改头像，密码等功能
- 提供xterm集成，支持绑定redis、mongodb等查询接口
- 支持自动检测已配置数据库实例的死锁检测功能，并提供推送


## 文档地址
https://github.com/lazzyfu/AuditSQL/wiki


## 页面展示(简单展示几处)

**登陆页面：**

![](https://github.com/lazzyfu/AuditSQL/blob/master/media/png/login.png)

**个人详情页面：**

![](https://github.com/lazzyfu/AuditSQL/blob/master/media/png/profile.png)

**上线版本页面：**

![](https://github.com/lazzyfu/AuditSQL/blob/master/media/png/version.png)

**DDL工单页面：**

![](https://github.com/lazzyfu/AuditSQL/blob/master/media/png/ddl.png)

**测试环境页面：**

![](https://github.com/lazzyfu/AuditSQL/blob/master/media/png/test.png)

**执行任务页面：**

![](https://github.com/lazzyfu/AuditSQL/blob/master/media/png/perform.png)

**使用pt-osc改表：**

![](https://github.com/lazzyfu/AuditSQL/blob/master/media/png/perform_ddl.png)

**使用gh-ost改表：**

![](https://github.com/lazzyfu/AuditSQL/blob/master/media/png/ddl_ghost.png)

**查询页面：**

![](https://github.com/lazzyfu/AuditSQL/blob/master/media/png/query.png)

**查看表结构和索引：**

![](https://github.com/lazzyfu/AuditSQL/blob/master/media/png/table.png)

**xterm页面：**

![](https://github.com/lazzyfu/AuditSQL/blob/master/media/png/webshell.png)


## 死锁钉钉通知：
```text
【警告 ◕﹏◕，探测到新的死锁记录，探测时间：2018-08-29 05:13:31】

## 死锁记录1 ##:
主机：192.168.10.19
时间: 2018-08-29 05:12:58+00:00
线程ID: 6102754
事务ID: 0
事务激活时间: 15
用户名: test_user
主机名:
IP: 192.168.0.100
库名: test_11
表名: data_realtime_trans_adc_statixxx
发生死锁的索引: idx_adc_id
锁类型: RECORD
锁模式: X
请求锁: w
是否回滚: 否
查询: REPLACE INTO data_realtime_trans_adc_statxxx (`adc_

主机：192.168.10.19
时间: 2018-08-29 05:12:58+00:00
线程ID: 6102764
事务ID: 0
事务激活时间: 5
用户名: test_user
主机名:
IP: 192.168.0.100
库名: test_11
表名: data_realtime_trans_adc_statixxx
发生死锁的索引: idx_adc_id
锁类型: RECORD
锁模式: X
请求锁: w
是否回滚: 是
查询: REPLACE INTO data_realtime_trans_adc_statixxx (`adc_

@所有人
```

## 联系方式

群号：710797678

E-mail: 1126227133@qq.com
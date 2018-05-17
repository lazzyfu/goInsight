AuditSQL介绍
===============

**AuditSQL(命名确实是个问题，^_^)是一个web审核平台，集成了inception，旨在降低DBA的运维成本。**

**欢迎大家的使用和好评，同时使用中遇到的问题，请指出，谢谢。**

**本项目会持续进行维护。**

## Wiki地址（详细部署和使用说明，请务必参考）

**https://github.com/lazzyfu/AuditSQL/wiki**

## 开发组件

- Python 3.6
- Django 2.0 
- celery 4.1.0
- django-channels
- AdminLTE

## 核心功能简介（更多功能，请使用时体验）

- 线上工单审核（近实时E-Mail通知）
  - 数据变更（DML和DDL审核）
  - 数据导出（xlsx格式，压缩和加密）

- 线下工单审核
  - SQL审核（DML和DDL审核）
  
- 工单历史
   - 审核记录：提供工单审核和操作历史的详细记录
   - 执行任务：提供线上/线下执行任务的记录和操作

- 定时任务
  - crontab：自定义crontab
  - 表结构监控
  - 数据库备份（mysqldump和xtrabackup）

- 数据库管理
   - 数据库账号的友好显示
   - 数据库账号的管理功能，支持新建用户、主机、变更权限操作
   - 提供指定数据库账号保护功能（被保护的账号无法被操作）

- 审核功能
   - SQL审核的流程化，规范化
   - SQL美化功能
   - SQL检测功能
   - 语法高亮功能
   - 注释识别功能
   - SQL语法自动补全(包括表名和列名)

- 数据导出功能
   - 根据提供的select语句，自动异步导出xlsx格式的数据
   - 支持导出数据的压缩和加密
   - 导出数据的邮件提醒功能

- 执行任务功能
   - 自动分片当前审核内容，并生产SQL执行任务列表
   - 提供回滚，inception执行日志预览功能
   - 提供DDL语句的OSC进度实时显示和停止功能

- 推送
   - 线上工单E-Mail近实时推送
   - 执行任务执行进度的近实时显示（使用websocket）
  
- 权限系统：
  - 独立`项目组`控制，每个`项目组`之间互相隔离（包括：用户、联系人，主机、角色等）
  - 自定义用户角色：如：DBA、项目经理、开发、产品等角色
  - 提供8种权限，可绑定到任意用户角色

- 其他
   - 支持LDAP认证登陆
   - 支持修改头像
   - 支持修改密码


## 线上审核简要流程图

![](https://github.com/lazzyfu/AuditSQL/blob/master/media/png/lct.png)

## 页面展示(随便展示几处)

![](https://github.com/lazzyfu/AuditSQL/blob/master/media/png/1.png)

![](https://github.com/lazzyfu/AuditSQL/blob/master/media/png/2.png)

![](https://github.com/lazzyfu/AuditSQL/blob/master/media/png/3.png)

![](https://github.com/lazzyfu/AuditSQL/blob/master/media/png/5.png)

![](https://github.com/lazzyfu/AuditSQL/blob/master/media/png/6.png)

![](https://github.com/lazzyfu/AuditSQL/blob/master/media/png/sql_tips.png)

![](https://github.com/lazzyfu/AuditSQL/blob/master/media/png/backup-list.png)

![](https://github.com/lazzyfu/AuditSQL/blob/master/media/png/backup_create.png)

![](https://github.com/lazzyfu/AuditSQL/blob/master/media/png/backup-preview.png)

![](https://github.com/lazzyfu/AuditSQL/blob/master/media/png/email_notice.png)

## 联系方式
   
QQ群: 710797678

E-mail: 1126227133@qq.com

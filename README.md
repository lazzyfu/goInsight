AuditSQL介绍
===============

**AuditSQL(命名确实是个问题，^_^)是基于Inception开发的一款web审核平台，旨在降低DBA的运维成本。**

**欢迎大家的使用和好评，同时使用中遇到的问题，请指出，谢谢。**

**本项目会持续进行维护，感谢！**

## Wiki地址（详细使用说明，请务必参考）

**https://github.com/lazzyfu/AuditSQL/wiki**

## 开发组件

- Python 3.6
- Django 2.0 
- django-celery
- django-channels
- AdminLTE

## 核心功能简介（更多功能，请使用时体验）

- 线上工单（近实时E-Mail通知）

- 线下工单

- SQL审核
   - SQL审核的流程化，规范化
   - SQL美化功能
   - SQL检测功能
   - 语法高亮功能
   - 注释识别功能

- 工单历史
   - 提供工单操作历史的详细记录

- 数据导出
   - 根据提供的select语句，自动异步导出xlsx格式的数据
   - 支持导出数据的压缩和加密
   - 导出数据的邮件提醒功能

- 自动化执行任务
   - 自动分片当前审核内容，并生产SQL执行任务列表
   - 提供回滚，inception执行日志预览功能
   - 提供DDL语句的OSC进度实时显示和停止功能

- 数据库管理
   - 数据库账号的友好显示
   - 数据库账号的管理功能，支持新建用户、主机、变更权限操作
   - 提供指定数据库账号保护功能（被保护的账号无法被操作）

- 推送
   - 线上工单E-Mail近实时推送
   - 执行任务执行进度的近实时显示（使用websocket）
  
- 其他
   - 支持LDAP认证登陆
   - 项目权限控制
   - 角色权限控制
   - 支持修改头像
   - 支持修改密码
  
- 扩展功能：
   - 支持数据库表结构变更自动E-Mail通知，并生成变更结果

## 简要流程图

![](https://github.com/lazzyfu/AuditSQL/blob/master/media/gif/liuchengtu.png)

## 页面展示(随便展示几处)

![](https://github.com/lazzyfu/AuditSQL/blob/master/media/show/show-1.png)

![](https://github.com/lazzyfu/AuditSQL/blob/master/media/show/show-2.png)

![](https://github.com/lazzyfu/AuditSQL/blob/master/media/show/show-3.png)

![](https://github.com/lazzyfu/AuditSQL/blob/master/media/show/show-4.png)

![](https://github.com/lazzyfu/AuditSQL/blob/master/media/show/show-5.png)


## 联系方式
   
QQ群: 710797678

E-mail: 1126227133@qq.com
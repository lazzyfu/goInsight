AuditSQL介绍
===============

**AuditSQL(命名确实是个问题，^_^)是一个web版的MySQL数据库审核平台，作为一个DBA，厌烦了SQL上线过程的语法审核、各个环境手动支持，反馈等重复且不讨好的工作，
于是开发此款工具，旨在降低DBA的运维成本，解放大家的生产力。**

**欢迎大家的使用，如果觉得不错，麻烦伸出您高贵的小手，给颗star，同时使用中遇到的问题，请联系我，谢谢（^_^）**

## 文档地址
https://github.com/lazzyfu/AuditSQL/wiki

## 开发组件

- Python 3.6+
- Django 2.0+
- celery 4.2.0
- django-channels
- AdminLTE

## 核心功能简介（更多功能，请使用时体验）
- 历史工单
  - 生产环境
  - 预发布环境
  - 测试环境
  - 执行任务记录

- SQL审核
  - DML和DDL语法审核(Inception支持)

- 工单
  - DML变更工单
  - DDL上线工单
  - 上线版本号

- 数据查询
  - 生产mysql
  - 非生产mysql
  - mongo和redis(xterm)

- 功能
   - SQL审核的流程化，规范化
   - SQL美化功能
   - SQL检测功能
   - 语法高亮功能
   - 注释识别功能
   - SQL语法自动补全(包括表名和列名)

- 执行任务功能
   - 自动分片，支持一键串行全部执行或有选择的执行
   - 提供回滚，inception执行日志预览功能
   - 提供DDL语句的OSC进度实时输出

- 推送
   - 实时钉钉推送，友好推送（用户更新手机号，直接@用户）
   - 执行任务执行进度的实时显示（websocket）

- 其他
   - 支持LDAP认证登陆
   - 支持修改头像
   - 支持用户自己修改密码（非ldap方式认证）

## 移除的功能
1. 考虑到钉钉的便捷性和及时性，仅支持钉钉推送，移除了邮件推送

## 页面展示(随便展示几处)
![](https://github.com/lazzyfu/AuditSQL/blob/master/media/png/test_env.png)

![](https://github.com/lazzyfu/AuditSQL/blob/master/media/png/staging_env.png)

![](https://github.com/lazzyfu/AuditSQL/blob/master/media/png/dml_gongdan.png)

![](https://github.com/lazzyfu/AuditSQL/blob/master/media/png/ddl_gongdan.png)

![](https://github.com/lazzyfu/AuditSQL/blob/master/media/png/task.png)

![](https://github.com/lazzyfu/AuditSQL/blob/master/media/png/perform_task.png)

![](https://github.com/lazzyfu/AuditSQL/blob/master/media/png/pro_query.png)

![](https://github.com/lazzyfu/AuditSQL/blob/master/media/png/offline_query.png)

![](https://github.com/lazzyfu/AuditSQL/blob/master/media/png/mongo_query.png)

![](https://github.com/lazzyfu/AuditSQL/blob/master/media/png/sys_config.png)

![](https://github.com/lazzyfu/AuditSQL/blob/master/media/png/dingding.png)


## 联系方式

群号：710797678

E-mail: 1126227133@qq.com
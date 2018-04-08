AuditSQL介绍
===============

**AuditSQL是基于Inception开发的一款web审核平台，旨在降低DBA的运维成本。**
**欢迎大家的使用和好评，同时有使用上的问题，请指出，谢谢**

Wiki地址（使用说明）
===============

**https://github.com/lazzyfu/AuditSQL/wiki**


流程图
===============

![](https://github.com/lazzyfu/AuditSQL/blob/master/media/gif/liuchengtu.png)

开发组件
===============
- Python 3.6
- Django 2.0 
- django-celery
- django-channels
- AdminLTE

功能
===============
- 线上工单审核（近实时E-Mail通知）
- 线下工单审核
- 执行任务
   - 显示当前任务SQL列表
   - 执行（执行当前SQL，实时显示当前执行的OSC任务进度）
   - 停止（关闭当前正在执行的OSC任务，并显示结果）
   - 结果（显示回滚SQL，inception执行日志）
   - 回滚（对执行的SQL进行回滚操作）

- SQL审核
   - 流程化
   - SQL美化功能
   - SQL检测功能
   - 审核历史记录
   - 语法高亮功能
   - 注释识别功能

- 进度推送：
   - 线上工单E-Mail近实时推送
   - 执行任务执行进度的近实时显示
  
- 其他：
   - 支持LDAP认证登陆
   - 项目权限控制
   - 角色权限控制
   - 支持修改头像
   - 支持修改密码
  
- 扩展功能：
   - 支持数据库表结构变更自动E-Mail通知，并生成变更结果


## 线上SQL审核提交流程(点击查看GIF动图）

![gif动图](https://github.com/lazzyfu/AuditSQL/blob/master/media/gif/2018-03-15%2009_31_03.gif)
![提交发送的邮件](https://github.com/lazzyfu/AuditSQL/blob/master/media/gif/mail-1.png)

## 线上SQL审批流程(点击查看GIF动图）

![审批流程](https://github.com/lazzyfu/AuditSQL/blob/master/media/gif/2.gif)

## 线下SQL执行任务流程(点击查看GIF动态)

![线下SQL流程](https://github.com/lazzyfu/AuditSQL/blob/master/media/gif/of-op.gif)

## 联系方式
   
   QQ群: 710797678
   
   E-mail: 1126227133@qq.com
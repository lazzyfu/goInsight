# AuditSQL

AuditSQL是web版本的MySQL数据库审核平台，旨在降低DBA的运维成本.
若是遇到使用上的问题或bug，请提出Issues，我会及时关注并给出解决办法.
本系统经过生产业务验证，稳定，可长期使用，支持rds，谢谢.

- __开发组件__
   * Python 3.6+
   * Django 2.0+
   * Celery 4.2.0+

- __功能__
   - 支持自定义级联环境（比如：测试环境--> 预发布环境 -->生产环境）
   - 语法审核、提示、美化、高亮、注释识别、补全等功能
   - 并提供上线版本号支持
   - 支持DDL、DML工单一键自动执行（可选单条执行或全部执行）
   - 支持gh-ost、原生alter改表（抛弃了pt-osc）
   - 支持工单流操作，提供：工单提交、审核、执行、反馈、关闭等
   - 支持钩子功能，可勾到定义的级联环境
   - 支持SQL查询，查询库授权、查询日志审核等功能
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


## 联系方式

群聊(欢迎加入)：

![](https://github.com/lazzyfu/AuditSQL/blob/master/media/png/ql.png)


打赏：


![](https://github.com/lazzyfu/AuditSQL/blob/master/media/png/ds.png)
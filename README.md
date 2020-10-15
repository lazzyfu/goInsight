# YaSQL简介
YaSQL是一款简单的MySQL/TiDB的数据库审核执行平台，旨在降低DBA的运维成本，规范线上SQL执行流程。

## 开发组件
* ![](https://img.shields.io/badge/Python_3.7-green.png)
* ![](https://img.shields.io/badge/Ant_Design_Vue-green.png)
* ![](https://img.shields.io/badge/Django_2.0+-green.png)
* ![](https://img.shields.io/badge/Djangorestframework-green.png)
* ![](https://img.shields.io/badge/Djangorestframework_jwt-green.png)
* ![](https://img.shields.io/badge/Celery-green.png)

## 功能简介
* 自定义工单环境
  * 可根据业务场景，自定义工单的环境，比如：测试环境、预发布环境、生产环境等等
* 支持的工单类型
  * DML工单
  * DDL工单
  * 导出工单（支持SELECT语句导出为CSV、XLSX格式）
* 支持的DB类型
  * MySQL（>=5.6）
  * TiDB
* 语法规则
  * 集成goInception，语法规则请参考：https://github.com/hanchuanchuan/goInception
  * 支持语法高亮、格式化、注释、补全
* 消息推送功能
  * 邮件
  * 钉钉
  * 企业微信
* 其他功能
  * 支持钩子功能，工单可以在各个环境内自由hook
  * MySQL DDL ALTER语句自动使用gh-ost改表
  * DML语句执行、支持获取锁定超时、事务封装、自动开启严格模式执行
  * 单个工单最大支持2048条SQL语句
  * MySQL DML语句支持自动生成回滚SQL（真实影响行数小于10W行）
  * 支持执行前台实时展示（基于websocket实现）
  * 支持自定义用户角色、可为每个角色绑定不同的工单权限
  * 支持集成LDAP（若支持LDAP密码修改，请自行实现相关接口）
  * 支持后台创建用户、密码修改（非LDAP模式）、修改头像
  
## WIKI地址(项目截图预览)
https://github.com/lazzyfu/YaSQL/wiki

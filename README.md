# YaSQL简介
YaSQL是MySQL/TiDB/ClickHouse的数据库审核/执行/查询平台，旨在降低DBA的运维成本，规范SQL操作流程。

## 开发组件
* ![](https://img.shields.io/static/v1?label=Python&message=3.7&color=green&?style=for-the-badge)
* ![](https://img.shields.io/static/v1?label=Vue&message=Ant_Design_Vue&color=green&?style=for-the-badge)
* ![](https://img.shields.io/static/v1?label=Django&message=Djangorestframework&color=green&?style=for-the-badge)
* ![](https://img.shields.io/static/v1?label=Jwt&message=Djangorestframework_jwt&color=green&?style=for-the-badge)
* ![](https://img.shields.io/static/v1?label=Celery&message=4&color=green&?style=for-the-badge) 

## 功能简介

* 自定义工单环境
  * 可根据业务场景，自定义工单的环境，比如：测试环境、预发布环境、生产环境等等
* 支持的工单类型
  * DML工单
  * DDL工单
  * 导出工单（支持SELECT语句导出为CSV、XLSX格式）
* 支持的DB类型
  * MySQL（>=5.6， Percona Server & 官方MySQL）
  * TiDB
  * Mariadb（暂不支持Mariadb生成备份语句）
  * ClickHouse（暂不支持生成备份语句）
* 语法规则
  * 集成goInception，语法规则请参考：https://github.com/hanchuanchuan/goInception
  * 前端支持语法高亮、格式化、注释、补全
* 支持MySQL/TIDB/ClickHouse/Redis查询（类似DMS查询功能）
  * 支持库表级别授权
  * 支持基于用户/组授权
  * 支持查询审计功能（自动记录用户执行的SQL）
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
 
## 使用文档
- [YaSQL预览](Home)
- YaSQL部署
  - [介绍](https://github.com/lazzyfu/YaSQL/wiki/%E4%BB%8B%E7%BB%8D)
  - [初始化环境](https://github.com/lazzyfu/YaSQL/wiki/%E5%88%9D%E5%A7%8B%E5%8C%96%E7%8E%AF%E5%A2%83)
  - [克隆项目](https://github.com/lazzyfu/YaSQL/wiki/%E5%85%8B%E9%9A%86%E9%A1%B9%E7%9B%AE)
  - [部署前端服务](https://github.com/lazzyfu/YaSQL/wiki/%E9%83%A8%E7%BD%B2%E5%89%8D%E7%AB%AF%E6%9C%8D%E5%8A%A1)
  - [部署后端服务](https://github.com/lazzyfu/YaSQL/wiki/%E9%83%A8%E7%BD%B2%E5%90%8E%E7%AB%AF%E6%9C%8D%E5%8A%A1)
  - [集成goInception](https://github.com/lazzyfu/YaSQL/wiki/%E9%9B%86%E6%88%90goInception)
  - [集成gh-ost](https://github.com/lazzyfu/YaSQL/wiki/%E9%9B%86%E6%88%90gh-ost)
- YaSQL用户管理
  - [集成LDAP](https://github.com/lazzyfu/YaSQL/wiki/%E9%9B%86%E6%88%90LDAPP)
  - [用户管理](https://github.com/lazzyfu/YaSQL/wiki/%E7%94%A8%E6%88%B7%E7%AE%A1%E7%90%86)
- YaSQL工单配置和使用
  - [配置工单环境](https://github.com/lazzyfu/YaSQL/wiki/%E9%85%8D%E7%BD%AE%E5%B7%A5%E5%8D%95%E7%8E%AF%E5%A2%83)
  - [配置审核数据库](https://github.com/lazzyfu/YaSQL/wiki/%E9%85%8D%E7%BD%AE%E5%AE%A1%E6%A0%B8%E6%95%B0%E6%8D%AE%E5%BA%93)
  - [工单钩子](https://github.com/lazzyfu/YaSQL/wiki/%E5%B7%A5%E5%8D%95%E9%92%A9%E5%AD%90)
  - [工单流程](https://github.com/lazzyfu/YaSQL/wiki/%E5%B7%A5%E5%8D%95%E6%B5%81%E7%A8%8B)
- YaSQL查询配置和使用
  - [查询配置](https://github.com/lazzyfu/YaSQL/wiki/%E6%9F%A5%E8%AF%A2%E9%85%8D%E7%BD%AE)
  - [配置权限](https://github.com/lazzyfu/YaSQL/wiki/%E9%85%8D%E7%BD%AE%E6%9D%83%E9%99%90)
  - [配置返回行数](https://github.com/lazzyfu/YaSQL/wiki/%E9%85%8D%E7%BD%AE%E8%BF%94%E5%9B%9E%E8%A1%8C%E6%95%B0)
  - [配置查询超时](https://github.com/lazzyfu/YaSQL/wiki/%E9%85%8D%E7%BD%AE%E6%9F%A5%E8%AF%A2%E8%B6%85%E6%97%B6)
- [升级](https://github.com/lazzyfu/YaSQL/wiki/%E5%8D%87%E7%BA%A7)
 

## QQ讨论群
<img src="https://github.com/lazzyfu/YaSQL/blob/master/example_pic/qq.png" alt="" align=center />

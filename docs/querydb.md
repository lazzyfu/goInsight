# 查询功能
1. 提供给开发快速检索各个环境下数据库的数据
2. 支持LIMIT规制改写
3. 支持基于用户或组级别的权限配置
4. 支持库或表级别的权限分配
5. 支持安全审计

#### 配置只读账号
`vim yasql/config.py`

```python
# 作为开发查询数据库使用
# 该账户仅设置为只读即可
# create user 'yasql_ro'@'%' identified by '1234.com'
# grant select on *.* to 'yasql_ro'@'%';
# 用户名和密码请进行自行修改，不要使用默认的
QUERY_USER = {
    'user': 'yasql_ro',
    'password': '1234.com'
}
```

将配置文件中的用户和密码改为自己设置的，然后保存，使用supervisorctl重启服务即可


#### 创建只读账号
请在需要执行查询的目标数据库上创建上面已配置的账号，并设置为只读（SELECT）权限
> 请务必保证当前主机能够访问远程数据库
> 日志位置：yasql/logs/celery.log


#### 配置主机
登录后台：{your domain_name}/admin/，找到：SQL工单配置 -> DB主机配置 -> 增加DB主机配置 -> 录入远程数据库实例信息
> 用途：SQL查询


#### 配置定时采集任务
登录后台：{your domain_name}/admin/，找到：PERIODIC TASKS -> Periodic tasks -> 增加PERIODIC TASKS

新增一个任务：
- name：DMS-同步库表信息
- Task(registered)：sqlquery.tasks.sqlquery_sync_schemas_tables
- Enabled：勾选
- Interval Schedule：10 minutes （这里随意配置，建议周期设置为分钟级别，比如：10分钟，集群主机越多，建议越大）

点击保存即可


#### 手动触发定时采集任务
如果您配置的是10分钟，则每10分钟触发同步一次。为了快速触发采集，您可以手动触发任务运行

登录后台：{your domain_name}/admin/，找到：PERIODIC TASKS -> Periodic tasks 

找到您刚才添加的任务名，点击前面的选择框、然后点击：动作 选择：Run Selected tasks，点击：执行 


#### 查看同步结果
定时任务运行后，会自动采集相关信息存放到表里面，在后台可以查看

登录后台：{your domain_name}/admin/，找到：SQLQUERY -> DB查询表  如果里面存在您目标库的库表信息，则表明同步成功

如果没有找到，可以通过日志看看是否存在网络不通、密码错误等问题。日志：logs/celery.log
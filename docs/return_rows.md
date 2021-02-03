# 配置返回行数
如果查询不指定limit 子句，这是一个很危险的操作，很容易搞死DB实例和YaSQL

### 配置
`vim yasql/config.py`

```python
# 查询LIMIT限制
# 如果查询没有LIMIT子句，默认加上limit default_return_rows
# 如果查询LIMIT子句返回行数大于max_return_rows，则改写为：limit max_return_rows
QUERY_LIMIT = {
    'default_return_rows': 100,
    'max_return_rows': 2000
}
```

您可以根据自己的需求，修改default_return_rows和max_return_rows即可，记得使用supervisorctl重启yasql服务

### 验证
登录前台查询页面，{your domain_name}/sqlquery/query ，输入sql验证即可

例如：
被改写为：
```sql
SELECT * from auditlog_logentry
被改写为：
SELECT * from auditlog_logentry LIMIT 100

SELECT * from auditlog_logentry limit 10000
被改写为：
SELECT * from auditlog_logentry LIMIT 2000

SELECT * from auditlog_logentry limit 10, 10000
被改写为：
SELECT  * from auditlog_logentry  LIMIT 10, 2000

SELECT * from auditlog_logentry limit 10000 offset 10
被改写为：
SELECT  * from auditlog_logentry  LIMIT 2000 OFFSET 10
```

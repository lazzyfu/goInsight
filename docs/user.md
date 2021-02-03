- [设置会话过期时间](#设置会话过期时间)
- [修改admin密码](#修改admin密码)
- [用户权限](#用户权限)
  - [权限](#权限)
  - [角色](#角色)
  - [用户如何绑定权限](#用户如何绑定权限)

### 设置会话过期时间
编辑文件：yasql/yasql/settings.py
```python
JWT_AUTH = {
    # 设置token过期时间为12小时
    'JWT_EXPIRATION_DELTA': datetime.timedelta(seconds=12 * 60 * 60),
    'JWT_AUTH_HEADER_PREFIX': 'JWT',
    'JWT_ALLOW_REFRESH': True,
    'JWT_SECRET_KEY': None,
    'JWT_GET_USER_SECRET_KEY': 'users.utils.jwt_get_user_secret',  # 为每个用户动态生成加密key
    'JWT_REFRESH_EXPIRATION_DELTA': datetime.timedelta(days=1),
}
```

修改：JWT_EXPIRATION_DELTA部分，然后重启服务(supervisorctl restart yasql-server)

### 修改admin密码
1. 登录后台http://xxx.xxx.xxx/admin
2. 输入系统默认的账号密码：admin / 1234.com
3. 依次选择：用户表 -> admin -> 密码 ，输入框填入明文密码，点击保存即可，系统会自动获取明文密码进行加密

### 用户权限
#### 权限
> 后台->权限表

> 权限不可删除，系统需要使用这些权限

* 删除上线版本权限
* 更新上线版本权限
* 创建上线版本权限
* 查看上线版本权限
* 执行工单权限
* 查看工单权限
* 审核工单权限
* 提交工单权限
* redis执行权限

#### 角色
> 后台->角色表
系统默认配置了3个角色（DBA、开发、测试），并绑定了权限

您也可以自己创建或删除现有的权限，并绑定不同的权限

#### 用户如何绑定权限
> 后台->用户表->用户[比如：lisi]->用户角色
给用户分配角色即可，用户绑定了指定的角色，就会继承当前角色绑定的权限

例如：用户lisi，角色为：DBA，那么lisi就用于DBA角色下面绑定的权限




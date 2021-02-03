# LDAP
YaSQL支持集成LDAP，可以使用ldap的账号进行登录系统

#### YaSQL支持2种用户认证方式
* 通过系统后台创建的本地用户
* 集成LDAP（直接对接企业账号）
  > 如果使用LDAP的认证形式，前台页面的修改密码功能将不起作用。需要您自己开发对接LDAP的密码修改的接口

#### 如何集成LDAP ？
> 每个公司的ldap目录结构不一样，请根据自己公司的LDAP目录结构进行配置

#### 示例
假设LDAP结构层级如下：
<img src="./pic/ldap.png" alt="" align=center />

##### OpenLDAP配置
> yasql/config.py
```python
# 启用LDAP
# LDAP配置如下，请按照自己公司的LDAP配置进行更正
LDAP_SUPPORT = {
    'enable': True,  # 为True启用LDAP，为False禁用LDAP
    'config': {
        'AUTH_LDAP_SERVER_URI': "ldap://192.168.3.221:389",       # ldap地址
        'AUTH_LDAP_BIND_DN': "uid=lisi,ou=people,dc=fzf,dc=com",  # 用于登录ldap的用户
        'AUTH_LDAP_BIND_PASSWORD': '1234.com',                    # uid=lisi的密码
        'AUTH_LDAP_USER_SEARCH': LDAPSearch(                      # ldap用户搜索目录，允许ou=people目录下的用户登录
            "ou=people,dc=fzf,dc=com", ldap.SCOPE_SUBTREE, "(uid=%(user)s)"
        ),
        # 用户属性映射(key:value) 
        # key为系统表字段，value为ldap字段(value部分对齐自己公司的ldap属性)
        'AUTH_LDAP_USER_ATTR_MAP': {
            'username': 'uid',
            'displayname': 'description',
            'email': 'email',
            'mobile': 'phone',
            'department': 'department'
        },
        'AUTH_LDAP_ALWAYS_UPDATE_USER': True,
        'AUTH_LDAP_START_TLS': False                            # 是否启用tls
    }
}
```

##### Windows活动目录配置
> yasql/config.py

> (uid=%(user)s) 需要改为 (sAMAccountName=%(user)s)

```python
# 启用LDAP
# LDAP配置如下，请按照自己公司的LDAP配置进行更正
LDAP_SUPPORT = {
    'enable': True,  # 为True启用LDAP，为False禁用LDAP
    'config': {
        'AUTH_LDAP_SERVER_URI': "ldap://192.168.3.221:389",       # ldap地址
        'AUTH_LDAP_BIND_DN': "uid=lisi,ou=people,dc=fzf,dc=com",  # 用于登录ldap的用户
        'AUTH_LDAP_BIND_PASSWORD': '1234.com',                    # uid=lisi的密码
        'AUTH_LDAP_USER_SEARCH': LDAPSearch(                      # ldap用户搜索目录，允许ou=people目录下的用户登录
            "ou=people,dc=fzf,dc=com", ldap.SCOPE_SUBTREE, "(sAMAccountName=%(user)s)"
        ),
        # 用户属性映射(key:value) 
        # key为系统表字段，value为ldap字段(value部分对齐自己公司的ldap属性)
        'AUTH_LDAP_USER_ATTR_MAP': {
            'username': 'uid',
            'displayname': 'description',
            'email': 'email',
            'mobile': 'phone',
            'department': 'department'
        },
        'AUTH_LDAP_ALWAYS_UPDATE_USER': True,
        'AUTH_LDAP_START_TLS': False                            # 是否启用tls
    }
}
```

#### 如何查看ldap日志
日志位置: yasql/logs/django.log

```text
'2020-11-10 14:08:39,143 [django_auth_ldap:370] [WARNING]- Caught LDAPError while authenticating zhangsan: INVALID_CREDENTIALS({'msgtype': 97, 'msgid': 1, 'result': 49, 'desc': 'Invalid credentials', 'ctrls': []})'
'2020-11-10 14:08:39,174 [django.channels.server:141] [INFO]- HTTP POST /api/users/login 200 [0.05, 127.0.0.1:56508]'
'2020-11-10 14:08:43,376 [django_auth_ldap:370] [WARNING]- Caught LDAPError while authenticating lisi: INVALID_CREDENTIALS({'msgtype': 97, 'msgid': 1, 'result': 49, 'desc': 'Invalid credentials', 'ctrls': []})'
'2020-11-10 14:08:43,474 [django.channels.server:141] [INFO]- HTTP POST /api/users/login 200 [0.10, 127.0.0.1:56515]'
'2020-11-10 14:08:45,067 [django_auth_ldap:370] [WARNING]- Caught LDAPError while authenticating lisi: INVALID_CREDENTIALS({'msgtype': 97, 'msgid': 1, 'result': 49, 'desc': 'Invalid credentials', 'ctrls': []})'
'2020-11-10 14:08:45,159 [django.channels.server:141] [INFO]- HTTP POST /api/users/login 200 [0.10, 127.0.0.1:56520]'
'2020-11-10 14:09:52,991 [django.utils.autoreload:597] [INFO]- Watching for file changes with StatReloader'
'2020-11-10 14:09:53,287 [daphne.server:111] [INFO]- HTTP/2 support not enabled (install the http2 and tls Twisted extras)'
'2020-11-10 14:09:53,287 [daphne.server:119] [INFO]- Configuring endpoint tcp:port=8000:interface=127.0.0.1'
'2020-11-10 14:09:53,288 [daphne.server:153] [INFO]- Listening on TCP address 127.0.0.1:8000'
'2020-11-10 14:09:53,368 [django.request:228] [WARNING]- Not Found: /'
'2020-11-10 14:09:53,370 [django.channels.server:149] [WARNING]- HTTP GET / 404 [0.04, 127.0.0.1:56576]'
'2020-11-10 14:10:04,595 [django_auth_ldap:655] [WARNING]- uid=lisi,ou=people,dc=fzf,dc=com does not have a value for the attribute description'
'2020-11-10 14:10:04,595 [django_auth_ldap:655] [WARNING]- uid=lisi,ou=people,dc=fzf,dc=com does not have a value for the attribute email'
'2020-11-10 14:10:04,595 [django_auth_ldap:655] [WARNING]- uid=lisi,ou=people,dc=fzf,dc=com does not have a value for the attribute phone'
'2020-11-10 14:10:04,595 [django_auth_ldap:655] [WARNING]- uid=lisi,ou=people,dc=fzf,dc=com does not have a value for the attribute department'
```

配置完成ldap后，记得重启Django服务。
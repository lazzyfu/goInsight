# AuditSQL

AuditSQL是基于Inception开发的一款web审核平台，旨在降低DBA的工作成本

为了提高审核效率和保证生产数据安全，数据库审核分解成如下2种模式：
- 线上审核

   开发 --> Inception审核SQL --> 生成工单 --> Leader审批 --> DBA手动执行 --> DBA反馈执行进度 --> 完成

- 线下审核

   开发 --> Inception审核SQL --> 生成执行任务 --> 执行任务 --> 实时显示进度 --> 查看结果


## 部署
拉取docker镜像：
```bash
docker pull lazzyfu/auditsql
```

启动docker：
```bash
docker run -itd -p 80:8000 --name=auditsql 459ad0efb89d /bin/bash
docker exec -it 459ad0efb89d /bin/bash
```

修改域名：

vim /etc/nginx/conf.d/nginx.conf
```bash
# 改成自己的域名
# 需要做域名解析或者自己本地hosts文件绑定宿主机IP
server_name sqlaudit.public.jbh.com;
```

系统配置：

vim /data/web/AuditSQL/AuditSQL/settings.py

修改下面部分
```python
#  如果不启用ldap登陆支持，注释下面
AUTHENTICATION_BACKENDS = [
    # 'django_auth_ldap.backend.LDAPBackend',
    'django.contrib.auth.backends.ModelBackend',
]

# 邮箱账号
## 修改成自己公司的邮箱账户
## 该账户用于实时发送审核邮件
EMAIL_HOST = 'smtp.163.com'
EMAIL_PORT = 25
EMAIL_HOST_USER = 'lazzyfu'
EMAIL_HOST_PASSWORD = '123.com'
EMAIL_USE_TLS = False
EMAIL_FROM = 'lazzyfu@163.com'

# LDAP配置
## 若不使用LDAP进行认证，注释上面的'django_auth_ldap.backend.LDAPBackend'
## ldap 服务器地址
AUTH_LDAP_SERVER_URI = "ldap://XXX.NET"
# AUTH_LDAP_START_TLS = True
AUTH_LDAP_ALWAYS_UPDATE_USER = True
## ldap绑定用户名，用户所在的绝对路径
AUTH_LDAP_BIND_DN = "CN=lazzyfu,OU=xx,DC=xx,DC=xx"
## ldap绑定用户名的密码
AUTH_LDAP_BIND_PASSWORD = "123.com"
## 搜索
AUTH_LDAP_USER_SEARCH = LDAPSearch("OU=xx科技,OU=xxx集团,DC=XXX,DC=NET", ldap.SCOPE_SUBTREE, "(CN=%(user)s)")

# 下面字段必须存在
## username：用户名
## email：邮箱地址
## displayname：对应的昵称或中文名
## key：为数据库字段
## value：为ldap对应字段
AUTH_LDAP_USER_ATTR_MAP = {"username": "cn", 'email': 'mail', "displayname": 'displayName'}

# 若需要调试ldap，不注释下面代码
# logger = logging.getLogger('django_auth_ldap')
# logger.addHandler(logging.StreamHandler())
# logger.setLevel(logging.DEBUG)
```

开启服务(麻烦，但是方便排查问题)：

```bash
chown -R mysql:mysql /var/lib/mysql
service mysql start
service redis start
uwsgi --ini /etc/nginx/conf.d/AuditSQL_uwsgi.ini
cd /data/web/AuditSQL
nohup daphne -b 0.0.0.0 -p 8001 -v2 AuditSQL.asgi:application --access-log=/var/log/daphnei.log &
service nginx start
/etc/init.d/celeryd start
```

## 已知的问题
- 当执行线下任务时，使用OSC执行时，inception会自动产生一个僵尸进程，不知道为什么

  解决办法：隔段时间自己进入docker容器手动清理下僵尸进程，对目标数据库没任何影响

## 开发组件
- Python 3.6
- Django 2.0 
- django-celery
- django-channels
- AdminLTE

## 功能：
- 线上审核
    - 数据变更提交
    - 工单记录
       - 审批 （Leader角色有权限操作）
       - 执行（DBA角色有权限操作）
       - 关闭（告知DBA不执行） 
    - 工单详情（显示当前工单的详情记录）

- 线下审核
    - SQL审核
    - 生成执行任务
    - 执行任务
       - 显示当前任务SQL列表
       - 执行，实时显示当前执行的OSC任务进度
       - 停止，关闭当前正在执行的OSC任务
       - 结果，显示回滚SQL，执行的日志（来自inception输出）
       
- SQL审核
   - 流程化
   - SQL美化
   - SQL检测
   - 审核历史记录
   - 语法高亮
   - 注释识别

- 进度推送：
   - 线上审核的每一步E-Mail实时推送
   - 线下执行任务进度的实时推送
  
- 其他：
   - 支持LDAP认证登陆
   - 项目权限控制
   - 角色权限控制
   - 支持修改头像
   - 支持修改密码
  
- 扩展功能：
   - 支持数据库表结构变更自动E-Mail通知，并生成变更结果


## 权限控制
用户角色（默认的3个角色）：
- Leader ：具有线上工单的审批权限
- DBA：具有线上工单的执行权限
- Developer

项目组（需要自己创建）：
- 联系人和用户可以属于多个项目
- 只要属于该项目的用户和联系人，才具有该项目的查看权限，用户权限继承项目权限

## 后台功能
![后台列表页](https://github.com/lazzyfu/AuditSQL/blob/master/media/gif/houtai-1.png)

账号配置：

配置用户账户、联系人、角色、项目组

以用户账户配置举例，如图：

![用户账户列表](https://github.com/lazzyfu/AuditSQL/blob/master/media/gif/user-1.png)

![用户账户详情](https://github.com/lazzyfu/AuditSQL/blob/master/media/gif/user-2.png)

## 线上SQL审核提交流程(点击查看GIF动图）
![gif动图](https://github.com/lazzyfu/AuditSQL/blob/master/media/gif/2018-03-15%2009_31_03.gif)
![提交发送的邮件](https://github.com/lazzyfu/AuditSQL/blob/master/media/gif/mail-1.png)

## 线上SQL审批流程(点击查看GIF动图）
![审批流程](https://github.com/lazzyfu/AuditSQL/blob/master/media/gif/2.gif)

## 线下SQL执行任务流程(点击查看GIF动态)
![线下SQL流程](https://github.com/lazzyfu/AuditSQL/blob/master/media/gif/11.gif)

## 联系方式
   
   QQ: 1126227133
   
   E-mail: 1126227133@qq.com
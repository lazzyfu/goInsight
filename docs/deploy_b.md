- [安装Django项目依赖包](#安装django项目依赖包)
- [安装UWSGI和GUNICORN服务](#安装uwsgi和gunicorn服务)
- [收集Django静态文件](#收集django静态文件)
- [部署Redis服务](#部署redis服务)
- [配置数据库连接](#配置数据库连接)
- [关闭DEBUG](#关闭debug)
- [初始化库表结构](#初始化库表结构)
- [初始化数据](#初始化数据)
- [部署supervisor服务](#部署supervisor服务)
- [配置supervisor服务](#配置supervisor服务)
- [启动supervisor服务](#启动supervisor服务)
- [查看服务状态](#查看服务状态)
- [日志位置](#日志位置)

#### 安装Django项目依赖包
```bash
cd /data/www/yasql/yasql
/venvyasql/bin/pip3.7 install -r requirements.txt -i https://mirrors.aliyun.com/pypi/simple
```
>请选择一个干净的系统，本地不要有自己安装的mysql包，否则在安装mysql-client时报ln类的错误，当然您也可以ln解决。

#### 安装UWSGI和GUNICORN服务
```bash
/venvyasql/bin/pip3.7 install gunicorn -i https://mirrors.aliyun.com/pypi/simple
/venvyasql/bin/pip3.7 install uwsgi -i https://mirrors.aliyun.com/pypi/simple
```

#### 收集Django静态文件
```bash
cd /data/www/yasql/yasql
mkdir static
mkdir logs
/venvyasql/bin/python3.7 manage.py collectstatic
```

#### 部署Redis服务
> redis版本需要且必须大于5.0

> 否则会抛出：aioredis.errors.ReplyError: ERR unknown command 'BZPOPMIN'

> 你也可以使用远程redis，本地不需要部署，我为了省事，就部署在本地了

```bash
wget https://download.redis.io/releases/redis-5.0.9.tar.gz
tar -zxf redis-5.0.9.tar.gz
cd redis-5.0.9/
make && make install

#更改redis配置
vim redis.conf
----------------------------
# 配置改为如下，此处设置密码
requirepass 1234.com
----------------------------

# 启动redis服务
nohup ./src/redis-server redis.conf &
```

#### 配置数据库连接
vim /data/www/yasql/yasql/config.py

> 请按照config.py文件里面的提示配置MySQL和Redis的地址和账号

#### 关闭DEBUG
vim /data/www/yasql/yasql/config.py

```python
# 关闭debug，本地开发时打开
# 生产环境请务必改为：DEBUG_ENABLED = False
DEBUG_ENABLED = False
```

#### 初始化库表结构
> 您已配置好了MySQL和Redis的地址和账号，并能够访问

执行migrate生成表结构，该操作会连接到上面的MySQL数据库创建表结构

>如果报数据库无法连接，请检查config.py里面的MySQL配置是否正确

`/venvyasql/bin/python3.7 manage.py migrate`

#### 初始化数据
`/venvyasql/bin/python3.7 manage.py loaddata initial_data.json`

#### 部署supervisor服务
```bash
/usr/local/bin/pip3.7 install supervisor
/usr/local/bin/echo_supervisord_conf > /etc/supervisord.conf
```

```editorconfig
vim /etc/supervisord.conf
------------------------------
# 更改为
[include]
files = supervisord.d/*.conf
------------------------------
```

```bash
mkdir /etc/supervisord.d
cd  /etc/supervisord.d
```


#### 配置supervisor服务
vim /etc/supervisord.d/yasql.conf
> 如果nginx里面upstream配置的不是127.0.0.1，此处需要将127.0.0.1改掉
```editorconfig
[program:yasql-server]
user=www
autorestart=true
environment=DJANGO_SETTINGS_MODULE="yasql.settings"
directory=/data/www/yasql/yasql
command=/venvyasql/bin/python3 /venvyasql/bin/gunicorn -w 8 -t 650 -b 127.0.0.1:8000 yasql.wsgi:application
redirect_stderr=true
stdout_logfile=/data/www/yasql/yasql/logs/yasql-server.log

[program:yasql-daphne]
user=www
autorestart=true
environment=DJANGO_SETTINGS_MODULE="yasql.settings"
directory=/data/www/yasql/yasql
numprocs=1
command=/venvyasql/bin/daphne -b 127.0.0.1 -p 8001 --proxy-headers -v2 yasql.asgi:application
redirect_stderr=true
stdout_logfile=/data/www/yasql/yasql/logs/yasql-daphne.log

[program:yasql-celery-beat]
user=www
autorestart=true
environment=DJANGO_SETTINGS_MODULE="yasql.settings"
directory=/data/www/yasql/yasql
command=/venvyasql/bin/celery beat -A yasql
redirect_stderr=true
stdout_logfile=/data/www/yasql/yasql/logs/yasql-celery-beat.log

[program:yasql-celery-default]
user=www
environment=DJANGO_SETTINGS_MODULE="yasql.settings"
directory=/data/www/yasql/yasql
command=/venvyasql/bin/celery worker -A yasql -n localhost -Q default -l info --time-limit=86400 --concurrency=5
redirect_stderr=true
stdout_logfile=/data/www/yasql/yasql/logs/yasql-celery-default.log
numprocs=1
autostart=true
autorestart=true
startsecs=10
stopwaitsecs = 600
stopasgroup=true
priority=1000

[program:yasql-celery-dbtask]
user=www
environment=DJANGO_SETTINGS_MODULE="yasql.settings"
directory=/data/www/yasql/yasql
command=/venvyasql/bin/celery worker -A yasql -n localhost -Q dbtask -l info --time-limit=86400 --concurrency=20
redirect_stderr=true
stdout_logfile=/data/www/yasql/yasql/logs/yasql-celery-dbtask.log
numprocs=1
autostart=true
autorestart=true
startsecs=10
stopwaitsecs = 600
stopasgroup=true
priority=1000
```

#### 启动supervisor服务
```bash
cd /data/www/yasql/yasql
 
# 创建日志目录
mkdir /data/www/yasql/yasql/{logs,static}

# 权限
chown -R www. /data/www/yasql/

# 启动supervisor进程
/usr/local/bin/supervisord -c /etc/supervisord.conf
```

#### 查看服务状态
```bash
[root@localhost yasql]# /usr/local/bin/supervisorctl status
goInception                      RUNNING   pid 4790, uptime 0:11:14
yasql-celery-beat                RUNNING   pid 4791, uptime 0:11:14
yasql-celery-dbtask              RUNNING   pid 4794, uptime 0:11:14
yasql-celery-default             RUNNING   pid 4795, uptime 0:11:14
yasql-daphne                     RUNNING   pid 4792, uptime 0:11:14
yasql-server                     RUNNING   pid 4793, uptime 0:11:14
```

如果启动异常，可查看日志报错，一般是目录权限问题

#### 日志位置
/data/www/yasql/yasql/logs
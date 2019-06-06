## 系统环境
- 系统：centos7
- 推荐配置：4核心/8GB内存
- 防火墙：关闭selinux 和 firewalld

## 配置文件
/data/web/opsql/config/config.py

只需要修改此文件即可，不需要再修改settings.py文件

## 安装部署
**安装系统依赖包**
```bash
yum -y install epel-release
yum -y install net-tools vim lsof lrzsz bzip2-devel wget \
gcc gcc-c++ make automake unzip curl curl-devel perl perl-devel \
expat expat-devel zlib zlib-devel asciidoc xmlto gettext-devel \
openssl-devel openssl mlocate python-devel openldap-devel \
readline-devel git mysql-devel
```

**安装Python-3.6**
```bash
# 不指定安装位置，默认安装到/usr/local/bin目录下
wget https://www.python.org/ftp/python/3.6.8/Python-3.6.8.tgz
tar -zxf Python-3.6.8.tgz
./configure
make && make install
```
 
**激活python虚拟环境**
```bash
/usr/local/bin/pip3.6 install --upgrade pip
/usr/local/bin/pip3.6 install virtualenv -i https://mirrors.aliyun.com/pypi/simple
/usr/local/bin/virtualenv /venv_py36 --python=/usr/local/bin/python3.6
echo "source /venv_py36/bin/activate" >> /root/.bashrc
source /root/.bashrc
```


**Clone项目代码**
```bash
mkdir /data/web -p
cd /data/web
# clone到本地，命名为opsql
git clone https://github.com/lazzyfu/AuditSQL.git opsql

```

**安装MySQL(可选)**
```bash
yum -y install https://repo.percona.com/yum/percona-release-latest.noarch.rpm
yum -y install Percona-Server-server-57 Percona-Server-devel-57

service mysql start
grep 'temporary password' /var/log/mysqld.log

# 修改密码和创建库（必须utf8，否则初始化失败）
alter user root@'localhost' identified by '123.com';
flush privileges;
```

**安装redis(可选)**
```bash
yum -y install redis
systemctl start redis.service  
```

**安装python依赖包**
```bash
cd /data/web/opsql/
pip install -r requirements.txt -i https://mirrors.aliyun.com/pypi/simple
```

**安装uwsgi**
```bash
pip install uwsgi
```

**安装nginx**
```bash
useradd nginx -s /bin/bash # 此处必须能登陆，celery服务需要使用
yum -y install nginx
chown -R nginx. /data/web
chown -R nginx. /venv_py36 # 必须设置虚拟环境的用户为nginx
```

**初始化数据**

需要修改配置文件(**/data/web/opsql/config/config.py**)中的数据库配置

```bash
cd /data/web/opsql

# 数据库创建库和用户，该用户必须要有with grant option权限
create database opsql character set utf8;
create user 'opsql'@'%' identified by '123.com';
grant all on *.* to 'opsql'@'10.10.1.201' with grant option;
flush privileges;

# 创建表结构
python manage.py migrate
# 导入数据
mysql -uopsql -p123.com opsql < documents/initial.sql
```

**处理静态文件**
```bash
python manage.py collectstatic
```

**编辑nginx配置文件**
```bash
vim /etc/nginx/conf.d/nginx.conf
### 内容如下 ###
server {
    listen      8000;
    server_name opsql.example.com;
    charset     utf-8;

    # max upload size
    client_max_body_size 75M;

    # django media directory
    location /media  {
        alias /data/web/opsql/media;
    }

    # django static directory
    location /static {
        alias /data/web/opsql/static;
    }

    # uwsgi 反向代理
    location / {
        uwsgi_pass  unix://///data/web/opsql_uwsgi.sock;
        uwsgi_read_timeout 600;

        # the uwsgi_params file you installed
        include     /etc/nginx/uwsgi_params;
    }

    # daphne 反向代理
    location /ws {
      proxy_pass http://0.0.0.0:8001;
      proxy_http_version 1.1;
      proxy_set_header Upgrade $http_upgrade;
      proxy_set_header Connection "upgrade";
      proxy_redirect     off;
      proxy_set_header   Host $host;
      proxy_set_header   X-Real-IP $remote_addr;
      proxy_set_header   X-Forwarded-For $proxy_add_x_forwarded_for;
      proxy_set_header   X-Forwarded-Host $server_name;
      proxy_read_timeout  36000s;
      proxy_send_timeout  36000s;
  }
}
```

**编辑uwsgi配置文件**
```bash
vim /etc/nginx/conf.d/opsql_uwsgi.ini
### 内容如下 ###
[uwsgi]
uid = nginx
chdir = /data/web/opsql
module = opsql.wsgi
home = /venv_py36
socket = /data/web/opsql_uwsgi.sock
processes = 8
master = true
max-requests = 6000
chmod-socket = 664
vacuum = true
enable-threads = true
single-interpreter = true
```

**安装supervisor服务**
```bash
pip install supervisor 
# 创建目录和生成配置文件
mkdir /etc/supervisor
mkdir /var/log/supervisord
echo_supervisord_conf > /etc/supervisor/supervisord.conf

vim /etc/supervisor/supervisord.conf
### 内容如下 ###
[program:uwsgi]
directory=/venv_py36/bin/
command=/venv_py36/bin/uwsgi --ini /etc/nginx/conf.d/opsql_uwsgi.ini
numprocs=1
user=root
startretries=3
startsecs=5
autostart=true
autorestart=true
stopsignal=INT
stopasgroup=true
killasgroup=true
redirect_stderr=true
stdout_logfile=/var/log/supervisord/uwsgi.log

[program:daphne]
directory=/data/web/opsql
command=/venv_py36/bin/daphne -b 0.0.0.0 -p 8001 --proxy-headers -v2 opsql.asgi:application
numprocs=1
user=root
startsecs=10
autostart=true
autorestart=true
redirect_stderr=true
stdout_logfile=/var/log/supervisord/daphne.log

[program:celery]
command=/venv_py36/bin/celery worker -A opsql --loglevel=INFO --time-limit=7200 --concurrency=10 --uid=nginx
directory=/data/web/opsql/
user=root
numprocs=1
stdout_logfile=/var/log/supervisord/celery_worker.log
stderr_logfile=/var/log/supervisord/celery_worker.log
autostart=true
autorestart=true
startsecs=10
stopwaitsecs = 600
stopasgroup=true
priority=1000

[program:celerybeat]
command=/venv_py36/bin/celery beat -A opsql --schedule /var/lib/celery/beat.db --loglevel=INFO --uid=nginx
directory=/data/web/opsql/
user=root
numprocs=1
stdout_logfile=/var/log/supervisord/beat.log
stderr_logfile=/var/log/supervisord/beat.log
autostart=true
autorestart=true
startsecs=10
stopasgroup=true
priority=999
```

**启动服务**
```bash
# 启动redis服务
systemctl start redis.service  

# 启动uwsgi服务
chown -R nginx. /data/web
supervisord -c /etc/supervisor/supervisord.conf 

supervisorctl status // 查看服务是否均已启动

# 启动nginx服务
systemctl restart nginx.service
```

**安装gh-ost(可以去github下载最新版本安装即可)**
```bash
rpm -ivh https://github.com/github/gh-ost/releases/download/v1.0.48/gh-ost-1.0.48-1.x86_64.rpm
```

**安装inception**

请自行安装，配置规则，然后修改config/config.py文件中的inception配置，重启uwsgi服务即可

inception的安装包位于: documents/inception-master.zip


**解决pymysql不兼容inception的问题**

若不解决，logs/all.log会输出下面错误
```
ValueError: invalid literal for int() with base 10: 'Inception2'
```

修改：
```bash
vim /venv_py36/lib/python3.6/site-packages/pymysql/connections.py
### 内容如下 ###
self.server_version = '5.7.18-16-log'  # 增加此行
if int(self.server_version.split('.', 1)[0]) >= 5:
    self.client_flag |= CLIENT.MULTI_RESULTS


vim /venv_py36/lib/python3.6/site-packages/pymysql/cursors.py
### 内容如下 ###
if not self._defer_warnings:
    #self._show_warnings()
    pass
```

**最后访问nginx里面配置的域名即可(http://opsql.example.com:8000)**

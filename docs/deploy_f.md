- [访问前台页面](#访问前台页面)

#### 最佳实践
强烈推荐您使用HTTPS的方式访问

HTTPS最佳实践：

浏览器 --> https://{your domain name} --> nginx+ssl --> nginx（yasql服务所在机器）

#### 安装Nginx
```bash
yum -y install nginx                    
useradd www -s /bin/bash
chown -R www. /data/www/
chown -R www. /venvyasql/
```

#### 编辑Nginx配置文件
vim /etc/nginx/conf.d/yasql.conf
> 请将下面的server_name替换为自己的域名

> nginx+ssl代理和nginx（yasql所在服务器）的机器的下面配置要一致

```editorconfig
server {
    listen      80;
    server_name yasql.examplexx.net;   # 此处更换为自己的域名
    charset     utf-8;

    root /data/www/yasql/yasql-fe/dist;
    index index.html;
    access_log   /var/log/nginx/yasql.log;

    location / {
       try_files $uri $uri/ /index.html;
    }

    location /api {
       proxy_read_timeout  650s;
       proxy_send_timeout  650s;
       proxy_connect_timeout 650s;
       proxy_pass http://127.0.0.1:8000;
    }

    location /api/media {
       expires 30d;
       alias /data/www/yasql/yasql/media;
    }

    location /admin {
       proxy_pass http://127.0.0.1:8000;
    }

    location /static {
       expires 30d;
       alias /data/www/yasql/yasql/static;
    }

    location /ws {
      proxy_pass http://127.0.0.1:8001;
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

#### 启动Nginx服务
```bash
systemctl enable nginx.service
systemctl start nginx.service
```

#### 访问前台页面
在浏览器访问：http://yasql.examplexx.net/ （此处应该是你在nginx里面配置的server_name）

>如果访问不了，本地先加下dns解析或者绑定下hosts

>如果nginx启动不了，检查下错误日志

>如果您基于https实现，请确保每台代理的nginx都配置如上面的location

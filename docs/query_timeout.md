#### 查询超时
单个查询执行时间默认为10分钟，即：600s。当执行时间超过10分钟，会被kill掉

检查nginx配置（如果您有多台nginx代理，则每台都需要配置）
> 确保存在proxy_read_timeout、proxy_send_timeout、proxy_connect_timeout
```nginx
  location /api {
       proxy_read_timeout  650s;     # 超时
       proxy_send_timeout  650s;     # 超时
       proxy_connect_timeout 650s;   # 超时
       proxy_pass http://127.0.0.1:8000;
    }
```

检查/etc/supervisord.d/yasql.conf配置
> 确保存在 -t 650

```bash
[program:yasql-server]
user=www
autorestart=true
environment=DJANGO_SETTINGS_MODULE="yasql.settings"
directory=/data/www/yasql/yasql
command=/venvyasql/bin/python3 /venvyasql/bin/gunicorn -w 8 -t 650 -b 127.0.0.1:8000 yasql.wsgi:application
redirect_stderr=true
stdout_logfile=/data/www/yasql/yasql/logs/yasql-server.log
```

由于默认是650s，如果您想继续调大，还需要修改前端文件
request.js

```javascript
// 创建 axios 实例
const request = axios.create({
  // API 请求的默认前缀
  baseURL: process.env.VUE_APP_API_BASE_URL,
  timeout: 650000 // 请求超时时间
})
```

#### 刷新页面
当您在执行查询的时候，如果触发了页面刷新或退出。当前的查询会被kill掉

#### 窗口自适应
支持前端查询界面窗口左右拉伸和上下拉伸


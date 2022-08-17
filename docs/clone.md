
# 两种方式
## 方式1: 克隆项目

从git上拉取项目代码，命令如下：

```bash
mkdir /data/www -p
cd /data/www
git clone https://github.com/lazzyfu/YaSQL.git yasql
```

## 方式2：下载压缩包
> 请选择最新的发行版本
```
wget https://github.com/lazzyfu/YaSQL/archive/refs/tags/v1.0.1.zip
```

## 目录结构说明
```bash
yasql                              # 项目根目录
├── doc                    
├── LICENSE
├── README.md
├── yasql                          # 后端代码目录
│   ├── apps
│   ├── config.py                  # django系统配置文件，mysql/redis/goInception/gh-ost/ldap/消息等
│   ├── __init__.py
│   ├── libs
│   ├── logs                       # 日志目录，需要自己手动创建
│   ├── static                     # 静态文件目录，需要自己手动创建
│   ├── logging.ini                # django日志配置
│   ├── manage.py
│   ├── media
│   ├── requirements.txt           # django依赖包
│   └── yasql
└── yasql-fe                       # 前端代码目录
    ├── babel.config.js
    ├── config
    ├── dist                       # 已编译发布的代码，Nginx指定的目录
    ...
```
# -*- coding:utf-8 -*-
# edit by fuzongfei

# 配置MySQL数据库，库必须先创建，且库的字符集必须为:utf8
# 存储django程序运行的系统库表等数据
# 权限：grant all on *.* to 'xxx'@'xxxhost' with grant options
DB = {
    'database': 'opsql',
    'user': 'root',
    'host': '127.0.0.1',
    'port': 3306,
    'password': '123.com',
}

# REDIS配置
# 存储session、消息队列等
REDIS = {
    'host': '127.0.0.1',
    'port': 6379
}

# 域名配置
# 用于邮件里面的提示
DOMAIN_NAME = {
    'value': 'http://127.0.0.1:8000'
}

# 邮箱配置
# 使用有效发送消息推送，可以和钉钉一起使用或单独使用
EMAIL = {
    'enable': True,
    'email_host': 'smtp.exmail.qq.com',
    'email_port': 465,
    'email_host_user': 'xxx@xxx.com',
    'email_host_password': '123.com',
    'email_use_ssl': True
}

# 钉钉通知
# 使用钉钉发送消息推送，可以和邮件一起使用或单独使用
DINGDING = {
    'enable': True,
    'webhook': "https://oapi.dingtalk.com/robot/send?access_token=74fc2cb89dea792"
               "ad276b336dec5e9fed0ee7484xxxxx"
}

# ALTER操作
# 是否启用Gh-ost工具改表，github地址：https://github.com/github/gh-ost
# Gh-ost工具对MySQL的ALTER操作有效
GHOST_TOOL = {
    'enable': True,
    # 每个参数后面必须跟上空格
    'arguments': "--allow-on-master --assume-rbr  "
                 "--initially-drop-ghost-table "
                 "--initially-drop-old-table  "
                 "-exact-rowcount "
                 "--approve-renamed-columns "
                 "--concurrent-rowcount=false "
                 "--hooks-path=/tmp/hook "
}

# 配置Inception主机地址
# 用于语法检查
# inception请自行安装
INCEPTION = {
    'inception_host': '10.10.1.202',
    'inception_port': 6033
}


# 查询限制
QUERY_LIMIT = {
    'enable': True,
    'default_limit': 100,
    'max_limit': 200
}

# SOAR测试数据库
SOAR_CONFIG = {
    # testenv需要指定
    # 该用户作为测试环境，需要all privileges权限
    # 请确保配置的测试环境数据库的版本大于生产环境的数据库版本
    "testenv": {
        "SOAR_HOST": "127.0.0.1",
        "SOAR_PORT": 3306,
        "SOAR_USER": "root",
        "SOAR_PASSWORD": "123.com"
    },
    "arguments": [
        # 额外的SOAR参数在下面逐行添加即可
        "-allow-online-as-test=false",  # 此参数不要修改，使用的是SOAR推荐的生产环境和测试环境的架构
        "-report-type=markdown",  # 此参数不要修改，前台对markdown格式进行了转码
        "-drop-test-temporary",
        "-max-join-table-count=5",
        "-max-group-by-cols-count=5",
        "-max-distinct-count=5",
        "-max-index-cols-count=5",
        "-only-syntax-check=false",
        "-allow-drop-index=false",
        "-explain=true",
        "-log-output=/tmp/soar.log"
    ]
}

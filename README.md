# goInsight
![](https://img.shields.io/static/v1?label=Go&message=1.20&color=green&?style=flat-square)
![](https://img.shields.io/static/v1?label=Vue&message=AntDesignVue&color=green&?style=flat-square)
![License](https://img.shields.io/github/license/lazzyfu/YaSQL?style=flat-square)

数据库工单审核和数据查询平台，旨在规范上线流程，解放DBA生产力。

## 语法审核工具
[gAudit](https://github.com/lazzyfu/gAudit)

## 功能
* 支持用户管理功能，2FA、定义任意层级的组织，并绑定用户
* 支持工单审批流，基于组织级别的权限控制；支持工单执行，工单执行记录
* 支持工单消息推送，支持钉钉/企业微信机器人/邮件推送
* 支持数据查询功能，支持用户定义表级别数据查询权限
* 支持语法审核，内置语法审核功能，默认集成gAudit接口
* 支持MySQL5.7+/TiDB Version 4+
* 支持提交DML/DDL工单
  * MySQL DDL自动集成gh-ost工具；
  * TiDB DDL ONLINE
  * 支持MySQL DML生成回滚SQL

## 限制
* 暂不支持TiDB DML语句生成回滚SQL

# 使用
## 下载
下载发行版本，解压到任意目录，并给可执行权限

## 编辑配置文件
按照配置文件备注调整即可，配置文件模版为：config.yaml.template

## 运行
您可以配置到supervisor或systemd中，也可以直接运行
```shell
nohup ./bin/goInsight -config config.yaml &
```

## 访问
访问http://ip:port
![](https://github.com/lazzyfu/goInsight/blob/master/docs/pics/login.png)
```
管理员账号密码，默认未开启2FA，登录后请及时修改密码。
默认管理员账号：admin
默认管理员密码：1234.Com!
```

- goInsight 后端代码
- goInsight-fe 前端代码

## 打包
### build前端
```
yarn build ---mode production --dest ../goInsight/dist
```
### build后端
```bash
GOOS=linux GOARCH=amd64 go build -o bin/goInsight -ldflags "-X main.version=1.0.0" main.go
```

## 启动服务
```
./bin/goInsight -config config.yaml
```

## 目录结构
```
goInsight
 - bin/goInsight 启动文件
 - config/config.yaml 配置文件
 - media 用户上传文件
 - logs 日志文件
```

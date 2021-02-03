# gh-ost
YaSQL集成了gh-ost工具，MySQL的Alter语句将会自动封装为gh-ost命令进行执行，可以支持在线热改。

> 由于pt-osc存在死锁的问题，因此本项目未集成

gh-ost地址：

https://github.com/github/gh-ost

#### 部署
```bash
wget https://github.com/github/gh-ost/releases/download/v1.1.0/gh-ost-1.1.0-1.x86_64.rpm
rpm -ivh gh-ost-1.1.0-1.x86_64.rpm
```

#### 配置gh-ost
> 通过下面的方式，您也可以自定义gh-ost的参数
vim yasql/config.py

```python
# gh-ost工具使用
GH_OST_ARGS = ['--allow-on-master',
               '--assume-rbr',
               '--initially-drop-ghost-table',
               '--initially-drop-old-table',
               '-exact-rowcount',
               '--approve-renamed-columns',
               '--concurrent-rowcount=false',
               '--chunk-size=800',
               '--hooks-path=/data/www/yasql/yasql/hook/']
```

记得重启django服务，否则参数将不会生效

`supervisorctl restart yasql-server`
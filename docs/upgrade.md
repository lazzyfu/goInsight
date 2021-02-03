## 备份配置文件
mv /data/www/yasql/yasql/config.py{,.bak}

## 拉取最新的代码
git pull origin master

## 将配置覆盖回去
> 请检查config.py是否有新增的配置，若有，请将配置添加到您自己的config.py里面

cp /data/www/yasql/yasql/config.py.bak /data/www/yasql/yasql/config.py

## 执行migrate
/venvyasql/bin/python3.7 manage.py migrate

## 重启后端服务
> 重启前，请确保没有在运行的gh-ost改表进程，否则改表将会失败，需要重新执行工单

supervisorctl restart all

## 重启nginx服务
chown -R www. /data/www/yasql/yasql-fe

systemctl restart nginx
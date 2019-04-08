# -*- coding:utf-8 -*-
# edit by fuzongfei
from celery import shared_task

from opsql.settings import DATABASES
from orders.models import MysqlConfig
from query.utils import SyncMysqlRemoteMeta


@shared_task
def periodic_sync_remote_schemameta_to_local():
    """
    定时任务
    同步sqlaudit_mysql_config表中只读数据库实例，在本地创建对应的库和表结构，不同步任何数据
    """
    remote_host = []
    # 排除django系统所在的库，避免采集数据的死循环
    for row in MysqlConfig.objects.filter(type=1).exclude(host=DATABASES.get('default').get('HOST')):
        remote_host.append({
            'id': row.id,
            'user': row.user,
            'password': row.password,
            'host': row.host,
            'port': row.port,
            'envi_id': row.envi_id,
            'comment': row.comment
        })
    op = SyncMysqlRemoteMeta(remote_host)
    op.run()


@shared_task
def mysql_privileges_operate(before_rule_id, after_rule_id, before_users, after_users):
    MySQLQueryRulesOperate(before_rule_id, after_rule_id, before_users, after_users).run()

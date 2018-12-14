# -*- coding:utf-8 -*-
# edit by fuzongfei
import logging

from celery import shared_task

from sqlaudit.settings import DATABASES
from sqlorders.models import MysqlConfig
from sqlquery.utils import MysqlQueryRemoteMetaOp, CreateLocalMysqlUser, MySQLQueryRulesOperate
from users.models import UserAccounts

logger = logging.getLogger(__name__)


@shared_task
def sync_remote_schemameta_to_local():
    """
    定时任务
    同步sqlaudit_mysql_config表中只读数据库实例，在本地创建对应的查询库和表结构
    """
    remote_host = []
    # 排除django系统所在的库，避免采集数据的死循环
    for row in MysqlConfig.objects.filter(is_type__in=(0, 2)).exclude(host=DATABASES.get('default').get('HOST')):
        remote_host.append({
            'id': row.id,
            'user': row.user,
            'password': row.password,
            'host': row.host,
            'port': row.port,
            'envi_id': row.envi_id,
            'comment': row.comment
        })
    aa = MysqlQueryRemoteMetaOp(remote_host)
    aa.run()


@shared_task
def create_local_mysql_user():
    """
    定时任务
    为每个后台用户创建一个本地的数据库账号
    此处的username必须为英文字符串，不能为中文
    """
    users = UserAccounts.objects.all().values_list('username', flat=True)
    aa = CreateLocalMysqlUser(users)
    aa.run()


@shared_task
def mysql_privileges_operate(before_rule_id, after_rule_id, before_users, after_users):
    MySQLQueryRulesOperate(before_rule_id, after_rule_id, before_users, after_users).run()

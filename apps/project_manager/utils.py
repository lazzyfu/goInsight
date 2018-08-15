# -*- coding:utf-8 -*-
# edit by fuzongfei
import json
import logging
import socket

import pymysql
from django.http import HttpResponse

from AuditSQL import settings
from mstats.models import MysqlSchemaInfo
from project_manager.inception.inception_api import sql_filter

logger = logging.getLogger(__name__)


def check_db_conn_status(host=None, port=None):
    """检测数据库是否可以连接"""
    port = int(port) if isinstance(port, str) else port
    obj = MysqlSchemaInfo.objects.filter(host=host, port=port).first()

    try:
        conn = pymysql.connect(user=obj.user,
                               host=obj.host,
                               password=obj.password,
                               port=obj.port,
                               use_unicode=True,
                               connect_timeout=1)

        if conn:
            return True, None
        conn.close()
    except pymysql.Error as err:
        logger.error(str(err))
        return False, str(err)


def check_incep_alive(fun):
    """检测inception进程是否运行"""

    def wapper(request, *args, **kwargs):
        inception_host = getattr(settings, 'INCEPTION_HOST')
        inception_port = getattr(settings, 'INCEPTION_PORT')

        sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        result = sock.connect_ex((inception_host, inception_port))

        if 0 == result:
            return fun(request, *args, **kwargs)
        else:
            context = {'status': 2, 'msg': 'Inception服务无法抵达，请联系管理员'}
            return HttpResponse(json.dumps(context))

    return wapper


def check_sql_filter(fun):
    """检查SQL类型，DML还是DDL操作"""

    def wapper(request, *args, **kwargs):
        sql = request.POST.get('contents')
        operate_type = request.POST.get('operate_type')

        filter_result = sql_filter(sql, operate_type)
        if filter_result.get('status') == 2:
            context = filter_result
            return HttpResponse(json.dumps(context))
        else:
            return fun(request, *args, **kwargs)

    return wapper

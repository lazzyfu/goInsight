# -*- coding:utf-8 -*-
# edit by fuzongfei
import logging
import re

import pymysql

from orders.models import MysqlConfig

logger = logging.getLogger('django')


def split_sqltype(sql, sql_type):
    # \s+ 匹配多个空字符，防止绕过
    ddl_filter = 'ALTER(\s+)TABLE|CREATE(\s+)TABLE|TRUNCATE(\s+)TABLE'
    dml_filter = 'INSERT(\s+)INTO|;UPDATE|^UPDATE|DELETE(\s+)FROM|\nUPDATE|\nDELETE|\nINSERT'

    if sql_type == 'DDL':
        if re.search(dml_filter, sql, re.I):
            msg = f'DDL模式下, 不支持SELECT|UPDATE|DELETE|INSERT语句'
            return False, msg
        else:
            return True, None
    if sql_type == 'DML':
        if re.search(ddl_filter, sql, re.I):
            msg = f'DML模式下, 不支持ALTER|CREATE|TRUNCATE语句'
            return False, msg
        else:
            return True, None


def checkdbstatus(host=None, port=None):
    # 检查数据库是否可以连接
    port = int(port) if isinstance(port, str) else port
    obj = MysqlConfig.objects.get(host=host, port=port)

    try:
        conn = pymysql.connect(user=obj.user,
                               password=obj.password,
                               host=host,
                               port=port,
                               use_unicode=True,
                               connect_timeout=1)

        if conn:
            return True, None
        conn.close()
    except pymysql.Error as err:
        logger.error(f'数据库连接检测错误，主机:{host}, 端口:{port}, 用户:{obj.user}')
        logger.error(str(err))
        return False, str(err)

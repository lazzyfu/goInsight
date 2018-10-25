# -*- coding:utf-8 -*-
# edit by fuzongfei
import json
import logging
import re
import socket

import pymysql
from django.http import HttpResponse

from sqlorders.models import MysqlSchemas
from sqlaudit import settings

logger = logging.getLogger('django')


def format_request(request):
    data = {}
    if request.method == 'GET':
        for key in request.GET.keys():
            values_list = ','.join(request.GET.getlist(key))
            data[key] = values_list if len(values_list) > 1 else values_list
    elif request.method == 'POST':
        for key in request.POST.keys():
            values_list = ','.join(request.POST.getlist(key))
            data[key] = values_list if len(values_list) > 1 else values_list
    return data


def check_db_conn_status(host=None, port=None):
    """检测数据库是否可以连接"""
    port = int(port) if isinstance(port, str) else port
    obj = MysqlSchemas.objects.filter(host=host, port=port).first()

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


class GetTableInfo(object):
    """获取指定主机的所有表"""

    def __init__(self, host, port, schema=None):
        # self.schema可以是单个库也可以是tuple
        self.host = host
        self.port = port
        self.schema = schema
        config = MysqlSchemas.objects.filter(host=self.host, port=self.port).first()
        self.conn = pymysql.connect(host=config.host,
                                    user=config.user,
                                    password=config.password,
                                    port=config.port,
                                    use_unicode=True,
                                    charset="utf8")

    IGNORED_PARAMS = ('information_schema', 'mysql', 'percona')

    def get_column_info(self):
        result = {}
        try:
            with self.conn.cursor() as cursor:
                tables_query = f"select TABLE_NAME,group_concat(COLUMN_NAME) from information_schema.COLUMNS " \
                               f"where table_schema not in {self.IGNORED_PARAMS} group by TABLE_NAME"
                cursor.execute(tables_query)
                tables = {}
                for table_name, column_name in cursor.fetchall():
                    tables[table_name] = list(column_name.split(','))

                result['tables'] = tables
        finally:
            self.conn.close()

        return result

    def get_online_tables(self):
        """
        返回格式：
        [{
        "id": 'test.tbl1',
        "icon": 'fa fa-table text-blue',
        "text": "tbl"
        }, ...]
        """
        result = []
        try:
            with self.conn.cursor() as cursor:
                query = f"select TABLE_NAME, concat_ws('.',TABLE_SCHEMA,TABLE_NAME) " \
                        f"from information_schema.COLUMNS where table_schema='{self.schema}' " \
                        f"group by TABLE_SCHEMA,TABLE_NAME"
                cursor.execute(query)
                for text, id in cursor.fetchall():
                    id = '___'.join((self.host, str(self.port), id))
                    result.append({'id': id,
                                   'text': text,
                                   "icon": 'fa fa-table text-blue'
                                   })
        finally:
            self.conn.close()
        return result

    def get_stru_info(self):
        """
        返回表结构和索引等信息
        """

        result = {}
        try:
            with self.conn.cursor() as cursor:
                stru_query = f"show create table {self.schema}"
                cursor.execute(stru_query)
                result['stru'] = cursor.fetchone()[1]

            self.conn.cursorclass = pymysql.cursors.DictCursor
            with self.conn.cursor() as cursor:
                try:
                    index_query = f"show index from {self.schema}"
                    # 获取字段
                    cursor.execute(index_query)
                    keys = cursor.fetchone().keys()
                    field = [{'field': j, 'title': j} for j in keys]

                    index_data = []
                    cursor.execute(index_query)
                    for i in cursor.fetchall():
                        index_data.append(i)

                    result['index'] = {'columnDefinition': field, 'data': index_data}
                except AttributeError as err:
                    result['index'] = {'columnDefinition': False, 'data': False}
        finally:
            self.conn.close()
        return result


# DDL和DML过滤
def sql_filter(sql, sql_type):
    # \s+ 匹配多个空字符，防止绕过
    ddl_filter = 'ALTER(\s+)TABLE|CREATE(\s+)TABLE|TRUNCATE(\s+)TABLE'
    dml_filter = 'INSERT(\s+)INTO|;UPDATE|^UPDATE|DELETE(\s+)FROM|\nUPDATE|\nDELETE|\nINSERT'

    if sql_type == 'DDL':
        if re.search(dml_filter, sql, re.I):
            context = {'status': 2, 'msg': f'DDL模式下, 不支持SELECT|UPDATE|DELETE|INSERT语句'}
        else:
            context = {'msg': '', 'status': 0, 'type': 'DDL'}
        return context

    elif sql_type == 'DML':
        if re.search(ddl_filter, sql, re.I):
            context = {'status': 2, 'msg': f'DML模式下, 不支持ALTER|CREATE|TRUNCATE语句'}
        else:
            context = {'msg': '', 'status': 0, 'type': 'DML'}
        return context


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
            context = {'status': 2, 'msg': '无法访问Inception服务无法，请联系管理员'}
            return HttpResponse(json.dumps(context))

    return wapper

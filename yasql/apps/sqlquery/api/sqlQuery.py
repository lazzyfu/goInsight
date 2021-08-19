# -*- coding:utf-8 -*-
# edit by fuzongfei
import datetime
import logging
import re
import time

import clickhouse_driver
import pymysql
import sqlparse
from asgiref.sync import async_to_sync
from channels.layers import get_channel_layer
from django_redis import get_redis_connection
from sqlparse.tokens import Whitespace, Keyword

from config import QUERY_LIMIT
from sqlquery.models import DbQueryLog

channel_layer = get_channel_layer()
logger = logging.getLogger('main')

NoneType = type(None)


def pull_msg(username=None, msg=None):
    """
    msg = {
    'type': 'query',
    'data': str(err)
    }
    """
    async_to_sync(channel_layer.group_send)(
        username,
        {
            "type": "user.message",
            'text': msg
        }
    )


class SqlQuery(object):
    """SQL查询接口"""

    def __init__(self, kwargs):
        """
        @param kwargs:
        {
        'config':
              {
                'host': '127.0.0.1',
                'password': 'xxx',
                'port': 3306,
                'user': 'yops_rw'
              },
         'rds_category': 1,
         'schema': 'broker',
         'tables: ['aa', 'bb'],
         'sql': 'select * from xxx',
         'username': 'zongfei.fu'
        }
        """
        self.kwargs = kwargs
        self.thread_id = None

    def _add_query_hash_to_redis(self):
        if self.thread_id:
            cnx_redis = get_redis_connection('default')
            value = ':'.join([str(self.thread_id), self.kwargs['config']['host'], str(self.kwargs['config']['port'])])
            cnx_redis.sadd(self.kwargs['query_hash'], value)

    def _del_query_hash_from_redis(self):
        if self.thread_id:
            cnx_redis = get_redis_connection('default')
            value = ':'.join([str(self.thread_id), self.kwargs['config']['host'], str(self.kwargs['config']['port'])])
            cnx_redis.srem(self.kwargs['query_hash'], value)

    def _remote_cnx(self):
        """连接到目标数据库"""
        config = self.kwargs.get('config')
        config['database'] = self.kwargs['schema']

        if self.kwargs['rds_category'] in [1, 2]:
            config.update(
                {
                    'max_allowed_packet': 1024 * 1024 * 1024,
                    'db': self.kwargs['schema'],
                    'read_timeout': 600,  # 设置最大查询时间600s
                    'cursorclass': pymysql.cursors.DictCursor
                }
            )
            # 先注释掉，兼容低版本，后续做下版本采集进行判断
            # if self.kwargs['rds_category'] in [1]:
            #     config['init_command'] = 'set session MAX_EXECUTION_TIME=600000'
            cnx = pymysql.connect(**config)
            return cnx
        if self.kwargs['rds_category'] in [3]:
            config.pop('charset')
            cnx = clickhouse_driver.connect(**config)
            return cnx

    def _operations_filter(self, sql):
        """检查语句的开头是否符合要求"""
        allowed_operations = [
            'SELECT',
            'SHOW',
            'DESC',
            'EXPLAIN'
        ]
        match_first_element = ''
        for i in sqlparse.parse(sql)[0].tokens:
            if i.ttype is not Whitespace:
                match_first_element = str(i).upper()
                break
        if match_first_element not in allowed_operations:
            return False, f"不被允许执行的语句：{match_first_element}"
        return True, sql

    def _limit_rules(self, sql):
        """检查SQL语句是否有LIMIT子句，并进行LIMIT限制"""
        default_return_rows = QUERY_LIMIT.get('default_return_rows')
        max_return_rows = QUERY_LIMIT.get('max_return_rows')

        # 从SQL语句中提取尾部的LIMIT子句
        # 不对SQL语句中的子查询的LIMIT做处理
        limit_sub = ""
        stmt = sqlparse.parse(sql)[0]
        seen = False
        # 匹配sql语句中的limit
        for token in stmt.tokens:
            if seen:
                limit_sub += token.value
            else:
                if token.ttype is Keyword and token.value.upper() == 'LIMIT':
                    seen = True
                    limit_sub += token.value.upper()

        rule_limit_n = re.compile(r'limit([\s]*\d+[\s]*)$', re.I | re.S)
        rule_limit_point = re.compile(r'limit([\s]*\d+[\s]*)(,)([\s]*\d+[\s]*)$', re.I | re.S)
        rule_limit_offset = re.compile(r'limit([\s]*\d+[\s]*)(offset)([\s]*\d+[\s]*)$', re.I | re.S)

        limit = re.compile(f'^SELECT(.*){limit_sub}', re.I | re.S)
        # SQL语句没有LIMIT子句
        if limit_sub == "":
            num = default_return_rows
            sql = limit.sub(
                r"SELECT \1 LIMIT {}".format(default_return_rows if int(num) >= max_return_rows else int(num)), sql)
            return sql
        # SQL语句中有LIMIT N子句
        if rule_limit_n.match(limit_sub):
            num = rule_limit_n.match(limit_sub).group(1)
            sql = limit.sub(
                r"SELECT \1 LIMIT {}".format(max_return_rows if int(num) >= max_return_rows else int(num)), sql)
            return sql
        # SQL语句中有LIMIT N, N子句
        if rule_limit_point.match(limit_sub):
            r = rule_limit_point.match(limit_sub)
            sql = limit.sub(
                r"SELECT \1 LIMIT {}, {}".format(
                    int(r.group(1)),
                    max_return_rows if int(r.group(3)) >= max_return_rows else int(r.group(3))
                ), sql)
            return sql
        # SQL语句中有LIMIT N OFFSET N子句
        if rule_limit_offset.match(limit_sub):
            r = rule_limit_offset.match(limit_sub)
            sql = limit.sub(
                r"SELECT \1 LIMIT {} OFFSET {}".format(
                    max_return_rows if int(r.group(1)) >= max_return_rows else int(r.group(1)),
                    int(r.group(3))
                ), sql)
            return sql

    def _escape_row(self, row):
        if isinstance(row, dict):
            for k in row:
                if isinstance(row[k], datetime.datetime):
                    # 时间类型转换为字符串，避免前端转时间的问题
                    row[k] = str(row[k])
                elif isinstance(row[k], datetime.timedelta):
                    row[k] = str(row[k])
                elif isinstance(row[k], NoneType):
                    row[k] = 'NULL'
                elif isinstance(row[k], int):
                    row[k] = str(row[k])
                elif isinstance(row[k], bytes):
                    # mysql列可能存在bit类型，转换成字符串和utf-8编码
                    row[k] = str(row[k], encoding='utf-8')
            return row
        if isinstance(row, tuple):
            new_row = []
            for i in row:
                if isinstance(i, datetime.datetime):
                    # 时间类型转换为字符串，避免前端转时间的问题
                    new_row.append(str(i))
                elif isinstance(i, datetime.timedelta):
                    new_row.append(str(i))
                elif isinstance(i, NoneType):
                    new_row.append('NULL')
                elif isinstance(i, int):
                    new_row.append(str(i))
                elif isinstance(i, bytes):
                    # mysql列可能存在bit类型，转换成字符串和utf-8编码
                    new_row.append(str(i, encoding='utf-8'))
                else:
                    new_row.append(i)
            return new_row

    def execute(self):
        msg = []
        status, _ = self._operations_filter(self.kwargs['sql'])
        if not status:
            return {'status': False, 'msg': _}
        # strip_comments防止下面SQL绕过limit
        # select * from noah_cmdb_instances;--
        rule_sql = self._limit_rules(sqlparse.format(_, strip_comments=True))

        # 记录执行SQL
        obj_record = DbQueryLog.objects.create(
            username=self.kwargs['username'],
            host=self.kwargs['config']['host'],
            schema=self.kwargs['schema'],
            tables=','.join(self.kwargs['tables']),
            query_sql=self.kwargs['sql']
        )

        # 建立连接
        try:
            cnx = self._remote_cnx()
            # 插入page_hash，设置过期时间为600s
            if self.kwargs['rds_category'] in [3]:
                self.thread_id = None
            else:
                # clickhouse拿不到thread_id
                self.thread_id = cnx.thread_id()
            self._add_query_hash_to_redis()

            with cnx.cursor() as cursor:
                start_time = time.time()
                cursor.execute(rule_sql)

                # 获取列名
                description = cursor.description
                fields = [{'field': 'state', 'checkbox': True}]
                fields.extend([{'field': x[0], 'title': x[0], 'escape': True} for x in description])

                # 获取数据
                data = []
                for row in cursor.fetchall():
                    escape_row = self._escape_row(row)
                    if isinstance(escape_row, list):
                        escape_row = dict(zip([i.name for i in description], escape_row))
                    data.append(escape_row)

                # 记录执行成功的语句，作为审计
                end_time = time.time()
                obj_record.affected_rows = cursor.rowcount
                obj_record.query_status = 'success'
                obj_record.query_consume_time = round(float(end_time - start_time), 3)
                obj_record.save()

                msg.append(f'{rule_sql}\n执行成功，影响行数：{obj_record.affected_rows}')
                msg.append(f'执行预估耗时：{obj_record.query_consume_time}')

                cursor.close()
            cnx.close()
            return {'status': True, 'data': {'columns': fields, 'data': data}}
        except Exception as err:
            err = err.args.__str__()[:1024]
            obj_record.query_status = err
            obj_record.save()
            msg.append(str(err))
            return {'status': False, 'msg': str(err)}
        finally:
            # 删除query_hash
            self._del_query_hash_from_redis()
            pull_msg(username=self.kwargs['username'], msg={'type': 'query', 'msg': '\n'.join(msg)})

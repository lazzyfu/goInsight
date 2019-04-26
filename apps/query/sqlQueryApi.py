# -*- coding:utf-8 -*-
# edit by fuzongfei

import datetime
import json
import re
import time

import pymysql
from asgiref.sync import async_to_sync
from channels.layers import get_channel_layer
from django_redis import get_redis_connection

from config.config import QUERY_LIMIT
from opsql.settings import DATABASES
from orders.models import MysqlConfig
from query.models import QueryBusinessGroup, MysqlUserGroupMap, MySQLQueryLog
from query.utils import LOCAL_QUERY_USER_PASSWORD

channel_layer = get_channel_layer()

NoneType = type(None)


class MySQLQueryApi(object):
    """MySQL查询接口"""

    def __init__(self, user=None, sqls=None, host=None, port=None, schema=None):
        self.user = user  # request.user.username
        self.sqls = sqls

        # 远程主机信息
        self.r_host = host
        self.r_port = port if isinstance(port, int) else int(port)
        self.r_schema = schema
        self.queryset = MysqlConfig.objects.get(host=self.r_host, port=self.r_port)

        # 本地主机信息
        self.local_host = DATABASES.get('default').get('HOST')
        self.local_port = DATABASES.get('default').get('PORT') if DATABASES.get('default').get('PORT') else 3306
        self.local_password = LOCAL_QUERY_USER_PASSWORD
        self.local_schema = '_'.join(['query', str(self.queryset.id), self.r_schema])

    def _local_cnx(self, map_user=None):
        """连接到本地数据库实例，用于SQL权限验证"""
        cnx = pymysql.connect(user=map_user,
                              password=self.local_password,
                              host=self.local_host,
                              port=self.local_port,
                              database=self.local_schema,
                              max_allowed_packet=1024 * 1024 * 1024,
                              charset='utf8',
                              cursorclass=pymysql.cursors.DictCursor)
        return cnx

    def _remote_cnx(self):
        """连接到远程数据库，执行SQL查询"""
        cnx = pymysql.connect(host=self.r_host,
                              port=self.r_port,
                              user=self.queryset.user,
                              password=self.queryset.password,
                              database=self.r_schema,
                              max_allowed_packet=1024 * 1024 * 1024,
                              charset='utf8',
                              cursorclass=pymysql.cursors.DictCursor)
        # 设置最大查询时间600s
        cnx._read_timeout = 600
        return cnx

    def cnx_redis(self):
        """连接到redis"""
        cnx = get_redis_connection('default')
        return cnx

    def _fmt(self, sqls):
        """格式化SQL语句，返回列表"""
        result = []
        for i in [re.sub('^\s+', '', i, re.S, re.I) for i in sqls.strip().split(';') if
                  re.sub('^\s+', '', i, re.S, re.I) != '']:
            # 匹配不以#开头的，此类为注释，不执行
            if re.search('^(?!#)', i, re.I):
                result.append(i)
        return result

    def _allowed(self, sqls):
        allowed_query = ['select', 'show', 'desc', 'explain']
        match_first_element = []
        for sql in sqls:
            match_first_element.append(sql.split(' ', 1)[0])

        # 转换为小写
        lower_match_first_element = [i.lower() for i in match_first_element]

        if not set(allowed_query) >= set(lower_match_first_element):
            no_support_r_query = list(set(lower_match_first_element).difference(set(allowed_query)))
            msg = '不支持如下SQL语句：{}'.format(','.join(no_support_r_query))
            return False, msg
        return True, sqls

    def _match(self, sqls):
        """对查询进行规则检测"""
        default_rows = 100
        max_rows = 200

        # 检查配置的值
        if QUERY_LIMIT['enable'] is True:
            default_rows = int(QUERY_LIMIT['default_limit'])
            max_rows = int(QUERY_LIMIT['max_limit'])

        # 对limit进行处理
        for sql in sqls:
            limit = re.compile('^SELECT(.*)LIMIT([\s]*)(.*)', re.I)
            no_limit = re.compile('^SELECT([\s\S]*)', re.I)

            if limit.match(sql) is None:
                sqls[sqls.index(sql)] = no_limit.sub(r"SELECT \1 LIMIT {}".format(default_rows), sql)
            else:
                limit_num = None
                value = limit.match(sql).group(3).upper()
                try:
                    limit_num = int(value)
                except ValueError as err:
                    if ',' in value:
                        limit_num, o = value.split(',')

                    if 'OFFSET' in value:
                        limit_num, o = value.split('OFFSET')
                sqls[sqls.index(sql)] = limit.sub(
                    r"SELECT \1 LIMIT {}".format(default_rows if int(limit_num) >= max_rows else int(limit_num)), sql)
        return sqls

    def get_map_mysql_user(self, sql):
        """获取映射的MySQL用户"""
        # 如果用户在对应业务组里面，且查询的表名属于该业务组，返回该业务组指定的映射的mysql用户
        obj = QueryBusinessGroup.objects.filter(user__username=self.user, config__id=self.queryset.pk)
        if obj.exists():
            tables = obj.values_list('tables', flat=True)
            for i in tables:
                for t in i.split(','):
                    if re.search(t.strip(), sql, re.I):
                        return obj.values_list('map_mysqluser__group', flat=True)[0]
        # 没有在业务组里面，查找所映射的用户
        queryset = MysqlUserGroupMap.objects.filter(comment__id=self.queryset.pk,
                                                    user__username=self.user,
                                                    schema__icontains=self.r_schema
                                                    ).values_list('group', flat=True)
        for i in queryset:
            if i.startswith('s_'):
                return i
            else:
                return queryset[0]

    def error_code_format(self, err):
        code, msg = err.args
        if code == 1142:
            r = re.compile('(.*) command denied (.*) for table (.*)', re.I)
            s = r.match(msg)
            return f"{s.group(1)} command denied for table {s.group(3)}\n" \
                f"Tips: 表不存在或您没有部分列的访问权限"

        elif code == 1143:
            r = re.compile('(.*) command denied (.*) for column (.*) in table (.*)', re.I)
            s = r.match(msg)
            return f"Tips: 您没有表{s.group(4)}中{s.group(3)}列的权限"
        else:
            return msg

    def check_permissions(self, sqls):
        # 检测sql是否有权限
        for sql in sqls:
            map_user = self.get_map_mysql_user(sql)
            cnx = self._local_cnx(map_user)
            try:
                with cnx.cursor() as cursor:
                    cursor.execute(sql)
            except Exception as err:
                return False, self.error_code_format(err)
            finally:
                cnx.close()
        return True, None

    def query(self, page_hash):
        fmt_sqls = self._fmt(self.sqls)
        s, d = self._allowed(fmt_sqls)
        if not s:
            push_msg = {'type': 1, 'msg': d}
            result = {'status': 2, 'msg': d}
        else:
            match_sqls = self._match(d)

            # 连接到远程数据库
            conn = self._remote_cnx()

            # 将page_hash写入到redis
            # pymysql查询的线程id作为值插入到redis的集合中
            cnx_redis = self.cnx_redis()
            thread_id = conn.thread_id()
            mi = '___'.join([str(thread_id),
                             self.queryset.host,
                             str(self.queryset.port),
                             self.queryset.user,
                             self.queryset.password])
            # 将页面hash作为key，pymysql查询的线程id作为值插入到redis的集合中
            cnx_redis.sadd(page_hash, mi)

            # 记录用户的查询
            querylog = MySQLQueryLog.objects.create(user=self.user,
                                                    host=self.r_host,
                                                    database=self.r_schema,
                                                    query_sql=self.sqls)

            status, msg = self.check_permissions(match_sqls)
            if status:
                try:
                    dynamic_table = {}
                    pull_msg = []
                    i = 1
                    for sql in match_sqls:
                        # 获取字段名
                        with conn.cursor() as cursor:
                            cursor.execute(sql)
                            keys = cursor.fetchone().keys()
                            field = [{'field': j, 'title': j} for j in keys]

                        # 获取数据
                        start_time = time.time()
                        with conn.cursor() as cursor:
                            cursor.execute(sql)
                            querylog.affect_rows = cursor.rowcount
                            querylog.query_status = '成功'
                            querylog.save()
                            pull_msg.append(f'{sql}\n执行成功，影响行数：{querylog.affect_rows}')
                            data = []
                            for j in cursor.fetchall():
                                for k in j:
                                    if isinstance(j[k], str):
                                        # 处理特殊字符，避免html进行转义
                                        v = j[k].replace('<', '&lt;').replace('>', '&gt;')
                                        j[k] = v.replace('\n', '<br>')
                                    elif isinstance(j[k], datetime.datetime):
                                        # 时间类型转换为字符串，避免前端转时间的问题
                                        j[k] = str(j[k])
                                    elif isinstance(j[k], datetime.timedelta):
                                        j[k] = str(j[k])
                                    elif isinstance(j[k], NoneType):
                                        j[k] = 'NULL'
                                    elif isinstance(j[k], int):
                                        j[k] = str(j[k])
                                    elif isinstance(j[k], bytes):
                                        # mysql列可能存在bit类型，转换成字符串和utf-8编码
                                        j[k] = str(j[k], encoding='utf-8')
                                data.append(j)
                            dynamic_table.update({f'{i}': {'columnDefinition': field, 'data': data}})
                            i += 1
                            end_time = time.time()
                            query_time = str(round(float(end_time - start_time), 3)) + 's'
                            querylog.query_time = query_time
                            querylog.save()
                            pull_msg.append(f'耗时：{query_time}\n')
                    push_msg = {'type': 1, 'msg': pull_msg}
                    result = {'status': 0, 'data': dynamic_table}
                except Exception as err:
                    if 'NoneType' in str(err):
                        msg = f'{sql[:30]} ...\n没有查询到记录'
                    else:
                        msg = str(err)
                    querylog.query_status = msg
                    querylog.save()
                    push_msg = {'type': 1, 'msg': msg}
                    result = {'status': 2, 'msg': msg}
                finally:
                    # 如果查询有返回或者查询抛出异常，从redis的集合里面移除本次查询的线程ID
                    cnx_redis.srem(page_hash, mi)
                    conn.close()
            else:
                push_msg = {'type': 1, 'msg': msg}
                result = {'status': 2, 'msg': msg}

        async_to_sync(channel_layer.group_send)(self.user, {"type": "user.message",
                                                            'text': json.dumps(push_msg)})
        return result

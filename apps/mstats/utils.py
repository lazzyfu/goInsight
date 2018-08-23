# -*- coding:utf-8 -*-
# edit by fuzongfei
import datetime
import json
import re
import time

import pymysql
from asgiref.sync import async_to_sync
from channels.layers import get_channel_layer

from mstats.models import MySQLQueryLog, MysqlSchemaInfo

channel_layer = get_channel_layer()


def get_mysql_user_info(host, port):
    data = MysqlSchemaInfo.objects.filter(host=host, port=port).first()
    conn = pymysql.connect(host=data.host,
                           user=data.user,
                           password=data.password,
                           charset='utf8',
                           port=data.port,
                           cursorclass=pymysql.cursors.DictCursor, )

    # 此方法可能受中间件影响，查询结果不准确
    # mysql_version = int(conn.server_version.split('.', 2)[1])
    version_query = 'select version();'

    try:
        with conn.cursor() as cursor:
            cursor.execute(version_query)
            mysql_version = int(cursor.fetchone().get('version()').split('.', 2)[1])

            id = 1

            user_query = "select user from mysql.user"
            cursor.execute(user_query)
            user_list = []

            for row in cursor.fetchall():
                user_list.append(row.get('user'))

            user_set = list(set(user_list))

            user_dict = []
            for i in user_set:
                user_dict.append({'id': id, 'pid': 0, 'privileges': '', 'schema': '', 'user': i})
                id += 1

            if mysql_version > 6:
                user_info_query = "select concat(\"'\",user,\"'\",'@',\"'\",host,\"'\") as username, " \
                                  "concat(left(authentication_string,5),'...',right(authentication_string,2)) " \
                                  "as password, password_expired from mysql.user"
            else:
                user_info_query = "select concat(\"'\",user,\"'\",'@',\"'\",host,\"'\") as username," \
                                  "concat(left(password,5),'...',right(password,2)) as password, " \
                                  "password_expired from mysql.user"
            cursor.execute(user_info_query)

            privileges_dict = []

            for row in cursor.fetchall():
                user = row.get('username')
                password = row.get('password')
                password_expired = row.get('password_expired')
                username = user.split('@')[0].replace("'", '').strip()
                user_host = user.split('@')[1].replace("'", '').strip()

                privileges_query = f"show grants for {user}"
                cursor.execute(privileges_query)
                pid = 0
                for i in cursor.fetchall():
                    for k, v in i.items():
                        data = re.split(r' TO ', v.replace('GRANT', '').strip())[0].split(r' ON ')
                        privileges = data[0]
                        schema = data[1]
                        for j in user_dict:
                            if username == j.get('user'):
                                pid = j.get('id')
                        privileges_dict.append({
                            'id': id,
                            'pid': pid,
                            'user': username,
                            'host': user_host,
                            'password': password,
                            'privileges': privileges,
                            'schema': schema,
                            'password_expired': password_expired
                        })
                        id += 1
            data = user_dict + privileges_dict
            return data
    finally:
        conn.close()


def mysql_rw_query(querys, rw='r'):
    """
    rw取值：'r' 和 'rw'
    """
    allowed_r_query = ['select', 'show', 'desc', 'explain']
    allowed_rw_query = ['select', 'show', 'desc', 'explain', 'update', 'delete', 'insert']

    match_first_element = []
    for sql in querys:
        match_first_element.append(sql.split(' ', 1)[0])

    # 转换为小写
    lower_match_first_element = [i.lower() for i in match_first_element]

    if rw == 'r':
        if not set(allowed_r_query) >= set(lower_match_first_element):
            no_support_r_query = list(set(lower_match_first_element).difference(set(allowed_r_query)))
            msg = '不支持如下SQL语句：{}'.format(','.join(no_support_r_query))
            return False, msg

    if rw == 'rw':
        if not set(allowed_rw_query) >= set(lower_match_first_element):
            no_support_rw_query = list(set(lower_match_first_element).difference(set(allowed_rw_query)))
            msg = '不支持如下SQL语句：{}'.format(','.join(no_support_rw_query))
            return False, msg

    return True, querys


def mysql_query_format(querys):
    """
    接收原始SQL
    格式化SQL语句，返回格式化后的SQL列表
    """

    format_querys = []

    # 匹配以\n开头和结尾且只包括\n的转换为''
    # 删除列表中的''元素
    for i in [re.sub('^\s+', '', i, re.S, re.I) for i in querys.strip().split(';') if
              re.sub('^\s+', '', i, re.S, re.I) != '']:
        # 多行匹配\n、\t、空格并替换为' '
        j = re.sub('\s+', ' ', i, re.S, re.I)
        # 匹配不以#开头的，此类为注释，不执行
        if re.search('^(?!#)', j, re.I):
            format_querys.append(j)

    return format_querys


def mysql_query_rules(querys):
    """
    对查询进行规则检测
    """

    # 判断是否有limit、没有增加limit 1000
    # 最大仅允许limit 1000
    for i in querys:
        limit = re.compile('^SELECT([\s\S]*) FROM ([\s\S]*) LIMIT (\d+)$', re.I)
        limit_offset = re.compile('^SELECT([\s\S]*) FROM ([\s\S]*) LIMIT (\d+) OFFSET (\d+)$', re.I)
        no_limit = re.compile('^SELECT([\s\S]*) FROM ([\s\S]*)', re.I)
        # select语句
        if re.match('^select', i, re.I):
            # 禁止limit N offset N语法
            if limit_offset.match(i) is None:
                if limit.match(i) is None:
                    # 当未匹配到select ... limit ...语句，重写查询
                    querys[querys.index(i)] = no_limit.sub(r"SELECT \1 FROM \2 LIMIT 500", i)
                else:
                    limit_num = limit.match(i)
                    if int(limit_num.group(3).replace(';', '')) > 500:
                        querys[querys.index(i)] = limit.sub(r"SELECT \1 FROM \2 LIMIT 1000", i)
            else:
                # 重新limit N offset N 为limit N语法
                querys[querys.index(i)] = limit_offset.sub(r'SELECT \1 FROM \2 LIMIT \3', i)

    print(querys)
    return querys


NoneType = type(None)


class MySQLQuery(object):
    """
    MySQL查询接口
    """

    def __init__(self, querys=None, host=None, port=None, schema=None, rw='r', envi=None):
        self.querys = querys
        self.host = host
        self.port = int(port)
        self.schema = schema
        self.envi = envi

        # 格式化SQL语句
        format_querys = mysql_query_format(self.querys)
        # 匹配查询规则，进行过滤
        limit_querys = mysql_query_rules(format_querys)
        # 判断是只读还是读写操作，依照环境而定
        self.status, self.data = mysql_rw_query(limit_querys, rw=rw)

        obj = MysqlSchemaInfo.objects.get(host=host, port=port, schema=schema)
        self.conn = pymysql.connect(host=self.host,
                                    port=self.port,
                                    user=obj.user,
                                    password=obj.password,
                                    charset='utf8',
                                    database=self.schema,
                                    cursorclass=pymysql.cursors.DictCursor)
        # 设置最大查询时间30s
        self.conn._read_timeout = 30

    def query(self, request):
        obj = MySQLQueryLog.objects.create(user=request.user.username,
                                           host=self.host,
                                           database=self.schema,
                                           envi=self.envi,
                                           query_sql=self.querys)
        if not self.status:
            obj.query_status = self.data
            obj.save()
            json_pull_data = {'type': 1, 'msg': self.data}
            result = {'status': 2, 'msg': self.data}
        else:
            try:
                start_time = time.time()
                dynamic_table = {}
                pull_msg = []
                i = 1
                for sql in self.data:
                    # 如果是DML语句中的update/insert/delete、执行并提交
                    # 此处统一将其转换为小写
                    first_element = sql.split(' ', 1)[0].lower()
                    if first_element in ('update', 'insert', 'delete'):
                        with self.conn.cursor() as cursor:
                            cursor.execute(sql)
                            obj.affect_rows = cursor.rowcount
                            obj.query_status = '成功'
                            obj.save()
                            pull_msg.append(f'{sql}\n执行成功，影响行数：{obj.affect_rows}\n')
                        self.conn.commit()
                    else:
                        # 非修改语句
                        # 获取字段
                        with self.conn.cursor() as cursor:
                            cursor.execute(sql)
                            keys = cursor.fetchone().keys()
                            field = [{'field': j, 'title': j} for j in keys]

                        # 获取数据
                        with self.conn.cursor() as cursor:
                            cursor.execute(sql)
                            obj.affect_rows = cursor.rowcount
                            obj.query_status = '成功'
                            obj.save()
                            pull_msg.append(f'{sql}\n执行成功，影响行数：{obj.affect_rows}\n')
                            data = []
                            for j in cursor.fetchall():
                                for k in j:
                                    if isinstance(j[k], str):
                                        # 处理特殊字符，避免html会进行转义
                                        v = j[k].replace('<', '&lt;').replace('>', '&gt;')
                                        j[k] = v.replace('\n', '<br>')
                                    elif isinstance(j[k], datetime.datetime):
                                        # 时间类型转换为字符串，避免前端转时间的问题
                                        j[k] = str(j[k])
                                    elif isinstance(j[k], datetime.timedelta):
                                        j[k] = str(j[k])
                                    elif isinstance(j[k], NoneType):
                                        j[k] = 'NULL'
                                    elif isinstance(j[k], bytes):
                                        # mysql列可能存在bit类型，转换成字符串和utf-8编码
                                        j[k] = str(j[k], encoding='utf-8')
                                data.append(j)
                        dynamic_table.update({f'{i}': {'columnDefinition': field, 'data': data}})
                        i += 1
                end_time = time.time()
                query_time = float(end_time - start_time)
                obj.query_time = round(query_time, 3)
                obj.save()
                json_pull_data = {'type': 1, 'msg': pull_msg}
                result = {'status': 0, 'data': dynamic_table}
            except Exception as err:
                if 'NoneType' in str(err):
                    error = f'{sql[:30]} ...\n没有查询到记录'
                else:
                    error = str(err)
                obj.query_status = error
                obj.save()
                json_pull_data = {'type': 1, 'msg': error}
                result = {'status': 2, 'msg': error}
            finally:
                self.conn.close()
        async_to_sync(channel_layer.group_send)(request.user.username, {"type": "user.message",
                                                                        'text': json.dumps(json_pull_data)})

        return result

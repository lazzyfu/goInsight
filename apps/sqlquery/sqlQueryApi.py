# -*- coding:utf-8 -*-
# edit by fuzongfei

import datetime
import json
import re
import time

import pymysql
from asgiref.sync import async_to_sync
from channels.layers import get_channel_layer

from sqlorders.models import MysqlSchemas, SysConfig, MysqlConfig
from sqlquery.models import MySQLQueryLog
from sqlaudit.settings import DATABASES

channel_layer = get_channel_layer()

NoneType = type(None)


class MySQLQuery(object):
    """
    MySQL查询接口
    """

    def __init__(self, user=None, querys=None, host=None, port=None, schema=None, rw='r'):
        # 用户连接本地数据库的账号
        # 本地数据库是指django系统库所在的数据库实例
        self.local_user = user
        self.local_password = 'LNjLJ6MeMJiZznL6'
        self.local_host = DATABASES.get('default').get('HOST')
        self.local_port = DATABASES.get('default').get('PORT') if DATABASES.get('default').get('PORT') else 3306

        # 目标数据库
        self.querys = querys
        self.remote_host = host
        self.remote_port = int(port)
        self.remote_schema = schema
        self.rw = rw

    def format_querys(self):
        """
            接收原始SQL
            格式化SQL语句，返回格式化后的SQL列表
            """

        result = []

        # 匹配以\n开头和结尾且只包括\n的转换为''
        # 删除列表中的''元素
        for i in [re.sub('^\s+', '', i, re.S, re.I) for i in self.querys.strip().split(';') if
                  re.sub('^\s+', '', i, re.S, re.I) != '']:
            # 多行匹配\n、\t、空格并替换为' '
            j = re.sub('\s+', ' ', i, re.S, re.I)
            # 匹配不以#开头的，此类为注释，不执行
            if re.search('^(?!#)', j, re.I):
                result.append(j)

        return result

    def filter_rulers(self, querys):
        """
            对查询进行规则检测
            """
        default_rows = 100
        max_rows = 200
        if SysConfig.objects.get(key='query_limit').is_enabled == '0':
            queryset = SysConfig.objects.get(key='query_limit').value
            a, b = queryset.split(',')
            default_rows = int(a.split('=')[1])
            max_rows = int(b.split('=')[1])
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
                        querys[querys.index(i)] = no_limit.sub(r"SELECT \1 FROM \2 LIMIT {}".format(default_rows), i)
                    else:
                        limit_num = limit.match(i)
                        if int(limit_num.group(3).replace(';', '')) > max_rows:
                            querys[querys.index(i)] = limit.sub(r"SELECT \1 FROM \2 LIMIT {}".format(max_rows), i)
                else:
                    # 重写limit N offset N 为limit N语法
                    limit_offset_match = limit_offset.match(i)
                    if int(limit_offset_match.group(3).replace(';', '')) > max_rows:
                        querys[querys.index(i)] = limit_offset.sub(r'SELECT \1 FROM \2 LIMIT {}'.format(max_rows), i)
                    else:
                        querys[querys.index(i)] = limit_offset.sub(r'SELECT \1 FROM \2 LIMIT \3', i)
        return querys

    def is_rw(self, querys, rw='r'):
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

    def _local_cnx(self):
        id = MysqlConfig.objects.get(host=self.remote_host, port=self.remote_port).id
        schema = '_'.join(['query', str(id), self.remote_schema])
        cnx = pymysql.connect(user=self.local_user,
                              password=self.local_password,
                              host=self.local_host,
                              port=self.local_port,
                              database=schema,
                              max_allowed_packet=1024 * 1024 * 1024,
                              charset='utf8',
                              cursorclass=pymysql.cursors.DictCursor)
        return cnx

    def _remote_cnx(self):
        obj = MysqlConfig.objects.get(host=self.remote_host, port=self.remote_port)
        cnx = pymysql.connect(host=self.remote_host,
                              port=self.remote_port,
                              user=obj.user,
                              password=obj.password,
                              charset='utf8',
                              database=self.remote_schema,
                              cursorclass=pymysql.cursors.DictCursor)
        # 设置最大查询时间600s
        cnx._read_timeout = 600
        return cnx

    def error_code_format(self, err):
        code, msg = err.args
        if code == 1142:
            r = re.compile('(.*) command denied (.*) for table (.*)', re.I)
            s = r.match(msg)
            return f"错误: {s.group(1)} command denied for table {s.group(3)}\n" \
                   f"提示：请检查是否有该表的{s.group(1)}权限"

        elif code == 1143:
            r = re.compile('(.*) command denied (.*) for column (.*) in table (.*)', re.I)
            s = r.match(msg)
            return f"错误: {s.group(1)} command denied for column {s.group(3)} in table {s.group(4)}\n" \
                   f"提示：您没有表{s.group(4)}中{s.group(3)}列的权限"
        else:
            return msg

    def check_permissions(self, querys):
        cnx = self._local_cnx()
        status = True
        msg = ''
        try:
            for sql in querys:
                with cnx.cursor() as cursor:
                    cursor.execute(sql)
        except Exception as err:
            status = False
            msg = self.error_code_format(err)
        finally:
            cnx.close()
            return status, msg

    def query(self):
        fquerys = self.format_querys()
        rquerys = self.filter_rulers(fquerys)
        status, after_querys = self.is_rw(rquerys, rw=self.rw)

        conn = self._remote_cnx()

        obj = MySQLQueryLog.objects.create(user=self.local_user,
                                           host=self.remote_host,
                                           database=self.remote_schema,
                                           query_sql=self.querys)
        if not status:
            obj.query_status = after_querys
            obj.save()
            json_pull_data = {'type': 1, 'msg': after_querys}
            result = {'status': 2, 'msg': after_querys}
        else:
            check_status, check_output = self.check_permissions(after_querys)
            if check_status:
                try:
                    dynamic_table = {}
                    pull_msg = []
                    i = 1
                    for sql in after_querys:
                        # 如果是DML语句中的update/insert/delete、执行并提交
                        # 此处统一将其转换为小写
                        first_element = sql.split(' ', 1)[0].lower()
                        if first_element in ('update', 'insert', 'delete'):
                            start_time = time.time()
                            with conn.cursor() as cursor:
                                cursor.execute(sql)
                                obj.affect_rows = cursor.rowcount
                                obj.query_status = '成功'
                                obj.save()
                                pull_msg.append(f'{sql}\n执行成功，影响行数：{obj.affect_rows}')
                            conn.commit()
                            end_time = time.time()
                            query_time = str(round(float(end_time - start_time), 3)) + 's'
                            obj.query_time = query_time
                            obj.save()
                            pull_msg.append(f'耗时：{query_time}\n')
                        else:
                            # 非修改语句
                            # 获取字段
                            with conn.cursor() as cursor:
                                cursor.execute(sql)
                                keys = cursor.fetchone().keys()
                                field = [{'field': j, 'title': j} for j in keys]

                            # 获取数据
                            start_time = time.time()
                            with conn.cursor() as cursor:
                                cursor.execute(sql)
                                obj.affect_rows = cursor.rowcount
                                obj.query_status = '成功'
                                obj.save()
                                pull_msg.append(f'{sql}\n执行成功，影响行数：{obj.affect_rows}')
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
                            obj.query_time = query_time
                            obj.save()
                            pull_msg.append(f'耗时：{query_time}\n')
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
                    conn.close()
            else:
                json_pull_data = {'type': 1, 'msg': check_output}
                result = {'status': 2, 'msg': check_output}
        async_to_sync(channel_layer.group_send)(self.local_user, {"type": "user.message",
                                                                  'text': json.dumps(json_pull_data)})

        return result

# -*- coding:utf-8 -*-
# edit by fuzongfei
import datetime
import json
import logging
import os
import re
import socket
import time
import zipfile

import pymysql
from asgiref.sync import async_to_sync
from channels.layers import get_channel_layer
from django.core.files import File
from django.http import HttpResponse
from openpyxl import Workbook
from openpyxl.cell.cell import ILLEGAL_CHARACTERS_RE

from sqlaudit import settings
from sqlorders.models import MysqlSchemas, MysqlConfig, SqlExportFiles, SqlOrdersExecTasks

logger = logging.getLogger('django')
channel_layer = get_channel_layer()


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


class ExportToExcel(object):
    """
    传入参数：user, sql, host, port, id
    """

    def __init__(self, id=None, user=None, sql=None, host=None, port=None, database=None):
        self.user = user
        self.sql = sql
        self.id = id
        self.encoding = 'gbk'
        self.type = 'excel'
        obj = MysqlConfig.objects.get(host=host, port=port)
        self.conn = pymysql.connect(host=obj.host,
                                    user=obj.user,
                                    password=obj.password,
                                    port=obj.port,
                                    db=database,
                                    max_allowed_packet=1024 * 1024 * 1024,
                                    charset='utf8')

        self.execute_log = []
        self.affected_row = 0

        # 文件名
        num = datetime.datetime.now().strftime("%Y%m%d%H%M%S")
        self.title = f'result_{num}'
        self.file = self.title + '.xlsx'
        self.zip_file = self.file + '.zip'

    def pull_msg(self, msg):
        # 推送消息
        msg = f"{msg} \n"
        pull_msg = {'status': 3, 'data': msg}
        async_to_sync(channel_layer.group_send)(self.user, {"type": "user.message",
                                                            'text': json.dumps(pull_msg)})

    def set_session_timeout(self):
        with self.conn.cursor() as cursor:
            cursor.execute("set session net_read_timeout=3600")

        with self.conn.cursor() as cursor:
            cursor.execute("set session net_write_timeout=3600")

    def get_count(self):
        # 查询当前SQL的返回的查询数量，返回分页的SQL列表
        status = True
        msg = None
        try:
            count_query = f"select count(*) as count from ({self.sql}) as subquery"
            self.conn.cursorclass = pymysql.cursors.DictCursor
            with self.conn.cursor() as cursor:
                cursor.execute(count_query)
                count = cursor.fetchone()
                self.affected_row = count['count']
            msg = f"SQL导出记录总数：{count['count']}"
        except Exception as err:
            msg = f"导出失败，发现错误：{str(err.args[1])}"
            status = False
        finally:
            self.execute_log.append(msg)
            self.pull_msg(msg)
            return status

    def compress_file(self):
        # 压缩文件
        msg = f'正在压缩文件：{self.file} ---> {self.zip_file}'
        self.pull_msg(msg)
        self.execute_log.append(msg)
        with zipfile.ZipFile(self.zip_file, 'w', allowZip64=True, compression=zipfile.ZIP_DEFLATED) as filezip:
            filezip.write(self.file)

        # 存储文件
        with open(self.zip_file, 'rb') as f:
            myfile = File(f)
            SqlExportFiles.objects.create(
                task_id=self.id,
                file_name=self.zip_file,
                file_size=os.path.getsize(self.zip_file),
                files=myfile,
                content_type='xlsx'
            )

        # 删除临时文件`
        msg = f'删除源文件：{self.file}'
        self.pull_msg(msg)
        self.execute_log.append(msg)
        os.remove(self.file) if os.path.exists(self.file) else None
        os.remove(self.zip_file) if os.path.exists(self.zip_file) else None

    def export_xlsx(self):
        # 导出成xlsx格式
        # num：保存文件的结尾_num标识，为str类型
        # 使用write_only能够有效降低内存的使用
        wb = Workbook(write_only=True)
        wb.encoding = f'{self.encoding}'
        ws = wb.create_sheet()
        ws.title = self.title

        # 获取列名作为标题
        self.conn.cursorclass = pymysql.cursors.DictCursor
        with self.conn.cursor() as cursor:
            cursor.execute(self.sql)
            title = []
            for column_name in cursor.fetchone():
                title.append(column_name)
        ws.append(title)

        # 获取数据，并写入到表格
        if self.affected_row <= 100000:
            # 当导出数据量小于10W时，使用fetchall直接读取到内存中
            self.conn.cursorclass = pymysql.cursors.Cursor
            with self.conn.cursor() as cursor:
                msg = f'正在导出SQL：{self.sql}'
                self.pull_msg(msg)
                self.execute_log.append(msg)

                cursor.execute(self.sql)
                rows = cursor.fetchall()

                msg = f'正在处理数据...'
                self.pull_msg(msg)
                self.execute_log.append(msg)

                for row in rows:
                    # 过滤掉特殊字符
                    filter_illegal_characters_row = list(
                        map(
                            (lambda x: ILLEGAL_CHARACTERS_RE.sub(r'', x) if isinstance(x, str) else x), row
                        )
                    )
                    ws.append(filter_illegal_characters_row)
            wb.save(self.file)
        else:
            # 当导出数据量大于10W时，使用SSCursor进行迭代读取，避免内存使用过大
            self.conn.cursorclass = pymysql.cursors.SSCursor
            with self.conn.cursor() as cursor:
                msg = f'正在导出SQL：{self.sql}'
                self.pull_msg(msg)
                self.execute_log.append(msg)

                cursor.execute(self.sql)
                while True:
                    row = cursor.fetchone()
                    if row:
                        # 过滤掉特殊字符
                        filter_illegal_characters_row = list(
                            map(
                                (lambda x: ILLEGAL_CHARACTERS_RE.sub(r'', x) if isinstance(x, str) else x), row
                            )
                        )
                        ws.append(filter_illegal_characters_row)
                    else:
                        break
            wb.save(self.file)
        self.compress_file()

    def run(self):
        queryset = SqlOrdersExecTasks.objects.get(id=self.id)
        status = self.get_count()
        if status:
            start_time = time.time()
            self.set_session_timeout()
            self.export_xlsx()
            end_time = time.time()
            consume_time = ''.join((str(round(end_time - start_time, 2)), 's'))
            msg = f'执行耗时：{consume_time}'
            self.execute_log.append(msg)
            self.pull_msg(msg)
            queryset.runtime = consume_time
            queryset.exec_status = '1'
        else:
            queryset.exec_status = '5'
        queryset.exec_log = '\n'.join(self.execute_log)
        queryset.save()

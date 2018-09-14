# -*- coding:utf-8 -*-
# edit by fuzongfei
import json
import logging

import pymysql
from asgiref.sync import async_to_sync
from channels.layers import get_channel_layer

from sqlorders.models import MysqlSchemas
from sqlaudit import settings

channel_layer = get_channel_layer()
logger = logging.getLogger('django')


class InceptionSqlApi(object):
    def __init__(self, host=None, port=None, database=None, contents=None, user=None):
        self.host = host
        self.port = port
        self.database = database
        self.sql_content = contents
        self.user = user

        self.inception_host = getattr(settings, 'INCEPTION_HOST')
        self.inception_port = int(getattr(settings, 'INCEPTION_PORT'))

        dst_server = MysqlSchemas.objects.get(host=self.host, port=self.port, schema=self.database)
        self.dst_host = dst_server.host
        self.dst_user = dst_server.user
        self.dst_password = dst_server.password
        self.dst_port = dst_server.port
        self.dst_database = self.database

    def conn_incep(self, sql):
        try:
            # 连接到inception
            conn = pymysql.connect(host=f"{self.inception_host}", user='root', password='', db='',
                                   port=self.inception_port, use_unicode=True, charset="utf8")
            cur = conn.cursor()
            cur.execute(sql)
            result = cur.fetchall()
            if result:
                field_names = [i[0] for i in cur.description]
                incep_data = []
                for row in result:
                    incep_data.append(dict(map(lambda x, y: [x, y], field_names, row)))
                cur.close()
                conn.close()
                return incep_data
        except pymysql.Error as err:
            raise EnvironmentError(err)

    def run_check(self):
        """对SQL进行审核"""
        sql = f"/*--user={self.dst_user};--password={self.dst_password};--host={self.dst_host};" \
              f"--enable-check=1;--port={self.dst_port};*/" \
              f"\ninception_magic_start;" \
              f"\nuse {self.dst_database};" \
              f"\n{self.sql_content}" \
              f"\ninception_magic_commit;"

        return {'status': 0, 'data': self.conn_incep(sql)}

    def run_exec(self, status, backup=None):
        """对SQL进行执行"""
        if backup == 'yes':
            sql = f"/*--user={self.dst_user};--password={self.dst_password};--host={self.dst_host};" \
                  f"--execute=1;--port={self.dst_port};*/" \
                  f"\ninception_magic_start;" \
                  f"\nuse {self.dst_database};" \
                  f"\n{self.sql_content}" \
                  f"\ninception_magic_commit;"
        else:
            sql = f"/*--user={self.dst_user};--password={self.dst_password};--host={self.dst_host};" \
                  f"--execute=1;--disable-remote-backup;--port={self.dst_port};*/" \
                  f"\ninception_magic_start;" \
                  f"\nuse {self.dst_database};" \
                  f"\n{self.sql_content}" \
                  f"\ninception_magic_commit;"

        exec_result = self.conn_incep(sql)
        pull_msg = {'status': status, 'data': exec_result}
        # 推送消息
        async_to_sync(channel_layer.group_send)(self.user, {"type": "user.message",
                                                            'text': json.dumps(pull_msg)})
        return exec_result

    def run_status(self, status):
        """执行inception命令"""
        sql = f"/*--user={self.dst_user};--password={self.dst_password};--host={self.dst_host};" \
              f"--execute=1;--port={self.dst_port};*/" \
              f"\n{self.sql_content}"
        exec_result = self.conn_incep(sql)
        pull_msg = {'status': status, 'data': exec_result}
        # 推送消息
        async_to_sync(channel_layer.group_send)(self.user, {"type": "user.message",
                                                            'text': json.dumps(pull_msg)})
        return exec_result

    def is_check_pass(self):
        """判断SQL是否通过审核"""
        check_data = self.run_check()['data']
        errlevel = [x['errlevel'] for x in check_data]
        if 1 in errlevel or 2 in errlevel:
            context = {'status': 2, 'msg': 'SQL语法检查未通过, 请执行语法检测'}
        else:
            context = {'status': 0, 'data': check_data}

        return context

    def make_sqlsha1(self):
        """
        将SQL切片成列表，分表进行审核并生成sqlsha1
        不可一起审核生成，否则执行时，两者SQL生成的sqlsha1不一致，无法获取进度
        """
        check_data = self.run_check()['data']
        sql_list = [x['SQL'] + ';' for x in check_data]
        result = []
        for sql in sql_list:
            self.sql_content = sql
            result.append(self.run_check()['data'][1])
        return result

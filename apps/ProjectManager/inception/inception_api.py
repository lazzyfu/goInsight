# -*- coding:utf-8 -*-
# edit by fuzongfei
import json
import logging
import re

import pymysql
from asgiref.sync import async_to_sync
from channels.layers import get_channel_layer

from AuditSQL import settings
from ProjectManager.models import InceptionHostConfig

channel_layer = get_channel_layer()
logger = logging.getLogger(__name__)


class GetBackupApi(object):
    """从备份主机上获取备份数据"""

    def __init__(self, data):
        self.inception_backup_host = getattr(settings, 'INCEPTION_BACKUP_HOST')
        self.inception_backup_user = getattr(settings, 'INCEPTION_BACKUP_USER')
        self.inception_backup_password = getattr(settings, 'INCEPTION_BACKUP_PASSWORD')
        self.inception_backup_port = getattr(settings, 'INCEPTION_BACKUP_PORT')

        self.backupdbName = data['backupdbName']
        self.sequence = data['sequence']

    def get_rollback_statement(self):
        conn = pymysql.connect(host=self.inception_backup_host, user=self.inception_backup_user,
                               password=self.inception_backup_password,
                               port=self.inception_backup_port, use_unicode=True, charset="utf8")

        cur = conn.cursor()

        rollback_statement = []

        if self.backupdbName != 'None':
            try:
                table_query = f"select tablename from {self.backupdbName}.$_$Inception_backup_information$_$ " \
                              f"where opid_time={self.sequence}"
                cur.execute(table_query)
                for row in cur.fetchall():
                    if row:
                        dst_table = row[0]

                        rollback_statement_query = f"select rollback_statement from {self.backupdbName}.{dst_table} " \
                                                   f"where opid_time={self.sequence}"
                        cur.execute(rollback_statement_query)

                        for i in cur.fetchall():
                            rollback_statement.append(i[0])

                        if rollback_statement:
                            return '\n'.join(rollback_statement)
                        else:
                            return False
                    else:
                        return False
            except conn.ProgrammingError as err:
                logger.warning(err)
                return False
            finally:
                cur.close()
                conn.close()


class GetDatabaseListApi(object):
    """获取目标主机的所有库"""

    def __init__(self, host):
        self.host = host

    IGNORED_PARAMS = ['information_schema', 'mysql', 'percona']

    def get_dbname(self):
        config = InceptionHostConfig.objects.get(host=self.host, is_enable=0)
        host = config.host
        user = config.user
        password = config.password
        port = config.port

        try:
            conn = pymysql.connect(host=host, user=user,
                                   password=password,
                                   port=port, use_unicode=True, charset="utf8")
            cur = conn.cursor()
            cur.execute("select schema_name from information_schema.schemata")
            db_list = []
            for i in cur.fetchall():
                db_list.append(i[0])

            for i in self.IGNORED_PARAMS:
                if i in db_list:
                    db_list.remove(i)

            cur.close()
            conn.close()
            return db_list
        except Exception as err:
            logger.warning(err)


# DDL和DML过滤
def sql_filter(sql_content, op_action):
    ddl_filter = 'ALTER TABLE|CREATE TABLE|TRUNCATE TABLE'
    dml_filter = 'INSERT INTO|;UPDATE|^UPDATE|DELETE FROM'

    if op_action == 'op_schema':
        if re.search(dml_filter, sql_content, re.I):
            context = {'status': 2, 'msg': f'DDL模式下, 不支持SELECT|UPDATE|DELETE|INSERT语句'}
        else:
            context = {'msg': '', 'status': 0, 'type': 'DDL'}
        return context

    elif op_action == 'op_data':
        if re.search(ddl_filter, sql_content, re.I):
            context = {'status': 2, 'msg': f'DML模式下, 不支持ALTER|CREATE|TRUNCATE语句'}
        else:
            context = {'msg': '', 'status': 0, 'type': 'DML'}
        return context


class IncepSqlCheck(object):
    def __init__(self, sql_content, host, database, user):
        self.sql_content = sql_content
        self.host = host
        self.database = database
        self.user = user
        self.inception_host = getattr(settings, 'INCEPTION_HOST')
        self.inception_port = int(getattr(settings, 'INCEPTION_PORT'))

        dst_server = InceptionHostConfig.objects.get(host=self.host, is_enable=0)
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

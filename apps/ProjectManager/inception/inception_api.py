# -*- coding:utf-8 -*-
# edit by fuzongfei
import json
import re

import pymysql
from asgiref.sync import async_to_sync
from channels.layers import get_channel_layer
from django_redis import cache

from AuditSQL import settings
from ProjectManager.models import InceptionHostConfig, IncepMakeExecTask
from ProjectManager.utils import update_tasks_status

channel_layer = get_channel_layer()


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
            dstTableQuery = f"select tablename from {self.backupdbName}.$_$Inception_backup_information$_$ where opid_time={self.sequence}"
            cur.execute(dstTableQuery)
            dstTable = cur.fetchone()[0]

            rollbackStatementQuery = f"select rollback_statement from {self.backupdbName}.{dstTable} where opid_time={self.sequence}"
            cur.execute(rollbackStatementQuery)

            for i in cur.fetchall():
                rollback_statement.append(i[0])
        cur.close()
        conn.close()

        if rollback_statement:
            return '\n'.join(rollback_statement)
        else:
            return '无记录'


class GetDatabaseListApi(object):
    """获取目标主机的所有库"""

    def __init__(self, host):
        self.host = host

    IGNORED_PARAMS = ['information_schema', 'mysql', 'percona']

    def get_dbname(self):
        master = InceptionHostConfig.objects.get(host=self.host, is_enable=0)
        masterHost = master.host
        masterUser = master.user
        masterPassword = master.password
        masterPort = master.port

        try:
            conn = pymysql.connect(host=masterHost, user=masterUser,
                                   password=masterPassword,
                                   port=masterPort, use_unicode=True, charset="utf8")
            cur = conn.cursor()
            dbQuery = "select schema_name from information_schema.schemata"
            cur.execute(dbQuery)
            dbList = []
            for i in cur.fetchall():
                dbList.append(i[0])

            for i in self.IGNORED_PARAMS:
                if i in dbList:
                    dbList.remove(i)

            cur.close()
            conn.close()
            return dbList
        except Exception as err:
            raise


# DDL和DML过滤
def sql_filter(sql_content, op_action):
    DDL_FILTER = 'ALTER TABLE|CREATE TABLE|TRUNCATE TABLE'
    DML_FILTER = 'INSERT INTO|;UPDATE|^UPDATE|DELETE FROM'

    if op_action == 'op_schema':
        if re.search(DML_FILTER, sql_content, re.I):
            context = {'errMsg': f'DDL模式下, 不支持SELECT|UPDATE|DELETE|INSERT语句', 'errCode': 400}
        else:
            context = {'errMsg': '', 'errCode': 200, 'type': 'DDL'}
        return context

    elif op_action == 'op_data':
        if re.search(DDL_FILTER, sql_content, re.I):
            context = {'errMsg': f'DML模式下, 不支持ALTER|CREATE|TRUNCATE语句', 'errCode': 400}
        else:
            context = {'errMsg': '', 'errCode': 200, 'type': 'DML'}
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
        sql = f"/*--user={self.dst_user};--password={self.dst_password};--host={self.dst_host};--enable-check=1;--port={self.dst_port};*/" \
              f"\ninception_magic_start;" \
              f"\nuse {self.dst_database};" \
              f"\n{self.sql_content}" \
              f"\ninception_magic_commit;"

        return {'errCode': 200, 'data': self.conn_incep(sql)}

    def run_exec(self, status):
        """对SQL进行执行"""
        sql = f"/*--user={self.dst_user};--password={self.dst_password};--host={self.dst_host};--execute=1;--port={self.dst_port};*/" \
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
        sql = f"/*--user={self.dst_user};--password={self.dst_password};--host={self.dst_host};--execute=1;--port={self.dst_port};*/" \
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
            context = {'errMsg': 'SQL语法检查未通过, 请执行语法检测', 'errCode': 400}
        else:
            context = {'data': check_data, 'errCode': 200}

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

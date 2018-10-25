# -*- coding:utf-8 -*-
# edit by fuzongfei
import json
import logging
import sys
import threading
import time

import pymysql
import sqlparse
from asgiref.sync import async_to_sync
from channels.layers import get_channel_layer

from sqlorders.api.generalRollbackApi import ReadRemoteBinlog

channel_layer = get_channel_layer()
logger = logging.getLogger('django')


class CnxStatusCheckThread(threading.Thread):
    def __init__(self, username, cnx, watch_thread_id):
        """
        监控被执行的SQL，是否被阻塞，有锁信息等
        cnx: 新的mysql连接
        watch_cnx: 被监控的pymysql执行SQL语句时建立的连接
        """
        threading.Thread.__init__(self)
        self.username = username
        self.cnx = cnx
        self.watch_thread_id = watch_thread_id

    def run(self):
        # 每秒查询一次状态
        # 当检测到lock，且锁定超过10s后，回滚该条SQL
        check_cmd = f"select * from information_schema.processlist where ID={self.watch_thread_id}"

        while True:
            with self.cnx.cursor() as cursor:
                cursor.execute(check_cmd)
                processlist_info = cursor.fetchone()
                if processlist_info:
                    # 返回的数据格式
                    # {'ID': 5703, 'USER': 'yops', 'HOST': '10.10.1.25:63032', 'DB': 'aa',
                    # 'COMMAND': 'Sleep', 'TIME': 0, 'STATE': '', 'INFO': None, 'TIME_MS': 44,
                    # 'ROWS_SENT': 0, 'ROWS_EXAMINED': 0}
                    pull_msg = {'status': 1, 'data': processlist_info}
                    async_to_sync(channel_layer.group_send)(self.username, {"type": "user.message",
                                                                            'text': json.dumps(pull_msg)})
                    print(processlist_info)
                else:
                    return False
            time.sleep(1)


class ExecuteSql(object):
    def __init__(self, host=None, port=None,
                 user=None, password=None, username=None,
                 database=None, charset='utf8mb4'):
        # 接收消息的用户
        self.username = username

        # 初始化sql，接收传入的单条SQL
        self.sql = None

        # 数据库连接配置
        self.host = host
        self.port = port
        self.user = user
        self.password = password
        self.database = database
        self.charset = charset

    def _connect(self):
        """新建连接"""
        try:
            cnx = pymysql.connect(host=self.host, port=self.port,
                                  user=self.user, password=self.password,
                                  charset=self.charset, database=self.database,
                                  cursorclass=pymysql.cursors.DictCursor
                                  )
            return cnx
        except Exception as err:
            print(err)
            sys.exit(1)

    def _sql_parser(self):
        """返回是DML还是DDL"""
        res = sqlparse.parse(self.sql)
        syntax_type = res[0].token_first().ttype.__str__()
        if syntax_type == 'Token.Keyword.DDL':
            type = 'DDL'
        elif syntax_type == 'Token.Keyword.DML':
            type = 'DML'
        else:
            # 非DML和DDL语句，比如：use db
            type = None
        return type

    def _get_position(self, cnx):
        """
        返回pos, file
        """
        position_cmd = 'show master status'

        with cnx.cursor() as cursor:
            cursor.execute(position_cmd)
            r = cursor.fetchone()
            return r['Position'], r['File']

    def _op_ddl(self):
        # 操作DDL语句
        pass

    def _op_dml(self, cnx):
        # 操作DML语句
        # 事务执行前，获取start position和binlog file
        start_pos, binlog_file = self._get_position(cnx)

        # 启动监控线程
        # 监控被执行的SQL是否有锁等待
        watch_thread_id = cnx.thread_id()
        t_cnx = CnxStatusCheckThread(self.username, self._connect(), watch_thread_id)
        t_cnx.setDaemon(True)
        t_cnx.start()

        # 执行SQL
        # 每条DML语句为作为一个事务执行
        try:
            start_time = time.time()
            with cnx.cursor() as cursor:
                cursor.execute(self.sql)
                affected_rows = cursor.rowcount
            cnx.commit()
            end_time = time.time()
            execute_time = str(round(float(end_time - start_time), 3)) + 's'

            # 事务执行完成后，获取end position
            end_pos, _ = self._get_position(cnx)

            # 判断影响行数
            if affected_rows > 0:
                data = ReadRemoteBinlog(binlog_file=binlog_file,
                                        start_pos=start_pos,
                                        end_pos=end_pos,
                                        host=self.host,
                                        port=self.port,
                                        user=self.user,
                                        password=self.password,
                                        affected_rows=affected_rows)
                # 返回回滚语句的列表
                result = {'status': 'success', 'rollbacksql': data.run_by_rows(), 'affected_rows': affected_rows,
                          'execute_time': execute_time}
            else:
                result = {'status': 'success', 'rollbacksql': '', 'affected_rows': f'影响行数：{affected_rows}',
                          'execute_time': execute_time}
        except Exception as err:
            result = {'status': 'fail', 'msg': str(err)}
        return result

    def run_by_sql(self, sql):
        cnx = self._connect()

        # 传入单条SQL语句
        self.sql = sql
        type = self._sql_parser()
        if type == 'DML':
            result = self._op_dml(cnx)
        elif type == 'DDL':
            # result = self._op_ddl()
            pass
        else:
            result = {'status': 'warn', 'msg': f'非DML和DDL语句，执行失败'}
        self._close(cnx)
        return result

    def _close(self, cnx):
        cnx.close()

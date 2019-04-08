# -*- coding:utf-8 -*-
# edit by fuzongfei
import json
import logging
import os
import re
import subprocess
import threading
import time

import pymysql
import sqlparse
from asgiref.sync import async_to_sync
from channels.layers import get_channel_layer

from config.config import GHOST_TOOL
from orders.api.extractTable import extract_tables
from orders.api.generalRollbackApi import ReadRemoteBinlog

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
                else:
                    return False
            time.sleep(1)


class ExecuteSql(object):
    def __init__(self,
                 host=None,
                 port=None,
                 user=None,
                 password=None,
                 username=None,
                 database=None,
                 database_type=None,
                 charset=None):
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

        # 判断数据库的类型，用于gh-ost
        # 分为阿里rds和非阿里的数据库
        self.database_type = database_type

    def _connect(self):
        """新建连接"""
        cnx = pymysql.connect(host=self.host,
                              port=self.port,
                              user=self.user,
                              password=self.password,
                              charset=self.charset,
                              database=self.database,
                              cursorclass=pymysql.cursors.DictCursor,
                              sql_mode="STRICT_TRANS_TABLES",
                              )
        cnx.show_warnings()
        return cnx

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

    def _remove_comment(self, sql):
        # 将语句中的注释和SQL分离
        sql_split = {}
        for stmt in sqlparse.split(sql):
            statement = sqlparse.parse(stmt)[0]
            comment = statement.token_first()
            if isinstance(comment, sqlparse.sql.Comment):
                sql_split = {'comment': comment.value, 'statement': statement.value.replace(comment.value, '')}
            else:
                sql_split = {'comment': '', 'statement': statement.value}

        # 获取不包含注释的SQL语句
        return sql_split['statement']

    def _check_mysql_environment(self, cnx):
        """
        检测mysql是否开启了binlog以及GTID
        若未开启，不进行备份
        """
        check_cmd = ["show variables like 'log_bin'",
                     "show variables like 'binlog_format'",
                     "show variables like 'server_id'"
                     ]
        rr = []
        with cnx.cursor() as cursor:
            for i in check_cmd:
                cursor.execute(i)
                data = cursor.fetchone()
                if data['Variable_name'] == 'log_bin' and data['Value'] == 'OFF':
                    rr.append('Binlog没有启用')

                if data['Variable_name'] == 'binlog_format' and data['Value'] != 'ROW':
                    rr.append('Binlog格式不为ROW')

                if data['Value'] == '0':
                    rr.append('server_id没有配置')
        self._close(cnx)
        return rr

    def _get_position(self, cnx):
        """
        返回pos, file
        """
        position_cmd = 'show master status'

        with cnx.cursor() as cursor:
            cursor.execute(position_cmd)
            r = cursor.fetchone()
        return r['Position'], r['File']

    def _get_processlist(self, watch_thread_id):
        """启动获取processlist的线程"""
        t_cnx = CnxStatusCheckThread(self.username, self._connect(), watch_thread_id)
        t_cnx.setDaemon(True)
        t_cnx.start()

    def _pull_msg(self, status, data):
        """实时推送消息"""
        pull_msg = {'status': status, 'data': data}
        async_to_sync(channel_layer.group_send)(
            self.username, {"type": "user.message",
                            'text': json.dumps(pull_msg)
                            })

    def _execute_sql(self, cnx):
        """执行SQL语句"""
        start_time = time.time()
        with cnx.cursor() as cursor:
            cursor.execute(self.sql)
            affected_rows = cursor.rowcount
            thread_id = cnx.thread_id()
        cnx.commit()

        end_time = time.time()
        runtime = str(round(float(end_time - start_time), 3)) + 's'
        exec_log = f"状态: 执行成功\n" \
            f"影响行数：{affected_rows}\n" \
            f"执行耗时：{runtime}\n"
        return affected_rows, runtime, exec_log, thread_id

    def _ghost_tool(self):
        syntaxcompile = re.compile('^ALTER(\s+)TABLE(\s+)([\S]*)(\s+)(ADD|CHANGE|RENAME|MODIFY|DROP|CONVERT)([\s\S]*)',
                                   re.I)
        syntaxmatch = syntaxcompile.match(self.sql)

        if syntaxmatch is not None:
            start_time = time.time()
            # 由于gh-ost不支持反引号，会被解析成命令，因此此处替换掉
            table = syntaxmatch.group(3).replace('`', '')
            # 将schema.table进行处理，这种情况gh-ost不识别，只保留table
            if len(table.split('.')) > 1:
                table = table.split('.')[1]

            # 处理反引号和将双引号处理成单引号
            value = ' '.join((syntaxmatch.group(5), syntaxmatch.group(6))).replace('`', '').replace('"', '\'')

            # 获取用户配置的gh-ost参数
            user_args = GHOST_TOOL['arguments']
            ghost_cmd = None

            if self.database_type == 0:
                # 非阿里云rds
                ghost_cmd = f"gh-ost {user_args} " \
                    f"--user={self.user} --password=\"{self.password}\" --host={self.host} --port={self.port} " \
                    f"--database=\"{self.database}\" --table=\"{table}\" --alter=\"{value}\" --execute"

            if self.database_type == 1:
                # 阿里云rds
                ghost_cmd = f"gh-ost {user_args} " \
                    f"--user={self.user} --password=\"{self.password}\" --host={self.host} --port={self.port} " \
                    f"--assume-master-host={self.host}:{self.port} --aliyun-rds " \
                    f"--database=\"{self.database}\" --table=\"{table}\" --alter=\"{value}\" --execute"

            # 记录执行的命令
            logger.info(ghost_cmd)

            # 删除sock，如果存在的话
            sock = os.path.join('/tmp', f'gh-ost.{self.database}.{table}.sock')
            os.remove(sock) if os.path.exists(sock) else None

            # 执行gh-ost命令
            p = subprocess.Popen(ghost_cmd, shell=True, stdout=subprocess.PIPE, stderr=subprocess.STDOUT)

            # 执行日志
            exec_log = ''

            # 检测子进程是否退出
            while p.poll() is None:
                data = p.stdout.readline().decode('utf8')
                if data:
                    exec_log += data
                    self._pull_msg(2, data)

            end_time = time.time()
            runtime = str(round(float(end_time - start_time), 3)) + 's'

            if p.returncode == 0:
                result = {'status': 'success', 'rollbacksql': '', 'affected_rows': 0, 'exec_log': exec_log,
                          'runtime': runtime}
            else:
                result = {'status': 'fail', 'exec_log': str(exec_log)}
        else:
            self._pull_msg(2, f'未成功匹配到SQL：{self.sql}，请检查语法是否存在问题')
            result = {'status': 'fail', 'exec_log': f'未成功匹配到SQL：{self.sql}，请检查语法是否存在问题'}
        return result

    def _extract_tables(self):
        """获取sql语句中的表名"""
        return [i.name for i in extract_tables(self.sql)]

    def _op_ddl(self, cnx):
        # 操作DDL语句
        # 匹配CREATE/DROP/RENAME/TRUNCATE语句
        # 此类语句直接执行
        origcompile = re.compile('^(CREATE|DROP|RENAME|TRUNCATE)([\s\S]*)', re.I)
        origmatch = origcompile.match(self.sql)
        result = None
        if origmatch is not None:
            # 启动监控线程
            # 监控被执行的SQL是否有锁等待
            watch_thread_id = cnx.thread_id()
            self._get_processlist(watch_thread_id)

            try:
                # 执行SQL
                affected_rows, runtime, exec_log, _ = self._execute_sql(cnx)
                self._close(cnx)
                result = {'status': 'success', 'rollbacksql': '', 'runtime': runtime,
                          'affected_rows': f'{affected_rows}',
                          'exec_log': exec_log}
            except Exception as err:
                exec_log = f"状态: 执行失败\n" \
                    f"错误信息：{str(err)}\n"
                result = {'status': 'fail', 'exec_log': exec_log}

        # 匹配ALTER语句
        # 此类语句需要分情况处理，alter直接执行或者使用工具gh-ost执行
        altercompile = re.compile('^(ALTER)([\s\S]*)', re.I)
        altermatch = altercompile.match(self.sql)
        if altermatch is not None:
            if GHOST_TOOL['enable'] is True:
                # 使用gh-ost工具执行ALTER语句
                result = self._ghost_tool()
            else:
                try:
                    # 直接执行ALTER语句
                    affected_rows, runtime, exec_log, _ = self._execute_sql(cnx)
                    self._close(cnx)
                    result = {'status': 'success', 'rollbacksql': '', 'affected_rows': f'{affected_rows}',
                              'runtime': runtime,
                              'exec_log': exec_log}
                except Exception as err:
                    exec_log = f"状态: 执行失败\n" \
                        f"错误信息：{str(err)}\n"
                    result = {'status': 'fail', 'exec_log': exec_log}
        return result

    def _op_dml(self, cnx):
        # 检查mysql相关配置
        # 如果check_result列表长度大于0，说明不支持备份
        binlog_file = start_pos = end_pos = result = None
        check_result = self._check_mysql_environment(self._connect())

        # 操作DML语句
        # 事务执行前，获取start position和binlog file
        if len(check_result) == 0:
            start_pos, binlog_file = self._get_position(cnx)

        # 启动监控线程
        # 监控被执行的SQL是否有锁等待
        watch_thread_id = cnx.thread_id()
        self._get_processlist(watch_thread_id)

        # 执行SQL
        # 每条DML语句为作为一个事务执行
        try:
            # 执行SQL
            affected_rows, runtime, exec_log, thread_id = self._execute_sql(cnx)

            # 事务执行完成后，获取end position
            if len(check_result) == 0:
                end_pos, _ = self._get_position(cnx)
            self._close(cnx)

            # 判断影响行数
            # 大于0小于20W时进行备份
            if 0 < affected_rows < 200000:
                # 获取回滚的SQL
                if len(check_result) == 0:
                    self._pull_msg(3, '正在执行当前SQL的备份，这可能需要花费些时间...')
                    # 备份时，传入schema和tables的列表
                    # 只读取指定schema和tables的binlog
                    rb_schemas = [self.database]
                    rb_tables = self._extract_tables()

                    data = ReadRemoteBinlog(binlog_file=binlog_file,
                                            start_pos=start_pos,
                                            end_pos=end_pos,
                                            host=self.host,
                                            port=self.port,
                                            user=self.user,
                                            password=self.password,
                                            thread_id=thread_id,
                                            only_schema=rb_schemas,
                                            only_tables=rb_tables)

                    # 接收数据格式
                    r_data = data.run_by_rows()
                    if r_data['status'] == 'success':
                        rollbacksql = '\n\n'.join(r_data['data'])
                        result = {'status': 'success', 'rollbacksql': rollbacksql, 'affected_rows': affected_rows,
                                  'runtime': runtime, 'exec_log': exec_log}
                        self._pull_msg(3, '备份成功，请点击结果按钮查看回滚SQL语句...')

                    elif r_data['status'] == 'fail':
                        result = {'status': 'fail', 'exec_log': exec_log}
                        self._pull_msg(3, f'备份失败, 失败原因：{exec_log}\n')
                else:
                    check_result = ', '.join(check_result)
                    exec_log = f"状态: 执行成功，备份失败\n" \
                        f"错误信息：{check_result}"
                    result = {'status': 'success', 'rollbacksql': '', 'affected_rows': f'{affected_rows}',
                              'runtime': runtime, 'exec_log': exec_log}
            elif affected_rows > 200000:
                result = {'status': 'success', 'rollbacksql': '', 'affected_rows': f'{affected_rows}',
                          'runtime': runtime, 'exec_log': ''.join([exec_log, '更新超过20W行，不进行备份'])}
            elif affected_rows == 0:
                result = {'status': 'success', 'rollbacksql': '', 'affected_rows': f'{affected_rows}',
                          'runtime': runtime, 'exec_log': exec_log}
        except Exception as err:
            exec_log = f"状态: 执行失败\n" \
                f"错误信息：{str(err)}\n"
            result = {'status': 'fail', 'exec_log': exec_log}
        return result

    def run_by_sql(self, sql):
        try:
            cnx = self._connect()

            # 传入单条SQL语句
            # 处理注释
            self.sql = self._remove_comment(sql)

            # 判断传入SQL的类型，为DML还是DDL
            type = self._sql_parser()

            if type == 'DML':
                result = self._op_dml(cnx)
            elif type == 'DDL':
                result = self._op_ddl(cnx)
            else:
                # 匹配USE语法，单独处理
                usecompile = re.compile('^(USE)([\s\S]*)', re.I)
                usematch = usecompile.match(self.sql)

                if usematch is not None:
                    exec_log = f"状态: 执行成功\n"
                    result = {'status': 'success', 'rollbacksql': '', 'affected_rows': 0,
                              'runtime': '0.00s', 'exec_log': exec_log}
                else:
                    # 其他语法的SQL，不执行，直接返回警告
                    exec_log = f"状态: 警告\n" \
                        f"错误信息：非DML和DDL语句，执行失败\n"
                    result = {'status': 'warn', 'exec_log': exec_log}
        except Exception as err:
            self._pull_msg(3, str(err))
            result = {'status': 'fail', 'exec_log': err}
        return result

    def _close(self, cnx):
        cnx.close()

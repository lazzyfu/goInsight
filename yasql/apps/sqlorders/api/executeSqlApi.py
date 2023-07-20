# -*- coding:utf-8 -*-
# edit by xff
import os
import re
import subprocess
import threading
import time
from warnings import filterwarnings

import pymysql
from clickhouse_driver import dbapi
from asgiref.sync import async_to_sync
from celery.utils.log import get_task_logger
from channels.layers import get_channel_layer

from config import GH_OST_ARGS
from sqlorders.api.extractTables import extract_tables
from sqlorders.api.generateRollbacksql import ReadRemoteBinlog

filterwarnings('ignore', category=pymysql.Warning)
channel_layer = get_channel_layer()

logger = get_task_logger('celery.logger')


def pull_msg(task_id=None, msg=None):
    """
    msg = {
    'type': 'execute',   // 取值：execute/processlist/ghost
    'data': str(err)
    }
    """
    async_to_sync(channel_layer.group_send)(
        task_id,
        {
            "type": "user.message",
            'text': msg
        }
    )


class CheckCnxStatusThread(threading.Thread):
    def __init__(self, task_id, cnx, watch_thread_id):
        """
        监控被执行的SQL，是否被阻塞，有锁信息等
        cnx: 新的mysql连接
        watch_cnx: 被监控的pymysql执行SQL语句时建立的连接
        """
        threading.Thread.__init__(self)
        self.task_id = task_id
        self.cnx = cnx
        self.watch_thread_id = watch_thread_id

    def run(self):
        # 每1秒查询一次状态
        check_cmd = f"select * from information_schema.processlist where ID={self.watch_thread_id}"

        while True:
            with self.cnx.cursor() as cursor:
                cursor.execute(check_cmd)
                processlist_info = cursor.fetchone()
                if processlist_info:
                    # 返回的数据格式
                    # {'ID': 5703, 'USER': 'yasql_rw', 'HOST': '10.10.1.25:63032', 'DB': 'aa',
                    # 'COMMAND': 'Sleep', 'TIME': 0, 'STATE': '', 'INFO': None, 'TIME_MS': 44,
                    # 'ROWS_SENT': 0, 'ROWS_EXAMINED': 0}
                    msg = {'type': 'processlist', 'data': processlist_info}
                    pull_msg(task_id=self.task_id, msg=msg)
                else:
                    self.cnx.close()
                    return False
            time.sleep(1)


class ExecuteSQL(object):
    def __init__(self, config):
        """
        {'host': '127.0.0.1', 'port': 3306, 'charset': 'utf8', 'rds_type': 3, 'database': 'test', rds_category: 1,
        'user': 'yasql_rw', 'password': '123.com', 'task_id': '1e0695520bb640e2ab9dcb8258aeb937', 'sql_type': 'DML'}
        """
        self.config = config
        self.sql = None
        self.clickhouse = True if self.config['rds_category'] == 3 else False

    def _cnx(self):
        # 新建连接
        cfg = self.config.copy()
        del cfg['rds_type']
        del cfg['task_id']
        del cfg['sql_type']
        del cfg['rds_category']
        if not self.clickhouse:
            cfg['cursorclass'] = pymysql.cursors.DictCursor
            # 执行SQL设置为严格模式
            cfg['sql_mode'] = "ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES"
            cnx = pymysql.connect(**cfg)
            # 设置锁定等待超时，避免一直等待下去
            with cnx.cursor() as cursor:
                cursor.execute("set session lock_wait_timeout = 10")
        else:
            del cfg['charset']
            cnx = dbapi.connect(**cfg)
            cnx.thread_id = lambda: 1

        return cnx

    def _extract_tables(self):
        """获取sql语句中的表名"""
        return [i.name for i in extract_tables(self.sql)]

    def _check_is_enabled_binlog(self, cnx):
        """检测mysql是否开启了binlog,若未开启，无法生成备份"""
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
                    rr.append('you must have binary logs enabled')

                if data['Variable_name'] == 'binlog_format' and data['Value'] != 'ROW':
                    rr.append('this binary log format must be row')

                if data['Value'] == '0':
                    rr.append('you must specify server_id')
        return len(rr) == 0, rr

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
        t_cnx = CheckCnxStatusThread(self.config['task_id'], self._cnx(), watch_thread_id)
        t_cnx.setDaemon(True)
        t_cnx.start()

    def _get_tidb_thread_id(self, cnx):
        """获取tidb执行SQL的线程ID"""
        with cnx.cursor() as cursor:
            cursor.execute("SELECT CONNECTION_ID() AS value")
            thread_id = cursor.fetchone()['value']
        return thread_id

    def _execute_sql(self, cnx):
        """执行SQL语句"""
        start_time = time.time()
        with cnx.cursor() as cursor:
            cursor.execute(self.sql)
            affected_rows = cursor.rowcount
            thread_id = cnx.thread_id()
        cnx.commit()

        end_time = time.time()
        consuming_time = round(float(end_time - start_time), 3)  # 简单计算耗时
        execute_log = f"状态: SUCCESS\n" \
                      f"影响行数：{affected_rows}\n" \
                      f"执行耗时：{consuming_time}\n"
        return affected_rows, consuming_time, execute_log, thread_id

    def _get_rollbacksql(self, affected_rows, start_read_binlog_file, start_pos, end_pos, thread_id):
        result = {'status': 'success',
                  'rollbacksql': '',
                  'execute_log': '',
                  'affected_rows': affected_rows
                  }
        # 根据影响行数判断是否进行备份
        if affected_rows == 0:
            msg = {'type': 'execute', 'data': '当前SQL影响行数为0，无可备份的数据'}
            pull_msg(task_id=self.config['task_id'], msg=msg)
            return result
        if 0 < affected_rows < 100000:
            # 获取回滚的SQL
            msg = {'type': 'execute', 'data': '正在执行当前SQL的备份，这可能需要花费些时间...'}
            pull_msg(task_id=self.config['task_id'], msg=msg)
            # 备份时，传入schema和tables的列表
            # 只读取指定schema和tables的binlog
            rb_schemas = self.config['database']
            rb_tables = self._extract_tables()

            rb = ReadRemoteBinlog(binlog_file=start_read_binlog_file,
                                  start_pos=start_pos,
                                  end_pos=end_pos,
                                  host=self.config['host'],
                                  port=self.config['port'],
                                  user=self.config['user'],
                                  password=self.config['password'],
                                  thread_id=thread_id,
                                  only_schema=rb_schemas,
                                  only_tables=rb_tables)

            # 接收数据格式
            rb_data = rb.run_by_rows()
            if rb_data['status'] == 'success':
                result['rollbacksql'] = '\n\n'.join(rb_data['data'])
                msg = {'type': 'execute', 'data': '备份成功，请点击结果按钮查看回滚SQL语句.'}
                pull_msg(task_id=self.config['task_id'], msg=msg)
                return result

            # 执行失败
            result['status'] = 'fail'
            result['execute_log'] = rb_data['msg']
            msg = {'type': 'execute', 'data': '备份失败, 失败原因：{rb_data["msg"]}\n'}
            pull_msg(task_id=self.config['task_id'], msg=msg)
            return result
        if affected_rows > 100000:
            result['status'] = 'fail'
            result['execute_log'] = '更新超过10W行，不进行备份'
            msg = {'type': 'execute', 'data': '更新超过10W行，不进行备份'}
            pull_msg(task_id=self.config['task_id'], msg=msg)
            return result

    def _op_mysql_dml(self, cnx):
        # 启动监控线程，监控被执行的SQL当前的会话状态
        self._get_processlist(cnx.thread_id())

        # 每条DML语句为作为一个事务执行
        try:
            # 事务执行前，获取start position和binlog file
            start_pos, start_read_binlog_file = self._get_position(cnx)
            # 执行SQL
            affected_rows, consuming_time, execute_log, thread_id = self._execute_sql(cnx)
            # 事务执行完成后，获取end position
            end_pos, _ = self._get_position(cnx)
            result = self._get_rollbacksql(affected_rows, start_read_binlog_file, start_pos, end_pos, thread_id)
            result['execute_log'] = '\n\n'.join([execute_log, result['execute_log']])
            result['consuming_time'] = consuming_time
            return result
        except Exception as err:
            execute_log = f"状态: Fail\n" \
                          f"错误信息：{str(err)}\n"
            result = {'status': 'fail', 'execute_log': execute_log}
            return result

    def _op_tidb_dml(self, cnx):
        # TiDB 需要执行 SQL 获取 thread_id
        thread_id = self._get_tidb_thread_id(cnx)
        # 监控进程
        self._get_processlist(thread_id)
        # 执行SQL
        try:
            affected_rows, consuming_time, execute_log, thread_id = self._execute_sql(cnx)
            result = {'status': 'success',
                      'rollbacksql': '',
                      'execute_log': execute_log,
                      'affected_rows': affected_rows,
                      'consuming_time': consuming_time
                      }
            return result
        except Exception as err:
            execute_log = f"状态: Fail\n" \
                          f"错误信息：{str(err)}\n"
            result = {'status': 'fail', 'execute_log': execute_log}
            return result

    def _op_clickhouse_dml(self, cnx):
        return self._op_tidb_dml(cnx)

    def _ghost_tool(self):
        syntaxcompile = re.compile(
            r'^ALTER(\s+)TABLE(\s+)([\S]*)(\s+)(ADD|CHANGE|RENAME|MODIFY|DROP|CONVERT|ENGINE)([\s\S]*)',
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
            args = GH_OST_ARGS

            # 当rds为阿里云RDS时，需要额外加上参数
            if self.config['rds_type'] == 1:
                args.append(f"--assume-master-host={self.config['host']}:{self.config['port']} --aliyun-rds")
            ghost_cmd = f"gh-ost {' '.join(args)} " \
                        f"--user={self.config['user']} --password=\"{self.config['password']}\" " \
                        f"--host={self.config['host']} --port={self.config['port']} " \
                        f"--database=\"{self.config['database']}\" --table=\"{table}\" " \
                        f"--alter=\"{value}\" --execute"

            # 输出到日志的命令，替换掉密码
            ghost_cmd_debug = re.sub(r"--password=(\S+)", '--password=xxx', ghost_cmd)
            logger.info(ghost_cmd_debug)

            # 删除sock，如果存在的话
            sock = f"/tmp/gh-ost.{self.config['database']}.{table}.sock"
            os.remove(sock) if os.path.exists(sock) else None

            # 执行gh-ost命令
            p = subprocess.Popen(ghost_cmd, shell=True, stdout=subprocess.PIPE, stderr=subprocess.STDOUT)

            execute_log = ''
            # 检测子进程是否退出
            while p.poll() is None:
                data = p.stdout.readline().decode('utf8')
                if data:
                    execute_log += data
                    msg = {'type': 'ghost', 'data': data}
                    pull_msg(task_id=self.config['task_id'], msg=msg)

            # 当进程退出时，读取剩余的输出
            if p.stdout:
                data = p.stdout.read().decode('utf8')
                execute_log += data
                msg = {'type': 'ghost', 'data': data}
                pull_msg(task_id=self.config['task_id'], msg=msg)

            # 计算耗时
            end_time = time.time()
            consuming_time = round(float(end_time - start_time), 3)

            if p.returncode == 0:
                return {'status': 'success',
                        'rollbacksql': '',
                        'affected_rows': 0,
                        'execute_log': execute_log,
                        'consuming_time': consuming_time}
            return {'status': 'fail', 'execute_log': execute_log}
        # 未匹配规则
        data = f"未成功匹配到SQL：{self.sql}，请检查语法是否存在问题"
        msg = {'type': 'ghost', 'data': data}
        pull_msg(task_id=self.config['task_id'], msg=msg)
        return {'status': 'fail', 'execute_log': data}

    def _op_mysql_ddl(self, cnx):
        # MySQL匹配到下面的SQL直接执行
        othercompile = re.compile(
            r'^('
            r'CREATE\s+TABLE|CREATE\s+VIEW|'
            r'DROP\s+TABLE|DROP\s+VIEW|DROP\s+TRIGGER|DROP\s+INDEX|'
            r'RENAME\s+TABLE|'
            r'TRUNCATE\s+TABLE'
            r')([\s\S]*)',
            re.I)
        if othercompile.match(self.sql) is not None:
            # 启动监控线程，监控被执行的SQL当前的会话状态
            self._get_processlist(cnx.thread_id())
            try:
                # 执行SQL
                affected_rows, consuming_time, execute_log, _ = self._execute_sql(cnx)
                result = {'status': 'success',
                          'rollbacksql': '',
                          'consuming_time': consuming_time,
                          'execute_log': execute_log,
                          'affected_rows': affected_rows
                          }
                return result
            except Exception as err:
                execute_log = f"状态: Fail\n" \
                              f"错误信息：{str(err)}\n"
                result = {'status': 'fail', 'execute_log': execute_log}
                return result

        # MySQL匹配ALTER语句
        # 此类语句需要分情况处理, alter直接执行或者使用工具gh-ost执行
        altercompile = re.compile(r'^(ALTER\s+TABLE)([\s\S]*)', re.I)
        if altercompile.match(self.sql) is not None:
            result = self._ghost_tool()
            time.sleep(10)  # 此处休眠10s，避免gh-ost资源未释放
            return result

        # 不满足上面的条件
        return {'status': 'fail', 'execute_log': '未成功匹配到规则，不被允许执行的DDL，请进行核对'}

    def _op_tidb_ddl(self, cnx):
        """tidb直接连接数据库执行"""
        sqlcompile = re.compile(
            r'^('
            r'CREATE\s+TABLE|CREATE\s+SEQUENCE|CREATE\s+INDEX|CREATE\s+OR|CREATE\s+VIEW|'
            r'DROP\s+TABLE|DROP\s+VIEW|DROP\s+TRIGGER|DROP\s+INDEX|'
            r'RENAME\s+TABLE|'
            r'TRUNCATE\s+TABLE|'
            r'ALTER\s+TABLE'
            r')([\s\S]*)',
            re.I)
        if sqlcompile.match(self.sql) is not None:
            # TiDB 需要执行 SQL 获取 thread_id
            thread_id = self._get_tidb_thread_id(cnx)
            # 启动监控线程，监控被执行的SQL当前的会话状态
            self._get_processlist(thread_id)
            try:
                # 执行SQL
                affected_rows, consuming_time, execute_log, _ = self._execute_sql(cnx)
                result = {'status': 'success',
                          'rollbacksql': '',
                          'consuming_time': consuming_time,
                          'execute_log': execute_log,
                          'affected_rows': affected_rows
                          }
                return result
            except Exception as err:
                execute_log = f"状态: Fail\n" \
                              f"错误信息：{str(err)}\n"
                result = {'status': 'fail', 'execute_log': execute_log}
                return result
        else:
            logger.error("函数: _op_tidb_ddl 原因: SQL语句未匹配正则sqlcompile")
            execute_log = f"状态: Fail\n" \
                f"错误信息：函数: _op_tidb_ddl 原因: SQL语句未匹配正则sqlcompile\n"
            return {'status': 'fail', 'execute_log': execute_log}

    def _op_clickhouse_ddl(self, cnx):
        return self._op_clickhouse_dml(cnx)

    def check_read_only(self, cnx):
        # 执行SQL之前，检查数据库是否关闭只读, read_only = 0 为只读关闭，=1为只读打开
        # 如果只读为ON，则不执行，请确认是否连接的是主库
        check_cmd = "select @@read_only as value"
        with cnx.cursor() as cursor:
            cursor.execute(check_cmd)
            if cursor.fetchone()['value'] == 0:
                return True
        return False

    def run_by_sql(self, sql):
        result = None
        try:
            cnx = self._cnx()

            # 检查mysql环境,rds_category=1 为mysql
            if self.config['rds_category'] == 1:
                # 检查是否支持备份
                status, result = self._check_is_enabled_binlog(cnx)
                if not status:
                    execute_log = f"状态: Fail\n" \
                                  f"错误信息：{', '.join(result)}\n"
                    msg = {'type': 'execute', 'data': execute_log}
                    pull_msg(task_id=self.config['task_id'], msg=msg)
                    result = {'status': 'fail', 'execute_log': execute_log}
                    return result

                if not self.check_read_only(cnx):
                    msg = {'type': 'execute', 'data': '当前READ_ONLY = ON，执行失败，请确认是否连接的是主库'}
                    pull_msg(task_id=self.config['task_id'], msg=msg)
                    result = {'status': 'fail', 'execute_log': '当前READ_ONLY = ON，执行失败，请确认是否连接的是主库'}
                    return result

            self.sql = sql
            # 判断传入SQL的类型，为DML还是DDL
            if self.config['sql_type'] == 'DML':
                if self.config['rds_category'] == 1:
                    result = self._op_mysql_dml(cnx)
                if self.config['rds_category'] == 2:
                    result = self._op_tidb_dml(cnx)
                if self.config['rds_category'] == 3:
                    result = self._op_clickhouse_dml(cnx)
            if self.config['sql_type'] == 'DDL':
                if self.config['rds_category'] == 1:
                    result = self._op_mysql_ddl(cnx)
                if self.config['rds_category'] == 2:
                    result = self._op_tidb_ddl(cnx)
                if self.config['rds_category'] == 3:
                    result = self._op_clickhouse_ddl(cnx)
            cnx.close()
        except Exception as err:
            msg = {'type': 'execute', 'data': str(err)}
            pull_msg(task_id=self.config['task_id'], msg=msg)
            result = {'status': 'fail', 'execute_log': str(err)}
        finally:
            return result

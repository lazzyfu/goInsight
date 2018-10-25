# -*- coding:utf-8 -*-
# edit by fuzongfei
import datetime
import json

import simplejson
from pymysqlreplication import BinLogStreamReader
from pymysqlreplication.event import GtidEvent
from pymysqlreplication.row_event import DeleteRowsEvent, UpdateRowsEvent, WriteRowsEvent

GTID_LOG_EVENT = 0x21


class ReadRemoteBinlog(object):
    """
    binlog_file：读取的binlog文件
    start_pos：开始读取的position
    end_pos：结束读取的position
    """

    def __init__(self, binlog_file=None, start_pos=None, end_pos=None,
                 host=None, port=None, user=None, password=None,
                 affected_rows=None, only_schema=None, only_tables=None):

        self.binlog_file = binlog_file
        self.start_pos = start_pos
        self.end_pos = end_pos
        self.affected_rows = affected_rows

        self.only_schemas = only_schema
        self.only_tables = only_tables

        # 目标数据库配置
        self.mysql_setting = {'host': host,
                              'port': port,
                              'user': user,
                              'passwd': password}

    def _handler_date(self, obj):
        """格式化时间"""
        if type(obj) == datetime.datetime:
            return '{0.year:04d}-{0.month:02d}-{0.day:02d} {0.hour:02d}:{0.minute:02d}:{0.second:02d}'.format(obj)
        if type(obj) == datetime.date:
            return '{0.year:04d}-{0.month:02d}-{0.day:02d} 00:00:00'.format(obj)
        if type(obj) == datetime.timedelta:
            return str(obj)

    def _val_join(self, items):
        """组合column name, column value"""
        k, v = items
        if v is None:
            return f"{k} IS NULL"
        else:
            if isinstance(v, int):
                return f"`{k}`={v}"
            else:
                return f"`{k}`=\"{v}\""

    def _del_join(self, items):
        """
        type == 'DELETE'类型
        对values进行处理
        """
        v = items
        NoneType = type(None)
        if isinstance(v, NoneType):
            return 'NULL'
        elif isinstance(v, int):
            return f'{v}'
        else:
            return f"'{v}'"

    def _upd_join(self, items):
        """
        type == 'UPDATE'类型
        组合column name, column value
        """
        k, v = items
        if v is None:
            return f"{k}=NUll"
        else:
            if isinstance(v, int):
                return f"`{k}`={v}"
            else:
                return f"`{k}`=\"{v}\""

    def _format_binlog(self, rows):
        return simplejson.dumps(rows, default=self._handler_date)

    def generate_rollback_sql(self, rows):
        rollback_statement = []
        for row in rows:
            format_row = json.loads(self._format_binlog(row))
            type = format_row['type']
            database = format_row['database']
            table = format_row['table']
            sql = ''

            if type == 'INSERT':
                where_expression = ' AND '.join(map(self._val_join, row['values'].items()))
                sql = f"DELETE FROM `{database}`.`{table}` WHERE {where_expression} LIMIT 1;"

            elif type == 'DELETE':
                column_name = ', '.join(map(lambda key: f'`{key}`', row['values'].keys()))
                column_value = ', '.join(map(self._del_join, row['values'].values()))
                sql = f"INSERT INTO `{database}`.`{table}`({column_name}) VALUES ({column_value});"

            elif type == 'UPDATE':
                before_values = ', '.join(map(self._upd_join, row['before'].items()))
                after_values = ' AND '.join(map(self._val_join, row['after'].items()))
                sql = f"UPDATE `{database}`.`{table}` SET {before_values} WHERE {after_values} LIMIT 1;"

            rollback_statement.append(sql)
        return rollback_statement

    def run_by_rows(self):
        rows = []
        stream = BinLogStreamReader(connection_settings=self.mysql_setting,
                                    server_id=100,
                                    only_events=[DeleteRowsEvent, WriteRowsEvent, UpdateRowsEvent, GtidEvent],
                                    log_file=self.binlog_file,
                                    log_pos=self.start_pos,
                                    only_schemas=self.only_schemas,
                                    only_tables=self.only_tables,
                                    resume_stream=True,
                                    blocking=True)

        for binlogevent in stream:
            log_pos = binlogevent.packet.log_pos

            if binlogevent.event_type == GTID_LOG_EVENT:
                # 此处获取每个事务的GTID
                # 不做处理
                gtid = binlogevent.gtid
                pass
            else:
                if log_pos > self.end_pos:
                    print('binlog syncer exit...')
                    break
                else:
                    # 获取当前事务的GTID的影响行数
                    affected_rows = len(binlogevent.rows)
                    # 判断当前事务的GTID的影响行数是否等于传入的影响行数
                    # 由于pymysql执行无法获取到当前事务的GTID以及pymysqlreplication无法获取到binlog的thread_id
                    # 所以无法实现精确定位，只能通过该方式实现，可能存在多备份数据的情况
                    if self.affected_rows == affected_rows:
                        for row in binlogevent.rows:
                            binlog = {'database': binlogevent.schema, 'table': binlogevent.table}
                            if isinstance(binlogevent, DeleteRowsEvent):
                                binlog['values'] = row["values"]
                                binlog['type'] = 'DELETE'
                            elif isinstance(binlogevent, UpdateRowsEvent):
                                binlog["before"] = row["before_values"]
                                binlog["after"] = row["after_values"]
                                binlog['type'] = 'UPDATE'
                            elif isinstance(binlogevent, WriteRowsEvent):
                                binlog['values'] = row["values"]
                                binlog['type'] = 'INSERT'
                            rows.append(binlog)
                    return self.generate_rollback_sql(rows)
            stream.close()

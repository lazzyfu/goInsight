# -*- coding:utf-8 -*-
# edit by fuzongfei
import datetime
import json
import logging

import simplejson
from pymysqlreplication import BinLogStreamReader
from pymysqlreplication.event import GtidEvent
from pymysqlreplication.row_event import DeleteRowsEvent, UpdateRowsEvent, WriteRowsEvent

logger = logging.getLogger('django')


class ReadRemoteBinlog(object):
    """
    binlog_file：读取的binlog文件
    start_pos：开始读取的position
    end_pos：结束读取的position
    sql_type：INSERT、UPDATE、DELETE

    返回数据：
    success: {'status': 'success', 'data': [{'gtid': gtid, 'rbsql': rbsql}, {'gtid': gtid, 'rbsql': rbsql} ...]}
    fail: result = {'status': 'fail', 'msg': str(err)}
    """

    def __init__(self, binlog_file=None, start_pos=None, end_pos=None,
                 host=None, port=None, user=None, password=None, sql_type=None,
                 affected_rows=None, only_schema=None, only_tables=None):

        self.binlog_file = binlog_file
        self.start_pos = start_pos
        self.end_pos = end_pos
        self.affected_rows = affected_rows
        self.sql_type = sql_type

        # only_schema和only_table必须为list类型
        self.only_schemas = only_schema
        self.only_tables = only_tables

        # 目标数据库配置
        self.mysql_setting = {'host': host,
                              'port': port,
                              'user': user,
                              'passwd': password,
                              'max_allowed_packet': 32 * 1024 * 1024
                              }

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

    def _generate_rollback_sql(self, rows):
        rollback_statement = []
        for row in rows:
            format_row = json.loads(self._format_binlog(row))
            type = format_row['type']
            database = format_row['database']
            table = format_row['table']
            # 主键可能由一个字段或多个字段组成
            primary_key = ([format_row.get('primary_key')] if isinstance(format_row.get('primary_key'), str) else list(
                format_row.get('primary_key'))) if format_row.get('primary_key') else []
            sql = ''

            if type == 'INSERT':
                if primary_key:
                    where = ' AND '.join(
                        ['='.join((primary, str(row['values'].get(primary)))) for primary in primary_key])
                else:
                    where = ' AND '.join(map(self._val_join, row['values'].items()))
                sql = f"DELETE FROM `{database}`.`{table}` WHERE {where} LIMIT 1;"

            elif type == 'DELETE':
                column_name = ', '.join(map(lambda key: f'`{key}`', row['values'].keys()))
                column_value = ', '.join(map(self._del_join, row['values'].values()))
                sql = f"INSERT INTO `{database}`.`{table}`({column_name}) VALUES ({column_value});"

            elif type == 'UPDATE':
                before_values = ', '.join(map(self._upd_join, row['before'].items()))
                if primary_key:
                    where = ' AND '.join(
                        ['='.join((primary, str(row['after'].get(primary)))) for primary in primary_key])
                else:
                    where = ' AND '.join(map(self._val_join, row['after'].items()))
                sql = f"UPDATE `{database}`.`{table}` SET {before_values} WHERE {where};"

            rollback_statement.append(sql)
        return rollback_statement

    def _filter_gtid(self, data):
        rollbacksql = []
        dlist = []
        for i in data:
            if isinstance(i, str):
                dlist.append(data.index(i))
        dlist.append(len(data))

        result = []
        i = 0
        while i <= len(dlist) - 1:
            try:
                rr = data[dlist[i]:dlist[i + 1]]
                gtid = rr[0]
                gtid_len = len(rr) - 1
                value = rr[1:]
                result.append({'gtid': gtid, 'gtid_len': gtid_len, 'value': value})
                i += 1
            except IndexError as err:
                break
        for j in result:
            if j['gtid_len'] == self.affected_rows:
                rbsql = self._generate_rollback_sql(j['value'])
                rollbacksql.append({'gtid': j['gtid'], 'rbsql': rbsql})
        return rollbacksql

    def run_by_rows(self):
        try:
            stream = BinLogStreamReader(connection_settings=self.mysql_setting,
                                        server_id=101213112,
                                        only_events=[DeleteRowsEvent, WriteRowsEvent, UpdateRowsEvent, GtidEvent],
                                        resume_stream=True,
                                        blocking=False,
                                        log_file=f'{self.binlog_file}',
                                        log_pos=self.start_pos,
                                        only_schemas=f'{self.only_schemas}',
                                        only_tables=f'{self.only_tables}'
                                        )
            rows = []
            for binlogevent in stream:
                log_pos = stream.log_pos
                if log_pos >= self.end_pos:
                    stream.close()
                    break
                else:
                    if isinstance(binlogevent, GtidEvent):
                        # 此处获取每个事务的GTID
                        # 不做处理
                        gtid = binlogevent.gtid
                        rows.append(gtid)
                    else:
                        # 判断当前事务的GTID的影响行数是否等于传入的影响行数
                        # 由于pymysql执行无法获取到当前事务的GTID以及pymysqlreplication无法获取到binlog的thread_id
                        # 所以无法实现精确定位，只能通过该方式实现，可能存在多备份数据的情况
                        for row in binlogevent.rows:
                            binlog = {'database': binlogevent.schema,
                                      'table': binlogevent.table,
                                      'primary_key': binlogevent.primary_key}
                            if self.sql_type == 'DELETE':
                                if isinstance(binlogevent, DeleteRowsEvent):
                                    binlog['values'] = row["values"]
                                    binlog['type'] = 'DELETE'
                                    rows.append(binlog)
                            if self.sql_type == 'UPDATE':
                                if isinstance(binlogevent, UpdateRowsEvent):
                                    binlog["before"] = row["before_values"]
                                    binlog["after"] = row["after_values"]
                                    binlog['type'] = 'UPDATE'
                                    rows.append(binlog)
                            if self.sql_type == 'INSERT':
                                if isinstance(binlogevent, WriteRowsEvent):
                                    binlog['values'] = row["values"]
                                    binlog['type'] = 'INSERT'
                                    rows.append(binlog)
            stream.close()
            result = {'status': 'success', 'data': self._filter_gtid(rows)}
        except Exception as err:
            result = {'status': 'fail', 'msg': str(err)}

        return result

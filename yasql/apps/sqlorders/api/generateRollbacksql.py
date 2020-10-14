# -*- coding:utf-8 -*-
# edit by fuzongfei
import datetime
import json
import logging

import simplejson
from pymysql import escape_string
from pymysql.constants import FIELD_TYPE
from pymysqlreplication import BinLogStreamReader
from pymysqlreplication.event import QueryEvent
from pymysqlreplication.row_event import DeleteRowsEvent, UpdateRowsEvent, WriteRowsEvent

logger = logging.getLogger('django')


class ReadRemoteBinlog(object):
    """
        binlog_file：读取的binlog文件
        start_pos：开始读取的position
        end_pos：结束读取的position
        trx_timestamp: 事务开始的时间
        affected_rows：事务影响的行数

        返回数据：
        success: {'status': 'success', 'data': [rollbacksql]}
        fail: {'status': 'fail', 'msg': str(err)}
        """

    def __init__(self, binlog_file=None, start_pos=None, end_pos=None,
                 host=None, port=None, user=None, password=None, thread_id=None,
                 only_schema=None, only_tables=None):

        self.binlog_file = binlog_file
        self.start_pos = start_pos
        self.end_pos = end_pos
        self.thread_id = thread_id
        # only_schema和only_table必须为list类型
        self.only_schemas = only_schema
        self.only_tables = only_tables

        # 目标数据库配置
        self.mysql_setting = {'host': host,
                              'port': port,
                              'user': user,
                              'passwd': password,
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
        if isinstance(v, type(None)):
            return 'NULL'
        elif isinstance(v, int):
            return f'{v}'
        elif isinstance(v, str):
            # 使用pymysql.escape_string对数据中的引号进行转义
            escape_v = escape_string(v)
            return f"\"{escape_v}\""
        else:
            return f"\"{v}\""

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
            elif isinstance(v, str):
                # 使用pymysql.escape_string对数据中的引号进行转义
                escape_v = escape_string(v)
                return f"`{k}`=\"{escape_v}\""
            else:
                return f"`{k}`=\"{v}\""

    def _format_binlog(self, row):
        return simplejson.dumps(row, default=self._handler_date)

    def _geometry(self, row):
        """解码Geometry类型
        列类型：{'column': 'GEO_LOCATION', 'type': 255}
        # pymysqlreplication返回的原始数据
        > a = b'\x00\x00\x00\x00\x01\x01\x00\x00\x00\xcd#\x7f0\xf0\x19]@\xb0\x1e\xf7\xad\xd6\xf3C@'
        > bytes.hex(a)
        Out[54]: '000000000101000000cd237f30f0195d40b01ef7add6f34340'

        # 在数据库存储的原始数据为：unhex('000000000101000000cd237f30f0195d40b01ef7add6f34340')
        # 需要将回滚语句里面的"unhex('xxx')"改写为unhex('xxx')插入即可
        mysql> select AsText(unhex('000000000101000000cd237f30f0195d40b01ef7add6f34340'));
        +------------------------------------------------------------------------+
        | AsText(unhex('000000000101000000cd237f30f0195d40b01ef7add6f34340')) |
        +------------------------------------------------------------------------+
        | POINT(116.405285 39.904989)                                            |
        +------------------------------------------------------------------------+
        然后在插入
        """
        for col in row['columns']:
            if col['type'] == FIELD_TYPE.GEOMETRY:
                name = col['column']
                if row['type'] in ['INSERT', 'DELETE']:
                    row['values'][name] = f"unhex('{bytes.hex(row['values'][name])}')"
                if row['type'] == 'UPDATE':
                    row['before'][name] = f"unhex('{bytes.hex(row['before'][name])}')"
                    row['after'][name] = f"unhex('{bytes.hex(row['after'][name])}')"
        return row

    def _generate_rollback_sql(self, rows):
        rollback_statement = []
        for row in rows:
            row = self._geometry(row)
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

    def run_by_rows(self):
        try:
            server_id = 6666666 + int(self.thread_id)
            stream = BinLogStreamReader(connection_settings=self.mysql_setting,
                                        server_id=server_id,
                                        only_events=[DeleteRowsEvent, WriteRowsEvent, UpdateRowsEvent, QueryEvent],
                                        resume_stream=True,
                                        blocking=False,
                                        log_file=f'{self.binlog_file}',
                                        log_pos=self.start_pos,
                                        only_schemas=f'{self.only_schemas}',
                                        only_tables=f'{self.only_tables}'
                                        )
            rows = []
            thread_id = query = None
            for binlogevent in stream:
                log_pos = binlogevent.packet.log_pos
                if log_pos >= self.end_pos:
                    # 当当前的binlogevent日志位置大于结束的binlog时，退出
                    stream.close()
                    break
                else:
                    if isinstance(binlogevent, QueryEvent):
                        thread_id = binlogevent.slave_proxy_id
                        query = binlogevent.query

                    if not isinstance(binlogevent, QueryEvent):
                        if self.thread_id == thread_id and query == 'BEGIN':
                            for row in binlogevent.rows:
                                columns = [{'column': x.name, 'type': x.type} for x in binlogevent.columns]
                                binlog = {'database': binlogevent.schema,
                                          'table': binlogevent.table,
                                          'primary_key': binlogevent.primary_key,
                                          'columns': columns
                                          }
                                if isinstance(binlogevent, DeleteRowsEvent):
                                    binlog['values'] = row["values"]
                                    binlog['type'] = 'DELETE'
                                    rows.append(binlog)
                                if isinstance(binlogevent, UpdateRowsEvent):
                                    binlog["before"] = row["before_values"]
                                    binlog["after"] = row["after_values"]
                                    binlog['type'] = 'UPDATE'
                                    rows.append(binlog)
                                if isinstance(binlogevent, WriteRowsEvent):
                                    binlog['values'] = row["values"]
                                    binlog['type'] = 'INSERT'
                                    rows.append(binlog)

            stream.close()
            result = {'status': 'success', 'data': self._generate_rollback_sql(rows)}
        except Exception as err:
            # print("Exception in user code:")
            # print('-' * 60)
            # traceback.print_exc(file=sys.stdout)
            # print('-' * 60)
            print(err)
            result = {'status': 'fail', 'msg': str(err)}
        return result

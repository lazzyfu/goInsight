# -*- coding:utf-8 -*-
# edit by fuzongfei
import pymysql

from sqlorders.models import MysqlSchemas


class GetSchemasGrantApi(object):
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
                    id = '-'.join((self.host, str(self.port), id))
                    result.append({'id': id,
                                   'text': text,
                                   "icon": 'fa fa-table text-blue'
                                   })
        finally:
            self.conn.close()
        return result

    def get_offline_tables(self):
        """
        返回格式：
        [{"id": 'host-port-schema', 'schema', 'children':
        {
        "id": 'test.tbl1',
        "icon": 'fa fa-table text-blue',
        "text": "tbl"
        }, ...}]
        """
        try:
            with self.conn.cursor() as cursor:
                result = []
                if len(self.schema) > 1:
                    query = f"select table_schema, TABLE_NAME from information_schema.COLUMNS" \
                            f" where table_schema in {self.schema}  group by TABLE_SCHEMA, TABLE_NAME"
                else:
                    query = f"select table_schema, TABLE_NAME from information_schema.COLUMNS" \
                            f" where table_schema='{self.schema[0]}'  group by TABLE_SCHEMA, TABLE_NAME"
                cursor.execute(query)

                data = {}
                s_schema = ''
                for schema, table in cursor.fetchall():
                    if schema == s_schema:
                        data[schema].append(table)
                    else:
                        data[schema] = [table]
                        s_schema = schema

                for k, v in data.items():
                    p_id = '-'.join((self.host, self.port, k))
                    p_text = k

                    c_data = []
                    for t in v:
                        c_id = '-'.join((self.host, str(self.port), '.'.join((k, t))))
                        c_text = t
                        c_data.append({'id': c_id,
                                       'text': c_text,
                                       "icon": 'fa fa-table text-blue'
                                       })
                    result.append({
                        'id': p_id,
                        'text': p_text,
                        'children': c_data
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

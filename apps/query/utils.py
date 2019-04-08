# -*- coding:utf-8 -*-
# edit by fuzongfei
import hashlib
import logging

import pymysql
from django.utils.crypto import get_random_string
from pymysql.constants import CLIENT

from opsql.settings import DATABASES
from orders.models import MysqlConfig
from query.models import MysqlUserGroupMap, MysqlPrivBlacklist

logger = logging.getLogger('django')
# 此处密码只针对本机的127.0.0.1授权，没太大意义，用作查询权限验证而已
LOCAL_QUERY_USER_PASSWORD = 'bFhHa_i4NIYFTwTnV0QhXqI5'


class SyncMysqlRemoteMeta(object):
    """
    同步指定的库表元数据信息到本地
    本地库格式：query_{id}_{table_schema}
    """

    def __init__(self, conn_config):
        self.conn_config = conn_config
        self.IGNORE_SCHEMA = ('INFORMATION_SCHEMA', 'information_schema',
                              'PERFORMANCE_SCHEMA', 'performance_schema',
                              'MYSQL', 'SYS', 'mysql', 'tidb_loader')
        self.LOCAL_SCHEMAMETA = []
        self.REMOTE_SCHEMAMETA = []

        self.REMOTE_QUERY = f"select TABLE_SCHEMA, TABLE_NAME, " \
            f"group_concat(concat_ws('&',TABLE_NAME,COLUMN_NAME,ifnull(COLUMN_TYPE, 'N')) SEPARATOR '#') as MD5 " \
            f"from information_schema.COLUMNS where table_schema not in {self.IGNORE_SCHEMA} " \
            f"group by table_schema,table_name"

        self.LOCAL_QUERY = f"select TABLE_SCHEMA, TABLE_NAME, " \
            f"group_concat(concat_ws('&',TABLE_NAME,COLUMN_NAME,ifnull(COLUMN_TYPE, 'N')) SEPARATOR '#') as MD5 " \
            f"from information_schema.COLUMNS where table_schema like 'query_%' " \
            f"group by table_schema,table_name;"

    def _md5sum(self, schema, data):
        # 在不同mysql版本间使用md5(group_concat( ... order by column))也存在md5不一致
        # 将数据转换为列表并按照字母排序
        to_sorted = ''.join(sorted([x.strip() for x in data.split('#')]))
        d = ''.join([schema, to_sorted])
        """校验字符串，生成MD5"""
        hash_md5 = hashlib.md5()
        hash_md5.update(d.encode('utf-8'))
        return hash_md5.hexdigest()

    def _local_cnx(self):
        """连接到本地数据库"""
        user = DATABASES.get('default').get('USER')
        host = DATABASES.get('default').get('HOST')
        password = DATABASES.get('default').get('PASSWORD')
        port = DATABASES.get('default').get('PORT')
        try:
            cnx = pymysql.connect(host=host,
                                  user=user,
                                  password=password,
                                  port=port,
                                  max_allowed_packet=1024 * 1024 * 1024,
                                  charset='utf8',
                                  client_flag=CLIENT.MULTI_STATEMENTS,
                                  cursorclass=pymysql.cursors.DictCursor)
            # 更改group_concat的默认长度(1024)，太短
            with cnx.cursor() as cursor:
                cursor.execute('set session group_concat_max_len=18446744073709551615;')
            return True, cnx
        except Exception as err:
            logger.error(err)
            return False, None

    def _remote_cnx(self, config):
        """连接到目标数据库"""
        try:
            cnx = pymysql.connect(user=config['user'],
                                  password=config['password'],
                                  host=config['host'],
                                  port=config['port'],
                                  max_allowed_packet=1024 * 1024 * 1024,
                                  charset='utf8',
                                  client_flag=CLIENT.MULTI_STATEMENTS,
                                  cursorclass=pymysql.cursors.DictCursor)
            # 更改group_concat的默认长度(1024)，太短
            with cnx.cursor() as cursor:
                cursor.execute('set session group_concat_max_len=18446744073709551615;')
            return True, cnx
        except Exception as err:
            logger.error(err)
            return False, None

    def make_schemameta(self, cnx, id=None, type=None):
        """获取并生成库表元数据信息"""
        try:
            with cnx.cursor() as cursor:
                if type == 'remote':
                    cursor.execute(self.REMOTE_QUERY)
                if type == 'local':
                    cursor.execute(self.LOCAL_QUERY)

                for i in cursor.fetchall():
                    table_schema = i['TABLE_SCHEMA']
                    table_name = i['TABLE_NAME']
                    # 统一格式化为: query_{id}_{table_schema}
                    fmt_table_schema = '_'.join(('query', str(id), table_schema)) if id else table_schema
                    table_md5 = self._md5sum(schema=fmt_table_schema, data=i['MD5'])
                    tbl_stru_query = f"show create table `{table_schema}`.`{table_name}`"
                    cursor.execute(tbl_stru_query)
                    try:
                        table_struc = cursor.fetchone()['Create Table']
                        result = self.REMOTE_SCHEMAMETA if type == 'remote' else self.LOCAL_SCHEMAMETA
                        result.append({
                            'table_schema': fmt_table_schema,
                            'table_name': table_name,
                            'table_md5': table_md5,
                            'columns': [x.split('&')[1] for x in i['MD5'].split('#')],
                            'table_struc': table_struc
                        })
                    except Exception as err:
                        continue
        except Exception as err:
            logger.error(err)
        finally:
            cnx.close()

    def get_remote_schememeta(self, row):
        """获取远程数据库的库表元数据信息"""
        status, cnx = self._remote_cnx(row)
        if status is True:
            self.make_schemameta(cnx=cnx, id=row['id'], type='remote')

    def get_local_schemameta(self):
        """获取本地库表元数据信息"""
        status, cnx = self._local_cnx()
        if status is True:
            self.make_schemameta(cnx=cnx, type='local')

    def drop_lst(self, data):
        """
        删除本地的库表结构
        """
        status, cnx = self._local_cnx()
        if status is True:
            for row in data:
                table_schema = row['table_schema']
                table_name = row['table_name']
                try:
                    with cnx.cursor() as cursor:
                        # 先删除表结构
                        cursor.execute(
                            f"SET SESSION FOREIGN_KEY_CHECKS=OFF;DROP TABLE IF EXISTS `{table_schema}`.`{table_name}`")
                        logger.info(f'检测到本地表结构不一致，删除本地表：DROP TABLE IF EXISTS `{table_schema}`.`{table_name}`')
                except Exception as err:
                    logger.error(err)
                    continue
        cnx.close()

    def cst_local_mysqluser(self, user):
        """在本地数据库中创建该用户"""
        status, cnx = self._local_cnx()
        if status is True:
            try:
                with cnx.cursor() as cursor:
                    cursor.execute(f"create user '{user}'@'%' identified by '{LOCAL_QUERY_USER_PASSWORD}'")
            except Exception as err:
                logger.error(err)

    def make_local_group_user_grant(self, normal_user=None, super_user=None, schema=None, table=None, columns=None):
        """
        在本地数据库中为用户授权
        先revoke在grant
        """
        status, cnx = self._local_cnx()
        revoke = ';'.join(
            [f"revoke select,update,delete,insert on `{schema}`.`{table}` from '{user}'@'%'" for user in
             [normal_user, super_user]])
        # 超级管理员不受黑名单的限制
        super_user_grant = f"grant select,update,delete,insert on `{schema}`.`{table}` to '{super_user}'@'%'"
        normal_user_grant = ''
        # 对普通用户进行权限控制
        # 如果黑名单中查找到记录，如果列定义的是*，则不为用户生成grant语句，如果列定义的是：id,name ，则生成除定义列之外的grant列权限
        try:
            # 列级别授权
            deny_columns = MysqlPrivBlacklist.objects.get(schema='_'.join(schema.split('_')[2:]), table=table).columns
            if deny_columns != '*':
                has_priv_columns = set(columns).difference(deny_columns.strip().split(','))
                fcolumn = ','.join([f"`{c}`" for c in has_priv_columns])
                normal_user_grant = f"grant select({fcolumn}) on table `{schema}`.{table} " \
                    f"to '{normal_user}'@'%'"
        except MysqlPrivBlacklist.DoesNotExist:
            normal_user_grant = f"grant select,update,delete,insert on `{schema}`.`{table}` to '{normal_user}'@'%'"
        finally:
            if status is True:
                try:
                    with cnx.cursor() as cursor:
                        cursor.execute(revoke)
                except Exception as err:
                    pass

                try:
                    with cnx.cursor() as cursor:
                        cursor.execute(';'.join([super_user_grant, normal_user_grant]))
                except Exception as err:
                    logger.error(err)

    def create_mysql_group_user(self, schema, table, columns):
        """
        schema格式: query_1_broker
        创建mysql组用户，每个库对应一个用户
        存储在数据库中的用户为:
        gn_query_1_broker：给普通开发者使用
        gs_query_1_broker：给超级管理者使用
        """
        id = schema.split('_')[1]
        db = '_'.join(schema.split('_')[2:])

        if not MysqlUserGroupMap.objects.filter(comment_id=id, schema=db).exists():
            # 用户名的长度为：15
            normal_groupuser = '_'.join(['n', get_random_string(14)])
            super_groupuser = '_'.join(['s', get_random_string(14)])
            # 将用户存入到表
            MysqlUserGroupMap.objects.create(comment_id=id, schema=db, group=normal_groupuser)
            MysqlUserGroupMap.objects.create(comment_id=id, schema=db, group=super_groupuser)
            # 创建用户
            self.cst_local_mysqluser(user=normal_groupuser)
            self.cst_local_mysqluser(user=super_groupuser)
        else:
            # 获取用户
            normal_groupuser = MysqlUserGroupMap.objects.get(comment_id=id, schema=db, group__icontains='n_').group
            super_groupuser = MysqlUserGroupMap.objects.get(comment_id=id, schema=db, group__icontains='s_').group
        # 为用户分配表或字段的权限
        self.make_local_group_user_grant(normal_user=normal_groupuser,
                                         super_user=super_groupuser,
                                         schema=schema,
                                         table=table,
                                         columns=columns)

    def create_lst(self, data):
        """
        表不存在时，本地数据库新建库表结构
        """
        status, cnx = self._local_cnx()
        if status is True:
            for row in data:
                table_schema = row['table_schema']
                table_name = row['table_name']
                table_struc = row['table_struc'].replace('\n', '')
                columns = row['columns']
                try:
                    with cnx.cursor() as cursor:
                        # 检查库是否存在，不存在创建
                        cursor.execute(f"CREATE DATABASE IF NOT EXISTS `{table_schema}`")
                        # 新建表
                        cursor.execute(f"set session foreign_key_checks=off; USE `{table_schema}`; {table_struc}")
                        logger.info(f'新建本地表：`{table_schema}`.`{table_name}`')

                        # 创建mysql组用户
                        self.create_mysql_group_user(schema=table_schema, table=table_name, columns=columns)
                except Exception as err:
                    logger.error(err)
                    continue
            cnx.close()

    def table_diff_inspect(self):
        rst_md5 = [x['table_md5'] for x in self.REMOTE_SCHEMAMETA]
        lst_md5 = [x['table_md5'] for x in self.LOCAL_SCHEMAMETA]

        # 远程删除库表
        # 本地需要删除这些表
        delete_md5 = list(set(lst_md5).difference(rst_md5))
        if delete_md5:
            delete_data = [x for x in self.LOCAL_SCHEMAMETA if x['table_md5'] in delete_md5]
            self.drop_lst(delete_data)

        # 远程数据库存在，本地数据库不存在的表
        # 本地创建该表
        add_md5 = list(set(rst_md5).difference(lst_md5))
        if add_md5:
            add_data = [x for x in self.REMOTE_SCHEMAMETA if x['table_md5'] in add_md5]
            self.drop_lst(data=add_data)
            self.create_lst(data=add_data)

    def schema_diff_inspect(self):
        """
        如果远程的库结构删除，本地的也删除
        """
        rst_schema = [x['table_schema'] for x in self.REMOTE_SCHEMAMETA]
        lst_schema = [x['table_schema'] for x in self.LOCAL_SCHEMAMETA]
        delete_schema = list(set(lst_schema).difference(set(rst_schema)))
        status, cnx = self._local_cnx()
        if status is True:
            with cnx.cursor() as cursor:
                for i in delete_schema:
                    try:
                        cursor.execute(f"drop database {i}")
                        logger.info(f'删除本地库: {i}')
                    except cnx.InternalError as err:
                        logger.error(f'删除本地库失败: {i}')
                        logger.error(err)
            cnx.close()

    def run(self):
        self.get_local_schemameta()
        for row in self.conn_config:
            self.get_remote_schememeta(row)

        # 对本地的库和远程的所有库的集合进行比较
        self.table_diff_inspect()
        self.schema_diff_inspect()

        # 重置
        self.LOCAL_SCHEMAMETA = []
        self.REMOTE_SCHEMAMETA = []


class GetGrantSchemaMeta(object):
    """获取当前用户授权的表信息"""

    def __init__(self, user=None, id=None, schema=''):
        # MysqlConfig的主键
        self.id = id
        self.local_user = user
        self.schema = schema
        self.local_schema = '_'.join(['query', str(id), self.schema])
        self.local_password = LOCAL_QUERY_USER_PASSWORD
        self.local_host = DATABASES.get('default').get('HOST')
        self.local_port = DATABASES.get('default').get('PORT') if DATABASES.get('default').get('PORT') else 3306

    def _local_conn(self):
        """本地数据库连接"""
        cnx = pymysql.connect(host=self.local_host,
                              user=self.local_user,
                              password=self.local_password,
                              port=self.local_port,
                              charset="utf8")
        # 更改group_concat的默认长度(1024)，太短
        with cnx.cursor() as cursor:
            cursor.execute('set session group_concat_max_len=18446744073709551615;')
        return cnx

    def _remote_conn(self):
        """远程数据库连接"""
        queryset = MysqlConfig.objects.get(id=self.id)
        cnx = pymysql.connect(host=queryset.host,
                              user=queryset.user,
                              password=queryset.password,
                              port=queryset.port,
                              charset="utf8")
        return cnx

    def get_tab_completion(self):
        """
        获取表和对应的列，tab补全使用
        """
        result = {}
        cnx = self._remote_conn()
        try:
            with cnx.cursor() as cursor:
                tables_query = f"select table_name,group_concat(column_name) from information_schema.COLUMNS " \
                    f"where table_schema='{self.schema}' group by table_name"
                cursor.execute(tables_query)
                tables = {}
                for table_name, column_name in cursor.fetchall():
                    tables[table_name] = list(column_name.split(','))

                result['tables'] = tables
        finally:
            cnx.close()
        return result

    def get_table(self):
        """
        返回格式：
        [{
        "id": 'test.tbl1',
        "icon": 'fa fa-table text-blue',
        "text": "tbl"
        }, ...]
        """
        cnx = self._local_conn()
        result = []
        try:
            with cnx.cursor() as cursor:
                query = f"select table_name, group_concat('<b>', column_name, '</b>', ' ', column_type) as column_name " \
                    f"from information_schema.columns where table_schema='{self.local_schema}' " \
                    f"group by table_schema, table_name"
                cursor.execute(query)
                for table, columns in cursor.fetchall():
                    queryset = MysqlConfig.objects.get(id=self.id)
                    id = '___'.join([queryset.host, str(queryset.port), '.'.join([self.schema, table])])
                    columns_children = [{'id': '.'.join([id, c]), 'text': c, 'icon': 'fa fa-columns text-blue'} for c in
                                        columns.split(',<b>')]
                    columns_length = len(columns.split(',<b>'))
                    result.append({'id': id,
                                   'text': table,
                                   'icon': 'fa fa-table text-blue',
                                   'children': [{
                                       'id': '.'.join([id, 'visible']),
                                       'text': f'字段 ({columns_length})',
                                       'icon': 'fa fa-columns text-blue',
                                       'children': columns_children
                                   }]
                                   })
        finally:
            cnx.close()
        return result

    def get_stru(self):
        """
        返回表结构信息
        """

        result = {}
        cnx = self._remote_conn()
        try:
            with cnx.cursor() as cursor:
                stru_query = f"show create table {self.schema}"
                cursor.execute(stru_query)
                result['stru'] = cursor.fetchone()[1]
        finally:
            cnx.close()
        return result

    def get_index(self):
        """
        返回表索引信息
        """

        result = {}
        cnx = self._remote_conn()
        try:
            cnx.cursorclass = pymysql.cursors.DictCursor
            with cnx.cursor() as cursor:
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
            cnx.close()
        return result

    def get_base(self):
        """
        返回表基本信息
        """

        cnx = self._remote_conn()
        base_data = []
        try:
            cnx.cursorclass = pymysql.cursors.DictCursor
            with cnx.cursor() as cursor:
                schema, table = self.schema.split('.')
                base_query = f"select TABLE_NAME as '表名', TABLE_TYPE as '表类型', ENGINE as '引擎', " \
                    f"ROW_FORMAT as '行格式', TABLE_ROWS as '表行数(估算值)', " \
                    f"round(DATA_LENGTH/1024, 2) as '数据大小(KB)', " \
                    f"round(INDEX_LENGTH/1024, 2) as '索引大小(KB)', " \
                    f"TABLE_COLLATION as '字符集校验规则', TABLE_COMMENT as '备注', " \
                    f"CREATE_TIME as '创建时间'  from information_schema.tables where " \
                    f"table_schema='{schema}' and table_name='{table}'"
                # 获取字段
                cursor.execute(base_query)
                for i in cursor.fetchall():
                    for k, v in i.items():
                        format_row = f"<tr><td>{k}</td><td>{v}</td></tr>"
                        base_data.append(format_row)
        finally:
            cnx.close()
        return base_data


class DbDictQueryApi(object):
    def __init__(self, host, port, schema):
        self.host = host
        self.port = port if isinstance(port, int) else int(port)
        self.schema = schema
        obj = MysqlConfig.objects.get(host=self.host, port=self.port)
        self.user = obj.user
        self.password = obj.password
        self.character = obj.character

    def connect(self):
        cnx = pymysql.connect(host=self.host,
                              user=self.user,
                              password=self.password,
                              port=self.port,
                              database='information_schema',
                              max_allowed_packet=1024 * 1024 * 1024,
                              charset=self.character)
        with cnx.cursor() as cursor:
            cursor.execute('set session group_concat_max_len=18446744073709551615;')
        return cnx

    def query(self):
        cnx = self.connect()
        query = f"select t.TABLE_NAME,if(t.TABLE_COMMENT!='',t.TABLE_COMMENT,'None'),t.CREATE_TIME," \
            f"group_concat(distinct concat_ws('<b>', c.COLUMN_NAME,c.COLUMN_TYPE,if(c.IS_NULLABLE='NO','NOT NULL','NULL')," \
            f"ifnull(c.COLUMN_DEFAULT, ''),ifnull(c.CHARACTER_SET_NAME,''), ifnull(c.COLLATION_NAME,'')," \
            f"ifnull(c.COLUMN_COMMENT, '')) separator '<a>') as COLUMNS_INFO," \
            f"group_concat(distinct concat_ws('<b>', s.INDEX_NAME,if(s.NON_UNIQUE=0,'唯一','不唯一'),s.Cardinality," \
            f"s.INDEX_TYPE,s.COLUMN_NAME) separator '<a>') as INDEX_INFO " \
            f"from COLUMNS c join TABLES t on c.TABLE_SCHEMA = t.TABLE_SCHEMA and c.TABLE_NAME = t.TABLE_NAME " \
            f"join STATISTICS s on c.TABLE_SCHEMA = s.TABLE_SCHEMA and c.TABLE_NAME = s.TABLE_NAME " \
            f"where t.TABLE_SCHEMA='{self.schema}' " \
            f"group by t.TABLE_NAME,t.TABLE_COMMENT,t.CREATE_TIME"
        with cnx.cursor() as cursor:
            cursor.execute(query)
            return cursor.fetchall()


class GetTableInfo(object):
    """获取指定主机的所有表"""

    def __init__(self, host, port, schema=None):
        # self.schema可以是单个库也可以是tuple
        self.host = host
        self.port = port
        self.schema = schema
        config = MysqlConfig.objects.get(host=self.host, port=self.port)
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
                    id = '___'.join((self.host, str(self.port), id))
                    result.append({'id': id,
                                   'text': text,
                                   "icon": 'fa fa-table text-blue'
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

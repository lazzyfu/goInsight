# -*- coding:utf-8 -*-
# edit by fuzongfei
import hashlib
import logging

import pymysql
from django.db.models import F
from pymysql.constants import CLIENT

from sqlaudit.settings import DATABASES
from sqlorders.models import MysqlConfig
from sqlquery.models import MysqlRulesChain

logger = logging.getLogger('django')


class MysqlQueryRemoteMetaOp(object):
    """
    获取远程数据库的库表元数据信息，存储到本地
    本地存储的库名：query_id_`schema_name`
    MysqlConfig的id
    """

    def __init__(self, conn_config):
        self.conn_config = conn_config
        self.IGNORE_SCHEMA = ('INFORMATION_SCHEMA', 'PERFORMANCE_SCHEMA', 'MYSQL', 'SYS')
        self.SCHEMAS = {'remote': [], 'local': []}
        self.TABLE_META = {'remote': [], 'local': []}

    def md5_sum(self, data):
        """校验字符串，生成MD5"""
        hash_md5 = hashlib.md5()
        hash_md5.update(data.encode('utf-8'))
        return hash_md5.hexdigest()

    def _local_cnx(self):
        """本地连接"""
        user = DATABASES.get('default').get('USER')
        host = DATABASES.get('default').get('HOST')
        password = DATABASES.get('default').get('PASSWORD')
        port = DATABASES.get('default').get('PORT') if DATABASES.get('default').get('PORT') else 3306
        try:
            cnx = pymysql.connect(host=host,
                                  user=user,
                                  password=password,
                                  port=port,
                                  max_allowed_packet=1024 * 1024 * 1024,
                                  charset='utf8',
                                  client_flag=CLIENT.MULTI_STATEMENTS,
                                  cursorclass=pymysql.cursors.DictCursor)
            return cnx
        except Exception as err:
            logger.error(err.args[1])

    def _remote_cnx(self, config):
        """远程连接"""
        try:
            cnx = pymysql.connect(user=config['user'],
                                  password=config['password'],
                                  host=config['host'],
                                  port=config['port'],
                                  charset='utf8',
                                  cursorclass=pymysql.cursors.DictCursor)
            return cnx
        except Exception as err:
            logger.error(err.args[1])

    def check_schema_exist(self, cursor, schema_name):
        """检查本地库名是否存在"""
        check_schema_exist = f"SELECT SCHEMA_NAME FROM INFORMATION_SCHEMA.SCHEMATA " \
                             f"WHERE SCHEMA_NAME='{schema_name}'"
        cursor.execute(check_schema_exist)
        if not cursor.fetchone():
            cursor.execute(f'CREATE DATABASE `{schema_name}`')

    def create_lst(self, data):
        """创建本地表结构"""
        cnx = self._local_cnx()
        for row in data:
            schema_name = row['schema_name']
            table_stru = row['table_stru'].replace('\n', '')
            try:
                with cnx.cursor() as cursor:
                    self.check_schema_exist(cursor, schema_name)
                    cursor.execute(f"USE `{schema_name}`;set session foreign_key_checks=off;{table_stru}")
                    logger.info(f'新建本地表结构：USE `{schema_name}`; {table_stru}')
            except Exception as err:
                logger.error(err)
        cnx.close()

    def drop_lst(self, data):
        """
        删除本地的表结构
        """
        cnx = self._local_cnx()
        for row in data:
            schema_name = row['schema_name']
            table_name = row['table_name']
            with cnx.cursor() as cursor:
                try:
                    cursor.execute(f"set session foreign_key_checks=off;use `{schema_name}`;drop table {table_name}")
                    logger.info(f'删除本地表：`{schema_name}`.{table_name}')
                except cnx.InternalError as err:
                    logger.error('执行删除表...')
                    logger.error(err)
        cnx.close()

    def update_mysql_rule_chains(self, data, type):
        # 更新规则链表的记录
        for row in data:
            schema_name = '_'.join(row['schema_name'].split('_')[2:])
            table_name = row['table_name']
            cid = row['schema_name'].split('_')[1]
            comment = MysqlConfig.objects.get(pk=cid).comment
            if type == 'add':
                # 为该表增加一条规则
                if not MysqlRulesChain.objects.filter(schema=schema_name, table=table_name).exists():
                    MysqlRulesChain.objects.create(cid_id=cid,
                                                   action='allow',
                                                   schema=schema_name,
                                                   table=table_name,
                                                   comment=f'{comment}-{schema_name}-{table_name}'
                                                   )
                    logger.info(f'新增表规则，表名：{schema_name}.{table_name}')

            if type == 'delete':
                # 删除规则
                if MysqlRulesChain.objects.filter(schema=schema_name, table=table_name).exists():
                    MysqlRulesChain.objects.get(schema=schema_name, table=table_name).delete()
                    logger.info(f'删除表规则，表名：{schema_name}.{table_name}')

    def check_diff_schema(self):
        """比对远端和本地数据库有哪些库不同，远程的库删除，本地的也删除"""
        diff_schemas = list(set(self.SCHEMAS['local']).difference(set(self.SCHEMAS['remote'])))
        cnx = self._local_cnx()
        with cnx.cursor() as cursor:
            for i in diff_schemas:
                try:
                    cursor.execute(f"drop database {i}")
                except cnx.InternalError as err:
                    logger.error('执行删除本地库...')
                    logger.error(err)
        cnx.close()

    def check_diff_table(self):
        rst_md5 = [x['schema_table_md5'] for x in self.TABLE_META['remote']]
        lst_md5 = [x['schema_table_md5'] for x in self.TABLE_META['local']]

        # 远程存在，本地不存在
        # 本地新建
        add_values = list(set(rst_md5).difference(lst_md5))
        if add_values:
            # 本地创建表
            add_table_meta = [x for x in self.TABLE_META['remote'] if x['schema_table_md5'] in add_values]
            self.create_lst(add_table_meta)
            self.update_mysql_rule_chains(add_table_meta, type='add')

        # 远程删除
        # 本地需要删除这些表
        delete_values = list(set(lst_md5).difference(rst_md5))
        if delete_values:
            delete_table_meta = [x for x in self.TABLE_META['local'] if x['schema_table_md5'] in delete_values]
            self.drop_lst(delete_table_meta)
            self.update_mysql_rule_chains(delete_table_meta, type='delete')

        # 远程和本地都存在的表，检测表结构是否相同
        equal_values = list(set(rst_md5).intersection(lst_md5))
        if equal_values:
            rts = [x['table_stru_md5'] for x in self.TABLE_META['remote'] if x['schema_table_md5'] in equal_values]
            lts = [x['table_stru_md5'] for x in self.TABLE_META['local'] if x['schema_table_md5'] in equal_values]
            diff_struc = list(set(lts).difference(rts))
            data = [x['schema_table_md5'] for x in self.TABLE_META['local'] if x['table_stru_md5'] in diff_struc]
            change_table_meta = [x for x in self.TABLE_META['remote'] if x['schema_table_md5'] in data]
            self.drop_lst(change_table_meta)
            self.create_lst(change_table_meta)

    def get_schema_meta(self, cnx, schema_query, id=None, type=None):
        with cnx.cursor() as cursor:
            cursor.execute(schema_query)
            for i in cursor.fetchall():
                # 库名格式：query_id_db，例如：query_11_testdb
                schema_name = '_'.join(('query', str(id), i['schema_name'])) if id else i['schema_name']
                self.SCHEMAS[type].append(schema_name)
                table_query = f"select table_name from information_schema.tables " \
                              f"where table_schema='{i['schema_name']}' and table_type in ('BASE TABLE')"
                cursor.execute(table_query)
                for s in cursor.fetchall():
                    table_name = s['table_name']
                    schema_table_md5 = self.md5_sum('_'.join((schema_name, table_name)))
                    # 获取表结构
                    table_stru_query = f"show create table `{i['schema_name']}`.{table_name}"
                    cursor.execute(table_stru_query)
                    for t in cursor.fetchall():
                        table_stru = t['Create Table']
                        table_stru_md5 = self.md5_sum(table_stru)
                        self.TABLE_META[type].append({
                            'schema_name': schema_name,
                            'table_name': table_name,
                            'table_stru': table_stru,
                            'table_stru_md5': table_stru_md5,
                            'schema_table_md5': schema_table_md5
                        })

    def get_remote(self, row):
        """获取远程库表元数据"""
        cnx = self._remote_cnx(row)
        try:
            id = row['id']
            schema_query = f"select schema_name from information_schema.schemata " \
                           f"where schema_name not in {self.IGNORE_SCHEMA}"
            self.get_schema_meta(cnx, schema_query, id=id, type='remote')
        except Exception as err:
            msg = f"任务：sync_remote_tablemeta 连接到数据库({row.get('host')})失败, {err.args[1]}"
            logger.error(msg)
        finally:
            cnx.close()

    def get_local(self, row):
        """获取本地库表元数据"""
        cnx = self._local_cnx()
        try:
            id = row['id']
            schema_query = f"select schema_name from information_schema.schemata " \
                           f"where schema_name like 'query_{id}\_%'"
            self.get_schema_meta(cnx, schema_query, type='local')
        except Exception as err:
            msg = f"任务：sync_remote_tablemeta 连接到数据库({row.get('host')})失败, {err.args[1]}"
            logger.error(msg)
        finally:
            cnx.close()

    def remove_not_exist_host_db(self):
        # 移除后台不存在的主机配置的本地Schema
        cnx = self._local_cnx()
        ids = list(MysqlConfig.objects.filter(is_type__in=(0, 2)).values_list('id', flat=True))
        schema_query = f"select schema_name from information_schema.schemata where schema_name regexp '^query_[0-9]'"
        with cnx.cursor() as cursor:
            cursor.execute(schema_query)
            for row in cursor.fetchall():
                schema = row['schema_name']
                try:
                    if int(schema.split('_')[1]) not in ids:
                        cursor.execute(f"drop database {schema}")
                        logger.info(f'删除本地不存在数据库配置的库{schema}')
                except Exception as err:
                    logger.error(err.args[1])
                    continue
        cnx.close()

    def run(self):
        for row in self.conn_config:
            self.get_remote(row)
            self.get_local(row)
            self.check_diff_table()
            self.check_diff_schema()
            self.SCHEMAS = {'remote': [], 'local': []}
            self.TABLE_META = {'remote': [], 'local': []}
        self.remove_not_exist_host_db()


class CreateLocalMysqlUser(object):
    def __init__(self, users):
        # 传入的普通用户列表
        self.users = users
        self.password = 'LNjLJ6MeMJiZznL6'

        # settings.py中配置的超级用户
        # 该用户应该拥有all privileges
        self.super_user = DATABASES.get('default').get('USER')
        self.super_password = DATABASES.get('default').get('PASSWORD')
        self.host = DATABASES.get('default').get('HOST')
        self.port = DATABASES.get('default').get('PORT') if DATABASES.get('default').get('PORT') else 3306

        self.DIFF_USERS = []

    def _super_conn(self):
        """
        以超级用户连接本地数据库
        本地数据库就是django应用库表所在的数据库
        """
        cnx = pymysql.connect(user=self.super_user,
                              password=self.super_password,
                              host=self.host,
                              port=self.port,
                              max_allowed_packet=1024 * 1024 * 1024,
                              charset='utf8',
                              client_flag=CLIENT.MULTI_STATEMENTS,
                              cursorclass=pymysql.cursors.DictCursor)
        return cnx

    def create_local_user(self):
        cnx = self._super_conn()
        with cnx.cursor() as cursor:
            for user in self.DIFF_USERS:
                query_create_user = f"create user '{user}'@'%' identified by '{self.password}'"
                cursor.execute(query_create_user)

    def check_user_exist(self):
        cnx = self._super_conn()
        query_user = "select user from mysql.user"
        local_users = []
        with cnx.cursor() as cursor:
            cursor.execute(query_user)
            for row in cursor.fetchall():
                local_users.append(row['user'])
        self.DIFF_USERS = list(set(self.users).difference(local_users))
        if self.DIFF_USERS:
            self.create_local_user()

    def run(self):
        self.check_user_exist()


class GetGrantSchemaMeta(object):
    """获取当前用户授权的表信息"""

    def __init__(self, user=None, id=None, schema=''):
        # MysqlConfig的主键
        self.id = id
        # user = request.user.username
        self.local_user = user
        self.schema = schema
        self.local_schema = '_'.join(['query', str(id), self.schema])
        self.local_password = 'LNjLJ6MeMJiZznL6'
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


class MySQLQueryRulesOperate(object):
    def __init__(self, before_rule_id, after_rule_id, before_users, after_users):
        self.before_rule_id = before_rule_id
        self.after_rule_id = after_rule_id
        self.before_users = before_users
        self.after_users = after_users

    def _local_cnx(self):
        """本地连接"""
        user = DATABASES.get('default').get('USER')
        host = DATABASES.get('default').get('HOST')
        password = DATABASES.get('default').get('PASSWORD')
        port = DATABASES.get('default').get('PORT') if DATABASES.get('default').get('PORT') else 3306
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
        return cnx

    def execute_statements(self, data):
        cnx = self._local_cnx()
        with cnx.cursor() as cursor:
            for s in data:
                try:
                    cursor.execute(s)
                    logger.info(f'执行授权: {s}')
                except Exception as err:
                    logger.warning(err)
                    continue

    def generate_table_statements(self, users=None, type=None, schema=None, table=None):
        # 表级授权，字段级授权移除了，主要是MySQL执行字段级授权太慢，不好控制
        statements = []
        st = '.'.join([f"`{schema}`", table])

        # 授权操作
        # 当type = add时，allow --> grant
        if type == 'add':
            statements.append(
                ';'.join(
                    [f"grant select,update,delete,insert on table {st} to '{user}'@'%'" for user in users]))
        else:
            statements.append(
                ';'.join(
                    [f"grant select,update,delete,insert on {st} to '{user}'@'%'" for user in users]))
        # 移除操作
        # 当type = remove时，allow --> revoke
        if type == 'remove':
            statements.append(
                ';'.join([f"revoke select,update,delete,insert on table {st} from '{user}'@'%'" for user in users]))

        statements.append('flush privileges;')
        self.execute_statements(statements)

    def analyze_priv(self, users=None, type=None, privs_id=None):
        queryset = MysqlRulesChain.objects.filter(id__in=privs_id).annotate(
            c_id=F('cid__id')
        )
        for row in queryset:
            schema = '_'.join(['query', str(row.c_id), row.schema])
            table = row.table

            self.generate_table_statements(users, type, schema, table)

    def run(self):
        # 处理新增加的用户
        # 给新增的用户增加after_rule_id的权限
        add_users = list(set(self.after_users).difference(set(self.before_users)))
        if add_users:
            self.analyze_priv(users=add_users, type='add', privs_id=self.after_rule_id)

        # 处理移除的用户
        # 给移除的用户移除before_rule_id的权限
        remove_users = list(set(self.before_users).difference(set(self.after_users)))
        if remove_users:
            self.analyze_priv(users=remove_users, type='remove', privs_id=self.after_rule_id)

        # 检查是否增加权限
        # 给before_users用户增加新增的权限
        add_privs_id = list(set(self.after_rule_id).difference(set(self.before_rule_id)))
        if add_privs_id:
            self.analyze_priv(users=self.before_users, type='add', privs_id=add_privs_id)

        # 检查是否移除权限
        # 给before_users用户删除移除的权限
        remove_privs_id = list(set(self.before_rule_id).difference(set(self.after_rule_id)))
        if remove_privs_id:
            self.analyze_priv(users=self.before_users, type='remove', privs_id=remove_privs_id)

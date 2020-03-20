# -*- coding:utf-8 -*-
# edit by fuzongfei
import datetime

import pymysql
from django import forms
from pymysql.constants import CLIENT

from opsql.settings import DATABASES
from orders.models import MysqlConfig
from orders.utils.tools import checkdbstatus
from query.models import MySQLQueryLog, MysqlUserGroupMap
from query.sqlQueryApi import MySQLQueryApi
from query.utils import GetGrantSchemaMeta, DbDictQueryApi, GetTableInfo


class GetGrantSchemaForm(forms.Form):
    id = forms.CharField(max_length=2048, required=False, label=u'jstree传入的node.id')
    text = forms.CharField(max_length=2048, required=False, label=u'jstree传入的node.text')

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

    def query(self, request):
        cdata = self.cleaned_data
        id = cdata.get('id')
        map_mysql_usergroup = list(
            MysqlUserGroupMap.objects.filter(user__username=request.user.username).values_list('group', flat=True))
        map_mysql_usergroup.extend(['None', 'None'])  # 防止列表为空转换为元组有问题
        context = []
        if id == '#':
            # 返回当前用户所有授权的库
            query = f"select distinct `db` from mysql.columns_priv where `user` in {tuple(map_mysql_usergroup)} " \
                f"and `db` regexp '^query_' union select distinct `db` from mysql.tables_priv " \
                f"where `user` in {tuple(map_mysql_usergroup)} and `db` regexp '^query_'"
            cnx = self._local_cnx()
            with cnx.cursor() as cursor:
                cursor.execute(query)
                for row in cursor.fetchall():
                    id = row['db'].split('_')[1]
                    schema = '_'.join(row['db'].split('_')[2:])
                    if MysqlConfig.objects.filter(pk=id).exists():
                        obj = MysqlConfig.objects.get(pk=id)
                        show_schema = '_'.join((obj.comment, schema))
                        context.append({
                            'id': '___'.join((obj.host, str(obj.port), schema)),
                            'text': show_schema,
                            'children': True,
                        })

        if len(id.split('___')) == 3:
            # 获取当前用户授权库的表
            host = id.split('___')[0]
            port = id.split('___')[1]
            queryset = MysqlConfig.objects.get(host=host, port=port)
            schema = id.split('___')[2]
            query_user = MysqlUserGroupMap.objects.filter(comment__id=queryset.pk, user__username=request.user.username,
                                                          schema=schema).values_list('group', flat=True)
            user = query_user[0]
            for i in query_user:
                if i.startswith('s_'):
                    user = i
            data = GetGrantSchemaMeta(user, queryset.id, schema).get_table()
            context = data

        return context


class GetTableStrucForm(forms.Form):
    schema = forms.CharField(max_length=1024, required=True)

    def query(self):
        cdata = self.cleaned_data
        host, port, schema = cdata.get('schema').split('___')
        id = MysqlConfig.objects.get(host=host, port=port).id
        if len(schema.split('.')) == 2:
            schema = schema
        if len(schema.split('.')) == 3:
            schema = '.'.join(schema.split('.')[:2])
        data = GetGrantSchemaMeta(id=id, schema=schema).get_stru()
        context = {'status': 0, 'data': data}
        return context


class GetTableIndexForm(forms.Form):
    schema = forms.CharField(max_length=1024, required=True)

    def query(self):
        cdata = self.cleaned_data
        host, port, schema = cdata.get('schema').split('___')
        id = MysqlConfig.objects.get(host=host, port=port).id
        if len(schema.split('.')) == 2:
            schema = schema
        if len(schema.split('.')) == 3:
            schema = '.'.join(schema.split('.')[:2])
        data = GetGrantSchemaMeta(id=id, schema=schema).get_index()
        context = {'status': 0, 'data': data}
        return context


class GetTableBaseForm(forms.Form):
    schema = forms.CharField(max_length=1024, required=True)

    def query(self):
        cdata = self.cleaned_data
        host, port, schema = cdata.get('schema').split('___')
        id = MysqlConfig.objects.get(host=host, port=port).id
        if len(schema.split('.')) == 2:
            schema = schema
        if len(schema.split('.')) == 3:
            schema = '.'.join(schema.split('.')[:2])
        data = GetGrantSchemaMeta(id=id, schema=schema).get_base()
        context = {'status': 0, 'data': data}
        return context


class ExecSqlQueryForm(forms.Form):
    schema = forms.CharField(max_length=1024)
    page_hash = forms.CharField(max_length=64, min_length=32, label=u'页面hash')
    contents = forms.CharField(widget=forms.Textarea, label=u'sql')

    def execute(self, request):
        cdata = self.cleaned_data
        schema = cdata.get('schema')
        page_hash = cdata.get('page_hash')
        contents = cdata.get('contents')

        host, port, schema = schema.split('___')
        if len(schema.split('.')) in [2, 3]:
            schema = schema.split('.')[0]

        mysql_query = MySQLQueryApi(user=request.user.username, sqls=contents, host=host, port=port, schema=schema)
        result = mysql_query.query(page_hash)
        return result


class GetTablesForm(forms.Form):
    schema = forms.CharField()

    def query(self):
        cdata = self.cleaned_data
        schema = cdata['schema']
        host, port, schema = schema.split(',')

        status, msg = checkdbstatus(host, port)
        if status:
            table_list = GetTableInfo(host, port, schema).get_column_info()
            context = {'status': 0, 'msg': '', 'data': table_list}
        else:
            context = {'status': 2, 'msg': f'无法连接到数据库，请联系管理员，\n主机: {host}\n端口: {port}'}
        return context


class GetHistorySqlForm(forms.Form):
    def query(self, request):
        queryset = MySQLQueryLog.objects.filter(user=request.user.username, query_status=u'成功').order_by(
            '-created_at').values('created_at', 'query_sql')[:1000]
        data = []
        for row in queryset:
            created_at = '时间：' + (row['created_at'] + datetime.timedelta(hours=8)).strftime("%Y-%m-%d %H:%M:%S")
            query_sql = 'SQL语句：' + row['query_sql']
            data.append('\n'.join((created_at, query_sql)))
        context = {'status': 0, 'data': data}
        return context


class GetFilterHistorySqlForm(forms.Form):
    contents = forms.CharField(required=False, max_length=128, label=u'搜索的内容')

    def query(self, request):
        cdata = self.cleaned_data
        contents = cdata.get('contents')

        if contents:
            queryset = MySQLQueryLog.objects.filter(user=request.user.username, query_status=u'成功',
                                                    query_sql__icontains=contents).order_by(
                '-created_at').values('created_at', 'query_sql')
        else:
            queryset = MySQLQueryLog.objects.filter(user=request.user.username,
                                                    query_status=u'成功').order_by(
                '-created_at').values('created_at', 'query_sql')[:1000]
        data = []
        for row in queryset:
            created_at = '时间：' + (row['created_at'] + datetime.timedelta(hours=8)).strftime("%Y-%m-%d %H:%M:%S")
            query_sql = 'SQL语句：' + row['query_sql']
            data.append('\n'.join((created_at, query_sql)))

        if len(data) == 0:
            data.append('未找到SQL记录')
        context = {'status': 0, 'data': data}
        return context


class DbDictForm(forms.Form):
    envi_id = forms.CharField(required=True, label=u'环境')
    database = forms.CharField(required=True, max_length=256, label=u'库名')

    def query(self):
        cdata = self.cleaned_data
        host, port, schema = cdata.get('database').split(',')
        dbquery = DbDictQueryApi(host, port, schema)
        return dbquery.query()

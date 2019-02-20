# -*- coding:utf-8 -*-
# edit by fuzongfei
import datetime

import pymysql
from django import forms
from django.core.exceptions import ValidationError
from pymysql.constants import CLIENT

from sqlorders.models import MysqlSchemas, MysqlConfig, SqlOrdersEnvironment
from sqlquery.models import MySQLQueryLog
from sqlquery.sqlQueryApi import MySQLQuery
from sqlquery.utils import GetGrantSchemaMeta, DbDictQueryApi
from sqlaudit.settings import DATABASES


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
        text = cdata.get('text')
        context = []
        if id == '#':
            # 返回当前用户所有授权的库
            query = f"select distinct `db` from mysql.tables_priv where `user`='{request.user.username}' and `db` like 'query_%'"
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
            data = GetGrantSchemaMeta(request.user.username, queryset.id, schema).get_table()
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
    contents = forms.CharField(widget=forms.Textarea, label=u'sql')
    
    def execute(self, request):
        cdata = self.cleaned_data
        schema = cdata.get('schema')
        contents = cdata.get('contents')
        
        host, port, schema = schema.split('___')
        if len(schema.split('.')) in [2, 3]:
            schema = schema.split('.')[0]
        
        # 判断是否是只读
        is_type = MysqlConfig.objects.get(host=host, port=port).is_type
        is_rw = None
        if is_type == 0:
            is_rw = 'r'
        if is_type == 2:
            is_rw = 'rw'
        mysql_query = MySQLQuery(user=request.user.username, querys=contents, host=host, port=port,
                                 schema=schema, rw=is_rw)
        result = mysql_query.query()
        return result


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


# class MysqlRulesChainForm(forms.ModelForm):
#     schema = forms.ChoiceField(choices=[(x[1], '_'.join(x))
#                                         for x in
#                                         MysqlSchemas.objects.filter(is_type__in=(0, 2)).exclude(
#                                             host=DATABASES.get('default')).values_list('comment',
#                                                                                        'schema')],
#                                label=u'库名')

def envi_validator(value):
    value = value if isinstance(value, int) else int(value)
    envi = [x for x in list(SqlOrdersEnvironment.objects.all().values_list('envi_id', flat=True))]
    if value not in envi:
        raise ValidationError('请选择正确的工单环境')


class DbDictForm(forms.Form):
    envi_id = forms.CharField(required=True, validators=[envi_validator], label=u'环境')
    database = forms.CharField(required=True, max_length=256, label=u'库名')
    
    def query(self):
        cdata = self.cleaned_data
        host, port, schema = cdata.get('database').split(',')
        dbquery = DbDictQueryApi(host, port, schema)
        return dbquery.query()

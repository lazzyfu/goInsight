import time

from django.core.exceptions import PermissionDenied
from django.db.models import F
from django.http import JsonResponse
from django.shortcuts import render
# Create your views here.
from django.views import View

from mstats.models import MysqlSchemaInfo, MysqlSchemaGrant, WebShellInfo
from mstats.utils import MySQLQuery
from project_manager.inception.inception_api import GetTableMetaInfo
from project_manager.utils import check_db_conn_status
from user_manager.permissions import permission_required
from utils.tools import format_request


class GetSchemaInfo(View):
    """获取schema列表"""

    def get(self, request):
        result = format_request(request)
        envi = result.get('envi').split(',')
        is_master = result.get('is_master')

        # host存在，非生产环境
        if result.get('host'):
            host, port = result.get('host').split(',')
            data = MysqlSchemaInfo.objects.filter(envi__in=envi, host=host, port=port).values('id', 'host', 'port',
                                                                                              'schema', 'comment')
        else:
            # 否则生产环境
            data = MysqlSchemaInfo.objects.filter(envi__in=envi, is_master=is_master).values('id', 'host', 'port',
                                                                                             'schema', 'comment')
        return JsonResponse(list(data), safe=False)


class GetStruInfoView(View):
    """返回表结构和索引等信息"""

    def get(self, request):
        result = format_request(request)
        host, port, schema = result.get('schema').split('-')
        if len(schema.split('.')) == 2:
            data = GetTableMetaInfo(host, port, schema).get_stru_info()
            context = {'status': 0, 'data': data}
        else:
            context = {'status': 2, 'msg': ''}
        return JsonResponse(context, safe=False)


class GetOfflineQuerySchemaTreeView(View):
    """
    获取非生产环境查询接口的schema和table信息，渲染成tree结构
    数据格式：192.168.203.16,3306
    """

    def get(self, request):
        result = format_request(request)
        envi = result.get('envi').split(',')
        host, port = result.get('hostname').split(',')
        obj = MysqlSchemaInfo.objects.filter(envi__in=envi, host=host, port=port).values_list('schema', flat=True)
        schema = tuple(obj)
        context = GetTableMetaInfo(host, port, schema).get_offline_tables()
        return JsonResponse(context, safe=False)


class GetOnlineQuerySchemaTreeView(View):
    """获取生产环境查询接口schema和表信息，提供给前端js，渲染成tree结构"""

    def get(self, request):
        # 线上查询，授权过滤
        query = f"select b.id, b.host, b.port, b.schema from auditsql_mysql_schema_grant a " \
                f"join auditsql_mysql_schema b on a.schema_id = b.schema_join join auditsql_useraccount c  " \
                f"on c.uid = a.user_id where c.uid={request.user.uid} and b.envi=3 and b.is_master=0"

        context = []
        for row in MysqlSchemaInfo.objects.raw(query):
            data = GetTableMetaInfo(row.host, row.port, row.schema).get_online_tables()
            context.append({
                'id': '-'.join((row.host, str(row.port), row.schema)),
                'text': row.schema,
                'children': data
            })
        return JsonResponse(context, safe=False)


# 线上环境mysql查询，只读
class ROnlineMySQLQueryView(View):
    @permission_required('can_online_mysql_query')
    def get(self, request):
        return render(request, 'ol_sql_query.html')


class OnlineMySQLQueryView(View):
    @permission_required('can_online_mysql_query')
    def post(self, request):
        data = format_request(request)
        host, port, schema = data.get('schema').split('-')
        if len(schema.split('.')) == 2:
            schema = schema.split('.')[0]
        querys = data.get('contents')

        # 验证传入的host是否合法
        schemas = MysqlSchemaGrant.objects.filter(user__uid=request.user.uid).annotate(
            schemas=F('schema__schema')).values_list('schemas', flat=True)
        if schema in schemas:
            mysql_query = MySQLQuery(querys=querys, host=host, port=port, schema=schema, rw='r', envi=0)
            result = mysql_query.query(request)
        else:
            raise PermissionDenied
        return JsonResponse(result, safe=False)


# 线下环境mysql查询
class ROfflineMySQLQueryView(View):
    @permission_required('can_offline_mysql_query')
    def get(self, request):
        return render(request, 'of_sql_query.html')


class OfflineMySQLQueryView(View):
    """
    接收的数据格式：
    172.17.101.40-3306-test.tbl 或 172.17.101.40-3306-test
    """

    @permission_required('can_offline_mysql_query')
    def post(self, request):
        data = format_request(request)
        host, port, schema = data.get('schema').split('-')
        if len(schema.split('.')) == 2:
            schema = schema.split('.')[0]
        querys = data.get('contents')

        mysql_query = MySQLQuery(querys=querys, host=host, port=port, schema=schema, rw='rw', envi=1)
        result = mysql_query.query(request)
        return JsonResponse(result, safe=False)


class WebSSHView(View):
    def get(self, request):
        return render(request, 'webssh.html')


class GetWebSSHCmdView(View):
    def get(self, request):
        query = f"select a.id, a.command, a.comment from auditsql_web_shell a " \
                f"join auditsql_web_shell_grant b on a.id = b.shell_id " \
                f"join auditsql_useraccount c on b.user_id = c.uid where c.uid = {request.user.uid}"
        data = []
        for row in MysqlSchemaInfo.objects.raw(query):
            data.append({
                'command': row.command,
                'comment': row.comment
            })
        return JsonResponse(list(data), safe=False)

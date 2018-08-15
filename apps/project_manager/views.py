import datetime
import json

import sqlparse
from channels.layers import get_channel_layer
from django.db.models import F
from django.http import JsonResponse, HttpResponse
from django.views import View

from apps.project_manager.inception.inception_api import GetTableMetaInfo, sql_filter, IncepSqlCheck
from mstats.models import MysqlSchemaInfo
from project_manager.forms import SyntaxCheckForm
from project_manager.utils import check_db_conn_status
from user_manager.models import PermissionDetail, RolesDetail
from utils.tools import format_request
from .models import AuditTasks

channel_layer = get_channel_layer()


class BeautifySQLView(View):
    """
    美化SQL
    判断SQL类型（DML还是DDL），并分别进行美化
    最后合并返回
    """

    def post(self, request):
        data = format_request(request)
        sql_content = data.get('sql_content').strip()

        sql_split = []
        for stmt in sqlparse.split(sql_content):
            sql = sqlparse.parse(stmt)[0]
            sql_comment = sql.token_first()
            if isinstance(sql_comment, sqlparse.sql.Comment):
                sql_split.append({'comment': sql_comment.value, 'sql': sql.value.replace(sql_comment.value, '')})
            else:
                sql_split.append({'comment': '', 'sql': sql.value})

        beautify_sql_list = []
        for row in sql_split:
            comment = row['comment']
            sql = row['sql']
            res = sqlparse.parse(sql)
            syntax_type = res[0].token_first().ttype.__str__()
            if syntax_type == 'Token.Keyword.DDL':
                sql_format = sqlparse.format(sql)
                beautify_sql_list.append(comment + sql_format)
            elif syntax_type == 'Token.Keyword.DML':
                sql_format = sqlparse.format(sql, reindent=True)
                beautify_sql_list.append(comment + sql_format)
            elif syntax_type == 'Token.Keyword':
                beautify_sql_list.append(comment + sql)
        context = {'data': '\n\n'.join(beautify_sql_list)}

        return HttpResponse(json.dumps(context))


class GetHostInfoView(View):
    """获取指定环境主机信息"""

    def get(self, request):
        result = format_request(request)
        envi = result.get('envi').split(',')
        is_master = result.get('is_master')

        data = MysqlSchemaInfo.objects.filter(envi__in=envi, is_master=is_master).filter(comment__isnull=False).values(
            'host', 'port', 'comment').distinct()

        return JsonResponse(list(data), safe=False)


class GetTableMetaInfoView(View):
    """获取mysql表的信息，提供给tab补全使用"""

    def post(self, request):
        data = format_request(request)
        print()
        host, port = data.get('schema').split(',')
        status, msg = check_db_conn_status(host, port)
        if status:
            meta_data = GetTableMetaInfo(host, port).get_column_info()
            context = {'status': 0, 'msg': '', 'data': meta_data}
        else:
            context = {'status': 2, 'msg': f'无法连接到数据库，请联系DBA\n主机: {host}\n端口: {port}'}
        return HttpResponse(json.dumps(context))


class GetTableInfo(View):
    """获取指定主机的所有表"""

    def post(self, request):
        data = format_request(request).get('schema')
        host, port, schema = data.split(',')

        status, msg = check_db_conn_status(host, port)
        if status:
            table_list = GetTableMetaInfo(host, port, schema).get_column_info()
            context = {'status': 0, 'msg': '', 'data': table_list}
        else:
            context = {'status': 2, 'msg': f'无法连接到数据库，请联系DBA\n主机: {host}\n端口: {port}'}
        return HttpResponse(json.dumps(context))


class GetAuditUserView(View):
    """获取DBA信息"""

    def get(self, request):
        result = []
        role_list = PermissionDetail.objects.annotate(
            role_name=F('role__role_name'),
            permission_name=F('permission__permission_name')).filter(
            permission__permission_name='can_approve'
        ).values_list('role_name')

        for i in role_list:
            role_name = i[0]
            data = RolesDetail.objects.annotate(
                username=F('user__username'),
                displayname=F('user__displayname'),
            ).filter(role__role_name=role_name).values('username', 'displayname').order_by('username')

            result.append({'user': list(data)})
        return JsonResponse(result, safe=False)


class SyntaxCheckView(View):
    """语法检查"""

    def post(self, request):
        form = SyntaxCheckForm(request.POST)
        if form.is_valid():
            cleaned_data = form.cleaned_data
            host, port, database = cleaned_data['database'].split(',')
            operate_type = cleaned_data.get('operate_type')
            contents = cleaned_data['contents']

            # 对检测的SQL类型进行区分
            filter_result = sql_filter(contents, operate_type)

            # 实例化
            of_audit = IncepSqlCheck(contents, host, port, database, request.user.username)

            if filter_result['status'] == 2:
                context = filter_result
            else:
                # SQL语法检查
                context = of_audit.run_check()

            return HttpResponse(json.dumps(context))
        else:
            error = "请选择主机、库名或项目组"
            context = {'status': 2, 'msg': error}
            return HttpResponse(json.dumps(context))


class GetOnlineAuditTasksList(View):
    def get(self, request):
        """
        如果当前任务的提交时间大于任务设置的过期时间，不允许选择该任务
        is_disable：是否禁用，0：否，1：是
        """
        before_14_days = (datetime.datetime.now() - datetime.timedelta(days=14)).strftime("%Y-%m-%d %H:%M:%S")
        query = f"select id,tasks,if(now()>expire_time,1,0) as is_disable from auditsql_audit_tasks " \
                f"where created_at >= '{before_14_days}' order by created_at desc"
        data = []
        for row in AuditTasks.objects.raw(query):
            data.append({'tasks': row.tasks, 'is_disable': row.is_disable})

        return JsonResponse(data, safe=False)

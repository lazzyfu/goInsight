import json

import sqlparse
from channels.layers import get_channel_layer
from django.db.models import F
from django.http import JsonResponse, HttpResponse
from django.views import View

from apps.project_manager.inception.inception_api import GetSchemaInfo, sql_filter, IncepSqlCheck
from project_manager.forms import SyntaxCheckForm
from project_manager.utils import check_mysql_conn
from user_manager.models import GroupsDetail, Contacts, PermissionDetail, RolesDetail
from utils.tools import format_request
from .models import InceptionHostConfigDetail, InceptionHostConfig

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
        try:
            for row in sql_split:
                comment = row['comment']
                sql = row['sql']
                res = sqlparse.parse(sql)
                if res[0].tokens[0].ttype[1] == 'DDL':
                    sql_format = sqlparse.format(sql)
                    beautify_sql_list.append(comment + sql_format)
                else:
                    sql_format = sqlparse.format(sql, reindent=True)
                    beautify_sql_list.append(comment + sql_format)
            context = {'data': '\n\n'.join(beautify_sql_list)}
        except Exception as err:
            context = {'status': 2, 'msg': '注释或语法错误'}

        return HttpResponse(json.dumps(context))


class IncepHostConfigView(View):
    """获取指定的数据库配置"""

    def get(self, request):
        data = format_request(request)
        print(data)
        config_type = 0 if data.get('type') is None else data.get('type')
        purpose = 0 if data.get('purpose') is None else data.get('purpose')
        user_in_group = request.session.get('groups')
        result = InceptionHostConfigDetail.objects.annotate(host=F('config__host'),
                                                            comment=F('config__comment')
                                                            ).filter(config__type=config_type). \
            filter(config__purpose=purpose). \
            filter(config__is_enable=0). \
            filter(group__group_id__in=user_in_group). \
            values('host', 'comment')
        return JsonResponse(list(result), safe=False)


class GetSchemaView(View):
    """获取选中环境的数据库库名"""

    def post(self, request):
        data = format_request(request)
        host = data['host']
        obj = InceptionHostConfig.objects.get(host=host)
        result = check_mysql_conn(obj.user, host, obj.password, obj.port)
        if result['status'] == 'INFO':
            db_list = GetSchemaInfo(host).get_values()
            context = {'status': 0, 'msg': '', 'data': db_list}
        else:
            context = {'status': 2, 'msg': f'获取列表失败，不能连接到mysql服务器：{host}'}
        return HttpResponse(json.dumps(context))


class GroupInfoView(View):
    def get(self, request):
        groups = GroupsDetail.objects.filter(
            user__uid=request.user.uid).annotate(
            group_id=F('group__group_id'), group_name=F('group__group_name')) \
            .values('group_id', 'group_name')

        return JsonResponse(list(groups), safe=False)


class GetAuditUserView(View):
    """获取指定项目组的批准人和执行人信息"""

    def post(self, request):
        data = format_request(request)
        group_id = data.get('group_id')
        result = []
        if group_id:
            role_list = PermissionDetail.objects.annotate(
                role_name=F('role__role_name'),
                permission_name=F('permission__permission_name')).filter(
                permission__permission_name__in=('can_approve', 'can_execute')
            ).values_list('role_name', 'permission_name')

            for i in role_list:
                role_name = i[0]
                can_priv = i[1]

                uid = RolesDetail.objects.annotate(uid=F('user__uid')).filter(
                    role__role_name=role_name).values_list(
                    'uid', flat=True)

                data = GroupsDetail.objects.annotate(
                    uid=F('user__uid'),
                    username=F('user__username'),
                    email=F('user__email'),
                ).filter(group__group_id=group_id).filter(user__uid__in=uid).values('uid', 'username', 'email')
                result.append({'priv': can_priv, 'user': list(data)})
        return JsonResponse(result, safe=False)


class ContactsInfoView(View):
    def post(self, request):
        """ 获取指定项目的联系人"""
        group_id = request.POST.get('group_id')

        result = []
        if group_id:

            query = f"select ac.contact_id, group_concat(concat_ws(':', ac.contact_name, ac.contact_id, " \
                    f"ac.contact_email)) as contact_info " \
                    f"from auditsql_contacts as ac JOIN auditsql_contacts_detail a ON ac.contact_id = a.contact_id " \
                    f"JOIN  auditsql_groups a2 " \
                    f"ON a.group_id = a2.group_id where a.group_id = {group_id} group by ac.contact_id;"

            for row in Contacts.objects.raw(query):
                result.append(row.contact_info)

        return JsonResponse(result, safe=False)


class SyntaxCheckView(View):
    """语法检查"""

    def post(self, request):
        form = SyntaxCheckForm(request.POST)
        if form.is_valid():
            cleaned_data = form.cleaned_data
            host = cleaned_data['host']
            database = cleaned_data['database']
            operate_type = cleaned_data.get('operate_type')
            contents = cleaned_data['contents']

            # 对检测的SQL类型进行区分
            filter_result = sql_filter(contents, operate_type)

            # 实例化
            incep_of_audit = IncepSqlCheck(contents, host, database, request.user.username)

            if filter_result['status'] == 2:
                context = filter_result
            else:
                # SQL语法检查
                context = incep_of_audit.run_check()

            return HttpResponse(json.dumps(context))
        else:
            error = "请选择主机、库名或项目组"
            context = {'status': 2, 'msg': error}
            return HttpResponse(json.dumps(context))

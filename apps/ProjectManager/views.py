import json

import sqlparse
from channels.layers import get_channel_layer
from django.db.models import F
from django.http import JsonResponse, HttpResponse
from django.views import View

from ProjectManager.utils import check_mysql_conn
from UserManager.models import GroupsDetail, UserAccount, Contacts
from apps.ProjectManager.inception.inception_api import GetDatabaseListApi
from utils.tools import format_request
from .models import Remark, InceptionHostConfigDetail, InceptionHostConfig

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
                if res[0].tokens[0].ttype[1] == 'DML':
                    sql_format = sqlparse.format(sql, reindent=True)
                    beautify_sql_list.append(comment + sql_format)
                elif res[0].tokens[0].ttype[1] == 'DDL':
                    sql_format = sqlparse.format(sql)
                    beautify_sql_list.append(comment + sql_format)
            context = {'data': '\n\n'.join(beautify_sql_list)}
        except Exception as err:
            raise OSError(err)
            context = {'errCode': 400, 'errMsg': "注释不合法, 请检查"}

        return HttpResponse(json.dumps(context))


class IncepHostConfigView(View):
    """获取指定的数据库配置"""

    def get(self, request):
        config_type = request.GET.get('type')
        user_in_group = self.request.session.get('groups')
        result = InceptionHostConfigDetail.objects.annotate(host=F('config__host'),
                                                            comment=F('config__comment')).filter(
            config__type=config_type).filter(group__group_id__in=user_in_group).values('host', 'comment')
        return JsonResponse(list(result), safe=False)


class GetDBListView(View):
    """列出选中环境的数据库库名"""

    def post(self, request):
        data = format_request(request)
        host = data['host']
        obj = InceptionHostConfig.objects.get(host=host)
        result = check_mysql_conn(obj.user, host, obj.password, obj.port)
        if result['status'] == 'INFO':
            db_list = GetDatabaseListApi(host).get_dbname()
            context = {'errCode': 200, 'errMsg': db_list}
        else:
            context = {'errCode': 400, 'errMsg': f'获取列表失败，不能连接到mysql服务器：{host}'}
        return HttpResponse(json.dumps(context))


class RemarkInfoView(View):
    def post(self, request):
        obj = Remark.objects.all().values('id', 'remark')
        return JsonResponse(list(obj), safe=False)


class GroupInfoView(View):
    def get(self, request):
        groups = GroupsDetail.objects.filter(
            user__uid=request.user.uid).annotate(
            group_id=F('group__group_id'), group_name=F('group__group_name')) \
            .values('group_id', 'group_name')

        return JsonResponse(list(groups), safe=False)


class AuditUserView(View):
    def post(self, request):
        """
        获取指定项目可用的dba和leader信息
        """
        group_id = request.POST.get('group_id')
        result = []
        if group_id:
            data = GroupsDetail.objects.annotate(
                uid=F('user__uid'),
                username=F('user__username'),
                email=F('user__email'),
            ).filter(group__group_id=group_id).values('uid', 'username', 'email')

            for i in data:
                uid = i['uid']
                user_role = UserAccount.objects.get(uid=uid).user_role()
                i['user_role'] = user_role
                result.append(i)

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

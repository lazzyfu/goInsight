import json

from django.db import transaction
from django.http import HttpResponse
from django.shortcuts import render
# Create your views here.
from django.utils.decorators import method_decorator
from django.views import View

from project_manager.models import InceptionHostConfig
from user_manager.permissions import permission_required
from mstats.forms import PrivModifyForm
from mstats.utils import get_mysql_user_info, check_mysql_conn_status, MySQLuser_manager
from utils.tools import format_request


class RenderMySQLUserView(View):
    @permission_required('can_mysqluser_view')
    def get(self, request):
        return render(request, 'mysql_user_manager.html')


class MySQLUserView(View):
    @permission_required('can_mysqluser_view')
    @method_decorator(check_mysql_conn_status)
    def get(self, request):
        data = format_request(request)
        host = data.get('host')
        data = get_mysql_user_info(host)

        return HttpResponse(json.dumps(data))


class MysqlUserManager(View):
    @permission_required('can_mysqluser_edit')
    @transaction.atomic
    def post(self, request):
        data = format_request(request)
        form = PrivModifyForm(data)
        context = {}
        if form.is_valid():
            cleaned_data = form.cleaned_data
            db_host = cleaned_data.get('db_host')
            user = cleaned_data.get('user')
            action = cleaned_data.get('action')

            host = data.get('host')
            password = data.get('password')
            schema = data.get('schema')
            privileges = data.get('privileges')

            username = user + '@' + '"' + host + '"'

            data = InceptionHostConfig.objects.get(host=db_host)
            protection_user = []
            if len(list(data.protection_user.split(','))) == 1:
                protection_user = data.protection_user.split(',')
                protection_user.append('')
            else:
                protection_user = data.protection_user.split(',')
            protection_user_tuple = tuple([x.strip() for x in protection_user])

            if user in protection_user_tuple:
                context = {'status': 1, 'msg': f'该用户({user})已被保护，无法操作'}
            else:
                mysql_user_mamager = MySQLuser_manager(locals())
                if action == "modify_privileges":
                    context = mysql_user_mamager.priv_modify()
                elif action == "new_host":
                    context = mysql_user_mamager.new_host()
                elif action == 'delete_host':
                    context = mysql_user_mamager.delete_host()
                elif action == 'new_user':
                    context = mysql_user_mamager.new_host()

            return HttpResponse(json.dumps(context))

        else:
            error = form.errors.as_text()
            context = {'status': 2, 'msg': error}
            return HttpResponse(json.dumps(context))

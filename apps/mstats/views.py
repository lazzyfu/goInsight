import json

from django.http import HttpResponse
from django.shortcuts import render
# Create your views here.
from django.utils.decorators import method_decorator
from django.views import View

from UserManager.permissions import check_dba_permission
from mstats.forms import PrivModifyForm
from mstats.utils import get_mysql_user_info, check_mysql_conn_status, MySQLUserManager
from utils.tools import format_request


class RenderMySQLUserManagerView(View):
    @method_decorator(check_dba_permission)
    def get(self, request):
        return render(request, 'mysql_user_manager.html')


class MySQLUserManagerView(View):
    @method_decorator(check_mysql_conn_status)
    @method_decorator(check_dba_permission)
    def get(self, request):
        data = format_request(request)
        host = data.get('host')
        data = get_mysql_user_info(host)

        return HttpResponse(json.dumps(data))


class MySQLPrivModifyView(View):
    def post(self, request):
        data = format_request(request)
        form = PrivModifyForm(data)
        if form.is_valid():
            cleaned_data = form.cleaned_data
            db_host = cleaned_data.get('db_host')
            user = cleaned_data.get('user')
            host = cleaned_data.get('host')
            action = cleaned_data.get('action')
            password = data.get('password')
            schema = cleaned_data.get('schema')
            privileges = cleaned_data.get('privileges')

            username = user + '@' + '"' + host + '"'

            mysql_user_mamager = MySQLUserManager(locals())

            context = {}
            if action == "modify_privileges":
                context = mysql_user_mamager.priv_modify()
            elif action == "new_host":
                context = mysql_user_mamager.new_host()

            return HttpResponse(json.dumps(context))

        else:
            error = form.errors.as_text()
            context = {'status': 1, 'msg': error}
            return HttpResponse(json.dumps(context))

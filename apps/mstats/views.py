import json

from django.http import HttpResponse
from django.shortcuts import render
# Create your views here.
from django.utils.decorators import method_decorator
from django.views import View

from UserManager.permissions import check_dba_permission
from mstats.utils import get_mysql_user_info, check_mysql_conn_status
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

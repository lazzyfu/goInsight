import json
import re

import sqlparse
from django.http import JsonResponse, HttpResponse
from django.shortcuts import render
from django.views import View
from django.views.generic import FormView

from ProjectManager.forms import InceptionSqlOperateForm
from apps.ProjectManager.inception.inception_api import GetDatabaseApi, InceptionApi
from utils.tools import format_request
from .models import InceptionHostConfig


class ProjectListView(View):
    def get(self, request):
        return render(request, 'index.html')


class InceptionSqlOperateView(FormView):
    """inception处理DML语句"""

    form_class = InceptionSqlOperateForm
    template_name = 'inception_sql_operate.html'

    def form_valid(self, form):
        cleaned_data = form.cleaned_data
        host = cleaned_data['host']
        database = cleaned_data['database']
        op_action = cleaned_data.get('op_action')
        op_type = cleaned_data['op_type']
        sql_content = cleaned_data['sql_content']

        DDL_FILTER = 'ALTER TABLE|CREATE TABLE|TRUNCATE TABLE'
        DML_FILTER = 'INSERT INTO|;UPDATE|^UPDATE|DELETE FROM'

        checkData = InceptionApi().sqlprepare(sqlcontent=sql_content, host=host, database=database,
                                              action='check')

        # 修改表结构
        if op_action == 'op_schema':
            if re.search(DML_FILTER, sql_content, re.I):
                context = {'errMsg': f'DDL模式下, 不支持SELECT|UPDATE|DELETE|INSERT语句', 'errCode': 400}
            else:
                if op_type == 'check':
                    context = {'data': checkData, 'errCode': 200}
                if op_type == 'commit':
                    if 1 in [x['errlevel'] for x in checkData] or 2 in [x['errlevel'] for x in checkData]:
                        context = {'errMsg': 'SQL语法检查未通过, 请执行语法检测', 'errCode': 400}
                    else:
                        executeData = InceptionApi().sqlprepare(sqlcontent=sql_content, host=host,
                                                                database=database,
                                                                action='execute')
                        context = form.is_save(self.request, executeData)

        # 修改数据
        if op_action == 'op_data':
            if re.search(DDL_FILTER, sql_content, re.I):
                context = {'errMsg': f'DML模式下, 不支持ALTER|CREATE|TRUNCATE语句', 'errCode': 400}
            else:
                if op_type == 'check':
                    context = {'data': checkData, 'errCode': 200}
                if op_type == 'commit':
                    if 1 in [x['errlevel'] for x in checkData] or 2 in [x['errlevel'] for x in checkData]:
                        context = {'errMsg': 'SQL语法检查未通过, 请执行语法检测', 'errCode': 400}
                    else:
                        executeData = InceptionApi().sqlprepare(sqlcontent=sql_content, host=host,
                                                                database=database,
                                                                action='execute')
                        context = form.is_save(self.request, executeData)

        return HttpResponse(json.dumps(context))

    def form_invalid(self, form):
        # error = form.errors.as_text()
        error = "请选择主机或库名"
        context = {'errCode': '400', 'errMsg': error}

        return HttpResponse(json.dumps(context))


class BeautifySQLView(View):
    """美化SQL"""

    def post(self, request):
        data = format_request(request)
        sqlContent = data['sql_content'].rstrip()
        sqlFormat = '\n'.join([line for line in sqlContent.split('\n') if line != ''])
        beautifySQL = sqlparse.format(sqlFormat, keyword_case='upper')
        context = {'data': beautifySQL}
        return HttpResponse(json.dumps(context))


class GetInceptionHostConfigView(View):
    def get(self, request):
        envResult = InceptionHostConfig.objects.all().values('host', 'comment')
        return JsonResponse(list(envResult), safe=False)


class GetDatabaseListView(View):
    def post(self, request):
        data = format_request(request)
        host = data['host']
        dbResult = GetDatabaseApi(host).get_dbname()
        return HttpResponse(json.dumps(dbResult))

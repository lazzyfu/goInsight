import json
import re
from ast import literal_eval

import sqlparse
from django.http import JsonResponse, HttpResponse
from django.shortcuts import render
from django.views import View
from datetime import datetime

from ProjectManager.forms import InceptionSqlOperateForm
from apps.ProjectManager.inception.inception_api import GetDatabaseApi, InceptionApi
from .models import InceptionHostConfig, InceptionSqlOperateRecord
from utils.tools import format_request


class ProjectListView(View):
    def get(self, request):
        return render(request, 'index.html')


class InceptionSqlOperateView(View):
    def get(self, request):
        return render(request, 'inception_sql_operate.html')

class InceptionSqlOperateDDLView(View):
    """inception处理DDL语句"""

    def post(self, request):
        data = format_request(request)
        form = InceptionSqlOperateForm(data)

        if form.is_valid():
            cleaned_data = form.cleaned_data

            host = cleaned_data['host']
            database = cleaned_data['database']
            op_type = cleaned_data['op_type']
            sql_content = cleaned_data['sql_content'].rstrip()
            op_user = request.user.username
            op_uid = request.user.uid
            workid = datetime.now().strftime("%Y%m%d%H%M%S%f")

            sql_filter = sql_content.replace('\n', '')
            DML_FILTER = 'INSERT INTO|;UPDATE|^UPDATE|DELETE FROM'

            if sql_content[-1] != ';':
                context = {'errMsg': 'SQL语句没用以;结尾, 请重新输入', 'errCode': 400}
            else:
                if re.search(DML_FILTER, sql_filter, re.I):
                    context = {'errMsg': f'DDL模式下, 不支持SELECT|UPDATE|DELETE|INSERT语句', 'errCode': 400}
                else:
                    data = InceptionApi().sqlprepare(sqlcontent=sql_content, host=host, database=database,
                                                     action='check')
                    if op_type == 'check':
                        context = {'data': data, 'errCode': 200}
                    elif op_type == 'commit':
                        if 1 in [x['errlevel'] for x in data] or 2 in [x['errlevel'] for x in data]:
                            context = {'errMsg': 'SQL语法检查未通过, 请执行语法检测', 'errCode': 400}
                        else:
                            checkData = InceptionApi().sqlprepare(sqlcontent=sql_content, host=host,
                                                                  database=database,
                                                                  action='execute')
                            for line in checkData:
                                sql = line['SQL']
                                step_id = line['ID']
                                stage = line['stage']
                                stagestatus = line['stagestatus']
                                errlevel = line['errlevel']
                                errormessage = line['errormessage']
                                Affected_rows = line['Affected_rows']
                                sequence = line['sequence']
                                backup_dbname = line['backup_dbname']
                                execute_time = line['execute_time']

                                InceptionSqlOperateRecord.objects.create(
                                    op_user=op_user,
                                    op_uid=op_uid,
                                    workid=workid,
                                    dst_host=host,
                                    dst_database=database,
                                    stagestatus=stagestatus,
                                    sql=sql,
                                    step_id=step_id,
                                    stage=stage,
                                    errlevel=errlevel,
                                    errormessage=errormessage,
                                    affected_rows=Affected_rows,
                                    sequence=literal_eval(sequence),
                                    backup_dbname=backup_dbname,
                                    execute_time=execute_time
                                )
                            context = {'data': checkData, 'errMsg': '执行完成', 'errCode': 200}
        else:
            # error = form.errors.as_text()
            error = "请选择主机或库名"
            context = {'errCode': '400', 'errMsg': error}
        return HttpResponse(json.dumps(context))


class InceptionSqlOperateDMLView(View):
    """inception DML check and output"""

    def get(self, request):
        return render(request, 'inception_sql_operate.html')


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

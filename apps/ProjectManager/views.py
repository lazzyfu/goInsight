import json
import re

import sqlparse
from django.http import JsonResponse, HttpResponse
from django.shortcuts import render, get_object_or_404
from django.views import View
from django.views.generic import FormView, ListView
from pure_pagination import PaginationMixin

from ProjectManager.forms import InceptionSqlOperateForm
from apps.ProjectManager.inception.inception_api import GetDatabaseListApi, InceptionApi, GetBackupApi
from utils.tools import format_request
from .models import InceptionHostConfig, InceptionSqlOperateRecord


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
    """获取inception审核的目标数据库配置"""
    def get(self, request):
        envResult = InceptionHostConfig.objects.all().values('host', 'comment')
        return JsonResponse(list(envResult), safe=False)


class GetDatabaseListView(View):
    """列出选中环境的数据库库名"""
    def post(self, request):
        data = format_request(request)
        host = data['host']
        dbResult = GetDatabaseListApi(host).get_dbname()
        return HttpResponse(json.dumps(dbResult))

class InceptionSqlRecords(PaginationMixin, ListView):
    """查看用户的工单记录"""
    paginate_by = 8
    context_object_name = 'sqlRecord'
    template_name = 'inception_sql_records.html'

    def get_queryset(self):
        workidQuery = "select workid,id,op_user,dst_host,op_time from sqlaudit_inception_sql_operate_record group by workid order by op_time desc"
        sqlRecord = []
        for row in InceptionSqlOperateRecord.objects.raw(workidQuery):
            workid = row.workid
            op_user = row.op_user
            dst_host = row.dst_host
            op_time = row.op_time
            singleRecord = InceptionSqlOperateRecord.objects.filter(op_uid=self.request.user.uid).filter(workid=workid).order_by(
                'op_time')
            sqlRecord.append({'workid': workid, 'op_user': op_user, 'dst_host': dst_host, 'op_time': op_time,
                              'record': singleRecord})
        return sqlRecord


class InceptionAllSqlDetailView(View):
    """查看当前用户会话执行的所有sql的详情"""
    def get(self, request, workid):
        sequenceResult = []
        originalSql = ''
        originalSqlQuery = InceptionSqlOperateRecord.objects.raw(
            f"select id,group_concat(`op_sql` separator '\n') as `op_sql` from sqlaudit_inception_sql_operate_record where workid={workid} group by workid")
        for i in originalSqlQuery:
            originalSql = i.op_sql

        sqlDetail = InceptionSqlOperateRecord.objects.filter(workid=workid)
        for row in sqlDetail:
            sequenceResult.append({'backupdbName': row.backup_dbname, 'sequence': row.sequence})
        rollbackSql = GetBackupApi(sequenceResult).get_backupinfo()

        return render(request, 'offline_all_sql_detail.html',
                      {'originalSql': originalSql, 'rollbackSql': rollbackSql})


class InceptionSingleSqlDetailView(View):
    """查看当前用户会话执行的每条sql的详情"""
    def get(self, request, sequence):
        sqlDetail = get_object_or_404(InceptionSqlOperateRecord, sequence=sequence)
        sequenceResult = [{'backupdbName': sqlDetail.backup_dbname, 'sequence': sqlDetail.sequence}]
        rollbackSql = GetBackupApi(sequenceResult).get_backupinfo()
        return render(request, 'offline_single_sql_detail.html',
                      {'sqlDetail': sqlDetail, 'rollbackSql': rollbackSql})
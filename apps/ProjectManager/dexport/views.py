# -*- coding:utf-8 -*-
# edit by fuzongfei
import json
from datetime import datetime

from django.db.models import Case, When, Value, CharField
from django.http import HttpResponse, JsonResponse
from django.shortcuts import render
from django.utils.decorators import method_decorator
from django.views import View

from ProjectManager.models import DataExport, Files
from ProjectManager.permissions import check_data_export_permission
from ProjectManager.tasks import make_export_file, send_data_export_mail
from utils.tools import format_request


class DataExportView(View):
    def get(self, request):
        return render(request, 'data_export.html')

    def post(self, request):
        data = format_request(request)

        title = '__'.join((data['title'], datetime.now().strftime("%Y%m%d%H%M%S")))

        DataExport.objects.create(
            proposer=request.user.username,
            group_id=data['group_id'],
            dst_host=data['host'],
            dst_database=data['database'],
            title=title,
            file_format=data['file_format'],
            file_coding=data['file_coding'],
            operate_dba=data['operate_dba'],
            email_cc=data['email_cc_id'],
            sql_contents=data['sql_content']
        )
        latest_id = DataExport.objects.latest('id').id
        send_data_export_mail.delay(latest_id=latest_id)
        context = {'errCode': 200, 'errMsg': '提交成功'}
        return HttpResponse(json.dumps(context))


class DataExportRecordsView(View):
    def get(self, request):
        return render(request, 'data_export_records.html')


class DataExportRecordsListView(View):
    def get(self, request):
        user_in_group = request.session['groups']

        records = DataExport.objects.all().annotate(
            exec_status=Case(
                When(status='0', then=Value('未生成')),
                When(status='1', then=Value('执行中')),
                When(status='2', then=Value('已生成')),
                output_field=CharField(),
            ),
            group_name=F('group__group_name'),
        ).filter(group_id__in=user_in_group).values('id', 'exec_status', 'group_name', 'title', 'dst_host',
                                                    'dst_database',
                                                    'proposer', 'operate_dba', 'file_coding', 'file_format',
                                                    'sql_contents', 'created_at').order_by('-created_at')
        return JsonResponse(list(records), safe=False)


class ExecDataExportView(View):
    """生成导出文件"""

    @method_decorator(check_data_export_permission)
    def post(self, request):
        id = request.POST.get('id')

        if DataExport.objects.get(pk=id).status in ('1', '2'):
            context = {'errCode': 400, 'errMsg': '数据正在执行或已完成，请不要重复操作'}
        else:
            make_export_file.delay(user=request.user.username, id=id)

            DataExport.objects.filter(pk=id).update(status='1')

            context = {'errCode': 200, 'errMsg': '已提交处理，请稍后'}
        return HttpResponse(json.dumps(context))


class DataExportDownloadView(View):
    def get(self, request):
        id = request.GET.get('id')

        if DataExport.objects.get(pk=id).status == '2':
            obj = Files.objects.get(export_id=id)
            file_size = str(round(obj.file_size / 1024, 2)) + 'KB'
            file_name = obj.file_name
            file_path = obj.files.url
            encryption_key = obj.encryption_key

            context = {'errCode': 200, 'file_size': file_size, 'file_path': file_path, 'file_name': file_name,
                       'encryption_key': encryption_key}
        else:
            context = {'errCode': 400, 'errMsg': '文件未生成，无法下载'}
        return HttpResponse(json.dumps(context))

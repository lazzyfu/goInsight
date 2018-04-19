# -*- coding:utf-8 -*-
# edit by fuzongfei
import json
from datetime import datetime

from django.db import transaction
from django.db.models import Case, When, Value, CharField, F
from django.http import HttpResponse, JsonResponse
from django.shortcuts import render
from django.utils.decorators import method_decorator
from django.views import View

from project_manager.dexport.forms import OlDataExportForm
from project_manager.models import OlDataExportDetail
from project_manager.tasks import make_export_file, send_data_export_mail
from user_manager.permissions import permission_required
from utils.tools import format_request


class OlDataExportView(View):
    def get(self, request):
        return render(request, 'data_export.html')

    @permission_required('can_commit')
    @transaction.atomic
    def post(self, request):
        data = format_request(request)
        form = OlDataExportForm(data)
        if form.is_valid():
            context = form.is_save(request)
        else:
            error = form.errors.as_text()
            context = {'status': 2, 'msg': error}
        return HttpResponse(json.dumps(context))


class OlDataExportDetailRecordsView(View):
    def get(self, request):
        return render(request, 'data_export_records.html')


class OlDataExportDetailRecordsListView(View):
    def get(self, request):
        user_in_group = request.session['groups']

        records = OlDataExportDetail.objects.all().annotate(
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


class ExecOlDataExportDetailView(View):
    """生成导出文件"""

    @permission_required('can_export')
    @transaction.atomic
    def post(self, request):
        id = request.POST.get('id')

        if OlDataExportDetail.objects.get(pk=id).status in ('1', '2'):
            context = {'status': 2, 'msg': '数据正在执行或已完成，请不要重复操作'}
        else:
            make_export_file.delay(user=request.user.username, id=id)

            OlDataExportDetail.objects.filter(pk=id).update(status='1')

            context = {'status': 0, 'msg': '已提交处理，请稍后'}
        return HttpResponse(json.dumps(context))


class OlDataExportDetailDownloadView(View):
    def get(self, request):
        id = request.GET.get('id')

        if OlDataExportDetail.objects.get(pk=id).status == '2':
            obj = Files.objects.get(export_id=id)
            file_size = str(round(obj.file_size / 1024, 2)) + 'KB'
            file_name = obj.file_name
            file_path = obj.files.url
            encryption_key = obj.encryption_key

            context = {'status': 0, 'msg': '',
                       'data': {'file_size': file_size, 'file_path': file_path, 'file_name': file_name,
                                'encryption_key': encryption_key}}
        else:
            context = {'status': 2, 'msg': '文件未生成，无法下载'}
        return HttpResponse(json.dumps(context))

# -*- coding:utf-8 -*-
# edit by fuzongfei
import json

from django.db import transaction
from django.http import HttpResponse
from django.shortcuts import render
from django.views import View

from project_manager.dexport.forms import OlDataExportForm
from project_manager.models import OlDataExportDetail, AuditContents
from project_manager.tasks import make_export_file
from user_manager.permissions import permission_required
from utils.tools import format_request


class OlDataExportView(View):
    """数据导出提交处理"""
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


class ExecOlDataExportView(View):
    """生成导出文件"""

    @permission_required('can_export')
    @transaction.atomic
    def post(self, request):
        id = request.POST.get('id')

        data = AuditContents.objects.get(pk=id)
        detail = OlDataExportDetail.objects.get(ol=id)

        # 当记录关闭时
        if data.progress == '5':
            context = {'status': 2, 'msg': '该记录已被关闭、请不要重复提交'}
        # 当记录未关闭时
        else:
            # 已批准
            if data.progress == '2':
                if detail.progress == '0':
                    make_export_file.delay(user=request.user.username, id=id)

                    # 更新进度为：导出中
                    OlDataExportDetail.objects.filter(ol=id).update(progress='1')

                    context = {'status': 0, 'msg': '已提交处理，请稍后'}
                else:
                    context = {'status': 2, 'msg': '操作失败、任务正在处理或已完成'}

            # 未批准
            elif data.progress == '1' or data.progress == '0':
                context = {'status': 2, 'msg': '操作失败、审核未通过'}
            # 其他情况
            else:
                context = {'status': 2, 'msg': '操作失败、请不要重复提交'}

        return HttpResponse(json.dumps(context))
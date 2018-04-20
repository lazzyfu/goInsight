# -*- coding:utf-8 -*-
# edit by fuzongfei
import json

from django.db import transaction
from django.http import HttpResponse
from django.shortcuts import render
from django.utils.decorators import method_decorator
from django.views import View

from project_manager.of.forms import IncepOfAuditForm
from project_manager.utils import check_incep_alive, check_sql_filter
from user_manager.permissions import permission_required
from utils.tools import format_request


class IncepOfAuditView(View):
    """线下工单生成执行任务"""
    def get(self, request):
        return render(request, 'incep_of_audit.html')

    @method_decorator(check_incep_alive)
    @method_decorator(check_sql_filter)
    @permission_required('can_commit')
    @transaction.atomic
    def post(self, request):
        data = format_request(request)
        form = IncepOfAuditForm(data)
        if form.is_valid():
            context = form.save(request)
            return HttpResponse(json.dumps(context))
        else:
            error = form.errors.as_text()
            context = {'status': 2, 'msg': error}

            return HttpResponse(json.dumps(context))

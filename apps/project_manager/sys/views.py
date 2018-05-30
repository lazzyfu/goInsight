# -*- coding:utf-8 -*-
# edit by fuzongfei
import json

from django.http import HttpResponse, JsonResponse
from django.shortcuts import render
from django.views import View

from project_manager.models import DomainName, InceptionHostConfig
from user_manager.permissions import permission_required
from utils.tools import format_request


class SysConfigView(View):
    @permission_required('can_admin')
    def get(self, request):
        return render(request, 'sys_config.html')


class GetDomainView(View):
    @permission_required('can_admin')
    def get(self, request):
        domain_name = 'None'
        if DomainName.objects.filter().first():
            domain_name = DomainName.objects.get().domain_name
        result = {'status': 0, 'data': domain_name}
        return HttpResponse(json.dumps(result))

    @permission_required('can_admin')
    def post(self, request):
        data = format_request(request)
        domain_name = data.get('domain_name')
        if DomainName.objects.filter():
            DomainName.objects.update(domain_name=domain_name)
        else:
            DomainName.objects.create(domain_name=domain_name)
        result = {'status': 0, 'msg': '域名修改成功'}
        return HttpResponse(json.dumps(result))


class GetDBAccountView(View):
    @permission_required('can_admin')
    def get(self, request):
        data = InceptionHostConfig.objects.all().values('user', 'password', 'host', 'port', 'type', 'is_enable',
                                                        'protection_user', 'comment')
        return JsonResponse(list(data), safe=False)

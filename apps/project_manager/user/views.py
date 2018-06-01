# -*- coding:utf-8 -*-
# edit by fuzongfei
import json

from django.db import transaction
from django.http import HttpResponse, JsonResponse
from django.shortcuts import render
from django.views import View

from project_manager.models import DomainName, InceptionHostConfig, Webhook
from project_manager.utils import check_db_account
from user_manager.models import Permission, Groups, Roles
from user_manager.permissions import permission_required
from utils.tools import format_request


class UserConfigView(View):
    @permission_required('can_admin')
    def get(self, request):
        return render(request, 'user_config.html')


class GetUserPermissionView(View):
    @permission_required('can_admin')
    def get(self, request):
        data = Permission.objects.all().values('id', 'permission_name', 'permission_desc', 'created_at')
        return JsonResponse(list(data), safe=False)


class GetProjectView(View):
    @permission_required('can_admin')
    def get(self, request):
        data = Groups.objects.all().values('group_id', 'group_name', 'created_at')
        return JsonResponse(list(data), safe=False)

    @permission_required('can_admin')
    def post(self, request):
        data = format_request(request)
        group_name = data.get('group_name')
        del data['csrfmiddlewaretoken']

        if data.get('action') == 'delete_row':
            for i in data.get('group_id').split(','):
                Groups.objects.get(pk=i).delete()
            context = {'status': 0, 'msg': '删除成功'}

        elif data.get('action') == 'modify_value':
            del data['0']
            del data['action']
            Groups.objects.filter(group_id=data.get('group_id')).update(**data)
            context = {'status': 0, 'msg': '修改成功'}

        elif data.get('action') == 'new_row':
            del data['action']
            Groups.objects.create(group_name=group_name)
            context = {'status': 0, 'msg': '创建成功'}

        return HttpResponse(json.dumps(context))


class GetRoleView(View):
    @permission_required('can_admin')
    def get(self, request):
        query = "select a.role_id, a.role_name, GROUP_CONCAT(c.permission_desc) as permission_desc " \
                "from auditsql_roles a left join auditsql_permission_detail b on a.role_id = b.role_id " \
                "left join auditsql_permission c on b.permission_id = c.id group by a.role_name, a.role_id"

        data = []
        for row in Roles.objects.raw(query):
            data.append({
                'role_id': row.role_id,
                'role_name': row.role_name,
                'permission_desc': row.permission_desc
            })

        return JsonResponse(list(data), safe=False)

    @permission_required('can_admin')
    def post(self, request):
        data = format_request(request)
        role_name = data.get('role_name')
        del data['csrfmiddlewaretoken']

        if data.get('action') == 'delete_row':
            for i in data.get('role_id').split(','):
                Roles.objects.get(pk=i).delete()
            context = {'status': 0, 'msg': '删除成功'}
        elif data.get('action') == 'modify_value':
            del data['0']
            del data['action']
            Roles.objects.filter(role_id=data.get('role_id')).update(**data)
            context = {'status': 0, 'msg': '修改成功'}
        elif data.get('action') == 'new_row':
            del data['action']
            Roles.objects.create(role_name=role_name)
            context = {'status': 0, 'msg': '创建成功'}

        return HttpResponse(json.dumps(context))

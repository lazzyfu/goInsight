# -*- coding:utf-8 -*-
# edit by fuzongfei
import json

from channels.layers import get_channel_layer

from django.core.exceptions import PermissionDenied
from django.http import HttpResponse
from django.shortcuts import get_object_or_404
from ProjectManager.models import OnlineAuditContents, IncepMakeExecTask

channel_layer = get_channel_layer()


def check_group_permission(fun):
    """
    验证用户是否属于指定的项目组
    如果用户不属于该项目，则返回：PermissionDenied
    """

    def wapper(request, *args, **kwargs):
        user_in_group = request.session['groups']
        group_id = request.POST.get('group_id')

        if len(user_in_group) > 0:
            if int(group_id) in request.session['groups']:
                return fun(request, *args, **kwargs)
            else:
                context = {'errCode': '403', 'errMsg': '权限拒绝，您不属于该项目组的成员'}
        else:
            raise PermissionDenied
        return HttpResponse(json.dumps(context))

    return wapper


def check_sql_detail_permission(fun):
    """
    :param fun: request
    :return: 验证用户是否有指定项目详情记录的访问权限
    """

    def wapper(request, *args, **kwargs):
        id = kwargs['id']
        group_id = int(kwargs['group_id'])

        # 检查该记录是否存在
        obj = get_object_or_404(OnlineAuditContents, pk=id)

        # 检查用户是否有该项目的权限
        if group_id not in request.session['groups']:
            raise PermissionDenied

        # 验证pk记录中的group_id是否和输入的group_id相同
        if obj.group_id == group_id:
            return fun(request, *args, **kwargs)
        else:
            raise PermissionDenied

    return wapper


def check_incep_tasks_permission(fun):
    """
    只要DBA角色的用户，才能操作线上执行任务
    """

    def wapper(request, *args, **kwargs):
        id = request.POST.get('id')
        category = IncepMakeExecTask.objects.get(pk=id).category
        user_role = request.user.user_role()
        if category == '1' and user_role == 'DBA':
            return fun(request, *args, **kwargs)
        if category == '0':
            return fun(request, *args, **kwargs)
        else:
            # raise PermissionDenied
            context = {'errCode': 400, 'errMsg': '权限拒绝，只要DBA可以操作'}
            return HttpResponse(json.dumps(context))

    return wapper


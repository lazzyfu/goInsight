# -*- coding:utf-8 -*-
# edit by fuzongfei
import json
from functools import wraps

from django.core.exceptions import PermissionDenied
from django.db.models import F
from django.http import HttpResponse
from django.shortcuts import get_object_or_404

from ProjectManager.models import OnlineAuditContents
from UserManager.models import PermissionDetail


def has_perm(perm_list, perm):
    if perm in perm_list:
        return True
    else:
        return False


def permission_required(*perm):
    """
    rewrite permission required
    不使用系统自带的permission_required
    使用方法：permission_required('can_view'), permission_required('can_view', 'can_edit')
    """

    def decorator(view_func):
        @wraps(view_func)
        def _wrapped_view(request, *args, **kwargs):
            if isinstance(perm, str):
                perms = (perm,)
            else:
                perms = perm
            user_role = request.request.user.user_role()
            perm_list = list(PermissionDetail.objects.annotate(
                permission_name=F('permission__permission_name')).filter(role__role_name=user_role).values_list(
                'permission_name', flat=True))
            if any(has_perm(perm_list, p) for p in perms):
                return view_func(request, *args, **kwargs)
            else:
                raise PermissionDenied

        return _wrapped_view

    return decorator


def check_group_permission(fun):
    """
    验证项目组权限
    如果用户不属于该项目，则返回：PermissionDenied
    """

    def wapper(request, *args, **kwargs):
        user_in_group = request.session['groups']
        group_id = request.POST.get('group_id')

        if user_in_group is not None:
            if int(group_id) in request.session['groups']:
                return fun(request, *args, **kwargs)
            else:
                context = {'status': 1, 'msg': '权限拒绝，您不属于该项目组的成员'}
                return HttpResponse(json.dumps(context))
        else:
            raise PermissionDenied

    return wapper


def check_record_details_permission(fun):
    """
    验证用户是否有指定项目详情记录的访问权限
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

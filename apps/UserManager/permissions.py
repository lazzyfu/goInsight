# -*- coding:utf-8 -*-
# edit by fuzongfei
import json
from functools import wraps

from django.core.exceptions import PermissionDenied
from django.db.models import F
from django.http import HttpResponse

from UserManager.models import PermissionDetail


def check_dba_permission(fun):
    """
    只要DBA角色的用户，才能执行生成导出任务
    """

    def wapper(request, *args, **kwargs):
        user_role = request.user.user_role()
        if user_role == 'DBA':
            return fun(request, *args, **kwargs)
        else:
            raise PermissionDenied

    return wapper


def permission_required(perm):
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
            if all(True for p in perms if p in perm_list):
                return view_func(request, *args, **kwargs)
            else:
                context = {'status': 1, 'msg': '权限拒绝, 您没有权限操作'}
                return HttpResponse(json.dumps(context))

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

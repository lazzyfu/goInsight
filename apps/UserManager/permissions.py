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
    """

    def decorator(view_func):
        @wraps(view_func)
        def _wrapped_view(request, *args, **kwargs):
            user_role = request.request.user.user_role()
            permission_list = list(PermissionDetail.objects.annotate(
                permission_name=F('permission__permission_name')).filter(role__role_name=user_role).values_list(
                'permission_name', flat=True))
            if perm in permission_list:
                return view_func(request, *args, **kwargs)
            else:
                context = {'status': 1, 'msg': '权限拒绝'}
                return HttpResponse(json.dumps(context))

        return _wrapped_view

    return decorator

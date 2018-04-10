# -*- coding:utf-8 -*-
# edit by fuzongfei

from django.core.exceptions import PermissionDenied


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

# -*- coding:utf-8 -*-
# edit by fuzongfei
from functools import wraps

from django.core.exceptions import PermissionDenied

from users.models import RolePermission


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
            role_name = request.request.user.user_role()
            perm_list = list(
                RolePermission.objects.filter(role__role_name=role_name).values_list('permission_name', flat=True))
            if any(has_perm(perm_list, p) for p in perms):
                return view_func(request, *args, **kwargs)
            else:
                raise PermissionDenied

        return _wrapped_view

    return decorator

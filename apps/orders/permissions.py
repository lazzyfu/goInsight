# -*- coding:utf-8 -*-
# edit by fuzongfei

from rest_framework import permissions

from users.models import RolePermissions


def anyof(*perm_classess):
    """实现OR权限控制"""
    class Or(permissions.BasePermission):
        def has_permission(*args):
            allowed = [p.has_permission(*args) for p in perm_classess]
            return any(allowed)

    return Or


class CanCommitPermission(permissions.BasePermission):
    """工单提交权限"""

    message = "您没有工单提交权限"

    def has_permission(self, request, view):
        roles = request.user.user_role()
        perm_list = list(
            RolePermissions.objects.filter(role__role_name=roles).values_list('permission_name', flat=True)
        )
        if 'can_commit' not in perm_list:
            return False
        return True


class CanAuditPermission(permissions.BasePermission):
    """工单审核权限"""

    message = "您没有工单审核权限"

    def has_permission(self, request, view):
        roles = request.user.user_role()
        perm_list = list(
            RolePermissions.objects.filter(role__role_name=roles).values_list('permission_name', flat=True)
        )
        if 'can_audit' not in perm_list:
            return False
        return True


class CanExecutePermission(permissions.BasePermission):
    """工单执行权限"""

    message = "您没有工单执行权限"

    def has_permission(self, request, view):
        roles = request.user.user_role()
        perm_list = list(
            RolePermissions.objects.filter(role__role_name=roles).values_list('permission_name', flat=True)
        )
        if 'can_execute' not in perm_list:
            return False
        return True

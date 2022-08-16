# -*- coding:utf-8 -*-
# edit by xff

from rest_framework import permissions

from users.models import RolePermissions

"""
系统默认的权限：
can_commit_orders: 提交工单的权限
can_audit_orders: 审核工单的权限
can_view_orders: 查看工单的权限
can_execute_orders: 执行工单的权限
can_view_versions: 查看上线版本的权限
can_create_versions: 创建上线版本的权限
can_update_versions: 更新上线版本的权限
can_delete_versions: 删除上线版本的权限
"""


def anyof(*perm_classess):
    """实现OR权限控制"""

    class Or(permissions.BasePermission):
        def has_permission(*args):
            allowed = [p.has_permission(*args) for p in perm_classess]
            return any(allowed)

    return Or


def check_permission(request, perm):
    roles = request.user.user_role()
    perm_list = list(
        RolePermissions.objects.filter(role__role_name=roles).values_list('permission_name', flat=True)
    )
    if perm not in perm_list:
        return False
    return True


class CanCommitOrdersPermission(permissions.BasePermission):
    message = "您没有工单提交权限"

    def has_permission(self, request, view):
        return check_permission(request, 'can_commit_orders')


class CanAuditOrdersPermission(permissions.BasePermission):
    message = "您没有工单审核权限"

    def has_permission(self, request, view):
        return check_permission(request, 'can_audit_orders')


class CanViewOrdersPermission(permissions.BasePermission):
    message = "您没有工单查看权限"

    def has_permission(self, request, view):
        return check_permission(request, 'can_view_orders')


class CanExecuteOrdersPermission(permissions.BasePermission):
    message = "您没有工单执行权限"

    def has_permission(self, request, view):
        return check_permission(request, 'can_execute_orders')


class CanViewVersionPermission(permissions.BasePermission):
    message = "您没有上线版本查看权限"

    def has_permission(self, request, view):
        return check_permission(request, 'can_view_versions')


class CanCreateVersionsPermission(permissions.BasePermission):
    message = "您没有上线版本创建权限"

    def has_permission(self, request, view):
        return check_permission(request, 'can_create_versions')


class CanUpdateVersionsPermission(permissions.BasePermission):
    message = "您没有上线版本更新权限"

    def has_permission(self, request, view):
        return check_permission(request, 'can_update_versions')


class CanDeleteVersionsPermission(permissions.BasePermission):
    message = "您没有上线版本删除权限"

    def has_permission(self, request, view):
        return check_permission(request, 'can_delete_versions')


class CanExecRedisPermission(permissions.BasePermission):
    message = "您没有Redis操作权限"

    def has_permission(self, request, view):
        return check_permission(request, 'can_exec_redis')
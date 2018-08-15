# -*- coding:utf-8 -*-
# edit by fuzongfei
import ast
import logging
from functools import wraps

from django.core.exceptions import PermissionDenied
from django.db.models import F
from django.shortcuts import get_object_or_404

from project_manager.models import AuditContents, IncepMakeExecTask
from user_manager.models import PermissionDetail
from utils.tools import format_request

logger = logging.getLogger(__name__)


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


def order_permission_required(perm):
    """
    验证线上工单执行权限
    可执行的用户：有指定的权限或工单申请人
    """

    def decorator(view_func):
        @wraps(view_func)
        def _wrapped_view(request, *args, **kwargs):
            if isinstance(perm, str):
                perms = (perm,)
            else:
                perms = perm
            user_role = request.request.user.user_role()
            id = request.request.POST.get('id')
            # 此处用作debug
            logger.info(f'钩子id: {id}, type: {type(id)}')
            if id != '':
                pk_id = int(id) if isinstance(id, str) else id
            obj = get_object_or_404(AuditContents, pk=pk_id)
            perm_list = list(PermissionDetail.objects.annotate(
                permission_name=F('permission__permission_name')).filter(role__role_name=user_role).values_list(
                'permission_name', flat=True))
            if any(has_perm(perm_list, p) for p in perms) or obj.proposer == request.request.user.username:
                return view_func(request, *args, **kwargs)
            else:
                raise PermissionDenied

        return _wrapped_view

    return decorator


def perform_tasks_permission_required(*perm):
    """
    执行任务权限检查
    可执行的用户：有指定的权限或工单申请人
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

            data = format_request(request.request)
            taskid = ast.literal_eval(data.get('taskid'))

            allowed_user = IncepMakeExecTask.objects.filter(taskid=taskid).first().user
            # 当执行任务为线上任务时或谁提交的，谁有权限执行
            if any(has_perm(perm_list, p) for p in perms) or allowed_user == request.request.user.username:
                return view_func(request, *args, **kwargs)
            else:
                raise PermissionDenied

        return _wrapped_view

    return decorator

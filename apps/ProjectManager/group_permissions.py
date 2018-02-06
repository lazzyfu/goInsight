# -*- coding:utf-8 -*-
# edit by fuzongfei
import json

from django.http import HttpResponse


def check_group_permission(fun):
    """
    验证用户是否属于指定的项目组
    如果用户不属于该项目，则返回：PermissionDenied
    """

    def wapper(request, *args, **kwargs):
        group_id = request.POST.get('group_id')

        if int(group_id) in request.session['groups']:
            return fun(request, *args, **kwargs)
        else:
            context = {'errCode': '403', 'errMsg': '权限拒绝，您不属于该项目组的成员'}
        return HttpResponse(json.dumps(context))

    return wapper

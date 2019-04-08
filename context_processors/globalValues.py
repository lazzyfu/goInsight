# -*- coding:utf-8 -*-
# edit by fuzongfei
from orders.models import SysEnvironment


def get_sys_enviroment(request):
    """返回工单环境"""
    queryset = SysEnvironment.objects.all()
    sys_environment = queryset.values('envi_id', 'envi_name')
    return locals()

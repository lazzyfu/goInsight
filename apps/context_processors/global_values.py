#!/usr/bin/env python3
# -*— coding: utf-8 -*-
# __author__ : pandonglin

from sqlorders.models import SqlOrdersEnvironment, SysConfig


def get_order_enviroment(request):
    """返回工单环境"""
    queryset = SqlOrdersEnvironment.objects.all()
    order_environment = queryset.values('envi_id', 'envi_name')
    return locals()


def get_mail_status(request):
    """是否启用邮件功能"""
    is_enable_mail = SysConfig.objects.get(key='email_push').is_enabled
    return locals()

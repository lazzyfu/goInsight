# -*- coding:utf-8 -*-
# edit by fuzongfei

from django.contrib.auth.decorators import login_required
from django.urls import path

from project_manager.of.views import OfflineWorkOrderAuditView, MOneWorkOrderAuditView, HookWorkOrderView

urlpatterns = [
    # 线下环境工单提交页面
    path(r'of_work_order_audit/', login_required(OfflineWorkOrderAuditView.as_view()), name='p_of_work_order_audit'),
    # 测试环境工单提交页面
    path(r'test_work_order_audit/', login_required(MOneWorkOrderAuditView.as_view()), name='p_test_work_order_audit'),
    # 钩子
    path(r'hook_work_order/', login_required(HookWorkOrderView.as_view()), name='p_hook_work_order'),
]

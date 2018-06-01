# -*- coding:utf-8 -*-
# edit by fuzongfei

from django.contrib.auth.decorators import login_required
from django.urls import path, include

from .views import BeautifySQLView, IncepHostConfigView, \
    GetSchemaView, GroupInfoView, GetAuditUserView, ContactsInfoView, SyntaxCheckView

urlpatterns = [
    path(r'host_config/', login_required(IncepHostConfigView.as_view()),
         name='p_host_config'),
    path(r'schema_info/', login_required(GetSchemaView.as_view())),
    path(r'group_info/', login_required(GroupInfoView.as_view())),
    path(r'get_audit_user/', login_required(GetAuditUserView.as_view())),
    path(r'contacts_info/', login_required(ContactsInfoView.as_view())),
    path(r'beautify_sql/', login_required(BeautifySQLView.as_view())),
    path(r'syntax_check/', login_required(SyntaxCheckView.as_view())),

    # 线上工单审核
    path(r'ol/', include('project_manager.ol.urls')),
    # 线下工单审核
    path(r'ol/', include('project_manager.of.urls')),
    # 执行任务
    path(r'pt/', include('project_manager.pt.urls')),
    # 系统配置
    path(r'sys/', include('project_manager.sys.urls')),
    # 用户配置
    path(r'user/', include('project_manager.user.urls')),
]

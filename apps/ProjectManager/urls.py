# -*- coding:utf-8 -*-
# edit by fuzongfei

from django.contrib.auth.decorators import login_required
from django.urls import path, include

from .views import BeautifySQLView, IncepHostConfigView, \
    GetDBListView, \
    RemarkInfoView, GroupInfoView, AuditUserView, ContactsInfoView

urlpatterns = [
    path(r'host_config/', login_required(IncepHostConfigView.as_view()),
         name='p_host_config'),
    path(r'db_list/', login_required(GetDBListView.as_view()), name='p_db_list'),
    path(r'remark_info/', login_required(RemarkInfoView.as_view()), name='p_remark_info'),
    path(r'group_info/', login_required(GroupInfoView.as_view()), name='p_group_info'),
    path(r'audit_user/', login_required(AuditUserView.as_view()), name='p_audit_user'),
    path(r'contacts_info/', login_required(ContactsInfoView.as_view()), name='p_contacts_info'),
    path(r'beautify_sql/', login_required(BeautifySQLView.as_view()), name='p_beautify_sql'),

    # 线上工单审核
    path(r'ol/', include('ProjectManager.ol.urls')),
    # 线下工单审核
    path(r'ol/', include('ProjectManager.of.urls')),
    # 执行任务
    path(r'pt/', include('ProjectManager.pt.urls')),
    # 数据导出
    path(r'de/', include('ProjectManager.dexport.urls')),
]

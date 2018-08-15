# -*- coding:utf-8 -*-
# edit by fuzongfei

from django.contrib.auth.decorators import login_required
from django.urls import path, include

from mstats.views import GetOnlineQuerySchemaTreeView, GetOfflineQuerySchemaTreeView
from .views import BeautifySQLView, GetAuditUserView, \
    SyntaxCheckView, GetOnlineAuditTasksList, \
    GetTableInfo, GetHostInfoView, GetTableMetaInfoView

urlpatterns = [
    path(r'get_host_info/', login_required(GetHostInfoView.as_view()), name='p_host_info'),
    path(r'get_table_info/', login_required(GetTableInfo.as_view())),
    path(r'get_table_meta_info/', login_required(GetTableMetaInfoView.as_view())),
    # 生产环境mysql查询
    path(r'get_online_query_schema_tree/', login_required(GetOnlineQuerySchemaTreeView.as_view()),
         name='p_get_online_query_schema_tree'),
    # 非生产环境mysql查询
    path(r'get_offline_query_schema_tree/', login_required(GetOfflineQuerySchemaTreeView.as_view()),
         name='p_get_offline_query_schema_tree'),
    path(r'get_audit_user/', login_required(GetAuditUserView.as_view())),
    path(r'beautify_sql/', login_required(BeautifySQLView.as_view())),
    path(r'syntax_check/', login_required(SyntaxCheckView.as_view())),
    path(r'get_tasks/', login_required(GetOnlineAuditTasksList.as_view())),

    # 线上工单审核
    path(r'ol/', include('project_manager.ol.urls')),
    # 线下工单审核
    path(r'ol/', include('project_manager.of.urls')),
    # 执行任务
    path(r'pt/', include('project_manager.pt.urls')),
    # 系统配置
    path(r'sys/', include('project_manager.sys.urls')),
]

# -*- coding:utf-8 -*-
# edit by fuzongfei

from django.contrib.auth.decorators import login_required
from django.urls import path

from mstats.views import ROnlineMySQLQueryView, OnlineMySQLQueryView, \
    ROfflineMySQLQueryView, OfflineMySQLQueryView, \
    GetSchemaInfo, WebSSHView, GetWebSSHCmdView, GetStruInfoView

urlpatterns = [
    # mysql online query
    path(r'r_online_query/', login_required(ROnlineMySQLQueryView.as_view()), name='p_r_online_query'),
    path(r'online_query/', login_required(OnlineMySQLQueryView.as_view()), name='p_online_query'),
    # mysql offline query
    path(r'r_offline_query/', login_required(ROfflineMySQLQueryView.as_view()), name='p_r_offline_query'),
    path(r'offline_query/', login_required(OfflineMySQLQueryView.as_view()), name='p_offline_query'),
    # 获取库名
    path(r'get_schema_info/', login_required(GetSchemaInfo.as_view()), name='p_get_schema_info'),
    # 获取表结构和索引信息
    path(r'get_stru_info/', login_required(GetStruInfoView.as_view()), name='p_get_stru_info'),
    # webssh
    path(r'web_ssh/', login_required(WebSSHView.as_view()), name='p_web_ssh'),
    path(r'get_ssh_cmd/', login_required(GetWebSSHCmdView.as_view()), name='p_get_ssh_cmd'),
]

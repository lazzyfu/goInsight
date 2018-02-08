# -*- coding:utf-8 -*-
# edit by fuzongfei

from django.urls import path, re_path
from django.contrib.auth.decorators import login_required

from .views import BeautifySQLView, GetInceptionHostConfigView, \
    GetDatabaseListView, IncepOfflineSqlRecords, IncepOfflineSingleSqlDetailView, \
    IncepOfflineAllSqlDetailView, GetRemarkInfo, GetGroupView, GetDbaLeaderView, GetContactsView, \
    IncepOnlineAuditRecordsView, IncepOnlineClickVerifyView, IncepOnlineClickFinishView, IncepOnlineClickCloseView, \
    OnlineAuditDetailView, \
    OnlineSqlReplyView, IncepOfflineSqlCheckView, IncepOnlineSqlCheckView

urlpatterns = [
    path(r'get_inception_hostconfig/', login_required(GetInceptionHostConfigView.as_view()),
         name='p_inception_hostconfig'),
    path(r'get_database/', login_required(GetDatabaseListView.as_view()), name='p_get_database'),
    path(r'beautify_sql/', login_required(BeautifySQLView.as_view()), name='p_beautify_sql'),

    # 线下工单
    path(r'incep_offline_sql_check/', login_required(IncepOfflineSqlCheckView.as_view()),
         name='p_incep_offline_sql_check'),
    path(r'incep_offline_sql_records/', login_required(IncepOfflineSqlRecords.as_view()),
         name='p_incep_offline_sql_records'),
    re_path(r'allsql_detail/(?P<workid>.*)/', login_required(IncepOfflineAllSqlDetailView.as_view())),
    re_path(r'singlesql_detail/(?P<sequence>.*)/', login_required(IncepOfflineSingleSqlDetailView.as_view())),

    # 线上工单
    path(r'incep_online_sql_commit/', login_required(IncepOnlineSqlCheckView.as_view()),
         name='p_incep_online_sql_commit'),
    path(r'get_remark_info/', login_required(GetRemarkInfo.as_view()), name='p_remark_info'),
    path(r'get_group_info/', login_required(GetGroupView.as_view()), name='p_get_group'),
    path(r'get_dba_leader/', login_required(GetDbaLeaderView.as_view()), name='p_get_dba_leader'),
    path(r'get_contacts/', login_required(GetContactsView.as_view()), name='p_get_contacts'),

    # 线上工单记录
    path(r'incep_online_sql_records/', login_required(IncepOnlineAuditRecordsView.as_view()),
         name='p_incep_incep_online_sql_records'),
    path(r'incep_online_click_verify/', login_required(IncepOnlineClickVerifyView.as_view()),
         name='p_incept_incep_online_click_verify'),
    path(r'incep_online_click_finish/', login_required(IncepOnlineClickFinishView.as_view()),
         name='p_incep_incep_online_click_finish'),
    path(r'incep_online_click_close/', login_required(IncepOnlineClickCloseView.as_view()),
         name='p_incep_incep_online_click_close'),

    # 线上工单记录详情
    re_path(r'online_sql_detail/(?P<id>\d+)/(?P<group_id>\d+)/', login_required(OnlineAuditDetailView.as_view()),
            name='p_online_sql_detail'),
    path(r'online_sql_reply/', login_required(OnlineSqlReplyView.as_view()), name='p_online_reply'),
]

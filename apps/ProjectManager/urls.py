# -*- coding:utf-8 -*-
# edit by fuzongfei

from django.urls import path, re_path
from django.contrib.auth.decorators import login_required

from .views import BeautifySQLView, GetInceptionHostConfigView, \
    GetDatabaseListView, InceptionSqlOperateView, InceptionSqlRecords, InceptionSingleSqlDetailView, \
    InceptionAllSqlDetailView, OnlineSqlCommitView, GetRemarkInfo, GetGroupView, GetDbaLeaderView, GetContactsView, \
    OnlineAuditRecordsView, OnlineClickVerifyView, OnlineClickFinishView, OnlineClickCloseView, OnlineAuditDetailView, \
    OnlineSqlReplyView

urlpatterns = [
    path(r'beautify_sql/', login_required(BeautifySQLView.as_view()), name='p_beautify_sql'),
    path(r'inception_sql_operate/', login_required(InceptionSqlOperateView.as_view()), name='p_inception_sql_operate'),
    path(r'get_inception_hostconfig/', login_required(GetInceptionHostConfigView.as_view()),
         name='p_inception_hostconfig'),
    path(r'get_database/', login_required(GetDatabaseListView.as_view()), name='p_get_database'),
    path(r'inception_sql_records/', login_required(InceptionSqlRecords.as_view()), name='p_inception_sql_records'),
    re_path(r'inception_all_sql_detail/(?P<workid>.*)/', login_required(InceptionAllSqlDetailView.as_view())),
    re_path(r'inception_single_sql_detail/(?P<sequence>.*)/', login_required(InceptionSingleSqlDetailView.as_view())),

    # 线上工单提交
    path(r'online_sql_commit/', login_required(OnlineSqlCommitView.as_view()), name='p_online_sql_commit'),
    path(r'get_remark_info/', login_required(GetRemarkInfo.as_view()), name='p_remark_info'),
    path(r'get_group_info/', login_required(GetGroupView.as_view()), name='p_get_group'),
    path(r'get_dba_leader/', login_required(GetDbaLeaderView.as_view()), name='p_get_dba_leader'),
    path(r'get_contacts/', login_required(GetContactsView.as_view()), name='p_get_contacts'),

    # 线上工单记录
    path(r'online_sql_records/', login_required(OnlineAuditRecordsView.as_view()), name='p_online_sql_records'),
    path(r'online_click_verify/', login_required(OnlineClickVerifyView.as_view()), name='p_online_click_verify'),
    path(r'online_click_finish/', login_required(OnlineClickFinishView.as_view()), name='p_online_click_finish'),
    path(r'online_click_close/', login_required(OnlineClickCloseView.as_view()), name='p_online_click_close'),

    # 线上工单记录详情
    re_path(r'online_sql_detail/(?P<id>\d+)/(?P<group_id>\d+)/', login_required(OnlineAuditDetailView.as_view()),
            name='p_online_sql_detail'),
    path(r'online_sql_reply/', login_required(OnlineSqlReplyView.as_view()), name='p_online_reply'),
]

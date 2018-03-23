# -*- coding:utf-8 -*-
# edit by fuzongfei

from django.urls import path, re_path
from django.contrib.auth.decorators import login_required

from .views import BeautifySQLView, GetInceptionHostConfigView, \
    GetDatabaseListView, \
    GetRemarkInfo, GetGroupView, GetDbaLeaderView, GetContactsView, \
    IncepOnlineSqlRecordsView, IncepOnlineClickVerifyView, IncepOnlineClickFinishView, IncepOnlineClickCloseView, \
    OnlineAuditDetailView, \
    OnlineSqlReplyView, IncepSqlCheckView, IncepOnlineSqlCheckView, IncepTasksRecordsListView, \
    IncepTasksDetailView, IncepExecTaskView, IncepTasksResultView, IncepTasksRecordsView, IncepTasksDetailListView, \
    IncepRollbackView, IncepStopView, IncepCreateOnlineTasksView

urlpatterns = [
    path(r'get_inception_hostconfig/', login_required(GetInceptionHostConfigView.as_view()),
         name='p_inception_hostconfig'),
    path(r'get_database/', login_required(GetDatabaseListView.as_view()), name='p_get_database'),

    # 线上工单
    path(r'incep_online_sql_check/', login_required(IncepOnlineSqlCheckView.as_view()),
         name='p_incep_online_sql_check'),
    path(r'incep_online_sql_records/', login_required(IncepOnlineSqlRecordsView.as_view()),
         name='p_incep_online_sql_records'),
    re_path(r'online_sql_detail/(?P<id>\d+)/(?P<group_id>\d+)/', login_required(OnlineAuditDetailView.as_view())),

    # 获取信息
    path(r'get_remark_info/', login_required(GetRemarkInfo.as_view()), name='p_remark_info'),
    path(r'get_group_info/', login_required(GetGroupView.as_view()), name='p_get_group'),
    path(r'get_dba_leader/', login_required(GetDbaLeaderView.as_view()), name='p_get_dba_leader'),
    path(r'get_contacts/', login_required(GetContactsView.as_view()), name='p_get_contacts'),

    # 线上工单审批
    path(r'incep_online_click_verify/', login_required(IncepOnlineClickVerifyView.as_view()),
         name='p_incept_online_click_verify'),
    path(r'incep_online_click_finish/', login_required(IncepOnlineClickFinishView.as_view()),
         name='p_incep_online_click_finish'),
    path(r'incep_online_click_close/', login_required(IncepOnlineClickCloseView.as_view()),
         name='p_incep_online_click_close'),
    path(r'online_sql_reply/', login_required(OnlineSqlReplyView.as_view()), name='p_online_reply'),
    path(r'incep_create_tasks/', login_required(IncepCreateOnlineTasksView.as_view()), name='p_incep_create_tasks'),

    # 线下工单
    path(r'incep_sql_check/', login_required(IncepSqlCheckView.as_view()), name='p_incep_sql_check'),
    path(r'incep_tasks_record/', login_required(IncepTasksRecordsView.as_view()), name='p_incep_tasks_records'),
    path(r'incep_tasks_record_list/', login_required(IncepTasksRecordsListView.as_view()),
         name='p_incep_tasks_records_list'),
    re_path(r'incep_tasks_detail/(?P<taskid>.*)/', login_required(IncepTasksDetailView.as_view())),
    path(r'incep_tasks_detail_list/', login_required(IncepTasksDetailListView.as_view()),
         name='p_incep_tasks_detail_list'),
    re_path(r'incep_tasks_result/', login_required(IncepTasksResultView.as_view()), name='p_incep_tasks_result'),
    path(r'beautify_sql/', login_required(BeautifySQLView.as_view()), name='p_beautify_sql'),

    path(r'incep_exec_task/', login_required(IncepExecTaskView.as_view()), name='p_incep_exec_task'),
    re_path(r'incep_rollback/', login_required(IncepRollbackView.as_view()), name='p_incep_rollback'),
    re_path(r'incep_stop/', login_required(IncepStopView.as_view()), name='p_incep_stop'),
]

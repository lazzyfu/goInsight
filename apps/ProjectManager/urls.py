# -*- coding:utf-8 -*-
# edit by fuzongfei

from django.contrib.auth.decorators import login_required
from django.urls import path, re_path

from .views import BeautifySQLView, IncepHostConfigView, \
    GetDBListView, \
    RemarkInfoView, GroupInfoView, AuditUserView, ContactsInfoView, \
    IncepOlRecordsView, IncepOlApproveView, IncepOlFeedbackView, IncepOlCloseView, \
    IncepOlDetailsView, \
    IncepOlReplyView, IncepOfAuditView, IncepOlAuditView, IncepOfRecordsListView, \
    IncepPerformView, IncepOfResultsView, IncepOfRecordsView, IncepOfDetailsView, \
    IncepRollbackView, IncepStopView, IncepGenerateTasksView, IncepOfDetailsListView

urlpatterns = [
    path(r'incep_host_config/', login_required(IncepHostConfigView.as_view()),
         name='p_incep_host_config'),
    path(r'db_list/', login_required(GetDBListView.as_view()), name='p_db_list'),
    path(r'remark_info/', login_required(RemarkInfoView.as_view()), name='p_remark_info'),
    path(r'group_info/', login_required(GroupInfoView.as_view()), name='p_group_info'),
    path(r'audit_user/', login_required(AuditUserView.as_view()), name='p_audit_user'),
    path(r'contacts_info/', login_required(ContactsInfoView.as_view()), name='p_contacts_info'),
    path(r'beautify_sql/', login_required(BeautifySQLView.as_view()), name='p_beautify_sql'),

    # 线上工单审核
    path(r'incep_ol_audit/', login_required(IncepOlAuditView.as_view()), name='p_incep_ol_audit'),
    path(r'incep_ol_records/', login_required(IncepOlRecordsView.as_view()), name='p_incep_ol_records'),
    re_path(r'incep_ol_details/(?P<id>\d+)/(?P<group_id>\d+)/', login_required(IncepOlDetailsView.as_view())),

    # 线上工单流程操作
    path(r'incep_ol_approve/', login_required(IncepOlApproveView.as_view()), name='p_incep_ol_approve'),
    path(r'incep_ol_feedback/', login_required(IncepOlFeedbackView.as_view()), name='p_incep_ol_feedback'),
    path(r'incep_ol_close/', login_required(IncepOlCloseView.as_view()), name='p_incep_ol_close'),
    path(r'incep_ol_reply/', login_required(IncepOlReplyView.as_view()), name='p_online_reply'),
    path(r'incep_generate_tasks/', login_required(IncepGenerateTasksView.as_view()), name='p_incep_generate_tasks'),

    # 线下工单审核
    path(r'incep_of_audit/', login_required(IncepOfAuditView.as_view()), name='p_incep_of_audit'),

    # 执行任务
    path(r'incep_of_records/', login_required(IncepOfRecordsView.as_view()), name='p_incep_of_records'),
    path(r'incep_of_records_l/', login_required(IncepOfRecordsListView.as_view()), name='p_incep_of_records_l'),
    re_path(r'incep_of_details/(?P<taskid>.*)/', login_required(IncepOfDetailsView.as_view())),
    path(r'incep_of_details_l/', login_required(IncepOfDetailsListView.as_view()), name='p_incep_of_details_l'),
    re_path(r'incep_of_results/', login_required(IncepOfResultsView.as_view()), name='p_incep_of_results'),

    # 执行任务流程操作
    path(r'incep_perform/', login_required(IncepPerformView.as_view()), name='p_incep_perform'),
    re_path(r'incep_rollback/', login_required(IncepRollbackView.as_view()), name='p_incep_rollback'),
    re_path(r'incep_stop/', login_required(IncepStopView.as_view()), name='p_incep_stop'),
]

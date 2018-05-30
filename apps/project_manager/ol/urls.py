# -*- coding:utf-8 -*-
# edit by fuzongfei

from django.contrib.auth.decorators import login_required
from django.urls import path, re_path

from project_manager.ol.views import OnlineAuditView, OnlineRecordsView, OnlineDetailsView, OnlineApproveView, \
    OnlineFeedbackView, OnlineCloseView, OnlineReplyView, OnlineGenerateTasksView, OnlineRecordsListView

urlpatterns = [
    path(r'ol_audit/', login_required(OnlineAuditView.as_view()), name='p_ol_audit'),
    path(r'ol_records/', login_required(OnlineRecordsView.as_view()), name='p_ol_records'),
    path(r'ol_records_list/', login_required(OnlineRecordsListView.as_view()), name='p_ol_records_list'),
    re_path(r'ol_details/(?P<id>\d+)/(?P<group_id>\d+)/', login_required(OnlineDetailsView.as_view())),

    path(r'ol_approve/', login_required(OnlineApproveView.as_view()), name='p_ol_approve'),
    path(r'ol_feedback/', login_required(OnlineFeedbackView.as_view()), name='p_ol_feedback'),
    path(r'ol_close/', login_required(OnlineCloseView.as_view()), name='p_ol_close'),
    path(r'ol_reply/', login_required(OnlineReplyView.as_view()), name='p_ol_reply'),
    path(r'ol_generate_tasks/', login_required(OnlineGenerateTasksView.as_view()), name='p_ol_generate_tasks'),
]
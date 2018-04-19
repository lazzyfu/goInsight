# -*- coding:utf-8 -*-
# edit by fuzongfei

from django.contrib.auth.decorators import login_required
from django.urls import path, re_path

from project_manager.ol.views import IncepOlAuditView, IncepOlRecordsView, IncepOlDetailsView, IncepOlApproveView, \
    IncepOlFeedbackView, IncepOlCloseView, IncepOlReplyView, IncepGenerateTasksView, IncepOlRecordsListView

urlpatterns = [
    path(r'incep_ol_audit/', login_required(IncepOlAuditView.as_view()), name='p_incep_ol_audit'),
    path(r'incep_ol_records/', login_required(IncepOlRecordsView.as_view()), name='p_incep_ol_records'),
    path(r'incep_ol_records_l/', login_required(IncepOlRecordsListView.as_view()), name='p_incep_ol_records_l'),
    re_path(r'incep_ol_details/(?P<id>\d+)/(?P<group_id>\d+)/', login_required(IncepOlDetailsView.as_view())),

    path(r'incep_ol_approve/', login_required(IncepOlApproveView.as_view()), name='p_incep_ol_approve'),
    path(r'incep_ol_feedback/', login_required(IncepOlFeedbackView.as_view()), name='p_incep_ol_feedback'),
    path(r'incep_ol_close/', login_required(IncepOlCloseView.as_view()), name='p_incep_ol_close'),
    path(r'incep_ol_reply/', login_required(IncepOlReplyView.as_view()), name='p_online_reply'),
    path(r'incep_generate_tasks/', login_required(IncepGenerateTasksView.as_view()), name='p_incep_generate_tasks'),
]
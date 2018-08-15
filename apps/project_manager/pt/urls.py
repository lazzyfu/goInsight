# -*- coding:utf-8 -*-
# edit by fuzongfei

from django.contrib.auth.decorators import login_required
from django.urls import path, re_path

from project_manager.pt.views import PerformExecView, PerformRollbackView, PerformStopView, PerformRecordsView, \
    PerformRecordsListView, PerformDetailsView, PerformDetailsListView, PerformResultsView, PerformFullExecView, \
    PerformRecordsSQLPreView

urlpatterns = [
    # 执行任务记录
    path(r'perform_records/', login_required(PerformRecordsView.as_view()), name='p_perform_records'),
    path(r'perform_records_list/', login_required(PerformRecordsListView.as_view()),
         name='p_perform_records_list'),
    path(r'perform_records_sqlpre/', login_required(PerformRecordsSQLPreView.as_view()),
         name='p_perform_records_sqlpre'),
    # 执行任务详情
    re_path(r'perform_details/(?P<taskid>.*)/', login_required(PerformDetailsView.as_view())),
    path(r'perform_details_list/', login_required(PerformDetailsListView.as_view()),
         name='p_perform_details_list'),
    # 执行任务结果
    path(r'perform_results/', login_required(PerformResultsView.as_view()), name='p_perform_results'),
    # 执行任务操作
    path(r'perform_full_exec/', login_required(PerformFullExecView.as_view()), name='p_perform_full_exec'),
    path(r'perform_exec/', login_required(PerformExecView.as_view()), name='p_perform_exec'),
    re_path(r'perform_rollback/', login_required(PerformRollbackView.as_view()), name='p_perform_rollback'),
    re_path(r'perform_stop/', login_required(PerformStopView.as_view()), name='p_perform_stop'),
]

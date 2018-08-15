# -*- coding:utf-8 -*-
# edit by fuzongfei

from django.contrib.auth.decorators import login_required
from django.urls import path, re_path

from project_manager.ol.views import WorkOrderAuditView, OnlineRecordsView, WorkOrderDetailsView, WorkOrderApproveView, \
    WorkOrderFeedbackView, WorkOrderCloseView, WorkOrderReplyView, GeneratePerformTasksView, AuditRecordsListView, \
    DeployTasksView, ROnlineAuditTasksView, OnlineAuditTasksListView, StagingRecordsView, MOneRecordsView, \
    OnlineWorkOrderAuditView, DingNoticeView

urlpatterns = [
    # 生产环境工单提交页面
    path(r'ol_work_order_audit/', login_required(OnlineWorkOrderAuditView.as_view()), name='p_ol_work_order_audit'),
    # 生产环境工单列表页
    path(r'ol_records/', login_required(OnlineRecordsView.as_view()), name='p_ol_records'),
    # 预发布环境工单列表页
    path(r'staging_records/', login_required(StagingRecordsView.as_view()), name='p_staging_records'),
    # M1环境工单列表页
    path(r'test_records/', login_required(MOneRecordsView.as_view()), name='p_test_records'),
    # 生成工单
    path(r'work_order_audit/', login_required(WorkOrderAuditView.as_view()), name='p_work_order_audit'),
    # 工单页面列表
    path(r'audit_records_list/', login_required(AuditRecordsListView.as_view()), name='p_audit_records_list'),
    # 工单详情
    re_path(r'work_order_details/(?P<id>\d+)/', login_required(WorkOrderDetailsView.as_view())),
    # 工单的审批、反馈、关闭、回复
    path(r'work_order_approve/', login_required(WorkOrderApproveView.as_view())),
    path(r'work_order_feedback/', login_required(WorkOrderFeedbackView.as_view())),
    path(r'work_order_close/', login_required(WorkOrderCloseView.as_view())),
    path(r'work_order_reply/', login_required(WorkOrderReplyView.as_view()), name='p_work_order_reply'),
    # 获取部署任务步骤信息
    path(r'deploy_tasks/', login_required(DeployTasksView.as_view()), name='p_deploy_tasks'),
    # 通知未完成工单的开发
    path(r'ding_notice/', login_required(DingNoticeView.as_view()), name='p_ding_notice'),
    # 工单转换成执行任务
    path(r'generate_perform_tasks/', login_required(GeneratePerformTasksView.as_view())),
    path(r'ol_tasks/', login_required(ROnlineAuditTasksView.as_view()), name='p_tasks'),
    path(r'ol_get_tasks/', login_required(OnlineAuditTasksListView.as_view()), name='p_get_tasks'),
]

# -*- coding:utf-8 -*-
# edit by fuzongfei

from django.contrib.auth.decorators import login_required
from django.urls import path, re_path

from orders.views.View import SQLOrdersView, OrdersCommitView, RenderOrdersEnviView, OrdersListView, OrdersDetailsView, \
    OrderReplyView, GetOrderReplyView, HookOrdersView, SQLExportView, OpsOrdersView, OnlineVersionView, MyOrdersView, \
    MyOrdersListView, OrderExecuteDetailsView
from orders.views.opView import OrderApproveView, OrderFeedbackView, OrderReviewView, OrderCloseView
from orders.views.tasksView import GenerateSubtasksView, RenderSubtasksView, SubTasksDetailView, FullExecuteView, \
    SingleExecuteView, GetTasksLogView, StopExecuteView
from orders.views.utilsView import BeautifySQLView, GetSchemasView, GetSysEnviView, OnlineVersionNoExpireView, \
    SyntaxCheckView, OnlineVersionDetailView, OnlineVersionListView

urlpatterns = [
    # 获取系统配置的环境
    path('get_sysenvi/', login_required(GetSysEnviView.as_view())),
    # 格式化SQL
    path('beautify_sql/', login_required(BeautifySQLView.as_view())),
    # 语法检查
    path('syntax_check/', login_required(SyntaxCheckView.as_view())),
    # 获取指定环境的schema
    path('get_schemas/', login_required(GetSchemasView.as_view())),
    # 渲染上线版本页面
    path('online_version/', login_required(OnlineVersionView.as_view()), name='p_online_version'),
    path('online_version/list/', login_required(OnlineVersionListView.as_view())),
    # 获取上线版本
    path('online_version/no_expire/', login_required(OnlineVersionNoExpireView.as_view())),
    path('online_version/detail/', login_required(OnlineVersionDetailView.as_view())),
    # 渲染DML和DDL工单页面
    path('sql_orders/', login_required(SQLOrdersView.as_view()), name='p_dml'),
    # 渲染SQL导出工单页面
    path('sql_export/', login_required(SQLExportView.as_view()), name='p_sql_export'),
    # 渲染运维工单页面
    path('ops_orders/', login_required(OpsOrdersView.as_view()), name='p_ops'),
    # 提交工单
    path('commit/', login_required(OrdersCommitView.as_view()), name='p_commit'),
    # 工单列表
    re_path('envi/(?P<envi_id>\d+)/', login_required(RenderOrdersEnviView.as_view())),
    path('list/', login_required(OrdersListView.as_view())),
    # 工单操作，如：审核、执行、反馈、复核等
    path('op/approve/', login_required(OrderApproveView.as_view())),
    path('op/feedback/', login_required(OrderFeedbackView.as_view())),
    path('op/review/', login_required(OrderReviewView.as_view())),
    path('op/close/', login_required(OrderCloseView.as_view())),
    # 工单详情
    re_path(r'detail/(?P<id>\d+)/', login_required(OrdersDetailsView.as_view())),
    path(r'execute_detail/', login_required(OrderExecuteDetailsView.as_view()),
         name='p_get_sql_exec_details'),
    # 工单回复
    path(r'reply/', login_required(OrderReplyView.as_view())),
    path(r'get_order_reply/', login_required(GetOrderReplyView.as_view())),
    # 工单HOOK功能
    path(r'hook/', login_required(HookOrdersView.as_view())),
    # SQL工单拆解成执行任务
    path(r'generate_subtasks/', login_required(GenerateSubtasksView.as_view())),
    re_path(r'subtasks/list/(?P<taskid>.*)/', login_required(RenderSubtasksView.as_view())),
    path(r'subtasks/detail/', login_required(SubTasksDetailView.as_view())),
    # 任务的执行
    path(r'subtasks/full/', login_required(FullExecuteView.as_view())),
    path(r'subtasks/single/', login_required(SingleExecuteView.as_view())),
    path(r'subtasks/stop/', login_required(StopExecuteView.as_view())),
    # 获取子任务的执行日志
    path(r'subtasks/getlog/', login_required(GetTasksLogView.as_view())),
    # 渲染我的工单页面
    path('my_orders/', login_required(MyOrdersView.as_view()), name='p_my_orders'),
    path(r'my_orders/list', login_required(MyOrdersListView.as_view())),
]

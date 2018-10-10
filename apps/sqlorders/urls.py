# -*- coding:utf-8 -*-
# edit by fuzongfei
from django.contrib.auth.decorators import login_required
from django.urls import path, re_path

from sqlorders.views import GetSqlOrdersEnviView, RenderSqlDmlOrdersView, GetAuditUserView, GetTablesView, \
    RenderSqlDdlOrdersView, GetOfflineSchemasView, SqlOrdersAuditView, RenderSqlOrdersListView, \
    SqlOrdersListView, SyntaxCheckView, BeautifySQLView, SqlOrdersDetailsView, SqlOrdersApproveView, \
    SqlOrdersFeedbackView, SqlOrdersCloseView, GetParentSchemasView, HookSqlOrdersView, GeneratePerformTasksView, \
    RenderPerformTasksView, PerformTasksDetailsView, PerformTasksSQLPreView, SinglePerformTasksView, \
    FullPerformTasksView, GetPerformTasksResultView, PerformTasksRollbackView, \
    SqlOrdersTasksVersionView, RenderSqlOrdersTasksVersionView, SqlOrdersTasksVersionListView, GetVersionOrdersList, \
    PerformTasksOpView, GetTargetSchemasView, CommitOrderReplyView, GetOrderReplyView

urlpatterns = [
    path('get_sql_orders_envi/', login_required(GetSqlOrdersEnviView.as_view()), name='p_get_sql_orders_envi'),
    path(r'render_sql_dml_orders/', login_required(RenderSqlDmlOrdersView.as_view()), name='p_render_sql_dml_orders'),
    path(r'render_sql_ddl_orders/', login_required(RenderSqlDdlOrdersView.as_view()), name='p_render_sql_ddl_orders'),
    path(r'get_audit_user/', login_required(GetAuditUserView.as_view()), name='p_get_audit_user'),
    # 获取schema列表
    path(r'get_product_schemas/', login_required(GetTargetSchemasView.as_view()), name='p_get_target_schemas'),
    path(r'get_offline_schemas/', login_required(GetOfflineSchemasView.as_view()), name='p_get_offline_schemas'),
    path(r'get_parent_schemas/', login_required(GetParentSchemasView.as_view()), name='p_get_parent_schemas'),
    path(r'get_tables/', login_required(GetTablesView.as_view()), name='p_get_tables'),
    # inceotion语法检查
    path(r'syntax_check/', login_required(SyntaxCheckView.as_view())),
    # sql美化
    path(r'beautify_sql/', login_required(BeautifySQLView.as_view())),
    # 处理提交的DDL、DML工单
    path(r'sql_order_audit/', login_required(SqlOrdersAuditView.as_view()), name='p_sql_orders_audit'),
    # 查看工单
    re_path(r'sql_orders_list/(?P<envi_id>\d+)/', login_required(RenderSqlOrdersListView.as_view())),
    path(r'get_sql_orders_list/', login_required(SqlOrdersListView.as_view()), name='p_get_sql_orders_list'),
    # 工单详情
    re_path(r'sql_orders_details/(?P<id>\d+)/', login_required(SqlOrdersDetailsView.as_view())),
    # 工单的审批、反馈、关闭
    path(r'sql_orders_approve/', login_required(SqlOrdersApproveView.as_view())),
    path(r'sql_orders_feedback/', login_required(SqlOrdersFeedbackView.as_view())),
    path(r'sql_orders_close/', login_required(SqlOrdersCloseView.as_view())),
    # 回复工单
    path(r'commit_order_reply/', login_required(CommitOrderReplyView.as_view()), name='p_commit_order_reply'),
    path(r'get_order_reply/', login_required(GetOrderReplyView.as_view()), name='p_get_order_reply'),
    # 钩子
    path(r'hook_sql_orders/', login_required(HookSqlOrdersView.as_view()), name='p_hook_sql_orders'),
    # 工单转换成执行任务
    path(r'generate_perform_tasks/', login_required(GeneratePerformTasksView.as_view())),
    # 执行任务
    re_path(r'perform_tasks/(?P<taskid>.*)/', login_required(RenderPerformTasksView.as_view())),
    path(r'perform_tasks_details/', login_required(PerformTasksDetailsView.as_view()),
         name='p_perform_tasks_details'),
    # 预览SQL
    path(r'perform_tasks_sqlpre/', login_required(PerformTasksSQLPreView.as_view()),
         name='p_perform_tasks_sqlpre'),
    # 执行任务，执行、回滚、停止
    path(r'full_perform_tasks/', login_required(FullPerformTasksView.as_view()), name='p_full_perform_tasks'),
    path(r'single_perform_tasks/', login_required(SinglePerformTasksView.as_view()), name='p_single_perform_tasks'),
    re_path(r'perform_rollback/', login_required(PerformTasksRollbackView.as_view()), name='p_perform_tasks_rollback'),
    re_path(r'perform_tasks_stop/', login_required(PerformTasksOpView.as_view()), name='p_perform_tasks_stop'),
    # 执行任务结果
    path(r'get_perform_tasks_result/', login_required(GetPerformTasksResultView.as_view()),
         name='p_get_perform_tasks_result'),
    # 上线任务版本
    path(r'render_sql_tasks_version/', login_required(RenderSqlOrdersTasksVersionView.as_view()),
         name='p_render_get_sql_tasks_version'),
    path(r'get_sql_tasks_version/', login_required(SqlOrdersTasksVersionView.as_view()),
         name='p_get_sql_tasks_version'),
    path(r'get_sql_tasks_version_list/', login_required(SqlOrdersTasksVersionListView.as_view()),
         name='p_get_sql_tasks_version_list'),
    # 获取任务版本内的工单信息
    path(r'get_version_orders_list/', login_required(GetVersionOrdersList.as_view())),
]

# -*- coding:utf-8 -*-
# edit by fuzongfei
from django.contrib.auth.decorators import login_required
from django.urls import path, re_path, include

from query.views import GetTableStrucView, \
    RenderSqlQueryView, GetGrantSchemaView, ExecSqlQueryView, GetHistorySqlView, GetTableIndexView, GetTableBaseView, \
    RenderDictView, GenerateHtmlView, GetTablesView, QueryStatusCheckView

urlpatterns = [
    # 渲染查询页面
    re_path(r'sql_query/', login_required(RenderSqlQueryView.as_view()), name='p_sqlquery'),
    # 获取指定环境授权给用户的schema信息
    path(r'get_schemas_grant/', login_required(GetGrantSchemaView.as_view()),
         name='p_get_grant_schema'),
    path(r'get_tables/', login_required(GetTablesView.as_view()), name='p_get_tables'),
    # 获取表结构和索引信息
    path(r'get_stru_info/', login_required(GetTableStrucView.as_view()), name='p_get_table_stru'),
    path(r'get_index_info/', login_required(GetTableIndexView.as_view()), name='p_get_table_index'),
    path(r'get_index_base/', login_required(GetTableBaseView.as_view()), name='p_get_table_base'),
    # 执行查询
    path(r'exec_query/', login_required(ExecSqlQueryView.as_view()), name='p_exec_sql_query'),
    # 获取当前用户执行的sql历史
    path(r'history_sql/', login_required(GetHistorySqlView.as_view()), name='p_get_history_sql'),
    # 查询状态检查
    path(r'query_status_check/', login_required(QueryStatusCheckView.as_view()), name='p_query_status_check'),
    # 数据字典
    path(r'render/', login_required(RenderDictView.as_view()), name='p_dbdict'),
    path(r'general_html/', login_required(GenerateHtmlView.as_view()), name='p_generate_html'),
    # xiao soar
    path('soar/', include('query.soar.urls')),
]

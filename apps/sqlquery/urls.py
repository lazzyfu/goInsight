# -*- coding:utf-8 -*-
# edit by fuzongfei
from django.contrib.auth.decorators import login_required
from django.urls import path, re_path

from sqlquery.views import GetStruInfoView, \
    RenderSqlQueryView, GetSchemasGrantView, ExecSqlQueryView, GetHistorySqlView

urlpatterns = [
    # mysql query api
    # 渲染查询页面
    re_path(r'sql_query/(?P<envi_id>\d+)/', login_required(RenderSqlQueryView.as_view())),
    # 获取指定环境授权给用户的schema信息
    path(r'get_schemas_grant/', login_required(GetSchemasGrantView.as_view()),
         name='p_get_schemas_grant'),
    # 获取表结构和索引信息
    path(r'get_stru_info/', login_required(GetStruInfoView.as_view()), name='p_get_stru_info'),
    # 执行查询
    path(r'exec_query/', login_required(ExecSqlQueryView.as_view()), name='p_exec_sql_query'),
    # 获取当前用户执行的sql历史
    path(r'history_sql/', login_required(GetHistorySqlView.as_view()), name='p_get_history_sql'),
]

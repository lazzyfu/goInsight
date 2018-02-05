# -*- coding:utf-8 -*-
# edit by fuzongfei

from django.urls import path, re_path
from django.contrib.auth.decorators import login_required

from .views import ProjectListView, BeautifySQLView, GetInceptionHostConfigView, \
    GetDatabaseListView, InceptionSqlOperateView, InceptionSqlRecords, InceptionSingleSqlDetailView, \
    InceptionAllSqlDetailView, OnlineSqlCommitView, GetRemarkInfo

urlpatterns = [
    path(r'index/', ProjectListView.as_view(), name='p_project'),
    path(r'beautify_sql/', login_required(BeautifySQLView.as_view()), name='p_beautify_sql'),
    path(r'inception_sql_operate/', login_required(InceptionSqlOperateView.as_view()), name='p_inception_sql_operate'),
    path(r'get_inception_hostconfig/', login_required(GetInceptionHostConfigView.as_view()),
         name='p_inception_hostconfig'),
    path(r'get_database/', login_required(GetDatabaseListView.as_view()), name='p_get_database'),
    path(r'inception_sql_records/', login_required(InceptionSqlRecords.as_view()), name='p_inception_sql_records'),
    re_path(r'inception_all_sql_detail/(?P<workid>.*)/', login_required(InceptionAllSqlDetailView.as_view())),
    re_path(r'inception_single_sql_detail/(?P<sequence>.*)/', login_required(InceptionSingleSqlDetailView.as_view())),

    # 线上工单
    path(r'online_sql_commit/', login_required(OnlineSqlCommitView.as_view()), name='p_online_sql_commit'),
    path(r'get_remark_info/', login_required(GetRemarkInfo.as_view()), name='p_remark_info'),
]

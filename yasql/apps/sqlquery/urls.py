# -*- coding:utf-8 -*-
# edit by xff

from django.urls import path, re_path, include

from sqlquery import views

sqlquery_patterns = [
    path('tree', views.GetTreeView.as_view(), name='v1.sqlquery.get.tree'),
    path('execute-query', views.ExecuteQueryView.as_view(), name='v1.sqlquery.execute.query'),
    path('delete-query-hash', views.DeleteQueryHashView.as_view(), name='v1.sqlquery.delete.query.hash'),
    path('get/tableinfo', views.GetTableInfoView.as_view(), name='v1.sqlquery.get.tableinfo'),
    path('get/history/sql', views.GetHistorySQLView.as_view(), name='v1.sqlquery.get.history.sql'),
    path('get/dbdict', views.GetDBDictLView.as_view(), name='v1.sqlquery.get.dbdict'),
]

urlpatterns = [
    re_path(r'sqlquery/', include(sqlquery_patterns)),
]

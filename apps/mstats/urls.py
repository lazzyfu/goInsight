# -*- coding:utf-8 -*-
# edit by fuzongfei

from django.contrib.auth.decorators import login_required
from django.urls import path

from mstats.views import RenderMySQLUserView, MySQLUserView, MySQLuser_managerView

urlpatterns = [
    path(r'r_mysql_user_manager/', login_required(RenderMySQLUserView.as_view()), name='p_r_mysql_user_manager'),
    path(r'mysql_user/', login_required(MySQLUserView.as_view()), name='p_mysql_user'),
    path(r'mysql_user_manager/', login_required(MySQLuser_managerView.as_view()), name='p_mysql_user_manager'),
]

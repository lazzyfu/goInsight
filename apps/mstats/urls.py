# -*- coding:utf-8 -*-
# edit by fuzongfei

from django.contrib.auth.decorators import login_required
from django.urls import path

from mstats.views import RenderMySQLUserManagerView, MySQLUserManagerView, MySQLPrivModifyView

urlpatterns = [
    path(r'r_mysql_user_manager/', login_required(RenderMySQLUserManagerView.as_view()), name='p_r_mysql_user_manager'),
    path(r'mysql_user_manager/', login_required(MySQLUserManagerView.as_view()), name='p_mysql_user_manager'),
    path(r'mysql_priv_modify/', login_required(MySQLPrivModifyView.as_view()), name='p_mysql_priv_modify'),
]

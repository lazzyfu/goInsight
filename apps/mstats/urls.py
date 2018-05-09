# -*- coding:utf-8 -*-
# edit by fuzongfei

from django.contrib.auth.decorators import login_required
from django.urls import path, re_path

from mstats.views import RenderMySQLUserView, MySQLUserView, MysqlUserManager, RBackupTaskView, BackupTaskView, \
    BackupTaskDetailView, BackupTaskPreviewView, BackupTaskPreviewListView

urlpatterns = [
    path(r'r_mysql_user_manager/', login_required(RenderMySQLUserView.as_view()), name='p_r_mysql_user_manager'),
    path(r'mysql_user/', login_required(MySQLUserView.as_view()), name='p_mysql_user'),
    path(r'mysql_user_manager/', login_required(MysqlUserManager.as_view()), name='p_mysql_user_manager'),
    path(r'rbackup_task/', login_required(RBackupTaskView.as_view()), name='p_rbackup_task'),
    path(r'backup_task/', login_required(BackupTaskView.as_view()), name='p_backup_task'),
    path(r'backup_task_detail/', login_required(BackupTaskDetailView.as_view()), name='p_backup_task_detail'),
    re_path(r'backup_task_preview/(?P<id>\d+)/', login_required(BackupTaskPreviewView.as_view())),
    re_path(r'backup_task_preview_list', login_required(BackupTaskPreviewListView.as_view()),
            name='p_backup_task_preview_list'),
]

# -*- coding:utf-8 -*-
# edit by fuzongfei

from django.contrib.auth.decorators import login_required
from django.urls import path

from .views import CrontabView, RCrontabView, \
    GetCeleryTasksView, GetCrontabView, DeletePeriodicTaskView, ModifyPeriodicTaskView

urlpatterns = [
    path(r'rcrontab/', login_required(RCrontabView.as_view()), name='p_rcrontab'),
    path(r'crontab/', login_required(CrontabView.as_view()), name='p_crontab'),
    path(r'delete_periodic_task/', login_required(DeletePeriodicTaskView.as_view()), name='p_delete_periodic_task'),
    path(r'modify_periodic_task/', login_required(ModifyPeriodicTaskView.as_view()), name='p_modify_periodic_task'),
    path(r'get_crontab/', login_required(GetCrontabView.as_view()), name='p_get_crontab'),
    path(r'get_celey_tasks/', login_required(GetCeleryTasksView.as_view()), name='p_get_celery_tasks'),
]

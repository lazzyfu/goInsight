# -*- coding:utf-8 -*-
# edit by fuzongfei

from django.contrib.auth.decorators import login_required
from django.urls import path

from scheduled_tasks.views import CrontabView, RCrontabView, RPeriodicTaskView, PeriodicTaskView, \
    GetCeleryTasksView, GetCrontabView

urlpatterns = [
    path(r'rcrontab/', login_required(RCrontabView.as_view()), name='p_rcrontab'),
    path(r'crontab/', login_required(CrontabView.as_view()), name='p_crontab'),
    path(r'rperiodic_task/', login_required(RPeriodicTaskView.as_view()), name='p_rperiodic_task'),
    path(r'periodic_task/', login_required(PeriodicTaskView.as_view()), name='p_periodic_task'),
    path(r'get_crontab/', login_required(GetCrontabView.as_view()), name='p_get_crontab'),
    path(r'get_celey_tasks/', login_required(GetCeleryTasksView.as_view()), name='p_get_celery_tasks'),
]

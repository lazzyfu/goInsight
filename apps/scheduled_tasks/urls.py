# -*- coding:utf-8 -*-
# edit by fuzongfei

from django.contrib.auth.decorators import login_required
from django.urls import path

from scheduled_tasks.views import CrontabView, RCrontabView

urlpatterns = [
    path(r'rcrontab/', login_required(RCrontabView.as_view()), name='p_rcrontab'),
    path(r'crontab/', login_required(CrontabView.as_view()), name='p_crontab'),
]

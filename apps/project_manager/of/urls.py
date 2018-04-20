# -*- coding:utf-8 -*-
# edit by fuzongfei

from django.contrib.auth.decorators import login_required
from django.urls import path

from project_manager.of.views import IncepOfAuditView

urlpatterns = [
    path(r'incep_of_audit/', login_required(IncepOfAuditView.as_view()), name='p_incep_of_audit'),
]

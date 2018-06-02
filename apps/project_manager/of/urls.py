# -*- coding:utf-8 -*-
# edit by fuzongfei

from django.contrib.auth.decorators import login_required
from django.urls import path

from project_manager.of.views import OflineAuditView

urlpatterns = [
    path(r'of_audit/', login_required(OflineAuditView.as_view()), name='p_of_audit'),
]

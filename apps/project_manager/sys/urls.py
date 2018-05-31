# -*- coding:utf-8 -*-
# edit by fuzongfei

from django.contrib.auth.decorators import login_required
from django.urls import path

from project_manager.sys.views import SysConfigView, GetDomainView, GetDBAccountView, ModifyDBAccountView, \
    GetWebhookView

urlpatterns = [
    path(r'config/', login_required(SysConfigView.as_view()), name='p_config'),
    path(r'get_domain/', login_required(GetDomainView.as_view()), name='p_get_domain'),
    path(r'get_webhook/', login_required(GetWebhookView.as_view()), name='p_get_webhook'),
    path(r'get_db_account/', login_required(GetDBAccountView.as_view()), name='p_get_db_account'),
    path(r'modify_db_account/', login_required(ModifyDBAccountView.as_view()), name='p_modify_db_account'),
]

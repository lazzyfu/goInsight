# -*- coding:utf-8 -*-
# edit by fuzongfei

from django.contrib.auth.decorators import login_required
from django.urls import path

from webshell.views import WebSSHView, GetWebSSHCmdView

urlpatterns = [
    path(r'web_ssh/', login_required(WebSSHView.as_view()), name='p_web_ssh'),
    path(r'get_ssh_cmd/', login_required(GetWebSSHCmdView.as_view()), name='p_get_ssh_cmd'),
]

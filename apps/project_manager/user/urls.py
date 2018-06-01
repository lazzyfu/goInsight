# -*- coding:utf-8 -*-
# edit by fuzongfei

from django.contrib.auth.decorators import login_required
from django.urls import path, re_path

from project_manager.user.views import UserConfigView, GetUserPermissionView, GetProjectView, GetRoleView

urlpatterns = [
    path(r'config/', login_required(UserConfigView.as_view()), name='p_user_config'),
    path(r'get_user_permission/', login_required(GetUserPermissionView.as_view()), name='p_get_user_permission'),
    path(r'get_project/', login_required(GetProjectView.as_view()), name='p_get_project'),
    path(r'get_role/', login_required(GetRoleView.as_view()), name='p_get_role'),
]

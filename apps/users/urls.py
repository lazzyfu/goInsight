# -*- coding:utf-8 -*-
# edit by fuzongfei
from django.contrib.auth.decorators import login_required
from django.urls import path

from users.views import LoginView, UserProfileView, LogoutView, ChangePasswordView, ChangeMobileView, ChangePicView, \
    GetUserMailView, VerifyView

urlpatterns = [
    path('login/', LoginView.as_view(), name='p_login'),
    path('verify', VerifyView.as_view(), name='p_verify'),
    path(r'logout/', login_required(LogoutView.as_view()), name='p_logout'),
    path(r'profile/', login_required(UserProfileView.as_view()), name='p_user_profile'),
    path(r'change_password/', login_required(ChangePasswordView.as_view()), name='p_change_password'),
    path(r'change_mobile/', login_required(ChangeMobileView.as_view()), name='p_change_mobile'),
    path(r'change_picture/', login_required(ChangePicView.as_view()), name='p_change_picture'),
    # 获取用户名和邮箱
    path(r'get_usermail/', login_required(GetUserMailView.as_view()), name='p_get_user_mail'),
]

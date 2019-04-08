# -*- coding:utf-8 -*-
# edit by fuzongfei
from django.contrib.auth.decorators import login_required
from django.urls import path

from users.views import LoginView, UserProfileView, LogoutView, VerifyCodeView, ChangePasswordView, ChangeMobileView, \
    ChangePicView, GetEmailCcView, GetAuditorView

urlpatterns = [
    path('login/', LoginView.as_view(), name='p_login'),
    path('verify', VerifyCodeView.as_view(), name='p_verifycode'),
    path('logout/', login_required(LogoutView.as_view()), name='p_logout'),
    path('profile/', login_required(UserProfileView.as_view()), name='p_profile'),
    path('change_password/', login_required(ChangePasswordView.as_view()), name='p_change_password'),
    path('change_mobile/', login_required(ChangeMobileView.as_view()), name='p_change_mobile'),
    path('change_picture/', login_required(ChangePicView.as_view()), name='p_change_picture'),
    path('get_email_cc/', login_required(GetEmailCcView.as_view())),
    path('get_auditor/', login_required(GetAuditorView.as_view())),
]

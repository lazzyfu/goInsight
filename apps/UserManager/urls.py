# -*- coding:utf-8 -*-
# edit by fuzongfei

from django.urls import path
from django.contrib.auth.decorators import login_required

from .views import LoginView, LogoutView

urlpatterns = [
    path(r'login/', LoginView.as_view(), name='p_login'),
    path(r'logout/', login_required(LogoutView.as_view()), name='p_logout'),
]
# -*- coding:utf-8 -*-
# edit by xff
from django.urls import path

from users import views

urlpatterns = [
    path('login', views.Login.as_view(), name='v1.users.login'),
    path('logout', views.Logout.as_view(), name='v1.users.logout'),
    path('info', views.UserInfo.as_view(), name='v1.users.info'),
    path('list', views.GetUsers.as_view(), name='v1.users.list'),
    path('update/<str:username>', views.UpdateUserInfo.as_view(), name='v1.users.update'),
    path('change/password', views.ChangePassword.as_view(), name='v1.users.change.password'),
    path('change/avatar', views.ChangeAvatar.as_view(), name='v1.users.change.avatar'),
    path('dashboard/sys', views.SysDashboard.as_view(), name='v1.users.dashboard.sys'),
    path('dashboard/self', views.SelfDashboard.as_view(), name='v1.users.dashboard.self'),
]

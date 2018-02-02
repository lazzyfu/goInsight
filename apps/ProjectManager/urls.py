# -*- coding:utf-8 -*-
# edit by fuzongfei

from django.urls import path
from django.contrib.auth.decorators import login_required

from .views import ProjectListView

urlpatterns = [
    path(r'index/', ProjectListView.as_view(), name='p_project'),
]
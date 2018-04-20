# -*- coding:utf-8 -*-
# edit by fuzongfei

from django.contrib.auth.decorators import login_required
from django.urls import path

from project_manager.dexport.views import OlDataExportView, ExecOlDataExportView

urlpatterns = [
    path(r'data_export/', login_required(OlDataExportView.as_view()), name='p_data_export'),
    path(r'exec_data_export/', login_required(ExecOlDataExportView.as_view()), name='p_exec_data_export'),
]

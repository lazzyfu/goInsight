# -*- coding:utf-8 -*-
# edit by fuzongfei

from django.contrib.auth.decorators import login_required
from django.urls import path

from project_manager.dexport.views import OlDataExportDetailRecordsView, \
    OlDataExportDetailRecordsListView, \
    ExecOlDataExportDetailView, OlDataExportDetailDownloadView, OlDataExportView

urlpatterns = [
    path(r'data_export/', login_required(OlDataExportView.as_view()), name='p_data_export'),
    path(r'data_export_records/', login_required(OlDataExportDetailRecordsView.as_view()), name='p_data_export_records'),
    path(r'data_export_records_l/', login_required(OlDataExportDetailRecordsListView.as_view()),
         name='p_data_export_records_l'),
    path(r'exec_data_export/', login_required(ExecOlDataExportDetailView.as_view()), name='p_exec_data_export'),
    path(r'data_export_download/', login_required(OlDataExportDetailDownloadView.as_view()), name='p_data_export_download'),
]

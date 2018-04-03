# -*- coding:utf-8 -*-
# edit by fuzongfei

from django.contrib.auth.decorators import login_required
from django.urls import path

from ProjectManager.dexport.views import DataExportView, DataExportRecordsView, DataExportRecordsListView, \
    ExecDataExportView, DataExportDownloadView

urlpatterns = [
    path(r'data_export/', login_required(DataExportView.as_view()), name='p_data_export'),
    path(r'data_export_records/', login_required(DataExportRecordsView.as_view()), name='p_data_export_records'),
    path(r'data_export_records_l/', login_required(DataExportRecordsListView.as_view()),
         name='p_data_export_records_l'),
    path(r'exec_data_export/', login_required(ExecDataExportView.as_view()), name='p_exec_data_export'),
    path(r'data_export_download/', login_required(DataExportDownloadView.as_view()), name='p_data_export_download'),
]

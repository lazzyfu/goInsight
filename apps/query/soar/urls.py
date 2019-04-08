# -*- coding:utf-8 -*-
# edit by fuzongfei
from django.contrib.auth.decorators import login_required
from django.urls import path

from query.soar.views import RenderSoarView, SoarAnalyzeView

urlpatterns = [
    # xiao soar
    path(r'render/', login_required(RenderSoarView.as_view()), name='p_soar'),
    path(r'analyze/', login_required(SoarAnalyzeView.as_view()), name='p_soar_analyze')
]
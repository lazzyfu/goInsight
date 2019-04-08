# -*- coding:utf-8 -*-
# edit by fuzongfei
from django.contrib.auth.decorators import login_required
from django.urls import path

from dash.views import RenderOrderChartView, GetOrderChartView, GetSqlCountView

urlpatterns = [
    # 获取工单图表信息
    path(r'chart/', login_required(RenderOrderChartView.as_view()), name='p_render_chart'),
    path(r'chart/orders/', login_required(GetOrderChartView.as_view())),
    path(r'chart/count/', login_required(GetSqlCountView.as_view()), name='p_get_sqlcount'),
]

# -*- coding:utf-8 -*-
# edit by fuzongfei
import logging

from django.http import JsonResponse
from django.shortcuts import render
# Create your views here.
from django.views import View

from sqlquery.forms import GetGrantSchemaForm, GetTableStrucForm, ExecSqlQueryForm, GetHistorySqlForm, \
    GetFilterHistorySqlForm, GetTableIndexForm, GetTableBaseForm

logger = logging.getLogger('django')


class RenderSqlQueryView(View):
    """渲染SQL query页面"""

    def get(self, request):
        return render(request, 'sqlquery/sql_query.html')


class GetGrantSchemaView(View):
    """获取指定环境授权给用户的schema信息"""

    def get(self, request):
        form = GetGrantSchemaForm(request.GET)
        context = None
        if form.is_valid():
            context = form.query(request)
        else:
            error = form.errors.as_text()
            logger.error(error)
        return JsonResponse(context, safe=False)


class GetTableStrucView(View):
    """返回表结构"""

    def get(self, request):
        form = GetTableStrucForm(request.GET)
        if form.is_valid():
            context = form.query()
        else:
            error = form.errors.as_text()
            context = {'status': 2, 'msg': error}

        return JsonResponse(context, safe=False)


class GetTableIndexView(View):
    """返回表索引"""

    def get(self, request):
        form = GetTableIndexForm(request.GET)
        if form.is_valid():
            context = form.query()
        else:
            error = form.errors.as_text()
            context = {'status': 2, 'msg': error}

        return JsonResponse(context, safe=False)


class GetTableBaseView(View):
    """返回表基本信息"""

    def get(self, request):
        form = GetTableBaseForm(request.GET)
        if form.is_valid():
            context = form.query()
        else:
            error = form.errors.as_text()
            context = {'status': 2, 'msg': error}

        return JsonResponse(context, safe=False)


class ExecSqlQueryView(View):
    """执行sql查询"""

    def post(self, request):
        form = ExecSqlQueryForm(request.POST)
        if form.is_valid():
            context = form.execute(request)
        else:
            error = form.errors.as_text()
            context = {'status': 2, 'msg': error}

        return JsonResponse(context, safe=False)


class GetHistorySqlView(View):
    """获取当前用户执行的SQL历史,返回前1000条"""

    def get(self, request):
        form = GetHistorySqlForm(request.GET)
        if form.is_valid():
            context = form.query(request)
        else:
            error = form.errors.as_text()
            context = {'status': 2, 'msg': error}

        return JsonResponse(context, safe=False)

    def post(self, request):
        form = GetFilterHistorySqlForm(request.POST)
        if form.is_valid():
            context = form.query(request)
        else:
            error = form.errors.as_text()
            context = {'status': 2, 'msg': error}

        return JsonResponse(context, safe=False)

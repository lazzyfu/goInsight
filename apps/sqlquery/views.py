from django.http import JsonResponse
from django.shortcuts import render

# Create your views here.
from django.views import View

from sqlorders.models import SqlOrdersEnvironment
from sqlquery.forms import GetSchemasGrantForm, GetStruInfoForm, ExecSqlQueryForm, GetHistorySqlForm, \
    GetFilterHistorySqlForm


class RenderSqlQueryView(View):
    """渲染SQL query页面"""

    def get(self, request, envi_id):
        envi_name = SqlOrdersEnvironment.objects.get(envi_id=envi_id).envi_name
        return render(request, 'sqlquery/sql_query.html', {'envi_id': envi_id, 'envi_name': envi_name})


class GetSchemasGrantView(View):
    """获取指定环境授权给用户的schema信息"""

    def post(self, request):
        form = GetSchemasGrantForm(request.POST)
        if form.is_valid():
            context = form.query(request)
        return JsonResponse(context, safe=False)


class GetStruInfoView(View):
    """返回表结构和索引等信息"""

    def get(self, request):
        form = GetStruInfoForm(request.GET)
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

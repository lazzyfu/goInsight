# Create your views here.
from django_filters.rest_framework import DjangoFilterBackend
from rest_framework import filters
from rest_framework.generics import ListAPIView
from rest_framework.views import APIView

from libs.RenderColumns import render_dynamic_columns
from libs.response import JsonResponseV1
from sqlquery import serializers, models
from libs.Pagination import Pagination


class GetTreeView(APIView):
    def post(self, request, *args, **kwargs):
        serializer = serializers.GetTreeSerializer(data=request.data)
        if serializer.is_valid():
            status, data = serializer.query(request)
            if status:
                return JsonResponseV1(data)
            return JsonResponseV1(message=data, code='0001')
        return JsonResponseV1(message=serializer.errors, code='0001')


class ExecuteQueryView(APIView):
    """执行查询"""

    def post(self, request, *args, **kwargs):
        serializer = serializers.ExecuteQuerySerializer(data=request.data, context={'request': request})
        if serializer.is_valid():
            result = serializer.execute()
            if result['status']:
                return JsonResponseV1(result['data'])
            return JsonResponseV1(message=result['msg'], code='0001')
        return JsonResponseV1(message=serializer.errors, code='0001', flat=True)


class DeleteQueryHashView(APIView):
    """删除查询hash对应的thread_id"""

    def post(self, request, *args, **kwargs):
        serializer = serializers.DeleteQueryHashSerializer(data=request.data, context={'request': request})
        if serializer.is_valid():
            serializer.del_session()
            return JsonResponseV1(code='0000')
        return JsonResponseV1(message=serializer.errors, code='0001')


class GetTableInfoView(APIView):
    """获取表元信息"""

    def post(self, request, *args, **kwargs):
        serializer = serializers.GetTableInfoSerializer(data=request.data, context={'request': request})
        if serializer.is_valid():
            status, data = serializer.query()
            if status:
                return JsonResponseV1(data)
            return JsonResponseV1(message=data, code='0001')
        return JsonResponseV1(message=serializer.errors, code='0001')


class GetHistorySQLView(ListAPIView):
    queryset = models.DbQueryLog.objects.all()
    serializer_class = serializers.GetHistorySQLSerializer
    pagination_class = Pagination
    filter_backends = [DjangoFilterBackend, filters.SearchFilter, filters.OrderingFilter]
    search_fields = ['schema', 'tables']
    ordering = ['-created_at']

    def get(self, request, *args, **kwargs):
        queryset = self.filter_queryset(self.get_queryset().filter(username=request.user.username))
        page = self.paginate_queryset(queryset)
        serializer = self.get_serializer(page, many=True)
        render_columns = [
            {'key': 'host', 'value': '主机'},
            {'key': 'schema', 'value': '库名'},
            {'key': 'tables', 'value': '表名', 'ellipsis': True},
            {'key': 'query_sql', 'value': 'SQL', 'ellipsis': True},
            {'key': 'query_consume_time', 'value': '预估耗时(秒)'},
            {'key': 'query_status', 'value': '状态', 'ellipsis': True},
            {'key': 'affected_rows', 'value': '返回行数'},
            {'key': 'created_at', 'value': '执行时间'},
        ]
        columns = render_dynamic_columns(render_columns)
        data = {'columns': columns, 'data': serializer.data}
        return self.get_paginated_response(data)


class GetDBDictLView(APIView):
    def post(self, request, *args, **kwargs):
        serializer = serializers.GetDBDictSerializer(data=request.data)
        if serializer.is_valid():
            status, data = serializer.query()
            if status:
                return JsonResponseV1(data)
            return JsonResponseV1(message=data, code='0001')
        return JsonResponseV1(message=serializer.errors, code='0001')

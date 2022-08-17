# -*- coding:utf-8 -*-
# edit by xff
import base64
import datetime
# Create your views here.
import json

from django.http import Http404, HttpResponse
from django.utils import timezone
from django_filters.rest_framework import DjangoFilterBackend
from rest_framework import filters
from rest_framework.exceptions import PermissionDenied
from rest_framework.generics import ListAPIView, GenericAPIView, CreateAPIView, UpdateAPIView, DestroyAPIView
from rest_framework.views import APIView
from rest_framework.viewsets import ViewSet

from libs import permissions
from libs.Pagination import Pagination
from libs.RenderColumns import render_dynamic_columns
from libs.response import JsonResponseV1
from sqlorders import models, serializers
from sqlorders.filters import SqlOrderListFilter, GetTasksListFilter


class GetDBEnvironment(ListAPIView):
    queryset = models.DbEnvironment.objects.all()
    serializer_class = serializers.DbEnvironmentSerializer

    # 获取工单环境
    def get(self, request, *args, **kwargs):
        serializer = self.get_serializer(self.get_queryset(), many=True)
        return JsonResponseV1(data=serializer.data)


class GetDbSchemas(APIView):
    # 获取指定环境指定用途的schemas列表
    def get(self, request):
        serializer = serializers.DbSchemasSerializer(data=request.query_params)
        if serializer.is_valid():
            return JsonResponseV1(data=serializer.query)
        return JsonResponseV1(message=serializer.errors, code='0001')


class IncepSyntaxCheckView(APIView):
    def post(self, request, *args, **kwargs):
        serializer = serializers.IncepSyntaxCheckSerializer(data=request.data)

        if serializer.is_valid():
            status, data, msg = serializer.check()
            
            render_columns = [
                {'key': 'level', 'value': '级别'},
                {'key': 'finger_id', 'value': '指纹'},
                {'key': 'type', 'value': '类型'},
                {'key': 'summary', 'value': '提示', 'width': '35%'},
                {'key': 'query', 'value': 'SQL内容', 'width': '25%', 'ellipsis': True},
                {'key': 'affected_rows', 'value': '影响/扫描行数'}
            ]
            columns = render_dynamic_columns(render_columns)
            message = '语法检查不通过,请根据【提示】进行更正'
            if status:
                message = "语法检查通过,可以提交"
            if status is False and data is None:
                message = msg
            innerd = {
                'status': 0 if status else 1,
                'data': data if ((status is False and msg is None) or (status is True)) else None,
            }
            data = {'columns': columns, 'data': innerd}
            return JsonResponseV1(data=data, message=message)
        return JsonResponseV1(message=serializer.errors, code='0001', flat=True)


class SqlOrdersCommit(GenericAPIView):
    permission_classes = (permissions.CanCommitOrdersPermission,)
    serializer_class = serializers.SqlOrdersCommitSerializer

    def post(self, request, *args, **kwargs):
        serializer = self.get_serializer(data=request.data)

        if serializer.is_valid():
            serializer.save()
            return JsonResponseV1(message="提交成功")
        return JsonResponseV1(message=serializer.errors, code='0001', flat=True)


class SqlOrdersList(ListAPIView):
    permission_classes = (permissions.CanViewOrdersPermission,)
    queryset = models.DbOrders.objects.all()
    serializer_class = serializers.SqlOrdersListSerializer
    pagination_class = Pagination
    filter_backends = [DjangoFilterBackend, filters.SearchFilter, filters.OrderingFilter]
    filter_class = SqlOrderListFilter
    ordering = ['-created_at']
    search_fields = ['title', 'database', 'remark', 'applicant', 'progress', 'contents']

    def get(self, request, *args, **kwargs):
        queryset = self.filter_queryset(self.get_queryset())
        page = self.paginate_queryset(queryset)
        serializer = self.get_serializer(page, many=True)
        render_columns = [
            {'key': 'progress', 'value': '进度', 'width': '8%'},
            {'key': 'applicant', 'value': '申请人'},
            {'key': 'department', 'value': '部门'},
            {'key': 'env_name', 'value': '环境'},
            {'key': 'escape_title', 'value': '标题', 'width': '18%', 'ellipsis': True},
            {'key': 'sql_type', 'value': '类型'},
            {'key': 'remark', 'value': '备注'},
            {'key': 'version', 'value': '版本'},
            {'key': 'host', 'value': '实例/库'},
            {'key': 'auditor', 'value': '审核人'},
            {'key': 'reviewer', 'value': '复核人'},
        ]
        columns = render_dynamic_columns(render_columns)
        data = {'columns': columns, 'data': serializer.data}
        return self.get_paginated_response(data)


class SqlOrdersDetail(ListAPIView):
    """SQL工单详情"""
    permission_classes = (permissions.CanViewOrdersPermission,)
    queryset = models.DbOrders.objects.all()
    serializer_class = serializers.SqlOrderDetailSerializer
    lookup_field = 'order_id'

    def get(self, request, *args, **kwargs):
        queryset = self.get_object()
        serializer = self.get_serializer(queryset, context={"request": request})
        return JsonResponseV1(data=serializer.data)


class OpSqlOrderView(ViewSet):
    """更新SQL工单状态，如：审核，关闭等"""
    permission_classes = (permissions.CanViewOrdersPermission,)

    def get_obj(self, pk):
        try:
            obj = models.DbOrders.objects.get(pk=pk)
            return obj
        except models.DbOrders.DoesNotExist:
            raise Http404

    def approve(self, request, pk):
        serializer = serializers.OpSqlOrderSerializer(instance=self.get_obj(pk),
                                                      data=request.data,
                                                      context={"request": request, "handler": "_approve"})

        if serializer.is_valid():
            serializer.save()
            return JsonResponseV1(data=serializer.data, message="操作成功")
        return JsonResponseV1(message=serializer.errors, code='0001')

    def feedback(self, request, pk):
        serializer = serializers.OpSqlOrderSerializer(instance=self.get_obj(pk),
                                                      data=request.data,
                                                      context={"request": request, "handler": "_feedback"})
        if serializer.is_valid():
            serializer.save()
            return JsonResponseV1(data=serializer.data, message="操作成功")
        return JsonResponseV1(message=serializer.errors, code='0001')

    def close(self, request, pk):
        serializer = serializers.OpSqlOrderSerializer(instance=self.get_obj(pk),
                                                      data=request.data,
                                                      context={"request": request, "handler": "_close"})
        if serializer.is_valid():
            serializer.save()
            return JsonResponseV1(data=serializer.data, message="操作成功")
        return JsonResponseV1(message=serializer.errors, code='0001')

    def review(self, request, pk):
        serializer = serializers.OpSqlOrderSerializer(instance=self.get_obj(pk),
                                                      data=request.data,
                                                      context={"request": request, "handler": "_review"})
        if serializer.is_valid():
            serializer.save()
            return JsonResponseV1(data=serializer.data, message="操作成功")
        return JsonResponseV1(message=serializer.errors, code='0001')


class GenerateTasksView(APIView):
    permission_classes = (permissions.CanExecuteOrdersPermission,)

    def post(self, request, *args, **kwargs):
        serializer = serializers.GenerateSqlOrdersTasksSerializer(data=request.data)

        if serializer.is_valid():
            data = serializer.save(request)
            return JsonResponseV1(data=data)
        return JsonResponseV1(message=serializer.errors, code='0001', flat=True)


class GetTaskIdView(APIView):
    def get(self, request, *args, **kwargs):
        """根据order id返回taskid"""
        order_id = kwargs.get('order_id')
        task_id = models.DbOrdersExecuteTasks.objects.filter(order_id=order_id).first().task_id
        return JsonResponseV1(data=task_id)


class GetTasksPreviewView(ListAPIView):
    permission_classes = (permissions.CanViewOrdersPermission,)
    queryset = models.DbOrdersExecuteTasks.objects.all()
    serializer_class = serializers.SqlOrdersTasksListSerializer
    pagination_class = Pagination
    filter_backends = [DjangoFilterBackend, filters.SearchFilter, filters.OrderingFilter]
    filter_class = GetTasksListFilter
    search_fields = ['sql']
    ordering = ['created_time']

    def get(self, request, *args, **kwargs):
        task_id = kwargs.get('task_id')
        queryset = self.filter_queryset(self.get_queryset().filter(task_id=task_id))

        # 数据隐藏按钮打开了
        # 仅允许申请人、审核人、复核人和超权用户查看数据
        obj = models.DbOrders.objects.get(
            pk=models.DbOrdersExecuteTasks.objects.filter(task_id=task_id).first().order_id
        )
        if obj.is_hide == 'ON' and not request.user.is_superuser:
            allowed_view_users = [obj.applicant]
            allowed_view_users.extend([x['user'] for x in json.loads(obj.auditor)])
            allowed_view_users.extend([x['user'] for x in json.loads(obj.reviewer)])
            if request.user.username not in allowed_view_users:
                raise PermissionDenied(detail='您没有权限查看该工单的数据，5s后，自动跳转到工单列表页面')

        origin_queryset = self.queryset.filter(task_id=task_id)
        total = origin_queryset.count()
        progress_0 = origin_queryset.filter(progress=0).count()
        progress_1 = origin_queryset.filter(progress=1).count()
        progress_3 = origin_queryset.filter(progress=3).count()

        page = self.paginate_queryset(queryset)
        serializer = self.get_serializer(page, context={'request': request}, many=True)
        render_columns = [
            {'key': 'num', 'value': '序号'},  # 自定义num，前台显示序号使用
            {'key': 'applicant', 'value': '申请人'},
            {'key': 'sql', 'value': 'SQL', 'ellipsis': True, 'width': '50%'},
            {'key': 'progress', 'value': '进度'},
            {'key': 'result', 'value': '查看结果'},  # 自定义result
        ]
        columns = render_dynamic_columns(render_columns)
        data = {'columns': columns,
                'data': {'data': serializer.data,
                         'total': total,
                         'progress_0': progress_0,
                         'progress_1': progress_1,
                         'progress_3': progress_3}}
        return self.get_paginated_response(data)


class GetTasksListView(ListAPIView):
    permission_classes = (permissions.CanViewOrdersPermission,)
    queryset = models.DbOrdersExecuteTasks.objects.all()
    serializer_class = serializers.SqlOrdersTasksListSerializer
    pagination_class = Pagination
    filter_backends = [DjangoFilterBackend, filters.SearchFilter, filters.OrderingFilter]
    filter_class = GetTasksListFilter
    search_fields = ['sql']
    ordering = ['created_time']

    def get(self, request, *args, **kwargs):
        task_id = kwargs.get('task_id')
        queryset = self.filter_queryset(self.get_queryset().filter(task_id=task_id))

        # 数据隐藏按钮打开了
        # 仅允许申请人、审核人、复核人和超权用户查看数据
        obj = models.DbOrders.objects.get(
            pk=models.DbOrdersExecuteTasks.objects.filter(task_id=task_id).first().order_id
        )
        if obj.is_hide == 'ON' and not request.user.is_superuser:
            allowed_view_users = [obj.applicant]
            allowed_view_users.extend([x['user'] for x in json.loads(obj.auditor)])
            allowed_view_users.extend([x['user'] for x in json.loads(obj.reviewer)])
            if request.user.username not in allowed_view_users:
                raise PermissionDenied(detail='您没有权限查看该工单的数据，5s后，自动跳转到工单列表页面')

        page = self.paginate_queryset(queryset)
        serializer = self.get_serializer(page, context={'request': request}, many=True)
        render_columns = [
            {'key': 'num', 'value': '序号'},  # 自定义num，前台显示序号使用
            {'key': 'applicant', 'value': '申请人'},
            {'key': 'sql', 'value': 'SQL', 'ellipsis': True, 'width': '50%'},
            {'key': 'progress', 'value': '进度'},
            {'key': 'execute', 'value': '执行'},  # 自定义execute
            {'key': 'result', 'value': '查看结果'},  # 自定义result
        ]
        if queryset.exists():
            if queryset.first().sql_type == 'DDL':
                render_columns.insert(-1, {'key': 'ghost_pause', 'value': '暂停(gh-ost)'})
                render_columns.insert(-1, {'key': 'ghost_recovery', 'value': '恢复(gh-ost)'})
        columns = render_dynamic_columns(render_columns)
        data = {'columns': columns, 'data': serializer.data}
        return self.get_paginated_response(data)


class ExecuteSingleTaskView(APIView):
    permission_classes = (permissions.CanExecuteOrdersPermission,)

    def post(self, request, *args, **kwargs):
        serializer = serializers.ExecuteSingleTaskSerializer(data=request.data)

        if serializer.is_valid():
            serializer.execute(request)
            return JsonResponseV1(message="任务提交成功，请查看输出")
        return JsonResponseV1(message=serializer.errors, code='0001', flat=True)


class ExecuteMultiTasksView(APIView):
    permission_classes = (permissions.CanExecuteOrdersPermission,)

    def post(self, request, *args, **kwargs):
        serializer = serializers.ExecuteMultiTasksSerializer(data=request.data)

        if serializer.is_valid():
            serializer.execute(request)
            return JsonResponseV1(message="任务提交成功，请查看输出")
        return JsonResponseV1(message=serializer.errors, code='0001', flat=True)


class ThrottleTaskView(APIView):
    permission_classes = (permissions.CanExecuteOrdersPermission,)

    def post(self, request, *args, **kwargs):
        serializer = serializers.ThrottleTaskSerializer(data=request.data)

        if serializer.is_valid():
            message = serializer.execute(request)
            return JsonResponseV1(message=message)
        return JsonResponseV1(message=serializer.errors, code='0001', flat=True)


class GetTasksResultView(ListAPIView):
    """SQL工单详情"""
    permission_classes = (permissions.CanViewOrdersPermission,)
    queryset = models.DbOrdersExecuteTasks.objects.all()
    serializer_class = serializers.GetTasksResultSerializer
    lookup_field = 'id'

    def get(self, request, *args, **kwargs):
        queryset = self.get_object()
        serializer = self.get_serializer(queryset, context={"request": request})
        return JsonResponseV1(data=serializer.data)


class HookSqlOrdersView(APIView):
    permission_classes = (permissions.anyof(permissions.CanCommitOrdersPermission,
                                            permissions.CanViewOrdersPermission,
                                            permissions.CanExecuteOrdersPermission,
                                            permissions.CanAuditOrdersPermission),
                          )

    def post(self, request, *args, **kwargs):
        serializer = serializers.HookSqlOrdersSerializer(data=request.data)

        if serializer.is_valid():
            serializer.save()
            return JsonResponseV1(message="任务提交成功，请查看输出")
        return JsonResponseV1(message=serializer.errors, code='0001', flat=True)


class DownloadExportFilesView(APIView):
    """下载导出文件"""
    permission_classes = (permissions.CanViewOrdersPermission,)

    def get(self, request, base64_filename):
        file_name = base64.b64decode(base64_filename).decode()

        if not models.DbExportFiles.objects.filter(file_name=file_name).exists():
            raise Http404

        obj = models.DbExportFiles.objects.get(file_name=file_name)
        if not models.DbOrdersExecuteTasks.objects.get(pk=obj.task_id).applicant == request.user.username:
            raise PermissionDenied(detail='您没有权限')

        fsock = open(f"media/{obj.files}", 'rb')
        response = HttpResponse(fsock, content_type="application/zip")
        response['Content-Disposition'] = f'attachment; filename={file_name}'
        return response


class ReleaseVersionsGet(APIView):
    """获取上线版本号，提交工单使用"""

    def get(self, request):
        before_30_days = (timezone.now() - datetime.timedelta(days=30))
        queryset = models.ReleaseVersions.objects.filter(
            expire_time__gte=before_30_days
        ).values('id', 'version', 'expire_time').order_by('-created_at')
        for row in queryset:
            row['disabled'] = 0
            if row['expire_time'] < datetime.datetime.date(timezone.now()):
                row['disabled'] = 1
        return JsonResponseV1(data=queryset)


class ReleaseVersionsList(ListAPIView):
    """获取上线版本号列表，管理上线版本号使用"""
    permission_classes = (permissions.CanViewVersionPermission,)
    queryset = models.ReleaseVersions.objects.all()
    serializer_class = serializers.ReleaseVersionsListSerializer
    pagination_class = Pagination
    filter_backends = [filters.SearchFilter, filters.OrderingFilter]
    search_fields = ['username', 'version', 'expire_time']
    ordering = ['-created_at']

    def get(self, request, *args, **kwargs):
        queryset = self.filter_queryset(self.get_queryset())
        page = self.paginate_queryset(queryset)
        serializer = self.get_serializer(page, many=True)
        render_columns = [
            {'key': 'version', 'value': '版本'},
            {'key': 'username', 'value': '创建人'},
            {'key': 'expire_time', 'value': '截止日期'},
            {'key': 'created_at', 'value': '创建时间'},
            {'key': 'key', 'value': '操作'},
            {'key': 'id', 'value': '详情'},
        ]
        columns = render_dynamic_columns(render_columns)
        data = {'columns': columns, 'data': serializer.data}
        return self.get_paginated_response(data)


class ReleaseVersionsCreate(CreateAPIView):
    """创建版本"""
    permission_classes = (permissions.CanCreateVersionsPermission,)
    serializer_class = serializers.ReleaseVersionsCreateSerializer

    def create(self, request, *args, **kwargs):
        serializer = self.get_serializer(data=request.data)
        if serializer.is_valid():
            self.perform_create(serializer)
            return JsonResponseV1(message="创建成功")
        return JsonResponseV1(code='0001', message=serializer.errors, flat=True)


class ReleaseVersionsUpdate(UpdateAPIView):
    """更新版本号，该类只更新单条记录"""
    permission_classes = (permissions.CanUpdateVersionsPermission,)

    def put(self, request, *args, **kwargs):
        serializer = serializers.ReleaseVersionsSerializer(
            instance=models.ReleaseVersions.objects.get(pk=kwargs['key']),  # 返回单条记录
            data=request.data
        )
        if serializer.is_valid():
            serializer.save()
            return JsonResponseV1(message="更新成功")
        return JsonResponseV1(code='0001', message=serializer.errors, flat=True)


class ReleaseVersionsDelete(DestroyAPIView):
    """删除版本"""
    permission_classes = (permissions.CanDeleteVersionsPermission,)
    queryset = models.ReleaseVersions.objects.all()
    lookup_field = 'id'  # 默认为主键，可不写

    def destroy(self, request, *args, **kwargs):
        instance = self.get_object()
        self.perform_destroy(instance)
        return JsonResponseV1(message="删除成功")


class ReleaseVersionsView(APIView):
    """获取指定版本内工单在所有环境的进度"""

    def get(self, request, *args, **kwargs):
        # 获取版本对应的主键
        version = kwargs.get('version')
        version_id = models.ReleaseVersions.objects.get(version=version).pk
        # 获取环境，行转为动态列
        obj = models.DbEnvironment.objects.values('id', 'name')
        row2columns = ''
        for row in obj:
            row2columns += f"max(if(env_id={row['id']}, progress, -1)) as {row['name']},"
        # 获取任务下所有工单分别在各个环境中的状态，此处的环境为动态环境
        # id没有实际意义
        query = f"select " + row2columns + \
                f"substring(MD5(RAND()),1,20) as id,title as escape_title,order_id, applicant " \
                f"from yasql_dborders where version_id='{version_id}' group by escape_title,order_id,applicant"
        rawquery = models.DbOrders.objects.raw(query)
        # 获取环境列名
        dynamic_columns = list(rawquery.columns)[:-4]
        data = []
        for row in rawquery:
            columns = {
                'id': row.id,
                'escape_title': row.escape_title,
                'order_id': row.order_id,
                'applicant': row.applicant,
            }
            for col in dynamic_columns:
                columns[col] = getattr(row, col)
            data.append(columns)

        render_columns = [
            {'key': 'escape_title', 'ellipsis': True, 'value': '标题'},
            {'key': 'applicant', 'value': '申请人'},
        ]
        render_columns.extend([{'key': x, 'value': x} for x in dynamic_columns])
        columns = render_dynamic_columns(render_columns)
        data = {'columns': columns, 'data': data}
        return JsonResponseV1(data=data)

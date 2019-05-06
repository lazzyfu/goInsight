# -*- coding:utf-8 -*-
# edit by fuzongfei
import json

from django.db.models import Case, When, Value, CharField
from django.utils import timezone
from rest_framework import status
from rest_framework.renderers import TemplateHTMLRenderer
from rest_framework.response import Response
from rest_framework.views import APIView

from config.config import EMAIL
from orders.models import SysEnvironment, Orders, OrdersTasks
from orders.permissions import CanCommitPermission, anyof, CanExecutePermission, CanAuditPermission
from orders.serializers.commitSerializers import OrdersCommitSerializer, OrderReplySerializer, HookOrdersSerializer
from orders.serializers.listSerializers import OrderListSerializer, GetOrderReplySerializer, MyOrderListSerializer


class OnlineVersionView(APIView):
    """渲染上线版本号页面"""
    renderer_classes = [TemplateHTMLRenderer]
    template_name = 'orders/online_version.html'

    def get(self, request):
        return Response()


class SQLOrdersView(APIView):
    """渲染DML和DDL工单页面"""
    renderer_classes = [TemplateHTMLRenderer]
    template_name = 'orders/sql_orders.html'

    def get(self, request):
        is_enable_mail = 0 if EMAIL['enable'] else 1
        return Response(data={'is_enable_mail': is_enable_mail})


class SQLExportView(APIView):
    """渲染SQL导出工单"""
    renderer_classes = [TemplateHTMLRenderer]
    template_name = 'orders/export_orders.html'

    def get(self, request):
        return Response()


class MyOrdersView(APIView):
    """渲染SQL导出工单"""
    renderer_classes = [TemplateHTMLRenderer]
    template_name = 'orders/my_orders.html'

    def get(self, request):
        return Response()


class OpsOrdersView(APIView):
    """渲染运维工单"""
    renderer_classes = [TemplateHTMLRenderer]
    template_name = 'orders/ops_orders.html'

    def get(self, request):
        return Response()


class OrdersCommitView(APIView):
    """提交工单"""

    permission_classes = (CanCommitPermission,)

    def post(self, request):
        serializer = OrdersCommitSerializer(data=request.data)

        if serializer.is_valid():
            s, data = serializer.save(request)
            if s:
                data = {'code': 0, 'data': data}
            else:
                data = {'code': 1, 'data': data}
            return Response(data=data, status=status.HTTP_200_OK)
        else:
            errors = [str(v[0]) for k, v in serializer.errors.items()]
            data = {'code': 2, 'data': '\n'.join(errors)}
            return Response(data=data, status=status.HTTP_400_BAD_REQUEST)


class RenderOrdersEnviView(APIView):
    """渲染工单页面"""
    renderer_classes = [TemplateHTMLRenderer]
    template_name = 'orders/orders_list.html'

    def get(self, request, envi_id):
        envi_name = SysEnvironment.objects.get(envi_id=envi_id).envi_name
        return Response(data={'envi_id': envi_id, 'envi_name': envi_name})


class OrdersListView(APIView):
    """分页获取指定环境的工单列表"""

    def post(self, request):
        serializer = OrderListSerializer(data=request.data)
        if serializer.is_valid():
            data = serializer.query()
            return Response(data=data, status=status.HTTP_200_OK)
        else:
            errors = [str(v[0]) for k, v in serializer.errors.items()]
            data = {'code': 2, 'data': '\n'.join(errors)}
            return Response(data=data, status=status.HTTP_400_BAD_REQUEST)


class MyOrdersListView(APIView):
    """分页获取指定环境的工单列表"""

    def post(self, request):
        serializer = MyOrderListSerializer(data=request.data)
        if serializer.is_valid():
            data = serializer.query(request)
            return Response(data=data, status=status.HTTP_200_OK)
        else:
            errors = [str(v[0]) for k, v in serializer.errors.items()]
            data = {'code': 2, 'data': '\n'.join(errors)}
            return Response(data=data, status=status.HTTP_400_BAD_REQUEST)


class OrdersDetailsView(APIView):
    """
    工单详情页面
    """

    renderer_classes = [TemplateHTMLRenderer]
    template_name = 'orders/orders_detail.html'

    def get(self, request, id):
        queryset = Orders.objects.annotate(
            progress_value=Case(
                When(progress='0', then=Value('待批准')),
                When(progress='1', then=Value('未批准')),
                When(progress='2', then=Value('已批准')),
                When(progress='3', then=Value('处理中')),
                When(progress='4', then=Value('已完成')),
                When(progress='5', then=Value('已关闭')),
                When(progress='6', then=Value('已复核')),
                When(progress='7', then=Value('已勾住')),
                output_field=CharField(),
            ),
        ).get(id=id)
        queryset.auditor = json.loads(queryset.auditor)
        queryset.reviewer = json.loads(queryset.reviewer)
        if queryset.close_info:
            queryset.close_info = json.loads(queryset.close_info)
        return Response(data={'contents': queryset}, status=status.HTTP_200_OK)


class OrderExecuteDetailsView(APIView):
    def get(self, request):
        id = request.GET.get('id')
        queryset = OrdersTasks.objects.annotate(
            status=Case(
                When(task_progress='0', then=Value('未执行')),
                When(task_progress='1', then=Value('已完成')),
                When(task_progress='2', then=Value('处理中')),
                When(task_progress='3', then=Value('失败')),
                When(task_progress='4', then=Value('异常')),
                output_field=CharField(),
            )
        ).filter(order__id=id).values('taskid', 'sql', 'status', 'affected_row', 'consume_time', 'executor',
                                      'execition_time')

        field = [{'field': 'sid', 'title': 'ID'},
                 {'field': 'sql', 'title': 'SQL'},
                 {'field': 'status', 'title': '状态'},
                 {'field': 'rollback_sql', 'title': '回滚SQL'},
                 {'field': 'consume_time', 'title': '耗时(s)'},
                 {'field': 'executor', 'title': '执行人'},
                 {'field': 'execition_time', 'title': '执行时间'},
                 ]

        i = 1
        data = []
        for row in queryset:
            data.append({
                'sid': i,
                'sql': row['sql'][:50] + ' ...',
                'status': row['status'],
                'affected_row': row['affected_row'],
                'rollback_sql': f"<a target='_bank' href='/orders/subtasks/list/{row['taskid']}'>查看</a>",
                'consume_time': row['consume_time'],
                'executor': row['executor'],
                'execition_time': timezone.localtime(row['execition_time']).strftime('%Y-%m-%d %H:%M:%S')

            })
            i += 1
        context = {'columns': field, 'data': data}
        return Response(data=context, status=status.HTTP_200_OK)


class OrderReplyView(APIView):
    """
    工单回复
    """
    # 拥有工单提交/执行/审核权限的用户可以操作
    permission_classes = (anyof(CanCommitPermission, CanExecutePermission, CanAuditPermission),)

    def post(self, request):
        serializer = OrderReplySerializer(data=request.data)
        if serializer.is_valid():
            s, data = serializer.save(request)
            return Response(data={'code': 0, 'data': data}, status=status.HTTP_200_OK)
        else:
            errors = [str(v[0]) for k, v in serializer.errors.items()]
            data = {'code': 2, 'data': '\n'.join(errors)}
            return Response(data=data, status=status.HTTP_400_BAD_REQUEST)


class GetOrderReplyView(APIView):
    """
    获取工单回复的内容
    """
    # 拥有工单提交/执行/审核权限的用户可以操作
    permission_classes = (anyof(CanCommitPermission, CanExecutePermission, CanAuditPermission),)

    def post(self, request):
        serializer = GetOrderReplySerializer(data=request.data)
        if serializer.is_valid():
            data = serializer.query()
            return Response(data={'code': 0, 'data': data}, status=status.HTTP_200_OK)
        else:
            errors = [str(v[0]) for k, v in serializer.errors.items()]
            data = {'code': 2, 'data': '\n'.join(errors)}
            return Response(data=data, status=status.HTTP_400_BAD_REQUEST)


class HookOrdersView(APIView):
    """
    工单的HOOK功能
    """
    # 拥有工单提交/执行/审核权限的用户可以操作
    permission_classes = (anyof(CanCommitPermission, CanExecutePermission, CanAuditPermission),)

    def post(self, request):
        serializer = HookOrdersSerializer(data=request.data)
        if serializer.is_valid():
            s, data = serializer.save(request)
            code = 0 if s else 2
            return Response(data={'code': code, 'data': data}, status=status.HTTP_200_OK)
        else:
            errors = [str(v[0]) for k, v in serializer.errors.items()]
            data = {'code': 2, 'data': '\n'.join(errors)}
            return Response(data=data, status=status.HTTP_400_BAD_REQUEST)

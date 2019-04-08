# Create your views here.
from rest_framework import status
from rest_framework.renderers import TemplateHTMLRenderer
from rest_framework.response import Response
from rest_framework.views import APIView

from dash.serializers import GetSqlCountSerializer
from orders.models import Orders


class RenderOrderChartView(APIView):
    """渲染图表页面"""
    renderer_classes = [TemplateHTMLRenderer]
    template_name = 'dash/dash.html'

    def get(self, request):
        return Response()


class GetOrderChartView(APIView):
    def get(self, request):
        my_queryset = Orders.objects.filter(applicant=request.user.username)
        my_order_count = my_queryset.count()
        my_order_dml_count = my_queryset.filter(sql_type='DML').count()
        my_order_ddl_count = my_queryset.filter(sql_type='DDL').count()
        my_order_ops_count = my_queryset.filter(sql_type='OPS').count()
        my_order_export_count = my_queryset.filter(sql_type='EXPORT').count()

        myordercharts = [
            {
                'value': my_order_dml_count,
                'color': '#f56954',
                'highlight': '#f56954',
                'label': 'DML工单'
            },
            {
                'value': my_order_ddl_count,
                'color': '#00a65a',
                'highlight': '#00a65a',
                'label': 'DDL工单'
            },
            {
                'value': my_order_ops_count,
                'color': '#00c0ef',
                'highlight': '#00c0ef',
                'label': '运维工单'
            },
            {
                'value': my_order_export_count,
                'color': '#f39c12',
                'highlight': '#f39c12',
                'label': '导出工单'
            }
        ]

        platform_order_count = Orders.objects.count()
        platform_order_dml_count = Orders.objects.filter(sql_type='DML').count()
        platform_order_ddl_count = Orders.objects.filter(sql_type='DDL').count()
        platform_order_ops_count = Orders.objects.filter(sql_type='OPS').count()
        platform_order_export_count = Orders.objects.filter(sql_type='EXPORT').count()
        platformordercharts = [
            {
                'value': platform_order_dml_count,
                'color': '#f56954',
                'highlight': '#f56954',
                'label': 'DML工单'
            },
            {
                'value': platform_order_ddl_count,
                'color': '#00a65a',
                'highlight': '#00a65a',
                'label': 'DDL工单'
            },
            {
                'value': platform_order_ops_count,
                'color': '#00c0ef',
                'highlight': '#00c0ef',
                'label': '运维工单'
            },
            {
                'value': platform_order_export_count,
                'color': '#f39c12',
                'highlight': '#f39c12',
                'label': '导出工单'
            }
        ]

        context = {'my_order_count': my_order_count,
                   'platform_order_count': platform_order_count,
                   'myordercharts': myordercharts,
                   'platformordercharts': platformordercharts}

        return Response(data=context, status=status.HTTP_200_OK)


class GetSqlCountView(APIView):
    def post(self, request):
        serializer = GetSqlCountSerializer(data=request.POST)
        if serializer.is_valid():
            data = serializer.query()
            return Response(data={'code': 0, 'data': data}, status=status.HTTP_200_OK)
        else:
            errors = [str(v[0]) for k, v in serializer.errors.items()]
            data = {'code': 2, 'data': '\n'.join(errors)}
            return Response(data=data, status=status.HTTP_400_BAD_REQUEST)

# -*- coding:utf-8 -*-
# edit by fuzongfei

from django.db.models import Case, When, Value, CharField, Q, F
from rest_framework import serializers

from orders.models import Orders, OrderReply
from orders.utils.validatorsCheck import envi_validator


class OrderListSerializer(serializers.Serializer):
    envi_id = serializers.CharField(required=True, validators=[envi_validator],
                                    error_messages={'required': 'envi_id字段不能为空'})
    limit_size = serializers.IntegerField(required=True, label=u'每页显示数量', error_messages={'required': '每页数量不能为空'})
    offset_size = serializers.IntegerField(required=True, label=u'分页偏移量', error_messages={'required': '分页偏移量不能为空'})
    search_content = serializers.CharField(max_length=128, required=False, label='搜索内容')

    def query(self):
        sdata = self.validated_data
        envi_id = sdata.get('envi_id')
        limit_size = sdata.get('limit_size')
        offset_size = sdata.get('offset_size')
        search_content = sdata.get('search_content')

        query = Orders.objects.filter(envi_id=envi_id).annotate(
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
            progress_color=Case(
                When(progress__in=('0',), then=Value('btn-primary')),
                When(progress__in=('2',), then=Value('btn-warning')),
                When(progress__in=('1', '5'), then=Value('btn-danger')),
                When(progress__in=('3',), then=Value('bg-navy')),
                When(progress__in=('4',), then=Value('btn-success')),
                When(progress__in=('6',), then=Value('bg-purple')),
                When(progress__in=('7',), then=Value('btn-default')),
                output_field=CharField(),
            ),
            task_version=Case(
                When(version__version__isnull=True, then=Value('')),
                When(version__version__isnull=False, then=F('version__version')),
                output_field=CharField(),
            )
        )
        if search_content:
            obj = query.filter(Q(version__icontains=search_content) | Q(title__icontains=search_content) | Q(
                applicant__icontains=search_content) | Q(auditor__icontains=search_content) | Q(
                reviewer__icontains=search_content) | Q(host__icontains=search_content) | Q(
                database__icontains=search_content) | Q(contents__icontains=search_content))
        else:
            obj = query

        total = obj.count()
        data = obj.values('id', 'envi_id', 'task_version', 'host', 'port', 'database', 'sql_type',
                          'title', 'progress_value', 'progress_color', 'remark',
                          'applicant', 'auditor', 'reviewer', 'created_at'
                          ).order_by('-created_at')[offset_size:limit_size]
        result = {'total': total, 'rows': data}
        return result


class MyOrderListSerializer(serializers.Serializer):
    limit_size = serializers.IntegerField(required=True, label=u'每页显示数量', error_messages={'required': '每页数量不能为空'})
    offset_size = serializers.IntegerField(required=True, label=u'分页偏移量', error_messages={'required': '分页偏移量不能为空'})
    search_content = serializers.CharField(max_length=128, required=False, label='搜索内容')

    def query(self, request):
        sdata = self.validated_data
        limit_size = sdata.get('limit_size')
        offset_size = sdata.get('offset_size')
        search_content = sdata.get('search_content')

        query = Orders.objects.filter(applicant=request.user.username).annotate(
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
            progress_color=Case(
                When(progress__in=('0',), then=Value('btn-primary')),
                When(progress__in=('2',), then=Value('btn-warning')),
                When(progress__in=('1', '5'), then=Value('btn-danger')),
                When(progress__in=('3',), then=Value('bg-navy')),
                When(progress__in=('4',), then=Value('btn-success')),
                When(progress__in=('6',), then=Value('bg-purple')),
                When(progress__in=('7',), then=Value('btn-default')),
                output_field=CharField(),
            ),
            task_version=Case(
                When(version__version__isnull=True, then=Value('')),
                When(version__version__isnull=False, then=F('version__version')),
                output_field=CharField(),
            ),
            envi_name=F('envi__envi_name')
        )
        if search_content:
            obj = query.filter(Q(version__icontains=search_content) | Q(title__icontains=search_content) | Q(
                applicant__icontains=search_content) | Q(auditor__icontains=search_content) | Q(
                reviewer__icontains=search_content) | Q(host__icontains=search_content) | Q(
                database__icontains=search_content) | Q(contents__icontains=search_content))
        else:
            obj = query

        total = obj.count()
        data = obj.values('id', 'envi_id', 'envi_name', 'task_version', 'host', 'port', 'database', 'sql_type',
                          'title', 'progress_value', 'progress_color', 'remark',
                          'applicant', 'auditor', 'reviewer', 'created_at'
                          ).order_by('-created_at')[offset_size:limit_size]
        result = {'total': total, 'rows': data}
        return result


class GetOrderReplySerializer(serializers.Serializer):
    reply_id = serializers.IntegerField(required=True, error_messages={'required': '回复工单id不能为空'})

    def query(self):
        sdata = self.validated_data
        reply_id = sdata.get('reply_id')
        queryset = OrderReply.objects.annotate(
            username=F('user__username'),
            avatar_file=F('user__avatar_file'),
        ).filter(reply__id=reply_id).values('username', 'avatar_file', 'reply_contents', 'created_at').order_by(
            '-created_at')
        return queryset

# -*- coding:utf-8 -*-
# by fuzongfei

import django_filters

from sqlorders import models


class SqlOrderListFilter(django_filters.rest_framework.FilterSet):
    username = django_filters.CharFilter(field_name="applicant", lookup_expr="iexact")
    progress = django_filters.NumberFilter(field_name="progress", lookup_expr="iexact")
    env = django_filters.NumberFilter(field_name="env__id", lookup_expr="iexact")
    start_created_at = django_filters.DateTimeFilter(field_name="created_at", lookup_expr="gte")
    end_created_at = django_filters.DateTimeFilter(field_name="created_at", lookup_expr="lte")

    class Meta:
        model = models.DbOrders
        fields = ["username", "progress", "env", "start_created_at", "end_created_at"]


class GetTasksListFilter(django_filters.rest_framework.FilterSet):
    progress = django_filters.NumberFilter(field_name="progress", lookup_expr="iexact")

    class Meta:
        model = models.DbOrdersExecuteTasks
        fields = ['progress']

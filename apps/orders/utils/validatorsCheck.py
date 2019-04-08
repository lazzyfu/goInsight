# -*- coding:utf-8 -*-
# edit by fuzongfei
from rest_framework import serializers

from orders.models import SysEnvironment, OnlineVersion


def envi_validator(value):
    value = value if isinstance(value, int) else int(value)
    envi = [x for x in list(SysEnvironment.objects.all().values_list('envi_id', flat=True))]
    if value not in envi:
        raise serializers.ValidationError('请选择正确的工单环境')


def sql_type_validator(value):
    allowed_sql_type = ['DML', 'DDL', 'OPS', 'EXPORT']
    if value not in allowed_sql_type:
        raise serializers.ValidationError('传入类型错误')


def file_foramt_validator(value):
    allowed_file_foramt = ['xlsx', 'csv']
    if value not in allowed_file_foramt:
        raise serializers.ValidationError('仅支持xlsx和csv文件类型')


def online_version_validator(value):
    versions = OnlineVersion.objects.all().values_list('version', flat=True)
    if value not in versions:
        raise serializers.ValidationError('上线版本号错误')


def subtask_stop_validator(value):
    allowed_actions = ['pause_ghost', 'recovery_ghost', 'stop_ghost']
    if value not in allowed_actions:
        raise serializers.ValidationError('子任务停止操作错误')

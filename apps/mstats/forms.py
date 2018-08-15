# -*- coding:utf-8 -*-
# edit by fuzongfei
import json

from django import forms
from django_celery_beat.models import PeriodicTask, CrontabSchedule


class PrivModifyForm(forms.Form):
    db_host = forms.CharField(max_length=128, min_length=3, required=True)
    user = forms.CharField(max_length=30, min_length=1, required=True)
    action = forms.ChoiceField(
        choices=(('modify_privileges', u'更改权限'),
                 ('new_host', u'新建主机'),
                 ('delete_host', u'删除主机'),
                 ('new_user', u'新建用户')))


class SchemaMonitorForm(forms.Form):
    name = forms.CharField(max_length=64, min_length=3, required=True)
    host = forms.CharField(max_length=64, min_length=4, required=True)
    schema = forms.CharField(max_length=64, min_length=1, required=True)
    crontab = forms.IntegerField(required=True)
    receiver = forms.CharField(max_length=20480, min_length=1, required=True)

    def is_save(self):
        cleaned_data = super(SchemaMonitorForm, self).clean()
        name = cleaned_data.get('name')
        host = cleaned_data.get('host')
        schema = cleaned_data.get('schema')
        crontab = cleaned_data.get('crontab')
        receiver = cleaned_data.get('receiver')
        task = 'scheduled_tasks.tasks.monitor_schema_modify'

        kwargs = {'host': host, 'schema': schema, 'receiver': receiver}

        if PeriodicTask.objects.filter(name=name).first():
            context = {'status': 2, 'msg': '同名任务已经存在'}
        else:
            PeriodicTask.objects.create(
                name=name,
                task=task,
                crontab=CrontabSchedule.objects.get(pk=crontab),
                kwargs=json.dumps(kwargs),
                description=u'表结构监控'
            )
            context = {'status': 0, 'msg': '任务创建成功'}

        return context

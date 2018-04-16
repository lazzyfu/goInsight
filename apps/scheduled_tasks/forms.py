# -*- coding:utf-8 -*-
# edit by fuzongfei
import json

from celery import current_app
from django import forms
from djcelery.models import CrontabSchedule, PeriodicTask

# 自动获取任务列表，并提供form选择
celery_app = current_app
tasks = list(sorted(name for name in celery_app.tasks if not name.startswith('celery.')))
tasks_choice = tuple(zip(tasks, tasks))


class PeriodicForm(forms.Form):
    name = forms.CharField(max_length=64, min_length=3, required=True)
    host = forms.CharField(max_length=64, min_length=4, required=True)
    schema = forms.CharField(max_length=64, min_length=1, required=True)
    crontab = forms.IntegerField(required=True)
    receiver = forms.CharField(max_length=256, min_length=1, required=True)
    task = forms.ChoiceField(choices=tasks_choice)
    enabled = forms.ChoiceField(choices=(('0', u'禁用'), ('1', u'启用')))

    def is_save(self):
        cleaned_data = super(PeriodicForm, self).clean()
        name = cleaned_data.get('name')
        host = cleaned_data.get('host')
        schema = cleaned_data.get('schema')
        crontab = cleaned_data.get('crontab')
        receiver = cleaned_data.get('receiver')
        task = cleaned_data.get('task')
        enabled = cleaned_data.get('enabled')

        kwargs = {'host': host, 'schema': schema, 'receiver': receiver}

        if PeriodicTask.objects.filter(name=name).first():
            context = {'status': 2, 'msg': '同名任务已经存在'}
        else:
            PeriodicTask.objects.create(
                name=name,
                task=task,
                enabled=enabled,
                crontab=CrontabSchedule.objects.get(pk=crontab),
                kwargs=json.dumps(kwargs)
            )
            context = {'status': 0, 'msg': '任务创建成功'}

        return context


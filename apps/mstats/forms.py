# -*- coding:utf-8 -*-
# edit by fuzongfei
import json

from django import forms
# from djcelery.models import PeriodicTask, CrontabSchedule
from django_celery_beat.models import PeriodicTask, CrontabSchedule

from mstats.utils import CheckCParserValid, GeneralBackupCmd


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


class BackupTaskForm(forms.Form):
    name = forms.CharField(max_length=64, min_length=3, required=True)
    ssh_host = forms.CharField(max_length=64, min_length=4, required=True)
    ssh_user = forms.CharField(max_length=64, min_length=4, required=True)
    ssh_password = forms.CharField(max_length=64, min_length=4, required=True)
    ssh_port = forms.IntegerField(required=True)
    crontab = forms.IntegerField(required=True)
    backup_method = forms.ChoiceField(choices=(
        ('mysqldump', 'mysqldump'),
        ('xtrabackup', 'xtrabackup'),
        ('mysqldump,xtrabackup', 'mysqldump,xtrabackup')
    ))
    backup_dir = forms.CharField(max_length=256, min_length=2, required=True)
    backup_args = forms.CharField(max_length=8192, min_length=0, required=True)

    def is_save(self):
        cleaned_data = super(BackupTaskForm, self).clean()
        name = cleaned_data.get('name')
        ssh_host = cleaned_data.get('ssh_host')
        ssh_user = cleaned_data.get('ssh_user')
        ssh_password = cleaned_data.get('ssh_password')
        ssh_port = cleaned_data.get('ssh_port')
        crontab = cleaned_data.get('crontab')
        task = 'mstats.tasks.backup_schema'
        backup_dir = cleaned_data.get('backup_dir')
        backup_args = cleaned_data.get('backup_args')

        check_pass = CheckCParserValid(ssh_user=ssh_user, ssh_password=ssh_password,
                                       ssh_host=ssh_host, ssh_port=ssh_port,
                                       backup_dir=backup_dir, parser_string=backup_args).run()
        if check_pass is True:
            backup_cmd = GeneralBackupCmd(ssh_user=ssh_user, ssh_password=ssh_password,
                                          ssh_host=ssh_host, ssh_port=ssh_port,
                                          backup_dir=backup_dir, parser_string=backup_args).run()
            kwargs = {
                'ssh_host': ssh_host,
                'ssh_user': ssh_user,
                'ssh_password': ssh_password,
                'ssh_port': ssh_port,
                'backup_dir': backup_dir,
                'backup_cmd': backup_cmd
            }

            if PeriodicTask.objects.filter(name=name).first():
                context = {'status': 2, 'msg': '同名任务已经存在'}
            else:
                PeriodicTask.objects.create(
                    name=name,
                    task=task,
                    crontab=CrontabSchedule.objects.get(pk=crontab),
                    kwargs=json.dumps(kwargs),
                    description=u'数据库备份'
                )
                context = {'status': 0, 'msg': '任务创建成功'}
        else:
            context = check_pass

        return context

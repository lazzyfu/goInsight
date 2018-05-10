import json
import re

import os
from django.db import transaction
from django.http import HttpResponse, JsonResponse
from django.shortcuts import render
# Create your views here.
from django.utils.decorators import method_decorator
from django.views import View
from djcelery.models import PeriodicTask, CrontabSchedule

from project_manager.models import InceptionHostConfig
from user_manager.permissions import permission_required
from mstats.forms import PrivModifyForm, BackupForm
from mstats.utils import get_mysql_user_info, check_mysql_conn_status, MySQLuser_manager, ParamikoOutput
from utils.tools import format_request


class RenderMySQLUserView(View):
    @permission_required('can_mysqluser_view')
    def get(self, request):
        return render(request, 'mysql_user_manager.html')


class MySQLUserView(View):
    @permission_required('can_mysqluser_view')
    @method_decorator(check_mysql_conn_status)
    def get(self, request):
        data = format_request(request)
        host = data.get('host')
        data = get_mysql_user_info(host)

        return HttpResponse(json.dumps(data))


class MysqlUserManager(View):
    @permission_required('can_mysqluser_edit')
    @transaction.atomic
    def post(self, request):
        data = format_request(request)
        form = PrivModifyForm(data)
        context = {}
        if form.is_valid():
            cleaned_data = form.cleaned_data
            db_host = cleaned_data.get('db_host')
            user = cleaned_data.get('user')
            action = cleaned_data.get('action')

            host = data.get('host')
            password = data.get('password')
            schema = data.get('schema')
            privileges = data.get('privileges')

            username = user + '@' + '"' + host + '"'

            data = InceptionHostConfig.objects.get(host=db_host)
            protection_user = []
            if len(list(data.protection_user.split(','))) == 1:
                protection_user = data.protection_user.split(',')
                protection_user.append('')
            else:
                protection_user = data.protection_user.split(',')
            protection_user_tuple = tuple([x.strip() for x in protection_user])

            if user in protection_user_tuple:
                context = {'status': 1, 'msg': f'该用户({user})已被保护，无法操作'}
            else:
                mysql_user_mamager = MySQLuser_manager(locals())
                if action == "modify_privileges":
                    context = mysql_user_mamager.priv_modify()
                elif action == "new_host":
                    context = mysql_user_mamager.new_host()
                elif action == 'delete_host':
                    context = mysql_user_mamager.delete_host()
                elif action == 'new_user':
                    context = mysql_user_mamager.new_host()

            return HttpResponse(json.dumps(context))

        else:
            error = form.errors.as_text()
            context = {'status': 2, 'msg': error}
            return HttpResponse(json.dumps(context))


class RBackupTaskView(View):
    @permission_required('can_scheduled_view')
    def get(self, request):
        return render(request, 'backup_task.html')


class BackupTaskDetailView(View):
    @permission_required('can_scheduled_view')
    def get(self, request):
        data = format_request(request)
        kwargs = PeriodicTask.objects.get(pk=data.get('id')).kwargs
        return HttpResponse(json.dumps(kwargs))


class BackupTaskView(View):
    @permission_required('can_scheduled_view')
    def get(self, request):
        data = PeriodicTask.objects.filter(description=u'数据库备份').values()
        result = []
        for i in data:
            crontab_value = CrontabSchedule.objects.get(id=i.get('crontab_id'))
            i['crontab_value'] = str(crontab_value)
            kwargs = json.loads(i['kwargs'])
            i['ssh_host'] = kwargs.get('ssh_host')
            result.append(i)

        return JsonResponse(result, safe=False)

    @permission_required('can_scheduled_edit')
    @transaction.atomic
    def post(self, request):
        data = format_request(request)
        form = BackupForm(data)
        if form.is_valid():
            context = form.is_save()
        else:
            error = form.errors.as_text()
            context = {'status': 2, 'msg': error}

        return HttpResponse(json.dumps(context))


class BackupTaskPreviewView(View):
    def get(self, request, id):
        kwargs = json.loads(PeriodicTask.objects.get(pk=id).kwargs)
        host = kwargs.get('ssh_host')
        return render(request, 'backup_task_preview.html', {'id': id, 'host': host})


class BackupTaskPreviewListView(View):
    """获取备份主机的备份目录列表"""

    def get(self, request):
        data = format_request(request)
        id = data.get('id')
        show_type = data.get('type')
        kwargs = json.loads(PeriodicTask.objects.get(pk=id).kwargs)
        backup_dir = os.path.join(kwargs.get('backup_dir'), show_type)

        cmd = f"du -sh {backup_dir}/* --time"

        paramiko_conn = ParamikoOutput(ssh_user=kwargs.get('ssh_user'),
                                       ssh_password=kwargs.get('ssh_password'),
                                       ssh_host=kwargs.get('ssh_host'),
                                       ssh_port=kwargs.get('ssh_port'))
        data = paramiko_conn.run(cmd)
        result = []
        for i in data:
            split_i = i.split('\t')
            file_size = split_i[0]
            file_time = split_i[1]
            file_name = split_i[2]
            result.append({'file_name': file_name, 'file_size': file_size, 'file_time': file_time})
        result.reverse()
        return HttpResponse(json.dumps(result))


class GetBackupDiskUsedView(View):
    """获取指定主机备份目录磁盘空间的使用详情"""

    def get(self, request):
        data = format_request(request)
        id = data.get('id')
        kwargs = json.loads(PeriodicTask.objects.get(pk=id).kwargs)
        backup_dir = kwargs.get('backup_dir')
        mysqldump_backup_dir = os.path.join(backup_dir, 'mysqldump')
        xtrabackup_backup_dir = os.path.join(backup_dir, 'xtrabackup')

        paramiko_conn = ParamikoOutput(ssh_user=kwargs.get('ssh_user'),
                                       ssh_password=kwargs.get('ssh_password'),
                                       ssh_host=kwargs.get('ssh_host'),
                                       ssh_port=kwargs.get('ssh_port'))

        cmd = f"du -sh {mysqldump_backup_dir} {xtrabackup_backup_dir} && df -h {backup_dir}"
        data = paramiko_conn.run(cmd)

        result = {}
        for i in data[:2]:
            result[i.split('\t')[1]] = i.split('\t')[0]

        df = [i for i in data[-1].split(' ') if i != '']
        result.update({'total_size': df[1],
                       'used_size': df[2],
                       'free_size': df[3],
                       'used_percent (%)': int(df[4].split('%')[0]),
                       'free_percent (%)': 100 - int(df[4].split('%')[0])
                       })

        return HttpResponse(json.dumps(result))

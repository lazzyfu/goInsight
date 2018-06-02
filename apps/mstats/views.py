import json
import os

import pymysql
from django.db import transaction
from django.http import HttpResponse, JsonResponse
from django.shortcuts import render
# Create your views here.
from django.utils.decorators import method_decorator
from django.views import View
# from djcelery.models import PeriodicTask, CrontabSchedule
from django_celery_beat.models import PeriodicTask, CrontabSchedule

from mstats.forms import PrivModifyForm, BackupTaskForm, SchemaMonitorForm
from mstats.utils import get_mysql_user_info, check_mysql_conn_status, MysqlUserManager, ParamikoOutput, MySQLQuery
from project_manager.models import InceptionHostConfig
from user_manager.permissions import permission_required
from utils.tools import format_request


class RenderMySQLUserView(View):
    @permission_required('can_mysql_user')
    def get(self, request):
        return render(request, 'mysql_user_manager.html')


class MySQLUserView(View):
    @permission_required('can_mysql_user')
    @method_decorator(check_mysql_conn_status)
    def get(self, request):
        data = format_request(request)
        host = data.get('host')
        data = get_mysql_user_info(host)

        return HttpResponse(json.dumps(data))


class MysqlUserManagerView(View):
    @permission_required('can_mysql_user')
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

            data = InceptionHostConfig.objects.get(comment=db_host)
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
                mysql_user_mamager = MysqlUserManager(locals())
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


class RSchemaMonitorTaskView(View):
    """渲染schema monitor页面"""

    @permission_required('can_scheduled')
    def get(self, request):
        return render(request, 'periodic_task.html')


class SchemaMonitorTaskView(View):
    """处理schema monitor数据"""

    @permission_required('can_scheduled')
    def get(self, request):
        data = PeriodicTask.objects.filter(description=u'表结构监控').values()
        result = []
        for i in data:
            crontab_value = CrontabSchedule.objects.get(id=i.get('crontab_id'))
            i['crontab_value'] = str(crontab_value)
            result.append(i)
            kwargs = json.loads(i['kwargs'])
            i['host'] = kwargs.get('host')
            i['schema'] = kwargs.get('schema')
            i['receiver'] = kwargs.get('receiver')

        return JsonResponse(result, safe=False)

    @permission_required('can_scheduled')
    @transaction.atomic
    def post(self, request):
        data = format_request(request)
        form = SchemaMonitorForm(data)
        if form.is_valid():
            context = form.is_save()
        else:
            error = form.errors.as_text()
            context = {'status': 2, 'msg': error}

        return HttpResponse(json.dumps(context))


class RBackupTaskView(View):
    """渲染backup task页面"""

    @permission_required('can_scheduled')
    def get(self, request):
        return render(request, 'backup_task.html')


class BackupTaskView(View):
    @permission_required('can_scheduled')
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

    @permission_required('can_scheduled')
    @transaction.atomic
    def post(self, request):
        data = format_request(request)
        form = BackupTaskForm(data)
        if form.is_valid():
            context = form.is_save()
        else:
            error = form.errors.as_text()
            context = {'status': 2, 'msg': error}

        return HttpResponse(json.dumps(context))


class BackupTaskDetailView(View):
    """获取备份任务详情"""

    @permission_required('can_scheduled')
    def get(self, request):
        data = format_request(request)
        kwargs = PeriodicTask.objects.get(pk=data.get('id')).kwargs
        return HttpResponse(json.dumps(kwargs))


class BackupTaskPreviewView(View):
    """渲染备份数据预览页面数据"""

    @permission_required('can_scheduled')
    def get(self, request, id):
        kwargs = json.loads(PeriodicTask.objects.get(pk=id).kwargs)
        host = kwargs.get('ssh_host')
        return render(request, 'backup_task_preview.html', {'id': id, 'host': host})


class BackupTaskPreviewListView(View):
    """获取备份主机的备份目录列表"""

    @permission_required('can_scheduled')
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
        if data['status'] == 0:
            result = []
            for i in data['data']:
                split_i = i.split('\t')
                file_size = split_i[0]
                file_time = split_i[1]
                file_name = split_i[2]
                result.append({'file_name': file_name, 'file_size': file_size, 'file_time': file_time})
            result.reverse()
            context = result
        else:
            context = []
        return HttpResponse(json.dumps(context))


class GetBackupDiskUsedView(View):
    """获取指定主机备份目录磁盘空间的使用详情"""

    @permission_required('can_scheduled')
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
        if data['status'] == 0:
            result = {}
            for i in data['data'][:2]:
                result[i.split('\t')[1]] = i.split('\t')[0]

            df = [i for i in data['data'][-1].split()]
            result.update({'total_size': df[-5],
                           'used_size': df[-4],
                           'free_size': df[-3],
                           'used_percent (%)': int(df[-2].split('%')[0]),
                           'free_percent (%)': 100 - int(df[-2].split('%')[0])
                           })
            context = {'status': 0, 'data': result}
        else:
            context = data
        return HttpResponse(json.dumps(context))


class RMySQLQueryView(View):
    @permission_required('can_mysql_query')
    def get(self, request):
        return render(request, 'sql_query.html')


class MySQLQueryView(View):
    @permission_required('can_mysql_query')
    def post(self, request):
        data = format_request(request)
        querys = data.get('contents')
        host = data.get('host')
        database = data.get('database')
        mysql_query = MySQLQuery(querys, host, database)
        result = mysql_query.query(request)
        return JsonResponse(result, safe=False)

import json

from celery import schedules
from django.db import transaction
from django.http import HttpResponse, JsonResponse
from django.shortcuts import render
from celery import current_app

# Create your views here.
from django.utils.decorators import method_decorator
from django.views import View
from djcelery.models import CrontabSchedule, PeriodicTask
from djcelery.schedulers import ModelEntry

from UserManager.permissions import check_dba_permission, permission_required
from scheduled_tasks.forms import PeriodicForm
from utils.tools import format_request


class RCrontabView(View):
    @method_decorator(check_dba_permission)
    def get(self, request):
        return render(request, 'crontab.html')


class CrontabView(View):
    @method_decorator(check_dba_permission)
    def get(self, request):
        data = CrontabSchedule.objects.values()
        return JsonResponse(list(data), safe=False)

    @transaction.atomic
    @method_decorator(check_dba_permission)
    def post(self, request):
        data = format_request(request)
        action = data['action']
        del data['csrfmiddlewaretoken']
        del data['action']

        if action == 'new_crontab':
            ndata = dict([(k, v.replace(' ', '')) for k, v in data.items()])
            crobj = schedules.crontab(**ndata)
            ModelEntry.to_model_schedule(crobj)
            context = {'status': 0, 'msg': '创建成功'}
        elif action == 'delete_crontab':
            id = data.get('id')
            for i in id.split(','):
                CrontabSchedule.objects.get(id=i).delete()
            context = {'status': 0, 'msg': '删除成功'}
        elif action == 'edit_crontab':
            # 删除无用的元素
            del data['0']
            CrontabSchedule.objects.filter(id=data.get('id')).update(**data)
            context = {'status': 0, 'msg': '修改成功'}
        else:
            context = {'status': 2, 'msg': '操作失败'}

        return HttpResponse(json.dumps(context))


class RPeriodicTaskView(View):
    @method_decorator(check_dba_permission)
    def get(self, request):
        return render(request, 'periodic_task.html')


class PeriodicTaskView(View):
    @method_decorator(check_dba_permission)
    def get(self, request):
        data = PeriodicTask.objects.values()
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

    @transaction.atomic
    @method_decorator(check_dba_permission)
    def post(self, request):
        data = format_request(request)
        form = PeriodicForm(data)
        if form.is_valid():
            context = form.is_save()
        else:
            error = form.errors.as_text()
            context = {'status': 2, 'msg': error}

        return HttpResponse(json.dumps(context))


class DeletePeriodicTaskView(View):
    @transaction.atomic
    @method_decorator(check_dba_permission)
    def post(self, request):
        data = format_request(request)
        id = data.get('id')
        for i in id.split(','):
            PeriodicTask.objects.get(pk=i).delete()
        context = {'status': 0, 'msg': '删除成功'}
        return HttpResponse(json.dumps(context))


class ModifyPeriodicTaskView(View):
    @transaction.atomic
    @method_decorator(check_dba_permission)
    def post(self, request):
        data = format_request(request)
        id = data.get('id')
        status = data.get('status')
        action = data.get('action')

        context = {}
        if action == 'modify_status':
            PeriodicTask.objects.filter(pk=id).update(enabled=status)
            context = {'status': 0, 'msg': '状态切换成功'}
        elif action == 'edit_periodic':
            del data['csrfmiddlewaretoken']
            del data['action']
            kwargs = json.loads(PeriodicTask.objects.get(pk=id).kwargs)
            kwargs['receiver'] = data.get('receiver')
            PeriodicTask.objects.filter(pk=id).update(name=data.get('name'), kwargs=json.dumps(kwargs))

            context = {'status': 0, 'msg': '修改成功'}

        return HttpResponse(json.dumps(context))


class GetCrontabView(View):
    @method_decorator(check_dba_permission)
    def get(self, request):
        result = []
        for i in CrontabSchedule.objects.all():
            result.append({'id': i.pk, 'crontab_value': str(i)})
        return HttpResponse(json.dumps(result))


class GetCeleryTasksView(View):
    """
    周期任务中的任务命名必须以:monitor开发
    """
    @method_decorator(check_dba_permission)
    def get(self, post):
        # 获取任务列表
        celery_app = current_app
        celery_app.loader.import_default_modules()
        tasks = list(sorted(name for name in celery_app.tasks if not name.startswith('celery.') if 'monitor' in name))
        return HttpResponse(json.dumps(tasks))

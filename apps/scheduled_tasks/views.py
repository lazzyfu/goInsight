import json

from celery import current_app
from celery import schedules
from django.db import transaction
from django.http import HttpResponse, JsonResponse
from django.shortcuts import render
# Create your views here.
from django.views import View
from django_celery_beat.models import CrontabSchedule, PeriodicTask, PeriodicTasks
from django_celery_beat.schedulers import ModelEntry

from scheduled_tasks.utils import refresh_periodic_tasks
from user_manager.permissions import permission_required
from utils.tools import format_request


class RCrontabView(View):
    @permission_required('can_scheduled_view')
    def get(self, request):
        return render(request, 'crontab.html')


class CrontabView(View):
    @permission_required('can_scheduled_view')
    def get(self, request):
        data = CrontabSchedule.objects.values()
        return JsonResponse(list(data), safe=False)

    @permission_required('can_scheduled_edit')
    @transaction.atomic
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
            # refresh_periodic_tasks()
            context = {'status': 0, 'msg': '删除成功'}
        elif action == 'edit_crontab':
            # 删除无用的元素
            del data['0']
            CrontabSchedule.objects.filter(id=data.get('id')).update(**data)
            # refresh_periodic_tasks()
            context = {'status': 0, 'msg': '修改成功'}
        else:
            context = {'status': 2, 'msg': '操作失败'}

        return HttpResponse(json.dumps(context))


class DeletePeriodicTaskView(View):
    @permission_required('can_scheduled_edit')
    @transaction.atomic
    def post(self, request):
        data = format_request(request)
        id = data.get('id')
        for i in id.split(','):
            PeriodicTask.objects.get(pk=i).delete()
        context = {'status': 0, 'msg': '删除成功'}
        return HttpResponse(json.dumps(context))


class ModifyPeriodicTaskView(View):
    @permission_required('can_scheduled_edit')
    @transaction.atomic
    def post(self, request):
        data = format_request(request)
        id = data.get('id')
        status = data.get('status')
        action = data.get('action')

        context = {}
        if action == 'modify_status':
            PeriodicTask.objects.filter(pk=id).update(enabled=status)
            refresh_periodic_tasks()
            context = {'status': 0, 'msg': '状态切换成功'}
        elif action == 'edit_periodic':
            del data['csrfmiddlewaretoken']
            del data['action']
            kwargs = json.loads(PeriodicTask.objects.get(pk=id).kwargs)
            kwargs['receiver'] = data.get('receiver')
            PeriodicTask.objects.filter(pk=id).update(name=data.get('name'), kwargs=json.dumps(kwargs))
            refresh_periodic_tasks()

            context = {'status': 0, 'msg': '修改成功'}

        return HttpResponse(json.dumps(context))


class GetCrontabView(View):
    @permission_required('can_scheduled_view')
    def get(self, request):
        result = []
        for i in CrontabSchedule.objects.all():
            result.append({'id': i.pk, 'crontab_value': str(i)})
        return HttpResponse(json.dumps(result))


class GetCeleryTasksView(View):
    @permission_required('can_scheduled_view')
    def get(self, post):
        # 获取任务列表
        celery_app = current_app
        celery_app.loader.import_default_modules()
        tasks = list(sorted(name for name in celery_app.tasks if not name.startswith('celery.')))
        return HttpResponse(json.dumps(tasks))

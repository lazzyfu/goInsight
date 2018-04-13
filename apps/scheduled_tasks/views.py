import json

from celery import schedules
from django.http import HttpResponse, JsonResponse
from django.shortcuts import render

# Create your views here.
from django.views import View
from djcelery.models import CrontabSchedule
from djcelery.schedulers import ModelEntry

from utils.tools import format_request


class RCrontabView(View):
    def get(self, request):
        return render(request, 'crontab.html')


class CrontabView(View):
    def get(self, request):
        data = CrontabSchedule.objects.values()
        return JsonResponse(list(data), safe=False)

    def post(self, request):
        data = format_request(request)
        action = data['action']
        del data['csrfmiddlewaretoken']
        del data['action']

        context = {}
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

        return HttpResponse(json.dumps(context))

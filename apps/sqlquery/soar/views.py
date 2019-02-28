# -*- coding:utf-8 -*-
# edit by fuzongfei
from django.http import JsonResponse
from django.shortcuts import render
from django.views import View

from sqlquery.soar.forms import SoarAnalyzeForm


class RenderSoarView(View):
    def get(self, request):
        return render(request, 'soar/soar.html')


class SoarAnalyzeView(View):
    def post(self, request):
        form = SoarAnalyzeForm(request.POST)
        if form.is_valid():
            context = form.analyze()
        else:
            msg = []
            for field, errors in form.errors.items():
                for error in errors:
                    msg.append(''.join([form.fields[field].label, error]))
            context = {'status': 2, 'msg': '\n'.join(msg)}
        return JsonResponse(context, safe=False)

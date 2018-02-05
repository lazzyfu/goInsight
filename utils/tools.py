# -*- coding:utf-8 -*-
# edit by fuzongfei


def format_request(request):
    data = {}
    if request.method == 'GET':
        for key in request.GET.keys():
            values_list = ','.join(request.GET.getlist(key))
            data[key] = values_list if len(values_list) > 1 else values_list
    elif request.method == 'POST':
        for key in request.POST.keys():
            values_list = ','.join(request.POST.getlist(key))
            data[key] = values_list if len(values_list) > 1 else values_list
    return data

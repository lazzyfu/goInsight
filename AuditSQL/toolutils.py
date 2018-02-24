# -*- coding:utf-8 -*-
# edit by fuzongfei

# 403页面
from django.shortcuts import render_to_response


def permission_denied(request):
    return render_to_response("403.html")

# 404页面
def page_not_found(request):
    return render_to_response("404.html")

# 500页面
def server_error(request):
    return render_to_response("500.html")
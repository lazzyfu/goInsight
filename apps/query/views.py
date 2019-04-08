# -*- coding:utf-8 -*-
# edit by fuzongfei
import json
import logging
import os
import time

import pymysql
from django.http import JsonResponse, HttpResponse
from django.shortcuts import render
# Create your views here.
from django.utils.crypto import get_random_string
from django.views import View
from django_redis import get_redis_connection

from config.config import DOMAIN_NAME
from query.forms import GetGrantSchemaForm, GetTableStrucForm, GetTableIndexForm, GetTableBaseForm, ExecSqlQueryForm, \
    GetHistorySqlForm, GetFilterHistorySqlForm, DbDictForm, GetTablesForm

logger = logging.getLogger('django')


class RenderSqlQueryView(View):
    """渲染SQL query页面"""

    def get(self, request):
        return render(request, 'query/sqlquery.html')


class GetGrantSchemaView(View):
    """获取指定环境授权给用户的schema信息"""

    def get(self, request):
        form = GetGrantSchemaForm(request.GET)
        context = None
        if form.is_valid():
            context = form.query(request)
        else:
            error = form.errors.as_text()
            logger.error(error)
        return JsonResponse(context, safe=False)


class GetTablesView(View):
    """获取指定主机的所有表"""

    def post(self, request):
        form = GetTablesForm(request.POST)
        if form.is_valid():
            context = form.query()
        return HttpResponse(json.dumps(context))


class GetTableStrucView(View):
    """返回表结构"""

    def get(self, request):
        form = GetTableStrucForm(request.GET)
        if form.is_valid():
            context = form.query()
        else:
            error = form.errors.as_text()
            context = {'status': 2, 'msg': error}

        return JsonResponse(context, safe=False)


class GetTableIndexView(View):
    """返回表索引"""

    def get(self, request):
        form = GetTableIndexForm(request.GET)
        if form.is_valid():
            context = form.query()
        else:
            error = form.errors.as_text()
            context = {'status': 2, 'msg': error}

        return JsonResponse(context, safe=False)


class GetTableBaseView(View):
    """返回表基本信息"""

    def get(self, request):
        form = GetTableBaseForm(request.GET)
        if form.is_valid():
            context = form.query()
        else:
            error = form.errors.as_text()
            context = {'status': 2, 'msg': error}

        return JsonResponse(context, safe=False)


class ExecSqlQueryView(View):
    """执行sql查询"""

    def post(self, request):
        form = ExecSqlQueryForm(request.POST)
        if form.is_valid():
            context = form.execute(request)
        else:
            error = form.errors.as_text()
            context = {'status': 2, 'msg': error}

        return JsonResponse(context, safe=False)


class GetHistorySqlView(View):
    """获取当前用户执行的SQL历史,返回前1000条"""

    def get(self, request):
        form = GetHistorySqlForm(request.GET)
        if form.is_valid():
            context = form.query(request)
        else:
            error = form.errors.as_text()
            context = {'status': 2, 'msg': error}

        return JsonResponse(context, safe=False)

    def post(self, request):
        form = GetFilterHistorySqlForm(request.POST)
        if form.is_valid():
            context = form.query(request)
        else:
            error = form.errors.as_text()
            context = {'status': 2, 'msg': error}

        return JsonResponse(context, safe=False)


class QueryStatusCheckView(View):
    """当用户刷新页面或者关闭页面时，检查是有正在执行的SQL，若有，kill掉"""

    def post(self, request):
        page_hash = request.POST.get('page_hash')
        cnx_redis = get_redis_connection('default')
        result = cnx_redis.smembers(page_hash)
        if result:
            for i in result:
                try:
                    thread_id, host, port, user, password = i.decode('utf-8').split('___')
                    conn = pymysql.connect(user=user,
                                           password=password,
                                           host=host,
                                           port=int(port),
                                           charset='utf8',
                                           cursorclass=pymysql.cursors.DictCursor)
                    with conn.cursor() as cursor:
                        cursor.execute(f'KILL CONNECTION {int(thread_id)}')
                        cursor.fetchone()
                        cursor.close()
                    logger.info(f'kill mysql query, thread_id:{thread_id}, user:{request.user.username}, host:{host}, \
                                port: {port}')
                    conn.close()
                    cnx_redis.srem(page_hash, i)
                except Exception as err:
                    logger.error(err)
        context = {'status': 0, 'msg': ''}
        return JsonResponse(context, safe=False)


class RenderDictView(View):
    def get(self, request):
        """渲染数据字典页面"""
        return render(request, 'query/mysql_dict.html')

    def post(self, request):
        form = DbDictForm(request.POST)
        if form.is_valid():
            data = form.query()
            if len(data) > 0:
                context = {'status': 0, 'data': data}
            else:
                context = {'status': 2, 'msg': '没有查询到数据'}
        else:
            error = form.errors.as_text()
            context = {'status': 2, 'msg': error}

        return JsonResponse(context, safe=False)


class GenerateHtmlView(View):
    def post(self, request):
        filename = '.'.join([get_random_string(24), 'html'])
        domain_name_tips = DOMAIN_NAME['value']
        data = request.POST.get('data')
        html_start = """
        <!DOCTYPE html>
        <html lang="en">
        <head>
            <meta charset="UTF-8">
            <title>数据字典</title>
            <style>
            html {
              font-family: sans-serif;
              -ms-text-size-adjust: 100%;
              -webkit-text-size-adjust: 100%;
            }
            body {
              margin: 0;
            }
            article,
            aside,
            details,
            figcaption,
            figure,
            footer,
            header,
            hgroup,
            main,
            menu,
            nav,
            section,
            summary {
              display: block;
            }
            audio,
            canvas,
            progress,
            video {
              display: inline-block;
              vertical-align: baseline;
            }
            audio:not([controls]) {
              display: none;
              height: 0;
            }
            [hidden],
            template {
              display: none;
            }
            a {
              background-color: transparent;
            }
            a:active,
            a:hover {
              outline: 0;
            }
            abbr[title] {
              border-bottom: 1px dotted;
            }
            b,
            strong {
              font-weight: bold;
            }
            dfn {
              font-style: italic;
            }
            h1 {
              font-size: 2em;
              margin: 0.67em 0;
            }
            mark {
              background: #ff0;
              color: #000;
            }
            small {
              font-size: 80%;
            }
            sub,
            sup {
              font-size: 75%;
              line-height: 0;
              position: relative;
              vertical-align: baseline;
            }
            sup {
              top: -0.5em;
            }
            sub {
              bottom: -0.25em;
            }
            img {
              border: 0;
            }
            svg:not(:root) {
              overflow: hidden;
            }
            figure {
              margin: 1em 40px;
            }
            hr {
              -webkit-box-sizing: content-box;
                 -moz-box-sizing: content-box;
                      box-sizing: content-box;
              height: 0;
            }
            pre {
              overflow: auto;
            }
            code,
            kbd,
            pre,
            samp {
              font-family: monospace, monospace;
              font-size: 1em;
            }
            button,
            input,
            optgroup,
            select,
            textarea {
              color: inherit;
              font: inherit;
              margin: 0;
            }
            button {
              overflow: visible;
            }
            button,
            select {
              text-transform: none;
            }
            button,
            html input[type="button"],
            input[type="reset"],
            input[type="submit"] {
              -webkit-appearance: button;
              cursor: pointer;
            }
            button[disabled],
            html input[disabled] {
              cursor: default;
            }
            button::-moz-focus-inner,
            input::-moz-focus-inner {
              border: 0;
              padding: 0;
            }
            input {
              line-height: normal;
            }
            input[type="checkbox"],
            input[type="radio"] {
              -webkit-box-sizing: border-box;
                 -moz-box-sizing: border-box;
                      box-sizing: border-box;
              padding: 0;
            }
            input[type="number"]::-webkit-inner-spin-button,
            input[type="number"]::-webkit-outer-spin-button {
              height: auto;
            }
            input[type="search"] {
              -webkit-appearance: textfield;
              -webkit-box-sizing: content-box;
                 -moz-box-sizing: content-box;
                      box-sizing: content-box;
            }
            input[type="search"]::-webkit-search-cancel-button,
            input[type="search"]::-webkit-search-decoration {
              -webkit-appearance: none;
            }
            fieldset {
              border: 1px solid #c0c0c0;
              margin: 0 2px;
              padding: 0.35em 0.625em 0.75em;
            }
            legend {
              border: 0;
              padding: 0;
            }
            textarea {
              overflow: auto;
            }
            optgroup {
              font-weight: bold;
            }
            table {
              border-collapse: collapse;
              border-spacing: 0;
            }
            td,
            th {
              padding: 0;
            }
            * {
              -webkit-box-sizing: border-box;
              -moz-box-sizing: border-box;
              box-sizing: border-box;
            }
            *:before,
            *:after {
              -webkit-box-sizing: border-box;
              -moz-box-sizing: border-box;
              box-sizing: border-box;
            }
            html {
              font-size: 10px;
              -webkit-tap-highlight-color: rgba(0, 0, 0, 0);
            }
            body {
              font-family: "Helvetica Neue", Helvetica, Arial, sans-serif;
              font-size: 14px;
              line-height: 1.42857143;
              color: #333333;
              background-color: #ffffff;
            }
            input,
            button,
            select,
            textarea {
              font-family: inherit;
              font-size: inherit;
              line-height: inherit;
            }
            a {
              color: #337ab7;
              text-decoration: none;
            }
            a:hover,
            a:focus {
              color: #23527c;
              text-decoration: underline;
            }
            a:focus {
              outline: 5px auto -webkit-focus-ring-color;
              outline-offset: -2px;
            }
            figure {
              margin: 0;
            }
            img {
              vertical-align: middle;
            }
            .img-responsive {
              display: block;
              max-width: 100%;
              height: auto;
            }
            .img-rounded {
              border-radius: 6px;
            }
            .img-thumbnail {
              padding: 4px;
              line-height: 1.42857143;
              background-color: #ffffff;
              border: 1px solid #dddddd;
              border-radius: 4px;
              -webkit-transition: all 0.2s ease-in-out;
              -o-transition: all 0.2s ease-in-out;
              transition: all 0.2s ease-in-out;
              display: inline-block;
              max-width: 100%;
              height: auto;
            }
            .img-circle {
              border-radius: 50%;
            }
            hr {
              margin-top: 20px;
              margin-bottom: 20px;
              border: 0;
              border-top: 1px solid #eeeeee;
            }
            .sr-only {
              position: absolute;
              width: 1px;
              height: 1px;
              margin: -1px;
              padding: 0;
              overflow: hidden;
              clip: rect(0, 0, 0, 0);
              border: 0;
            }
            .sr-only-focusable:active,
            .sr-only-focusable:focus {
              position: static;
              width: auto;
              height: auto;
              margin: 0;
              overflow: visible;
              clip: auto;
            }
            [role="button"] {
              cursor: pointer;
            }
            table {
              background-color: transparent;
            }
            caption {
              padding-top: 8px;
              padding-bottom: 8px;
              color: #777777;
              text-align: left;
            }
            th {
              text-align: left;
            }
            .table {
              width: 100%;
              max-width: 100%;
              margin-bottom: 20px;
            }
            .table > thead > tr > th,
            .table > tbody > tr > th,
            .table > tfoot > tr > th,
            .table > thead > tr > td,
            .table > tbody > tr > td,
            .table > tfoot > tr > td {
              padding: 8px;
              line-height: 1.42857143;
              vertical-align: top;
              border-top: 1px solid #dddddd;
            }
            .table > thead > tr > th {
              vertical-align: bottom;
              border-bottom: 2px solid #dddddd;
            }
            .table > caption + thead > tr:first-child > th,
            .table > colgroup + thead > tr:first-child > th,
            .table > thead:first-child > tr:first-child > th,
            .table > caption + thead > tr:first-child > td,
            .table > colgroup + thead > tr:first-child > td,
            .table > thead:first-child > tr:first-child > td {
              border-top: 0;
            }
            .table > tbody + tbody {
              border-top: 2px solid #dddddd;
            }
            .table .table {
              background-color: #ffffff;
            }
            .table-condensed > thead > tr > th,
            .table-condensed > tbody > tr > th,
            .table-condensed > tfoot > tr > th,
            .table-condensed > thead > tr > td,
            .table-condensed > tbody > tr > td,
            .table-condensed > tfoot > tr > td {
              padding: 5px;
            }
            .table-bordered {
              border: 1px solid #dddddd;
            }
            .table-bordered > thead > tr > th,
            .table-bordered > tbody > tr > th,
            .table-bordered > tfoot > tr > th,
            .table-bordered > thead > tr > td,
            .table-bordered > tbody > tr > td,
            .table-bordered > tfoot > tr > td {
              border: 1px solid #dddddd;
            }
            .table-bordered > thead > tr > th,
            .table-bordered > thead > tr > td {
              border-bottom-width: 2px;
            }
            .table-striped > tbody > tr:nth-of-type(odd) {
              background-color: #f9f9f9;
            }
            .table-hover > tbody > tr:hover {
              background-color: #f5f5f5;
            }
            table col[class*="col-"] {
              position: static;
              float: none;
              display: table-column;
            }
            table td[class*="col-"],
            table th[class*="col-"] {
              position: static;
              float: none;
              display: table-cell;
            }
            .table > thead > tr > td.active,
            .table > tbody > tr > td.active,
            .table > tfoot > tr > td.active,
            .table > thead > tr > th.active,
            .table > tbody > tr > th.active,
            .table > tfoot > tr > th.active,
            .table > thead > tr.active > td,
            .table > tbody > tr.active > td,
            .table > tfoot > tr.active > td,
            .table > thead > tr.active > th,
            .table > tbody > tr.active > th,
            .table > tfoot > tr.active > th {
              background-color: #f5f5f5;
            }
            .table-hover > tbody > tr > td.active:hover,
            .table-hover > tbody > tr > th.active:hover,
            .table-hover > tbody > tr.active:hover > td,
            .table-hover > tbody > tr:hover > .active,
            .table-hover > tbody > tr.active:hover > th {
              background-color: #e8e8e8;
            }
            .table > thead > tr > td.success,
            .table > tbody > tr > td.success,
            .table > tfoot > tr > td.success,
            .table > thead > tr > th.success,
            .table > tbody > tr > th.success,
            .table > tfoot > tr > th.success,
            .table > thead > tr.success > td,
            .table > tbody > tr.success > td,
            .table > tfoot > tr.success > td,
            .table > thead > tr.success > th,
            .table > tbody > tr.success > th,
            .table > tfoot > tr.success > th {
              background-color: #dff0d8;
            }
            .table-hover > tbody > tr > td.success:hover,
            .table-hover > tbody > tr > th.success:hover,
            .table-hover > tbody > tr.success:hover > td,
            .table-hover > tbody > tr:hover > .success,
            .table-hover > tbody > tr.success:hover > th {
              background-color: #d0e9c6;
            }
            .table > thead > tr > td.info,
            .table > tbody > tr > td.info,
            .table > tfoot > tr > td.info,
            .table > thead > tr > th.info,
            .table > tbody > tr > th.info,
            .table > tfoot > tr > th.info,
            .table > thead > tr.info > td,
            .table > tbody > tr.info > td,
            .table > tfoot > tr.info > td,
            .table > thead > tr.info > th,
            .table > tbody > tr.info > th,
            .table > tfoot > tr.info > th {
              background-color: #d9edf7;
            }
            .table-hover > tbody > tr > td.info:hover,
            .table-hover > tbody > tr > th.info:hover,
            .table-hover > tbody > tr.info:hover > td,
            .table-hover > tbody > tr:hover > .info,
            .table-hover > tbody > tr.info:hover > th {
              background-color: #c4e3f3;
            }
            .table > thead > tr > td.warning,
            .table > tbody > tr > td.warning,
            .table > tfoot > tr > td.warning,
            .table > thead > tr > th.warning,
            .table > tbody > tr > th.warning,
            .table > tfoot > tr > th.warning,
            .table > thead > tr.warning > td,
            .table > tbody > tr.warning > td,
            .table > tfoot > tr.warning > td,
            .table > thead > tr.warning > th,
            .table > tbody > tr.warning > th,
            .table > tfoot > tr.warning > th {
              background-color: #fcf8e3;
            }
            .table-hover > tbody > tr > td.warning:hover,
            .table-hover > tbody > tr > th.warning:hover,
            .table-hover > tbody > tr.warning:hover > td,
            .table-hover > tbody > tr:hover > .warning,
            .table-hover > tbody > tr.warning:hover > th {
              background-color: #faf2cc;
            }
            .table > thead > tr > td.danger,
            .table > tbody > tr > td.danger,
            .table > tfoot > tr > td.danger,
            .table > thead > tr > th.danger,
            .table > tbody > tr > th.danger,
            .table > tfoot > tr > th.danger,
            .table > thead > tr.danger > td,
            .table > tbody > tr.danger > td,
            .table > tfoot > tr.danger > td,
            .table > thead > tr.danger > th,
            .table > tbody > tr.danger > th,
            .table > tfoot > tr.danger > th {
              background-color: #f2dede;
            }
            .table-hover > tbody > tr > td.danger:hover,
            .table-hover > tbody > tr > th.danger:hover,
            .table-hover > tbody > tr.danger:hover > td,
            .table-hover > tbody > tr:hover > .danger,
            .table-hover > tbody > tr.danger:hover > th {
              background-color: #ebcccc;
            }
            .table-responsive {
              overflow-x: auto;
              min-height: 0.01%;
            }
            @media screen and (max-width: 767px) {
              .table-responsive {
                width: 100%;
                margin-bottom: 15px;
                overflow-y: hidden;
                -ms-overflow-style: -ms-autohiding-scrollbar;
                border: 1px solid #dddddd;
              }
              .table-responsive > .table {
                margin-bottom: 0;
              }
              .table-responsive > .table > thead > tr > th,
              .table-responsive > .table > tbody > tr > th,
              .table-responsive > .table > tfoot > tr > th,
              .table-responsive > .table > thead > tr > td,
              .table-responsive > .table > tbody > tr > td,
              .table-responsive > .table > tfoot > tr > td {
                white-space: nowrap;
              }
              .table-responsive > .table-bordered {
                border: 0;
              }
              .table-responsive > .table-bordered > thead > tr > th:first-child,
              .table-responsive > .table-bordered > tbody > tr > th:first-child,
              .table-responsive > .table-bordered > tfoot > tr > th:first-child,
              .table-responsive > .table-bordered > thead > tr > td:first-child,
              .table-responsive > .table-bordered > tbody > tr > td:first-child,
              .table-responsive > .table-bordered > tfoot > tr > td:first-child {
                border-left: 0;
              }
              .table-responsive > .table-bordered > thead > tr > th:last-child,
              .table-responsive > .table-bordered > tbody > tr > th:last-child,
              .table-responsive > .table-bordered > tfoot > tr > th:last-child,
              .table-responsive > .table-bordered > thead > tr > td:last-child,
              .table-responsive > .table-bordered > tbody > tr > td:last-child,
              .table-responsive > .table-bordered > tfoot > tr > td:last-child {
                border-right: 0;
              }
              .table-responsive > .table-bordered > tbody > tr:last-child > th,
              .table-responsive > .table-bordered > tfoot > tr:last-child > th,
              .table-responsive > .table-bordered > tbody > tr:last-child > td,
              .table-responsive > .table-bordered > tfoot > tr:last-child > td {
                border-bottom: 0;
              }
            }
            .clearfix:before,
            .clearfix:after {
              content: " ";
              display: table;
            }
            .clearfix:after {
              clear: both;
            }
            .center-block {
              display: block;
              margin-left: auto;
              margin-right: auto;
            }
            .pull-right {
              float: right !important;
            }
            .pull-left {
              float: left !important;
            }
            .hide {
              display: none !important;
            }
            .show {
              display: block !important;
            }
            .invisible {
              visibility: hidden;
            }
            .text-hide {
              font: 0/0 a;
              color: transparent;
              text-shadow: none;
              background-color: transparent;
              border: 0;
            }
            .hidden {
              display: none !important;
            }
            .affix {
              position: fixed;
            }
            </style>
        </head>
        <body>
        <div style='padding: 10px; margin: 10px'>
        """

        html_end = """
        </div>
        </body>
        </html>
        """
        # 如果tmp临时目录不存在，创建
        if not os.path.exists('media/tmp/'):
            os.makedirs('media/tmp/')
        # 自动删除2小时之前的临时文件
        for root, dirs, files in os.walk('media/tmp/'):
            for file in files:
                filename = os.path.join(root, file)
                file_ctime = os.path.getctime(filename)
                if int(time.time() - file_ctime) > 7200:
                    os.remove(filename)
        # 存储临时文件
        with open(f'media/tmp/{filename}', 'a', encoding='utf-8') as f:
            f.write(html_start)
            f.write(data)
            f.write(html_end)
        context = {'status': 0, 'jump_url': f'{domain_name_tips}/media/tmp/{filename}'}
        return JsonResponse(context, safe=False)

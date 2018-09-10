# -*- coding:utf-8 -*-
# edit by fuzongfei
import ast
import json
from ast import literal_eval

from channels.layers import get_channel_layer
from django.core.cache import cache
from django.db import transaction
from django.http import JsonResponse, HttpResponse
from django.shortcuts import render
from django.utils.decorators import method_decorator
from django.views import View

from apps.project_manager.inception.inception_api import GetBackupApi, IncepSqlCheck
from project_manager.models import IncepMakeExecTask, AuditContents
from project_manager.tasks import incep_async_tasks, \
    stop_incep_osc, get_osc_percent, incep_multi_tasks, update_audit_content_progress
from project_manager.utils import check_incep_alive
from user_manager.permissions import perform_tasks_permission_required
from utils.tools import format_request

channel_layer = get_channel_layer()


class PerformRecordsView(View):
    """渲染执行任务列表页"""

    def get(self, request):
        return render(request, 'perform_records.html')


class PerformRecordsListView(View):
    """渲染执行任务列表页表格数据"""

    def get(self, request):
        exec_tasks = []
        query = f"select a.id,a.user,a.taskid,a.dst_host,a.dst_database,a.make_time,a.envi_desc," \
                f"b.title,b.tasks " \
                f"from auditsql_incep_tasks as a left join auditsql_work_order as b " \
                f"on a.related_id = b.id where a.user='{request.user.username}' " \
                f"group by a.taskid order by a.make_time desc"

        for row in IncepMakeExecTask.objects.raw(query):
            exec_tasks.append({'user': row.user,
                               'taskid': row.taskid,
                               'title': row.title,
                               'tasks': row.tasks,
                               'envi_desc': row.envi_desc,
                               'dst_host': row.dst_host,
                               'dst_database': row.dst_database,
                               'make_time': row.make_time})
        return JsonResponse(list(exec_tasks), safe=False)


class PerformRecordsSQLPreView(View):
    """获取执行任务的SQL列表，进行预览展示"""
    def get(self, request):
        data = format_request(request)
        taskid = ast.literal_eval(data.get('taskid'))
        print(taskid)
        result = IncepMakeExecTask.objects.filter(taskid=taskid).values_list('sql_content', flat=True)
        return HttpResponse(json.dumps({'data': list(result)}))


class PerformResultsView(View):
    """返回执行任务执行结果和备份信息"""

    def get(self, request):
        id = request.GET.get('id')
        if IncepMakeExecTask.objects.get(id=id).exec_status in ('1', '4'):
            sql_detail = IncepMakeExecTask.objects.get(id=id)
            sequence_result = {'backupdbName': sql_detail.backup_dbname, 'sequence': sql_detail.sequence}
            rollback_sql = GetBackupApi(sequence_result).get_rollback_statement()

            exec_log = sql_detail.exec_log if sql_detail.exec_log else ''

            # 此处要将exec_log去字符串处理，否则无法转换为json
            data = {'rollback_log': rollback_sql, 'exec_log': literal_eval(exec_log)}
            context = {'status': 0, 'msg': '', 'data': data}
        else:
            context = {'status': 2, 'msg': '该SQL未被执行，无法查询状态信息'}

        return HttpResponse(json.dumps(context))


class PerformDetailsView(View):
    """渲染指定执行任务详情页面"""

    def get(self, request, taskid):
        return render(request, 'perform_details.html', {'taskid': taskid})


class PerformDetailsListView(View):
    """渲染指定执行任务页面数据"""

    def get(self, request):
        taskid = request.GET.get('taskid')

        query = f"select id,user,sqlsha1,sql_content,taskid,case exec_status " \
                f"when '0' then '未执行' when '1' then '已完成' when '2' then '处理中' when '3' then '回滚中' " \
                f"when '4' then '已回滚' " \
                f"when '6' then '异常' " \
                f"when '5' then '失败' end as exec_status " \
                f"from auditsql_incep_tasks where taskid={taskid}".format(taskid=taskid)
        i = 0
        task_details = []
        for row in IncepMakeExecTask.objects.raw(query):
            task_details.append({
                'sid': i,
                'id': row.id,
                'user': row.user,
                'envi_desc': row.envi_desc,
                'sqlsha1': row.sqlsha1,
                'sql_content': row.sql_content,
                'taskid': row.taskid,
                'exec_status': row.exec_status
            })
            i += 1
        return HttpResponse(json.dumps(task_details))


class PerformFullExecView(View):
    """执行任务-全部执行"""

    @method_decorator(check_incep_alive)
    @perform_tasks_permission_required('can_execute')
    @transaction.atomic
    def post(self, request):
        data = format_request(request)
        taskid = data.get('taskid')
        query = f"select * from auditsql_incep_tasks where taskid={taskid} order by id asc"

        key = ast.literal_eval(taskid)
        if 'run' == cache.get(key):
            context = {'status': 1, 'msg': '当前任务正在运行，请不要重复执行'}
        else:
            cache.set(key, 'run', timeout=600)
            incep_multi_tasks.delay(username=request.user.username,
                                    query=query,
                                    key=key)
            context = {'status': 1, 'msg': '任务已提交，请查看输出'}
        return HttpResponse(json.dumps(context))


class PerformExecView(View):
    """执行任务-手动单条执行"""

    @method_decorator(check_incep_alive)
    @perform_tasks_permission_required('can_execute')
    @transaction.atomic
    def post(self, request):
        data = format_request(request)
        id = data.get('id')
        obj = IncepMakeExecTask.objects.get(id=id)
        host = obj.dst_host
        port = obj.dst_port
        database = obj.dst_database
        sql = obj.sql_content + ';'

        key = ast.literal_eval(obj.taskid)
        if 'run' == cache.get(key):
            context = {'status': 1, 'msg': '正在自动化操作，请不要手动执行'}
        else:
            status = ''
            query = f"select id,group_concat(exec_status) as exec_status from auditsql_incep_tasks " \
                    f"where taskid={obj.taskid} group by taskid"
            for row in IncepMakeExecTask.objects.raw(query):
                status = row.exec_status.split(',')

            # 每次只能执行一条任务，不可同时执行，避免数据库压力
            if '2' in status or '3' in status:
                context = {'status': 2, 'msg': '请等待当前任务执行完成'}
            else:
                # 避免任务重复点击执行
                if obj.exec_status not in ('0', '5', '6'):
                    context = {'status': 2, 'msg': '请不要重复操作任务'}
                else:
                    # 将任务进度设置为：处理中
                    obj.exec_status = 2
                    obj.save()

                    # 如果sqlsha1存在，使用pt-online-schema-change执行
                    if obj.sqlsha1:
                        # 异步执行SQL任务
                        r = incep_async_tasks.delay(user=request.user.username,
                                                    id=id,
                                                    sql=sql,
                                                    host=host,
                                                    port=port,
                                                    database=database,
                                                    sqlsha1=obj.sqlsha1,
                                                    backup='yes',
                                                    exec_status='2')
                        task_id = r.task_id
                        # 将celery task_id写入到表
                        obj.celery_task_id = task_id
                        obj.save()
                        # 获取OSC执行进度
                        get_osc_percent.delay(task_id=task_id)

                        context = {'status': 1, 'msg': '任务已提交，请查看输出'}

                    else:
                        # 当affected_row>2000时，只执行不备份
                        if obj.affected_row > 2000:
                            incep_async_tasks.delay(user=request.user.username,
                                                    id=id,
                                                    sql=sql,
                                                    host=host,
                                                    port=port,
                                                    database=database,
                                                    exec_status='2')
                        else:
                            # 当affected_row<=2000时，执行并备份
                            incep_async_tasks.delay(user=request.user.username,
                                                    id=id,
                                                    backup='yes',
                                                    sql=sql,
                                                    host=host,
                                                    port=port,
                                                    database=database,
                                                    exec_status='2')

                        context = {'status': 1, 'msg': '任务已提交，请查看输出'}
            # 更新父任务进度
            update_audit_content_progress(request.user.username, obj.taskid)
        return HttpResponse(json.dumps(context))


class PerformStopView(View):
    """
    执行任务-停止OSC执行
    只支持停止修改表结构的操作
    """

    @method_decorator(check_incep_alive)
    @perform_tasks_permission_required('can_execute')
    @transaction.atomic
    def post(self, request):
        id = request.POST.get('id')
        obj = IncepMakeExecTask.objects.get(id=id)
        celery_task_id = obj.celery_task_id

        if obj.exec_status in ('0', '1', '4'):
            context = {'status': 2, 'msg': '请不要重复操作任务'}
        else:
            # 关闭正在执行的任务
            stop_incep_osc.delay(user=request.user.username,
                                 id=id,
                                 celery_task_id=celery_task_id)
            context = {'status': 1, 'msg': '任务已提交，请查看输出'}
        return HttpResponse(json.dumps(context))


class PerformRollbackView(View):
    """
    执行任务-回滚操作
    回滚操作不会进行再次进行备份
    """

    @method_decorator(check_incep_alive)
    @perform_tasks_permission_required('can_execute')
    @transaction.atomic
    def post(self, request):
        data = format_request(request)
        id = data.get('id')
        obj = IncepMakeExecTask.objects.get(id=id)
        host = obj.dst_host
        port = obj.dst_port
        database = obj.dst_database

        if obj.exec_status in ('0', '3', '4'):
            context = {'status': 2, 'msg': '请不要重复操作'}
        else:
            # 获取回滚语句
            rollback_sql = GetBackupApi(
                {'backupdbName': obj.backup_dbname, 'sequence': obj.sequence}).get_rollback_statement()
            if rollback_sql is None:
                context = {'status': 2, 'msg': '没有找到备份记录，回滚失败'}
            else:
                of_audit = IncepSqlCheck(rollback_sql, obj.dst_host, obj.dst_port, obj.dst_database,
                                         request.user.username)
                result = of_audit.make_sqlsha1()[1:]

                for row in result:
                    rollback_sql = row['SQL'] + ';'
                    rollback_sqlsha1 = row['sqlsha1']

                    # 将任务进度设置为：回滚中
                    obj.exec_status = 3
                    obj.rollback_sqlsha1 = rollback_sqlsha1
                    obj.save()

                    if row['sqlsha1']:
                        # 异步执行SQL任务
                        r = incep_async_tasks.delay(user=request.user.username,
                                                    id=id,
                                                    host=host,
                                                    port=port,
                                                    database=database,
                                                    sql=rollback_sql,
                                                    sqlsha1=rollback_sqlsha1,
                                                    exec_status='3')
                        task_id = r.task_id
                        # 将celery task_id写入到表
                        obj.celery_task_id = task_id
                        obj.save()
                        # 获取OSC执行进度
                        get_osc_percent.delay(task_id=task_id)

                        context = {'status': 1, 'msg': '任务已提交，请查看输出'}
                    else:
                        incep_async_tasks.delay(user=request.user.username,
                                                id=id,
                                                sql=rollback_sql,
                                                host=host,
                                                port=port,
                                                database=database,
                                                exec_status='3')

                        context = {'status': 1, 'msg': '任务已提交，请查看输出'}
        return HttpResponse(json.dumps(context))

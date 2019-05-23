# -*- coding:utf-8 -*-
# edit by fuzongfei
import os
import re
import subprocess
import uuid

import sqlparse
from django.core.cache import cache
from django.db.models import Case, When, Value, CharField
from django.utils import timezone
from django.utils.html import linebreaks
from rest_framework import serializers

from config.config import GHOST_TOOL
from orders.tasks import async_full_execute, async_execute_sql, async_export_tasks
from orders.models import Orders, OrdersTasks
from orders.utils.task import update_orders_progress
from orders.utils.tools import checkdbstatus
from orders.utils.validatorsCheck import envi_validator, subtask_stop_validator


class GenerateSubtasksSerializer(serializers.Serializer):
    id = serializers.CharField(required=True, error_messages={'required': 'id字段不能为空'})
    envi_id = serializers.CharField(required=True, validators=[envi_validator],
                                    error_messages={'required': 'envi_id字段不能为空'})

    def save(self, request):
        sdata = self.validated_data
        id = sdata.get('id')
        envi_id = sdata.get('envi_id')

        data = Orders.objects.get(pk=id)
        s, msg = checkdbstatus(host=data.host, port=data.port)
        if not s:
            return False, msg
        if s:
            if data.progress not in ('0', '1'):
                # 判断当前工单的子任务是否存在，存在直接跳转
                if OrdersTasks.objects.filter(order_id=id).exists():
                    taskid = OrdersTasks.objects.filter(order_id=id).first().taskid
                    jump_url = f'/orders/subtasks/list/{taskid}'
                    return True, jump_url
                else:
                    # 分割SQL，转换成sql列表
                    # 移除sql头尾的分号;
                    split_sqls = [sql.strip(';') for sql in sqlparse.split(data.contents, encoding='utf8')]
                    taskid = uuid.uuid1().hex

                    # 生成子任务
                    for sql in split_sqls:
                        OrdersTasks.objects.create(
                            applicant=data.applicant,
                            taskid=taskid,
                            host=data.host,
                            port=data.port,
                            database=data.database,
                            sql=sql.strip(';'),
                            sql_type=data.sql_type,
                            envi_id=envi_id,
                            file_format=data.file_format,
                            order_id=id
                        )
                    jump_url = f'/orders/subtasks/list/{taskid}'
                    return True, jump_url
            else:
                return False, '工单未完成，无法执行'


class SubtasksDetailSerializer(serializers.Serializer):
    taskid = serializers.CharField(required=True, min_length=32, max_length=32, error_messages={
        'required': 'taskid不能为空',
        'min_length': 'taskid长度必须为32个字符',
        'max_length': 'taskid长度必须为32个字符',
    })

    def query(self):
        sdata = self.validated_data
        taskid = sdata.get('taskid')
        queryset = OrdersTasks.objects.annotate(
            progress=Case(
                When(task_progress='0', then=Value('未执行')),
                When(task_progress='1', then=Value('已完成')),
                When(task_progress='2', then=Value('处理中')),
                When(task_progress='3', then=Value('失败')),
                When(task_progress='4', then=Value('异常')),
                output_field=CharField(),
            )
        ).filter(taskid=taskid).values('id', 'applicant', 'sql', 'taskid', 'progress', 'sql_type')

        i = 1
        task_details = []

        for row in queryset:
            task_details.append({
                'sid': i,
                'id': row['id'],
                'applicant': row['applicant'],
                'sql': row['sql'],
                'taskid': row['taskid'],
                'progress': row['progress'],
                'sql_type': row['sql_type']
            })
            i += 1
        return task_details


class FullExecuteSerializer(serializers.Serializer):
    taskid = serializers.CharField(required=True, min_length=32, max_length=32, error_messages={
        'required': 'taskid不能为空',
        'min_length': 'taskid长度必须为32个字符',
        'max_length': 'taskid长度必须为32个字符',
    })

    def execute(self, request):
        sdata = self.validated_data
        taskid = sdata.get('taskid')

        query = f"select * from auditsql_orders_tasks where taskid=\"{taskid}\" order by id asc"

        key = taskid
        if 'run' == cache.get(key):
            return False, '当前任务正在运行，请不要重复执行'
        else:
            cache.set(key, 'run', timeout=10)
            async_full_execute.delay(username=request.user.username,
                                     query=query,
                                     key=key)
            return True, '任务已提交到后台执行，请查看输出'


class SingleExecuteSerializer(serializers.Serializer):
    id = serializers.IntegerField(required=True, error_messages={'required': '任务ID不能为空'})

    def execute(self, request):
        sdata = self.validated_data
        id = sdata.get('id')

        obj = OrdersTasks.objects.get(id=id)
        host = obj.host
        port = obj.port
        database = obj.database
        sql = obj.sql

        key = obj.taskid
        if 'run' == cache.get(key):
            return False, '当前工单正在全部执行，请不要手动执行'
        else:
            status = ''
            query = f"select id,group_concat(task_progress) as task_progress from auditsql_orders_tasks " \
                f"where taskid=\"{obj.taskid}\" group by taskid"
            for row in OrdersTasks.objects.raw(query):
                status = row.task_progress.split(',')

            # 每次只能执行一条任务，不可同时执行，避免数据库压力
            if '2' in status:
                return False, '请等待当前任务执行完成'

            if obj.task_progress == '1':
                return False, '请不要重复执行任务'

            if obj.task_progress in ('0', '3', '4'):
                # 将任务进度设置为：处理中
                obj.executor = request.user.username
                obj.execition_time = timezone.now()
                obj.task_progress = '2'
                obj.save()

                # 获取工单的备注
                sql_type = Orders.objects.get(pk=obj.order_id).sql_type
                # 如果type为EXPORT
                if sql_type == 'EXPORT':
                    async_export_tasks.delay(username=request.user.username,
                                             id=id,
                                             sql=sql,
                                             host=obj.host,
                                             port=obj.port,
                                             database=obj.database)
                # 如果type为DML和DDL
                if sql_type in ['DML', 'DDL']:
                    async_execute_sql.delay(
                        username=request.user.username,
                        id=id,
                        sql=sql,
                        host=host,
                        port=port,
                        database=database,
                        task_progress='2')

                # 更新父任务进度
                update_orders_progress(request.user.username, obj.taskid)
                return True, '任务已提交到后台执行，请查看输出'


class StopExecuteSerializer(serializers.Serializer):
    id = serializers.IntegerField(required=True, error_messages={'required': '任务ID不能为空'})
    action = serializers.CharField(required=True, validators=[subtask_stop_validator], error_messages={
        'required': '停止动作值不能为空'
    })

    def op(self):
        sdata = self.validated_data
        id = sdata.get('id')
        action = sdata.get('action')

        obj = OrdersTasks.objects.get(id=id)
        if obj.sql_type != 'DDL':
            return False, '非DDL语句，无法执行'

        # 判断是否使用gh-ost执行
        if GHOST_TOOL['enable'] is True:
            # 获取gh-ost的sock文件
            # 将语句中的注释和SQL分离
            sql_split = {}
            for stmt in sqlparse.split(obj.sql):
                sql = sqlparse.parse(stmt)[0]
                sql_comment = sql.token_first()
                if isinstance(sql_comment, sqlparse.sql.Comment):
                    sql_split = {'comment': sql_comment.value, 'sql': sql.value.replace(sql_comment.value, '')}
                else:
                    sql_split = {'comment': '', 'sql': sql.value}

            # 获取不包含注释的SQL语句
            sql = sql_split['sql']
            formatsql = re.compile('^ALTER(\s+)TABLE(\s+)([\S]*)(\s+)(ADD|CHANGE|REMAME|MODIFY|DROP)([\s\S]*)',
                                   re.I)
            match = formatsql.match(sql)
            # 由于gh-ost不支持反引号，会被解析成命令，因此此处替换掉
            table = match.group(3).replace('`', '')
            # 将schema.table进行处理，这种情况gh-ost不识别，只保留table
            if len(table.split('.')) > 1:
                table = table.split('.')[1]
            sock = os.path.join('/tmp', f"gh-ost.{obj.database}.{table}.sock")
            # 判断sock是否存在
            if os.path.exists(sock):
                if action == 'pause_ghost':
                    pause_cmd = f"echo throttle | nc -U {sock}"
                    p = subprocess.Popen(pause_cmd, shell=True)
                    p.wait()
                    return True, '暂停动作已执行，请查看输出'

                if action == 'recovery_ghost':
                    recovery_cmd = f"echo no-throttle | nc -U {sock}"
                    p = subprocess.Popen(recovery_cmd, shell=True)
                    p.wait()
                    return True, '恢复动作已执行，请查看输出'

                if action == 'stop_ghost':
                    stop_cmd = f"echo panic | nc -U {sock}"
                    p = subprocess.Popen(stop_cmd, shell=True)
                    p.wait()
                    return True, '终止动作已执行，请查看输出'
            else:
                return False, f'不能找到文件{sock}, 操作失败'
        else:
            return False, '当前SQL未使用Gh-ost工具执行'


class GetTasksLogSerializer(serializers.Serializer):
    id = serializers.IntegerField(required=True, error_messages={'required': '任务ID不能为空'})

    def query(self):
        sdata = self.validated_data
        id = sdata.get('id')

        obj = OrdersTasks.objects.get(id=id)
        task_execlog = linebreaks(obj.task_execlog)

        if obj.task_progress in ('1', '3', '4'):
            data = obj.rollback_sql if obj.is_ghost == 0 else '无'
            return True, {'log': task_execlog, 'data': data}
        else:
            return False, '当前SQL未被执行，无法查询状态信息'

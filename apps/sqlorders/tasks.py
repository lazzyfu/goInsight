# -*- coding:utf-8 -*-
# edit by fuzongfei

"""
status = 0: 推送执行结果
status = 1: 推送执行进度
status = 2: 推送inception processlist
"""

import ast
import json
import logging
import os
import re
import subprocess
import time

import sqlparse
from asgiref.sync import async_to_sync
from celery import shared_task
from celery.result import AsyncResult
from channels.layers import get_channel_layer
from django.core.cache import cache
from django.utils import timezone

from sqlorders.inceptionApi import InceptionSqlApi
from sqlorders.models import SqlOrdersExecTasks, SqlOrdersContents, MysqlConfig, SysConfig
from sqlorders.msgNotice import SqlOrdersMsgPull

channel_layer = get_channel_layer()
logger = logging.getLogger('django')


def update_tasks_status(id=None, exec_result=None, exec_status=None):
    """
    更新任务进度
    更新备份信息
    """

    data = SqlOrdersExecTasks.objects.get(id=id)
    errlevel = [x['errlevel'] for x in exec_result] if exec_result is not None else []

    if exec_result is None:
        # 若inception没有返回结果，标记为异常
        data.exec_status = '6'
        data.save()
    else:
        # 执行失败
        if 1 in errlevel or 2 in errlevel:
            # 状态变为失败
            data.exec_status = '5'
            data.exec_log = exec_result
            data.save()
        else:
            # 执行成功
            # 执行状态为处理中时，状态变为已完成
            if exec_status == '2':
                data.exec_status = '1'
                data.sequence = exec_result[1]['sequence']
                data.backup_dbname = exec_result[1]['backup_dbname']
                data.exec_log = exec_result
                data.save()
            # 执行状态为回滚中时，状态变为已回滚
            elif exec_status == '3':
                data.exec_status = '4'
                data.save()


def update_audit_content_progress(username, taskid):
    # 检查任务是否都执行完成，如果执行完成，将父任务进度设置为已完成
    obj = SqlOrdersExecTasks.objects.filter(taskid=taskid)
    exec_status = obj.values_list('exec_status', flat=True)
    related_id = obj.first().related_id

    if related_id:
        if all([False for i in list(exec_status) if i != '1']):
            data = SqlOrdersContents.objects.get(id=related_id)
            if data.progress != '4':
                data.progress = '4'
                data.save()

                data.updated_at = timezone.now()
                data.save()
                # 发送邮件
                msg_pull = SqlOrdersMsgPull(id=related_id, user=username, type='feedback')
                msg_pull.run()


"""
status = 0: 推送执行结果
status = 1: 推送执行进度
status = 2: 推送inception processlist
"""


@shared_task
def get_osc_percent(task_id):
    """实时获取pt-online-schema-change执行进度"""
    task = AsyncResult(task_id)

    while task.state in ('PENDING', 'STARTED', 'PROGRESS'):
        while task.state == 'PROGRESS':
            user = task.result.get('user')
            host = task.result.get('host')
            port = task.result.get('port')
            database = task.result.get('database')
            sqlsha1 = task.result.get('sqlsha1')

            sql = f"inception get osc_percent '{sqlsha1}'"
            of_audit = InceptionSqlApi(host, port, database, sql, user)

            # 执行SQL
            of_audit.run_status(1)

            # 每1s获取一次
            time.sleep(1)
        else:
            continue


@shared_task(bind=True)
def incep_async_tasks(self, id=None, user=None, sql=None, sqlsha1=None, host=None, port=None, database=None,
                      exec_status=None,
                      backup=None):
    # 更新任务状态为: PROGRESS
    self.update_state(state="PROGRESS",
                      meta={'user': user, 'host': host, 'port': port, 'database': database, 'sqlsha1': sqlsha1})

    of_audit = InceptionSqlApi(host, port, database, sql, user)

    # 执行SQL
    exec_result = of_audit.run_exec(0, backup)

    # 更新任务进度
    update_tasks_status(id=id, exec_result=exec_result, exec_status=exec_status)

    # 更新任务状态为: SUCCESS
    self.update_state(state="SUCCESS")


@shared_task
def stop_incep_osc(user, id=None, celery_task_id=None):
    obj = SqlOrdersExecTasks.objects.get(id=id)
    host = obj.dst_host
    port = obj.dst_port
    database = obj.dst_database

    exec_status = None
    if obj.exec_status == '2':
        sqlsha1 = obj.sqlsha1
        exec_status = 0
    elif obj.exec_status == '3':
        sqlsha1 = obj.rollback_sqlsha1
        exec_status = 1

    sql = f"inception stop alter '{sqlsha1}'"

    # 执行SQL
    task = AsyncResult(celery_task_id)
    if task.state == 'PROGRESS':
        of_audit = InceptionSqlApi(host, port, database, sql, user)
        of_audit.run_status(0)

        # 更新任务进度
        update_tasks_status(id=id, exec_status=exec_status)


@shared_task
def incep_multi_tasks(username, query, key):
    taskid = key
    for row in SqlOrdersExecTasks.objects.raw(query):
        id = row.id
        host = row.host
        port = row.port
        database = row.database
        sqlsha1 = row.sqlsha1
        sql = row.sql + ';'

        obj = SqlOrdersExecTasks.objects.get(id=id)
        if obj.exec_status not in ('1', '2', '3', '4'):
            # 将任务进度设置为: 处理中
            obj.exec_status = '2'
            obj.save()

            # 如果sqlsha1存在，说明是大表，需要使用工具进行修改
            # inception_osc_min_table_size默认为16M
            # 如果此处向走gh-ost，则设置inception_osc_min_table_size=0
            if sqlsha1:
                # 判断是否使用gh-ost执行
                if SysConfig.objects.get(key='is_ghost').is_enabled == '0':
                    r = ghost_async_tasks.delay(user=username,
                                                id=id,
                                                sql=sql,
                                                host=obj.host,
                                                port=obj.port,
                                                database=obj.database)
                    task_id = r.task_id
                else:
                    # 异步执行SQL任务
                    r = incep_async_tasks.delay(user=username,
                                                id=id,
                                                sql=sql,
                                                host=host,
                                                port=port,
                                                database=database,
                                                sqlsha1=sqlsha1,
                                                backup='yes',
                                                exec_status='2')
                    task_id = r.task_id
                    # 将celery task_id写入到表
                    obj.celery_task_id = task_id
                    obj.save()
                    # 获取OSC执行进度
                    get_osc_percent.delay(task_id=task_id)
            else:
                # 当affected_row>1000000时，只执行不备份
                if obj.affected_row > 1000000:
                    r = incep_async_tasks.delay(user=username,
                                                id=id,
                                                sql=sql,
                                                host=host,
                                                port=port,
                                                database=database,
                                                exec_status='2')
                else:
                    # 当affected_row<=2000时，执行并备份
                    r = incep_async_tasks.delay(user=username,
                                                id=id,
                                                backup='yes',
                                                sql=sql,
                                                host=host,
                                                port=port,
                                                database=database,
                                                exec_status='2')
                task_id = r.task_id
            # 判断当前任务是否执行完成，执行完成后，执行下一个任务
            # 否则会变为并行异步执行
            task = AsyncResult(task_id)
            while task.state != 'SUCCESS':
                time.sleep(0.2)
                continue

    cache.delete(key)
    # 更新父任务进度
    update_audit_content_progress(username, ast.literal_eval(taskid))


@shared_task
def stop_incep_osc(user, id=None, celery_task_id=None):
    obj = SqlOrdersExecTasks.objects.get(id=id)
    host = obj.dst_host
    port = obj.dst_port
    database = obj.dst_database

    exec_status = None
    if obj.exec_status == '2':
        sqlsha1 = obj.sqlsha1
        exec_status = 0
    elif obj.exec_status == '3':
        sqlsha1 = obj.rollback_sqlsha1
        exec_status = 1

    sql = f"inception stop alter '{sqlsha1}'"

    # 执行SQL
    task = AsyncResult(celery_task_id)
    if task.state == 'PROGRESS':
        of_audit = InceptionSqlApi(host, port, database, sql, user)
        of_audit.run_status(0)

        # 更新任务进度
        update_tasks_status(id=id, exec_status=exec_status)


@shared_task
def ghost_async_tasks(user=None, id=None, sql=None, host=None, port=None, database=None):
    """ghost改表"""
    """ghost改表"""
    # 将语句中的注释和SQL分离
    sql_split = {}
    for stmt in sqlparse.split(sql):
        sql = sqlparse.parse(stmt)[0]
        sql_comment = sql.token_first()
        if isinstance(sql_comment, sqlparse.sql.Comment):
            sql_split = {'comment': sql_comment.value, 'sql': sql.value.replace(sql_comment.value, '')}
        else:
            sql_split = {'comment': '', 'sql': sql.value}

    # 获取不包含注释的SQL语句
    sql = sql_split['sql']

    formatsql = re.compile('^ALTER(\s+)TABLE(\s+)([\S]*)(\s+)(ADD|CHANGE|REMAME|MODIFY|DROP)([\s\S]*)', re.I)
    match = formatsql.match(sql)

    queryset = SqlOrdersExecTasks.objects.get(id=id)

    if match is not None:
        # 由于gh-ost不支持反引号，会被解析成命令，因此此处替换掉
        table = match.group(3).replace('`', '')
        # 将schema.table进行处理，这种情况gh-ost不识别，只保留table
        if len(table.split('.')) > 1:
            table = table.split('.')[1]

        # 处理反引号和将双引号处理成单引号
        value = ' '.join((match.group(5), match.group(6))).replace('`', '').replace('"', '\'')

        obj = MysqlConfig.objects.get(host=host, port=port)
        # 获取用户配置的gh-ost参数
        user_args = SysConfig.objects.get(key='is_ghost').value

        ghost_cmd = f"gh-ost {user_args} " \
                    f"--user={obj.user} --password=\"{obj.password}\" --host={host} --port={port} " \
                    f"--assume-master-host={host} " \
                    f"--database=\"{database}\" --table=\"{table}\" --alter=\"{value}\" --execute"

        # 删除sock，如果存在的话
        sock = os.path.join('/tmp', f"gh-ost.{database}.{table}.sock")
        os.remove(sock) if os.path.exists(sock) else None
        logger.info(f'删除sock：{sock}')

        # 执行gh-ost命令
        p = subprocess.Popen(ghost_cmd, shell=True, stdout=subprocess.PIPE, stderr=subprocess.STDOUT)
        logger.info(f'执行命令：{ghost_cmd}')

        # 记录ghost进程的pid
        queryset.ghost_pid = p.pid
        queryset.save()

        # 执行日志
        execute_log = ''

        # 检测子进程是否退出
        while p.poll() is None:
            data = p.stdout.readline().decode('utf8')
            if data:
                execute_log += data
                pull_msg = {'status': 3, 'data': data}
                # 推送消息
                async_to_sync(channel_layer.group_send)(user, {"type": "user.message",
                                                               'text': json.dumps(pull_msg)})

        if p.returncode == 0:
            # 返回状态为0，设置状态为成功
            queryset.exec_status = '1'
        else:
            # 返回状态不为0，设置状态为失败
            queryset.exec_status = '5'
        # 记录日志和标记为使用gh-ost执行
        queryset.exec_log = execute_log
        queryset.is_ghost = 1
        queryset.save()
    else:
        pull_msg = {'status': 3, 'data': f'未成功匹配到SQL：{sql}'}
        logger.error(f'无法匹配SQL：{sql}')
        # 推送消息
        async_to_sync(channel_layer.group_send)(user, {"type": "user.message",
                                                       'text': json.dumps(pull_msg)})

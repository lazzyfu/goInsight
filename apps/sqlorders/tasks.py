# -*- coding:utf-8 -*-
# edit by fuzongfei

"""
status = 0: 推送执行结果
status = 1: 推送执行进度
status = 2: 推送inception processlist
"""

import ast
import logging
import time

from celery import shared_task
from channels.layers import get_channel_layer
from django.core.cache import cache
from django.utils import timezone

from sqlorders.api.executeStatementApi import ExecuteSql
from sqlorders.models import SqlOrdersExecTasks, SqlOrdersContents, SysConfig, MysqlSchemas
from sqlorders.msgNotice import SqlOrdersMsgPull

channel_layer = get_channel_layer()
logger = logging.getLogger('django')


def upd_current_task_status(id=None, exec_result=None, exec_status=None):
    """更新当前任务的进度"""
    # exec_result的数据格式
    # {'status': 'success', 'rollbacksql': [sql,], 'affected_rows': 1, 'runtime': '1.000s', 'exec_log': ''}
    # 或
    # {'status': 'fail', 'exec_log': ''}
    data = SqlOrdersExecTasks.objects.get(id=id)
    if exec_result['status'] in ['fail', 'warn']:
        # 标记为失败
        data.exec_status = '5'
        data.exec_log = exec_result.get('exec_log')
        data.save()
    elif exec_result['status'] == 'success':
        print(exec_result)
        # 执行状态为处理中时，状态变为已完成
        if exec_status == '2':
            data.exec_status = '1'
            data.affected_row = int(exec_result.get('affected_rows'))
            data.rollback_sql = '\n'.join(exec_result.get('rollbacksql'))
            data.runtime = exec_result.get('runtime')
            data.exec_log = exec_result.get('exec_log')
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


@shared_task
def async_execute_sql(id=None, username=None, sql=None, host=None, port=None, database=None, exec_status=None):
    """执行SQL"""
    dst_server = MysqlSchemas.objects.get(host=host, port=port, schema=database)
    dst_host = dst_server.host
    dst_user = dst_server.user
    dst_password = dst_server.password
    dst_port = dst_server.port
    dst_database = database

    execute_sql = ExecuteSql(host=dst_host, port=dst_port,
                             user=dst_user, password=dst_password,
                             database=dst_database, username=username)
    result = execute_sql.run_by_sql(sql)

    # 更新任务进度
    upd_current_task_status(id=id, exec_result=result, exec_status=exec_status)


@shared_task
def async_execute_multi_sql(username, query, key):
    """批量执行SQL"""
    taskid = key
    for row in SqlOrdersExecTasks.objects.raw(query):
        id = row.id
        host = row.host
        port = row.port
        database = row.database
        sql = row.sql + ';'

        obj = SqlOrdersExecTasks.objects.get(id=id)
        if obj.exec_status not in ('1', '2'):
            # 将任务进度设置为: 处理中
            obj.executor = username
            obj.execition_time = timezone.now()
            obj.exec_status = '2'
            obj.save()

            # 执行SQL
            async_execute_sql.delay(
                username=username,
                id=id,
                sql=sql,
                host=host,
                port=port,
                database=database,
                exec_status='2')
        time.sleep(0.02)
    cache.delete(key)
    # 更新父任务进度
    update_audit_content_progress(username, ast.literal_eval(taskid))

# @shared_task
# def ghost_async_tasks(user=None, id=None, sql=None, host=None, port=None, database=None):
#     """ghost改表"""
#     # 将语句中的注释和SQL分离
#     sql_split = {}
#     for stmt in sqlparse.split(sql):
#         sql = sqlparse.parse(stmt)[0]
#         sql_comment = sql.token_first()
#         if isinstance(sql_comment, sqlparse.sql.Comment):
#             sql_split = {'comment': sql_comment.value, 'sql': sql.value.replace(sql_comment.value, '')}
#         else:
#             sql_split = {'comment': '', 'sql': sql.value}
#
#     # 获取不包含注释的SQL语句
#     sql = sql_split['sql']
#
#     formatsql = re.compile('^ALTER(\s+)TABLE(\s+)([\S]*)(\s+)(ADD|CHANGE|REMAME|MODIFY|DROP)([\s\S]*)', re.I)
#     match = formatsql.match(sql)
#
#     queryset = SqlOrdersExecTasks.objects.get(id=id)
#
#     if match is not None:
#         # 由于gh-ost不支持反引号，会被解析成命令，因此此处替换掉
#         table = match.group(3).replace('`', '')
#         # 将schema.table进行处理，这种情况gh-ost不识别，只保留table
#         if len(table.split('.')) > 1:
#             table = table.split('.')[1]
#
#         # 处理反引号和将双引号处理成单引号
#         value = ' '.join((match.group(5), match.group(6))).replace('`', '').replace('"', '\'')
#
#         obj = MysqlConfig.objects.get(host=host, port=port)
#         # 获取用户配置的gh-ost参数
#         user_args = SysConfig.objects.get(key='is_ghost').value
#
#         ghost_cmd = f"gh-ost {user_args} " \
#                     f"--user={obj.user} --password=\"{obj.password}\" --host={host} --port={port} " \
#                     f"--assume-master-host={host} " \
#                     f"--database=\"{database}\" --table=\"{table}\" --alter=\"{value}\" --execute"
#
#         # 删除sock，如果存在的话
#         sock = os.path.join('/tmp', f"gh-ost.{database}.{table}.sock")
#         os.remove(sock) if os.path.exists(sock) else None
#         logger.info(f'删除sock：{sock}')
#
#         # 执行gh-ost命令
#         p = subprocess.Popen(ghost_cmd, shell=True, stdout=subprocess.PIPE, stderr=subprocess.STDOUT)
#         logger.info(f'执行命令：{ghost_cmd}')
#
#         # 记录ghost进程的pid
#         queryset.ghost_pid = p.pid
#         queryset.save()
#
#         # 执行日志
#         execute_log = ''
#
#         # 检测子进程是否退出
#         while p.poll() is None:
#             data = p.stdout.readline().decode('utf8')
#             if data:
#                 execute_log += data
#                 pull_msg = {'status': 3, 'data': data}
#                 # 推送消息
#                 async_to_sync(channel_layer.group_send)(user, {"type": "user.message",
#                                                                'text': json.dumps(pull_msg)})
#
#         if p.returncode == 0:
#             # 返回状态为0，设置状态为成功
#             queryset.exec_status = '1'
#         else:
#             # 返回状态不为0，设置状态为失败
#             queryset.exec_status = '5'
#         # 记录日志和标记为使用gh-ost执行
#         queryset.exec_log = execute_log
#         queryset.is_ghost = 1
#         queryset.save()
#     else:
#         pull_msg = {'status': 3, 'data': f'未成功匹配到SQL：{sql}'}
#         logger.error(f'无法匹配SQL：{sql}')
#         # 推送消息
#         async_to_sync(channel_layer.group_send)(user, {"type": "user.message",
#                                                        'text': json.dumps(pull_msg)})


# -*- coding:utf-8 -*-
# edit by fuzongfei
import ast
import time

from celery import shared_task
from celery.result import AsyncResult
from channels.layers import get_channel_layer
from dingtalkchatbot.chatbot import DingtalkChatbot
from django.core.cache import cache
from django.utils import timezone

from project_manager.inception.inception_api import IncepSqlCheck
from project_manager.models import AuditContents, IncepMakeExecTask, \
    DomainName, Webhook
from user_manager.models import UserAccount

channel_layer = get_channel_layer()


def update_tasks_status(id=None, exec_result=None, exec_status=None):
    """
    更新任务进度
    更新备份信息
    """

    data = IncepMakeExecTask.objects.get(id=id)
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
    obj = IncepMakeExecTask.objects.filter(taskid=taskid)
    exec_status = obj.values_list('exec_status', flat=True)
    related_id = obj.first().related_id

    if related_id:
        if all([False for i in list(exec_status) if i != '1']):
            data = AuditContents.objects.get(id=related_id)
            if data.progress != '4':
                data.progress = '4'
                data.save()

                data.updated_at = timezone.now()
                data.save()
                # 发送钉钉推送（包括任务进度）
                xiaoding_pull.delay(user=username, id=related_id, type='feedback')


@shared_task
def ding_notice_pull(data, weekday):
    """at未完成工单的开发，告知其尽快处理工单"""
    if Webhook.objects.filter().first():
        webhook_addr = Webhook.objects.get().webhook_addr
        xiaoding = DingtalkChatbot(webhook_addr)

        text = ''
        mobile = []
        tasks = data[0].get('tasks') + '\n'
        for i in data:
            proposer = i.get('proposer')
            title = i.get('title')
            text += f'工单负责人：{proposer}\t\t标题：{title}\n'
            mobile.append(i.get('mobile'))

        if weekday == 2:
            title = "【预发布环境未处理工单提醒】\n请下面相关工单负责人及时处理、以免影响发布，◕‿◕\n"
        if weekday == 3:
            title = "【生产环境未处理工单提醒】\n请下面相关工单负责人及时处理、以免影响发布，◕‿◕\n"

        content = ''.join((title, tasks, text))

        xiaoding.send_text(msg=content, at_mobiles=mobile)


@shared_task
def xiaoding_pull(user, id, type, addition_info=None, sleep=3):
    """
    id: 记录的id
    type: commit、close、approve、feedback、hook
    addition_info: 用户附加的消息
    """
    # 休眠3s后执行
    time.sleep(sleep)

    # 获取域名
    domain_name = DomainName.objects.get().domain_name if DomainName.objects.filter().first() else None

    if Webhook.objects.filter().first():
        webhook_addr = Webhook.objects.get().webhook_addr
        xiaoding = DingtalkChatbot(webhook_addr)

        data = AuditContents.objects.get(pk=id)
        tasks = data.tasks if data.tasks else None
        jira_url = data.url if data.url else None
        host = ':'.join((data.host, str(data.port)))
        addition_info = addition_info if addition_info else None
        if data.envi_desc == 0:
            envi_desc = '测试环境'
            if domain_name:
                jump_url = '/'.join((domain_name, 'projects/ol/test_records/'))
        elif data.envi_desc == 1:
            envi_desc = '预发布环境'
            if domain_name:
                jump_url = '/'.join((domain_name, 'projects/ol/staging_records/'))
        elif data.envi_desc == 2:
            envi_desc = '生产环境'
            if domain_name:
                jump_url = '/'.join((domain_name, 'projects/ol/ol_records/'))

        # 如果用户手机号存在，钉钉直接@mobile
        # 如果手机号不存在，钉钉直接@all
        proposer_mobile = UserAccount.objects.get(username=data.proposer).mobile
        operator_mobile = UserAccount.objects.get(username=data.operator).mobile

        # 提交
        if type == 'commit':
            text = f"您好、{user}提交了审核内容，◕‿◕\n" \
                   f"标题: {data.title}\n" \
                   f"环境: {envi_desc}\n" \
                   f"类型: {data.operate_type}\n" \
                   f"主机: {host}\n" \
                   f"库名: {data.database}\n" \
                   f"审核DBA: {data.operator}\n" \
                   f"部署版本: {tasks}\n" \
                   f"JIRA需求链接: {jira_url}\n" \
                   f"URL: {jump_url} \n" \
                   f"提交时间: {timezone.localtime(data.created_at).strftime('%Y-%m-%d %H:%M:%S')}\n"

            if operator_mobile:
                xiaoding.send_text(msg=text, at_mobiles=[operator_mobile])
            else:
                xiaoding.send_text(msg=text, is_at_all=True)
        # 审核
        elif type == 'approve':
            if data.progress == '2':
                text = f"您好、{user}审核已通过，◕‿◕\n" \
                       f"标题: {data.title}\n" \
                       f"环境: {envi_desc}\n" \
                       f"附加信息: {addition_info}\n" \
                       f"URL: {jump_url} \n" \
                       f"审核时间: {timezone.localtime(data.operate_time).strftime('%Y-%m-%d %H:%M:%S')}\n"
            elif data.progress == '1':
                text = f"您好、{user}审核未通过，◕﹏◕\n" \
                       f"标题: {data.title}\n" \
                       f"环境: {envi_desc}\n" \
                       f"附加信息: {addition_info}\n" \
                       f"URL: {jump_url} \n" \
                       f"审核时间: {timezone.localtime(data.operate_time).strftime('%Y-%m-%d %H:%M:%S')}\n"
            if proposer_mobile:
                xiaoding.send_text(msg=text, at_mobiles=[proposer_mobile])
            else:
                xiaoding.send_text(msg=text, is_at_all=True)
        # 反馈
        elif type == 'feedback':
            if data.progress == '3':
                text = f"您好、{user}正在处理中，请稍后，◕‿◕\n" \
                       f"标题: {data.title}\n" \
                       f"环境: {envi_desc}\n" \
                       f"附加信息: {addition_info}\n" \
                       f"URL: {jump_url} \n" \
                       f"处理时间: {timezone.localtime(data.updated_at).strftime('%Y-%m-%d %H:%M:%S')}\n"
            elif data.progress == '4':
                text = f"您好、{user}处理完成，◕‿◕\n" \
                       f"标题: {data.title}\n" \
                       f"环境: {envi_desc}\n" \
                       f"附加信息: {addition_info}\n" \
                       f"URL: {jump_url} \n" \
                       f"完成时间: {timezone.localtime(data.updated_at).strftime('%Y-%m-%d %H:%M:%S')}\n"
            if proposer_mobile and operator_mobile:
                xiaoding.send_text(msg=text, at_mobiles=[proposer_mobile, operator_mobile])
            else:
                xiaoding.send_text(msg=text, is_at_all=True)
        # 关闭
        elif type == 'close':
            if data.progress == '5':
                text = f"您好、{user}关闭了记录，请不要处理，◕‿◕\n" \
                       f"标题: {data.title}\n" \
                       f"环境: {envi_desc}\n" \
                       f"附加信息: {addition_info}\n" \
                       f"URL: {jump_url} \n" \
                       f"关闭时间: {timezone.localtime(data.close_time).strftime('%Y-%m-%d %H:%M:%S')}\n"
            if proposer_mobile and operator_mobile:
                xiaoding.send_text(msg=text, at_mobiles=[proposer_mobile, operator_mobile])
            else:
                xiaoding.send_text(msg=text, is_at_all=True)
        # 回复
        elif type == 'reply':
            text = f"您好、{user}回复了内容，◕‿◕\n" \
                   f"标题: {data.title}\n" \
                   f"环境: {envi_desc}\n" \
                   f"回复内容: {addition_info}\n" \
                   f"URL: {jump_url} \n" \
                   f"回复时间: {timezone.localtime(data.updated_at).strftime('%Y-%m-%d %H:%M:%S')}\n"
            if proposer_mobile and operator_mobile:
                xiaoding.send_text(msg=text, at_mobiles=[proposer_mobile, operator_mobile])
            else:
                xiaoding.send_text(msg=text, is_at_all=True)

        # 钩子
        elif type == 'hook':
            text = f"您好、{user}扭转了工单，◕‿◕\n" \
                   f"标题: {data.title}\n" \
                   f"新环境: {envi_desc}\n" \
                   f"类型: {data.operate_type}\n" \
                   f"主机: {host}\n" \
                   f"库名: {data.database}\n" \
                   f"审核DBA: {data.operator}\n" \
                   f"部署版本: {tasks}\n" \
                   f"JIRA需求链接: {jira_url}\n" \
                   f"URL: {jump_url} \n" \
                   f"扭转时间: {timezone.localtime(data.updated_at).strftime('%Y-%m-%d %H:%M:%S')}\n"
            if proposer_mobile and operator_mobile:
                xiaoding.send_text(msg=text, at_mobiles=[proposer_mobile, operator_mobile])
            else:
                xiaoding.send_text(msg=text, is_at_all=True)


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
            of_audit = IncepSqlCheck(sql, host, port, database, user)

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

    of_audit = IncepSqlCheck(sql, host, port, database, user)

    # 执行SQL
    exec_result = of_audit.run_exec(0, backup)

    # 更新任务进度
    update_tasks_status(id=id, exec_result=exec_result, exec_status=exec_status)

    # 更新任务状态为: SUCCESS
    self.update_state(state="SUCCESS")


@shared_task
def stop_incep_osc(user, id=None, celery_task_id=None):
    obj = IncepMakeExecTask.objects.get(id=id)
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
        of_audit = IncepSqlCheck(sql, host, port, database, user)
        of_audit.run_status(0)

        # 更新任务进度
        update_tasks_status(id=id, exec_status=exec_status)


@shared_task
def incep_multi_tasks(username, query, key):
    taskid = key
    for row in IncepMakeExecTask.objects.raw(query):
        id = row.id
        host = row.dst_host
        port = row.dst_port
        database = row.dst_database
        sqlsha1 = row.sqlsha1
        sql = row.sql_content + ';'

        obj = IncepMakeExecTask.objects.get(id=id)
        if obj.exec_status not in ('1', '2', '3', '4'):
            # 将任务进度设置为: 处理中
            obj.exec_status = '2'
            obj.save()

            # 如果sqlsha1存在，使用pt-online-schema-change执行
            if sqlsha1:
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
                # 当affected_row>2000时，只执行不备份
                if obj.affected_row > 2000:
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

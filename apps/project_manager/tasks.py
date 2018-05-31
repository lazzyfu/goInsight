# -*- coding:utf-8 -*-
# edit by fuzongfei

import time

from celery import shared_task
from celery.result import AsyncResult
from channels.layers import get_channel_layer
from django.core.mail import EmailMessage
from django.db.models import F
from django.template.loader import render_to_string

from AuditSQL.settings import EMAIL_FROM
from project_manager.inception.inception_api import IncepSqlCheck
from project_manager.models import AuditContents, OlAuditContentsReply, IncepMakeExecTask, \
    DomainName, OlAuditDetail, Webhook
from project_manager.utils import update_tasks_status
from user_manager.utils import GetEmailAddr
from dingtalkchatbot.chatbot import DingtalkChatbot

channel_layer = get_channel_layer()


@shared_task
def xiaoding_pull(user, title, type, progress=None):
    """
    type: commit、 close、approve、feedback
    """
    if Webhook.objects.filter().first():
        webhook_addr = Webhook.objects.get().webhook_addr
        xiaoding = DingtalkChatbot(webhook_addr)

        if type == 'commit':
            xiaoding.send_text(msg=f"您好、{user}提交了审核内容\n标题：{title}")
        elif type == 'approve':
            if progress == '2':
                xiaoding.send_text(msg=f"您好、{user}已批准，请DBA处理\n标题：{title}")
            elif progress == '1':
                xiaoding.send_text(msg=f"您好、{user}审核未通过\n标题：{title}")
        elif type == 'feedback':
            if progress == '3':
                xiaoding.send_text(msg=f"您好、{user}正在处理中，请稍后\n标题：{title}")
            elif progress == '4':
                xiaoding.send_text(msg=f"您好、{user}处理完成\n标题：{title}")
        elif type == 'close':
            if progress == '5':
                xiaoding.send_text(msg=f"您好、{user}关闭了记录，请DBA不要处理\n标题：{title}")


@shared_task
def send_commit_mail(**kwargs):
    latest_id = kwargs['latest_id']
    userinfo = GetEmailAddr(AuditContents, latest_id)

    receiver = userinfo.get_user_email('proposer', 'verifier', 'operator')
    cc = userinfo.get_contact_email()
    bcc = userinfo.get_bcc_email()

    # 向_commit_mail.html渲染data数据
    domain_name = ''
    if DomainName.objects.filter().first():
        domain_name = DomainName.objects.get().domain_name
    data = AuditContents.objects.annotate(group_name=F('group__group_name')).get(pk=latest_id)
    detail = OlAuditDetail.objects.get(ol=latest_id)
    email_html_body = render_to_string('_send_commit_mail.html', {
        'data': data,
        'detail': detail,
        'domain_name': domain_name
    })

    # 发送邮件
    msg = EmailMessage(subject=data.title,
                       body=email_html_body,
                       from_email=EMAIL_FROM,
                       to=receiver,
                       cc=cc,
                       bcc=bcc,
                       )
    msg.content_subtype = "html"

    msg.send()


@shared_task
def send_verify_mail(**kwargs):
    latest_id = kwargs['latest_id']
    userinfo = GetEmailAddr(AuditContents, latest_id)

    receiver = userinfo.get_user_email('proposer', 'verifier', 'operator')
    cc = userinfo.get_contact_email()
    bcc = userinfo.get_bcc_email()

    # 向mail_template.html渲染data数据
    data = AuditContents.objects.get(pk=latest_id)
    email_html_body = render_to_string('_send_verify_mail.html', {
        'data': data,
        'type': kwargs.get('type'),
        'user_role': kwargs.get('user_role'),
        'username': kwargs.get('username'),
        'addition_info': kwargs.get('addition_info')
    })

    # 发送邮件
    headers = {'Reply: ': receiver}
    title = 'Re: ' + data.title
    msg = EmailMessage(subject=title,
                       body=email_html_body,
                       from_email=EMAIL_FROM,
                       to=receiver,
                       cc=cc,
                       bcc=bcc,
                       headers=headers)
    msg.content_subtype = "html"
    msg.send()


@shared_task
def send_reply_mail(**kwargs):
    latest_id = kwargs['latest_id']
    reply_id = kwargs['reply_id']
    userinfo = GetEmailAddr(AuditContents, latest_id)

    receiver = userinfo.get_user_email('proposer', 'verifier', 'operator')
    cc = userinfo.get_contact_email()
    bcc = userinfo.get_bcc_email()

    title = AuditContents.objects.get(pk=reply_id).title
    reply_record = OlAuditContentsReply.objects.get(pk=latest_id)

    # 向mail_template.html渲染data数据
    email_html_body = render_to_string('_send_reply_mail.html', {
        'reply_record': reply_record,
        'username': kwargs['username'],
    })

    # 发送邮件
    headers = {'Reply: ': receiver}
    title = 'Re: ' + title
    msg = EmailMessage(subject=title,
                       body=email_html_body,
                       from_email=EMAIL_FROM,
                       to=receiver,
                       cc=cc,
                       bcc=bcc,
                       headers=headers)
    msg.content_subtype = "html"
    msg.send()


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
            database = task.result.get('database')
            sqlsha1 = task.result.get('sqlsha1')

            sql = f"inception get osc_percent '{sqlsha1}'"
            of_audit = IncepSqlCheck(sql, host, database, user)

            # 执行SQL
            of_audit.run_status(1)

            # 每1s获取一次
            time.sleep(1)
        else:
            continue


@shared_task(bind=True)
def incep_async_tasks(self, id=None, user=None, sql=None, sqlsha1=None, host=None, database=None, exec_status=None,
                      backup=None):
    # 更新任务状态为：PROGRESS
    self.update_state(state="PROGRESS", meta={'user': user, 'host': host, 'database': database, 'sqlsha1': sqlsha1})

    of_audit = IncepSqlCheck(sql, host, database, user)

    # 执行SQL
    exec_result = of_audit.run_exec(0, backup)

    # 更新任务进度
    update_tasks_status(id=id, exec_result=exec_result, exec_status=exec_status)


@shared_task
def stop_incep_osc(user, id=None, celery_task_id=None):
    obj = IncepMakeExecTask.objects.get(id=id)
    host = obj.dst_host
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
        of_audit = IncepSqlCheck(sql, host, database, user)
        of_audit.run_status(0)

        # 更新任务进度
        update_tasks_status(id=id, exec_status=exec_status)

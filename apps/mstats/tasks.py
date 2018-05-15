# -*- coding:utf-8 -*-
# edit by fuzongfei
import time

import sys
from celery import shared_task, states, task
import logging

from celery.result import AsyncResult
from django.core.mail import EmailMessage
from django.template.loader import render_to_string

from AuditSQL.settings import EMAIL_FROM
from mstats.utils import ParamikoOutput
from multiprocessing import Pool

logger = logging.getLogger('django')


@shared_task
def backup_schema(**kwargs):
    ssh_host = kwargs.get('ssh_host')
    ssh_user = kwargs.get('ssh_user')
    ssh_password = kwargs.get('ssh_password')
    ssh_port = kwargs.get('ssh_port')
    backup_cmd = kwargs.get('backup_cmd')

    paramiko_output = ParamikoOutput(ssh_host, ssh_port, ssh_user, ssh_password)

    # 执行xtrabackup备份
    xtrabackup_cmd = backup_cmd.get('xtrabackup_cmd')
    if xtrabackup_cmd:
        res = paramiko_output.run(xtrabackup_cmd)
        logger.info(res)

    # 执行mysqldump全备
    mysqldump_full_cmd = backup_cmd.get('mysqldump_full_cmd')
    if mysqldump_full_cmd:
        for cmd in mysqldump_full_cmd:
            res = paramiko_output.run(cmd)
        logger.info(res)

    # 执行mysqldump单表备份
    # mysqldump_single_cmd = backup_cmd.get('mysqldump_single_cmd')
    # if mysqldump_single_cmd:
    #     for cmd in mysqldump_single_cmd:
    #         res = paramiko_output.run(cmd)
    #     logger.info(res)
    # out_history = paramiko_output(ssh_host, ssh_port, ssh_user, ssh_password, history_cmd)
    # del out_history[0]
    # del out_history[-1]

    # email_html_body = render_to_string('_send_backup_mail.html', {
    #     'data': '\n'.join(out_history),
    # })
    #
    # title = f'mysql backup for {ssh_host}'
    #
    # # 发送邮件
    # msg = EmailMessage(subject=title,
    #                    body=email_html_body,
    #                    from_email=EMAIL_FROM,
    #                    to=receiver,
    #                    )
    # msg.content_subtype = "html"
    #
    # msg.send()


@shared_task(bind=True)
def test_mes(self, id, user):
    self.update_state(state="PROGRESS", meta={'user': user, 'id': id, 'gas': 'ddd'})
    print('start...')
    time.sleep(10)
    print('end...')
    return 'ss'


@shared_task()
def aaa(task_id):
    print(task_id)
    task = AsyncResult(task_id)
    print(task.result)
    print(task.state)
    if task.state == 'PROGRESS':
        id = task.result['id']
        user = task.result['user']
        for i in range(1, 11):
            time.sleep(1)
            print(f'\r{user} 任务ID:{id} 进度: {i*10}%')

    # result.get(on_message=on_raw_message, propagate=False)

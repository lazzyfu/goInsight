# -*- coding:utf-8 -*-
# edit by fuzongfei
from celery import shared_task
from django.core.mail import EmailMessage
from django.template.loader import render_to_string

from AuditSQL.settings import EMAIL_FROM
from project_manager.models import DomainName


@shared_task
def send_create_user_mail(**kwargs):
    receiver = []
    username = kwargs.get('username')
    receiver.append(kwargs.get('email'))
    password = kwargs.get('password')

    # 向_commit_mail.html渲染data数据
    domain_name = ''
    if DomainName.objects.filter().first():
        domain_name = DomainName.objects.get().domain_name
    email_html_body = render_to_string('_send_create_user_mail.html', {
        'username': username,
        'password': password,
        'domain_name': domain_name
    })

    # 发送邮件
    msg = EmailMessage(subject=u'账号开通',
                       body=email_html_body,
                       from_email=EMAIL_FROM,
                       to=receiver,
                       )
    msg.content_subtype = "html"

    msg.send()

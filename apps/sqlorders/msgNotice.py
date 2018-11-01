# -*- coding:utf-8 -*-
# edit by fuzongfei

from celery import shared_task
from dingtalkchatbot.chatbot import DingtalkChatbot
from django.core.mail import EmailMessage
from django.template.loader import render_to_string
from django.utils import timezone

from sqlorders.models import SqlOrdersContents, SqlOrdersEnvironment, SysConfig, SqlOrderReply
from users.models import UserAccounts
from sqlaudit.settings import EMAIL_FROM


def get_user_email(id):
    obj = SqlOrdersContents.objects.get(pk=id)
    proposer = list(UserAccounts.objects.filter(username=obj.proposer).values_list('email', flat=True))
    auditor = list(UserAccounts.objects.filter(username=obj.auditor).values_list('email', flat=True))
    email_cc = list(obj.email_cc.split(','))
    return proposer, auditor, email_cc


class SqlOrdersMsgPull(object):
    def __init__(self, id, user, type, addition_info=None):
        """
        :param user: 用户
        :param id: 记录id
        :param type: 类型，包括：提交、审核、反馈、关闭等
        """
        self.user = user
        self.id = id
        self.type = type
        self.addition_info = addition_info

        # 发送通知里面的域名提示
        if SysConfig.objects.get(key='domain_name_tips').is_enabled == '0':
            self.domain_name_tips = SysConfig.objects.get(key='domain_name_tips').value

    def mail_notice(self):
        if self.type == 'commit':
            send_commit_mail.delay(id=self.id, domain_name_tips=self.domain_name_tips)
        else:
            send_verify_mail.delay(id=self.id, user=self.user, type=self.type,
                                   addition_info=self.addition_info, domain_name_tips=self.domain_name_tips)

    def dingding_notice(self, webhook):
        dingding_push.delay(id=self.id, user=self.user, type=self.type, webhook=webhook,
                            addition_info=self.addition_info, domain_name_tips=self.domain_name_tips)

    def weixin_notice(self):
        pass

    def run(self):
        # 判断系统是否开启了相关通知
        if SysConfig.objects.get(key='email_push').is_enabled == '0':
            self.mail_notice()

        if SysConfig.objects.get(key='dingding_push').is_enabled == '0':
            webhook = SysConfig.objects.get(key='dingding_push').value
            self.dingding_notice(webhook)

        if SysConfig.objects.get(key='weixin_push').is_enabled == '0':
            self.weixin_notice()


@shared_task
def send_commit_mail(id, user=None, type=None, addition_info=None, domain_name_tips=None):
    """发送提交工单邮件"""
    proposer, auditor, email_cc = get_user_email(id)
    # 收件人
    email_receiver = list(set(proposer + auditor))
    # 抄送人
    email_cc = email_cc

    # 向_commit_mail.html渲染data数据
    data = SqlOrdersContents.objects.get(pk=id)

    # 查询工单环境
    envi_name = SqlOrdersEnvironment.objects.get(envi_id=data.envi_id).envi_name

    email_html_body = render_to_string('mailnotice/_sqlorders_commit_mail.html', {
        'data': data,
        'domain_name_tips': domain_name_tips,
        'envi_name': envi_name})

    # 发送邮件
    msg = EmailMessage(subject=data.title,
                       body=email_html_body,
                       from_email=EMAIL_FROM,
                       to=email_receiver,
                       cc=email_cc
                       )
    msg.content_subtype = "html"
    msg.send()


@shared_task
def send_verify_mail(id, user=None, type=None, addition_info=None, domain_name_tips=None):
    envi_name = None
    if type == 'reply':
        data = SqlOrderReply.objects.get(pk=id)
        id = data.reply_id
        title = SqlOrdersContents.objects.get(pk=id).title
    else:
        data = SqlOrdersContents.objects.get(pk=id)
        title = data.title
        # 查询工单环境
        envi_name = SqlOrdersEnvironment.objects.get(envi_id=data.envi_id).envi_name

    """发送工单确认邮件"""
    proposer, auditor, email_cc = get_user_email(id)
    # 收件人
    email_receiver = list(set(proposer + auditor))
    # 抄送人
    email_cc = email_cc

    # 向mail_template.html渲染data数据
    email_html_body = render_to_string('mailnotice/_sqlorders_verify_mail.html', {
        'data': data,
        'type': type,
        'envi_name': envi_name,
        'domain_name_tips': domain_name_tips,
        'addition_info': addition_info,
        'username': user
    })

    # 发送邮件
    headers = {'Reply: ': email_receiver}
    title = 'Re: ' + title
    msg = EmailMessage(subject=title,
                       body=email_html_body,
                       from_email=EMAIL_FROM,
                       to=email_receiver,
                       cc=email_cc,
                       reply_to=email_receiver,
                       headers=headers)
    msg.content_subtype = "html"
    msg.send()


@shared_task
def dingding_push(id, user=None, type=None, addition_info=None, webhook=None, domain_name_tips=None):
    xiaoding = DingtalkChatbot(webhook)
    # 查询工单环境
    envi_name = None

    if type == 'reply':
        obj = SqlOrderReply.objects.get(pk=id)
        record_id = obj.reply_id
        reply_contents = obj.reply_contents
        data = SqlOrdersContents.objects.get(pk=record_id)
    else:
        # 获取数据
        data = SqlOrdersContents.objects.get(pk=id)
        envi_name = SqlOrdersEnvironment.objects.get(envi_id=data.envi_id).envi_name
    # 如果用户手机号存在，钉钉直接@mobile
    # 如果手机号不存在，钉钉直接@all
    proposer_mobile = UserAccounts.objects.get(username=data.proposer).mobile
    auditor_mobile = UserAccounts.objects.get(username=data.auditor).mobile

    # 提交
    if type == 'commit':
        text = f"您好、{user}提交了审核内容，◕‿◕\n" \
               f"标题: {data.title}\n" \
               f"环境: {envi_name}\n" \
               f"类型: {data.sql_type}\n" \
               f"主机: {data.host}\n" \
               f"端口: {data.port}\n" \
               f"库名: {data.database}\n" \
               f"审核人: {data.auditor}\n" \
               f"上线版本号: {data.task_version}\n" \
               f"需求: {data.description}\n" \
               f"URL: {domain_name_tips}/sqlorders/sql_orders_list/{data.envi_id} \n" \
               f"提交时间: {timezone.localtime(data.created_at).strftime('%Y-%m-%d %H:%M:%S')}\n"

        if auditor_mobile:
            xiaoding.send_text(msg=text, at_mobiles=[auditor_mobile])
        else:
            xiaoding.send_text(msg=text, is_at_all=True)
    # 审核
    elif type == 'approve':
        if data.progress == '2':
            text = f"您好、{user}审核已通过，◕‿◕\n" \
                   f"标题: {data.title}\n" \
                   f"环境: {envi_name}\n" \
                   f"附加信息: {addition_info}\n" \
                   f"URL: {domain_name_tips}/sqlorders/sql_orders_list/{data.envi_id} \n" \
                   f"审核时间: {timezone.localtime(data.operate_time).strftime('%Y-%m-%d %H:%M:%S')}\n"
        elif data.progress == '1':
            text = f"您好、{user}审核未通过，◕﹏◕\n" \
                   f"标题: {data.title}\n" \
                   f"环境: {envi_name}\n" \
                   f"附加信息: {addition_info}\n" \
                   f"URL: {domain_name_tips}/sqlorders/sql_orders_list/{data.envi_id} \n" \
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
                   f"环境: {envi_name}\n" \
                   f"附加信息: {addition_info}\n" \
                   f"URL: {domain_name_tips}/sqlorders/sql_orders_list/{data.envi_id} \n" \
                   f"处理时间: {timezone.localtime(data.updated_at).strftime('%Y-%m-%d %H:%M:%S')}\n"
        elif data.progress == '4':
            text = f"您好、{user}处理完成，◕‿◕\n" \
                   f"标题: {data.title}\n" \
                   f"环境: {envi_name}\n" \
                   f"附加信息: {addition_info}\n" \
                   f"URL: {domain_name_tips}/sqlorders/sql_orders_list/{data.envi_id} \n" \
                   f"完成时间: {timezone.localtime(data.updated_at).strftime('%Y-%m-%d %H:%M:%S')}\n"
        if proposer_mobile and auditor_mobile:
            xiaoding.send_text(msg=text, at_mobiles=[proposer_mobile, auditor_mobile])
        else:
            xiaoding.send_text(msg=text, is_at_all=True)
    # 关闭
    elif type == 'close':
        if data.progress == '5':
            text = f"您好、{user}关闭了记录，请不要处理，◕‿◕\n" \
                   f"标题: {data.title}\n" \
                   f"环境: {envi_name}\n" \
                   f"附加信息: {addition_info}\n" \
                   f"URL: {domain_name_tips}/sqlorders/sql_orders_list/{data.envi_id} \n" \
                   f"关闭时间: {timezone.localtime(data.close_time).strftime('%Y-%m-%d %H:%M:%S')}\n"
        if proposer_mobile and auditor_mobile:
            xiaoding.send_text(msg=text, at_mobiles=[proposer_mobile, auditor_mobile])
        else:
            xiaoding.send_text(msg=text, is_at_all=True)

    # 回复
    elif type == 'reply':
        text = f"您好、{user}回复了内容，◕‿◕\n" \
               f"标题: {data.title}\n" \
               f"回复内容: {reply_contents}\n" \
               f"URL: {domain_name_tips}/sqlorders/sql_orders_details/{data.pk}/ \n" \
               f"回复时间: {timezone.localtime(data.updated_at).strftime('%Y-%m-%d %H:%M:%S')}\n"
        if proposer_mobile and auditor_mobile:
            xiaoding.send_text(msg=text, at_mobiles=[proposer_mobile, auditor_mobile])
        else:
            xiaoding.send_text(msg=text, is_at_all=True)

    # 钩子
    elif type == 'hook':
        text = f"您好、{user}扭转了工单，◕‿◕\n" \
               f"标题: {data.title}\n" \
               f"新环境: {envi_name}\n" \
               f"类型: {data.sql_type}\n" \
               f"主机: {data.host}\n" \
               f"端口: {data.port}\n" \
               f"库名: {data.database}\n" \
               f"审核DBA: {data.auditor}\n" \
               f"上线版本号: {data.task_version}\n" \
               f"需求: {data.description}\n" \
               f"URL: {domain_name_tips}/sqlorders/sql_orders_list/{data.envi_id} \n" \
               f"扭转时间: {timezone.localtime(data.updated_at).strftime('%Y-%m-%d %H:%M:%S')}\n"
        if proposer_mobile and auditor_mobile:
            xiaoding.send_text(msg=text, at_mobiles=[proposer_mobile, auditor_mobile])
        else:
            xiaoding.send_text(msg=text, is_at_all=True)

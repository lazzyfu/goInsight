# -*- coding:utf-8 -*-
# edit by fuzongfei
import json

from celery import shared_task
from dingtalkchatbot.chatbot import DingtalkChatbot
from django.core.mail import EmailMessage
from django.template.loader import render_to_string
from django.utils import timezone

from config.config import EMAIL, DINGDING
from opsql.settings import EMAIL_FROM, DOMAIN_NAME
from orders.models import Orders, OrderReply
from users.models import UserAccounts


def get_user_email(id):
    obj = Orders.objects.get(pk=id)
    # 获取申请人的邮箱地址
    applicant_email = list(UserAccounts.objects.filter(username=obj.applicant).values_list('email', flat=True))
    # 获取审核人的邮箱地址
    auditor = [i['user'] for i in json.loads(obj.auditor)]
    auditor_email = list(set(UserAccounts.objects.filter(username__in=auditor).values_list('email', flat=True)))
    # 获取复核人的邮箱地址
    reviewer = [i['user'] for i in json.loads(obj.reviewer)]
    reviewer_email = list(set(UserAccounts.objects.filter(username__in=reviewer).values_list('email', flat=True)))
    # 获取抄送人的邮箱地址
    email_cc_email = list(
        set(UserAccounts.objects.filter(username__in=obj.email_cc.split(',')).values_list('email', flat=True)))
    return applicant_email, auditor_email, reviewer_email, email_cc_email


class MsgPush(object):
    def __init__(self, id=None, user=None, type=None, msg=None):
        """
        :param user: 用户
        :param id: 记录id
        :param type: ['commit', 'approve', 'feedback', 'hook', 'review, 'close', 'reply']
        """
        self.user = user
        self.id = id
        self.type = type
        self.msg = msg

        # 消息推送里面的域名提示
        self.domain_name_tips = DOMAIN_NAME['value']

    def mail_notice(self):
        if self.type == 'commit':
            send_commit_mail.delay(id=self.id, domain_name_tips=self.domain_name_tips)
        else:
            send_op_mail.delay(id=self.id, user=self.user, type=self.type,
                               msg=self.msg, domain_name_tips=self.domain_name_tips)

    def dingding_notice(self, webhook):
        dingding_push.delay(id=self.id, user=self.user, type=self.type, webhook=webhook,
                            msg=self.msg, domain_name_tips=self.domain_name_tips)

    def send(self):
        # 判断是否启用了邮件通知
        if EMAIL['enable'] is True:
            self.mail_notice()

        # 判断是否启用了钉钉通知
        if DINGDING['enable'] is True:
            webhook = DINGDING['webhook']
            self.dingding_notice(webhook)


@shared_task
def send_commit_mail(id=None, domain_name_tips=None):
    """发送工单提交邮件"""
    applicant_email, auditor_email, reviewer_email, email_cc_email = get_user_email(id)
    # 收件人
    receiver_email = list(set(applicant_email + auditor_email + reviewer_email))
    # 向_commit_mail.html渲染data数据
    data = Orders.objects.get(pk=id)
    # 获取工单的审核人和复核人
    auditor = ', '.join([x['user'] for x in json.loads(data.auditor)])
    reviewer = ', '.join([x['user'] for x in json.loads(data.reviewer)])

    email_html_body = render_to_string(
        'mail/orders_commit_mail.html',
        {
            'data': data,
            'auditor': auditor,
            'reviewer': reviewer,
            'domain_name_tips': domain_name_tips,
        }
    )

    # 发送邮件
    msg = EmailMessage(subject=data.title,
                       body=email_html_body,
                       from_email=EMAIL_FROM,
                       to=receiver_email,
                       cc=email_cc_email
                       )
    msg.content_subtype = "html"
    msg.send()


@shared_task
def send_op_mail(id, user=None, type=None, msg=None, domain_name_tips=None):
    if type == 'reply':
        data = OrderReply.objects.get(pk=id)
        title = Orders.objects.get(pk=data.reply_id).title
        # 工单id，用于查找收件人
        order_id = data.reply_id
    else:
        data = Orders.objects.get(pk=id)
        title = data.title
        order_id = id

    applicant_email, auditor_email, reviewer_email, email_cc_email = get_user_email(order_id)
    # 收件人
    receiver_email = list(set(applicant_email + auditor_email + reviewer_email))

    # 向mail_template.html渲染data数据
    email_html_body = render_to_string('mail/orders_op_mail.html', {
        'data': data,
        'type': type,
        'domain_name_tips': domain_name_tips,
        'msg': msg,
        'username': user
    })

    # 发送邮件
    headers = {'Reply: ': receiver_email}
    title = 'Re: ' + title
    msg = EmailMessage(subject=title,
                       body=email_html_body,
                       from_email=EMAIL_FROM,
                       to=receiver_email,
                       cc=email_cc_email,
                       reply_to=receiver_email,
                       headers=headers)
    msg.content_subtype = "html"
    msg.send()


@shared_task
def dingding_push(id, user=None, type=None, msg=None, webhook=None, domain_name_tips=None):
    xiaoding = DingtalkChatbot(webhook)

    if type == 'reply':
        obj = OrderReply.objects.get(pk=id)
        reply_contents = obj.reply_contents
        data = Orders.objects.get(pk=obj.reply_id)
    else:
        # 获取数据
        data = Orders.objects.get(pk=id)
    # 如果用户手机号存在，钉钉直接@mobile
    # 如果手机号不存在，钉钉直接@all
    applicant_mobile = UserAccounts.objects.get(username=data.applicant).mobile
    auditor = [i['user'] for i in json.loads(data.auditor)]
    auditor_mobile = list(set(UserAccounts.objects.filter(username__in=auditor).values_list('mobile', flat=True)))
    reviewer = ','.join([i['user'] for i in json.loads(data.reviewer)])

    # 提交
    if type == 'commit':
        text = f"您好、{user}提交了审核内容，◕‿◕\n" \
            f"标题: {data.title}\n" \
            f"环境: {data.envi}\n" \
            f"类型: {data.sql_type}\n" \
            f"主机: {data.host}\n" \
            f"端口: {data.port}\n" \
            f"库名: {data.database}\n" \
            f"审核人: {auditor}\n" \
            f"复核人: {reviewer}\n" \
            f"上线版本号: {data.version}\n" \
            f"需求: {data.description}\n" \
            f"URL: {domain_name_tips}/orders/detail/{data.pk}/ \n" \
            f"提交时间: {timezone.localtime(data.created_at).strftime('%Y-%m-%d %H:%M:%S')}\n"

        if auditor_mobile:
            xiaoding.send_text(msg=text, at_mobiles=auditor_mobile)
        else:
            xiaoding.send_text(msg=text, is_at_all=True)
    # 审核
    elif type == 'approve':
        if data.progress == '2':
            text = f"您好、{user}审核已通过，◕‿◕\n" \
                f"标题: {data.title}\n" \
                f"环境: {data.envi}\n" \
                f"附加信息: {msg}\n" \
                f"URL: {domain_name_tips}/orders/envi/{data.envi_id}/ \n" \
                f"审核时间: {timezone.localtime(data.updated_at).strftime('%Y-%m-%d %H:%M:%S')}\n"
        elif data.progress == '1':
            text = f"您好、{user}审核未通过，◕﹏◕\n" \
                f"标题: {data.title}\n" \
                f"环境: {data.envi}\n" \
                f"附加信息: {msg}\n" \
                f"URL: {domain_name_tips}/orders/envi/{data.envi_id}/ \n" \
                f"审核时间: {timezone.localtime(data.updated_at).strftime('%Y-%m-%d %H:%M:%S')}\n"
        if applicant_mobile:
            xiaoding.send_text(msg=text, at_mobiles=[applicant_mobile])
        else:
            xiaoding.send_text(msg=text, is_at_all=True)
    # 反馈
    elif type == 'feedback':
        if data.progress == '3':
            text = f"您好、{user}正在处理中，请稍后，◕‿◕\n" \
                f"标题: {data.title}\n" \
                f"环境: {data.envi}\n" \
                f"附加信息: {msg}\n" \
                f"URL: {domain_name_tips}/orders/envi/{data.envi_id}/ \n" \
                f"处理时间: {timezone.localtime(data.updated_at).strftime('%Y-%m-%d %H:%M:%S')}\n"
        elif data.progress == '4':
            text = f"您好、{user}处理完成，◕‿◕\n" \
                f"标题: {data.title}\n" \
                f"环境: {data.envi}\n" \
                f"附加信息: {msg}\n" \
                f"URL: {domain_name_tips}/orders/envi/{data.envi_id}/ \n" \
                f"完成时间: {timezone.localtime(data.updated_at).strftime('%Y-%m-%d %H:%M:%S')}\n"
        if applicant_mobile and auditor_mobile:
            xiaoding.send_text(msg=text, at_mobiles=[applicant_mobile, auditor_mobile])
        else:
            xiaoding.send_text(msg=text, is_at_all=True)
    # 关闭
    elif type == 'close':
        if data.progress == '5':
            text = f"您好、{user}关闭了记录，请不要处理，◕‿◕\n" \
                f"标题: {data.title}\n" \
                f"环境: {data.envi}\n" \
                f"附加信息: {msg}\n" \
                f"URL: {domain_name_tips}/orders/envi/{data.envi_id}/ \n" \
                f"关闭时间: {timezone.localtime(data.close_time).strftime('%Y-%m-%d %H:%M:%S')}\n"
        if applicant_mobile and auditor_mobile:
            xiaoding.send_text(msg=text, at_mobiles=[applicant_mobile, auditor_mobile])
        else:
            xiaoding.send_text(msg=text, is_at_all=True)

    # 回复
    elif type == 'reply':
        text = f"您好、{user}回复了内容，◕‿◕\n" \
            f"标题: {data.title}\n" \
            f"回复内容: {reply_contents}\n" \
            f"URL: {domain_name_tips}/orders/detail/{data.pk}/ \n" \
            f"回复时间: {timezone.localtime(data.updated_at).strftime('%Y-%m-%d %H:%M:%S')}\n"
        if applicant_mobile and auditor_mobile:
            xiaoding.send_text(msg=text, at_mobiles=[applicant_mobile, auditor_mobile])
        else:
            xiaoding.send_text(msg=text, is_at_all=True)

    # 复核
    elif type == 'review':
        text = f"您好、{user}已复核，◕‿◕\n" \
            f"标题: {data.title}\n" \
            f"附加信息: {msg}\n" \
            f"复核时间: {timezone.localtime(data.updated_at).strftime('%Y-%m-%d %H:%M:%S')}\n"
        if applicant_mobile and auditor_mobile:
            xiaoding.send_text(msg=text, at_mobiles=[applicant_mobile, auditor_mobile])
        else:
            xiaoding.send_text(msg=text, is_at_all=True)

    # 钩子
    elif type == 'hook':
        text = f"您好、{user}扭转了工单，◕‿◕\n" \
            f"标题: {data.title}\n" \
            f"新环境: {data.envi}\n" \
            f"类型: {data.sql_type}\n" \
            f"主机: {data.host}\n" \
            f"端口: {data.port}\n" \
            f"库名: {data.database}\n" \
            f"上线版本号: {data.version}\n" \
            f"需求: {data.description}\n" \
            f"URL: {domain_name_tips}/orders/envi/{data.envi_id}/ \n" \
            f"扭转时间: {timezone.localtime(data.updated_at).strftime('%Y-%m-%d %H:%M:%S')}\n"
        if applicant_mobile and auditor_mobile:
            xiaoding.send_text(msg=text, at_mobiles=[applicant_mobile, auditor_mobile])
        else:
            xiaoding.send_text(msg=text, is_at_all=True)

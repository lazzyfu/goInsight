# -*- coding:utf-8 -*-
# edit by fuzongfei
import json

import requests
from celery.utils.log import get_task_logger
from dingtalkchatbot.chatbot import DingtalkChatbot
from django.core.mail import EmailMessage

from config import NOTICE_URL, NOTICE
from sqlorders import models, utils
from users.models import UserAccounts

logger = get_task_logger('celery.logger')


class MsgNotice(object):
    """
    工单通知模块，支持：钉钉/邮件/企业微信
    注：钉钉和企业微信需要拉群，如果用户的手机号存在，则在群里面at用户，否则at all
    进阶：当然您也可以直接对接钉钉或企业微信的企业接口，这个需要您自己开发接口
          如此，消息会直接at用户，不需要拉群，会发送到个人消息中心
    """

    def __init__(self, **kwargs):
        # 工单的基本信息
        self.pk = kwargs.get('pk')  # 工单的主键
        self.op = kwargs.get('op')  # 操作，如：commit/approve/feedback等
        self.username = kwargs.get('username')  # 用户

        # 导出文件信息
        self.export_file_name = kwargs.get('export_file_name')
        self.export_file_encryption_key = kwargs.get('export_file_encryption_key')
        self.export_sql = kwargs.get('export_sql')

        # 获取工单信息
        self.obj = models.DbOrders.objects.get(pk=self.pk)
        self.auditor = [x['user'] for x in json.loads(self.obj.auditor)]
        self.reviewer = [x['user'] for x in json.loads(self.obj.reviewer)]
        self.email_cc = list(filter(str, [x for x in self.obj.email_cc.split(',')]))
        self.notice_url = f"{NOTICE_URL}/sqlorders/detail/{self.obj.order_id}"
        self.export_notice_url = f"{NOTICE_URL}/sqlorders/export/download/{self.export_file_name}"

        # 获取申请人mobile&email
        self.applicant_mobile = list(
            filter(None, UserAccounts.objects.filter(username=self.obj.applicant).values_list('mobile', flat=True))
        )
        self.applicant_email = list(
            filter(None, UserAccounts.objects.filter(username=self.obj.applicant).values_list('email', flat=True))
        )

        # 获取审核人mobile&email
        self.auditor_mobile = list(
            filter(None, UserAccounts.objects.filter(username__in=self.auditor).values_list('mobile', flat=True))
        )
        self.auditor_email = list(
            filter(None, UserAccounts.objects.filter(username__in=self.auditor).values_list('email', flat=True))
        )

        # 获取复核人mobile&email
        self.reviewer_mobile = list(
            filter(None, UserAccounts.objects.filter(username__in=self.reviewer).values_list('mobile', flat=True))
        )
        self.reviewer_email = list(
            filter(None, UserAccounts.objects.filter(username__in=self.reviewer).values_list('email', flat=True))
        )

        # 获取抄送人mobile&email
        self.email_cc_mobile = list(
            filter(None, UserAccounts.objects.filter(username__in=self.email_cc).values_list('mobile', flat=True))
        )
        self.email_cc_email = list(
            filter(None, UserAccounts.objects.filter(username__in=self.email_cc).values_list('email', flat=True))
        )

        # 格式化工单进度
        self.fmt_progress = dict(utils.sqlProgressChoice).get(self.obj.progress)

    def push(self, content=None, mobiles=None, emails=None):
        if NOTICE['DINGDING']['enabled']:
            dingding_content = content.copy()
            dingding_content.append(f"请访问 {self.notice_url} 查看详情")
            self.dingding(dingding_content, mobiles)
        if NOTICE['WEIXIN']['enabled']:
            weixin_content = content.copy()
            weixin_content.append(f"\n请访问 [{self.notice_url}]({self.notice_url}) 查看详情")
            self.weixin(weixin_content, mobiles)
        if NOTICE['MAIL']['enabled']:
            mail_content = content.copy()
            mail_content.append(f"请访问 {self.notice_url} 查看详情")
            self.mail(mail_content, emails)

    def dingding(self, content=None, mobiles=None):
        # 此处您可以改写下面的接口对接企业接口，消息会推送到个人的消息中心
        key = NOTICE['DINGDING']['key']
        xiaoding = DingtalkChatbot(NOTICE['DINGDING']['webhook'])
        content = '\n'.join(['\n\n'.join(content), key])
        # 如果通知人mobile存在，at指定的用户，否则at all
        if mobiles:
            xiaoding.send_markdown(title='新的工单', text=content, at_mobiles=mobiles)
        else:
            xiaoding.send_markdown(title='新的工单', text=content, is_at_all=True)

    def weixin(self, content, mobiles=None):
        # 此处您可以改写下面的接口对接企业接口，消息会推送到个人的消息中心
        webhook = NOTICE['WEIXIN']['webhook']
        headers = {'Content-Type': 'application/x-www-form-urlencoded'}
        data = {
            "msgtype": "markdown",
            "markdown": {"content": '\n'.join(content)},
            "mentioned_mobile_list": mobiles
        }
        request = requests.post(
            url=webhook,
            data=json.dumps(data),
            timeout=3,
            headers=headers
        )
        # 记录下请求响应的信息
        logger.info(request.json())

    def mail(self, content, emails=None):
        try:
            if self.op in ['_commit', '_export']:
                # 发送邮件
                msg = EmailMessage(subject=self.obj.title,
                                   body='<br>'.join(content),
                                   from_email=NOTICE['MAIL']['email_host_user'],
                                   to=emails,
                                   )
            else:
                # 发送邮件
                headers = {'Reply: ': emails}
                title = 'Re: ' + self.obj.title
                msg = EmailMessage(subject=title,
                                   body='<br>'.join(content),
                                   from_email=NOTICE['MAIL']['email_host_user'],
                                   to=emails,
                                   reply_to=emails,
                                   headers=headers)
            msg.content_subtype = "html"
            msg.send()
        except Exception as err:
            logger.error(err)

    def _commit(self):
        """提交工单时，发送消息"""
        content = [
            f"您好，{self.username}提交了{self.obj.sql_type}工单，^_^",
            f">标题: {self.obj.title}",
            f">申请人: {self.obj.applicant}",
            f">审核人: {','.join(self.auditor)}",
            f">复核人: {','.join(self.reviewer)}",
            f">备注: {self.obj.remark}",
            f">环境: {self.obj.env}",
            f">DB类型: {dict(utils.rdsCategory).get(int(self.obj.rds_category))}",
            f">库名: {self.obj.database}"
        ]

        self.push(
            content=content,
            mobiles=list(set(
                self.applicant_mobile +
                self.auditor_mobile +
                self.reviewer_mobile +
                self.email_cc_mobile
            )),
            emails=list(set(
                self.applicant_email +
                self.auditor_email +
                self.reviewer_email +
                self.email_cc_email
            ))
        )

    def _approve(self):
        """审核操作，发送信息"""
        # 发送审核动作通知
        finish_status = []
        for i in json.loads(self.obj.auditor):
            # 增加超级审核人
            if i['status'] in (1, 2):
                if i['is_superuser'] == 1:
                    finish_status.append(f"{i['user']}(超级审核人)")
                else:
                    finish_status.append(f"{i['user']}")
        content = [
            f"您好，{self.username}审核了工单，请关注，^_^",
            f">工单标题: {self.obj.title}",
            f">审核人: {','.join(finish_status)}",
            f">工单备注: {self.obj.remark}",
            f">工单状态: {self.fmt_progress}",
            f">环境: {self.obj.env}",
            f">DB类型: {dict(utils.rdsCategory).get(int(self.obj.rds_category))}",
            f">库名: {self.obj.database}"
        ]
        self.push(
            content=content,
            mobiles=list(set(
                self.applicant_mobile +
                self.auditor_mobile
            )),
            emails=list(set(
                self.applicant_email +
                self.auditor_email
            ))
        )

        # 审核通过
        if self.obj.progress == 2:
            content = [
                f"您好，工单已审核通过，请关注，^_^",
                f">工单标题: {self.obj.title}",
                f">申请人: {self.obj.applicant}",
                f">工单备注: {self.obj.remark}",
                f">工单状态: {self.fmt_progress}",
                f">环境: {self.obj.env}",
                f">DB类型: {dict(utils.rdsCategory).get(int(self.obj.rds_category))}",
                f">库名: {self.obj.database}"
            ]

            self.push(
                content=content,
                mobiles=list(set(
                    self.applicant_mobile +
                    self.auditor_mobile
                )),
                emails=list(set(
                    self.applicant_email +
                    self.auditor_email
                ))
            )

    def _feedback(self):
        """反馈"""
        if self.obj.progress == 3:
            # 处理中
            content = [
                f"您好，工单正在执行中，请耐心等待，^_^",
                f">工单标题: {self.obj.title}",
                f">申请人: {self.obj.applicant}",
                f">工单备注: {self.obj.remark}",
                f">工单状态: {self.fmt_progress}",
                f">环境: {self.obj.env}",
                f">DB类型: {dict(utils.rdsCategory).get(int(self.obj.rds_category))}",
                f">库名: {self.obj.database}",
                f">执行人: {self.username}"
            ]
            self.push(
                content=content,
                mobiles=list(set(
                    self.applicant_mobile
                )),
                emails=list(set(
                    self.applicant_email
                ))
            )
        if self.obj.progress == 4:
            # 已完成
            content = [
                f"您好，工单已处理完成，请您及时核对[核对完成后，请点击「复核」按钮更改工单状态为: 已复核]，^_^",
                f">工单标题: {self.obj.title}",
                f">申请人: {self.obj.applicant}",
                f">工单备注: {self.obj.remark}",
                f">工单状态: {self.fmt_progress}",
                f">环境: {self.obj.env}",
                f">DB类型: {dict(utils.rdsCategory).get(int(self.obj.rds_category))}",
                f">库名: {self.obj.database}",
                f">执行人: {self.username}"
            ]
            self.push(
                content=content,
                mobiles=list(set(
                    self.reviewer_mobile
                )),
                emails=list(set(
                    self.reviewer_email
                ))
            )

    def _close(self):
        """关闭"""
        content = [
            f"您好，工单已被关闭 (⊙︿⊙)",
            f">工单标题: {self.obj.title}",
            f">申请人: {self.obj.applicant}",
            f">工单备注: {self.obj.remark}",
            f">工单状态: {self.fmt_progress}",
            f">环境: {self.obj.env}",
            f">DB类型: {dict(utils.rdsCategory).get(int(self.obj.rds_category))}",
            f">库名: {self.obj.database}",
            f">关闭人: {self.username}"
        ]
        self.push(
            content=content,
            mobiles=list(set(
                self.applicant_mobile
            )),
            emails=list(set(
                self.applicant_email
            ))
        )

    def _export(self):
        content = [f"您好，导出已完成，请拷贝链接到浏览器下载【机密信息、请勿外传!】 ^_^",
                   f">工单标题: {self.obj.title}",
                   f">申请人: {self.obj.applicant}",
                   f">解密秘钥: {self.export_file_encryption_key}",
                   f">下载地址: {self.export_notice_url}",
                   f">执行人: {self.username}",
                   f">原始SQL(截断，保留512字节): `{self.export_sql[:512]}`"
                   ]
        self.push(
            content=content,
            mobiles=list(set(
                self.applicant_mobile
            )),
            emails=list(set(
                self.applicant_email
            ))
        )

    def run(self):
        if self.op == '_commit':
            self._commit()
        if self.op == '_approve':
            self._approve()
        if self.op == '_feedback':
            self._feedback()
        if self.op == '_close':
            self._close()
        if self.op == '_export':
            self._export()

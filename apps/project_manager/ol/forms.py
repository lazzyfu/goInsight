# -*- coding:utf-8 -*-
# edit by fuzongfei
from datetime import datetime

from django import forms
from django.utils import timezone

from project_manager.inception.inception_api import IncepSqlCheck
from project_manager.models import AuditContents, OlAuditDetail, OlAuditContentsReply
from project_manager.models import operate_type_choice
from project_manager.tasks import send_commit_mail, send_verify_mail, send_reply_mail, xiaoding_pull


class OnlineAuditForm(forms.Form):
    title = forms.CharField(max_length=100, required=True, label=u'标题')
    verifier = forms.CharField(required=True, label=u'批准人')
    operator = forms.CharField(required=True, label=u'执行人')
    group_id = forms.CharField(required=True, label=u'项目组id')
    host = forms.CharField(required=True, label=u'数据库主机')
    database = forms.CharField(required=True, max_length=64, label=u'数据库')
    operate_type = forms.ChoiceField(choices=operate_type_choice, label=u'操作类型，是DDL还是DML')
    contents = forms.CharField(widget=forms.Textarea, label=u'审核内容')

    def is_save(self, request):
        cleaned_data = super(OnlineAuditForm, self).clean()
        title = cleaned_data.get('title') + '_[' + datetime.now().strftime("%Y%m%d%H%M%S") + ']'
        verifier = cleaned_data.get('verifier')
        operator = cleaned_data.get('operator')
        email_cc = self.data.get('email_cc_id')
        group_id = cleaned_data.get('group_id')
        host = cleaned_data.get('host')
        database = cleaned_data.get('database')
        operate_type = cleaned_data.get('operate_type')
        contents = cleaned_data.get('contents')

        result = IncepSqlCheck(contents, host, database, request.user.username).is_check_pass()
        if result.get('status') == 2:
            context = result
        else:
            AuditContents.objects.create(
                title=title,
                operate_type=operate_type,
                host=host,
                database=database,
                group_id=group_id,
                proposer=request.user.username,
                operator=operator,
                verifier=verifier,
                email_cc=email_cc
            )

            # 向子表插入关联数据
            OlAuditDetail.objects.create(ol=AuditContents.objects.latest('id'), contents=contents)

            # 发送通知邮件
            latest_id = AuditContents.objects.latest('id').id
            send_commit_mail.delay(latest_id=latest_id)

            # 发送钉钉推送
            xiaoding_pull.delay(user=request.user.username, title=title, type='commit')

            context = {'status': 0, 'jump_url': '/projects/ol/ol_records/'}
        return context


class OnlineAuditRecordForm(forms.Form):
    limit_size = forms.IntegerField(required=True, label=u'每页显示数量')
    offset_size = forms.IntegerField(required=True, label=u'分页偏移量')
    search_content = forms.CharField(max_length=128, required=False, label='搜索内容')


class OnlineApproveForm(forms.Form):
    id = forms.IntegerField(required=True)
    status = forms.CharField(max_length=10, required=True)
    addition_info = forms.CharField(required=False)

    def is_save(self, request):
        cleaned_data = super(OnlineApproveForm, self).clean()
        id = cleaned_data.get('id')
        status = cleaned_data.get('status')
        addition_info = cleaned_data.get('addition_info')

        data = AuditContents.objects.get(pk=id)

        context = {}
        # 当记录关闭时
        if data.progress == '5':
            context = {'status': 2, 'msg': '该记录已被关闭、请不要重复提交'}
        # 当记录未关闭时
        else:
            if data.progress == '0' or data.progress == '1':
                # 当用户点击的是通过, 状态变为：已批准
                if status == u'通过':
                    data.progress = '2'
                    data.verifier_time = timezone.now()
                    data.save()
                    send_verify_mail.delay(latest_id=id,
                                           type='approve',
                                           username=request.user.username,
                                           user_role=request.user.user_role(),
                                           addition_info=addition_info)
                    # 发送钉钉推送
                    xiaoding_pull.delay(user=request.user.username, title=data.title, type='approve', progress='2')
                    context = {'status': 0, 'msg': '操作成功、审核通过'}

                # 当用户点击的是不通过, 状态变为：未批准
                elif status == u'不通过':
                    data.progress = '1'
                    data.verifier_time = timezone.now()
                    data.save()
                    send_verify_mail.delay(latest_id=id,
                                           type='approve',
                                           username=request.user.username,
                                           user_role=request.user.user_role(),
                                           addition_info=addition_info)
                    # 发送钉钉推送
                    xiaoding_pull.delay(user=request.user.username, title=data.title, type='approve', progress='1')
                    context = {'status': 0, 'msg': '操作成功、审核未通过'}

            # 其他情况
            else:
                context = {'status': 2, 'msg': '操作失败、请不要重复提交'}
        return context


class OnlineFeedbackForm(forms.Form):
    id = forms.IntegerField(required=True)
    status = forms.CharField(max_length=10, required=True)
    addition_info = forms.CharField(required=False)

    def is_save(self, request):
        cleaned_data = super(OnlineFeedbackForm, self).clean()
        id = cleaned_data.get('id')
        status = cleaned_data.get('status')
        addition_info = cleaned_data.get('addition_info')

        data = AuditContents.objects.get(pk=id)

        context = {}
        # 当记录关闭时
        if data.progress == '5':
            context = {'status': 2, 'msg': '该记录已被关闭、请不要重复提交'}
        # 当记录未关闭时
        else:
            # 当进度状态为：已批准或处理中时
            if data.progress == '2' or data.progress == '3':
                # 当用户点击的是处理中, 状态变为：处理中
                if status == u'处理中':
                    data.progress = '3'
                    data.save()
                    send_verify_mail.delay(latest_id=id,
                                           type='feedback',
                                           username=request.user.username,
                                           user_role=request.user.user_role(),
                                           addition_info=addition_info)
                    # 发送钉钉推送
                    xiaoding_pull.delay(user=request.user.username, title=data.title, type='feedback', progress='3')
                    context = {'status': 0, 'msg': '操作成功、正在处理中'}

                # 当用户点击的是已完成, 状态变为：已完成
                elif status == u'已完成':
                    data.progress = '4'
                    data.operate_time = timezone.now()
                    data.save()
                    send_verify_mail.delay(latest_id=id,
                                           type='feedback',
                                           username=request.user.username,
                                           user_role=request.user.user_role(),
                                           addition_info=addition_info)
                    # 发送钉钉推送
                    xiaoding_pull.delay(user=request.user.username, title=data.title, type='feedback', progress='4')
                    context = {'status': 0, 'msg': '操作成功、处理完成'}

            # 未批准
            elif data.progress == '1' or data.progress == '0':
                context = {'status': 2, 'msg': '操作失败、审核未通过'}
            # 其他情况
            else:
                context = {'status': 2, 'msg': '操作失败、请不要重复提交'}
        return context


class OnlineCloseForm(forms.Form):
    id = forms.IntegerField(required=True)
    status = forms.CharField(max_length=10, required=True)
    addition_info = forms.CharField(required=False)

    def is_save(self, request):
        cleaned_data = super(OnlineCloseForm, self).clean()
        id = cleaned_data.get('id')
        status = cleaned_data.get('status')
        addition_info = cleaned_data.get('addition_info')

        data = AuditContents.objects.get(pk=id)

        context = {}
        # 当记录关闭时
        if data.progress == '5':
            context = {'status': 2, 'msg': '该记录已被关闭、请不要重复提交'}
        # 当记录未关闭时
        else:
            if len(addition_info) >= 5:
                # 当进度为：处理中或已完成时
                if status == u'提交':
                    if data.progress == '3' or data.progress == '4':
                        context = {'status': 2, 'msg': '操作失败、数据正在处理中或已完成'}
                    else:
                        data.progress = '5'
                        data.close_user = request.user.username
                        data.close_reason = addition_info
                        data.close_time = timezone.now()
                        data.save()
                        send_verify_mail.delay(latest_id=id,
                                               type='close',
                                               username=request.user.username,
                                               user_role=request.user.user_role(),
                                               addition_info=addition_info)
                        # 发送钉钉推送
                        xiaoding_pull.delay(user=request.user.username, title=data.title, type='close', progress='5')
                        context = {'status': 0, 'msg': '操作成功、记录关闭成功'}

                elif status == u'结束':
                    context = {'status': 2, 'msg': '操作失败、关闭窗口'}
            else:
                context = {'status': 2, 'msg': '操作失败、<关闭原因>不能少于5个字符'}
        return context


class OnlineReplyForm(forms.Form):
    reply_id = forms.IntegerField(required=True)
    reply_contents = forms.CharField(widget=forms.Textarea, min_length=5)

    def is_save(self, request):
        cleaned_data = super(OnlineReplyForm, self).clean()
        reply_id = cleaned_data.get('reply_id')
        reply_contents = cleaned_data.get('reply_contents')
        OlAuditContentsReply.objects.create(
            reply_id=reply_id,
            user_id=request.user.uid,
            reply_contents=reply_contents)

        context = {'status': 0, 'msg': '回复成功'}
        latest_id = OlAuditContentsReply.objects.latest('id').id
        send_reply_mail.delay(latest_id=latest_id,
                              reply_id=reply_id,
                              username=request.user.username,
                              email=request.user.email)
        return context

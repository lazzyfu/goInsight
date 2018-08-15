# -*- coding:utf-8 -*-
# edit by fuzongfei
from datetime import datetime

from django import forms
from django.utils import timezone

from project_manager.inception.inception_api import IncepSqlCheck
from project_manager.models import AuditContents, OlAuditContentsReply
from project_manager.models import operate_type_choice
from project_manager.tasks import xiaoding_pull


class WorkOrderAuditForm(forms.Form):
    """提交工单"""
    title = forms.CharField(max_length=100, required=True, label=u'标题')
    url = forms.CharField(max_length=1024, required=False, label=u'上线的Confluence链接')
    tasks = forms.CharField(max_length=256, required=False, label=u'任务')
    operator = forms.CharField(required=True, label=u'DBA')
    database = forms.CharField(required=True, max_length=64, label=u'数据库')
    envi_desc = forms.IntegerField(required=True, label=u'环境：0：测试、1：staging、2：生产、3：线下其他环境')
    remark = forms.IntegerField(required=True, label=u'工单备注')
    jump_url = forms.CharField(required=True, max_length=128, label=u'跳转到的工单页面的url')
    operate_type = forms.ChoiceField(choices=operate_type_choice, label=u'操作类型，是DDL还是DML')
    contents = forms.CharField(widget=forms.Textarea, label=u'审核内容')

    def is_save(self, request):
        cleaned_data = super(WorkOrderAuditForm, self).clean()
        title = cleaned_data.get('title') + '_[' + datetime.now().strftime("%Y%m%d%H%M%S") + ']'
        url = cleaned_data.get('url')
        tasks = cleaned_data.get('tasks')
        operator = cleaned_data.get('operator')
        envi_desc = cleaned_data.get('envi_desc')
        remark = cleaned_data.get('remark')
        jump_url = cleaned_data.get('jump_url')
        host, port, database = cleaned_data.get('database').split(',')
        operate_type = cleaned_data.get('operate_type')
        contents = cleaned_data.get('contents')

        result = IncepSqlCheck(contents, host, port, database, request.user.username).is_check_pass()
        if result.get('status') == 2:
            context = result
        else:
            obj = AuditContents.objects.create(
                title=title,
                url=url,
                tasks=tasks,
                operate_type=operate_type,
                host=host,
                database=database,
                port=port,
                envi_desc=envi_desc,
                remark=remark,
                proposer=request.user.username,
                operator=operator,
                contents=contents
            )

            # 发送钉钉推送
            xiaoding_pull.delay(user=request.user.username, id=obj.id, type='commit')

            # 跳转到工单记录页面
            context = {'status': 0, 'jump_url': jump_url}
        return context


class OnlineAuditRecordForm(forms.Form):
    envi_desc = forms.IntegerField(required=True, label=u'环境')
    limit_size = forms.IntegerField(required=True, label=u'每页显示数量')
    offset_size = forms.IntegerField(required=True, label=u'分页偏移量')
    search_content = forms.CharField(max_length=128, required=False, label='搜索内容')


class WorkOrderApproveForm(forms.Form):
    id = forms.IntegerField(required=True)
    status = forms.CharField(max_length=10, required=True)
    addition_info = forms.CharField(required=False)

    def is_save(self, request):
        cleaned_data = super(WorkOrderApproveForm, self).clean()
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
                    data.operate_time = timezone.now()
                    data.save()
                    # 发送钉钉推送
                    xiaoding_pull.delay(user=request.user.username, id=id, type='approve',
                                        addition_info=addition_info)
                    context = {'status': 0, 'msg': '操作成功、审核通过'}

                # 当用户点击的是不通过, 状态变为：未批准
                elif status == u'不通过':
                    data.progress = '1'
                    data.operate_time = timezone.now()
                    data.save()
                    # 发送钉钉推送
                    xiaoding_pull.delay(user=request.user.username, id=id, type='approve',
                                        addition_info=addition_info)
                    context = {'status': 0, 'msg': '操作成功、审核未通过'}

            # 其他情况
            else:
                context = {'status': 2, 'msg': '操作失败、请不要重复提交'}
        return context


class WorkOrderFeedbackForm(forms.Form):
    id = forms.IntegerField(required=True)
    status = forms.CharField(max_length=10, required=True)
    addition_info = forms.CharField(required=False)

    def is_save(self, request):
        cleaned_data = super(WorkOrderFeedbackForm, self).clean()
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
                    data.updated_at = timezone.now()
                    data.save()
                    # 发送钉钉推送
                    xiaoding_pull.delay(user=request.user.username, id=id, type='feedback',
                                        addition_info=addition_info)
                    context = {'status': 0, 'msg': '操作成功、正在处理中'}

                # 当用户点击的是已完成, 状态变为：已完成
                elif status == u'已完成':
                    data.progress = '4'
                    data.updated_at = timezone.now()
                    data.save()
                    # 发送钉钉推送（包括任务进度）
                    xiaoding_pull.delay(user=request.user.username, id=id, type='feedback',
                                        addition_info=addition_info)
                    context = {'status': 0, 'msg': '操作成功、处理完成'}

            # 未批准
            elif data.progress == '1' or data.progress == '0':
                context = {'status': 2, 'msg': '操作失败、审核未通过'}
            # 其他情况
            else:
                context = {'status': 2, 'msg': '操作失败、请不要重复提交'}
        return context


class WorkOrderCloseForm(forms.Form):
    id = forms.IntegerField(required=True)
    status = forms.CharField(max_length=10, required=True)
    addition_info = forms.CharField(required=False)

    def is_save(self, request):
        cleaned_data = super(WorkOrderCloseForm, self).clean()
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
                        # 发送钉钉推送
                        xiaoding_pull.delay(user=request.user.username, id=id, type='close',
                                            addition_info=addition_info)
                        context = {'status': 0, 'msg': '操作成功、记录关闭成功'}

                elif status == u'结束':
                    context = {'status': 2, 'msg': '操作失败、关闭窗口'}
            else:
                context = {'status': 2, 'msg': '操作失败、<关闭原因>输入不能少于5个字符'}
        return context


class WorkOrderReplyForm(forms.Form):
    reply_id = forms.IntegerField(required=True)
    reply_contents = forms.CharField(widget=forms.Textarea, min_length=5)

    def is_save(self, request):
        cleaned_data = super(WorkOrderReplyForm, self).clean()
        reply_id = cleaned_data.get('reply_id')
        reply_contents = cleaned_data.get('reply_contents')
        OlAuditContentsReply.objects.create(
            reply_id=reply_id,
            user_id=request.user.uid,
            reply_contents=reply_contents)

        data = AuditContents.objects.get(pk=reply_id)
        data.updated_at = timezone.now()
        data.save()
        # 发送钉钉推送
        xiaoding_pull.delay(user=request.user.username, id=id, type='reply', addition_info=reply_contents)
        context = {'status': 0, 'msg': '回复成功'}
        return context

# -*- coding:utf-8 -*-
# edit by fuzongfei

import json
from datetime import datetime

from django.db import transaction
from django.db.models import Case, When, Value, CharField, F
from django.http import HttpResponse
from django.shortcuts import render, get_object_or_404
from django.utils import timezone
from django.utils.decorators import method_decorator
from django.views import View
from django.views.generic import ListView, FormView
from pure_pagination import PaginationMixin

from ProjectManager.inception.inception_api import IncepSqlCheck
from ProjectManager.models import OnlineAuditContents, OnlineAuditContentsReply, IncepMakeExecTask
from ProjectManager.ol.forms import IncepOlAuditForm, IncepOlReplyForm, IncepOlApproveForm
from ProjectManager.permissions import check_group_permission, check_sql_detail_permission
from ProjectManager.tasks import send_verify_mail, \
    send_reply_mail
from ProjectManager.utils import check_incep_alive, check_sql_filter
from utils.tools import format_request


class IncepOlAuditView(View):
    def get(self, request):
        """渲染线上审核页面"""
        return render(request, 'incep_ol_audit.html')

    @method_decorator(check_incep_alive)
    @method_decorator(check_sql_filter)
    @transaction.atomic
    def post(self, request):
        """线上审核内容提交处理"""
        data = format_request(request)
        form = IncepOlAuditForm(data)
        if form.is_valid():
            context = form.save(request)
            return HttpResponse(json.dumps(context))
        else:
            error = form.errors.as_text()
            context = {'status': 2, 'msg': error}
            return HttpResponse(json.dumps(context))


class IncepOlRecordsView(PaginationMixin, ListView):
    paginate_by = 8
    context_object_name = 'audit_records'
    template_name = 'incep_ol_records.html'

    obj = OnlineAuditContents.objects.all().annotate(
        progress_value=Case(
            When(progress_status='0', then=Value('待批准')),
            When(progress_status='1', then=Value('未批准')),
            When(progress_status='2', then=Value('已批准')),
            When(progress_status='3', then=Value('处理中')),
            When(progress_status='4', then=Value('已完成')),
            When(progress_status='5', then=Value('已关闭')),
            output_field=CharField(),
        ),
        progress_color=Case(
            When(progress_status__in=('0',), then=Value('btn-primary')),
            When(progress_status__in=('2',), then=Value('btn-warning')),
            When(progress_status__in=('1', '5'), then=Value('btn-danger')),
            When(progress_status__in=('3',), then=Value('btn-info')),
            When(progress_status__in=('4',), then=Value('btn-success')),
            output_field=CharField(),
        ),
        group_name=F('group__group_name'),
        group_id=F('group__group_id'),
    )

    def get_queryset(self):
        user_in_group = self.request.session['groups']
        search_content = self.request.GET.get('search_content')

        if search_content:
            audit_records = self.obj.filter(
                contents__contains=search_content
            ).filter(group_id__in=user_in_group). \
                values('group_name',
                       'progress_color',
                       'progress_value', 'id', 'group_id',
                       'title',
                       'proposer', 'operate_dba', 'verifier',
                       'created_at').order_by('-created_at')
        else:
            audit_records = self.obj.filter(group_id__in=user_in_group). \
                values('group_name', 'progress_color',
                       'progress_value', 'id', 'group_id',
                       'title',
                       'proposer', 'operate_dba', 'verifier',
                       'created_at').order_by('-created_at')

        return audit_records


class IncepOlApproveView(FormView):
    form_class = IncepOlApproveForm

    @method_decorator(check_group_permission)
    def dispatch(self, request, *args, **kwargs):
        return super(IncepOlApproveView, self).dispatch(request, *args, **kwargs)

    def form_valid(self, form):
        cleaned_data = form.cleaned_data
        id = cleaned_data.get('id')
        status = cleaned_data.get('status')
        addition_info = cleaned_data.get('addition_info')

        data = OnlineAuditContents.objects.get(pk=id)
        context = {}
        # 当记录关闭时
        if data.progress_status == '5':
            context = {'status': '400', 'msg': '该记录已被关闭、请不要重复提交'}
        # 当记录未关闭时
        else:
            # 角色为Leader的用户可以审批
            if self.request.user.user_role() == 'Leader':
                if data.progress_status == '0' or data.progress_status == '1':
                    # 当用户点击的是通过, 状态变为：已批准
                    with transaction.atomic():
                        if status == u'通过':
                            data.progress_status = '2'
                            data.fact_verifier = self.request.user.username
                            data.verifier_time = timezone.now()
                            data.save()
                            context = {'status': '200', 'msg': '操作成功、审核通过'}
                            send_verify_mail.delay(latest_id=id,
                                                   username=self.request.user.username,
                                                   user_role=self.request.user.user_role(),
                                                   addition_info=addition_info)
                        # 当用户点击的是不通过, 状态变为：未批准
                        elif status == u'不通过':
                            data.progress_status = '1'
                            data.fact_verifier = self.request.user.username
                            data.verifier_time = timezone.now()
                            data.save()
                            context = {'status': '200', 'msg': '操作成功、审核未通过'}
                            send_verify_mail.delay(latest_id=id,
                                                   username=self.request.user.username,
                                                   user_role=self.request.user.user_role(),
                                                   addition_info=addition_info)
                # 其他情况
                else:
                    context = {'status': '400', 'msg': '操作失败、请不要重复提交'}
            else:
                context = {'status': '403', 'msg': '权限拒绝, 您没有权限操作'}
        return HttpResponse(json.dumps(context))

    def form_invalid(self, form):
        error = form.errors.as_text()
        context = {'status': '400', 'msg': error}

        return HttpResponse(json.dumps(context))


class IncepOlFeedbackView(FormView):
    form_class = IncepOlApproveForm

    @method_decorator(check_group_permission)
    def dispatch(self, request, *args, **kwargs):
        return super(IncepOlFeedbackView, self).dispatch(request, *args, **kwargs)

    def form_valid(self, form):
        cleaned_data = form.cleaned_data
        id = cleaned_data.get('id')
        status = cleaned_data.get('status')
        addition_info = cleaned_data.get('addition_info')

        data = OnlineAuditContents.objects.get(pk=id)
        context = {}
        # 当记录关闭时
        if data.progress_status == '5':
            context = {'status': '400', 'msg': '该记录已被关闭、请不要重复提交'}
        # 当记录未关闭时
        else:
            # 角色为DBA的才能进行操作
            if self.request.user.user_role() == 'DBA':
                # 当进度状态为：已批准或处理中时
                if data.progress_status == '2' or data.progress_status == '3':
                    # 当用户点击的是处理中, 状态变为：处理中
                    with transaction.atomic():
                        if status == u'处理中':
                            data.progress_status = '3'
                            data.save()
                            context = {'status': '200', 'msg': '操作成功、正在处理中'}
                            send_verify_mail.delay(latest_id=id,
                                                   username=self.request.user.username,
                                                   user_role=self.request.user.user_role(),
                                                   addition_info=addition_info)
                        # 当用户点击的是已完成, 状态变为：已完成
                        elif status == u'已完成':
                            data.progress_status = '4'
                            data.fact_operate_dba = self.request.user.username
                            data.operate_time = timezone.now()
                            data.save()
                            context = {'status': '200', 'msg': '操作成功、处理完成'}
                            send_verify_mail.delay(latest_id=id,
                                                   username=self.request.user.username,
                                                   user_role=self.request.user.user_role(),
                                                   addition_info=addition_info)
                # 未批准
                elif data.progress_status == '1' or data.progress_status == '0':
                    context = {'status': '400', 'msg': '操作失败、审核未通过'}
                # 其他情况
                else:
                    context = {'status': '400', 'msg': '操作失败、请不要重复提交'}
            else:
                context = {'status': '403', 'msg': '权限拒绝、只有DBA角色可以操作'}
        return HttpResponse(json.dumps(context))

    def form_invalid(self, form):
        error = form.errors.as_text()
        context = {'status': '400', 'msg': error}

        return HttpResponse(json.dumps(context))


class IncepOlCloseView(FormView):
    form_class = IncepOlApproveForm

    @method_decorator(check_group_permission)
    def dispatch(self, request, *args, **kwargs):
        return super(IncepOlCloseView, self).dispatch(request, *args, **kwargs)

    def form_valid(self, form):
        cleaned_data = form.cleaned_data
        id = cleaned_data.get('id')
        status = cleaned_data.get('status')
        addition_info = cleaned_data.get('addition_info')

        data = OnlineAuditContents.objects.get(pk=id)
        context = {}
        # 当记录关闭时
        if data.progress_status == '5':
            context = {'status': '400', 'msg': '该记录已被关闭、请不要重复提交'}
        # 当记录未关闭时
        else:
            if len(addition_info) >= 5:
                # 当进度为：处理中或已完成时
                if status == u'提交':
                    if data.progress_status == '3' or data.progress_status == '4':
                        context = {'status': '400', 'msg': '操作失败、数据正在处理中或已完成'}
                    else:
                        with transaction.atomic():
                            data.progress_status = '5'
                            data.close_user = self.request.user.username
                            data.close_reason = addition_info
                            data.close_time = timezone.now()
                            data.save()
                            context = {'status': '200', 'msg': '操作成功、记录关闭成功'}
                            send_verify_mail.delay(latest_id=id,
                                                   username=self.request.user.username,
                                                   user_role=self.request.user.user_role(),
                                                   addition_info=addition_info)
                elif status == u'结束':
                    context = {'status': '400', 'msg': '操作失败、关闭窗口'}
            else:
                context = {'status': '400', 'msg': '操作失败、<关闭原因>不能少于5个字符'}

        return HttpResponse(json.dumps(context))

    def form_invalid(self, form):
        error = form.errors.as_text()
        context = {'status': '400', 'msg': error}

        return HttpResponse(json.dumps(context))


class IncepOlDetailsView(View):
    @method_decorator(check_sql_detail_permission)
    def get(self, request, id, group_id):
        obj = OnlineAuditContents.objects.annotate(
            progress_value=Case(
                When(progress_status='0', then=Value('待批准')),
                When(progress_status='1', then=Value('未批准')),
                When(progress_status='2', then=Value('已批准')),
                When(progress_status='3', then=Value('处理中')),
                When(progress_status='4', then=Value('已完成')),
                When(progress_status='5', then=Value('已关闭')),
                output_field=CharField(),
            )
        )

        contents = obj.get(id=id)
        group = OnlineAuditContents.objects.filter(id=id).annotate(group_name=F('group__group_name')).values(
            'group_name').first()
        reply_contents = OnlineAuditContentsReply.objects.annotate(
            username=F('user__username'),
            avatar_file=F('user__avatar_file'),
        ).filter(reply__id=id).values('username', 'avatar_file', 'reply_contents', 'created_at')
        return render(request, 'incep_ol_details.html',
                      {'contents': contents, 'group': group, 'reply_contents': reply_contents})


class IncepOlReplyView(FormView):
    """处理用户的回复信息"""

    form_class = IncepOlReplyForm

    def form_valid(self, form):
        cleaned_data = form.cleaned_data
        reply_id = cleaned_data['reply_id']
        reply_contents = cleaned_data['reply_contents']
        OnlineAuditContentsReply.objects.create(
            reply_id=reply_id,
            user_id=self.request.user.uid,
            reply_contents=reply_contents)
        context = {'status': '200', 'msg': '回复成功'}
        latest_id = OnlineAuditContentsReply.objects.latest('id').id
        send_reply_mail.delay(latest_id=latest_id,
                              reply_id=reply_id,
                              username=self.request.user.username,
                              email=self.request.user.email)
        return HttpResponse(json.dumps(context))

    def form_invalid(self, form):
        error = form.errors.as_text()
        context = {'status': '400', 'msg': error}

        return HttpResponse(json.dumps(context))


class IncepGenerateTasksView(View):
    @method_decorator(check_group_permission)
    def post(self, request):
        id = request.POST.get('id')

        if IncepMakeExecTask.objects.filter(related_id=id).first():
            taskid = IncepMakeExecTask.objects.filter(related_id=id).first().taskid
            context = {'status': 201,
                       'jump_url': f'/projects/incep_perform_records/incep_perform_details/{taskid}'}
        else:
            obj = get_object_or_404(OnlineAuditContents, pk=id)

            # 只要leader批准后，才能执行生成执行任务
            if obj.progress_status in ('2', '3', '4', '5'):
                host = obj.dst_host
                database = obj.dst_database
                sql_content = obj.contents

                # 实例化
                incep_of_audit = IncepSqlCheck(sql_content, host, database, request.user.username)

                # 对OSC执行的SQL生成sqlsha1

                result = incep_of_audit.make_sqlsha1()
                taskid = datetime.now().strftime("%Y%m%d%H%M%S%f")
                # 生成执行任务记录
                for row in result:
                    IncepMakeExecTask.objects.create(
                        uid=request.user.uid,
                        user=obj.proposer,
                        taskid=taskid,
                        dst_host=host,
                        dst_database=database,
                        sql_content=row['SQL'],
                        sqlsha1=row['sqlsha1'],
                        affected_rows=row['Affected_rows'],
                        type=obj.type,
                        category='1',
                        related_id=id,
                        group_id=obj.group_id
                    )

                context = {'status': 201,
                           'jump_url': f'/projects/incep_perform_records/incep_perform_details/{taskid}'}
            else:
                context = {'status': 400,
                           'msg': 'Leader审核未通过'}

        return HttpResponse(json.dumps(context))

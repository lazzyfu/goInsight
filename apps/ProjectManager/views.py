import json
import re
from datetime import datetime

import sqlparse
from django.core.exceptions import PermissionDenied
from django.db import transaction
from django.db.models import F, When, Value, CharField, Case
from django.http import JsonResponse, HttpResponse
from django.shortcuts import render, get_object_or_404
from django.utils import timezone
from django.utils.decorators import method_decorator
from django.views import View
from django.views.generic import FormView, ListView
from pure_pagination import PaginationMixin

from ProjectManager.forms import InceptionSqlOperateForm, OnlineAuditCommitForm, VerifyCommitForm, ReplyContentForm
from ProjectManager.group_permissions import check_group_permission, check_sql_detail_permission
from UserManager.models import GroupsDetail, UserAccount, Contacts
from apps.ProjectManager.inception.inception_api import GetDatabaseListApi, InceptionApi, GetBackupApi, IncepSqlOperate
from utils.tools import format_request
from .models import InceptionHostConfig, InceptionSqlOperateRecord, Remark, OnlineAuditContents, \
    OnlineAuditContentsReply


class IncepOfflineSqlCheckView(FormView):
    form_class = InceptionSqlOperateForm
    template_name = 'incep_offline_sql_check.html'

    def form_valid(self, form):
        cleaned_data = form.cleaned_data
        host = cleaned_data['host']
        database = cleaned_data['database']
        op_action = cleaned_data.get('op_action')
        op_type = cleaned_data['op_type']
        sql_content = cleaned_data['sql_content']

        context = {}
        if op_type == 'check':
            context = IncepSqlOperate(sql_content, host, database).run_check(op_action)

        elif op_type == 'execute':
            context = IncepSqlOperate(sql_content, host, database).run_execute(op_action, form, self.request)
        return HttpResponse(json.dumps(context))

    def form_invalid(self, form):
        error = "请选择主机或库名"
        context = {'errCode': '400', 'errMsg': error}

        return HttpResponse(json.dumps(context))


class BeautifySQLView(View):
    """美化SQL"""

    def post(self, request):
        data = format_request(request)
        sqlContent = data['sql_content'].rstrip()
        sqlFormat = '\n'.join([line for line in sqlContent.split('\n') if line != ''])
        beautifySQL = sqlparse.format(sqlFormat, keyword_case='upper')
        context = {'data': beautifySQL}
        return HttpResponse(json.dumps(context))


class GetInceptionHostConfigView(View):
    """获取inception指定的目标数据库配置"""

    def get(self, request):
        type = request.GET.get('type')
        envResult = InceptionHostConfig.objects.filter(type=type).values('host', 'comment')
        return JsonResponse(list(envResult), safe=False)


class GetDatabaseListView(View):
    """列出选中环境的数据库库名"""

    def post(self, request):
        data = format_request(request)
        host = data['host']
        dbList = GetDatabaseListApi(host).get_dbname()
        return HttpResponse(json.dumps(dbList))


class IncepOfflineSqlRecords(PaginationMixin, ListView):
    """查看用户的工单记录"""
    paginate_by = 8
    context_object_name = 'sqlRecord'
    template_name = 'inception_sql_records.html'

    def get_queryset(self):
        workidQuery = "select workid,id,op_user,dst_host,op_time from sqlaudit_inception_sql_operate_record group by workid order by op_time desc"
        sqlRecord = []
        for row in InceptionSqlOperateRecord.objects.raw(workidQuery):
            workid = row.workid
            op_user = row.op_user
            dst_host = row.dst_host
            op_time = row.op_time
            singleRecord = InceptionSqlOperateRecord.objects.filter(op_uid=self.request.user.uid).filter(
                workid=workid).order_by(
                'op_time')
            sqlRecord.append({'workid': workid, 'op_user': op_user, 'dst_host': dst_host, 'op_time': op_time,
                              'record': singleRecord})
        return sqlRecord


class IncepOfflineAllSqlDetailView(View):
    """查看当前用户会话执行的所有sql的详情"""

    def get(self, request, workid):
        sequenceResult = []
        originalSql = ''
        originalSqlQuery = InceptionSqlOperateRecord.objects.raw(
            f"select id,group_concat(`op_sql` separator '\n') as `op_sql` from sqlaudit_inception_sql_operate_record where workid={workid} group by workid")
        for i in originalSqlQuery:
            originalSql = i.op_sql

        sqlDetail = InceptionSqlOperateRecord.objects.filter(workid=workid)
        for row in sqlDetail:
            sequenceResult.append({'backupdbName': row.backup_dbname, 'sequence': row.sequence})
        rollbackSql = GetBackupApi(sequenceResult).get_backupinfo()

        return render(request, 'allsql_detail.html',
                      {'originalSql': originalSql, 'rollbackSql': rollbackSql})


class IncepOfflineSingleSqlDetailView(View):
    """查看当前用户会话执行的每条sql的详情"""

    def get(self, request, sequence):
        sqlDetail = get_object_or_404(InceptionSqlOperateRecord, sequence=sequence)
        sequenceResult = [{'backupdbName': sqlDetail.backup_dbname, 'sequence': sqlDetail.sequence}]
        rollbackSql = GetBackupApi(sequenceResult).get_backupinfo()
        return render(request, 'singlesql_detail.html',
                      {'sqlDetail': sqlDetail, 'rollbackSql': rollbackSql})


class IncepOnlineSqlCheckView(FormView):
    """
    处理用户提交的审核内容
    """
    form_class = OnlineAuditCommitForm
    template_name = 'incep_online_sql_commit.html'

    def form_valid(self, form):
        cleaned_data = form.cleaned_data
        title = cleaned_data.get('title') + '__[' + datetime.now().strftime("%Y%m%d%H%M%S") + ']'
        remark = cleaned_data.get('remark')
        verifier = cleaned_data.get('verifier')
        operate_dba = cleaned_data.get('operate_dba')
        group_id = cleaned_data.get('group_id')
        email_cc = cleaned_data.get('email_cc_id')
        host = cleaned_data['host']
        database = cleaned_data['database']
        op_action = cleaned_data.get('op_action')
        sql_content = cleaned_data['sql_content']

        result = IncepSqlOperate(sql_content, host, database).check_valid(op_action, form, self.request)
        if result.get('errCode') == 400:
            context = result
        elif result.get('errCode') == 200:
            with transaction.atomic():
                OnlineAuditContents.objects.create(
                    title=title,
                    group_id=group_id,
                    remark=remark,
                    proposer=self.request.user.username,
                    operate_dba=operate_dba,
                    verifier=verifier,
                    email_cc=email_cc,
                    contents=sql_content
                )

                # 发送通知邮件
                # latest_id = AuditContents.objects.latest('id').id
                #                 send_commit_mail.delay(latest_id=latest_id)
            context = {'errCode': '200', 'errMsg': '提交成功'}
        return HttpResponse(json.dumps(context))

    def form_invalid(self, form):
        error = form.errors.as_text()
        context = {'errCode': '400', 'errMsg': error}
        return HttpResponse(json.dumps(context))


class GetRemarkInfo(View):
    def post(self, request):
        obj = Remark.objects.all().values('id', 'remark')
        return JsonResponse(list(obj), safe=False)


class GetGroupView(View):
    def get(self, request):
        groups = GroupsDetail.objects.filter(
            user__uid=request.user.uid).annotate(
            group_id=F('group__group_id'), group_name=F('group__group_name')) \
            .values('group_id', 'group_name')

        return JsonResponse(list(groups), safe=False)


class GetDbaLeaderView(View):
    def post(self, request):
        """
        获取指定项目可用的dba和leader信息
        """
        group_id = request.POST.get('group_id')
        result = []
        if group_id:
            data = GroupsDetail.objects.annotate(
                uid=F('user__uid'),
                username=F('user__username'),
                email=F('user__email'),
            ).filter(group__group_id=group_id).values('uid', 'username', 'email')

            for i in data:
                uid = i['uid']
                user_role = UserAccount.objects.get(uid=uid).user_role()
                i['user_role'] = user_role
                result.append(i)

        return JsonResponse(result, safe=False)


class GetContactsView(View):
    def post(self, request):
        """ 获取指定项目的联系人"""
        group_id = request.POST.get('group_id')

        result = []
        if group_id:

            query = f"select ac.contact_id, group_concat(concat_ws(':', ac.contact_name, ac.contact_id, ac.contact_email)) as contact_info " \
                    f"from auditsql_contacts as ac JOIN auditsql_contacts_detail a ON ac.contact_id = a.contact_id JOIN  auditsql_groups a2 " \
                    f"ON a.group_id = a2.group_id where a.group_id = {group_id} group by ac.contact_id;"

            for row in Contacts.objects.raw(query):
                result.append(row.contact_info)

        return JsonResponse(result, safe=False)


class IncepOnlineAuditRecordsView(PaginationMixin, ListView):
    paginate_by = 8
    context_object_name = 'audit_records'
    template_name = 'incep_online_sql_records.html'

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
        # 返回用户所在项目组的记录
        # 如果用户不属于任何项目组, 返回没有权限查询记录
        user_in_group = self.request.session['groups']
        if len(user_in_group) == 0:
            raise PermissionDenied
        else:
            contents = self.request.GET.get('search_contents')

            if contents:
                audit_records = self.obj.filter(
                    contents__contains=contents
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


class IncepOnlineClickVerifyView(FormView):
    form_class = VerifyCommitForm

    @method_decorator(check_group_permission)
    def dispatch(self, request, *args, **kwargs):
        return super(IncepOnlineClickVerifyView, self).dispatch(request, *args, **kwargs)

    def form_valid(self, form):
        cleaned_data = form.cleaned_data
        id = cleaned_data.get('id')
        status = cleaned_data.get('status')
        addition_info = cleaned_data.get('addition_info')

        data = OnlineAuditContents.objects.get(pk=id)
        # 当记录关闭时
        if data.progress_status == '5':
            context = {'errCode': '400', 'errMsg': '该记录已被关闭、请不要重复提交'}
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
                            context = {'errCode': '200', 'errMsg': '操作成功、审核通过'}
                            # send_verify_mail.delay(latest_id=id, username=request.user.username,
                            #                        user_role=request.user.user_role())
                        # 当用户点击的是不通过, 状态变为：未批准
                        elif status == u'不通过':
                            data.progress_status = '1'
                            data.fact_verifier = self.request.user.username
                            data.verifier_time = timezone.now()
                            data.save()
                            context = {'errCode': '200', 'errMsg': '操作成功、审核未通过'}
                            # send_verify_mail.delay(latest_id=id, username=request.user.username,
                            #                        user_role=request.user.user_role())
                # 其他情况
                else:
                    context = {'errCode': '400', 'errMsg': '操作失败、请不要重复提交'}
            else:
                context = {'errCode': '403', 'errMsg': '权限拒绝, 您没有权限操作'}
        return HttpResponse(json.dumps(context))

    def form_invalid(self, form):
        error = form.errors.as_text()
        context = {'errCode': '400', 'errMsg': error}

        return HttpResponse(json.dumps(context))


class IncepOnlineClickFinishView(FormView):
    form_class = VerifyCommitForm

    @method_decorator(check_group_permission)
    def dispatch(self, request, *args, **kwargs):
        return super(IncepOnlineClickFinishView, self).dispatch(request, *args, **kwargs)

    def form_valid(self, form):
        cleaned_data = form.cleaned_data
        id = cleaned_data.get('id')
        status = cleaned_data.get('status')
        addition_info = cleaned_data.get('addition_info')

        data = OnlineAuditContents.objects.get(pk=id)
        # 当记录关闭时
        if data.progress_status == '5':
            context = {'errCode': '400', 'errMsg': '该记录已被关闭、请不要重复提交'}
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
                            context = {'errCode': '200', 'errMsg': '操作成功、正在处理中'}
                            # send_verify_mail.delay(latest_id=id, username=request.user.username,
                            #                        user_role=request.user.user_role())
                        # 当用户点击的是已完成, 状态变为：已完成
                        elif status == u'已完成':
                            data.progress_status = '4'
                            data.fact_operate_dba = self.request.user.username
                            data.operate_time = timezone.now()
                            data.save()
                            context = {'errCode': '200', 'errMsg': '操作成功、处理完成'}
                            # send_verify_mail.delay(latest_id=id, username=request.user.username,
                            #                        user_role=request.user.user_role())
                # 未批准
                elif data.progress_status == '1' or data.progress_status == '0':
                    context = {'errCode': '400', 'errMsg': '操作失败、审核未通过'}
                # 其他情况
                else:
                    context = {'errCode': '400', 'errMsg': '操作失败、请不要重复提交'}
            else:
                context = {'errCode': '403', 'errMsg': '权限拒绝、只有DBA角色可以操作'}
        return HttpResponse(json.dumps(context))

    def form_invalid(self, form):
        error = form.errors.as_text()
        context = {'errCode': '400', 'errMsg': error}

        return HttpResponse(json.dumps(context))


class IncepOnlineClickCloseView(FormView):
    form_class = VerifyCommitForm

    @method_decorator(check_group_permission)
    def dispatch(self, request, *args, **kwargs):
        return super(IncepOnlineClickCloseView, self).dispatch(request, *args, **kwargs)

    def form_valid(self, form):
        cleaned_data = form.cleaned_data
        id = cleaned_data.get('id')
        status = cleaned_data.get('status')
        addition_info = cleaned_data.get('addition_info')

        data = OnlineAuditContents.objects.get(pk=id)
        # 当记录关闭时
        if data.progress_status == '5':
            context = {'errCode': '400', 'errMsg': '该记录已被关闭、请不要重复提交'}
        # 当记录未关闭时
        else:
            if len(addition_info) >= 5:
                # 当进度为：处理中或已完成时
                if status == u'提交':
                    if data.progress_status == '3' or data.progress_status == '4':
                        context = {'errCode': '400', 'errMsg': '操作失败、数据正在处理中或已完成'}
                    else:
                        with transaction.atomic():
                            data.progress_status = '5'
                            data.close_user = self.request.user.username
                            data.close_reason = addition_info
                            data.close_time = timezone.now()
                            data.save()
                            context = {'errCode': '200', 'errMsg': '操作成功、记录关闭成功'}
                            # send_verify_mail.delay(latest_id=id, username=request.user.username,
                            #                        user_role=request.user.user_role())
                elif status == u'结束':
                    context = {'errCode': '400', 'errMsg': '操作失败、关闭窗口'}
            else:
                context = {'errCode': '400', 'errMsg': '操作失败、<关闭原因>不能少于5个字符'}

        return HttpResponse(json.dumps(context))

    def form_invalid(self, form):
        error = form.errors.as_text()
        context = {'errCode': '400', 'errMsg': error}

        return HttpResponse(json.dumps(context))


class OnlineAuditDetailView(View):
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
        return render(request, 'incep_online_sql_detail.html',
                      {'contents': contents, 'group': group, 'reply_contents': reply_contents})


class OnlineSqlReplyView(FormView):
    """处理用户的回复信息"""

    form_class = ReplyContentForm

    def form_valid(self, form):
        cleaned_data = form.cleaned_data
        reply_id = cleaned_data['reply_id']
        reply_contents = cleaned_data['reply_contents']
        OnlineAuditContentsReply.objects.create(
            reply_id=reply_id,
            user_id=self.request.user.uid,
            reply_contents=reply_contents)
        context = {'status': '200', 'msg': '回复成功'}
        return HttpResponse(json.dumps(context))
        latest_id = OnlineAuditContentsReply.objects.latest('id').id
        # send_reply_mail.delay(latest_id=latest_id, reply_id=reply_id, username=request.user.username,
        #                       email=request.user.email)

    def form_invalid(self, form):
        error = form.errors.as_text()
        context = {'status': '400', 'msg': error}

        return HttpResponse(json.dumps(context))

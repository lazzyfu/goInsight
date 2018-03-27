import json
from ast import literal_eval
from datetime import datetime

import sqlparse
from channels.layers import get_channel_layer
from django.core.cache import cache
from django.db import transaction
from django.db.models import F, When, Value, CharField, Case, Count
from django.http import JsonResponse, HttpResponse
from django.shortcuts import render, get_object_or_404
from django.utils import timezone
from django.utils.decorators import method_decorator
from django.views import View
from django.views.generic import FormView, ListView
from pure_pagination import PaginationMixin

from ProjectManager.forms import InceptionSqlCheckForm, OnlineAuditCommitForm, VerifyCommitForm, ReplyContentForm
from ProjectManager.permissions import check_group_permission, check_sql_detail_permission, check_incep_tasks_permission
from ProjectManager.utils import update_tasks_status, check_incep_alive
from UserManager.models import GroupsDetail, UserAccount, Contacts
from apps.ProjectManager.inception.inception_api import GetDatabaseListApi, GetBackupApi, IncepSqlCheck, \
    sql_filter
from utils.tools import format_request
from .models import Remark, OnlineAuditContents, \
    OnlineAuditContentsReply, InceptionHostConfigDetail, IncepMakeExecTask
from .tasks import send_commit_mail, send_verify_mail, send_reply_mail, get_osc_percent, incep_async_tasks, \
    stop_incep_osc

channel_layer = get_channel_layer()


class IncepOfAuditView(View):
    def get(self, request):
        return render(request, 'incep_of_audit.html')

    @method_decorator(check_incep_alive)
    def post(self, request):
        form = InceptionSqlCheckForm(request.POST)
        if form.is_valid():
            cleaned_data = form.cleaned_data
            host = cleaned_data['host']
            database = cleaned_data['database']
            op_action = cleaned_data.get('op_action')
            op_type = cleaned_data['op_type']
            group_id = cleaned_data['group_id']
            sql_content = cleaned_data['sql_content']

            # 对检测的SQL类型进行区分
            filter_result = sql_filter(sql_content, op_action)

            # 实例化
            incep_of_audit = IncepSqlCheck(sql_content, host, database, self.request.user.username)

            if filter_result['errCode'] == 400:
                context = filter_result
            else:
                # SQL语法检查
                if op_type == 'check':
                    context = incep_of_audit.run_check()

                # 生成执行任务
                elif op_type == 'make':
                    # 生成执行任务之前，检测是否审核通过
                    check_result = incep_of_audit.is_check_pass()
                    if check_result['errCode'] == 400:
                        context = check_result
                    else:
                        # 对OSC执行的SQL生成sqlsha1
                        result = incep_of_audit.make_sqlsha1()
                        taskid = datetime.now().strftime("%Y%m%d%H%M%S%f")
                        # 生成执行任务记录
                        for row in result:
                            IncepMakeExecTask.objects.create(
                                uid=self.request.user.uid,
                                user=self.request.user.username,
                                taskid=taskid,
                                dst_host=host,
                                group_id=group_id,
                                dst_database=database,
                                sql_content=row['SQL'],
                                sqlsha1=row['sqlsha1'],
                                type=filter_result['type']
                            )
                        context = {'errCode': 201,
                                   'dst_url': f'/projects/incep_of_records/incep_of_details/{taskid}'}
            return HttpResponse(json.dumps(context))
        else:
            error = "请选择主机、库名和项目组"
            context = {'errCode': 400, 'errMsg': error}

            return HttpResponse(json.dumps(context))


class BeautifySQLView(View):
    """
    美化SQL
    判断SQL类型（DML还是DDL），并分别进行美化
    最后合并返回
    """

    def post(self, request):
        data = format_request(request)
        sqlContent = data.get('sql_content').strip()

        sqlSplit = []
        for stmt in sqlparse.split(sqlContent):
            sql = sqlparse.parse(stmt)[0]
            sql_comment = sql.token_first()
            if isinstance(sql_comment, sqlparse.sql.Comment):
                sqlSplit.append({'comment': sql_comment.value, 'sql': sql.value.replace(sql_comment.value, '')})
            else:
                sqlSplit.append({'comment': '', 'sql': sql.value})

        beautifySQL_list = []
        try:
            for row in sqlSplit:
                comment = row['comment']
                sql = row['sql']
                res = sqlparse.parse(sql)
                if res[0].tokens[0].ttype[1] == 'DML':
                    sqlFormat = sqlparse.format(sql, reindent=True)
                    beautifySQL_list.append(comment + sqlFormat)
                elif res[0].tokens[0].ttype[1] == 'DDL':
                    sqlFormat = sqlparse.format(sql)
                    beautifySQL_list.append(comment + sqlFormat)
            beautifySQL = '\n\n'.join(beautifySQL_list)
            context = {'data': beautifySQL}
        except Exception as err:
            raise OSError(err)
            context = {'errCode': 400, 'errMsg': "注释不合法, 请检查"}

        return HttpResponse(json.dumps(context))


class IncepHostConfigView(View):
    """获取inception指定的目标数据库配置"""

    def get(self, request):
        type = request.GET.get('type')
        user_in_group = self.request.session.get('groups')
        envResult = InceptionHostConfigDetail.objects.annotate(host=F('config__host'),
                                                               comment=F('config__comment')).filter(
            config__type=type).filter(group__group_id__in=user_in_group).values('host', 'comment')
        return JsonResponse(list(envResult), safe=False)


class GetDBListView(View):
    """列出选中环境的数据库库名"""

    def post(self, request):
        data = format_request(request)
        host = data['host']
        dbList = GetDatabaseListApi(host).get_dbname()
        return HttpResponse(json.dumps(dbList))


class IncepOfResultsView(View):
    def get(self, request):
        id = request.GET.get('id')
        if IncepMakeExecTask.objects.get(id=id).exec_status in ('1', '4'):
            sqlDetail = IncepMakeExecTask.objects.get(id=id)
            sequenceResult = {'backupdbName': sqlDetail.backup_dbname, 'sequence': sqlDetail.sequence}
            rollback_sql = GetBackupApi(sequenceResult).get_rollback_statement()

            exec_log = sqlDetail.exec_log if sqlDetail.exec_log else '无记录'

            # 此处要将exec_log去字符串处理，否则无法转换为json
            context = {'rollback_log': rollback_sql, 'exec_log': literal_eval(exec_log), 'errCode': 200}
        else:
            context = {'errCode': 400, 'errMsg': '该SQL未被执行，无法查询状态信息'}

        return HttpResponse(json.dumps(context))


class IncepOlAuditView(View):
    def get(self, request):
        return render(request, 'incep_ol_audit.html')

    @method_decorator(check_incep_alive)
    def post(self, request):
        form = OnlineAuditCommitForm(request.POST)
        if form.is_valid():
            cleaned_data = form.cleaned_data
            title = cleaned_data.get('title') + '__[' + datetime.now().strftime("%Y%m%d%H%M%S") + ']'
            remark = cleaned_data.get('remark')
            verifier = cleaned_data.get('verifier')
            operate_dba = cleaned_data.get('operate_dba')
            group_id = cleaned_data.get('group_id')
            email_cc = ','.join(self.request.POST.getlist('email_cc_id'))
            host = cleaned_data['host']
            database = cleaned_data['database']
            op_action = cleaned_data.get('op_action')
            sql_content = cleaned_data['sql_content']

            result = IncepSqlCheck(sql_content, host, database, self.request.user.username).is_check_pass()
            if result.get('errCode') == 400:
                context = result
            elif result.get('errCode') == 200:
                with transaction.atomic():
                    OnlineAuditContents.objects.create(
                        title=title,
                        op_action='数据修改' if op_action == 'op_data' else '表结构变更',
                        type='DML' if op_action == 'op_data' else 'DDL',
                        dst_host=host,
                        dst_database=database,
                        group_id=group_id,
                        remark=remark,
                        proposer=self.request.user.username,
                        operate_dba=operate_dba,
                        verifier=verifier,
                        email_cc=email_cc,
                        contents=sql_content
                    )

                # 发送通知邮件
                latest_id = OnlineAuditContents.objects.latest('id').id
                send_commit_mail.delay(latest_id=latest_id)
                context = {'errCode': '200', 'errMsg': '提交成功, 跳转到工单页面'}
            return HttpResponse(json.dumps(context))
        else:
            error = form.errors.as_text()
            context = {'errCode': '400', 'errMsg': error}
            return HttpResponse(json.dumps(context))


class RemarkInfoView(View):
    def post(self, request):
        obj = Remark.objects.all().values('id', 'remark')
        return JsonResponse(list(obj), safe=False)


class GroupInfoView(View):
    def get(self, request):
        groups = GroupsDetail.objects.filter(
            user__uid=request.user.uid).annotate(
            group_id=F('group__group_id'), group_name=F('group__group_name')) \
            .values('group_id', 'group_name')

        return JsonResponse(list(groups), safe=False)


class AuditUserView(View):
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


class ContactsInfoView(View):
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
    form_class = VerifyCommitForm

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
                            context = {'errCode': '200', 'errMsg': '操作成功、审核未通过'}
                            send_verify_mail.delay(latest_id=id,
                                                   username=self.request.user.username,
                                                   user_role=self.request.user.user_role(),
                                                   addition_info=addition_info)
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


class IncepOlFeedbackView(FormView):
    form_class = VerifyCommitForm

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
                            context = {'errCode': '200', 'errMsg': '操作成功、处理完成'}
                            send_verify_mail.delay(latest_id=id,
                                                   username=self.request.user.username,
                                                   user_role=self.request.user.user_role(),
                                                   addition_info=addition_info)
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


class IncepOlCloseView(FormView):
    form_class = VerifyCommitForm

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
                            send_verify_mail.delay(latest_id=id,
                                                   username=self.request.user.username,
                                                   user_role=self.request.user.user_role(),
                                                   addition_info=addition_info)
                elif status == u'结束':
                    context = {'errCode': '400', 'errMsg': '操作失败、关闭窗口'}
            else:
                context = {'errCode': '400', 'errMsg': '操作失败、<关闭原因>不能少于5个字符'}

        return HttpResponse(json.dumps(context))

    def form_invalid(self, form):
        error = form.errors.as_text()
        context = {'errCode': '400', 'errMsg': error}

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


class IncepOfRecordsView(View):
    def get(self, request):
        return render(request, 'incep_of_records.html')


class IncepOfRecordsListView(View):
    def get(self, request):
        exec_tasks = []
        user_in_group = '(' + str(request.session['groups'][0]) + ')' if len(request.session['groups']) == 1 else tuple(
            request.session['groups'])
        query = f"select a.id,a.user,a.taskid,a.dst_host,a.dst_database,a.make_time, b.group_name," \
                f"case a.category when '0' then '线下任务' when '1' then '线上任务' end as category " \
                f"from sqlaudit_incep_tasks as a join auditsql_groups as b " \
                f"on a.group_id = b.group_id where b.group_id in {user_in_group} group by a.taskid order by a.make_time  desc"
        for row in IncepMakeExecTask.objects.raw(query):
            exec_tasks.append({'user': row.user,
                               'taskid': row.taskid,
                               'group_name': row.group_name,
                               'category': row.category,
                               'dst_host': row.dst_host,
                               'dst_database': row.dst_database,
                               'make_time': row.make_time})
        return JsonResponse(list(exec_tasks), safe=False)


class IncepOfDetailsView(View):
    def get(self, request, taskid):
        return render(request, 'incep_of_details.html', {'taskid': taskid})


class IncepOfDetailsListView(View):
    def get(self, request):
        taskid = request.GET.get('taskid')

        query = "select id,user,sqlsha1,sql_content,taskid,case exec_status " \
                "when '0' then '未执行' when '1' then '已完成' when '2' then '处理中' when '3' then '回滚中' " \
                "when '4' then '已回滚' end as exec_status," \
                "case category when '0' then '线下任务' when '1' then '线上任务' end as category" \
                " from sqlaudit_incep_tasks where taskid={taskid}".format(taskid=taskid)
        i = 0
        task_details = []
        for row in IncepMakeExecTask.objects.raw(query):
            task_details.append({
                'sid': i,
                'id': row.id,
                'user': row.user,
                'category': row.category,
                'sqlsha1': row.sqlsha1,
                'sql_content': row.sql_content,
                'taskid': row.taskid,
                'exec_status': row.exec_status
            })
            i += 1
        del task_details[0]
        return HttpResponse(json.dumps(task_details))


class IncepPerformView(View):
    """
    执行任务
    """

    @method_decorator(check_incep_alive)
    @method_decorator(check_incep_tasks_permission)
    def post(self, request):
        id = request.POST.get('id')
        obj = IncepMakeExecTask.objects.get(id=id)

        query = f"select id,group_concat(exec_status) as exec_status from sqlaudit_incep_tasks where taskid={obj.taskid} group by taskid"
        for row in IncepMakeExecTask.objects.raw(query):
            status = row.exec_status.split(',')

        # 每次只能执行一条任务，不可同时执行，避免数据库压力
        key = '-'.join(('django', str(request.user.uid), obj.sqlsha1))
        if '2' in status or '3' in status:
            context = {'errCode': 400, 'errMsg': '请等待当前其他任务执行完成'}
        else:
            # 避免任务重复点击执行
            if obj.exec_status != '0':
                context = {'errCode': 400, 'errMsg': '请不要重复操作任务'}
            else:
                # 如果sqlsha1存在，使用OSC执行
                if obj.sqlsha1:
                    # 在redis里面存储key，用于celery后台线程通信
                    cache.set(key, 'start', timeout=None)

                    # 将任务进度设置为：处理中
                    obj.exec_status = 2
                    obj.save()

                    # 执行异步线程
                    # 执行SQL任务
                    incep_async_tasks.delay(user=request.user.username,
                                            redis_key=key,
                                            id=id,
                                            exec_status=1)
                    # 执行获取进度任务
                    get_osc_percent.delay(user=request.user.username,
                                          id=id,
                                          redis_key=key)

                    context = {'errCode': 200, 'errMsg': '提交处理，请查看输出'}

                # 如果sqlsha1不存在，直接执行
                else:
                    incep_of_audit = IncepSqlCheck(obj.sql_content + ';',
                                                   obj.dst_host,
                                                   obj.dst_database,
                                                   request.user.username)
                    exec_result = incep_of_audit.run_exec(0)
                    # 更新任务状态
                    update_tasks_status(id=id,
                                        exec_result=exec_result,
                                        exec_status=1)

                    context = {'errCode': 200, 'errMsg': '提交处理，请查看输出'}
        return HttpResponse(json.dumps(context))


class IncepStopView(View):
    """
    停止OSC执行
    """

    @method_decorator(check_incep_alive)
    @method_decorator(check_incep_tasks_permission)
    def post(self, request):
        id = request.POST.get('id')

        obj = IncepMakeExecTask.objects.get(id=id)
        key = '-'.join(('django', str(request.user.uid), obj.sqlsha1))

        if obj.exec_status in ('0', '1', '4'):
            context = {'errCode': 400, 'errMsg': '请不要重复操作任务'}
        else:
            # 关闭正在执行的任务
            stop_incep_osc.delay(user=request.user.username,
                                 redis_key=key,
                                 id=id)
            context = {'errCode': 200, 'errMsg': '提交处理，请查看输出'}
        return HttpResponse(json.dumps(context))


class IncepRollbackView(View):
    """
    回滚操作
    """

    @method_decorator(check_incep_alive)
    @method_decorator(check_incep_tasks_permission)
    def post(self, request):
        id = request.POST.get('id')

        obj = IncepMakeExecTask.objects.get(id=id)
        context = {}

        if obj.exec_status in ('0', '3', '4'):
            context = {'errCode': 400, 'errMsg': '请不要重复操作'}
        else:
            # 获取回滚语句
            rollback_sql = GetBackupApi(
                {'backupdbName': obj.backup_dbname, 'sequence': obj.sequence}).get_rollback_statement()
            if rollback_sql == u'无记录':
                context = {'errCode': 400, 'errMsg': '没有找到备份记录，回滚失败'}
            else:
                incep_of_audit = IncepSqlCheck(rollback_sql, obj.dst_host, obj.dst_database, request.user.username)
                result = incep_of_audit.make_sqlsha1()[1]

                if obj.type == 'DML':
                    incep_of_audit = IncepSqlCheck(rollback_sql, obj.dst_host, obj.dst_database, request.user.username)
                    exec_result = incep_of_audit.run_exec(0)
                    # 更新任务状态
                    update_tasks_status(id=id, exec_result=exec_result, exec_status=4)

                    context = {'errCode': 200, 'errMsg': '提交处理，请查看输出'}

                elif obj.type == 'DDL':
                    if result['sqlsha1']:
                        key = '-'.join(('django', str(request.user.uid), result['sqlsha1']))
                        # 在redis里面存储key，用于celery后台线程通信
                        cache.set(key, 'start', timeout=None)

                        # 将任务进度设置为：处理中
                        obj.exec_status = 3
                        obj.rollback_sqlsha1 = result['sqlsha1']
                        obj.save()

                        # 执行SQL任务
                        incep_async_tasks.delay(user=request.user.username,
                                                redis_key=key,
                                                sql=result['SQL'] + ';',
                                                id=id,
                                                exec_status=4)

                        # 执行获取进度任务
                        get_osc_percent.delay(user=request.user.username,
                                              sqlsha1=result['sqlsha1'],
                                              id=id,
                                              redis_key=key)
                        context = {'errCode': 200, 'errMsg': '提交处理，请查看输出'}
                    else:
                        incep_of_audit = IncepSqlCheck(result['SQL'] + ';', obj.dst_host, obj.dst_database,
                                                       request.user.username)
                        exec_result = incep_of_audit.run_exec(0)
                        # 更新任务状态
                        update_tasks_status(id=id, exec_result=exec_result, exec_status=4)

                        context = {'errCode': 200, 'errMsg': '提交处理，请查看输出'}
        return HttpResponse(json.dumps(context))


class IncepGenerateTasksView(View):
    @method_decorator(check_group_permission)
    def post(self, request):
        id = request.POST.get('id')

        if IncepMakeExecTask.objects.filter(related_id=id).first():
            taskid = IncepMakeExecTask.objects.filter(related_id=id).first().taskid
            context = {'errCode': 201,
                       'dst_url': f'/projects/incep_of_records/incep_of_details/{taskid}'}
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
                        uid=self.request.user.uid,
                        user=obj.proposer,
                        taskid=taskid,
                        dst_host=host,
                        dst_database=database,
                        sql_content=row['SQL'],
                        sqlsha1=row['sqlsha1'],
                        type=obj.type,
                        category='1',
                        related_id=id,
                        group_id=obj.group_id
                    )

                context = {'errCode': 201,
                           'dst_url': f'/projects/incep_of_records/incep_of_details/{taskid}'}
            else:
                context = {'errCode': 400,
                           'errMsg': 'Leader审核未通过'}

        return HttpResponse(json.dumps(context))

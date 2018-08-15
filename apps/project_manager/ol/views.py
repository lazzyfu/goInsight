# -*- coding:utf-8 -*-
# edit by fuzongfei

import json
from datetime import datetime

from django.db import transaction
from django.db.models import Case, When, Value, CharField, F, Q, Count
from django.http import HttpResponse, JsonResponse
from django.shortcuts import render, get_object_or_404
from django.utils.decorators import method_decorator
from django.views import View
from django.views.generic import FormView

from project_manager.inception.inception_api import IncepSqlCheck
from project_manager.models import AuditContents, OlAuditContentsReply, IncepMakeExecTask, AuditTasks
from project_manager.ol.forms import WorkOrderAuditForm, WorkOrderReplyForm, WorkOrderApproveForm, \
    OnlineAuditRecordForm, \
    WorkOrderFeedbackForm, WorkOrderCloseForm
from project_manager.tasks import ding_notice_pull
from project_manager.utils import check_incep_alive, check_sql_filter, check_db_conn_status
from user_manager.permissions import permission_required, order_permission_required
from utils.tools import format_request


class OnlineWorkOrderAuditView(View):
    """渲染生成环境工单提交页面"""

    def get(self, request):
        return render(request, 'ol_work_order_audit.html')


class WorkOrderAuditView(View):
    """工单提交、处理"""

    @permission_required('can_commit')
    @method_decorator(check_incep_alive)
    @method_decorator(check_sql_filter)
    @transaction.atomic
    def post(self, request):
        data = format_request(request)
        form = WorkOrderAuditForm(data)
        if form.is_valid():
            context = form.is_save(request)
            return HttpResponse(json.dumps(context))
        else:
            error = form.errors.as_text()
            context = {'status': 2, 'msg': error}
            return HttpResponse(json.dumps(context))


class OnlineRecordsView(View):
    def get(self, request):
        return render(request, 'ol_records.html')


# 预发布环境工单
class StagingRecordsView(View):
    def get(self, request):
        return render(request, 'staging_records.html')


# 测试环境工单
class MOneRecordsView(View):
    def get(self, request):
        return render(request, 'test_records.html')


class AuditRecordsListView(View):
    def get(self, request):
        data = format_request(request)
        form = OnlineAuditRecordForm(data)
        result = {}
        if form.is_valid():
            cleaned_data = form.cleaned_data
            envi_desc = cleaned_data.get('envi_desc')
            limit_size = cleaned_data.get('limit_size')
            offset_size = cleaned_data.get('offset_size')
            search_content = cleaned_data.get('search_content')
            query = AuditContents.objects.filter(envi_desc=envi_desc).annotate(
                progress_value=Case(
                    When(progress='0', then=Value('待批准')),
                    When(progress='1', then=Value('未批准')),
                    When(progress='2', then=Value('已批准')),
                    When(progress='3', then=Value('处理中')),
                    When(progress='4', then=Value('已完成')),
                    When(progress='5', then=Value('已关闭')),
                    When(progress='6', then=Value('已勾住')),
                    output_field=CharField(),
                ),
                progress_color=Case(
                    When(progress__in=('0',), then=Value('btn-primary')),
                    When(progress__in=('2',), then=Value('btn-warning')),
                    When(progress__in=('1', '5'), then=Value('btn-danger')),
                    When(progress__in=('3',), then=Value('btn-info')),
                    When(progress__in=('4',), then=Value('btn-success')),
                    When(progress__in=('6',), then=Value('btn-default')),
                    output_field=CharField(),
                ),
                remark_desc=Case(
                    When(remark=0, then=Value('周三上线')),
                    When(remark=1, then=Value('紧急上线')),
                    When(remark=2, then=Value('数据修复')),
                    output_field=CharField(),
                ),
            )
            if search_content:
                obj = query.filter(Q(tasks__icontains=search_content) | Q(title__icontains=search_content) | Q(
                    proposer__icontains=search_content) | Q(envi_desc__icontains=search_content) | Q(
                    host__icontains=search_content) | Q(host__icontains=search_content) | Q(
                    database__icontains=search_content) | Q(contents__icontains=search_content))
            else:
                obj = query

            ol_total = obj.count()

            ol_records = obj.values('progress_color', 'tasks', 'host', 'operate_type',
                                    'database', 'progress_value', 'id', 'envi_desc',
                                    'title', 'proposer', 'operator',
                                    'created_at', 'remark_desc'
                                    ).order_by('-created_at')[offset_size:limit_size]
            result = {'total': ol_total, 'rows': list(ol_records)}

        return JsonResponse(result, safe=False)


class DeployTasksView(View):
    """获取部署任务列表"""

    def get(self, request):
        data = format_request(request)
        tasks = data.get('tasks')
        # 获取任务下所有工单分别在三个环境中的状态
        query = f"select id,title,proposer,tasks," \
                f"max(if(envi_desc=0, progress, -1)) as test," \
                f"max(if(envi_desc=1, progress, -1)) as staging," \
                f"max(if(envi_desc=2, progress, -1)) as product " \
                f"from auditsql_work_order where tasks='{tasks}' group by title order by id desc"

        result = []
        for row in AuditContents.objects.raw(query):
            result.append({
                'id': row.id,
                'title': row.title,
                'proposer': row.proposer,
                'tasks': row.tasks,
                'test': row.test,
                'staging': row.staging,
                'product': row.product
            })
        return JsonResponse(result, safe=False)

    def post(self, request):
        data = format_request(request)
        tasks = data.get('tasks')
        # 获取任务工单总数
        query = f"select id,count(x.id) as num from" \
                f"(select id from auditsql_work_order where tasks='{tasks}' group by title) as x"
        for row in AuditContents.objects.raw(query):
            count_num = row.num

        query_finish = f"select id, " \
                       f"max(if(envi_desc = 0 , progress , - 1)) as test, " \
                       f"max(if(envi_desc = 1 , progress , - 1)) as staging, " \
                       f"max(if(envi_desc = 2 , progress , - 1)) as product " \
                       f"from auditsql_work_order " \
                       f"where tasks = '{tasks}' group by title"
        status = []
        for row in AuditContents.objects.raw(query_finish):
            status.append([row.test, row.staging, row.product])

        # 计算完成的数量
        i = 0
        finish_status = ['4', '6']
        for x in status:
            if all(False for y in x if y not in finish_status):
                i += 1

            # 关闭状态的也计入已完成状态
            if '5' in x:
                i += 1

        finish_value = '/'.join((str(i), str(count_num)))

        return JsonResponse({'num': finish_value}, safe=False)


class DingNoticeView(View):
    """at通知未完成工单的开发"""

    def post(self, request):
        data = format_request(request)
        tasks = data.get('tasks')

        query = f"select id,title,proposer,tasks,b.mobile, " \
                f"max(if(envi_desc=0, progress, -1)) as test, " \
                f"max(if(envi_desc=1, progress, -1)) as staging, " \
                f"max(if(envi_desc=2, progress, -1)) as product " \
                f"from auditsql_work_order a join auditsql_useraccount b on a.proposer = b.username " \
                f"where tasks='{tasks}' group by title order by id desc"

        result = []
        for row in AuditContents.objects.raw(query):
            result.append({
                'id': row.id,
                'title': row.title,
                'mobile': row.mobile,
                'proposer': row.proposer,
                'tasks': row.tasks,
                'test': row.test,
                'staging': row.staging,
                'product': row.product
            })

        # 获取今天为周几
        today = datetime.now().weekday() + 1
        # 如果是周二，通知预发布环境未完成工单的开发
        analyze_result = []
        if today == 2:
            for row in result:
                if not row['staging'] in ('4', '5', '6'):
                    # 剔除test环境未关闭的工单，已关闭的不做处理
                    if row['test'] != '5':
                        analyze_result.append(row)
            if analyze_result:
                ding_notice_pull.delay(analyze_result, today)
                context = {'status': 0, 'msg': '已钉'}
            else:
                context = {'status': 2, 'msg': '预发布环境任务已完成，无法钉'}

        # 如果是周三，通知生产环境布未完成工单的开发
        elif today == 3:
            for row in result:
                if not row['product'] in ('4', '5', '6'):
                    # 剔除test环境未关闭的工单，已关闭的不做处理
                    if row['staging'] not in ('-1', '5'):
                        analyze_result.append(row)
            if analyze_result:
                ding_notice_pull.delay(analyze_result, today)
                context = {'status': 0, 'msg': '已钉'}
            else:
                context = {'status': 2, 'msg': '生产环境任务已完成，无法钉'}
        else:
            context = {'status': 2, 'msg': '只能在指定周二或周三使用该功能'}

        return HttpResponse(json.dumps(context))


class WorkOrderApproveView(FormView):
    """线上工单审批操作，需要can_approve权限"""
    form_class = WorkOrderApproveForm

    def dispatch(self, request, *args, **kwargs):
        return super(WorkOrderApproveView, self).dispatch(request, *args, **kwargs)

    @permission_required('can_approve')
    @transaction.atomic
    def form_valid(self, form):
        context = form.is_save(self.request)
        return HttpResponse(json.dumps(context))

    def form_invalid(self, form):
        error = form.errors.as_text()
        context = {'status': 2, 'msg': error}
        return HttpResponse(json.dumps(context))


class WorkOrderFeedbackView(FormView):
    """线上工单反馈，反馈执行进度"""
    form_class = WorkOrderFeedbackForm

    @order_permission_required('can_execute')
    @transaction.atomic
    def form_valid(self, form):
        context = form.is_save(self.request)
        return HttpResponse(json.dumps(context))

    def form_invalid(self, form):
        error = form.errors.as_text()
        context = {'status': 2, 'msg': error}
        return HttpResponse(json.dumps(context))


class WorkOrderCloseView(FormView):
    """关闭记录"""
    form_class = WorkOrderCloseForm

    @order_permission_required('can_approve')
    @transaction.atomic
    def form_valid(self, form):
        context = form.is_save(self.request)
        return HttpResponse(json.dumps(context))

    def form_invalid(self, form):
        error = form.errors.as_text()
        context = {'status': 2, 'msg': error}
        return HttpResponse(json.dumps(context))


class WorkOrderDetailsView(View):
    """查看审核工单详情"""

    def get(self, request, id):
        obj = AuditContents.objects.annotate(
            progress_value=Case(
                When(progress='0', then=Value('待批准')),
                When(progress='1', then=Value('未批准')),
                When(progress='2', then=Value('已批准')),
                When(progress='3', then=Value('处理中')),
                When(progress='4', then=Value('已完成')),
                When(progress='5', then=Value('已关闭')),
                When(progress='6', then=Value('已勾住')),
                output_field=CharField(),
            ),
        )

        contents = obj.get(id=id)
        reply_contents = OlAuditContentsReply.objects.annotate(
            username=F('user__username'),
            avatar_file=F('user__avatar_file'),
        ).filter(reply__id=id).values('username', 'avatar_file', 'reply_contents', 'created_at')

        return render(request, 'work_order_details.html',
                      {'contents': contents,
                       'reply_contents': reply_contents})


class WorkOrderReplyView(FormView):
    """处理用户的回复信息"""

    form_class = WorkOrderReplyForm

    def form_valid(self, form):
        context = form.is_save(self.request)
        return HttpResponse(json.dumps(context))

    def form_invalid(self, form):
        error = form.errors.as_text()
        context = {'status': 2, 'msg': error}
        return HttpResponse(json.dumps(context))


class GeneratePerformTasksView(View):
    """工单转换成执行任务"""

    @order_permission_required('can_execute')
    def post(self, request):
        id = request.POST.get('id')
        envi_desc = request.POST.get('envi_desc')

        obj = get_object_or_404(AuditContents, pk=id)

        status, msg = check_db_conn_status(obj.host, obj.port)
        if status:

            # 只要审核通过后，才能生成执行任务
            if obj.progress in ('2', '3', '4', '6'):
                if IncepMakeExecTask.objects.filter(related_id=id).first():
                    taskid = IncepMakeExecTask.objects.filter(related_id=id).first().taskid
                    context = {'status': 0,
                               'jump_url': f'/projects/pt/perform_records/perform_details/{taskid}'}
                else:
                    host = obj.host
                    database = obj.database
                    port = obj.port
                    sql_content = obj.contents

                    # 实例化
                    of_audit = IncepSqlCheck(sql_content, host, port, database, request.user.username)

                    # 对OSC执行的SQL生成sqlsha1
                    result = of_audit.make_sqlsha1()
                    taskid = datetime.now().strftime("%Y%m%d%H%M%S%f")

                    # 生成执行任务记录
                    for row in result:
                        IncepMakeExecTask.objects.create(
                            uid=request.user.uid,
                            user=obj.proposer,
                            taskid=taskid,
                            dst_host=host,
                            dst_port=port,
                            dst_database=database,
                            sql_content=row['SQL'],
                            sqlsha1=row['sqlsha1'],
                            affected_row=row['Affected_rows'],
                            type=obj.operate_type,
                            envi_desc=envi_desc,
                            related_id=id
                        )

                    context = {'status': 0,
                               'jump_url': f'/projects/pt/perform_records/perform_details/{taskid}'}
            else:
                context = {'status': 2, 'msg': '审核未通过或任务已关闭'}
        else:
            context = {'status': 2, 'msg': f'无法连接到数据库，请联系DBA\n主机: {obj.host}\n端口: {obj.port}'}

        return HttpResponse(json.dumps(context))


class ROnlineAuditTasksView(View):
    def get(self, request):
        return render(request, 'ol_tasks.html')


class OnlineAuditTasksListView(View):
    def get(self, request):
        data = AuditTasks.objects.all().values('id', 'tasks', 'username', 'expire_time', 'created_at')
        return JsonResponse(list(data), safe=False)

    # 开发可以自己创建
    def post(self, request):
        data = format_request(request)
        tasks = data.get('tasks')
        expire_time = data.get('expire_time')
        action = data.get('action')
        id = data.get('id')
        if action == 'new_tasks':
            if AuditTasks.objects.filter(tasks=tasks).first():
                context = {'status': 2, 'msg': '记录已存在，不能重复创建'}
            else:
                AuditTasks.objects.create(tasks=tasks, expire_time=expire_time, username=request.user.displayname)
                context = {'status': 0, 'msg': '创建成功'}
        elif action == 'delete_tasks':
            for i in id.split(','):
                AuditTasks.objects.get(pk=i).delete()
            context = {'status': 0, 'msg': '删除成功'}
        return HttpResponse(json.dumps(context))

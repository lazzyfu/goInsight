import json
import re

import sqlparse
from django.db.models import F
from django.http import JsonResponse, HttpResponse
from django.shortcuts import render, get_object_or_404
from django.views import View
from django.views.generic import FormView, ListView
from pure_pagination import PaginationMixin

from ProjectManager.forms import InceptionSqlOperateForm, OnlineAuditCommitForm
from UserManager.models import GroupsDetail, UserAccount, ContactsDetail, Contacts
from apps.ProjectManager.inception.inception_api import GetDatabaseListApi, InceptionApi, GetBackupApi
from utils.tools import format_request
from .models import InceptionHostConfig, InceptionSqlOperateRecord, Remark


class ProjectListView(View):
    def get(self, request):
        return render(request, 'index.html')


class InceptionSqlOperateView(FormView):
    form_class = InceptionSqlOperateForm
    template_name = 'inception_sql_operate.html'

    def form_valid(self, form):
        cleaned_data = form.cleaned_data
        host = cleaned_data['host']
        database = cleaned_data['database']
        op_action = cleaned_data.get('op_action')
        op_type = cleaned_data['op_type']
        sql_content = cleaned_data['sql_content']

        DDL_FILTER = 'ALTER TABLE|CREATE TABLE|TRUNCATE TABLE'
        DML_FILTER = 'INSERT INTO|;UPDATE|^UPDATE|DELETE FROM'

        checkData = InceptionApi().sqlprepare(sqlcontent=sql_content, host=host, database=database,
                                              action='check')

        # 修改表结构
        if op_action == 'op_schema':
            if re.search(DML_FILTER, sql_content, re.I):
                context = {'errMsg': f'DDL模式下, 不支持SELECT|UPDATE|DELETE|INSERT语句', 'errCode': 400}
            else:
                if op_type == 'check':
                    context = {'data': checkData, 'errCode': 200}
                if op_type == 'commit':
                    if 1 in [x['errlevel'] for x in checkData] or 2 in [x['errlevel'] for x in checkData]:
                        context = {'errMsg': 'SQL语法检查未通过, 请执行语法检测', 'errCode': 400}
                    else:
                        executeData = InceptionApi().sqlprepare(sqlcontent=sql_content, host=host,
                                                                database=database,
                                                                action='execute')
                        context = form.is_save(self.request, executeData)

        # 修改数据
        if op_action == 'op_data':
            if re.search(DDL_FILTER, sql_content, re.I):
                context = {'errMsg': f'DML模式下, 不支持ALTER|CREATE|TRUNCATE语句', 'errCode': 400}
            else:
                if op_type == 'check':
                    context = {'data': checkData, 'errCode': 200}
                if op_type == 'commit':
                    if 1 in [x['errlevel'] for x in checkData] or 2 in [x['errlevel'] for x in checkData]:
                        context = {'errMsg': 'SQL语法检查未通过, 请执行语法检测', 'errCode': 400}
                    else:
                        executeData = InceptionApi().sqlprepare(sqlcontent=sql_content, host=host,
                                                                database=database,
                                                                action='execute')
                        context = form.is_save(self.request, executeData)

        return HttpResponse(json.dumps(context))

    def form_invalid(self, form):
        # error = form.errors.as_text()
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
    """获取inception审核的目标数据库配置"""

    def get(self, request):
        envResult = InceptionHostConfig.objects.all().values('host', 'comment')
        return JsonResponse(list(envResult), safe=False)


class GetDatabaseListView(View):
    """列出选中环境的数据库库名"""

    def post(self, request):
        data = format_request(request)
        host = data['host']
        dbList = GetDatabaseListApi(host).get_dbname()
        return HttpResponse(json.dumps(dbList))


class InceptionSqlRecords(PaginationMixin, ListView):
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


class InceptionAllSqlDetailView(View):
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

        return render(request, 'inception_all_sql_detail.html',
                      {'originalSql': originalSql, 'rollbackSql': rollbackSql})


class InceptionSingleSqlDetailView(View):
    """查看当前用户会话执行的每条sql的详情"""

    def get(self, request, sequence):
        sqlDetail = get_object_or_404(InceptionSqlOperateRecord, sequence=sequence)
        sequenceResult = [{'backupdbName': sqlDetail.backup_dbname, 'sequence': sqlDetail.sequence}]
        rollbackSql = GetBackupApi(sequenceResult).get_backupinfo()
        return render(request, 'inception_single_sql_detail.html',
                      {'sqlDetail': sqlDetail, 'rollbackSql': rollbackSql})


class OnlineSqlCommitView(FormView):
    """
    处理用户提交的审核内容
    """

    def get(self, request):
        return render(request, 'online_sql_commit.html')

    def post(self, request):
        data = format_request(request)
        form = OnlineAuditCommitForm(data)
    # def post(self, request):
    #     data = format_request(request)
    #     upload_files = request.FILES.getlist('files')
    #     form = AuditCommitForm(data)
    #     if form.is_valid():
    #         cleaned_data = form.cleaned_data
    #         title = cleaned_data['title'] + '__[' + datetime.now().strftime("%Y%m%d%H%M") + ']'
    #         remark = ','.join(
    #             Remark.objects.filter(id__in=cleaned_data['remark'].split(',')).values_list("remark", flat=True))
    #         verifier = UserAccount.objects.get(uid=cleaned_data['verifier'])
    #         operate_dba = UserAccount.objects.get(uid=cleaned_data['operate_dba'])
    #         with transaction.atomic():
    #             check_title_unique = AuditContents.objects.filter(title=title).first()
    #             if check_title_unique:
    #                 result = {'status': '400', 'msg': '标题名重复, 请更换, 谢谢'}
    #             else:
    #                 # 插入提交的数据
    #                 AuditContents.objects.create(
    #                     title=title,
    #                     items_id=cleaned_data['project'],
    #                     remark=remark,
    #                     proposer=request.user.username,
    #                     verifier=verifier,
    #                     operate_dba=operate_dba,
    #                     email_cc=cleaned_data['email_cc'],
    #                     environment=cleaned_data['environment'],
    #                     contents=cleaned_data['contents']
    #                 )
    #
    #                 # 生成一条详情记录
    #                 latest_id = AuditContents.objects.latest('id').id
    #                 AuditContentsDetail.objects.create(content_id=latest_id)
    #
    #                 # 处理上传文件
    #                 if upload_files:
    #                     for file in upload_files:
    #                         file_instance = UploadFiles(file_name=file.name, files=file, file_size=file.size,
    #                                                     content_id=latest_id,
    #                                                     content_type=file.content_type)
    #                         file_instance.save()
    #
    #                 # 分配权限：只有verifier才能执行批准操作
    #                 assign_perm('leader_verify', verifier, AuditContents.objects.get(id=latest_id))
    #
    #                 # 发送通知邮件
    #                 send_commit_mail.delay(latest_id=latest_id)
    #
    #                 result = {'status': '200', 'msg': '提交成功, 正在跳转到首页'}
    #     else:
    #         error = form.errors.as_text()
    #         result = {'status': '400', 'msg': error}
    #     return HttpResponse(json.dumps(result))


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

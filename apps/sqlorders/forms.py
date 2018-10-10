# -*- coding:utf-8 -*-
# edit by fuzongfei
import json
import os
from datetime import datetime

import psutil
import sqlparse
from django import forms
from django.core.serializers.json import DjangoJSONEncoder
from django.db.models import Max, Case, When, Value, CharField, Q
from django.shortcuts import get_object_or_404

from sqlorders.models import sql_type_choice, SqlOrdersEnvironment, envi_choice, MysqlSchemas, \
    SqlOrdersTasksVersions, SysConfig, SqlOrderReply
from sqlorders.utils import check_db_conn_status, GetTableInfo, sql_filter, GetInceptionBackupApi
from users.models import RolePermission
from .tasks import *


class GetTablesForm(forms.Form):
    schema = forms.CharField()

    def query(self):
        cdata = self.cleaned_data
        schema = cdata['schema']
        host, port, schema = schema.split(',')

        status, msg = check_db_conn_status(host, port)
        if status:
            table_list = GetTableInfo(host, port, schema).get_column_info()
            context = {'status': 0, 'msg': '', 'data': table_list}
        else:
            context = {'status': 2, 'msg': f'无法连接到数据库，请联系管理员，\n主机: {host}\n端口: {port}'}
        return context


class GetParentSchemasForm(forms.Form):
    envi_id = forms.ChoiceField(required=True, choices=envi_choice, label=u'环境')

    def query(self):
        cdata = self.cleaned_data
        envi_id = cdata.get('envi_id')
        parent_id = SqlOrdersEnvironment.objects.get(envi_id=envi_id).parent_id

        queryset = MysqlSchemas.objects.filter(envi_id=parent_id).filter(is_master=1).values('host', 'port', 'schema',
                                                                                             'comment')
        serialize_data = json.dumps(list(queryset), cls=DjangoJSONEncoder)
        return serialize_data


class SqlOrdersAuditForm(forms.Form):
    title = forms.CharField(max_length=100, required=True, label=u'标题')
    description = forms.CharField(max_length=1024, required=False, label=u'需求url或描述性文字')
    task_version = forms.CharField(max_length=256, required=False, label=u'上线版本号')
    auditor = forms.CharField(required=True, label=u'工单审核人')
    email_cc = forms.CharField(required=False, label=u'抄送联系人')
    database = forms.CharField(required=True, max_length=64, label=u'数据库')
    remark = forms.CharField(required=True, max_length=256, min_length=1, label=u'工单备注')
    sql_type = forms.ChoiceField(choices=sql_type_choice, label=u'操作类型，是DDL还是DML')
    contents = forms.CharField(widget=forms.Textarea, label=u'审核内容')
    envi_id = forms.ChoiceField(choices=envi_choice, required=False)

    def save(self, request):
        cdata = self.cleaned_data
        title = cdata.get('title') + '_[' + datetime.now().strftime("%Y%m%d%H%M%S") + ']'
        description = cdata.get('description')
        task_version = cdata.get('task_version')
        auditor = cdata.get('auditor')
        email_cc = ','.join(self.data.getlist('email_cc'))
        remark = cdata.get('remark')
        host, port, database = cdata.get('database').split(',')
        sql_type = cdata.get('sql_type')
        contents = cdata.get('contents')
        envi_id = cdata.get('envi_id')

        if sql_type == 'DDL':
            max_parent_id = SqlOrdersEnvironment.objects.all().aggregate(Max('parent_id'))['parent_id__max']
            envi_id = SqlOrdersEnvironment.objects.get(parent_id=max_parent_id).envi_id
        elif sql_type == 'DML':
            envi_id = envi_id

        result = InceptionSqlApi(host, port, database, contents, request.user.username).is_check_pass()
        if result.get('status') == 2:
            context = result
        else:
            obj = SqlOrdersContents.objects.create(
                title=title,
                description=description,
                task_version=task_version,
                sql_type=sql_type,
                host=host,
                database=database,
                port=port,
                envi_id=envi_id,
                remark=remark,
                proposer=request.user.username,
                auditor=auditor,
                email_cc=email_cc,
                contents=contents
            )

            # 发送邮件
            msg_pull = SqlOrdersMsgPull(id=obj.id, user=request.user.username, type='commit')
            msg_pull.run()

            # 跳转到工单记录页面
            context = {'status': 0, 'jump_url': f'/sqlorders/sql_orders_list/{envi_id}'}
            print(context)
        return context


class SyntaxCheckForm(forms.Form):
    host = forms.CharField(required=True, max_length=64)
    sql_type = forms.ChoiceField(choices=sql_type_choice, label=u'操作类型，是DDL还是DML')
    contents = forms.CharField(widget=forms.Textarea)

    def query(self, request):
        cdata = self.cleaned_data
        host, port, database = cdata.get('host').split(',')
        sql_type = cdata.get('sql_type')
        contents = cdata.get('contents')

        # 对检测的SQL类型进行区分
        filter_result = sql_filter(contents, sql_type)

        # 实例化
        of_audit = InceptionSqlApi(host, port, database, contents, request.user.username)

        if filter_result['status'] == 2:
            context = filter_result
        else:
            # SQL语法检查
            context = of_audit.run_check()
        return context


class BeautifySQLForm(forms.Form):
    contents = forms.CharField(widget=forms.Textarea)

    def beautify(self):
        cdata = self.cleaned_data
        contents = cdata.get('contents').strip()

        sql_split = []
        for stmt in sqlparse.split(contents):
            sql = sqlparse.parse(stmt)[0]
            sql_comment = sql.token_first()
            if isinstance(sql_comment, sqlparse.sql.Comment):
                sql_split.append({'comment': sql_comment.value, 'sql': sql.value.replace(sql_comment.value, '')})
            else:
                sql_split.append({'comment': '', 'sql': sql.value})

        beautify_sql_list = []
        for row in sql_split:
            comment = row['comment']
            sql = row['sql']
            res = sqlparse.parse(sql)
            syntax_type = res[0].token_first().ttype.__str__()
            if syntax_type == 'Token.Keyword.DDL':
                sql_format = sqlparse.format(sql)
                beautify_sql_list.append(comment + sql_format)
            elif syntax_type == 'Token.Keyword.DML':
                sql_format = sqlparse.format(sql, reindent=True)
                beautify_sql_list.append(comment + sql_format)
            elif syntax_type == 'Token.Keyword':
                beautify_sql_list.append(comment + sql)
        context = {'data': '\n\n'.join(beautify_sql_list)}
        return context


class SqlOrderListForm(forms.Form):
    envi_id = forms.IntegerField(required=True, label=u'环境')
    limit_size = forms.IntegerField(required=True, label=u'每页显示数量')
    offset_size = forms.IntegerField(required=True, label=u'分页偏移量')
    search_content = forms.CharField(max_length=128, required=False, label='搜索内容')

    def query(self, request):
        cdata = self.cleaned_data
        envi_id = cdata.get('envi_id')
        limit_size = cdata.get('limit_size')
        offset_size = cdata.get('offset_size')
        search_content = cdata.get('search_content')

        # 获取用户的权限，用于前端表格的列的显示
        role_name = request.user.user_role()
        perm_list = list(
            RolePermission.objects.filter(role__role_name=role_name).values_list('permission_name', flat=True))

        permissions = {'permissions': perm_list}

        query = SqlOrdersContents.objects.filter(envi_id=envi_id).annotate(
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
        )
        if search_content:
            obj = query.filter(Q(task_version__icontains=search_content) | Q(title__icontains=search_content) | Q(
                proposer__icontains=search_content) | Q(
                host__icontains=search_content) | Q(host__icontains=search_content) | Q(
                database__icontains=search_content) | Q(contents__icontains=search_content))
        else:
            obj = query

        ol_total = obj.count()

        ol_records = obj.values('progress_color', 'task_version', 'host', 'port', 'sql_type',
                                'database', 'progress_value', 'id', 'envi_id',
                                'title', 'proposer', 'auditor',
                                'created_at', 'remark'
                                ).order_by('-created_at')[offset_size:limit_size]
        rows = []
        for x in list(ol_records):
            x.update(permissions)
            rows.append(x)
        result = {'total': ol_total, 'rows': rows}
        return result


class SqlOrdersApproveForm(forms.Form):
    id = forms.IntegerField(required=True)
    status = forms.CharField(max_length=10, required=True)
    addition_info = forms.CharField(required=False)

    def save(self, request):
        cdata = self.cleaned_data
        id = cdata.get('id')
        status = cdata.get('status')
        addition_info = cdata.get('addition_info')

        data = SqlOrdersContents.objects.get(pk=id)

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
                    # 发送邮件
                    msg_pull = SqlOrdersMsgPull(id=id, user=request.user.username, type='approve',
                                                addition_info=addition_info)
                    msg_pull.run()
                    context = {'status': 0, 'msg': '操作成功、审核通过'}

                # 当用户点击的是不通过, 状态变为：未批准
                elif status == u'不通过':
                    data.progress = '1'
                    data.operate_time = timezone.now()
                    data.save()
                    # 发送邮件
                    msg_pull = SqlOrdersMsgPull(id=id, user=request.user.username, type='approve',
                                                addition_info=addition_info)
                    msg_pull.run()
                    context = {'status': 0, 'msg': '操作成功、审核未通过'}
            # 其他情况
            else:
                context = {'status': 2, 'msg': '操作失败、请不要重复提交'}
        return context


class SqlOrdersFeedbackForm(forms.Form):
    id = forms.IntegerField(required=True)
    status = forms.CharField(max_length=10, required=True)
    addition_info = forms.CharField(required=False)

    def save(self, request):
        cdata = self.cleaned_data
        id = cdata.get('id')
        status = cdata.get('status')
        addition_info = cdata.get('addition_info')

        data = SqlOrdersContents.objects.get(pk=id)

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
                    # 发送邮件
                    msg_pull = SqlOrdersMsgPull(id=id, user=request.user.username, type='feedback',
                                                addition_info=addition_info)
                    msg_pull.run()
                    context = {'status': 0, 'msg': '操作成功、正在处理中'}

                # 当用户点击的是已完成, 状态变为：已完成
                elif status == u'已完成':
                    data.progress = '4'
                    data.updated_at = timezone.now()
                    data.save()
                    # 发送邮件
                    msg_pull = SqlOrdersMsgPull(id=id, user=request.user.username, type='feedback',
                                                addition_info=addition_info)
                    msg_pull.run()
                    context = {'status': 0, 'msg': '操作成功、处理完成'}

            # 未批准
            elif data.progress == '1' or data.progress == '0':
                context = {'status': 2, 'msg': '操作失败、审核未通过'}
            # 其他情况
            else:
                context = {'status': 2, 'msg': '操作失败、请不要重复提交'}
        return context


class SqlOrdersCloseForm(forms.Form):
    id = forms.IntegerField(required=True)
    status = forms.CharField(max_length=10, required=True)
    addition_info = forms.CharField(required=False)

    def save(self, request):
        cdata = self.cleaned_data
        id = cdata.get('id')
        status = cdata.get('status')
        addition_info = cdata.get('addition_info')

        data = SqlOrdersContents.objects.get(pk=id)

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
                        # 发送邮件
                        msg_pull = SqlOrdersMsgPull(id=id, user=request.user.username, type='close',
                                                    addition_info=addition_info)
                        msg_pull.run()
                        context = {'status': 0, 'msg': '操作成功、记录关闭成功'}

                elif status == u'结束':
                    context = {'status': 2, 'msg': '操作失败、关闭窗口'}
            else:
                context = {'status': 2, 'msg': '操作失败、<关闭原因>输入不能少于5个字符'}
        return context


class HookSqlOrdersForm(forms.Form):
    id = forms.CharField(required=True, label=u'审核内容id')
    database = forms.CharField(required=True)
    envi_id = forms.ChoiceField(required=True, choices=envi_choice, label=u'环境')

    def save(self, request):
        cdata = self.cleaned_data
        host, port, database = cdata['database'].split(',')
        id = cdata.get('id')
        parent_id = SqlOrdersEnvironment.objects.get(envi_id=cdata['envi_id']).parent_id
        jump_url = f'/sqlorders/sql_orders_list/{parent_id}'

        data = SqlOrdersContents.objects.get(pk=id)
        if SqlOrdersContents.objects.filter(title=data.title, envi_id=parent_id).exists():
            # 如果指定的环境存在已被钩的工单，直接跳转
            context = {'status': 0, 'jump_url': jump_url}
        else:
            # 工单状态必须为已完成
            if data.progress in ('4', '6'):
                obj = SqlOrdersContents.objects.create(
                    title=data.title,
                    description=data.description,
                    task_version=data.task_version,
                    sql_type=data.sql_type,
                    host=host,
                    database=database,
                    port=port,
                    envi_id=parent_id,
                    progress='2',
                    remark=data.remark,
                    proposer=data.proposer,
                    auditor=data.auditor,
                    contents=data.contents,
                    updated_at=timezone.now()
                )

                # 更新状态为：已勾住
                SqlOrdersContents.objects.filter(pk=id).update(progress='6')

                # 发送邮件
                msg_pull = SqlOrdersMsgPull(id=obj.id, user=request.user.username, type='hook')
                msg_pull.run()

                # 跳转到工单记录页面
                context = {'status': 0, 'jump_url': jump_url}
            else:
                context = {'status': 2, 'msg': '当前工单进度：未完成，无法勾住'}

        return context


class GeneratePerformTasksForm(forms.Form):
    id = forms.CharField(required=True, label=u'审核内容id')
    envi_id = forms.ChoiceField(required=True, choices=envi_choice, label=u'环境')

    def save(self, request):
        cdata = self.cleaned_data
        id = cdata.get('id')
        envi_id = cdata.get('envi_id')

        obj = get_object_or_404(SqlOrdersContents, pk=id)
        status, msg = check_db_conn_status(obj.host, obj.port)
        if status:
            # 只要审核通过后，才能生成执行任务
            if obj.progress in ('2', '3', '4', '6'):
                if SqlOrdersExecTasks.objects.filter(related_id=id).exists():
                    taskid = SqlOrdersExecTasks.objects.filter(related_id=id).first().taskid
                    context = {'status': 0,
                               'jump_url': f'/sqlorders/perform_tasks/{taskid}'}
                else:
                    host = obj.host
                    database = obj.database
                    port = obj.port
                    sql_content = obj.contents

                    # 实例化
                    incep_audit = InceptionSqlApi(host, port, database, sql_content, request.user.username)

                    # 对OSC执行的SQL生成sqlsha1
                    result = incep_audit.make_sqlsha1()
                    taskid = datetime.now().strftime("%Y%m%d%H%M%S%f")

                    # 生成执行任务记录
                    for row in result:
                        SqlOrdersExecTasks.objects.create(
                            uid=request.user.uid,
                            user=obj.proposer,
                            taskid=taskid,
                            host=host,
                            port=port,
                            database=database,
                            sql=row['SQL'],
                            sqlsha1=row['sqlsha1'],
                            affected_row=row['Affected_rows'],
                            sql_type=obj.sql_type,
                            envi_id=envi_id,
                            related_id=id
                        )

                    context = {'status': 0,
                               'jump_url': f'/sqlorders/perform_tasks/{taskid}'}
            else:
                context = {'status': 2, 'msg': '审核未通过或任务已关闭'}
        else:
            context = {'status': 2, 'msg': f'无法连接到数据库，请联系系统管理员\n主机: {obj.host}\n端口: {obj.port}'}

        return context


class SinglePerformTasksForm(forms.Form):
    id = forms.IntegerField()

    def exec(self, request):
        cdata = self.cleaned_data
        id = cdata.get('id')

        obj = SqlOrdersExecTasks.objects.get(id=id)
        host = obj.host
        port = obj.port
        database = obj.database
        sql = obj.sql + ';'

        key = ast.literal_eval(obj.taskid)
        if 'run' == cache.get(key):
            context = {'status': 1, 'msg': '正在自动化操作，请不要手动执行'}
        else:
            status = ''
            query = f"select id,group_concat(exec_status) as exec_status from sqlaudit_sql_orders_execute_tasks " \
                    f"where taskid={obj.taskid} group by taskid"
            for row in SqlOrdersExecTasks.objects.raw(query):
                status = row.exec_status.split(',')

            # 每次只能执行一条任务，不可同时执行，避免数据库压力
            if '2' in status or '3' in status:
                context = {'status': 2, 'msg': '请等待当前任务执行完成'}
            else:
                # 避免任务重复点击执行
                if obj.exec_status not in ('0', '5', '6'):
                    context = {'status': 2, 'msg': '请不要重复操作任务'}
                else:
                    # 将任务进度设置为：处理中
                    obj.exec_status = '2'
                    obj.save()
                    print(obj.exec_status)

                    # 如果sqlsha1存在，说明是大表，需要使用工具进行修改
                    # inception_osc_min_table_size默认为16M
                    # 如果此处向走gh-ost，则设置inception_osc_min_table_size=0
                    if obj.sqlsha1:
                        # 判断是否使用gh-ost执行
                        if SysConfig.objects.get(key='is_ghost').is_enabled == '0':
                            ghost_async_tasks.delay(user=request.user.username,
                                                    id=id,
                                                    sql=sql,
                                                    host=obj.host,
                                                    port=obj.port,
                                                    database=obj.database)
                            context = {'status': 1, 'msg': '任务已提交，请查看输出'}
                        else:
                            # 使用pt-online-schema-change执行
                            # 异步执行SQL任务
                            r = incep_async_tasks.delay(user=request.user.username,
                                                        id=id,
                                                        sql=sql,
                                                        host=host,
                                                        port=port,
                                                        database=database,
                                                        sqlsha1=obj.sqlsha1,
                                                        backup='yes',
                                                        exec_status='2')
                            task_id = r.task_id
                            # 将celery task_id写入到表
                            obj.celery_task_id = task_id
                            obj.save()
                            # 获取OSC执行进度
                            get_osc_percent.delay(task_id=task_id)

                            context = {'status': 1, 'msg': '任务已提交，请查看输出'}

                    else:
                        # 当affected_row>100000时，只执行不备份
                        if obj.affected_row > 1000000:
                            incep_async_tasks.delay(user=request.user.username,
                                                    id=id,
                                                    sql=sql,
                                                    host=host,
                                                    port=port,
                                                    database=database,
                                                    exec_status='2')
                        else:
                            # 当affected_row<=100000时，执行并备份
                            incep_async_tasks.delay(user=request.user.username,
                                                    id=id,
                                                    backup='yes',
                                                    sql=sql,
                                                    host=host,
                                                    port=port,
                                                    database=database,
                                                    exec_status='2')

                        context = {'status': 1, 'msg': '任务已提交，请查看输出'}
            # 更新父任务进度
            update_audit_content_progress(request.user.username, obj.taskid)
        return context


class PerformTasksOpForm(forms.Form):
    id = forms.IntegerField(required=True)
    action = forms.ChoiceField(
        choices=(
            ('pause_ghost', 'pause_ghost'),
            ('recovery_ghost', 'recovery_ghost'),
            ('stop_ghost', 'stop_ghost'),
            ('stop_ptosc', 'stop_ptosc')
        ), error_messages={'required': '传入的值错误, 不接受非法的值'}
    )

    def op(self, request):
        cdata = self.cleaned_data
        id = cdata.get('id')
        action = cdata.get('action')

        obj = SqlOrdersExecTasks.objects.get(id=id)
        celery_task_id = obj.celery_task_id

        context = {}
        if obj.exec_status in ('0', '1', '4'):
            context = {'status': 2, 'msg': '请不要重复操作任务'}
        else:
            # 判断是否使用gh-ost执行
            if SysConfig.objects.get(key='is_ghost').is_enabled == '0':
                # 获取gh-ost的sock文件
                # 将语句中的注释和SQL分离
                sql_split = {}
                for stmt in sqlparse.split(obj.sql):
                    sql = sqlparse.parse(stmt)[0]
                    sql_comment = sql.token_first()
                    if isinstance(sql_comment, sqlparse.sql.Comment):
                        sql_split = {'comment': sql_comment.value, 'sql': sql.value.replace(sql_comment.value, '')}
                    else:
                        sql_split = {'comment': '', 'sql': sql.value}

                # 获取不包含注释的SQL语句
                sql = sql_split['sql']
                formatsql = re.compile('^ALTER(\s+)TABLE(\s+)([\S]*)(\s+)(ADD|CHANGE|REMAME|MODIFY|DROP)([\s\S]*)',
                                       re.I)
                match = formatsql.match(sql)
                # 由于gh-ost不支持反引号，会被解析成命令，因此此处替换掉
                table = match.group(3).replace('`', '')
                # 将schema.table进行处理，这种情况gh-ost不识别，只保留table
                if len(table.split('.')) > 1:
                    table = table.split('.')[1]
                sock = os.path.join('/tmp', f"gh-ost.{obj.database}.{table}.sock")
                # 判断程序是否允许
                if psutil.pid_exists(obj.ghost_pid):
                    if os.path.exists(sock):
                        if action == 'pause_ghost':
                            pause_cmd = f"echo throttle | nc -U {sock}"
                            p = subprocess.Popen(pause_cmd, shell=True)
                            p.wait()
                            context = {'status': 1, 'msg': '暂停动作已执行，请查看输出'}

                        if action == 'recovery_ghost':
                            recovery_cmd = f"echo no-throttle | nc -U {sock}"
                            p = subprocess.Popen(recovery_cmd, shell=True)
                            p.wait()
                            context = {'status': 1, 'msg': '恢复动作已执行，请查看输出'}

                        if action == 'stop_ghost':
                            stop_cmd = f"echo panic | nc -U {sock}"
                            p = subprocess.Popen(stop_cmd, shell=True)
                            p.wait()
                            context = {'status': 1, 'msg': '终止动作已执行，请查看输出'}
                    else:
                        context = {'status': 2, 'msg': f'不能找到文件{sock}, 操作失败'}
                else:
                    os.remove(sock) if os.path.exists(sock) else None
                    context = {'status': 2, 'msg': '进程不存在，操作失败'}
            else:
                # pt-osc操作
                if action == 'stop_ptosc':
                    # 停止pt-osc进程
                    stop_incep_osc.delay(user=request.user.username,
                                         id=id,
                                         celery_task_id=celery_task_id)
                    context = {'status': 2, 'msg': '终止动作已执行，请查看输出'}
                else:
                    context = {'status': 2, 'msg': '非ghost任务，操作失败'}
        return context


class FullPerformTasksForm(forms.Form):
    taskid = forms.CharField()

    def exec(self, request):
        cdata = self.cleaned_data
        taskid = cdata.get('taskid')

        query = f"select * from sqlaudit_sql_orders_execute_tasks where taskid={taskid} order by id asc"

        key = ast.literal_eval(taskid)
        if 'run' == cache.get(key):
            context = {'status': 1, 'msg': '当前任务正在运行，请不要重复执行'}
        else:
            cache.set(key, 'run', timeout=600)
            incep_multi_tasks.delay(username=request.user.username,
                                    query=query,
                                    key=key)
            context = {'status': 1, 'msg': '任务已提交，请查看输出'}

        return context


class PerformTasksRollbackForm(forms.Form):
    id = forms.IntegerField(required=True)

    def rollback(self, request):
        cdata = self.cleaned_data
        id = cdata.get('id')
        obj = SqlOrdersExecTasks.objects.get(id=id)
        host = obj.host
        port = obj.port
        database = obj.database

        if obj.exec_status in ('0', '3', '4'):
            context = {'status': 2, 'msg': '请不要重复操作'}
        else:
            # 获取回滚语句
            rollback_sql = GetInceptionBackupApi(
                {'backupdbName': obj.backup_dbname, 'sequence': obj.sequence}).get_rollback_statement()
            if rollback_sql is None:
                context = {'status': 2, 'msg': '没有找到备份记录，回滚失败'}
            else:
                of_audit = InceptionSqlApi(host, port, database, rollback_sql, request.user.username)
                result = of_audit.make_sqlsha1()[1:]

                for row in result:
                    rollback_sql = row['SQL'] + ';'
                    rollback_sqlsha1 = row['sqlsha1']

                    # 将任务进度设置为：回滚中
                    obj.exec_status = 3
                    obj.rollback_sqlsha1 = rollback_sqlsha1
                    obj.save()

                    if row['sqlsha1']:
                        # 异步执行SQL任务
                        r = incep_async_tasks.delay(user=request.user.username,
                                                    id=id,
                                                    host=host,
                                                    port=port,
                                                    database=database,
                                                    sql=rollback_sql,
                                                    sqlsha1=rollback_sqlsha1,
                                                    exec_status='3')
                        task_id = r.task_id
                        # 将celery task_id写入到表
                        obj.celery_task_id = task_id
                        obj.save()
                        # 获取OSC执行进度
                        get_osc_percent.delay(task_id=task_id)

                        context = {'status': 1, 'msg': '任务已提交，请查看输出'}
                    else:
                        incep_async_tasks.delay(user=request.user.username,
                                                id=id,
                                                sql=rollback_sql,
                                                host=host,
                                                port=port,
                                                database=database,
                                                exec_status='3')

                        context = {'status': 1, 'msg': '任务已提交，请查看输出'}

        return context


class SqlOrdersTasksVersionForm(forms.Form):
    id = forms.CharField(required=False)
    tasks_version = forms.CharField(required=False)
    expire_time = forms.CharField(required=False)
    action = forms.ChoiceField(choices=(('new', 'new'), ('delete', 'delete')))

    def save(self, request):
        cdata = self.cleaned_data
        action = cdata.get('action')

        if action == 'new':
            tasks_version = cdata.get('tasks_version')
            expire_time = cdata.get('expire_time')
            if SqlOrdersTasksVersions.objects.filter(tasks_version=tasks_version).exists():
                context = {'status': 2, 'msg': '记录已存在，不能重复创建'}
            else:
                SqlOrdersTasksVersions.objects.create(tasks_version=tasks_version, expire_time=expire_time,
                                                      username=request.user.displayname)
                context = {'status': 0, 'msg': '创建成功'}
        elif action == 'delete':
            id = self.data.getlist('id')[0]
            for i in id.split(','):
                SqlOrdersTasksVersions.objects.get(pk=i).delete()
            context = {'status': 0, 'msg': '删除成功'}
        return context


class CommitOrderReplyForm(forms.Form):
    reply_id = forms.IntegerField(required=True)
    reply_contents = forms.CharField(widget=forms.Textarea, min_length=2,
                                     error_messages={'required': '回复内容不能为空', 'min_length': '回复至少输入2个字符'})

    def is_save(self, request):
        cdata = self.cleaned_data
        reply_id = cdata.get('reply_id')
        reply_contents = cdata.get('reply_contents')
        obj = SqlOrderReply.objects.create(
            reply_id=reply_id,
            user_id=request.user.uid,
            reply_contents=reply_contents)
        # 发送钉钉推送
        msg_pull = SqlOrdersMsgPull(id=obj.id, user=request.user.username, type='reply')
        msg_pull.run()
        context = {'status': 0, 'msg': ''}
        return context

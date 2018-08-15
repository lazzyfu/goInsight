# -*- coding:utf-8 -*-
# edit by fuzongfei
from datetime import datetime

from django import forms
from django.utils import timezone

from project_manager.inception.inception_api import IncepSqlCheck
from project_manager.models import IncepMakeExecTask, AuditContents
from project_manager.tasks import xiaoding_pull


class OfflineAuditForm(forms.Form):
    host = forms.CharField(required=True)
    database = forms.CharField(required=True, max_length=64)
    operate_type = forms.CharField(required=True)
    envi_desc = forms.IntegerField(required=True, label=u'环境：0：测试、1：staging、2：生产、3：线下其他环境')
    contents = forms.CharField(widget=forms.Textarea)

    def save(self, request):
        cleaned_data = super(OfflineAuditForm, self).clean()

        host, port, database = cleaned_data['database'].split(',')
        operate_type = cleaned_data.get('operate_type')
        envi_desc = cleaned_data.get('envi_desc')
        sql_content = cleaned_data['contents']

        # 实例化
        of_audit = IncepSqlCheck(sql_content, host, port, database, request.user.username)

        # 生成执行任务
        check_result = of_audit.is_check_pass()
        if check_result['status'] == 2:
            context = check_result
        else:
            # 对OSC执行的SQL生成sqlsha1
            result = of_audit.make_sqlsha1()
            taskid = datetime.now().strftime("%Y%m%d%H%M%S%f")
            # 生成执行任务记录
            for row in result:
                IncepMakeExecTask.objects.create(
                    uid=request.user.uid,
                    user=request.user.username,
                    taskid=taskid,
                    dst_host=host,
                    dst_port=port,
                    dst_database=database,
                    sql_content=row['SQL'],
                    sqlsha1=row['sqlsha1'],
                    affected_row=row['Affected_rows'],
                    type=operate_type,
                    envi_desc=envi_desc
                )
            context = {'status': 0, 'msg': '',
                       'jump_url': f'/projects/pt/perform_records/perform_details/{taskid}'}
        return context


class HookWorkOrderForm(forms.Form):
    id = forms.CharField(required=True, label=u'审核内容id')
    database = forms.CharField(required=True)
    envi_desc = forms.IntegerField(required=True)
    jump_url = forms.CharField(required=True, max_length=128, label=u'跳转到的工单页面的url')

    def save(self, request):
        cleaned_data = super(HookWorkOrderForm, self).clean()
        host, port, database = cleaned_data['database'].split(',')
        id = cleaned_data.get('id')
        envi_desc = cleaned_data['envi_desc']
        jump_url = cleaned_data.get('jump_url')
        progress_choices = '2' if envi_desc == 1 else '0'

        data = AuditContents.objects.get(pk=id)
        if AuditContents.objects.filter(title=data.title, envi_desc=envi_desc).first():
            # 如果指定的环境存在已被钩的工单，直接跳转
            context = {'status': 0, 'jump_url': jump_url}
        else:
            # 工单状态必须为已完成
            if data.progress in ('4', '6'):
                obj = AuditContents.objects.create(
                    title=data.title,
                    url=data.url,
                    tasks=data.tasks,
                    operate_type=data.operate_type,
                    host=host,
                    database=database,
                    port=port,
                    envi_desc=envi_desc,
                    progress=progress_choices,
                    remark=data.remark,
                    proposer=data.proposer,
                    operator=data.operator,
                    contents=data.contents,
                    updated_at=timezone.now()
                )

                # 更新状态为：已勾住
                AuditContents.objects.filter(pk=id).update(progress='6')

                # 发送钉钉推送
                xiaoding_pull.delay(user=request.user.username, id=obj.id, type='hook')

                # 跳转到工单记录页面
                context = {'status': 0, 'jump_url': jump_url}
            else:
                context = {'status': 2, 'msg': '当前工单进度：未完成，无法勾住'}

        return context

# -*- coding:utf-8 -*-
# edit by fuzongfei
from datetime import datetime

from django import forms
from django.utils import timezone

from ProjectManager.inception.inception_api import IncepSqlCheck
from ProjectManager.models import OnlineAuditContents
from ProjectManager.tasks import send_commit_mail, send_verify_mail


class IncepOlAuditForm(forms.Form):
    title = forms.CharField(max_length=100, required=True, label=u'标题')
    remark = forms.CharField(required=True, label=u'备注的id，以逗号分隔')
    verifier = forms.CharField(required=True, label=u'批准的leader的uid')
    operate_dba = forms.CharField(required=True, label=u'执行dba的uid')
    group_id = forms.CharField(required=True, label=u'项目组id')
    host = forms.CharField(required=True, label=u'目标数据库主机')
    database = forms.CharField(required=True, max_length=64)
    op_action = forms.CharField(required=True, label=u'操作类型，是DDL还是DML')
    sql_content = forms.CharField(widget=forms.Textarea)

    def save(self, request):
        cleaned_data = super(IncepOlAuditForm, self).clean()
        title = cleaned_data.get('title') + '__[' + datetime.now().strftime("%Y%m%d%H%M%S") + ']'
        remark = cleaned_data.get('remark')
        verifier = cleaned_data.get('verifier')
        operate_dba = cleaned_data.get('operate_dba')
        group_id = cleaned_data.get('group_id')
        email_cc = self.data.get('email_cc_id')
        host = cleaned_data.get('host')
        database = cleaned_data.get('database')
        op_action = cleaned_data.get('op_action')
        sql_content = cleaned_data.get('sql_content')

        result = IncepSqlCheck(sql_content, host, database, request.user.username).is_check_pass()
        if result.get('status') == 2:
            context = result
        else:
            OnlineAuditContents.objects.create(
                title=title,
                op_action='数据修改' if op_action == 'op_data' else '表结构变更',
                type='DML' if op_action == 'op_data' else 'DDL',
                dst_host=host,
                dst_database=database,
                group_id=group_id,
                remark=remark,
                proposer=request.user.username,
                operate_dba=operate_dba,
                verifier=verifier,
                email_cc=email_cc,
                contents=sql_content
            )

            # 发送通知邮件
            latest_id = OnlineAuditContents.objects.latest('id').id
            send_commit_mail.delay(latest_id=latest_id)
            context = {'status': 0, 'msg': '', 'jump_url': '/projects/ol/incep_ol_records/'}
        return context


class IncepOlApproveForm(forms.Form):
    id = forms.IntegerField(required=True)
    status = forms.CharField(max_length=10, required=True)
    addition_info = forms.CharField(required=False)


class IncepOlReplyForm(forms.Form):
    reply_id = forms.IntegerField(required=True)
    reply_contents = forms.CharField(widget=forms.Textarea, min_length=5)

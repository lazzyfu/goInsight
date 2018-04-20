# -*- coding:utf-8 -*-
# edit by fuzongfei
from datetime import datetime

from django import forms
from project_manager.tasks import send_commit_mail

from project_manager.models import audit_type_choice, operate_type_choice, AuditContents, OlDataExportDetail


class OlDataExportForm(forms.Form):
    title = forms.CharField(max_length=100, required=True, label=u'标题')
    verifier = forms.CharField(required=True, label=u'批准人')
    operator = forms.CharField(required=True, label=u'执行人')
    group_id = forms.CharField(required=True, label=u'项目组id')
    host = forms.CharField(required=True, label=u'数据库主机')
    database = forms.CharField(required=True, max_length=64, label=u'数据库')
    audit_type = forms.ChoiceField(choices=audit_type_choice, label=u'审核类型')
    file_format = forms.CharField(required=True, max_length=32, label=u'文件格式，如：xlsx，txt等')
    file_coding = forms.CharField(required=True, max_length=32, label=u'文件格式，如：utf-8, gbk等')
    contents = forms.CharField(widget=forms.Textarea, label=u'审核内容')

    def is_save(self, request):
        cleaned_data = super(OlDataExportForm, self).clean()
        title = cleaned_data.get('title') + '_[' + datetime.now().strftime("%Y%m%d%H%M%S") + ']'
        verifier = cleaned_data.get('verifier')
        operator = cleaned_data.get('operator')
        email_cc = self.data.get('email_cc_id')
        group_id = cleaned_data.get('group_id')
        host = cleaned_data.get('host')
        database = cleaned_data.get('database')
        audit_type = cleaned_data.get('audit_type')
        file_format = cleaned_data.get('file_format')
        file_coding = cleaned_data.get('file_coding')
        contents = cleaned_data.get('contents')

        AuditContents.objects.create(
            title=title,
            audit_type=audit_type,
            host=host,
            database=database,
            group_id=group_id,
            proposer=request.user.username,
            operator=operator,
            verifier=verifier,
            email_cc=email_cc
        )

        # 向子表插入关联数据
        OlDataExportDetail.objects.create(ol=AuditContents.objects.latest('id'),
                                          file_format=file_format,
                                          file_coding=file_coding,
                                          contents=contents)

        latest_id = AuditContents.objects.latest('id').id
        send_commit_mail.delay(latest_id=latest_id)
        context = {'status': 0, 'jump_url': '/projects/ol/incep_ol_records/'}
        return context

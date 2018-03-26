# -*- coding:utf-8 -*-
# edit by fuzongfei

from django import forms


class IncepSqlCheckForm(forms.Form):
    host = forms.CharField(required=True)
    database = forms.CharField(required=True, max_length=64)
    op_action = forms.CharField(required=True)
    op_type = forms.CharField(required=True)
    sql_content = forms.CharField(widget=forms.Textarea)


class InceptionSqlCheckForm(forms.Form):
    host = forms.CharField(required=True)
    database = forms.CharField(required=True, max_length=64)
    op_action = forms.CharField(required=True)
    op_type = forms.CharField(required=True)
    group_id = forms.IntegerField(required=True)
    sql_content = forms.CharField(widget=forms.Textarea)

class OnlineAuditCommitForm(forms.Form):
    title = forms.CharField(max_length=100, required=True, label=u'标题')
    remark = forms.CharField(required=True, label=u'备注的id，以逗号分隔')
    verifier = forms.CharField(required=True, label=u'批准的leader的uid')
    operate_dba = forms.CharField(required=True, label=u'执行dba的uid')
    group_id = forms.CharField(required=True, label=u'项目组id')
    host = forms.CharField(required=True)
    database = forms.CharField(required=True, max_length=64)
    op_action = forms.CharField(required=True)
    sql_content = forms.CharField(widget=forms.Textarea)

class VerifyCommitForm(forms.Form):
    id = forms.IntegerField(required=True)
    status = forms.CharField(max_length=10, required=True)
    addition_info = forms.CharField(required=False)

class ReplyContentForm(forms.Form):
    reply_id = forms.IntegerField(required=True)
    reply_contents = forms.CharField(widget=forms.Textarea, min_length=5)
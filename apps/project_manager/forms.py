# -*- coding:utf-8 -*-
# edit by fuzongfei

from django import forms
from .models import operate_type_choice


class IncepSqlCheckForm(forms.Form):
    host = forms.CharField(required=True)
    database = forms.CharField(required=True, max_length=64)
    op_action = forms.CharField(required=True)
    op_type = forms.CharField(required=True)
    sql_content = forms.CharField(widget=forms.Textarea)


class SyntaxCheckForm(forms.Form):
    database = forms.CharField(required=True, max_length=64)
    operate_type = forms.ChoiceField(choices=operate_type_choice, label=u'操作类型，是DDL还是DML')
    contents = forms.CharField(widget=forms.Textarea)

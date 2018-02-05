# -*- coding:utf-8 -*-
# edit by fuzongfei

from django import forms

class InceptionSqlOperateForm(forms.Form):
    host = forms.CharField(required=True)
    database = forms.CharField(required=True, max_length=64)
    op_type = forms.CharField(required=True)
    sql_content = forms.CharField(widget=forms.Textarea)
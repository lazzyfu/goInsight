# -*- coding:utf-8 -*-
# edit by fuzongfei

from django import forms


class PrivModifyForm(forms.Form):
    db_host = forms.CharField(max_length=128, min_length=3, required=True)
    user = forms.CharField(max_length=30, min_length=1, required=True)
    host = forms.CharField(max_length=30, min_length=1, required=True)
    action = forms.ChoiceField(choices=(('modify_privileges', u'更改权限'), ('new_host', u'新建主机')))
    schema = forms.CharField(max_length=64, min_length=1, required=True)
    privileges = forms.CharField(max_length=1024, min_length=3, required=True)

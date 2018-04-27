# -*- coding:utf-8 -*-
# edit by fuzongfei
from django import forms


class LoginForm(forms.Form):
    username = forms.CharField(required=True, max_length=30,
                               widget=forms.TextInput(attrs={'class': 'form-control', 'placeholder': u'用户名'}))
    password = forms.CharField(required=True, max_length=30,
                               widget=forms.PasswordInput(attrs={'class': 'form-control', 'placeholder': u'密码'}))


class ChangePasswordForm(forms.Form):
    old_password = forms.CharField(max_length=30, required=True)
    new_password = forms.CharField(max_length=30, min_length=7, required=True)
    verify_password = forms.CharField(max_length=30, min_length=7, required=True)

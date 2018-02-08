# -*- coding:utf-8 -*-
# edit by fuzongfei

from django import forms
from django.contrib.auth import authenticate


class LoginForm(forms.Form):
    username = forms.CharField(required=True, max_length=30)
    password = forms.CharField(required=True, max_length=30)

    def is_verify(self):
        cleaned_data = super(LoginForm, self).clean()
        username = cleaned_data.get('username')
        password = cleaned_data.get('password')

        user = authenticate(username=username, password=password)
        return user


class ChangePasswordForm(forms.Form):
    old_password = forms.CharField(max_length=30, required=True)
    new_password = forms.CharField(max_length=30, min_length=7, required=True)
    verify_password = forms.CharField(max_length=30, min_length=7, required=True)

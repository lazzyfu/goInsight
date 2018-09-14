# -*- coding:utf-8 -*-
# edit by fuzongfei
from django import forms
from django.contrib.auth import authenticate, login
from django.contrib.auth.hashers import make_password
from django.forms import ModelForm

from users.models import UserAccounts, UserRoles, RolePermission
from users.utils import check_ldap_connection


class LoginForm(forms.Form):
    username = forms.CharField(required=True, max_length=30)
    password = forms.CharField(required=True, max_length=30, min_length=7)

    def authentication(self, request):
        cdata = self.cleaned_data
        username = cdata.get('username')
        password = cdata.get('password')

        status, msg = check_ldap_connection()
        if status:
            try:
                user = authenticate(username=username, password=password)
                obj = UserAccounts.objects.get(username=username)
                if not obj.check_password(password):
                    result = {'status': False, 'msg': f'用户{username}密码错误'}
                if not obj.is_active:
                    result = {'status': False, 'msg': f'用户{username}被禁用，请联系管理员'}
                else:
                    if user is not None:
                        login(request, user)
                        result = {'status': True}
            except UserAccounts.DoesNotExist:
                result = {'status': False, 'msg': '用户不存在，请联系管理员'}
        else:
            result = {'status': False, 'msg': msg}
        return result


class ChangePasswordForm(forms.Form):
    old_password = forms.CharField(max_length=30, required=True)
    new_password = forms.CharField(max_length=30, min_length=7, required=True)
    verify_password = forms.CharField(max_length=30, min_length=7, required=True)

    def change_pass(self, request):
        cdata = self.cleaned_data
        old_password = cdata['old_password']
        new_password = cdata['new_password']
        verify_password = cdata['verify_password']

        user = UserAccounts.objects.get(uid=request.user.uid)
        if new_password == verify_password:
            if user.check_password(old_password):
                if old_password != new_password:
                    user.password = make_password(new_password)
                    user.save()
                    context = {'status': 0, 'msg': '密码修改成功'}
                else:
                    context = {'status': 2, 'msg': '新旧密码一致，请重新输入'}
            else:
                context = {'status': 2, 'msg': '旧密码错误，请重新输入'}
        else:
            context = {'status': 2, 'msg': '密码不匹配，请重新输入'}
        return context


class ChangeMobileForm(forms.Form):
    mobile = forms.RegexField(regex=r'^\+?1?\d{9,15}$', required=True, min_length=11, max_length=11)

    def change_mobile(self, request):
        cdata = self.cleaned_data
        mobile = cdata['mobile']
        UserAccounts.objects.filter(uid=request.user.uid).update(mobile=mobile)
        context = {'status': 0, 'msg': '手机号修改成功', 'data': mobile}
        return context


class UserSaveForm(forms.Form):
    uid = forms.IntegerField()
    username = forms.CharField(max_length=30, min_length=2, required=True)
    password = forms.CharField(max_length=128, min_length=7, required=True)
    email = forms.EmailField()
    displayname = forms.CharField(max_length=30, min_length=1)
    mobile = forms.RegexField(regex=r'^\+?1?\d{9,15}$', min_length=11, max_length=11)
    is_active = forms.ChoiceField(choices=(('0', 'disable'), ('1', 'active')))
    user_role = forms.IntegerField(required=True)

    def save(self):
        cdata = self.cleaned_data
        uid = cdata.get('uid')
        password = cdata.pop('password')
        rid = cdata.pop('user_role')
        # 保存用户的基本信息
        UserAccounts.objects.update_or_create(uid=uid, defaults=cdata)
        if len(password) < 50:
            UserAccounts.objects.filter(uid=uid).update(password=make_password(password))

        # 设置用户角色
        new_role = UserRoles.objects.get(rid=rid)
        user = UserAccounts.objects.get(uid=uid)
        # 判断新旧用户角色是否相等
        if user.user_role != new_role.role_name:
            # 判断用户角色是否存在
            if UserRoles.objects.filter(user=uid).exists():
                # 存在先移除
                old_role = UserRoles.objects.get(user=uid)
                old_role.user.remove(user)
            # 添加新的角色
            new_role.user.add(user)
        return 'ok'


class UserDeleteForm(forms.Form):
    uid = forms.CharField()

    def delete(self):
        cdata = self.cleaned_data
        uid = cdata.get('uid')
        for i in uid.split(','):
            UserAccounts.objects.get(uid=i).delete()
        context = {'status': 0, 'msg': '用户删除成功'}
        return context


class RolesChangeForm(forms.Form):
    id = forms.CharField(required=True)
    role_name = forms.CharField()
    # field预留字段，可用于判断修改的字段
    field = forms.ChoiceField(choices=(('id', 'id'),))
    oldvalue = forms.CharField()

    def change(self):
        cdata = self.cleaned_data
        role_name = cdata['role_name']
        oldvalue = cdata['oldvalue']
        id = cdata['id']

        role = UserRoles.objects.get(role_name=role_name)
        if oldvalue:
            for i in oldvalue.split(','):
                permission = RolePermission.objects.get(id=i)
                permission.role.remove(role)
        if id:
            for j in id.split(','):
                permission = RolePermission.objects.get(id=j)
                permission.role.add(role)

        context = {'status': 0, 'msg': f"角色[{role_name}]权限修改成功"}
        return context

# -*- coding:utf-8 -*-
# edit by fuzongfei
from django.contrib.auth import authenticate, login
from django.contrib.auth.hashers import make_password
from django.db.models import F
from rest_framework import serializers

from users.models import UserAccounts, RolePermissions, UserRoles
from users.utils.validatorsCheck import password_validate, audit_permission_validator


class LoginSerializer(serializers.Serializer):
    username = serializers.CharField(required=True, max_length=30, min_length=2, label='用户名',
                                     error_messages={'required': '用户名不能为空',
                                                     'min_length': '用户名至少7个字符',
                                                     'max_lenght': '用户名最大长度为30个字符'
                                                     }
                                     )
    password = serializers.CharField(required=True, max_length=30, min_length=7, label='密码',
                                     error_messages={'required': '密码不能为空',
                                                     'min_length': '密码至少7个字符',
                                                     'max_lenght': '密码最大长度为30个字符'
                                                     }
                                     )
    verifycode = serializers.CharField(required=True, label='验证码', error_messages={'required': '验证码不能为空'})

    def authentication(self, request):
        sdata = self.validated_data
        username = sdata.get('username')
        password = sdata.get('password')
        try:
            user = authenticate(username=username, password=password)
            if user is not None:
                obj = UserAccounts.objects.get(username=username)
                if not obj.is_active:
                    msg = f'用户{self.username}被禁用，请联系管理员'
                    return False, msg
                else:
                    login(request, user)
                    return True, None
            else:
                msg = '用户名或密码错误'
                return False, msg
        except UserAccounts.DoesNotExist:
            msg = f'当前用户:{username}不存在，请联系系统管理员'
            return False, msg


class ChangePasswordSerializer(serializers.Serializer):
    old_password = serializers.CharField(required=True, error_messages={'required': '原密码不能为空'})
    new_password = serializers.CharField(required=True, min_length=7, max_length=30,
                                         error_messages={
                                             'required': '新密码不能为空',
                                             'min_length': '新密码至少7个字符',
                                             'max_lenght': '新密码最大长度为30个字符'
                                         }, validators=[password_validate]
                                         )
    verify_password = serializers.CharField(required=True, min_length=7, max_length=30,
                                            error_messages={
                                                'required': '确认密码不能为空',
                                                'min_length': '确认密码至少7个字符',
                                                'max_lenght': '确认密码最大长度为30个字符'
                                            }, validators=[password_validate]
                                            )

    def change(self, request):
        sdata = self.validated_data
        old_password = sdata['old_password']
        new_password = sdata['new_password']
        verify_password = sdata['verify_password']

        user = UserAccounts.objects.get(uid=request.user.uid)
        if new_password == verify_password:
            if user.check_password(old_password):
                if old_password != new_password:
                    user.password = make_password(new_password)
                    user.save()
                    msg = '密码修改成功'
                    return True, msg
                else:
                    msg = '新旧密码一致，请重新输入'
                    return False, msg
            else:
                msg = '旧密码错误，请重新输入'
                return False, msg
        else:
            msg = '密码不匹配，请重新输入'
            return False, msg


class ChangeMobileSerializer(serializers.Serializer):
    mobile = serializers.RegexField(regex=r'^\+?1?\d{9,15}$', required=True, min_length=11, max_length=11,
                                    error_messages={
                                        'required': '手机号不能为空',
                                        'min_length': '手机号最小长度为11',
                                        'max_length': '手机号最大长度为11'
                                    })

    def change(self, request):
        sdata = self.validated_data
        mobile = sdata['mobile']
        UserAccounts.objects.filter(uid=request.user.uid).update(mobile=mobile)
        msg = '手机号修改成功'
        return True, msg


class GetAuditorSerializer(serializers.Serializer):
    permission = serializers.CharField(required=True, validators=[audit_permission_validator],
                                       error_messages={'required': '验证权限字段不能为空'})

    def query(self):
        sdata = self.validated_data
        permission = sdata.get('permission')
        roles = RolePermissions.objects.filter(permission_name=permission).values_list(
            'role__role_name', flat=True)
        queryset = UserRoles.objects.filter(role_name__in=roles).filter(
            user__username__isnull=False
        ).values(
            username=F('user__username'),
            displayname=F('user__displayname')
        ).order_by('username')
        return queryset

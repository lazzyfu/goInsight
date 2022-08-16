# -*- coding:utf-8 -*-
# edit by xff
from uuid import uuid4

from django.contrib.auth import authenticate
from django.contrib.auth.hashers import make_password
from rest_framework import serializers

from users import models
from users.utils import validate_password


class LoginSerializer(serializers.Serializer):
    username = serializers.CharField(required=True, max_length=30, min_length=3, label='用户名',
                                     error_messages={'required': '用户名不能为空',
                                                     'min_length': '用户名至少3个字符',
                                                     'max_lenght': '用户名最大长度为30个字符',
                                                     'blank': '输入用户名不能为空，至少3个字符，最多30个字符'
                                                     }
                                     )
    password = serializers.CharField(required=True, max_length=30, min_length=7, label='密码',
                                     error_messages={'required': '密码不能为空',
                                                     'min_length': '密码至少7个字符',
                                                     'max_lenght': '密码最大长度为30个字符',
                                                     'blank': '输入密码不能为空，至少7个字符，最多30个字符'
                                                     }
                                     )

    def login(self):
        username = self.validated_data.get('username')
        password = self.validated_data.get('password')

        user = authenticate(username=username, password=password)
        if not user:
            return False, '用户名或密码错误'

        try:
            user = models.UserAccounts.objects.get(username=username)
        except models.UserAccounts.DoesNotExist:
            return False, '用户不存在'

        if not user.is_active:
            return False, '用户被禁用，请联系管理员'

        token = user.token
        return True, token


class UserInfoSerializer(serializers.ModelSerializer):
    class Meta:
        model = models.UserAccounts
        fields = ['uid', 'displayname', 'mobile', 'avatar_file', 'username', 'is_superuser', 'email']


class UsersListSerializer(serializers.ModelSerializer):
    class Meta:
        model = models.UserAccounts
        fields = ['username', 'displayname', 'email']


class UpdateUserInfoSerializer(serializers.ModelSerializer):
    class Meta:
        model = models.UserAccounts
        fields = ['displayname', 'email', 'mobile']

    def update(self, instance, validated_data):
        # 更新instance的记录
        for k, v in validated_data.items():
            setattr(instance, k, v)
        instance.save()
        return instance


class ChangePasswordSerializer(serializers.Serializer):
    current_password = serializers.CharField(required=True,
                                             error_messages={'required': '原密码不能为空'})
    new_password = serializers.CharField(required=True, min_length=7, max_length=30,
                                         error_messages={
                                             'required': '新密码不能为空',
                                             'min_length': '新密码至少7个字符',
                                             'max_lenght': '新密码最大长度为30个字符'
                                         }, validators=[validate_password]
                                         )
    verify_password = serializers.CharField(required=True, min_length=7, max_length=30,
                                            error_messages={
                                                'required': '确认密码不能为空',
                                                'min_length': '确认密码至少7个字符',
                                                'max_lenght': '确认密码最大长度为30个字符'
                                            }, validators=[validate_password]
                                            )

    def validate(self, attrs):
        if attrs['new_password'] != attrs['verify_password']:
            raise serializers.ValidationError('两次密码不匹配，请重新输入')

        if attrs['current_password'] == attrs['new_password']:
            raise serializers.ValidationError('新旧密码相同，请重新输入')

        return super(ChangePasswordSerializer, self).validate(attrs)

    def change(self, request):
        vdata = self.validated_data
        current_password = vdata['current_password']
        new_password = vdata['new_password']

        user = models.UserAccounts.objects.get(uid=request.user.uid)
        if not user.check_password(current_password):
            return False, '原密码错误，请重新输入'

        user.password = make_password(new_password)
        user.user_secret = uuid4()  # 生成新的secret，让token过期
        user.save()
        return True, '密码修改成功'


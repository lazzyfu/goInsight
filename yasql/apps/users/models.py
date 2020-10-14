# -*- coding:utf-8 -*-
# edit by fuzongfei
from uuid import uuid4

from django.contrib.auth.models import AbstractUser
from django.db import models
from rest_framework_jwt.settings import api_settings


class UserAccounts(AbstractUser):
    uid = models.BigAutoField(primary_key=True, verbose_name=u'uid')
    user_secret = models.UUIDField(default=uuid4(), verbose_name=u'用户JWT秘钥')
    is_active = models.BooleanField(default=True, verbose_name=u'是否激活')
    displayname = models.CharField(max_length=128, default='', verbose_name=u'别名')
    mobile = models.CharField(max_length=11, null=True, blank=True, verbose_name=u'用户手机号')
    department = models.CharField(max_length=128, null=False, blank=True, default='', verbose_name=u'部门')
    avatar_file = models.ImageField(upload_to='img/%Y/%m/%d/',
                                    default=u'img/avatar1.png',
                                    verbose_name=u'用户头像')

    def __str__(self):
        return self.username

    def user_role(self):
        try:
            role_name = UserRoles.objects.get(user=self.uid).role_name
        except UserRoles.DoesNotExist:
            role_name = None
        return role_name

    user_role.short_description = '角色'

    @property
    def token(self):
        # 生成token
        jwt_payload_handler = api_settings.JWT_PAYLOAD_HANDLER
        jwt_encode_handler = api_settings.JWT_ENCODE_HANDLER
        obj = UserAccounts.objects.get(uid=self.uid)
        payload = jwt_payload_handler(obj)
        token = jwt_encode_handler(payload)
        return token

    class Meta:
        verbose_name = u'用户表'
        verbose_name_plural = verbose_name

        default_permissions = ()
        app_label = 'users'
        db_table = 'yasql_user_accounts'


class UserRoles(models.Model):
    """用户角色表"""
    rid = models.AutoField(primary_key=True, verbose_name=u'主键')
    role_name = models.CharField(max_length=30, null=False, verbose_name=u'角色名')
    user = models.ManyToManyField('UserAccounts')
    created_at = models.DateTimeField(auto_now_add=True, verbose_name=u'创建时间')
    updated_at = models.DateTimeField(auto_now=True, verbose_name=u'更新时间')

    def __str__(self):
        return self.role_name

    def permission(self):
        permission_desc = RolePermissions.objects.filter(
            role__rid=self.rid).values_list('permission_desc', flat=True)
        return ', '.join(permission_desc)

    permission.short_description = '权限'

    class Meta:
        verbose_name = u'角色表'
        verbose_name_plural = verbose_name

        default_permissions = ()
        app_label = 'users'
        db_table = 'yasql_user_roles'


class RolePermissions(models.Model):
    """角色权限表"""
    id = models.AutoField(primary_key=True, verbose_name=u'主键')
    permission_name = models.CharField(max_length=30, null=False, verbose_name=u'权限名')
    permission_desc = models.CharField(max_length=30, default='', null=False, verbose_name=u'权限描述')
    role = models.ManyToManyField('UserRoles')
    created_at = models.DateTimeField(auto_now_add=True, verbose_name=u'创建时间')
    updated_at = models.DateTimeField(auto_now=True, verbose_name=u'更新时间')

    def __str__(self):
        return self.permission_desc

    class Meta:
        verbose_name = u'权限表'
        verbose_name_plural = verbose_name

        default_permissions = ()
        app_label = 'users'
        db_table = 'yasql_role_permissions'

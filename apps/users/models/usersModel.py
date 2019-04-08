# -*- coding:utf-8 -*-
# edit by fuzongfei

from django.contrib.auth.models import AbstractUser
from django.db import models

from .rolesModel import UserRoles


class UserAccounts(AbstractUser):
    uid = models.BigAutoField(primary_key=True, verbose_name=u'uid')
    is_active = models.BooleanField(default=True, verbose_name=u'是否激活')
    displayname = models.CharField(max_length=128, default='', verbose_name=u'别名')
    mobile = models.CharField(max_length=11, null=True, blank=True, verbose_name=u'用户手机号')
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

    class Meta:
        verbose_name = u'用户表'
        verbose_name_plural = verbose_name

        default_permissions = ()
        app_label = 'users'
        db_table = 'auditsql_useraccounts'

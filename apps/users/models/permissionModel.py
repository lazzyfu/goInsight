# -*- coding:utf-8 -*-
# edit by fuzongfei
from django.db import models


class RolePermissions(models.Model):
    """角色权限表"""
    id = models.AutoField(primary_key=True, verbose_name=u'主键')
    permission_name = models.CharField(max_length=30, null=False, verbose_name=u'权限名')
    permission_desc = models.CharField(max_length=30, default='', null=False, verbose_name=u'权限描述')
    role = models.ManyToManyField('users.UserRoles')
    created_at = models.DateTimeField(auto_now_add=True, verbose_name=u'创建时间')
    updated_at = models.DateTimeField(auto_now=True, verbose_name=u'更新时间')

    def __str__(self):
        return self.permission_desc

    class Meta:
        verbose_name = u'权限表'
        verbose_name_plural = verbose_name

        default_permissions = ()
        app_label = 'users'
        db_table = 'auditsql_role_permissions'

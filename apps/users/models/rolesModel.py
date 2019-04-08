# -*- coding:utf-8 -*-
# edit by fuzongfei

from django.db import models

from users.models import permissionModel


class UserRoles(models.Model):
    """用户角色表"""
    rid = models.AutoField(primary_key=True, verbose_name=u'主键')
    role_name = models.CharField(max_length=30, null=False, verbose_name=u'角色名')
    user = models.ManyToManyField('users.UserAccounts')
    created_at = models.DateTimeField(auto_now_add=True, verbose_name=u'创建时间')
    updated_at = models.DateTimeField(auto_now=True, verbose_name=u'更新时间')

    def __str__(self):
        return self.role_name

    def permission(self):
        permission_desc = permissionModel.RolePermissions.objects.filter(
            role__rid=self.rid).values_list('permission_desc', flat=True)
        return ', '.join(permission_desc)

    permission.short_description = '权限'

    class Meta:
        verbose_name = u'角色表'
        verbose_name_plural = verbose_name

        default_permissions = ()
        app_label = 'users'
        db_table = 'auditsql_user_roles'

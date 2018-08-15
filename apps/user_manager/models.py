from django.contrib.auth.models import AbstractUser
from django.db import models
from django.db.models import F


# Create your models here.


class UserAccount(AbstractUser):
    """继承系统用户表"""
    uid = models.AutoField(primary_key=True, verbose_name=u'用户uid')
    avatar_file = models.ImageField(upload_to='img/%Y/%m/%d/', default=u'img/avatar1.png',
                                    verbose_name=u'用户头像')
    displayname = models.CharField(max_length=128, default='', verbose_name=u'用户显示名')
    is_active = models.BooleanField(default=False, verbose_name=u'是否激活')
    mobile = models.CharField(max_length=11, null=False, default='', verbose_name=u'用户手机号')

    def __str__(self):
        return self.username

    def fullname(self):
        return '_'.join((self.username, self.displayname))

    def user_role(self):
        # 返回用户角色名
        return RolesDetail.objects.annotate(role_name=F('role__role_name')).get(user__uid=self.uid).role_name

    def user_schema(self):
        # 返回用户授权的schema
        from mstats.models import MysqlSchemaGrant
        result = MysqlSchemaGrant.objects.annotate(schemas=F('schema__schema')).filter(
            user__uid=self.uid).values_list('schemas', flat=True)
        return ','.join(list(result))

    def user_shell(self):
        # 返回用户授权的shell
        from mstats.models import WebShellGrant
        result = WebShellGrant.objects.annotate(shell_name=F('shell__comment')).filter(
            user__uid=self.uid).values_list('shell_name', flat=True)
        return ','.join(list(result))

    class Meta:
        verbose_name = u'账户'
        verbose_name_plural = verbose_name

        default_permissions = ()
        db_table = 'auditsql_useraccount'


class Roles(models.Model):
    """用户角色表"""
    role_id = models.AutoField(primary_key=True, verbose_name=u'主键')
    role_name = models.CharField(max_length=30, null=False, verbose_name=u'角色名')
    created_at = models.DateTimeField(auto_now_add=True, verbose_name=u'创建时间')
    updated_at = models.DateTimeField(auto_now=True, verbose_name=u'更新时间')

    def __str__(self):
        return self.role_name

    def permission(self):
        permission_desc = PermissionDetail.objects.annotate(permission_desc=F('permission__permission_desc')).filter(
            role__role_id=self.role_id).values_list('permission_desc', flat=True)
        return ', '.join(permission_desc)

    class Meta:
        verbose_name = u'角色'
        verbose_name_plural = verbose_name

        default_permissions = ()
        db_table = 'auditsql_roles'


class RolesDetail(models.Model):
    """
    用户角色详情表
    对应用户角色分组
    """
    id = models.AutoField(primary_key=True, verbose_name=u'主键')
    user = models.ForeignKey(UserAccount, on_delete=models.CASCADE, null=False)
    role = models.ForeignKey(Roles, on_delete=models.CASCADE, null=False)
    created_at = models.DateTimeField(auto_now_add=True, verbose_name=u'创建时间')
    updated_at = models.DateTimeField(auto_now=True, verbose_name=u'更新时间')

    class Meta:
        verbose_name = u'用户角色分组详情表'
        verbose_name_plural = verbose_name

        default_permissions = ()
        db_table = 'auditsql_roles_detail'

        unique_together = ('user',)


class RolePermission(models.Model):
    """权限表"""
    id = models.AutoField(primary_key=True, verbose_name=u'主键')
    permission_name = models.CharField(max_length=30, null=False, verbose_name=u'权限名')
    permission_desc = models.CharField(max_length=30, default='', null=False, verbose_name=u'权限描述')
    created_at = models.DateTimeField(auto_now_add=True, verbose_name=u'创建时间')
    updated_at = models.DateTimeField(auto_now=True, verbose_name=u'更新时间')

    def __str__(self):
        return self.permission_desc

    class Meta:
        verbose_name = u'内置权限'
        verbose_name_plural = verbose_name

        default_permissions = ()
        db_table = 'auditsql_permission'


class PermissionDetail(models.Model):
    """角色对应的权限，多对多关系"""
    id = models.AutoField(primary_key=True, verbose_name=u'主键')
    permission = models.ForeignKey(RolePermission, on_delete=models.CASCADE, null=False)
    role = models.ForeignKey(Roles, on_delete=models.CASCADE, null=False)
    created_at = models.DateTimeField(auto_now_add=True, verbose_name=u'创建时间')
    updated_at = models.DateTimeField(auto_now=True, verbose_name=u'更新时间')

    class Meta:
        verbose_name = u'权限详情表'
        verbose_name_plural = verbose_name

        default_permissions = ()
        db_table = 'auditsql_permission_detail'


class SystemMsg(models.Model):
    """系统消息推送表"""
    id = models.AutoField(primary_key=True, verbose_name=u'主键')
    title = models.CharField(max_length=64, null=False, default='', verbose_name=u'标题')
    content = models.TextField(default='', verbose_name=u'内容')
    created_at = models.DateTimeField(auto_now_add=True, verbose_name=u'创建时间')

    def __str__(self):
        return self.title

    class Meta:
        verbose_name = u'系统消息推送表'
        verbose_name_plural = verbose_name

        default_permissions = ()
        db_table = 'auditsql_system_msg'


class SystemMsgDetails(models.Model):
    """系统消息推送详情表，标识用户是否已读"""
    id = models.AutoField(primary_key=True, verbose_name=u'主键')
    user = models.ForeignKey(UserAccount, on_delete=models.CASCADE, null=False)
    msg = models.ForeignKey(SystemMsg, on_delete=models.CASCADE, null=False)
    created_at = models.DateTimeField(auto_now_add=True, verbose_name=u'创建时间')

    class Meta:
        verbose_name = u'系统消息推送详情表'
        verbose_name_plural = verbose_name

        default_permissions = ()
        db_table = 'auditsql_system_msg_details'

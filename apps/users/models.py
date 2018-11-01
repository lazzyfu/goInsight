# Create your models here.
from django.contrib.auth.models import AbstractUser
from django.db import models


# Create your models here.
from django.db.models import F


class UserAccounts(AbstractUser):
    uid = models.BigAutoField(primary_key=True, verbose_name=u'用户uid')
    is_active = models.BooleanField(default=True, verbose_name=u'激活')
    displayname = models.CharField(max_length=128, default='', verbose_name=u'别名')
    mobile = models.CharField(max_length=11, null=True, blank=True, verbose_name=u'手机号')
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

    user_role.short_description = '用户角色'

    def user_schema(self):
        # 返回用户授权的schema
        from sqlquery.models import MysqlSchemasGrant
        result = MysqlSchemasGrant.objects.annotate(schemas=F('schema__schema'), host=F('schema__host'),
                                                    port=F('schema__port')).filter(
            user__uid=self.uid).values_list('schema__comment', 'schemas')
        return ','.join(['_'.join(map(str, i)) for i in list(result)])

    user_schema.short_description = '授权库'

    class Meta:
        verbose_name = u'用户'
        verbose_name_plural = verbose_name

        default_permissions = ()
        db_table = 'sqlaudit_user_accounts'


class UserRoles(models.Model):
    """用户角色表"""
    rid = models.AutoField(primary_key=True, verbose_name=u'主键')
    role_name = models.CharField(max_length=30, null=False, verbose_name=u'角色名')
    user = models.ManyToManyField(UserAccounts)
    created_at = models.DateTimeField(auto_now_add=True, verbose_name=u'创建时间')
    updated_at = models.DateTimeField(auto_now=True, verbose_name=u'更新时间')

    def __str__(self):
        return self.role_name

    def permission(self):
        permission_desc = RolePermission.objects.filter(
            role__rid=self.rid).values_list('permission_desc', flat=True)
        return ', '.join(permission_desc)

    class Meta:
        verbose_name = u'角色'
        verbose_name_plural = verbose_name

        default_permissions = ()
        db_table = 'sqlaudit_user_roles'


class RolePermission(models.Model):
    """角色权限表"""
    id = models.AutoField(primary_key=True, verbose_name=u'主键')
    permission_name = models.CharField(max_length=30, null=False, verbose_name=u'权限名')
    permission_desc = models.CharField(max_length=30, default='', null=False, verbose_name=u'权限描述')
    role = models.ManyToManyField(UserRoles)
    created_at = models.DateTimeField(auto_now_add=True, verbose_name=u'创建时间')
    updated_at = models.DateTimeField(auto_now=True, verbose_name=u'更新时间')

    def __str__(self):
        return self.permission_desc

    class Meta:
        verbose_name = u'权限'
        verbose_name_plural = verbose_name

        default_permissions = ()
        db_table = 'sqlaudit_permissions'
